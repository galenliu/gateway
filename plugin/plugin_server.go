package plugin

//	plugin api
import (
	"context"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/gofiber/websocket/v2"
	json "github.com/json-iterator/go"
	"sync"
)

type PluginsServer struct {
	Plugins sync.Map
	manager *Manager
	logger  logging.Logger
}

func NewPluginServer(manager *Manager) *PluginsServer {
	ctx := context.TODO()
	s := &PluginsServer{}
	s.logger = manager.logger
	s.manager = manager
	//s.ipc = ipc.NewIPCServer(s, manager.config.IPCPort, manager.config.UserProfile, s.logger)
	wsChan := NewIpcServer(manager.config.IPCPort)
	go func() {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		for {
			select {
			case ws, ok := <-wsChan:
				if !ok {
					s.logger.Infof("ip server closed channel")
					return
				}
				go s.handleRegister(ctx, ws)
			}
		}
	}()
	return s
}

func (s *PluginsServer) handleRegister(ctx context.Context, client *wsConnection) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	m, err := client.wsRead()
	if err != nil {
		s.logger.Errorf("plugin register err:%s", err.Error())
		return
	}
	var message messages.PluginRegisterRequestJson
	err = json.Unmarshal(m.data, &message)
	if err != nil {
		s.logger.Errorf("plugin register:bad message err:%s", err.Error)
		return
	}
	plugin := s.registerPlugin(message.Data.PluginId)
	if client.registered {
		plugin.handleWs(ctx, client)
		return
	}
	msg := messages.PluginRegisterResponseJson{
		Data: messages.PluginRegisterResponseJsonData{
			GatewayVersion: "",
			PluginId:       message.Data.PluginId,
			Preferences:    s.manager.preferences,
			UserProfile:    *s.manager.GetUserProfile(),
		},
		MessageType: int(messages.MessageType_PluginRegisterResponse),
	}
	data, err := json.Marshal(msg)
	if err != nil {
		s.logger.Errorf(err.Error())
		return
	}
	err = client.wsWrite(websocket.TextMessage, data)
	if err != nil {
		s.logger.Errorf("plugin register:bad message err:%s", err.Error)
		return
	}

	s.logger.Infof("plugin: %s register success", message.Data.PluginId)
	client.registered = true
	plugin.handleWs(ctx, client)
}

//func (s *PluginsServer) RegisterPlugin(connection ipc.Connection) (ipc.PluginHandler, error) {
//	m, err := connection.ReadMessage()
//	if err != nil {
//		return nil, err
//	}
//	data, err := json.Marshal(m.Data)
//	if m.MessageType != messages.MessageType_PluginRegisterRequest {
//		return nil, fmt.Errorf("MessageType need PluginRegisterRequest")
//	}
//	var registerMessage messages.PluginRegisterRequestJsonData
//	err = json.Unmarshal(data, &registerMessage)
//	if err != nil {
//		return nil, err
//	}
//	responseData := messages.PluginRegisterResponseJsonData{
//		GatewayVersion: constant.Version,
//		PluginId:       registerMessage.PluginId,
//		Preferences:    *s.manager.GetPreferences(),
//		UserProfile:    *s.manager.config.UserProfile,
//	}
//	err = connection.WriteMessage(messages.MessageType_PluginRegisterResponse, responseData)
//	if err != nil {
//		return nil, err
//	}
//	plugin := s.registerPlugin(registerMessage.PluginId)
//	plugin.register(connection)
//	return plugin, nil
//}

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
