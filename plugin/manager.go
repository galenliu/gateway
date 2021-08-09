package plugin

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/util"
	wot "github.com/galenliu/gateway/pkg/wot/definitions/core"
	"github.com/galenliu/gateway/plugin/internal"
	"github.com/galenliu/gateway/server/models"
	json "github.com/json-iterator/go"
	"time"

	"io"
	"os"
	"path"
	"strings"
	"sync"
)

type eventBus interface {
	Publish(string, ...interface{})
	Subscribe(string, interface{})
	Unsubscribe(string, interface{})
}

type Options struct {
	AddonDirs []string
	IPCPort   string
	UserProfile
	Preferences
}

type Preferences struct {
}

type UserProfile struct {
	BaseDir        string
	DataDir        string
	AddonsDir      string
	ConfigDir      string
	UploadDir      string
	MediaDir       string
	LogDir         string
	GatewayVersion string
}

type Extension struct {
	Extensions string
	Resources  string
}

type Manager struct {
	options       Options
	configPath    string
	pluginServer  *PluginsServer

	devices       map[string]*internal.Device
	adapters      map[string]*Adapter
	installAddons map[string]*AddonInfo
	extensions    map[string]Extension

	bus          eventBus
	addonsLoaded bool
	isPairing    bool
	pluginCancel context.CancelFunc
	locker       *sync.Mutex
	running      bool
	verbose      bool
	logger       logging.Logger
	actions      map[string]*wot.ActionAffordance

	settingsStore models.SettingsStore
}

func NewAddonsManager(options Options, settingStore models.SettingsStore, bus eventBus, log logging.Logger) *Manager {
	am := &Manager{}
	am.options = options
	am.logger = log
	am.addonsLoaded = false
	am.isPairing = false
	am.running = false
	am.devices = make(map[string]*internal.Device, 50)
	am.installAddons = make(map[string]*AddonInfo, 50)
	am.adapters = make(map[string]*Adapter, 20)
	am.extensions = make(map[string]Extension)
	am.bus = bus
	am.settingsStore = settingStore
	am.locker = new(sync.Mutex)
	am.loadAddons()
	return am
}

func (m *Manager) addNewThing(pairingTimeout float64) error {
	if m.isPairing {
		return fmt.Errorf("add already in progress")
	}
	for _, adapter := range m.adapters {
		adapter.pairing(pairingTimeout)
	}
	m.isPairing = true
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Duration(pairingTimeout)*time.Millisecond)
	var handlePairingTimeout = func() {
		for {
			select {
			case <-ctx.Done():
				cancelFn()
				m.CancelAddNewThing()
				//bus.Publish(util.PairingTimeout)
				return
			}
		}
	}
	go handlePairingTimeout()
	return nil
}

func (m *Manager) CancelAddNewThing() {
	if !m.isPairing {
		return
	}
	for _, adapter := range m.adapters {
		adapter.cancelPairing()
	}
	m.isPairing = false
	return
}

func (m *Manager) actionNotify(action *internal.Action) {
	m.bus.Publish(constant.ActionStatus, nil)
}

func (m *Manager) eventNotify(event *internal.Event) {
	m.bus.Publish(constant.EVENT, nil)
}

func (m *Manager) connectedNotify(device *internal.Device, connected bool) {
	m.bus.Publish(constant.CONNECTED, connected)
}

func (m *Manager) handleAdapterAdded(adapter *Adapter) {
	m.locker.Lock()
	defer m.locker.Unlock()
	m.adapters[adapter.id] = adapter
	m.logger.Debug(fmt.Sprintf("adapter：(%s) added", adapter.id))
}

func (m *Manager) handleDeviceAdded(device *internal.Device) {
	m.locker.Lock()
	defer m.locker.Unlock()
	m.devices[device.GetId()] = device
	data, err := json.MarshalIndent(device, "", "  ")
	if err != nil {
		m.logger.Info("device marshal err")
	}
	m.bus.Publish(constant.DeviceAdded, data)
}

func (m *Manager) getAdapter(adapterId string) *Adapter {
	adapter, ok := m.adapters[adapterId]
	if !ok {
		return nil
	}
	return adapter
}

func (m *Manager) handleSetProperty(deviceId, propName string, setValue interface{}) error {
	device := m.devices[deviceId]
	if device == nil {
		return fmt.Errorf("device id err")
	}
	adapter := m.getAdapter(device.AdapterId)
	if adapter == nil {
		return fmt.Errorf("adapter id err")
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
	go adapter.Send(internal.DeviceSetPropertyCommand, nil)
	return nil
}

func (m *Manager) getDevice(deviceId string) *internal.Device {
	d, ok := m.devices[deviceId]
	if !ok {
		return nil
	}
	return d
}

//tar package to addon from the temp dir,
func (m *Manager) installAddon(packageId, packagePath string, enabled bool) error {
	if !m.addonsLoaded {
		return fmt.Errorf(`Cannot install add-on before other add-ons have been loaded.`)
	}
	m.logger.Info("execute install package id: %s ", packageId)
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

		localPath := m.options.AddonDirs[0] + string(os.PathSeparator) + p
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
	saved, err := m.settingsStore.GetSetting(key)
	if err == nil && saved != "" {
		var old AddonInfo
		ee := json.UnmarshalFromString(saved, &old)
		if ee != nil {
			old.Enabled = enabled
			newAddonInfo, err := json.MarshalToString(old)
			if err != nil {
				ee := m.settingsStore.SetSetting(key, newAddonInfo)
				if ee != nil {
					m.logger.Error(ee.Error())
				}
			}

		}
	}
	if enabled {
		return m.loadAddon(packageId)
	}
	return nil
}

func (m *Manager) loadAddons() {
	if m.addonsLoaded {
		return
	}
	m.addonsLoaded = true
	m.pluginServer = NewPluginServer(m)
	err := m.pluginServer.Start()
	if err != nil {
		m.logger.Error("Plugin Server Start Failed. Err: %s", err.Error())
		return
	}

	for _, d := range m.options.AddonDirs {
		fs, err := os.ReadDir(d)
		if err != nil {
			m.logger.Warning("load addon form path: %s ,err: %s", d, err.Error())
			continue
		}

		for _, fi := range fs {
			if fi.IsDir() {
				addonId := fi.Name()
				err = m.loadAddon(addonId)
				if err != nil {
					m.logger.Error("load addon id:%s failed err:%s", addonId, err.Error())
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

	addonInfo, obj, err := LoadManifest(addonPath, packageId)

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
	err = m.settingsStore.SetSetting(GetAddonKey(addonInfo.ID), info)
	if err != nil {
		return err
	}
	savedConfig, e := m.settingsStore.GetSetting(configKey)
	if e != nil && savedConfig == "" {
		if cfg != "" {
			eee := m.settingsStore.SetSetting(configKey, cfg)
			if eee != nil {
				m.logger.Error(eee.Error())
			}
		}
	}
	if savedConfig == "" && cfg != "" {
		eee := m.settingsStore.SetSetting(configKey, cfg)
		if eee != nil {
			return eee
		}
	}
	m.installAddons[packageId] = addonInfo
	if !addonInfo.Enabled {
		return fmt.Errorf("addon disenabled")
	}
	if addonInfo.ContentScripts != "" && addonInfo.WSebAccessibleResources != "" {
		var ext = Extension{
			Extensions: addonInfo.ContentScripts,
			Resources:  addonInfo.WSebAccessibleResources,
		}
		m.extensions[addonInfo.ID] = ext
	}
	if addonInfo.Exec == "" {
		return nil
	}

	err = util.EnsureDir(path.Join(path.Join(m.options.BaseDir, "data"), addonInfo.ID))
	if err != nil {
		return err
	}

	m.pluginServer.loadPlugin(addonPath, addonInfo.ID, addonInfo.Exec)

	return nil
}

func (m *Manager) unloadAddon(packageId string) error {
	if !m.addonsLoaded {
		return nil
	}
	_, ok := m.extensions[packageId]
	if ok {
		delete(m.extensions, packageId)
	}
	plugin := m.pluginServer.Plugins[packageId]
	if plugin == nil {
		return fmt.Errorf("can not found plugin: %s", packageId)
	}

	for key, adapter := range m.adapters {
		if adapter.id == plugin.pluginId {
			for _, dev := range adapter.devices {
				adapter.handleDeviceRemoved(dev)
			}
			delete(m.adapters, key)
		}
	}
	m.pluginServer.uninstallPlugin(packageId)
	return nil
}

func (m *Manager) findPluginPath(packageId string) string {
	for _, dir := range m.options.AddonDirs {
		_, e := os.Stat(path.Join(dir, packageId))
		if os.IsNotExist(e) {
			continue
		}
		return path.Join(dir, packageId)
	}
	return ""
}

func (m *Manager) Start() error {
	var err error
	go func() {
		err = m.pluginServer.Start()
		if err != nil {
			m.running = false
			m.logger.Error(err.Error())
		}
	}()
	m.running = true
	if err == nil {
		m.bus.Publish(constant.AddonManagerStarted)
	}
	return err
}

func (m *Manager) Stop() error {
	err := m.pluginServer.Stop()
	if err != nil {
		return err
	}
	m.bus.Publish(constant.AddonManagerStopped)
	m.running = false
	return nil
}
