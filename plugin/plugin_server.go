package plugin

//	plugin server
import (
	"fmt"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/plugin/internal"
	json "github.com/json-iterator/go"
	"sync"
)

type PluginsServer struct {
	Plugins   map[string]*Plugin
	locker    *sync.Mutex
	manager   *Manager
	ipc       *IpcServer
	closeChan chan struct{}
	logger    logging.Logger
	verbose   bool
}

func NewPluginServer(manager *Manager, log logging.Logger) *PluginsServer {
	server := &PluginsServer{}
	server.logger = log
	server.closeChan = make(chan struct{})
	server.Plugins = make(map[string]*Plugin, 30)
	server.locker = new(sync.Mutex)
	server.manager = manager
	server.ipc = NewIpcServer()
	return server
}

func (s *PluginsServer) messageHandler(data []byte, c *Connection) {

	//如果是注册请求的话，调用registerPlugin处理注册
	var m = json.Get(data, "messageType")

	if err := m.LastError(); err != nil {
		logging.Info("messageType err")
		return
	}
	messageType := m.ToInt()
	s.logger.Debug("%s: \t\n %s", internal.MessageTypeToString(messageType), string(data))

	if messageType == internal.PluginRegisterRequest {
		s.registerHandler(data, c)
	} else {
		//获取Plugin，并且把消息交由对应的Plugin处理
		pluginId := json.Get(data, "data", "pluginId").ToString()
		plugin := s.registerPlugin(pluginId)
		go plugin.handleConnection(c, data)
	}
}

func (s *PluginsServer) registerHandler(data []byte, c *Connection) {
	pluginId := json.Get(data, "data", "pluginId").ToString()
	plugin := s.registerPlugin(pluginId)
	plugin.registerAndHandleConnection(c)
}

//此处开启新协程，传入一个新的websocket连接,把读到的消息给MessageHandler
func (s *PluginsServer) handlerConnection(c *Connection) {
	d, err := c.readMessage()
	if err != nil {
		logging.Info("plugin connection err:", err.Error())
		return
	}
	s.messageHandler(d, c)
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
		logging.Error("plugin not exist")
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
	go s.ipc.Serve()
	if !s.manager.running {
		return fmt.Errorf("addon Manager stoped")
	}
	go func() {
		for {
			select {
			//每一个连接都开一个协程处理了
			case conn := <-s.ipc.wsChan:
				go s.handlerConnection(conn)
			case <-s.closeChan:
				fmt.Print("plugin server closed")
			}
		}
	}()
	return nil
}

// Stop if server stop, also need to stop all of package
func (s *PluginsServer) Stop() {
	s.ipc.close()
	s.closeChan <- struct{}{}
	for _, p := range s.Plugins {
		p.unload()
	}
}

