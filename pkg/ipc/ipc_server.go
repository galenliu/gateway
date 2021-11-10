package ipc

import (
	"context"
	"github.com/fasthttp/websocket"
	"github.com/galenliu/gateway-grpc"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/util"
	json "github.com/json-iterator/go"
	"net/http"
	"sync"
	"time"
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 1 * time.Second,
	Subprotocols:     []string{"websocket"},
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type IPCServer struct {
	logger       logging.Logger
	path         string
	port         string
	locker       *sync.Mutex
	pluginServer PluginServer
	userProfile  *rpc.UsrProfile
}

func NewIPCServer(pluginServer PluginServer, port string, userProfile *rpc.UsrProfile, log logging.Logger) *IPCServer {
	ipc := &IPCServer{}
	ipc.pluginServer = pluginServer
	ipc.logger = log
	ipc.port = port
	ipc.userProfile = userProfile
	go ipc.Run()
	return ipc
}

func (s *IPCServer) Run() {

	http.HandleFunc("/", s.handle)
	http.HandleFunc("/ws", s.handle)
	s.logger.Infof("ipc listen addr: %s", s.port)
	err := http.ListenAndServe("127.0.0.1"+s.port, nil)
	if err != nil {
		s.logger.Errorf("ipcServer Listen err : %s", err.Error())
	}
	if err != nil {
		s.logger.Errorf("ipcServer Listen err : %s", err.Error())
	}
}

//处理IPC客户端请求
func (s *IPCServer) handle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Errorf("upgrade websocket failed err: %s", err.Error())
		return
	}
	s.logger.Infof("ipc connection addr: %s", conn.RemoteAddr().String())
	go s.readLoop(conn)
}

// g
func (s *IPCServer) readLoop(conn *websocket.Conn) {
	conn.SetPongHandler(func(appData string) error {
		s.logger.Info("ping request: %s", appData)
		return conn.WriteMessage(websocket.PongMessage, nil)
	})
	clint := &connection{Conn: conn, logger: s.logger}
	pluginHandler, err := s.pluginServer.RegisterPlugin(clint)
	if err != nil {
		s.logger.Infof("register err: %s", err.Error())
		return
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer func() {
		cancelFunc()
		err := conn.Close()
		s.logger.Errorf("clint %s closed", clint.GetPluginId())
		if err != nil {
			s.logger.Errorf("clint %s close err: %s", clint.pluginId, err.Error())
			return
		}
	}()
	for {
		select {
		case <-ctx.Done():
			s.logger.Infof("ipc server exit")
			_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "ipc exit"))
			_ = clint.Close()
			return

		default:
			message, err := clint.ReadMessage()
			if err != nil {
				s.logger.Errorf("clint %s read err : %s", clint.GetPluginId(), err.Error())
				return
			}
			pluginHandler.OnMsg(message.MessageType, message.Data)
		}
	}
}

type connection struct {
	*websocket.Conn
	logger   logging.Logger
	pluginId string
}

func (c *connection) WriteMessage(message *rpc.BaseMessage) error {
	baseMessage := BaseMessage{MessageType: int(message.MessageType)}
	err := json.Unmarshal(message.Data, &baseMessage.Data)
	if err != nil {
		return err
	}
	err = c.Conn.WriteJSON(baseMessage)
	if err != nil {
		return err
	}
	if c.GetPluginId() == "" {
		c.logger.Debugf("Plugin send :%s", util.JsonIndent(baseMessage))
	} else {
		c.logger.Debugf("Plugin %s send :%s", c.GetPluginId(), util.JsonIndent(baseMessage))
	}
	return nil
}

func (c *connection) ReadMessage() (*rpc.BaseMessage, error) {
	//var msg rpc.BaseMessage
	var msg BaseMessage
	_, data, err := c.Conn.ReadMessage()
	err = json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}
	if c.GetPluginId() == "" {
		c.logger.Debugf("Plugin register message :%s", util.JsonIndent(msg))
	} else {
		c.logger.Debugf("Plugin %s read: %s", c.GetPluginId(), util.JsonIndent(msg))
	}

	marshal, err := json.Marshal(msg.Data)
	if err != nil {
		return nil, err
	}
	return &rpc.BaseMessage{MessageType: rpc.MessageType(msg.MessageType), Data: marshal}, nil
}

func (c *connection) SetPluginId(id string) {
	c.pluginId = id
}

func (c *connection) GetPluginId() string {
	return c.pluginId
}

type BaseMessage struct {
	MessageType int                    `json:"messageType"`
	Data        map[string]interface{} `json:"data"`
}
