package plugin

//	plugin server
import (
	"github.com/galenliu/gateway/pkg/constant"
	ipc "github.com/galenliu/gateway/pkg/ipc_server"
	"github.com/galenliu/gateway/pkg/logging"
	json "github.com/json-iterator/go"
	"sync"
)

type PluginsServer struct {
	Plugins   map[string]*Plugin
	locker    *sync.Mutex
	manager   *Manager
	ipc       *ipc.IPCServer
	closeChan chan struct{}
	logger    logging.Logger
	verbose   bool
}

func NewPluginServer(manager *Manager) *PluginsServer {
	server := &PluginsServer{}
	server.logger = manager.logger
	server.closeChan = make(chan struct{})
	server.Plugins = make(map[string]*Plugin, 30)
	server.locker = new(sync.Mutex)
	server.manager = manager
	server.ipc = ipc.NewAndStartIPCServer(manager.options.IPCPort, manager.logger)
	return server
}

func (s *PluginsServer) messageHandler(data []byte, c *ipc.Connection) {

	//如果是注册请求的话，调用registerPlugin处理注册
	t := json.Get(data, "messageType")
	if t.ValueType() != json.NumberValue{
	 	s.logger.Info("plugin message failed")
		 return
	 }
	messageType := t.ToInt()
	if messageType == constant.PluginRegisterRequest {
		s.registerHandler(data, c)
	} else {
		//获取Plugin，并且把消息交由对应的Plugin处理
		pluginId := json.Get(data, "data", "pluginId").ToString()
		plugin := s.registerPlugin(pluginId)
		go plugin.handleConnection(c, data)
	}
}

//并发，此处开启新协程，传入一个新的websocket连接,把读到的消息给MessageHandler
func (s *PluginsServer) handleConnection(c *ipc.Connection) {
	d, err := c.Read()
	if err != nil {
		s.logger.Info("plugin connection err:", err.Error())
		return
	}
	s.messageHandler(d, c)
}

func (s *PluginsServer) registerHandler(data []byte, c *ipc.Connection) {
	pluginId := json.Get(data, "data", "pluginId").ToString()
	plugin := s.registerPlugin(pluginId)
	go plugin.registerAndHandleConnection(c)
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
		s.logger.Error("plugin not exist")
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
	go func() {
		err := s.ipc.Start()
		if err != nil {
			s.logger.Error("IPC Start Failed. Err: %s", err.Error())
			return
		}
		for {
			select {
			//每一个连接都开一个协程处理了
			case conn := <-s.ipc.Connections:
				go s.handleConnection(conn)
			case <-s.closeChan:
				return
			}
		}
	}()
	return nil
}

// Stop if server stop, also need to stop all of package
func (s *PluginsServer) Stop() error {
	err := s.ipc.Close()
	if err != nil {
		return err
	}
	s.closeChan <- struct{}{}
	s.logger.Info("Plugin server stopped")
	for _, p := range s.Plugins {
		p.unload()
	}
	return nil
}
