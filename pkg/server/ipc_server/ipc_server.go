package ipc_server

import (
	"context"
	"fmt"
	"github.com/galenliu/gateway-grpc"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/server"
	"github.com/gorilla/websocket"
	json "github.com/json-iterator/go"
	ws "golang.org/x/net/websocket"
	"net/http"
	"sync"
	"time"
)

var upgrade = websocket.Upgrader{
	//ReadBufferSize:   1024,
	//WriteBufferSize:  1024,
	//HandshakeTimeout: 5 * time.Second,
	CheckOrigin: func(r *http.Request) bool {
		fmt.Printf("method: %s, url: %s , path: %s\n", r.Method, r.URL.String(), r.Host)
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
	ctx          context.Context
}

func NewIPCServer(ctx context.Context, pluginServer server.PluginServer, port string, userProfile *rpc.UsrProfile, log logging.Logger) *IPCServer {
	ipc := &IPCServer{}
	ipc.pluginServer = pluginServer
	ipc.logger = log
	ipc.ctx = ctx
	ipc.port = port
	ipc.userProfile = userProfile
	go ipc.Run()
	return ipc
}

func (s *IPCServer) Run() {
	http.HandleFunc("/", s.handle)
	http.HandleFunc("/ws", s.handle)
	s.logger.Infof("IPC server addr: %s", s.port)
	err := http.ListenAndServe(s.port, nil)
	if err != nil {
		s.logger.Errorf("ipc Listen err : %s", err.Error())
	}
	if err != nil {
		s.logger.Errorf("ipc start err : %s", err.Error())
	}
}

//处理IPC客户端请求
func (s *IPCServer) handleWs(ws *ws.Conn) {

	s.readLoopWs(ws)
}

//处理IPC客户端请求
func (s *IPCServer) handle(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrade.Upgrade(w, r, nil)
	s.logger.Info("ipc new connection:", conn.RemoteAddr().String())
	if err != nil {
		s.logger.Errorf("升级为websocket失败", err.Error())
		return
	}
	s.readLoop(conn)
}

func (s *IPCServer) readLoop(conn *websocket.Conn) {

	var err error
	var message rpc.PluginRegisterRequestMessage
	i, data, err := conn.ReadMessage()
	s.logger.Infof("-----messageType: %s data:%s ,err: %s", i, data, err.Error())
	err = json.Unmarshal(data, &message)
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
		case <-s.ctx.Done():
			_ = clint.Close()
			return
		}
	}
}

func (s *IPCServer) readLoopWs(conn *ws.Conn) {

	defer func(conn *ws.Conn) {
		err := conn.Close()
		if err != nil {
			s.logger.Info(err.Error())
		}
	}(conn)
	var err error
	var message rpc.PluginRegisterRequestMessage
	time.Sleep(1 * time.Second)

	var data []byte
	err = ws.Message.Receive(conn, &message)
	s.logger.Infof("----- data:%s ,err: %s", data, err.Error())
	err = json.Unmarshal(data, &message)
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
	d, err := json.Marshal(responseMessage)
	_, err = conn.Write(d)
	if err != nil {
		return
	}
	//clint := NewClint(message.Data.PluginId, conn)
	//
	//var pluginHandler server.PluginHandler
	//pluginHandler = s.pluginServer.RegisterPlugin(message.Data.PluginId, clint)
	//for {
	//	message, err := clint.Read()
	//	if err != nil {
	//		return
	//	}
	//	err = pluginHandler.OnMsg(message.MessageType, message.Data)
	//	if err != nil {
	//		return
	//	}
	//	select {
	//	case <-s.doneChan:
	//		return
	//	}
	//}
}
