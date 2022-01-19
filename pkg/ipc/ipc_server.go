package ipc

import (
	"fmt"
	"github.com/fasthttp/websocket"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"github.com/galenliu/gateway/pkg/logging"
	"net/http"
	"sync"
)

var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketServer struct {
	logger       logging.Logger
	path         string
	port         string
	locker       *sync.Mutex
	closeChan    chan struct{}
	pluginServer PluginServer
	userProfile  *messages.PluginRegisterResponseJsonDataUserProfile
}

func NewIPCServer(pluginServer PluginServer, port string, userProfile *messages.PluginRegisterResponseJsonDataUserProfile, log logging.Logger) *WebSocketServer {
	ipc := &WebSocketServer{}
	ipc.pluginServer = pluginServer
	ipc.logger = log
	ipc.closeChan = make(chan struct{})
	ipc.port = port
	ipc.userProfile = userProfile
	go ipc.Run()
	return ipc
}

func (s *WebSocketServer) Run() {
	for {
		http.HandleFunc("/", s.handle)
		http.HandleFunc("/ws", s.handle)
		fmt.Printf("ipc listen addr: %s \t\n", s.port)
		err := http.ListenAndServe("localhost"+s.port, nil)
		if err != nil {
			s.logger.Errorf("ipcServer Listen err : %s", err.Error())
		}
		select {
		case _ = <-s.closeChan:
			return
		}
	}
}

//处理IPC客户端请求
func (s *WebSocketServer) handle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Errorf("upgrade websocket failed err: %s", err.Error())
		return
	}
	s.logger.Infof("ipc connection addr: %s", conn.RemoteAddr().String())
	go s.readLoop(conn)
}

func (s *WebSocketServer) readLoop(conn *websocket.Conn) {
	//conn.SetPongHandler(func(appData string) error {
	//	s.logger.Info("ping request: %s", appData)
	//	return conn.WriteMessage(websocket.PongMessage, nil)
	//})
	con := newConnection(conn, s.logger)
	pluginHandler, err := s.pluginServer.RegisterPlugin(con)
	if err != nil {
		s.logger.Infof("register err: %s", err.Error())
		return
	}
	if pluginHandler == nil {
		s.logger.Infof("register err: %s", "pluginHandler nil")
		return
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			s.logger.Errorf("clint %s close err: %s", con.getPluginId(), err.Error())
			return
		}
	}()
	for {
		mt, data, err := con.ReadMessage()
		if err != nil {
			if mt == messages.MessageType_MashalERROR {
				s.logger.Infof(err.Error())
				continue
			}
			s.logger.Errorf("plugin read err : %s", err.Error())
			return
		}
		pluginHandler.OnMsg(mt, data)
	}
}

func (s *WebSocketServer) Close() {
	select {
	case s.closeChan <- struct{}{}:
	}
}
