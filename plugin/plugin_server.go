package plugin

//	plugin server
import (
	rpc "github.com/galenliu/gateway-grpc"
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
	s.ipc = ipc.NewIPCServer(s, manager.config.IPCPort, manager.config.UserProfile, s.logger)
	s.rpc = rpc_server.NewRPCServer(s, manager.config.RPCPort, manager.config.UserProfile, manager.logger)
	s.Start()
	return s
}

func (s *PluginsServer) RegisterPlugin(pluginId string, clint server.Clint) server.PluginHandler {
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

func (s *PluginsServer) unregisterPlugin(id string) {
	s.Plugins.Delete(id)
}

// loadPlugin
//  @Description:
//  @receiver s
//  @param addonPath   package所以的目录
//  @param id
//  @param exec
func (s *PluginsServer) loadPlugin(pluginId, packagePath, exec string) {
	plugin := s.registerPlugin(pluginId)
	plugin.exec = exec
	plugin.execPath = packagePath
	plugin.start()
}

func (s *PluginsServer) registerPlugin(pluginId string) *Plugin {
	plugin := s.getPlugin(pluginId)
	if plugin == nil {
		plugin = NewPlugin(pluginId, s.manager, s, s.logger)
	}
	s.Plugins.Store(pluginId, plugin)
	return plugin
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
func (s *PluginsServer) Start() {
	go func() {
		err := s.rpc.Run()
		if err != nil {
			s.logger.Errorf("rpc server err:", err.Error())
		}
	}()
	go func() {
		err := s.ipc.Run()
		if err != nil {
			s.logger.Errorf("ipc server err:", err.Error())
		}
	}()
	return
}

func (s *PluginsServer) GetPreferences() *rpc.Preferences {
	r := &rpc.Preferences{Language: "en-US", Units: &rpc.Preferences_Units{Temperature: "degree celsius"}}
	lang, err := s.manager.storage.GetSetting("localization.language")
	if err != nil {
		r.Language = lang
	}
	temp, err := s.manager.storage.GetSetting("llocalization.units.temperature")
	if err != nil {
		r.Units.Temperature = temp
	}
	return r
}

func (s *PluginsServer) shutdown() {

}
