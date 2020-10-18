package plugin

//	plugin server
import (
	"context"
	messages "github.com/galeuliu/gateway-schema"
	json "github.com/json-iterator/go"
	"go.uber.org/zap"
	"sync"
)

const PATH = "/v1"
const ADDR = ":9500"

type PluginsServer struct {
	Plugins      map[string]*Plugin
	locker       *sync.Mutex
	addonManager *AddonsManager
	ipc          *IpcServer
	ctx          context.Context
	verbose      bool
	logger       **zap.Logger
}

func NewPluginServer(manager *AddonsManager, _ctx context.Context) *PluginsServer {
	server := &PluginsServer{}
	server.ctx = _ctx
	server.locker = new(sync.Mutex)
	server.addonManager = manager
	ctx, _ := context.WithCancel(server.ctx)
	server.ipc = NewIpcServer(ctx, ADDR, PATH)
	return server
}

func (s *PluginsServer) messageHandler(m messages.BaseMessage, c *Connection) {

	//首先验证message是否合法,并且序列化
	data, err := messages.CheckMessage(m)
	if err != nil {

		return

	}
	log.Info("plugin server rev message ", zap.Int("MessageType", m.MessageType), zap.Any("data", m.Data))

	//如果是注册请求的话，调用registerPlugin处理注册
	if m.MessageType == messages.MessageTypePluginRegisterRequest {
		s.registerHandler(m, data, c)

	} else {
		//获取Plugin，并且把消息交由对应的Plugin处理
		any := json.Get(data, "plugin_id")
		pluginId := any.ToString()
		plugin := s.getPlugin(pluginId)
		plugin.OnMessage(m, data)
	}
}

func (s *PluginsServer) registerHandler(message messages.BaseMessage, data []byte, c *Connection) {

	pluginId := json.Get(data, "data").Get("plugin_id").ToString()
	plugin := s.getPlugin(pluginId)
	plugin.ws = c
	r := messages.PluginRegisterResponse{
		PluginId: pluginId,
		UserProfile: &messages.UserProfile{
			GatewayVersion: s.addonManager.userProfile.GatewayVersion,
			BaseDir:        s.addonManager.userProfile.BaseDir,
			AddonsDir:      s.addonManager.userProfile.AddonsDir,
			ConfigDir:      s.addonManager.userProfile.ConfigDir,
			MediaDir:       s.addonManager.userProfile.MediaDir,
			LogDir:         s.addonManager.userProfile.LogDir,
		},
		Preferences: &messages.Preferences{
			Language: s.addonManager.preferences.Language,
			Units:    messages.Units{Temperature: s.addonManager.preferences.Units.Temperature},
		},
	}
	m := messages.BaseMessage{
		MessageType: messages.MessageTypePluginRegisterResponse,
		Data:        r,
	}
	plugin.ws.SendMessage(m)
	plugin.registered = true

}

func (s *PluginsServer) getPlugin(pluginId string) *Plugin {

	//通过读写锁获取plugin
	s.locker.Lock()
	defer s.locker.Unlock()
	p := s.Plugins[pluginId]
	if p == nil {
		p = NewPlugin(pluginId, s, s.ctx)
	}
	return p
}

//此处开启新协程，传入一个新的websocket连接,把读到的消息给MessageHandler
func (s *PluginsServer) newConn(c *Connection) {
	c.connected = true
	for {
		if !c.connected {
			log.Info("lost connection", zap.Any("clint", c.ws.RemoteAddr()))
			return
		}
		m, err := c.ReadMessage()
		if err != nil {
			c.connected = false
			continue
		}
		s.messageHandler(m, c)
	}
}

func (s *PluginsServer) sendMsg() {

}

//开启goroutines处理ipc_server中ws
func (s *PluginsServer) Run() {
	go s.ipc.Serve()
	select {
	case conn := <-s.ipc.wsChan:
		go s.newConn(conn)
	case <-s.ctx.Done():
		s.close()
	}
}

func (s *PluginsServer) close() {
	close(s.ipc.wsChan)
}
