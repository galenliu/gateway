package ipc_server

import (
	"github.com/galenliu/gateway-grpc"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/server"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
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
	userProfile  *rpc.PluginRegisterResponseMessage_Data_UsrProfile
	preferences  *rpc.PluginRegisterResponseMessage_Data_Preferences
	doneChan     chan struct{}
}

func NewIPCServer(pluginServer server.PluginServer, port string, userProfile *rpc.PluginRegisterResponseMessage_Data_UsrProfile, preferences *rpc.PluginRegisterResponseMessage_Data_Preferences, log logging.Logger) *IPCServer {
	ipc := &IPCServer{}
	ipc.pluginServer = pluginServer
	ipc.logger = log
	ipc.doneChan = make(chan struct{})
	ipc.port = port
	ipc.userProfile = userProfile
	ipc.preferences = preferences
	return ipc
}

func (s *IPCServer) Start() error {
	go func() {
		err := func() error {
			http.HandleFunc("/", s.handle)
			s.logger.Infof("IPC server run addr: %s", s.port)
			err := http.ListenAndServe(s.port, nil)
			if err != nil {
				return err
			}
			return nil
		}()
		if err != nil {
			s.logger.Errorf("ipc start err : %s", err.Error())
		}
	}()
	return nil
}

func (s *IPCServer) Stop() error {
	return nil
}

func (s *IPCServer) handle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrade.Upgrade(w, r, nil)
	s.logger.Debug("accept new connection")
	if conn == nil {
		return
	}
	var message rpc.PluginRegisterRequestMessage
	err = conn.ReadJSON(&message)
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
			Preferences:    s.preferences,
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
