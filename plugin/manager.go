package plugin

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/galenliu/gateway-grpc"
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/container"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/util"
	wot "github.com/galenliu/gateway/pkg/wot/definitions/core"
	"github.com/galenliu/gateway/plugin/addon"
	json "github.com/json-iterator/go"
	"github.com/robfig/cron"
	"io"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

const (
	ActionPair   = "pair"
	ActionUnpair = "unpair"
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
	UserProfile     *rpc.UsrProfile
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
	container     container.Container
	Eventbus      *Eventbus
	addonsLoaded  bool
	isPairing     bool
	running       bool
	pairTask      chan struct{}
	locker        *sync.Mutex
	logger        logging.Logger
	actions       map[string]*wot.ActionAffordance
	storage       managerStore
	ctx           context.Context
}

func NewAddonsManager(ctx context.Context, conf Config, s managerStore, bus bus.Controller, log logging.Logger) *Manager {
	am := &Manager{}
	am.config = conf
	am.ctx = ctx
	am.logger = log
	am.addonsLoaded = false
	am.isPairing = false
	am.running = false
	am.Eventbus = NewEventBus(bus)
	am.storage = s
	am.locker = new(sync.Mutex)
	am.loadAddons()
	return am
}

func (m *Manager) GetPropertiesValue(thingId string) (map[string]interface{}, error) {
	panic("implement me")
}

func (m *Manager) AddNewThings(timeout int) error {
	if m.pairTask != nil {
		return fmt.Errorf("add new things already in progress")
	}
	m.logger.Info("pairing.....")
	m.pairTask = make(chan struct{})
	timeoutChan := time.After(time.Duration(timeout) * time.Millisecond)
	var handlePairingTimeout = func() {
		for _, adapter := range m.getAdapters() {
			adapter.pairing(timeout)
		}
		for {
			select {
			case <-timeoutChan:
				m.logger.Info("pairing timeout")
				m.Eventbus.bus.Publish(constant.PairingTimeout)
				m.CancelAddNewThing()
			case <-m.pairTask:
				m.logger.Info("pairing cancel")
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

func (m *Manager) actionNotify(action *addon.Action) {
	m.Eventbus.bus.Publish(constant.ActionStatus, nil)
}

func (m *Manager) eventNotify(event *addon.Event) {
	m.Eventbus.bus.Publish(constant.EVENT, nil)
}

func (m *Manager) connectedNotify(device *addon.Device, connected bool) {
	m.Eventbus.bus.Publish(constant.CONNECTED, connected)
}

func (m *Manager) addAdapter(adapter *Adapter) {
	m.adapters.Store(adapter.ID, adapter)
	m.Eventbus.bus.Publish(constant.AdapterAdded, adapter)
	m.logger.Debug(fmt.Sprintf("adapter：(%s) added", adapter.ID))
}

func (m *Manager) addService(service *Service) {
	m.adapters.Store(service.ID, service)
	m.Eventbus.bus.Publish(constant.ServiceAdded, service)
	m.logger.Debug(fmt.Sprintf("service：(%s) added", service.ID))
}

func (m *Manager) handleDeviceAdded(device *Device) {
	m.devices.Store(device.ID, device)
	m.Eventbus.bus.Publish(constant.DeviceAdded, device.Device)
}

func (m *Manager) handleDeviceRemoved(device *Device) {
	m.devices.Delete(device.ID)
	data, err := json.MarshalIndent(device, "", "  ")
	if err != nil {
		m.logger.Info("device marshal err")
	}
	m.Eventbus.bus.Publish(constant.DeviceAdded, data)
}

func (m *Manager) handleSetProperty(deviceId, propName string, setValue interface{}) error {
	device := m.getDevice(deviceId)
	if device == nil {
		return fmt.Errorf("device ID err")
	}
	adapter := m.getAdapter(device.ID)
	if adapter == nil {
		return fmt.Errorf("adapter ID err")
	}
	property := device.GetProperty(propName)

	if property == nil {
		return fmt.Errorf("property err")
	}
	return nil
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

func (m *Manager) getExtension(id string) *addon.Extension {
	a, ok := m.extensions.Load(id)
	ext, ok := a.(*addon.Extension)
	if !ok {
		return nil
	}
	return ext
}

func (m *Manager) getExtensions() (adapters []*addon.Extension) {
	m.extensions.Range(func(key, value interface{}) bool {
		ext, ok := value.(*addon.Extension)
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

func (m *Manager) getDevices() (devices []*Device) {
	m.devices.Range(func(key, value interface{}) bool {
		device, ok := value.(*Device)
		if ok {
			devices = append(devices, device)
		}
		return true
	})
	return
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
	m.logger.Infof("start install %s", packageId)
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
		m.pluginServer = NewPluginServer(m.ctx,m)
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
		err = addonInfo.Save()
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
		var ext = addon.Extension{
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
	m.adapters.Delete(adapter.ID)
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

func (m *Manager) handleOutletRemoved(device *addon.Outlet) {

}

func (m *Manager) removeApiHandler(id int) {

}

// 定时任务，更新Add-on
func (m *Manager) updateAddons() {
	m.logger.Infof("time task: checking addons update time:", time.Now().String())
}
