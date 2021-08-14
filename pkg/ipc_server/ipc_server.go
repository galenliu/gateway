package ipc_server

import (
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/rpc"
	"github.com/galenliu/gateway/plugin"
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

type PluginHandler interface {
	MessageHandler(mt rpc.MessageType, data []byte) error
}

type IPCServer struct {
	logger       logging.Logger
	addr         string
	path         string
	port         string
	locker       *sync.Mutex
	pluginServer *plugin.PluginsServer
	userProfile  []byte
	preferences  []byte
	doneChan     chan struct{}
}

func NewIPCServer(server *plugin.PluginsServer, port string, userProfile []byte, preferences []byte, log logging.Logger) *IPCServer {
	ipc := &IPCServer{}
	ipc.pluginServer = server
	ipc.logger = log
	ipc.doneChan = make(chan struct{})
	ipc.port = port
	ipc.userProfile = userProfile
	ipc.preferences = preferences
	return ipc
}

func (s *IPCServer) Start() error {
	http.HandleFunc("/", s.handle)
	s.logger.Info("IPC server run addr: %s", s.addr)
	err := http.ListenAndServe(s.addr, nil)
	if err != nil {
		s.logger.Error("ipc s fail,err: %s", err.Error())
		return err
	}
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
	err = conn.WriteJSON(rpc.PluginRegisterResponseMessage{
		MessageType: 0,
	})
	if err != nil {
		return
	}

	clint := NewClint(message.Data.PluginId, conn)
	var pluginHandler PluginHandler
	pluginHandler = s.pluginServer.RegisterPlugin(message.Data.PluginId, clint)
	for {
		message, err := clint.Read()
		if err != nil {
			return
		}
		err = pluginHandler.MessageHandler(message.MessageType, message.Data)
		if err != nil {
			return
		}
		select {
		case <-s.doneChan:
			return
		}
	}
}
