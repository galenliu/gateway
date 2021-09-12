package plugin

//	plugin server
import (
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/rpc"
	ipc "github.com/galenliu/gateway/pkg/rpc/ipc_server"
	"github.com/galenliu/gateway/pkg/rpc/rpc_server"
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
	server.ipc = ipc.NewIPCServer(server, manager.config.IPCPort, manager.config.UserProfile, manager.config.Preferences, manager.logger)
	server.rpc = rpc_server.NewRPCServer(server, manager.config.RPCPort, manager.config.UserProfile, manager.config.Preferences, manager.logger)
	return server
}

func (s *PluginsServer) RegisterPlugin(pluginId string, clint rpc.Clint) rpc.PluginHandler {
	plugin := s.getPlugin(pluginId)
	if plugin == nil {
		plugin = NewPlugin(pluginId, s.manager, s, s.logger)
		s.Plugins.Store(pluginId, plugin)
	}
	if clint != nil {
		plugin.Clint = clint
	}
	return plugin
}

func (s *PluginsServer) loadPlugin(addonPath, id, exec string) {
	plugin := s.getPlugin(id)
	if plugin == nil {
		plugin = NewPlugin(id, s.manager, s, s.logger)
		s.Plugins.Store(id, plugin)
	}
	plugin.exec = exec
	plugin.execPath = addonPath
	plugin.start()
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
