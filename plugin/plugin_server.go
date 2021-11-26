package plugin

//	plugin server
import (
	"context"
	"fmt"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/ipc"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"github.com/galenliu/gateway/pkg/logging"
	json "github.com/json-iterator/go"
	"sync"
)

type PluginsServer struct {
	Plugins   sync.Map
	manager   *Manager
	ipc       *ipc.IPCServer
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
	mt, j, err := clint.ReadMessage()
	if err != nil {
		return nil, err
	}
	if mt != messages.MessageType_PluginRegisterRequest {
		return nil, fmt.Errorf("MessageType need PluginRegisterRequest")
	}
	data, _ := json.Marshal(j)
	var registerMessage messages.PluginRegisterRequestJsonData
	err = json.Unmarshal(data, &registerMessage)
	if err != nil {
		return nil, fmt.Errorf("bad data")
	}

	responseData := messages.PluginRegisterResponseJsonData{
		GatewayVersion: constant.Version,
		PluginId:       registerMessage.PluginId,
		Preferences:    *s.getPreferences(),
		UserProfile:    *s.manager.config.UserProfile,
	}

	clint.SetPluginId(registerMessage.PluginId)
	err = clint.WriteMessage(messages.MessageType_PluginRegisterResponse, responseData)
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

func (s *PluginsServer) getPreferences() *messages.PluginRegisterResponseJsonDataPreferences {
	r := &messages.PluginRegisterResponseJsonDataPreferences{
		Language: "en-US",
		Units: messages.PluginRegisterResponseJsonDataPreferencesUnits{
			Temperature: "degree celsius",
		},
	}
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
