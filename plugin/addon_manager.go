package plugin

import (
	"context"
	messages "github.com/galeuliu/gateway-schema"
	"go.uber.org/zap"
	"io/ioutil"
	"path"
)

var log *zap.Logger

type IGateway interface {
	GetUserProfile() *messages.UserProfile
	GetPreferences() *messages.Preferences
	EnsureConfigPath(dir string, dirs ...string)
}

type AddonsManager struct {
	userProfile   *messages.UserProfile
	preferences   *messages.Preferences
	configPath    string
	pluginServer  *PluginsServer
	devices       map[string]*DeviceProxy
	installAddons map[string]interface{} // {addonId: manifest}

	isLoaded     bool
	addonPath    string
	databasePath string

	iGateway IGateway
}

func NewAddonsManager(gateway IGateway, _log *zap.Logger) *AddonsManager {

	ctx := context.Background()
	am := &AddonsManager{}
	am.iGateway = gateway
	am.isLoaded = false
	am.devices = make(map[string]*DeviceProxy)
	am.userProfile = gateway.GetUserProfile()
	am.preferences = gateway.GetPreferences()
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
	addonConfig, err := LoadManifest(manager.addonPath, PluginDirName)
	if err != nil {
		log.Error("load manifest fail", zap.String("plugin name", PluginDirName), zap.Error(err))
		return
	}
	manager.iGateway.EnsureConfigPath(path.Join(manager.userProfile.DataDir, addonConfig.ID))
	plugin := manager.pluginServer.getPlugin(addonConfig.ID)
	plugin.exec = addonConfig.Exec
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

func (manager *AddonsManager) RemoveThing(thingId string) {
	dev := manager.getDevice(thingId)
	adapter := dev.GetAdapter()
	adapter.removeThing(dev)
}

func (manager *AddonsManager) GetInstallAddons() map[string]interface{} {
	return manager.installAddons
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
