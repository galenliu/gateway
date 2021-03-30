package plugin

import (
	"addon"
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"gateway/config"
	"gateway/pkg/bus"
	"gateway/pkg/database"
	"gateway/pkg/log"
	"gateway/pkg/util"
	"gateway/server/models/thing"
	json "github.com/json-iterator/go"
	"github.com/xiam/to"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
)

var instance *AddonManager

type ThingAdded func(*thing.Thing)
type ThingRemoved func(*thing.Thing)
type PropertyChanged func(property addon.Property)

type Extension struct {
	Extensions string
	Resources  string
}

type AddonManager struct {
	configPath   string
	pluginServer *PluginsServer

	devices       map[string]*addon.Device
	adapters      map[string]*AdapterProxy
	installAddons map[string]*AddonInfo

	extensions map[string]Extension

	addonsLoaded bool
	isPairing    bool

	pluginCancel context.CancelFunc

	AddonsDir string
	DataDir   string
	locker    *sync.Mutex
	running   bool
}

func NewAddonsManager() *AddonManager {
	am := &AddonManager{}
	instance = am

	am.AddonsDir = config.Conf.AddonsDir
	am.DataDir = config.Conf.DataDir
	am.addonsLoaded = false
	am.isPairing = false
	am.running = false
	am.devices = make(map[string]*addon.Device, 50)
	am.installAddons = make(map[string]*AddonInfo, 50)
	am.adapters = make(map[string]*AdapterProxy, 20)
	am.extensions = make(map[string]Extension)

	_ = bus.Subscribe(bus.SetProperty, am.handleSetProperty)
	_ = bus.Subscribe(bus.GetDevices, am.handleGetDevices)
	_ = bus.Subscribe(bus.GetThings, am.handleGetThings)

	am.locker = new(sync.Mutex)
	am.LoadAddons()
	return am
}

func (manager *AddonManager) handleDeviceAdded(device *addon.Device) {
	manager.devices[device.ID] = device
	bus.Publish(util.ThingAdded, asWebThing(device))

}

func (manager *AddonManager) handleDeviceRemoved(device *addon.Device) {
	delete(manager.devices, device.ID)
	bus.Publish(util.ThingRemoved, asWebThing(device))
}

func (manager *AddonManager) actionNotify(action *addon.Action) {
	bus.Publish(util.ActionStatus, action)
}

func (manager *AddonManager) eventNotify(event *addon.Event) {
	bus.Publish(util.EVENT, event)
}

func (manager *AddonManager) connectedNotify(device *addon.Device, connected bool) {
	bus.Publish(util.CONNECTED, connected)
}

func (manager *AddonManager) addAdapter(adapter *AdapterProxy) {
	manager.locker.Lock()
	defer manager.locker.Unlock()
	manager.adapters[adapter.ID] = adapter
	log.Debug(fmt.Sprintf("adapter：(%s) added", adapter.ID))
}

func (manager *AddonManager) getAdapter(adapterId string) *AdapterProxy {
	adapter, ok := manager.adapters[adapterId]
	if !ok {
		return nil
	}
	return adapter
}

func (manager *AddonManager) findAdapter(adapterId string) (*AdapterProxy, error) {
	if adapterId == "" {
		return nil, fmt.Errorf("adapter id none")
	}
	adapter, ok := manager.adapters[adapterId]
	if !ok {
		return nil, fmt.Errorf("adapter(%s) not found", adapterId)
	}
	return adapter, nil
}

func (manager *AddonManager) handleSetProperty(deviceId, propName string, setValue interface{}) error {
	adapter := manager.getAdapterByDeviceId(deviceId)
	if adapter == nil {
		return fmt.Errorf("adapter not found")
	}
	property := adapter.GetDevice(deviceId).GetProperty(propName)
	var newValue interface{}
	if property.Type == addon.TypeBoolean {
		newValue = to.Bool(setValue)
	}
	if property.Type == addon.TypeInteger || property.Type == addon.TypeNumber {
		newValue = to.Int64(to.Bytes(setValue))
	}
	if property.Type == addon.TypeString {
		newValue = to.String(setValue)
	}

	if property == nil {
		return fmt.Errorf("device or property not found")
	}

	go adapter.handleSetPropertyValue(property, newValue)
	return nil
}

func (manager *AddonManager) handleGetDevices(devs *[]*addon.Device) {
	for _, d := range manager.devices {
		*devs = append(*devs, d)
	}
}

func (manager *AddonManager) handleGetThings(ts *[]*thing.Thing) {
	for _, d := range manager.devices {
		var t = asWebThing(d)
		*ts = append(*ts, t)
	}
}

func (manager *AddonManager) getAdapterByDeviceId(deviceId string) *AdapterProxy {
	device := manager.getDevice(deviceId)
	if device == nil {
		return nil
	}
	adapter, err1 := manager.findAdapter(device.AdapterId)
	if err1 != nil {
		log.Error(err1.Error())
		return nil
	}
	return adapter
}

func (manager *AddonManager) getDevice(deviceId string) *addon.Device {
	d, ok := manager.devices[deviceId]
	if !ok {
		return nil
	}
	return d
}

// get package from url, checksum
func (manager *AddonManager) installAddonFromUrl(id, url, checksum string, enabled bool) error {

	destPath := path.Join(os.TempDir(), id+".tar.gz")

	log.Info("fetching add-on %s as %s", url, destPath)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Download addon err,pakage id:%s err:%s", id, err.Error()))
	}
	defer func() {
		resp.Body.Close()
		err := os.Remove(destPath)
		if err != nil {
			log.Info("remove temp file failed ,err:%s", err.Error())
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
func (manager *AddonManager) installAddon(packageId, packagePath string, enabled bool) error {

	if !manager.addonsLoaded {
		return fmt.Errorf(`Cannot install add-on before other add-ons have been loaded.`)
	}
	log.Info("execute install package id: %s ", packageId)
	f, err := os.Open(packagePath)
	if err != nil {
		return err
	}
	defer f.Close()

	zp, _ := gzip.NewReader(f)
	tr := tar.NewReader(zp)

	for hdr, err := tr.Next(); err != io.EOF; hdr, err = tr.Next() {
		if err != nil {
			continue
		}
		// 读取文件信息
		fi := hdr.FileInfo()
		p := strings.Replace(hdr.Name, "package", packageId, 1)
		localPath := manager.AddonsDir + string(os.PathSeparator) + p
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
		var old AddonInfo
		ee := json.UnmarshalFromString(saved, &old)
		if ee != nil {
			old.Enabled = enabled
			newAddonInfo, err := json.MarshalToString(old)
			if err != nil {
				ee := database.SetSetting(key, newAddonInfo)
				if ee != nil {
					log.Error(ee.Error())
				}
			}

		}
	}
	if enabled {
		return manager.loadAddon(packageId)
	}
	return nil
}

func (manager *AddonManager) LoadAddons() {
	if manager.addonsLoaded {
		return
	}
	manager.addonsLoaded = true

	fs, err := os.ReadDir(manager.AddonsDir)
	if err != nil {
		log.Error("load addon err: %s", err.Error())
		return
	}
	manager.pluginServer = NewPluginServer(manager)

	for _, fi := range fs {
		if fi.IsDir() {
			addonId := fi.Name()
			err = manager.loadAddon(addonId)
			if err != nil {
				log.Error("load addon id:%s failed err:%s", addonId, err.Error())
			}
		}
	}
	return
}

func (manager *AddonManager) loadAddon(packageId string) error {

	if !manager.addonsLoaded {
		return nil
	}

	addonPath := path.Join(manager.AddonsDir, packageId)

	addonInfo, obj, err := loadManifest(addonPath, packageId)

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
				log.Error(eee.Error())
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

	err = util.EnsureDir(path.Join(manager.DataDir, addonInfo.ID))
	if err != nil {
		return err
	}

	manager.pluginServer.loadPlugin(addonPath, addonInfo.ID, addonInfo.Exec)

	return nil
}

func (manager *AddonManager) unloadAddon(packageId string) error {
	if !manager.addonsLoaded {
		return nil
	}
	_, ok := manager.extensions[packageId]
	if ok {
		delete(manager.extensions, packageId)
	}
	plugin, ok1 := manager.pluginServer.Plugins[packageId]
	if !ok1 {
		return fmt.Errorf("can not found plugin: %s", packageId)
	}

	for key, adapter := range manager.adapters {
		if adapter.PackageName == plugin.pluginId {
			for _, dev := range adapter.Devices {
				adapter.handleDeviceRemoved(dev)
			}
			delete(manager.adapters, key)
		}
	}
	manager.pluginServer.uninstallPlugin(packageId)
	return nil
}

func (manager *AddonManager) Start() {
	manager.pluginServer.Start()
	manager.running = true
}

func (manager *AddonManager) Close() {
	manager.pluginServer.Stop()
	manager.running = false
}
