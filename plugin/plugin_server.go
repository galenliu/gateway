package plugin

//	plugin server
import (
	"fmt"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/server"
	ipc "github.com/galenliu/gateway/pkg/server/ipc_server"
	"github.com/galenliu/gateway/pkg/server/rpc_server"
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
	s := &PluginsServer{}
	s.logger = manager.logger
	s.closeChan = make(chan struct{})
	s.manager = manager
	s.ipc = ipc.NewIPCServer(s, manager.config.IPCPort, manager.config.UserProfile, manager.config.Preferences, manager.logger)
	s.rpc = rpc_server.NewRPCServer(s, manager.config.RPCPort, manager.config.UserProfile, manager.config.Preferences, manager.logger)
	return s
}

func (s *PluginsServer) RegisterPlugin(pluginId string, clint server.Clint) server.PluginHandler {
	plugin := s.findPlugin(pluginId)
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
	plugin := s.findPlugin(id)
	if plugin == nil {
		plugin = NewPlugin(id, s.manager, s, s.logger)
		s.Plugins.Store(id, plugin)
	}
	plugin.exec = exec
	plugin.execPath = addonPath
	plugin.start()
}

func (s *PluginsServer) unloadPlugin(packageId string) error {
	plugin := s.findPlugin(packageId)
	if plugin == nil {
		return fmt.Errorf("plugin not exist")
	}
	plugin.unload()
	s.Plugins.Delete(packageId)
	return nil
}

func (s *PluginsServer) findPlugin(id string) *Plugin {
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

func (s *PluginsServer)shutdown(){
	_=s.ipc.Stop()
	_=s.rpc.Stop()
}


func (s *PluginsServer) unregisterPlugin(id string) {
	s.Plugins.Delete(id)
}
