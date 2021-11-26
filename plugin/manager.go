package plugin

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/galenliu/gateway/pkg/addon"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/container"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/util"
	wot "github.com/galenliu/gateway/pkg/wot/definitions/core"
	"github.com/robfig/cron"
	"io"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

type managerStore interface {
	AddonsStore
	GetSetting(key string) (string, error)
}

type Config struct {
	AddonsDir       string
	AttachAddonsDir string
	IPCPort         string
	RPCPort         string
	UserProfile     *messages.PluginRegisterResponseJsonDataUserProfile
}

type Manager struct {
	config        Config
	configPath    string
	pluginServer  *PluginsServer
	devices       sync.Map
	adapters      sync.Map
	services      sync.Map
	installAddons sync.Map
	extensions    sync.Map
	container     *container.ThingsContainer
	bus           managerBus
	addonsLoaded  bool
	isPairing     bool
	running       bool
	pairTask      chan struct{}
	locker        *sync.Mutex
	logger        logging.Logger
	actions       map[string]*wot.ActionAffordance
	storage       managerStore
	ctx           context.Context

	deferredRemove sync.Map
}

type managerBus interface {
	PublishDeviceAdded(device *addon.Device)
	PublishDeviceRemoved(deviceId string)
	PublishPairingTimeout()
	PublishActionStatus(action interface{})
	PublishConnected(thingId string, connected bool)
	PublishEvent(event *addon.Event)
	PublishPropertyChanged(thingId string, property *addon.PropertyDescription)
	AddPropertyChangedSubscription(thingId string, fn func(p *addon.PropertyDescription)) func()
}

func NewAddonsManager(ctx context.Context, conf Config, s managerStore, bus managerBus, log logging.Logger) *Manager {
	am := &Manager{}
	am.config = conf
	am.ctx = ctx
	am.logger = log
	am.addonsLoaded = false
	am.isPairing = false
	am.running = false
	am.bus = bus
	am.storage = s
	am.locker = new(sync.Mutex)
	am.loadAddons()
	return am
}

func (m *Manager) RequestAction(deviceId, actionName string, input map[string]interface{}) error {
	device := m.getDevice(deviceId)
	if device == nil {
		return fmt.Errorf("device %s not found", deviceId)
	}
	return nil
}

func (m *Manager) AddNewThings(timeout int) error {
	if m.pairTask != nil {
		return fmt.Errorf("add new things already in progress")
	}
	m.pairTask = make(chan struct{})
	timeoutChan := time.After(time.Duration(timeout) * time.Millisecond)
	var handlePairingTimeout = func() {
		for _, adapter := range m.getAdapters() {
			adapter.startPairing(timeout)
		}
		for {
			select {
			case <-timeoutChan:
				m.logger.Info("startPairing timeout")
				m.bus.PublishPairingTimeout()
				m.CancelAddNewThing()
			case <-m.pairTask:
				m.logger.Info("startPairing cancel")
				return
			}
		}
	}
	go handlePairingTimeout()
	return nil
}

func (m *Manager) RemoveThing(deviceId string) error {

	task, ok := m.deferredRemove.LoadOrStore(deviceId, make(chan struct{}))
	if ok {
		return fmt.Errorf("remove already progress")
	}
	device := m.getDevice(deviceId)
	if device == nil {
		m.logger.Infof("thing %s removed", deviceId)
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
			m.cancelRemoveThing(deviceId)
		case <-removeTask:
			m.deferredRemove.Delete(deviceId)
			m.logger.Infof("thing %s removed", deviceId)
			return
		}
	}()
	adapter.removeThing(device)
	return nil
}

func (m *Manager) CancelAddNewThing() {
	if m.pairTask != nil {
		select {
		case m.pairTask <- struct{}{}:
		default:
			break
		}
	}
	for _, adapter := range m.getAdapters() {
		adapter.cancelPairing()
	}
	m.pairTask = nil
	return
}

func (m *Manager) cancelRemoveThing(deviceId string) {
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

func (m *Manager) actionNotify(action *addon.Action) {
	m.bus.PublishActionStatus(action)
}

func (m *Manager) eventNotify(event *addon.Event) {
	m.bus.PublishEvent(event)
}

func (m *Manager) addAdapter(adapter *Adapter) {
	m.adapters.Store(adapter.id, adapter)
}

func (m *Manager) handleDeviceAdded(device *Device) {
	m.devices.Store(device.GetId(), device)
	m.logger.Infof("Device added: %s", util.JsonIndent(container.AsWebOfThing(device.Device)))
	go m.bus.PublishDeviceAdded(device.Device)
}

func (m *Manager) handleDeviceRemoved(device *Device) {
	m.devices.Delete(device.GetId())
	task, ok := m.deferredRemove.Load(device.GetId())
	taskChan := task.(chan struct{})
	if ok {
		select {
		case taskChan <- struct{}{}:
			m.logger.Infof("handle device removed")
		}
	}
	m.bus.PublishDeviceRemoved(device.GetId())
}

func (m *Manager) getAdapter(adapterId string) *Adapter {
	a, ok := m.adapters.Load(adapterId)
	adapter, ok := a.(*Adapter)
	if !ok {
		return nil
	}
	return adapter
}

func (m *Manager) getAdapters() (adapters []*Adapter) {
	m.adapters.Range(func(key, value interface{}) bool {
		adapter, ok := value.(*Adapter)
		if ok {
			adapters = append(adapters, adapter)
		}
		return true
	})
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
	m.extensions.Range(func(key, value interface{}) bool {
		ext, ok := value.(*Extension)
		if ok {
			adapters = append(adapters, ext)
		}
		return true
	})
	return
}

func (m *Manager) getDevice(deviceId string) *Device {
	d, ok := m.devices.Load(deviceId)
	device, ok := d.(*Device)
	if !ok {
		return nil
	}
	return device
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
	m.installAddons.Range(func(key, value interface{}) bool {
		ad, ok := value.(*Addon)
		if ok {
			addons = append(addons, ad)
		}
		return true
	})
	return
}

//tar package to addon from the temp dir,
func (m *Manager) installAddon(packageId, packagePath string) error {
	if !m.addonsLoaded {
		return fmt.Errorf("cannot install add-on before other add-ons have been loaded")
	}
	m.logger.Infof("run install %s", packageId)
	f, err := os.Open(packagePath)
	if err != nil {
		return err
	}
	defer func() {
		e := f.Close()
		if e != nil {
			m.logger.Error(e.Error())
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
	m.loadAddon(packageId)
	return nil
}

func (m *Manager) loadAddons() {
	if m.addonsLoaded {
		return
	}
	m.logger.Info("starting loading addons.")
	m.addonsLoaded = true
	if m.pluginServer == nil {
		m.pluginServer = NewPluginServer(m)
	}
	load := func(dir string) {
		fs, err := os.ReadDir(dir)
		if err != nil {
			m.logger.Warningf("load addon  %s ,err: %s", dir, err.Error())
			return
		}
		for _, fi := range fs {
			if fi.IsDir() {
				m.loadAddon(fi.Name())

			}
		}
	}
	load(m.config.AddonsDir)
	if m.config.AttachAddonsDir != "" {
		load(m.config.AttachAddonsDir)
	}
	//每天的23：00更新一次Add-on
	c := cron.New()
	m.logger.Infof("time task for update addons every day 23:00")
	_ = c.AddFunc("\"0 0 23 * * ?\"", func() {
		m.updateAddons()
	})
	c.Start()
	return
}

func (m *Manager) loadAddon(packageId string) {
	m.logger.Infof("starting loading addon %s.", packageId)
	var addonInfo *Addon
	var obj interface{}
	var err error
	packageDir := m.getAddonPath(packageId)
	addonInfo, obj, err = LoadManifest(packageDir, packageId, m.storage)
	if err != nil {
		m.logger.Errorf("load file %s%s"+"manifest.json  err:", os.PathSeparator, packageId)
		return
	}
	saved, err := m.storage.LoadAddonSetting(packageId)
	if err == nil && saved != "" {
		addonInfo = NewAddonSettingFromString(saved, m.storage)
	} else {
		err = addonInfo.save()
		if err != nil {
			m.logger.Errorf("addon save err: %s", err.Error())
		}
	}
	addonConf, err := m.storage.LoadAddonConfig(packageId)
	if err != nil && addonConf == "" {
		if obj != nil {
			err := m.storage.StoreAddonsConfig(packageId, obj)
			if err != nil {
				m.logger.Errorf("store addon config err: %s", err.Error())
			}
		}
	}

	m.installAddons.Store(packageId, addonInfo)

	if addonInfo.ContentScripts != "" && addonInfo.WSebAccessibleResources != "" {
		var ext = Extension{
			Extensions: addonInfo.ContentScripts,
			Resources:  addonInfo.WSebAccessibleResources,
		}
		m.extensions.Store(addonInfo.ID, ext)
	}
	util.EnsureDir(m.logger, path.Join(m.config.UserProfile.DataDir, packageId))
	if addonInfo.Exec == "" {
		m.logger.Errorf("addon %s has not exec", addonInfo.ID)
		return
	}

	if !addonInfo.Enabled {
		m.logger.Errorf("addon %s disabled", packageId)
		return
	}
	m.pluginServer.loadPlugin(addonInfo.ID, packageDir, addonInfo.Exec)
}

func (m *Manager) removeAdapter(adapter *Adapter) {
	m.adapters.Delete(adapter.id)
}

func (m *Manager) getAddonPath(packageId string) string {
	for _, dir := range []string{m.config.AddonsDir, m.config.AttachAddonsDir} {
		if dir == "" {
			continue
		}
		_, e := os.Stat(path.Join(dir, packageId))
		if os.IsNotExist(e) {
			continue
		}
		return path.Join(dir, packageId)
	}
	return ""
}

func (m *Manager) removeNotifier(notifierId string) {

}

func (m *Manager) handleOutletRemoved(outlet *Outlet) {

}

func (m *Manager) removeApiHandler(id int) {

}

// 定时任务，更新Add-on
func (m *Manager) updateAddons() {
	m.logger.Infof("time task: addons upgrade %s", time.Now().String())
}
