package plugin

//	plugin server
import (
	"context"
	ipc "github.com/galenliu/gateway/pkg/ipc_server"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/rpc_server"
	"sync"
)

type PluginsServer struct {
	Plugins   map[string]*Plugin
	locker    *sync.Mutex
	manager   *Manager
	ipc       *ipc.IPCServer
	rpc       *rpc_server.RPCServer
	closeChan chan struct{}
	logger    logging.Logger
	ctx       context.Context
	verbose   bool
}

func NewPluginServer(manager *Manager, userProfile []byte, preferences []byte) *PluginsServer {
	server := &PluginsServer{}
	server.logger = manager.logger
	server.closeChan = make(chan struct{})
	server.Plugins = make(map[string]*Plugin, 30)
	server.locker = new(sync.Mutex)
	server.manager = manager
	server.ipc = ipc.NewIPCServer(server, manager.options.IPCPort, userProfile, preferences, manager.logger)
	server.rpc = rpc_server.NewRPCServer(server, manager.options.RPCPort, userProfile, preferences, manager.logger)
	return server
}

func (s *PluginsServer) RegisterPlugin(pluginId string, clint IClint) *Plugin {
	plugin, ok := s.Plugins[pluginId]
	if !ok {
		plugin = NewPlugin(s, pluginId, s.logger)
		s.Plugins[pluginId] = plugin
	}
	plugin.Clint = clint
	return plugin
}

func (s *PluginsServer) loadPlugin(addonPath, id, exec string) {
	plugin := s.registerPlugin(id)
	plugin.exec = exec
	plugin.execPath = addonPath
	go plugin.execute()
}

func (s *PluginsServer) uninstallPlugin(packageId string) {
	plugin := s.Plugins[packageId]
	if plugin == nil {
		s.logger.Error("plugin not exist")
		return
	}
	plugin.unload()
	delete(s.Plugins, packageId)

}

func (s *PluginsServer) registerPlugin(packageId string) *Plugin {
	plugin, ok := s.Plugins[packageId]
	if !ok {
		plugin = NewPlugin(s, packageId, s.logger)
		s.Plugins[packageId] = plugin
	}
	return plugin
}

// Start create goroutines handle ipc massage
func (s *PluginsServer) Start() error {

	err := s.rpc.Start()
	if err != nil {
		s.logger.Errorf("rpc server start err:%s", err.Error())
		return err
	}
	err = s.ipc.Start()
	if err != nil {
		s.logger.Errorf("ipc server start err:%s", err.Error())
		return err
	}
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
	for _, p := range s.Plugins {
		p.unload()
	}
	return nil
}
