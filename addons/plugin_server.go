package addons

//	addons server
import (
	"context"
	messages "gitee.com/liu_guilin/WebThings-schema"
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
	server.Plugins = make(map[string]*Plugin)
	server.locker = new(sync.Mutex)
	server.addonManager = manager
	ctx, _ := context.WithCancel(server.ctx)
	server.ipc = NewIpcServer(ctx, ADDR, PATH)
	go server.Start()
	return server
}

func (s *PluginsServer) messageHandler(m messages.BaseMessage, c *Connection) {

	//首先验证message是否合法,并且序列化
	data, err := messages.CheckMessage(m)
	if err != nil {

		return

	}
	log.Info("addons server rev message ")

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

	m := messages.BaseMessage{
		MessageType: messages.MessageTypePluginRegisterResponse,
	}
	plugin.ws.SendMessage(m)
	plugin.registered = true

}

func (s *PluginsServer) getPlugin(pluginId string) *Plugin {

	//通过读写锁获取plugin
	s.locker.Lock()
	defer s.locker.Unlock()
	p := s.Plugins[pluginId]
	return p
}

//此处开启新协程，传入一个新的websocket连接,把读到的消息给MessageHandler
func (s *PluginsServer) newConn(c *Connection) {
	c.connected = true
	for {
		if !c.connected {
			log.Info("lost connection")
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

func (s *PluginsServer) loadPlugin(packageId, exec string, enabled bool) {
	plugin := s.registerPlugin(packageId, exec)
	if enabled {
		plugin.start()
	} else {
		plugin.Stop()
	}
}

func (s *PluginsServer) uninstallPlugin(packageId string) {
	pkg := s.getPlugin(packageId)
	if pkg == nil {
		return
	} else {
		pkg.Stop()
		delete(s.Plugins, packageId)
	}
}

func (s *PluginsServer) registerPlugin(packageId string, exec string) *Plugin {
	plugin := s.Plugins[packageId]
	if plugin == nil {
		plugin = NewPlugin(s, packageId, exec, s.ctx)
		s.Plugins[packageId] = plugin
	}
	return plugin
}

//create goroutines handle ipc massage
func (s *PluginsServer) Start() {
	go s.ipc.Serve()
	select {
	case conn := <-s.ipc.wsChan:
		go s.newConn(conn)
	case <-s.ctx.Done():
		s.Stop()
	}
}

//if server stop, also need to stop all of package
func (s *PluginsServer) Stop() {
	close(s.ipc.wsChan)
	for _, v := range s.Plugins {
		v.cancelFunc()
	}
}
