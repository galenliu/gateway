package plugin

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/galenliu/gateway-addon"
	"github.com/galenliu/gateway/pkg/database"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/util"
	wot "github.com/galenliu/gateway/pkg/wot/definitions/core"
	"github.com/galenliu/gateway/plugin/internal"
	json "github.com/json-iterator/go"

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
	DataDir   string
}

type Extension struct {
	Extensions string
	Resources  string
}

type manager struct {
	options      Options
	configPath   string
	pluginServer *PluginsServer

	devices       map[string]addon.IDevice
	adapters      map[string]*Adapter
	installAddons map[string]*internal.AddonInfo

	extensions map[string]Extension

	bus eventBus

	addonsLoaded bool
	isPairing    bool

	pluginCancel context.CancelFunc

	locker  *sync.Mutex
	running bool
	verbose bool
	logger  logging.Logger
	actions map[string]*wot.ActionAffordance
}

func NewAddonsManager(options Options, bus eventBus, log logging.Logger) AddonManager {
	am := &manager{}
	am.options = options
	am.logger = log
	am.addonsLoaded = false
	am.isPairing = false
	am.running = false
	am.devices = make(map[string]addon.IDevice, 50)
	am.installAddons = make(map[string]*internal.AddonInfo, 50)
	am.adapters = make(map[string]*Adapter, 20)
	am.extensions = make(map[string]Extension)
	am.bus = bus

	//def addon action
	//action := wot.NewActionAffordance()
	//obj := schema.NewObjectSchema()
	//timeout := schema.NewIntegerSchema()
	//timeout.Minimum = 1000
	//timeout.Maximum = 10000
	//obj.Properties["timeout"] = timeout
	//action.Input = obj
	//am.actions = make(map[string]*wot.ActionAffordance)
	//am.actions["pair"] = action

	am.locker = new(sync.Mutex)
	am.loadAddons()
	return am
}

func (m *manager) handleDeviceAdded(device *addon.Device) {
	m.devices[device.GetID()] = device
	//d, err := json.MarshalIndent(device, "", " ")
	d := device.AsDict()
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		logging.Info("device marshal err")
	}
	m.bus.Publish(util.DeviceAdded, data)
}

func (m *manager) actionNotify(action *addon.Action) {
	m.bus.Publish(util.ActionStatus, action.MarshalJson())
}

func (m *manager) eventNotify(event *addon.Event) {
	m.bus.Publish(util.EVENT, event.MarshalJson())
}

func (m *manager) connectedNotify(device *addon.Device, connected bool) {
	m.bus.Publish(util.CONNECTED, connected)
}

func (m *manager) addAdapter(adapter *Adapter) {
	m.locker.Lock()
	defer m.locker.Unlock()
	m.adapters[adapter.id] = adapter
	logging.Debug(fmt.Sprintf("adapter???(%s) added", adapter.id))
}

func (m *manager) getAdapter(adapterId string) *Adapter {
	adapter, ok := m.adapters[adapterId]
	if !ok {
		return nil
	}
	return adapter
}

func (m *manager) handleSetProperty(deviceId, propName string, setValue interface{}) error {
	device := m.getDevice(deviceId)
	if device == nil {
		return fmt.Errorf("device id err")
	}
	adapter := m.getAdapter(device.GetAdapterId())
	if adapter == nil {
		return fmt.Errorf("adapter id err")
	}
	property := device.GetProperty(propName)

	if property == nil {
		return fmt.Errorf("property err")
	}

	var newValue = property.ToValue(setValue)
	data := make(map[string]interface{})
	data[addon.Did] = device.GetID()
	data["propertyName"] = property.GetName()
	data["propertyValue"] = newValue
	go adapter.Send(internal.DeviceSetPropertyCommand, data)
	return nil
}

func (m *manager) getDevice(deviceId string) addon.IDevice {
	d, ok := m.devices[deviceId]
	if !ok {
		return nil
	}
	return d
}

//tar package to addon from the temp dir,
func (m *manager) installAddon(packageId, packagePath string, enabled bool) error {

	if !m.addonsLoaded {
		return fmt.Errorf(`Cannot install add-on before other add-ons have been loaded.`)
	}
	logging.Info("execute install package id: %s ", packageId)
	f, err := os.Open(packagePath)
	if err != nil {
		return err
	}
	defer func() {
		e := f.Close()
		if e != nil {
			logging.Error(e.Error())
		}
	}()

	zp, _ := gzip.NewReader(f)
	tr := tar.NewReader(zp)

	for hdr, err := tr.Next(); err != io.EOF; hdr, err = tr.Next() {
		if err != nil {
			continue
		}
		// ??????????????????
		fi := hdr.FileInfo()
		p := strings.Replace(hdr.Name, "package", packageId, 1)

		localPath := m.options.AddonDirs[0] + string(os.PathSeparator) + p
		if fi.IsDir() {
			_ = os.MkdirAll(localPath, os.ModePerm)
			continue
		}
		// ??????????????????????????????????????????????????????
		fw, err := os.Create(localPath)
		if err != nil {
			continue
		}
		if _, err := io.Copy(fw, tr); err != nil {
			return err
		}
		//TODO ???????????????????????????????????????
		_ = os.Chmod(fi.Name(), 777)
		_ = fw.Close()

	}
	var key = "addons." + packageId
	saved, err := database.GetSetting(key)
	if err == nil && saved != "" {
		var old internal.AddonInfo
		ee := json.UnmarshalFromString(saved, &old)
		if ee != nil {
			old.Enabled = enabled
			newAddonInfo, err := json.MarshalToString(old)
			if err != nil {
				ee := database.SetSetting(key, newAddonInfo)
				if ee != nil {
					logging.Error(ee.Error())
				}
			}

		}
	}
	if enabled {
		return m.loadAddon(packageId)
	}
	return nil
}

func (m *manager) loadAddons() {
	if m.addonsLoaded {
		return
	}
	m.addonsLoaded = true

	for _, d := range m.options.AddonDirs {
		fs, err := os.ReadDir(d)
		if err != nil {
			logging.Error("load addon err: %s", err.Error())
			return
		}
		m.pluginServer = NewPluginServer(m)

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

func (m *manager) loadAddon(packageId string) error {

	if !m.addonsLoaded {
		return nil
	}

	addonPath, err := m.findPlugin(packageId)
	if err != nil {
		return err
	}

	addonInfo, obj, err := internal.LoadManifest(addonPath, packageId)

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

	err = addonInfo.UpdateFromDB()
	if err != nil {
		return err
	}

	savedConfig, e := database.GetSetting(configKey)
	if e != nil && savedConfig == "" {
		if cfg != "" {
			eee := database.SetSetting(configKey, cfg)
			if eee != nil {
				logging.Error(eee.Error())
			}
		}
	}
	if savedConfig == "" && cfg != "" {
		eee := database.SetSetting(configKey, cfg)
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

	err = util.EnsureDir(path.Join(path.Join(m.options.DataDir, "data"), addonInfo.ID))
	if err != nil {
		return err
	}

	m.pluginServer.loadPlugin(addonPath, addonInfo.ID, addonInfo.Exec)

	return nil
}

func (m *manager) unloadAddon(packageId string) error {
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

func (m *manager) findPlugin(packageId string) (string, error) {
	for _, dir := range m.options.AddonDirs {
		_, e := os.Stat(path.Join(dir, packageId))
		if os.IsNotExist(e) {
			continue
		}
		return path.Join(dir, packageId), nil
	}
	return "", fmt.Errorf("addon is not exist")
}

func (m *manager) Start() error {
	var err error
	go func() {
		err = m.pluginServer.Start()
		if err != nil {
			m.running = false
			logging.Error(err.Error())
		}
	}()
	m.running = true
	if err == nil {
		m.bus.Publish(util.AddonManagerStarted)
	}
	return err
}

func (m *manager) Stop() error {
	m.pluginServer.Stop()
	m.bus.Publish(util.AddonManagerStopped)
	m.running = false
	return nil
}
