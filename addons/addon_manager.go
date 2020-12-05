package addons

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"gateway/pkg/logger"
	"gateway/pkg/runtime"
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

var addonsManager *AddonsManager
var log logger.Logger

type AddonsManager struct {
	configPath    string
	pluginServer  *PluginsServer
	devices       map[string]*DeviceProxy
	installAddons map[string]*AddonInfo // {addonId: manifest}

	IsRunning bool
	IsLoaded  bool

	pluginCancel context.CancelFunc

	AddonsDir string
	DataDir   string
	ctx       context.Context
	locker    *sync.Mutex
}

func NewAddonsManager(ctx context.Context) (*AddonsManager, error) {
	log = logger.GetLog()
	am := &AddonsManager{}
	addonsManager = am
	am.ctx = ctx
	am.AddonsDir = runtime.RuntimeConf.AddonsDir
	am.DataDir = runtime.RuntimeConf.DataDir

	am.IsRunning = false
	am.devices = make(map[string]*DeviceProxy)
	am.installAddons = make(map[string]*AddonInfo)
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

func (manager *AddonsManager) LoadAddons() error {

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

func (manager *AddonsManager) loadAddon(packageId string, enabled bool) error {
	//db := database.GetDB()

	manifest, err := LoadManifest(manager.AddonsDir, packageId)
	if err != nil {
		return err
	}

	err = util.EnsureDir(path.Join(manager.DataDir, manifest.ID))
	if err != nil {
		return err
	}
	saveAddonInfo := GetAddonsInfoByIDFromDB(packageId)
	if saveAddonInfo != nil {
		manifest.Enable = enabled
	}
	addonInfo := SaveAddonManifestToDB(manifest, true)
	manager.installAddons[packageId] = addonInfo
	manager.pluginServer.loadPlugin(addonInfo.ID, manager.installAddons[packageId].Exec, enabled)

	return nil
}

func (manager *AddonsManager) unloadAddon(packageId string) error {

	return nil
}

// get package from url, checksum
func (manager *AddonsManager) InstallAddonFromUrl(id, url, checksum string, enabled bool) error {

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
	manager.installAddon(id, destPath, enabled)
	return nil
}

//tar package to addon from the temp dir,
func (manager *AddonsManager) installAddon(packageId, packagePath string, enabled bool) error {

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

func (manager *AddonsManager) GetProperty(thingId, propName string) interface{} {
	dev := manager.getDevice(thingId)
	return dev.GetProperty(propName)
}

func (manager *AddonsManager) SetProperty(thingId, propName string, value interface{}) {
	dev := manager.getDevice(thingId)
	prop := dev.FindProperty(propName)
	prop.setValue(value)

}

func (manager *AddonsManager) RemoveThing(thingId string) {
	dev := manager.getDevice(thingId)
	adapter := dev.GetAdapter()
	adapter.removeThing(dev)
}

func (manager *AddonsManager) GetInstallAddons() map[string]interface{} {
	return nil
}

func (manager *AddonsManager) getDevice(thingId string) *DeviceProxy {
	if len(manager.devices) > 0 {
		dev := manager.devices[thingId]
		if dev != nil {
			return dev
		}
	}
	return nil
}

func (manager *AddonsManager) Start() {

	manager.IsRunning = true
}

func (manager *AddonsManager) Stop() {
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
	addonInfo.Enabled = true
	SaveAddonInfo(*addonInfo)
	err := addonsManager.loadAddon(addonId, true)
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
	addonInfo.Enabled = false
	SaveAddonInfo(*addonInfo)
	err := addonsManager.unloadAddon(addonId)
	if err != nil {
		return err
	}
	return nil
}

func InstallAddonFromUrl(id, url, checksum string, enable bool) {
	_ = addonsManager.InstallAddonFromUrl(id, url, checksum, enable)
}
