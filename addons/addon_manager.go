package addons

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"gateway/config"
	"gateway/event"
	"gateway/pkg/log"
	"gateway/pkg/util"
	json "github.com/json-iterator/go"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
)

var addonsManager *AddonManager

type IEvent interface {
	OnPropertyChanged(thingId, propName string, value interface{})
}

type AddonManager struct {
	configPath   string
	pluginServer *PluginsServer

	devices       map[string]*DeviceProxy
	adapters      map[string]*AdapterProxy
	installAddons map[string]*AddonInfo

	IsRunning bool
	IsLoaded  bool

	pluginCancel context.CancelFunc

	EventBus IEvent

	AddonsDir string
	DataDir   string
	ctx       context.Context
	locker    *sync.Mutex
}

func NewAddonsManager(ctx context.Context) (*AddonManager, error) {
	am := &AddonManager{}
	addonsManager = am
	am.ctx = ctx
	am.AddonsDir = config.Conf.AddonsDir
	am.DataDir = config.Conf.DataDir

	am.IsRunning = false
	am.devices = make(map[string]*DeviceProxy, 50)
	am.installAddons = make(map[string]*AddonInfo, 50)
	am.adapters = make(map[string]*AdapterProxy)

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

	fs, err := ioutil.ReadDir(manager.AddonsDir)
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
	manager.IsLoaded = true
	return nil
}

func (manager *AddonManager) loadAddon(packageId string, enabled bool) error {
	//db := database.GetDB()

	addonInfo, err := LoadManifest(manager.AddonsDir, packageId)
	if err != nil {
		return err
	}

	err = util.EnsureDir(path.Join(manager.DataDir, addonInfo.ID))
	if err != nil {
		return err
	}
	err = addonInfo.UpdateOrCreateFormDb()
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

func (manager *AddonManager) handlerDeviceAdded(dev *DeviceProxy) {
	if dev.ID != "" {
		manager.devices[dev.ID] = dev
	}
	event.FireDiscoverNewDevice(dev.Device)
}

func (manager *AddonManager) addAdapter(proxy *AdapterProxy) {
	manager.locker.Lock()
	defer manager.locker.Unlock()
	manager.adapters[proxy.ID] = proxy
	event.FireAdapterAdded(proxy)
}

// get package from url, checksum
func (manager *AddonManager) InstallAddonFromUrl(id, url, checksum string, enabled bool) error {

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

func (manager *AddonManager) GetInstallAddons() map[string]interface{} {
	return nil
}

func (manager *AddonManager) Start() {
	go manager.pluginServer.Start()
	manager.IsRunning = true
}

func (manager *AddonManager) Stop() {
	//停止
	manager.pluginServer.Stop()
	manager.IsRunning = false
}

func GetInstallAddons() ([]byte, error) {
	addonsManager.locker.Lock()
	defer addonsManager.locker.Unlock()
	var addons []*AddonInfo
	for _, v := range addonsManager.installAddons {
		addons = append(addons, v)
	}
	data, err := json.Marshal(addons)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func EnableAddon(addonId string) error {
	addonsManager.locker.Lock()
	defer addonsManager.locker.Unlock()
	addonInfo := addonsManager.installAddons[addonId]
	if addonInfo == nil {
		return fmt.Errorf("addon not exit")
	}

	err := addonInfo.UpdateAddonInfoToDB(true)
	err = addonsManager.loadAddon(addonId, true)
	if err != nil {
		return err
	}
	return nil
}

func DisableAddon(addonId string) error {
	addonsManager.locker.Lock()
	defer addonsManager.locker.Unlock()
	addonInfo := addonsManager.installAddons[addonId]
	if addonInfo == nil {
		return fmt.Errorf("addon not installed")
	}
	err := addonInfo.UpdateAddonInfoToDB(false)
	err = addonsManager.unloadAddon(addonId)
	if err != nil {
		return err
	}
	return nil
}

func GetThings() []*DeviceProxy {
	var devs = make([]*DeviceProxy, len(addonsManager.devices))
	for _, proxy := range addonsManager.devices {
		devs = append(devs, proxy)
	}
	return devs
}

func InstallAddonFromUrl(id, url, checksum string, enabled bool) {
	_ = addonsManager.InstallAddonFromUrl(id, url, checksum, enabled)
}

func SetDeviceProperty(thingId, propName string, value interface{}) (interface{}, error) {
	device, ok := addonsManager.devices[thingId]
	if !ok {
		return value, fmt.Errorf("invalid thingId")
	}
	return device.SetProperty(propName, value)
}

func AddNewThing(pairingTimeout int) error {
	for _, plugin := range addonsManager.pluginServer.Plugins {
		for _, adapter := range plugin.adapters {
			adapter.Pairing(pairingTimeout)
		}
	}
	return nil
}
