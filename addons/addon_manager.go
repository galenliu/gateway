package addons

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"gateway/util"
	"gateway/util/logger"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)



var log logger.Logger

type ManagerConfig struct {
	AddonsDir string
	DataDir string

}


type Manager struct {
	configPath    string
	pluginServer  *PluginsServer
	devices       map[string]*DeviceProxy
	installAddons map[string]*AddonManifest // {addonId: manifest}



	IsRunning bool
	IsLoaded  bool

	AddonsDir string
	DataDir  string
	ctx       context.Context
}

func NewAddonsManager(config ManagerConfig) *Manager {
	am := &Manager{}
	am.ctx = context.Background()
	am.AddonsDir = config.AddonsDir
	am.DataDir = config.DataDir


	am.IsLoaded = false
	am.IsRunning = false
	am.devices = make(map[string]*DeviceProxy)
	am.installAddons = make(map[string]*AddonManifest)
	am.pluginServer = NewPluginServer(am, am.ctx)
	log = logger.GetLog()
	return am
}

func (manager *Manager) LoadAddons() {
	if manager.IsLoaded {
		log.Info("addon manager is loaded")
		return
	}
	fs, err := ioutil.ReadDir(manager.AddonsDir)
	if err != nil {
		log.Error("read addon dir err:", zap.Error(err))
	}
	for _, fi := range fs {
		if fi.IsDir() {
			addonId := fi.Name()
			err = manager.loadAddon(addonId)
			if err != nil {
				log.Error(fmt.Sprintf("load add-ons: %v err:", addonId), zap.Error(err))
			}
		}
	}
}

func (manager *Manager) loadAddon(packageId string) error {
	//db := database.GetDB()

	manifest, err := LoadManifest(manager.AddonsDir, packageId)
	if err != nil {
		return err
	}

	////get saved form db,
	//var saveManifest *AddonManifest
	//db.First(&saveManifest, packageId)
	//
	////update
	//if saveManifest != nil {
	//	manifest.Enable = saveManifest.Enable
	//}
	////saved
	//db.Create(manifest)s
	//
	//if !manifest.Enable {
	//	err = fmt.Errorf("add-on is not enabled:%v", manifest.ID)
	//	return err
	//}
	//if manifest.GatewaySpecificSettings["webthings"].Exec == "" {
	//	err = fmt.Errorf("add-on exec nil:%v", manifest.ID)
	//	return err
	//}

	//create addon date dir
	err = util.EnsureDir(path.Join(manager.DataDir, manifest.ID))
	if err != nil {
		return err
	}

	manager.installAddons[packageId] = manifest

	manager.pluginServer.loadPlugin(manifest.ID, manifest.GatewaySpecificSettings["webthings"].Exec)

	//plugin := manager.pluginServer.getPlugin(manifest.ID)
	//plugin.exec = manifest.GatewaySpecificSettings.Exec
	//plugin.start()
	return nil
}

// get package from url, checksum
func (manager *Manager) InstallAddonFromUrl(id, url, checksum string, enable bool) error {

	destPath := path.Join(os.TempDir(), id+".tar.gz")
	log.Info(fmt.Sprintf("fetching add-on %s as %s", url, destPath))
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		log.Error("Download addon err:", zap.String("HTTP States:", resp.Status), zap.String("url", url))
		return err
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	_ = ioutil.WriteFile(destPath, data, os.ModePerm)

	manager.installAddon(id, destPath, enable)
	return nil
}

//tar package to addon from the temp dir,
func (manager *Manager) installAddon(packageId, packagePath string, enable bool) {

	if !manager.IsLoaded {
		log.Warn(fmt.Sprintf("Cannot install add-on:%v before other add-ons have been loaded", packageId))
	}
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
			fmt.Print("mkdir:" + p)
			_ = os.MkdirAll(localPath, os.ModePerm)
			continue
		}
		// 创建一个空文件，用来写入解包后的数据
		fw, err := os.Create(localPath)
		if err != nil {
			log.Error("create file err")
			continue
		}
		if _, err := io.Copy(fw, tr); err != nil {
			log.Error("create file err")
		}
		//TODO 给下载的文件增加可执行权限
		_ = os.Chmod(fi.Name(), 777)
		_ = fw.Close()
	}
	err := manager.loadAddon(packageId)
	manager.IsLoaded = false

	if err != nil {
		log.Warn(fmt.Sprintf("install add-ons:%v happand err:", packageId), zap.Error(err))
	}
}

func (manager *Manager) GetProperty(thingId, propName string) interface{} {
	dev := manager.getDevice(thingId)
	return dev.GetProperty(propName)
}

func (manager *Manager) SetProperty(thingId, propName string, value interface{}) {
	dev := manager.getDevice(thingId)
	prop := dev.FindProperty(propName)
	prop.setValue(value)

}

func (manager *Manager) RemoveThing(thingId string) {
	dev := manager.getDevice(thingId)
	adapter := dev.GetAdapter()
	adapter.removeThing(dev)
}

func (manager *Manager) GetInstallAddons() map[string]interface{} {
	return nil
}

func (manager *Manager) getDevice(thingId string) *DeviceProxy {
	if len(manager.devices) > 0 {
		dev := manager.devices[thingId]
		if dev != nil {
			return dev
		}
	}
	return nil
}

func (manager *Manager) Start() {

	manager.IsRunning = true
}

func (manager *Manager) Stop() {
	//停止
	manager.pluginServer.close()
	manager.IsRunning = false
}
