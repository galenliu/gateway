package plugin

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/galenliu/gateway"
	"github.com/galenliu/gateway-addon"
	"github.com/galenliu/gateway-addon/wot"
	"github.com/galenliu/gateway/pkg/database"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/util"
	data_schema2 "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	"github.com/galenliu/gateway/plugin/internal"
	json "github.com/json-iterator/go"

	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
)

type AddonManager interface {
	gateway.Program
	GetDevices() []addon.IDevice
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

	addonsLoaded bool
	isPairing    bool

	pluginCancel context.CancelFunc

	locker  *sync.Mutex
	running bool
	verbose bool
	logger  logging.Logger
	actions map[string]*wot.ActionAffordance
}

func NewAddonsManager(options Options, log logging.Logger) AddonManager {
	am := &manager{}
	instance = am
	am.options = options
	am.logger = log
	am.addonsLoaded = false
	am.isPairing = false
	am.running = false
	am.devices = make(map[string]addon.IDevice, 50)
	am.installAddons = make(map[string]*internal.AddonInfo, 50)
	am.adapters = make(map[string]*Adapter, 20)
	am.extensions = make(map[string]Extension)

	//def addon action
	action := wot.NewActionAffordance()
	obj := data_schema2.NewObjectSchema()
	timeout := data_schema2.NewIntegerSchema()
	timeout.Minimum = 1000
	timeout.Maximum = 10000
	obj.Properties["timeout"] = timeout
	action.Input = obj
	am.actions = make(map[string]*wot.ActionAffordance)
	am.actions["pair"] = action

	am.locker = new(sync.Mutex)
	am.loadAddons()
	return am
}

func (manager *manager) GetDevices() (device []addon.IDevice) {
	for _, dev := range manager.devices {
		device = append(device, dev)
	}
	return
}

func (manager *manager) handleDeviceAdded(device *addon.Device) {
	manager.devices[device.GetID()] = device
	//d, err := json.MarshalIndent(device, "", " ")
	d := device.AsDict()
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		logging.Info("device marshal err")
	}
	Publish(util.ThingAdded, data)
}

func (manager *manager) actionNotify(action *addon.Action) {
	Publish(util.ActionStatus, action.MarshalJson())
}

func (manager *manager) eventNotify(event *addon.Event) {
	Publish(util.EVENT, event.MarshalJson())
}

func (manager *manager) connectedNotify(device *addon.Device, connected bool) {
	Publish(util.CONNECTED, connected)
}

func (manager *manager) addAdapter(adapter *Adapter) {
	manager.locker.Lock()
	defer manager.locker.Unlock()
	manager.adapters[adapter.id] = adapter
	logging.Debug(fmt.Sprintf("adapter：(%s) added", adapter.id))
}

func (manager *manager) getAdapter(adapterId string) *Adapter {
	adapter, ok := manager.adapters[adapterId]
	if !ok {
		return nil
	}
	return adapter
}

func (manager *manager) handleSetProperty(deviceId, propName string, setValue interface{}) error {
	device := manager.getDevice(deviceId)
	if device == nil {
		return fmt.Errorf("device id err")
	}
	adapter := manager.getAdapter(device.GetAdapterId())
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

func (manager *manager) getDevice(deviceId string) addon.IDevice {
	d, ok := manager.devices[deviceId]
	if !ok {
		return nil
	}
	return d
}

// get package from url, checksum
func (manager *manager) installAddonFromUrl(id, url, checksum string, enabled bool) error {

	destPath := path.Join(os.TempDir(), id+".tar.gz")

	logging.Info("fetching add-on %s as %s", url, destPath)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Download addon err,pakage id:%s err:%s", id, err.Error()))
	}
	defer func() {
		_ = resp.Body.Close()
		err := os.Remove(destPath)
		if err != nil {
			logging.Info("remove temp file failed ,err:%s", err.Error())
		}
	}()
	data, _ := ioutil.ReadAll(resp.Body)
	_ = ioutil.WriteFile(destPath, data, 777)
	if !util.CheckSum(destPath, checksum) {
		return fmt.Errorf(fmt.Sprintf("checksum err,pakage id:%s", id))
	}
	err = manager.installAddon(id, destPath, enabled)
	if err != nil {
		return err
	}
	return nil
}

//tar package to addon from the temp dir,
func (manager *manager) installAddon(packageId, packagePath string, enabled bool) error {

	if !manager.addonsLoaded {
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
		// 读取文件信息
		fi := hdr.FileInfo()
		p := strings.Replace(hdr.Name, "package", packageId, 1)

		localPath := manager.options.AddonDirs[0] + string(os.PathSeparator) + p
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
		return manager.loadAddon(packageId)
	}
	return nil
}

func (manager *manager) loadAddons() {
	if manager.addonsLoaded {
		return
	}
	manager.addonsLoaded = true

	for _, d := range manager.options.AddonDirs {
		fs, err := os.ReadDir(d)
		if err != nil {
			logging.Error("load addon err: %s", err.Error())
			return
		}
		manager.pluginServer = NewPluginServer(manager)

		for _, fi := range fs {
			if fi.IsDir() {
				addonId := fi.Name()
				err = manager.loadAddon(addonId)
				if err != nil {
					manager.logger.Error("load addon id:%s failed err:%s", addonId, err.Error())
				}
			}
		}
	}
	return
}

func (manager *manager) loadAddon(packageId string) error {

	if !manager.addonsLoaded {
		return nil
	}

	addonPath, err := manager.findPlugin(packageId)
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
	manager.installAddons[packageId] = addonInfo
	if !addonInfo.Enabled {
		return fmt.Errorf("addon disenabled")
	}
	if addonInfo.ContentScripts != "" && addonInfo.WSebAccessibleResources != "" {
		var ext = Extension{
			Extensions: addonInfo.ContentScripts,
			Resources:  addonInfo.WSebAccessibleResources,
		}
		manager.extensions[addonInfo.ID] = ext
	}
	if addonInfo.Exec == "" {
		return nil
	}

	err = util.EnsureDir(path.Join(path.Join(manager.options.DataDir, "data"), addonInfo.ID))
	if err != nil {
		return err
	}

	manager.pluginServer.loadPlugin(addonPath, addonInfo.ID, addonInfo.Exec)

	return nil
}

func (manager *manager) unloadAddon(packageId string) error {
	if !manager.addonsLoaded {
		return nil
	}
	_, ok := manager.extensions[packageId]
	if ok {
		delete(manager.extensions, packageId)
	}
	plugin := manager.pluginServer.Plugins[packageId]
	if plugin == nil {
		return fmt.Errorf("can not found plugin: %s", packageId)
	}

	for key, adapter := range manager.adapters {
		if adapter.id == plugin.pluginId {
			for _, dev := range adapter.devices {
				adapter.handleDeviceRemoved(dev)
			}
			delete(manager.adapters, key)
		}
	}
	manager.pluginServer.uninstallPlugin(packageId)
	return nil
}

func (manager *manager) findPlugin(packageId string) (string, error) {
	for _, dir := range manager.options.AddonDirs {
		_, e := os.Stat(path.Join(dir, packageId))
		if os.IsNotExist(e) {
			continue
		}
		return path.Join(dir, packageId), nil
	}
	return "", fmt.Errorf("addon is not exist")
}

func (manager *manager) Start() error {
	var err error
	go func() {
		err = manager.pluginServer.Start()
		if err != nil {
			manager.running = false
			logging.Error(err.Error())
		}
	}()
	manager.running = true
	if err == nil {
		go Publish(util.PluginServerStarted)
	}
	return err
}

func (manager *manager) Stop() error {
	manager.pluginServer.Stop()
	Publish(util.PluginServerStopped)
	manager.running = false
	return nil
}
