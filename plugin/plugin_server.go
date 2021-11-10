package plugin

//	plugin server
import (
	"context"
	"fmt"
	rpc "github.com/galenliu/gateway-grpc"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/ipc"
	"github.com/galenliu/gateway/pkg/logging"
	json "github.com/json-iterator/go"
	"sync"
)

type PluginsServer struct {
	Plugins   sync.Map
	manager   *Manager
	ipc       *ipc.IPCServer
	rpc       *ipc.RPCServer
	closeChan chan struct{}
	ctx       context.Context
	logger    logging.Logger
}

func NewPluginServer(manager *Manager) *PluginsServer {
	s := &PluginsServer{}

	s.logger = manager.logger
	s.closeChan = make(chan struct{})
	s.manager = manager
	s.ipc = ipc.NewIPCServer(s, manager.config.IPCPort, manager.config.UserProfile, s.logger)
	//s.rpc = ipc.NewRPCServer(ctx, s, manager.config.RPCPort, manager.config.UserProfile, manager.logger)
	return s
}

func (s *PluginsServer) RegisterPlugin(clint ipc.Clint) (ipc.PluginHandler, error) {
	message, err := clint.ReadMessage()
	if err != nil {
		return nil, err
	}
	if message.MessageType != rpc.MessageType_PluginRegisterRequest {
		return nil, fmt.Errorf("MessageType need PluginRegisterRequest")
	}
	var registerMessage rpc.PluginRegisterRequestMessage_Data
	err = json.Unmarshal(message.Data, &registerMessage)
	if err != nil {
		return nil, fmt.Errorf("message failed err: %s", err.Error())
	}

	responseMessage :=
		&rpc.PluginRegisterResponseMessage_Data{
			PluginId:       json.Get(message.Data, "pluginId").ToString(),
			GatewayVersion: constant.Version,
			UserProfile:    s.manager.config.UserProfile,
			Preferences:    s.getPreferences(),
		}
	data, _ := json.Marshal(responseMessage)
	clint.SetPluginId(registerMessage.PluginId)
	err = clint.WriteMessage(&rpc.BaseMessage{MessageType: rpc.MessageType_PluginRegisterResponse, Data: data})
	if err != nil {
		return nil, err
	}

	plugin := s.registerPlugin(registerMessage.PluginId)
	plugin.clint = clint
	return plugin, nil
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
	go plugin.run()
}

func (s *PluginsServer) registerPlugin(pluginId string) *Plugin {
	plugin := s.getPlugin(pluginId)
	if plugin == nil {
		plugin = NewPlugin(pluginId, s.manager, s, s.logger)
	} else {
		return plugin
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

func (s *PluginsServer) getPreferences() *rpc.Preferences {
	r := &rpc.Preferences{Language: "en-US", Units: &rpc.Preferences_Units{Temperature: "degree celsius"}}
	lang, err := s.manager.storage.GetSetting("localization.language")
	if err == nil {
		r.Language = lang
	}
	temp, err := s.manager.storage.GetSetting("localization.units.temperature")
	if err == nil {
		r.Units.Temperature = temp
	}
	return r
}

func (s *PluginsServer) shutdown() {

}
