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

	isLoaded     bool
	addonPath    string
	databasePath string
}

func NewAddonsManager(u messages.UserProfile, p messages.Preferences, _log *zap.Logger) *AddonsManager {

	ctx := context.Background()
	am := &AddonsManager{}
	am.isLoaded = false
	am.userProfile = u
	am.preferences = p
	am.addonPath = am.userProfile.AddonsDir
	am.pluginServer = NewPluginServer(am, ctx)
	log = _log
	return am

}

func (am *AddonsManager) LoadAddons() {
	if am.isLoaded {
		return
	}
	fs, err := ioutil.ReadDir(am.addonPath)
	if err != nil {
		return
	}
	for _, fi := range fs {
		if fi.IsDir() {
			am.loadAddon(fi.Name())
		}

	}

}

func (am *AddonsManager) loadAddon(PluginDirName string) {
	manifest, err := LoadManifest(am.addonPath, PluginDirName)
	if err != nil {
		log.Error("load manifest fail", zap.String("plugin name", PluginDirName), zap.Error(err))
		return
	}
	plugin := am.pluginServer.getPlugin(manifest.ID)
	plugin.exec = manifest.Settings.Exec
	plugin.execPath = path.Join(am.addonPath, PluginDirName)
	plugin.run()
}

func (am *AddonsManager) Run() {
	//启动plugin server,接受来自plugin的请求
	am.pluginServer.Run()
}
