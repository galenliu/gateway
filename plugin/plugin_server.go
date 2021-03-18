package plugin

//	plugin server
import (
	"fmt"
	"gateway/config"
	"gateway/pkg/log"
	json "github.com/json-iterator/go"
	"strconv"
	"sync"
)

type PluginsServer struct {
	Plugins   map[string]*Plugin
	locker    *sync.Mutex
	manager   *AddonManager
	ipc       *IpcServer
	closeChan chan struct{}
	verbose   bool
}

func NewPluginServer(manager *AddonManager) *PluginsServer {
	server := &PluginsServer{}
	server.closeChan = make(chan struct{})
	server.Plugins = make(map[string]*Plugin, 30)
	server.locker = new(sync.Mutex)
	server.manager = manager
	server.ipc = NewIpcServer(":" + strconv.Itoa(config.Conf.Ports["ipc"]))
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
		plugin := s.registerPlugin(pluginId)
		go plugin.handleMessage(data)
	}
}

func (s *PluginsServer) registerHandler(data []byte, c *Connection) {
	pluginId := json.Get(data, "data", "pluginId").ToString()
	plugin := s.registerPlugin(pluginId)
	plugin.handleConnection(c)
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

func (s *PluginsServer) loadPlugin(addonPath, id, exec string) {
	plugin := s.registerPlugin(id)
	plugin.exec = exec
	plugin.execPath = addonPath
	go plugin.start()
}

func (s *PluginsServer) uninstallPlugin(packageId string) {
	plugin, ok := s.Plugins[packageId]
	if !ok {
		log.Error("plugin not exist")
		return
	}
	plugin.Stop()
	delete(s.Plugins, packageId)

}

func (s *PluginsServer) registerPlugin(packageId string) *Plugin {
	plugin, ok := s.Plugins[packageId]
	if !ok {
		plugin = NewPlugin(s, packageId)
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
		case <-s.closeChan:
			log.Debug("plugin server closed")
			return
		}
	}
}

//if server stop, also need to stop all of package
func (s *PluginsServer) Stop() {
	s.ipc.close()
	s.closeChan <- struct{}{}
	for _, p := range s.Plugins {
		p.Stop()
	}
}

func (s *PluginsServer) addAdapter(adapter *AdapterProxy) {
	s.manager.addAdapter(adapter)
}
