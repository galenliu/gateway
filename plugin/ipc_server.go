package plugin

import (
	"context"
	"fmt"
	"gateway/pkg/log"
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

type IpcServer struct {
	addr   string
	path   string
	ctx    context.Context
	wsChan chan *Connection
	locker *sync.Mutex
}

func NewIpcServer(_ctx context.Context, _addr string) *IpcServer {
	ipc := &IpcServer{
		addr:   _addr,
		ctx:    _ctx,
		wsChan: make(chan *Connection),
	}
	return ipc
}

type Connection struct {
	locker    *sync.Mutex
	ws        *websocket.Conn
	connected bool
}

func (c *Connection) send(data []byte) {
	c.locker.Lock()
	defer c.locker.Unlock()
	log.Debug(fmt.Sprintf("plugin send message :\t\n %s", string(data)))
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
	log.Debug("accept new connection")
	if conn == nil {
		return
	}
	//升级协议时可能发生的错误
	if err != nil {
		log.Error("ipc server upgrade failed,err: ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var ws = &Connection{
		ws:        conn,
		connected: true,
		locker:    new(sync.Mutex),
	}
	server.wsChan <- ws
}

func (server *IpcServer) Serve() {
	http.HandleFunc("/", server.handle)
	log.Info("plugin server start on port: %s", server.addr)
	err := http.ListenAndServe(server.addr, nil)
	log.Info(fmt.Sprintf("ipc server listening addr: %s", server.addr))
	if err != nil {
		log.Error("ipc server fail,err: %s", err.Error())
	}

}
