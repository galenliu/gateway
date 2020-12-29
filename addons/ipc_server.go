package addons

import (
	"context"
	"fmt"
	"gateway/pkg/log"
	"github.com/gorilla/websocket"
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

func (c *Connection) send(data []byte) {
	log.Debug(fmt.Sprintf("send date : %s", string(data)))
	err := c.ws.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		c.connected = false
	}
}

func (c *Connection) ReadMessage() (data []byte, err error) {

	_, data, err = c.ws.ReadMessage()
	if err != nil {
		log.Error("connection read message err:", err.Error())
		c.connected = false
	}
	return
}

func (server *IpcServer) handle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrade.Upgrade(w, r, nil)

	//升级协议时可能发生的错误
	if err != nil {
		log.Error("ipc server upgrade failed,err: ", err.Error())
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
	http.HandleFunc("/", server.handle)
	log.Info(fmt.Sprintf("ipc server listening addr: %s", server.addr))

	err := sev.ListenAndServe()
	if err != nil {
		log.Error("ipc server fail,err: ", err.Error())
	}

}
