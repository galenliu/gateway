package plugin

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"gateway/config"
	"gateway/pkg/bus"
	"gateway/pkg/log"
	"gateway/pkg/util"
	"gateway/server/models/thing"

	"addon"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
)

var manager *AddonManager

type ThingAdded func(*thing.Thing)
type ThingRemoved func(*thing.Thing)
type PropertyChanged func(property addon.Property)

type AddonManager struct {
	addon.AddonManager
	configPath   string
	pluginServer *PluginsServer

	devices       map[string]*addon.Device
	adapters      map[string]*AdapterProxy
	installAddons map[string]*AddonInfo

	isRunning bool
	isLoaded  bool

	isPairing bool

	//onDeviceAddedFuncs     map[interface{}]func(device *addon.Device)
	//onDeviceConnectedFuncs map[interface{}]func(device *addon.Device, connected bool)
	//onActionUpdateFuncs    map[interface{}]func(action *addon.Action)
	//onEventFuncs           map[interface{}]func(action *addon.Event)
	//onDeviceRemovedFuncs   map[interface{}]func(device *addon.Device)
	//onPropertyChangedFuncs map[interface{}]func(property *addon.Property)

	pluginCancel context.CancelFunc

	AddonsDir string
	DataDir   string
	ctx       context.Context
	locker    *sync.Mutex
}

func NewAddonsManager(ctx context.Context) (*AddonManager, error) {
	am := &AddonManager{}
	manager = am
	am.ctx = ctx
	am.AddonsDir = config.Conf.AddonsDir
	am.DataDir = config.Conf.DataDir

	//am.onDeviceAddedFuncs = make(map[interface{}]func(device *addon.Device))
	//am.onDeviceConnectedFuncs = make(map[interface{}]func(device *addon.Device, connected bool))
	//am.onActionUpdateFuncs = make(map[interface{}]func(action *addon.Action))
	//am.onEventFuncs = make(map[interface{}]func(action *addon.Event))
	//am.onDeviceRemovedFuncs = make(map[interface{}]func(device *addon.Device))
	//am.onPropertyChangedFuncs = make(map[interface{}]func(property *addon.Property))

	am.isRunning = false
	am.isPairing = false
	am.devices = make(map[string]*addon.Device, 50)
	am.installAddons = make(map[string]*AddonInfo, 50)
	am.adapters = make(map[string]*AdapterProxy, 20)

	var c context.Context
	c, am.pluginCancel = context.WithCancel(am.ctx)
	am.pluginServer = NewPluginServer(am, c)
	am.locker = new(sync.Mutex)
	err := am.LoadAddons()
	if err != nil {
		return nil, err
	}
	return am, nil
}

func (manager *AddonManager) LoadAddons() error {

	fs, err := os.ReadDir(manager.AddonsDir)
	if err != nil {
		return err
	}
	for _, fi := range fs {
		if fi.IsDir() {
			addonId := fi.Name()
			err = manager.loadAddon(addonId, true)
			if err != nil {
				log.Error(fmt.Sprintf("load add-ons: %v err:", addonId, addonId, err.Error()))
			}
		}
	}
	manager.isLoaded = true
	return nil
}

func (manager *AddonManager) loadAddon(packageId string, enabled bool) error {
	//db := database.GetDB()

	addonInfo, err := loadManifest(manager.AddonsDir, packageId)
	if err != nil {
		return err
	}

	err = util.EnsureDir(path.Join(manager.DataDir, addonInfo.ID))
	if err != nil {
		return err
	}
	err = addonInfo.UpdateFromDB()
	if err != nil {
		return err
	}
	manager.installAddons[packageId] = addonInfo
	manager.pluginServer.loadPlugin(addonInfo.ID, manager.installAddons[packageId].Exec, enabled)

	return nil
}

func (manager *AddonManager) unloadAddon(packageId string) error {

	return nil
}

func (manager *AddonManager) HandleDeviceAdded(device *addon.Device) {
	manager.devices[device.ID] = device
	if bus.HasCallBack(util.ThingAdded) {
		bus.Publish(util.ThingAdded, asThing(device))
	}
}

func (manager *AddonManager) HandleDeviceRemoved(device *addon.Device) {
	delete(manager.devices, device.ID)
	bus.Publish(util.ThingRemoved, asThing(device))
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

// get package from url, checksum
func (manager *AddonManager) installAddonFromUrl(id, url, checksum string, enabled bool) error {

	destPath := path.Join(os.TempDir(), id+".tar.gz")
	log.Info(fmt.Sprintf("fetching add-on %s as %s", url, destPath))
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return fmt.Errorf(fmt.Sprintf("Download addon err,pakage id:%s err:%s", id, err.Error()))
	}
	defer resp.Body.Close()
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

	log.Info(fmt.Sprintf("start instll package id: %s ", packageId))
	f, _ := os.Open(packagePath)
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
	return manager.loadAddon(packageId, enabled)

}

func (manager *AddonManager) handleSetPropertyValue(deviceId, propName string, setValue interface{}) {
	adapter := manager.getAdapterByDeviceId(deviceId)
	property := adapter.GetDevice(deviceId).GetProperty(propName)
	adapter.setPropertyValue(property, setValue)
}

func (manager *AddonManager) getAdapterByDeviceId(deviceId string) *AdapterProxy {
	device, err := manager.findDevice(deviceId)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	adapter, err1 := manager.findAdapter(device.AdapterId)
	if err1 != nil {
		log.Error(err1.Error())
		return nil
	}
	return adapter
}

func (manager *AddonManager) findDevice(deviceId string) (*addon.Device, error) {
	if deviceId == "" {
		return nil, fmt.Errorf("devices id none")
	}
	d, ok := manager.devices[deviceId]
	if !ok {
		return nil, fmt.Errorf("addon manager can not found devices(id:%s)", deviceId)
	}
	return d, nil
}

func (manager *AddonManager) Start() {
	go manager.pluginServer.Start()
	manager.isRunning = true
}

func (manager *AddonManager) Stop() {
	//停止
	manager.pluginServer.Stop()
	manager.isRunning = false
}
