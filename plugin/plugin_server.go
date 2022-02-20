package plugin

//	plugin api
import (
	"encoding/json"
	"fmt"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/ipc"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"github.com/galenliu/gateway/pkg/logging"
	"sync"
)

type PluginsServer struct {
	Plugins sync.Map
	manager *Manager
	ipc     *ipc.WebSocketServer
	logger  logging.Logger
}

func NewPluginServer(manager *Manager) *PluginsServer {
	s := &PluginsServer{}
	s.logger = manager.logger
	s.manager = manager
	s.ipc = ipc.NewIPCServer(s, manager.config.IPCPort, manager.config.UserProfile, s.logger)
	return s
}

func (s *PluginsServer) RegisterPlugin(connection ipc.Connection) (ipc.PluginHandler, error) {
	m, err := connection.ReadMessage()
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(m.Data)
	if m.MessageType != messages.MessageType_PluginRegisterRequest {
		return nil, fmt.Errorf("MessageType need PluginRegisterRequest")
	}
	var registerMessage messages.PluginRegisterRequestJsonData
	err = json.Unmarshal(data, &registerMessage)
	if err != nil {
		return nil, err
	}
	responseData := messages.PluginRegisterResponseJsonData{
		GatewayVersion: constant.Version,
		PluginId:       registerMessage.PluginId,
		Preferences:    *s.manager.GetPreferences(),
		UserProfile:    *s.manager.config.UserProfile,
	}
	err = connection.WriteMessage(messages.MessageType_PluginRegisterResponse, responseData)
	if err != nil {
		return nil, err
	}
	plugin := s.registerPlugin(registerMessage.PluginId)
	plugin.register(connection)
	return plugin, nil
}

func (s *PluginsServer) unregisterPlugin(id string) {
	s.Plugins.Delete(id)
}

func (s *PluginsServer) loadPlugin(pluginId, packagePath, exec string) {
	plugin := s.registerPlugin(pluginId)
	plugin.exec = exec
	plugin.execPath = packagePath
	go plugin.start()
}

func (s *PluginsServer) registerPlugin(pluginId string) *Plugin {
	plugin := s.getPlugin(pluginId)
	if plugin != nil {
		return plugin
	}
	plugin = NewPlugin(pluginId, s.manager, s.logger)
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
	s.Plugins.Range(func(key, value any) bool {
		p, ok := value.(*Plugin)
		if ok {
			plugins = append(plugins, p)
		}
		return true
	})
	return
}

func (s *PluginsServer) shutdown() {

}
