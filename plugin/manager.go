package plugin

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/rpc"
	"github.com/galenliu/gateway/pkg/util"
	wot "github.com/galenliu/gateway/pkg/wot/definitions/core"
	"github.com/galenliu/gateway/plugin/internal"
	json "github.com/json-iterator/go"
	"io"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

type Store interface {
	GetSetting(key string) (string, error)
	SetSetting(key, value string) error
}

type Config struct {
	AddonDirs   []string
	IPCPort     string
	RPCPort     string
	UserProfile *rpc.PluginRegisterResponseMessage_Data_UsrProfile
	Preferences *rpc.PluginRegisterResponseMessage_Data_Preferences
}

type Manager struct {
	config       Config
	configPath   string
	pluginServer *PluginsServer

	devices       sync.Map
	adapters      sync.Map
	installAddons sync.Map
	extensions    sync.Map

	Eventbus     *Eventbus
	addonsLoaded bool
	isPairing    bool
	running      bool

	locker  *sync.Mutex
	logger  logging.Logger
	actions map[string]*wot.ActionAffordance

	store Store
}

func NewAddonsManager(conf Config, settingStore Store, bus bus, log logging.Logger) *Manager {
	am := &Manager{}
	am.config = conf
	am.logger = log
	am.addonsLoaded = false
	am.isPairing = false
	am.running = false
	am.Eventbus = NewEventBus(bus)
	am.store = settingStore
	am.locker = new(sync.Mutex)
	bus.SubscribeAsync(constant.GatewayStart, am.Start)
	return am
}

func (m *Manager) UnloadAddon(id string) error {
	panic("implement me")
}

func (m *Manager) LoadAddon(id string) error {
	panic("implement me")
}

func (m *Manager) GetPropertiesValue(thingId string) (map[string]interface{}, error) {
	panic("implement me")
}

func (m *Manager) addNewThing(pairingTimeout float64) error {
	if m.isPairing {
		return fmt.Errorf("add already in progress")
	}
	for _, adapter := range m.getAdapters() {
		adapter.pairing(pairingTimeout)
	}
	m.isPairing = true
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Duration(pairingTimeout)*time.Millisecond)
	var handlePairingTimeout = func() {
		for {
			select {
			case <-ctx.Done():
				cancelFn()
				m.cancelAddNewThing()
				//bus.Publish(util.PairingTimeout)
				return
			}
		}
	}
	go handlePairingTimeout()
	return nil
}

func (m *Manager) cancelAddNewThing() {
	if !m.isPairing {
		return
	}
	for _, adapter := range m.getAdapters() {
		adapter.cancelPairing()
	}
	m.isPairing = false
	return
}

func (m *Manager) actionNotify(action *internal.Action) {
	m.Eventbus.bus.Publish(constant.ActionStatus, nil)
}

func (m *Manager) eventNotify(event *internal.Event) {
	m.Eventbus.bus.Publish(constant.EVENT, nil)
}

func (m *Manager) connectedNotify(device *internal.Device, connected bool) {
	m.Eventbus.bus.Publish(constant.CONNECTED, connected)
}

func (m *Manager) handleAdapterAdded(adapter *Adapter) {
	m.adapters.Store(adapter.ID, adapter)
	m.logger.Debug(fmt.Sprintf("adapter：(%s) added", adapter.ID))
}

func (m *Manager) handleDeviceAdded(device *Device) {
	m.devices.Store(device.ID, device)
	data, err := json.MarshalIndent(device, "", "  ")
	if err != nil {
		m.logger.Info("device marshal err")
	}
	m.Eventbus.bus.Publish(constant.DeviceAdded, data)
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
	adapter := m.getAdapter(device.AdapterId)
	if adapter == nil {
		return fmt.Errorf("adapter ID err")
	}
	property := device.GetProperty(propName)

	if property == nil {
		return fmt.Errorf("property err")
	}

	//property.SetValue(setValue)
	//var newValue = property.ToValue(setValue)
	//data := make(map[string]interface{})
	//data[addon.Did] = device.GetID()
	//data["propertyName"] = property.GetName()
	//data["propertyValue"] = newValue
	go adapter.sendMessage(internal.DeviceSetPropertyCommand, nil)
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

func (m *Manager) getInstallAddon(addonId string) *AddonInfo {
	a, ok := m.installAddons.Load(addonId)
	addon, ok := a.(*AddonInfo)
	if !ok {
		return nil
	}
	return addon
}

func (m *Manager) getInstallAddons() (addons []*AddonInfo) {
	m.installAddons.Range(func(key, value interface{}) bool {
		addon, ok := value.(*AddonInfo)
		if ok {
			addons = append(addons, addon)
		}
		return true
	})
	return
}

//tar package to addon from the temp dir,
func (m *Manager) installAddon(packageId, packagePath string, enabled bool) error {
	if !m.addonsLoaded {
		return fmt.Errorf(`Cannot install add-on before other add-ons have been loaded.`)
	}
	m.logger.Info("execute install package ID: %s ", packageId)
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

		localPath := m.config.AddonDirs[0] + string(os.PathSeparator) + p
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
	var key = "addons." + packageId
	saved, err := m.store.GetSetting(key)
	if err == nil && saved != "" {
		var old AddonInfo
		ee := json.UnmarshalFromString(saved, &old)
		if ee != nil {
			old.Enabled = enabled
			newAddonInfo, err := json.MarshalToString(old)
			if err != nil {
				ee := m.store.SetSetting(key, newAddonInfo)
				if ee != nil {
					m.logger.Error(ee.Error())
				}
			}

		}
	}
	if enabled {
		//return m.loadAddon(packageId)
	}
	return nil
}

func (m *Manager) loadAddons() {
	if m.addonsLoaded {
		return
	}
	m.logger.Infof("starting loading addons.")
	m.addonsLoaded = true
	m.pluginServer = NewPluginServer(m)
	_ = m.pluginServer.Start()
	for _, d := range m.config.AddonDirs {
		fs, err := os.ReadDir(d)
		if err != nil {
			m.logger.Warningf("load addon  %s ,err: %s", d, err.Error())
			continue
		}

		for _, fi := range fs {
			if fi.IsDir() {
				addonId := fi.Name()
				err = m.loadAddon(addonId)
				if err != nil {
					m.logger.Error("load addon ID:%s failed err:%s", addonId, err.Error())
				}
			}
		}
	}
	return
}

func (m *Manager) loadAddon(packageId string) error {

	if !m.addonsLoaded {
		return nil
	}

	addonPath := m.findPluginPath(packageId)

	addonInfo, obj, err := LoadManifest(addonPath, packageId, m.store)

	if err != nil {
		return err
	}

	configKey := "addons.config." + packageId
	var cfg string
	if obj != nil {
		var ee error
		cfg, ee = json.MarshalToString(obj)
		if ee != nil {
			return err
		}
	}
	info, err := json.MarshalToString(addonInfo)
	err = m.store.SetSetting(GetAddonKey(addonInfo.ID), info)
	if err != nil {
		return err
	}
	savedConfig, e := m.store.GetSetting(configKey)
	if e != nil && savedConfig == "" {
		if cfg != "" {
			eee := m.store.SetSetting(configKey, cfg)
			if eee != nil {
				m.logger.Error(eee.Error())
			}
		}
	}
	if savedConfig == "" && cfg != "" {
		eee := m.store.SetSetting(configKey, cfg)
		if eee != nil {
			return eee
		}
	}
	m.installAddons.Store(packageId, addonInfo)

	if !addonInfo.Enabled {
		return fmt.Errorf("addon disenabled")
	}
	if addonInfo.ContentScripts != "" && addonInfo.WSebAccessibleResources != "" {
		var ext = Extension{
			Extensions: addonInfo.ContentScripts,
			Resources:  addonInfo.WSebAccessibleResources,
		}
		m.extensions.Store(addonInfo.ID, ext)
	}
	if addonInfo.Exec == "" {
		return nil
	}

	err = util.EnsureDir(path.Join(path.Join(m.config.UserProfile.BaseDir, "data"), addonInfo.ID))
	if err != nil {
		return err
	}

	m.pluginServer.loadPlugin(addonPath, addonInfo.ID, addonInfo.Exec)

	return nil
}

func (m *Manager) unloadAddon(pluginId string) error {
	if !m.addonsLoaded {
		return nil
	}
	ext := m.getExtension(pluginId)
	if ext != nil {
		ext.unload()
		m.extensions.Delete(pluginId)
	}
	for _, adapter := range m.getAdapters() {
		if adapter.pluginId == pluginId {
			for _, device := range adapter.getDevices() {
				if device.AdapterId == adapter.ID {
					adapter.handleDeviceRemoved(device)
				}
			}
		}
	}
	m.pluginServer.unloadPlugin(pluginId)
	return nil
}

func (m *Manager) findPluginPath(packageId string) string {
	for _, dir := range m.config.AddonDirs {
		_, e := os.Stat(path.Join(dir, packageId))
		if os.IsNotExist(e) {
			continue
		}
		return path.Join(dir, packageId)
	}
	return ""
}

func (m *Manager) Start() error {
	m.running = true
	m.loadAddons()
	go m.Eventbus.bus.Publish(constant.AddonManagerStarted)
	return nil
}

func (m *Manager) Stop() error {
	err := m.pluginServer.Stop()
	if err != nil {
		return err
	}
	go m.Eventbus.bus.Publish(constant.AddonManagerStopped)
	m.running = false
	return nil
}
