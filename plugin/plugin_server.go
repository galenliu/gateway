package plugin

//	plugin server
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


type PluginsServer struct {
	Plugins map[string]*Plugin
	locker  *sync.Mutex
	manager *AddonManager
	ipc     *IpcServer
	ctx     context.Context
	verbose bool
	logger  **zap.Logger
}

func NewPluginServer(manager *AddonManager, _ctx context.Context) *PluginsServer {
	server := &PluginsServer{}
	server.ctx = _ctx
	server.Plugins = make(map[string]*Plugin, 30)
	server.locker = new(sync.Mutex)
	server.manager = manager
	ctx, _ := context.WithCancel(server.ctx)
	server.ipc = NewIpcServer(ctx, ":"+strconv.Itoa(config.Conf.Ports["ipc"]))
	return server
}

func (s *PluginsServer) messageHandler(data []byte, c *Connection) {

	log.Debug(fmt.Sprintf("plugin rev message: \t\n %s", string(data)))

	//如果是注册请求的话，调用registerPlugin处理注册
	var m = json.Get(data, "messageType")
	if err := m.LastError(); err != nil {
		log.Info("messageType err")
		return
	}
	messageType := m.ToInt()

	if messageType == PluginRegisterRequest {
		s.registerHandler(data, c)
	} else {
		//获取Plugin，并且把消息交由对应的Plugin处理
		pluginId := json.Get(data, "data", "pluginId").ToString()
		plugin, err := s.getPlugin(pluginId)
		if err != nil {
			log.Error(err.Error())
			c.connected = false
			return
		}
		go plugin.handleMessage(data)
	}
}

func (s *PluginsServer) registerHandler(data []byte, c *Connection) {
	pluginId := json.Get(data, "data", "pluginId").ToString()
	plugin, ok := s.Plugins[pluginId]
	if !ok {
		plugin = NewPlugin(s, pluginId)
		s.Plugins[pluginId] = plugin
	}
	plugin.handleConnection(c)
}

func (s *PluginsServer) getPlugin(pluginId string) (*Plugin, error) {
	//通过读写锁获取plugin
	s.locker.Lock()
	defer s.locker.Unlock()
	p, ok := s.Plugins[pluginId]
	if !ok {
		return nil, fmt.Errorf("plugin(%s) not found", pluginId)
	}
	return p, nil
}

//此处开启新协程，传入一个新的websocket连接,把读到的消息给MessageHandler
func (s *PluginsServer) readConnectionLoop(c *Connection) {
	for {
		if !c.connected {
			log.Info("lost connection")
			return
		}
		d, err := c.ReadMessage()
		if err != nil {
			log.Info("read connection err:", err.Error())
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
		plugin = NewPlugin(s, packageId)
		plugin.exec = exec
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
			go s.readConnectionLoop(conn)
		case <-s.ctx.Done():
			log.Debug("connection close")
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

func (s *PluginsServer) addAdapter(adapter *AdapterProxy) {
	s.manager.addAdapter(adapter)
}
