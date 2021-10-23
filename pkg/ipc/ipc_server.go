package ipc

import (
	"context"
	"github.com/fasthttp/websocket"
	"github.com/galenliu/gateway-grpc"
	"github.com/galenliu/gateway/pkg/logging"
	json "github.com/json-iterator/go"
	"net/http"
	"sync"
	"time"
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 1 * time.Second,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type IPC struct {
	logger       logging.Logger
	path         string
	port         string
	locker       *sync.Mutex
	pluginServer PluginServer
	userProfile  *rpc.UsrProfile
	ctx          context.Context
}

func NewIPCServer(ctx context.Context, pluginServer PluginServer, port string, userProfile *rpc.UsrProfile, log logging.Logger) *IPC {
	ipc := &IPC{}
	ipc.pluginServer = pluginServer
	ipc.logger = log
	ipc.ctx = ctx
	ipc.port = port
	ipc.userProfile = userProfile
	go ipc.Run()
	return ipc
}

func (s *IPC) Run() {
	http.HandleFunc("/", s.handle)
	http.HandleFunc("/ws", s.handle)
	s.logger.Infof("IPC server addr: %s", s.port)
	err := http.ListenAndServe(s.port, nil)
	if err != nil {
		s.logger.Errorf("IPC Listen err : %s", err.Error())
	}
	if err != nil {
		s.logger.Errorf("IPC Listen err : %s", err.Error())
	}
}

//处理IPC客户端请求
func (s *IPC) handle(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Errorf("upgrade websocket failed err: %s", err.Error())
		return
	}
	s.logger.Infof("IPC ipcConnection addr: %s", conn.RemoteAddr().String())
	go s.readLoop(conn)
}

func (s *IPC) readLoop(conn *websocket.Conn) {

	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			s.logger.Infof("ipcConnection closed")
		}
	}(conn)

	clint := &ipcConnection{conn}

	pluginHandler, err := s.pluginServer.RegisterPlugin(clint)
	if err != nil {
		s.logger.Infof("ipcConnection err: %s", err.Error())
		return
	}

	ctx, cancelFunc := context.WithCancel(s.ctx)
	for {
		select {
		case <-ctx.Done():
			s.logger.Infof("ipc server exit")
			_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "ipc exit"))
			_ = clint.Close()
			cancelFunc()
			return

		default:
			message, err := clint.ReadMessage()
			if err != nil {
				cancelFunc()
				return
			}
			s.logger.Debugf("ipc server rev: %+v", message)
			err = pluginHandler.OnMsg(message.MessageType, message.Data)
			if err != nil {
				cancelFunc()
				return
			}
		}
	}
}

type ipcConnection struct {
	*websocket.Conn
}

func (c *ipcConnection) WriteMessage(message *rpc.BaseMessage) error {
	baseMessage := BaseMessage{MessageType: int(message.MessageType)}
	err := json.Unmarshal(message.Data, &baseMessage.Data)
	if err != nil {
		return err
	}
	return c.WriteJSON(baseMessage)
}

func (c *ipcConnection) ReadMessage() (*rpc.BaseMessage, error) {
	//var msg rpc.BaseMessage
	var msg BaseMessage
	err := c.ReadJSON(&msg)
	data, err := json.Marshal(msg.Data)
	if err != nil {
		return nil, err
	}
	return &rpc.BaseMessage{MessageType: rpc.MessageType(msg.MessageType), Data: data}, nil
}

type BaseMessage struct {
	MessageType int                    `json:"messageType"`
	Data        map[string]interface{} `json:"data"`
}
