package plugin

//	plugin server
import (
	"github.com/galenliu/gateway/pkg/constant"
	ipc "github.com/galenliu/gateway/pkg/ipc_server"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/rpc"
	"github.com/galenliu/gateway/pkg/rpc_server"
	"sync"
)

type PluginsServer struct {
	Plugins   sync.Map
	manager   *Manager
	ipc       *ipc.IPCServer
	rpc       *rpc_server.RPCServer
	closeChan chan struct{}
	logger    logging.Logger
}

func NewPluginServer(manager *Manager) *PluginsServer {
	server := &PluginsServer{}
	server.logger = manager.logger
	server.closeChan = make(chan struct{})
	server.manager = manager
	u := &rpc.PluginRegisterResponseMessage_Data_UsrProfile{
		AddonsDir:  manager.config.UserProfile.AddonsDir,
		ConfigDir:  manager.config.UserProfile.ConfigDir,
		DataDir:    manager.config.UserProfile.DataDir,
		MediaDir:   manager.config.UserProfile.MediaDir,
		LogDir:     manager.config.UserProfile.LogDir,
		GatewayDir: constant.Version,
	}
	p := &rpc.PluginRegisterResponseMessage_Data_Preferences{
		Language: manager.config.Preferences.Language,
		Units:    &rpc.PluginRegisterResponseMessage_Data_Preferences_Units{Temperature: manager.config.Preferences.Units.Temperature},
	}
	server.ipc = ipc.NewIPCServer(server, manager.config.IPCPort, u, p, manager.logger)
	server.rpc = rpc_server.NewRPCServer(server, manager.config.RPCPort, u, p, manager.logger)
	return server
}

func (s *PluginsServer) RegisterPlugin(pluginId string, clint rpc.Clint) rpc.PluginHandler {
	plugin := s.getPlugin(pluginId)
	if plugin == nil {
		plugin = NewPlugin(s, pluginId, s.logger)
		s.Plugins.Store(pluginId, plugin)
	}
	plugin.Clint = clint
	return plugin
}

func (s *PluginsServer) loadPlugin(addonPath, id, exec string) {
	_ = s.RegisterPlugin(id, nil)
	//plugin.exec = exec
	//plugin.execPath = addonPath
	//go plugin.execute()
}

func (s *PluginsServer) unloadPlugin(packageId string) {
	plugin := s.getPlugin(packageId)
	if plugin == nil {
		s.logger.Error("plugin not exist")
		return
	}
	plugin.unload()
	s.Plugins.Delete(packageId)
}

func (s *PluginsServer) getPlugin(id string) *Plugin {
	p, ok := s.Plugins.Load(id)
	plugin, ok := p.(*Plugin)
	if !ok {
		return nil
	}
	return plugin
}

func (s *PluginsServer) getPlugins() (plugins []*Plugin) {
	s.Plugins.Range(func(key, value interface{}) bool {
		p, ok := value.(*Plugin)
		if ok {
			plugins = append(plugins, p)
		}
		return true
	})
	return
}

// Start create goroutines handle ipc massage
func (s *PluginsServer) Start() error {
	_ = s.rpc.Start()
	_ = s.ipc.Start()
	return nil
}

// Stop if server stop, also need to stop all of package
func (s *PluginsServer) Stop() error {
	err := s.ipc.Stop()
	err = s.rpc.Stop()
	if err != nil {
		return err
	}
	s.closeChan <- struct{}{}
	s.logger.Info("Plugin server stopped")
	for _, p := range s.getPlugins() {
		p.unload()
	}
	return nil
}
