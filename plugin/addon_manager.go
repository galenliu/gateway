package plugin

import (
	"context"
	messages "github.com/galeuliu/gateway-schema"
	"go.uber.org/zap"
	"io/ioutil"
	"path"
)

var log *zap.Logger

type AddonsManager struct {
	userProfile  messages.UserProfile
	preferences  messages.Preferences
	configPath   string
	pluginServer *PluginsServer
	devices      map[string]*DeviceProxy

	isLoaded     bool
	addonPath    string
	databasePath string
}

func NewAddonsManager(u messages.UserProfile, p messages.Preferences, _log *zap.Logger) *AddonsManager {

	ctx := context.Background()
	am := &AddonsManager{}
	am.isLoaded = false
	am.devices = make(map[string]*DeviceProxy)
	am.userProfile = u
	am.preferences = p
	am.addonPath = am.userProfile.AddonsDir
	am.pluginServer = NewPluginServer(am, ctx)
	log = _log
	return am

}

func (manager *AddonsManager) LoadAddons() {
	if manager.isLoaded {
		return
	}
	fs, err := ioutil.ReadDir(manager.addonPath)
	if err != nil {
		return
	}
	for _, fi := range fs {
		if fi.IsDir() {
			manager.loadAddon(fi.Name())
		}

	}

}

func (manager *AddonsManager) loadAddon(PluginDirName string) {
	manifest, err := LoadManifest(manager.addonPath, PluginDirName)
	if err != nil {
		log.Error("load manifest fail", zap.String("plugin name", PluginDirName), zap.Error(err))
		return
	}
	plugin := manager.pluginServer.getPlugin(manifest.ID)
	plugin.exec = manifest.Settings.Exec
	plugin.execPath = path.Join(manager.addonPath, PluginDirName)
	plugin.run()
}

func (manager *AddonsManager) Run() {
	//启动plugin server,接受来自plugin的请求
	manager.pluginServer.Run()
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

func (manager *AddonsManager) getDevice(thingId string) *DeviceProxy {
	if len(manager.devices) > 0 {
		dev := manager.devices[thingId]
		if dev != nil {
			return dev
		}
	}
	return nil
}
