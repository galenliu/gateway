package addons

//	addons server
import (
	"context"
	"fmt"
	"gateway/config"
	"gateway/pkg/log"
	json "github.com/json-iterator/go"
	"go.uber.org/zap"
	"strconv"
	"sync"
)

const PATH = "/"

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
	server.Plugins = make(map[string]*Plugin, 30)
	server.locker = new(sync.Mutex)
	server.addonManager = manager
	ctx, _ := context.WithCancel(server.ctx)
	server.ipc = NewIpcServer(ctx, ":"+strconv.Itoa(config.Conf.Ports["ipc"]), PATH)
	return server
}

func (s *PluginsServer) messageHandler(data []byte, c *Connection) {

	//如果是注册请求的话，调用registerPlugin处理注册
	var messageType = json.Get(data, "messageType").ToInt()

	log.Debug(fmt.Sprintf("plugin servier rev message:%s", string(data)))

	if messageType == PluginRegisterRequest {
		s.registerHandler(data, c)
	} else {
		//获取Plugin，并且把消息交由对应的Plugin处理
		pluginId := json.Get(data, "data", "pluginId").ToString()
		plugin, ok := s.Plugins[pluginId]
		if !ok {
			log.Error(fmt.Sprintf("plugin(%s) can not find", pluginId))
			return
		}
		plugin.OnMessage(data)
	}
}

func (s *PluginsServer) registerHandler(data []byte, c *Connection) {

	pluginId := json.Get(data, "data", "pluginId").ToString()
	plugin, ok := s.Plugins[pluginId]
	if !ok {
		log.Error(fmt.Sprintf("plugin: %s no loading", pluginId))
		return
	}
	plugin.conn = c

	m := struct {
		MessageType int         `json:"messageType"`
		Data        interface{} `json:"data"`
	}{
		MessageType: PluginRegisterResponse,
		Data: struct {
			PluginID    string              `json:"pluginId"`
			UserProfile *config.UserProfile `json:"user_profile"`
			Preferences *config.Preferences `json:"preferences"`
		}{
			PluginID:    pluginId,
			UserProfile: config.GetUserProfile(),
			Preferences: config.GetPreferences(),
		},
	}
	data, _ = json.Marshal(m)
	plugin.sendData(data)
	plugin.registered = true
	log.Info(fmt.Sprintf("plugin: %s registered", pluginId))
}

func (s *PluginsServer) getPlugin(pluginId string) *Plugin {

	//通过读写锁获取plugin
	s.locker.Lock()
	defer s.locker.Unlock()
	p := s.Plugins[pluginId]
	return p
}

//此处开启新协程，传入一个新的websocket连接,把读到的消息给MessageHandler
func (s *PluginsServer) handleNewConnection(c *Connection) {
	c.connected = true
	for {
		if !c.connected {
			log.Info("lost connection")
			return
		}
		d, err := c.ReadMessage()
		if err != nil {
			c.connected = false
			continue
		}
		s.messageHandler(d, c)
	}
}

func (s *PluginsServer) loadPlugin(packageId, exec string, enabled bool) {
	plugin := s.registerPlugin(packageId, exec)
	if enabled {
		go plugin.start()
	} else {
		plugin.Stop()
	}
}

func (s *PluginsServer) uninstallPlugin(packageId string) {
	plugin, ok := s.Plugins[packageId]
	if !ok {
		log.Error("plugin not exist")
	}
	plugin.Stop()
	delete(s.Plugins, packageId)

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
	for {
		select {
		case conn := <-s.ipc.wsChan:
			go s.handleNewConnection(conn)
		case <-s.ctx.Done():
			s.Stop()
		}
	}
}

//if server stop, also need to stop all of package
func (s *PluginsServer) Stop() {
	close(s.ipc.wsChan)
	for _, v := range s.Plugins {
		v.Stop()
	}
}
