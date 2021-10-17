package ipc_server

import (
	"github.com/galenliu/gateway-grpc"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/server"
	"github.com/gorilla/websocket"
	json "github.com/json-iterator/go"
	"net/http"
	"sync"
	"time"
)

var upgrade = websocket.Upgrader{
	//ReadBufferSize:   1024,
	//WriteBufferSize:  1024,
	HandshakeTimeout: 5 * time.Second,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type IPCServer struct {
	logger       logging.Logger
	path         string
	port         string
	locker       *sync.Mutex
	pluginServer server.PluginServer
	userProfile  *rpc.UsrProfile
	doneChan     chan struct{}
}

func NewIPCServer(pluginServer server.PluginServer, port string, userProfile *rpc.UsrProfile, log logging.Logger) *IPCServer {
	ipc := &IPCServer{}
	ipc.pluginServer = pluginServer
	ipc.logger = log
	ipc.doneChan = make(chan struct{})
	ipc.port = port
	ipc.userProfile = userProfile
	return ipc
}

func (s *IPCServer) Run() error {
	http.HandleFunc("/", s.handle)
	http.HandleFunc("/ws", s.handle)

	s.logger.Infof("IPC server run addr: %s", s.port)
	err := http.ListenAndServe(s.port, nil)
	if err != nil {
		return err
	}
	if err != nil {
		s.logger.Errorf("ipc start err : %s", err.Error())
	}
	return nil
}

//处理IPC客户端请求
func (s *IPCServer) handle(w http.ResponseWriter, r *http.Request) {
	var b []byte
	_,_=r.Body.Read(b)
	conn, err := upgrade.Upgrade(w, r, nil)
	s.logger.Info("ipc new connection:", conn.RemoteAddr().String())
	if err != nil {
		s.logger.Errorf("升级为websocket失败", err.Error())
		return
	}
	go s.readLoop(conn)
}

func (s *IPCServer) readLoop(conn *websocket.Conn) {
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			s.logger.Infof("ipc disconnection:", conn.RemoteAddr().String())
		}
	}(conn)
	var message rpc.PluginRegisterRequestMessage
	_, data, _ := conn.ReadMessage()
	err := json.Unmarshal(data, &message)
	if err != nil {
		return
	}
	if message.MessageType != rpc.MessageType_PluginRegisterRequest {
		return
	}
	responseMessage := &rpc.PluginRegisterResponseMessage{
		MessageType: rpc.MessageType_PluginRegisterResponse,
		Data: &rpc.PluginRegisterResponseMessage_Data{
			PluginId:       message.Data.PluginId,
			GatewayVersion: constant.Version,
			UserProfile:    s.userProfile,
			Preferences:    s.pluginServer.GetPreferences(),
		},
	}
	err = conn.WriteJSON(responseMessage)
	if err != nil {
		return
	}
	clint := NewClint(message.Data.PluginId, conn)

	var pluginHandler server.PluginHandler
	pluginHandler = s.pluginServer.RegisterPlugin(message.Data.PluginId, clint)
	for {
		message, err := clint.Read()
		if err != nil {
			return
		}
		err = pluginHandler.OnMsg(message.MessageType, message.Data)
		if err != nil {
			return
		}
		select {
		case <-s.doneChan:
			return
		}
	}
}
