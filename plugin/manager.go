package plugin

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	things "github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/addon/manager"
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/bus/topic"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/errors"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/galenliu/gateway/pkg/util"
	"github.com/robfig/cron"

	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type managerStore interface {
	AddonsStore
	GetSetting(key string) (string, error)
	SetSetting(key, value string) error
}

type ThingsContainer interface {
	Subscribe(topic topic.Topic, fn any) error
	Unsubscribe(topic topic.Topic, fn any)
	GetThings() []*things.Thing
}

type Config struct {
	AddonsDir       string
	AttachAddonsDir string
	IPCPort         string
	RPCPort         string
	UserProfile     *messages.PluginRegisterResponseJsonDataUserProfile
}

type Manager struct {
	*manager.Manager
	*bus.EventBus
	config         Config
	configPath     string
	pluginServer   *PluginsServer
	preferences    messages.PluginRegisterResponseJsonDataPreferences
	outlets        sync.Map
	installAddons  sync.Map
	extensions     sync.Map
	addonsLoaded   bool
	isPairing      bool
	running        bool
	pairTask       chan bool
	locker         *sync.Mutex
	storage        managerStore
	ctx            context.Context
	thingContainer ThingsContainer
	deferredRemove sync.Map
}

func NewAddonsManager(ctx context.Context, conf Config, s managerStore) *Manager {
	am := &Manager{}
	am.Manager = manager.NewManager()
	am.config = conf
	am.ctx = ctx
	am.addonsLoaded = false
	am.isPairing = false
	am.running = false
	am.EventBus = bus.NewBus()
	am.storage = s
	am.locker = new(sync.Mutex)
	am.UpdatePreferences()
	am.loadAddons()
	return am
}

func (m *Manager) SetThingsContainer(thingContainer ThingsContainer) {
	m.thingContainer = thingContainer
}

func (m *Manager) RequestAction(ctx context.Context, thingId, actionName string, input map[string]any) error {
	device := m.getDevice(thingId)
	if device == nil {
		return errors.NotFoundError("device %s not found", thingId)
	}
	return device.requestAction(ctx, thingId, actionName, input)
}

func (m *Manager) RemoveAction(deviceId, actionId, actionName string) error {
	return nil
}

func (m *Manager) AddNewThings(timeout int) error {
	if m.pairTask != nil {
		return fmt.Errorf("add new things already in progress")
	}
	m.pairTask = make(chan bool)
	timeoutChan := time.After(time.Duration(timeout) * time.Millisecond)
	var handlePairingTimeout = func() {
		for _, adapter := range m.getAdapters() {
			log.Infof("%s to call startPairing on", adapter.GetId)
			adapter.startPairing(timeout)
		}
		for {
			select {
			case <-timeoutChan:
				m.pairTask = nil
				log.Info("pairing timeout")
				m.Publish(topic.PairingTimeout)
				m.CancelAddNewThing()
				return
			case <-m.pairTask:
				log.Info("pairing cancel")
				m.pairTask = nil
				return
			}
		}
	}
	go handlePairingTimeout()
	return nil
}

func (m *Manager) CancelAddNewThing() {
	if m.pairTask != nil {
		select {
		case m.pairTask <- true:
		}
	}
	for _, adapter := range m.getAdapters() {
		adapter.cancelPairing()
	}
	return
}

func (m *Manager) RemoveThing(deviceId string) error {

	task, ok := m.deferredRemove.LoadOrStore(deviceId, make(chan struct{}))
	if ok {
		return fmt.Errorf("remove already progress")
	}
	device := m.getDevice(deviceId)
	if device == nil {
		log.Infof("thing %s removed", deviceId)
		return nil
	}
	adapter := device.getAdapter()
	if adapter == nil {
		return fmt.Errorf("adapter not found")
	}
	m.deferredRemove.Store(deviceId, make(chan struct{}))
	removeTask, _ := task.(chan struct{})
	timeout := time.After(constant.DeviceRemovalTimeout * time.Millisecond)
	go func() {
		select {
		case <-timeout:
			m.CancelRemoveThing(deviceId)
		case <-removeTask:
			m.deferredRemove.Delete(deviceId)
			log.Infof("thing %s removed", deviceId)
			return
		}
	}()
	adapter.removeThing(device)
	return nil
}

func (m *Manager) CancelRemoveThing(deviceId string) {
	task, _ := m.deferredRemove.Load(deviceId)
	removeTask, ok := task.(chan struct{})
	if ok {
		select {
		case removeTask <- struct{}{}:
		}
	}
	device := m.getDevice(deviceId)
	if device == nil {
		return
	}
	adapter := device.getAdapter()
	if adapter == nil {
		return
	}
	adapter.cancelRemoveThing(device)
}

func (m *Manager) handleAdapterAdded(adapter *Adapter) {
	m.StoreAdapter(adapter)
	m.Publish(topic.AdapterAdded, adapter)
}

func (m *Manager) handleAdapterUnload(adapterId string) {
	m.DeleteAdapter(adapterId)
}

func (m *Manager) handleDeviceAdded(device *device) {
	m.StoreDevice(device)
	log.Debugf("Thing added:\t\n %s ", util.JsonIndent(things.AsWebOfThing(*device.Device)))
	m.Publish(topic.DeviceAdded, topic.DeviceAddedMessage{
		DeviceId: device.GetId(),
		Device:   *device.Device,
	})
}

func (m *Manager) handleDeviceRemoved(device *device) {
	m.DeleteDevice(device.GetId())
	task, ok := m.deferredRemove.Load(device.GetId())
	if ok {
		taskChan := task.(chan struct{})
		select {
		case taskChan <- struct{}{}:
			log.Infof("handle device removed")
		}
	}
	m.Publish(topic.DeviceRemoved, device.GetId())
}

func (m *Manager) getAdapter(adapterId string) *Adapter {
	a := m.GetAdapter(adapterId)
	adapter, ok := a.(*Adapter)
	if !ok {
		return nil
	}
	return adapter
}

func (m *Manager) getAdapters() (adapters []*Adapter) {
	adapters = make([]*Adapter, 0)
	for _, a := range m.GetAdapters() {
		adp, ok := a.(*Adapter)
		if ok {
			adapters = append(adapters, adp)
		}
	}
	return
}

func (m *Manager) getExtension(id string) *Extension {
	a, ok := m.extensions.Load(id)
	ext, ok := a.(*Extension)
	if !ok {
		return nil
	}
	return ext
}

func (m *Manager) getExtensions() (adapters []*Extension) {
	m.extensions.Range(func(key, value any) bool {
		ext, ok := value.(*Extension)
		if ok {
			adapters = append(adapters, ext)
		}
		return true
	})
	return
}

func (m *Manager) getDevice(deviceId string) *device {
	d := m.GetDevice(deviceId)
	device, ok := d.(*device)
	if ok {
		return device
	}
	return nil
}

func (m *Manager) getInstallAddon(addonId string) *Addon {
	a, ok := m.installAddons.Load(addonId)
	ad, ok := a.(*Addon)
	if !ok {
		return nil
	}
	return ad
}

func (m *Manager) getPlugin(packetId string) *Plugin {
	if m.pluginServer != nil {
		return m.pluginServer.getPlugin(packetId)
	}
	return nil
}

func (m *Manager) getInstallAddons() (addons []*Addon) {
	addons = make([]*Addon, 0)
	m.installAddons.Range(func(key, value any) bool {
		ad, ok := value.(*Addon)
		if ok {
			addons = append(addons, ad)
		}
		return true
	})
	return
}

//tar package to addon from the temp dir,
func (m *Manager) installAddon(packageId, packageDir string) error {
	if !m.addonsLoaded {
		return fmt.Errorf("cannot install add-on before other add-ons have been loaded")
	}
	log.Infof("start install %s", packageId)

	f, err := os.Open(packageDir)
	defer func() {
		e := f.Close()
		if e != nil {
			log.Error(e.Error())
		}
	}()
	zp, _ := gzip.NewReader(f)
	tr := tar.NewReader(zp)

	for hdr, err := tr.Next(); err != io.EOF; hdr, err = tr.Next() {
		if err != nil {
			continue
		}
		// 读取文件信息
		fi := hdr.FileInfo()
		p := strings.Replace(hdr.Name, "package", packageId, 1)

		localPath := m.config.AddonsDir + string(os.PathSeparator) + p
		if fi.IsDir() {
			_ = os.MkdirAll(localPath, os.ModePerm)
			continue
		}
		// 创建一个空文件，用来写入解包后的数据
		fw, err := os.Create(localPath)
		if err != nil {
			continue
		}
		if _, err := io.Copy(fw, tr); err != nil {
			return err
		}
		//TODO 给下载的文件增加可执行权限
		_ = os.Chmod(fi.Name(), 777)
		_ = fw.Close()
	}
	err = m.LoadAddon(packageId)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) loadAddons() {
	if m.addonsLoaded {
		return
	}
	log.Info("starting loading addons.")

	if m.pluginServer == nil {
		m.pluginServer = NewPluginServer(m)
	}
	load := func(dir string) {
		fs, err := os.ReadDir(dir)
		if err != nil {
			log.Warningf("load addon  %s ,err: %s", dir, err.Error())
			return
		}
		for _, fi := range fs {
			if fi.IsDir() {
				err := m.LoadAddon(fi.Name())
				if err != nil {
					log.Errorf(err.Error())
					continue
				}
			}
		}
	}
	load(m.config.AddonsDir)
	m.addonsLoaded = true
	//每天的23：00更新一次Add-on
	c := cron.New()
	log.Infof("time task for update addons every day 23:00")
	_ = c.AddFunc("\"0 0 23 * * ?\"", func() {
		m.updateAddons()
	})
	c.Start()
	return
}

func (m *Manager) removeAdapter(adapter *Adapter) {
	m.DeleteAdapter(adapter.GetId())
}

func (m *Manager) removeNotifier(notifierId string) {
	//TODO
}

func (m *Manager) handleOutletRemoved(outlet *Outlet) {
	m.outlets.Delete(outlet.getId())
}

func (m *Manager) removeApiHandler(id int) {

}

// 定时任务，更新Add-on
func (m *Manager) updateAddons() {
	log.Infof("time task: addons upgrade %s", time.Now().String())
}

func (m *Manager) GetPreferences() *messages.PluginRegisterResponseJsonDataPreferences {
	return &m.preferences
}

func (m *Manager) GetUserProfile() *messages.PluginRegisterResponseJsonDataUserProfile {
	return m.config.UserProfile
}

func (m *Manager) UpdatePreferences() {
	r := messages.PluginRegisterResponseJsonDataPreferences{
		Language: "en-US",
		Units: messages.PluginRegisterResponseJsonDataPreferencesUnits{
			Temperature: "degree celsius",
		},
	}
	lang, err := m.storage.GetSetting("localization.language")
	if err == nil {
		r.Language = lang
	} else {
		_ = m.storage.SetSetting("localization.language", r.Language)
	}
	temp, err := m.storage.GetSetting("localization.units.temperature")
	if err == nil {
		r.Units.Temperature = temp
	} else {
		_ = m.storage.SetSetting("localization.units.temperature", r.Units.Temperature)
	}
	m.preferences = r
}

func (m *Manager) GetLanguage() string {
	return m.GetPreferences().Language
}
