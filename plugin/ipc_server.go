package plugin

import (
	"context"
	messages "github.com/galenliu/smartassistant-ipc"
	"github.com/gorilla/websocket"
	json "github.com/json-iterator/go"
	"go.uber.org/zap"
	"net/http"
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type IpcServer struct {
	addr   string
	path   string
	ctx    context.Context
	wsChan chan *Connection
}

func NewIpcServer(_ctx context.Context, _addr string, _path string) *IpcServer {
	ipc := &IpcServer{
		addr:   _addr,
		path:   _path,
		ctx:    _ctx,
		wsChan: make(chan *Connection),
	}
	return ipc
}

type Connection struct {
	ws        *websocket.Conn
	connected bool
}

func (c *Connection) SendMessage(message messages.BaseMessage) {

	data, _ := json.MarshalIndent(message, "", "	")
	err := c.ws.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Error("send message err", zap.Error(err))
		c.connected = false
	}
}

func (c *Connection) ReadMessage() (m messages.BaseMessage, err error) {

	err = c.ws.ReadJSON(&m)
	if err != nil {
		log.Error("connetion read message err", zap.Error(err))
		c.connected = false
	}
	return
}

func (server *IpcServer) handle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrade.Upgrade(w, r, nil)

	//升级协议时可能发生的错误
	if err != nil {
		log.Error("ipc server upgrade faild,err: %v", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var ws = &Connection{
		ws:        conn,
		connected: true,
	}
	server.wsChan <- ws
}

func (server *IpcServer) Serve() {
	sev := &http.Server{
		Addr: server.addr,
	}
	http.HandleFunc(server.path, server.handle)
	log.Info("listening", zap.String("addr", server.addr))

	err := sev.ListenAndServe()
	if err != nil {
		log.Error("ipc server fail", zap.Error(err))
	}

}
