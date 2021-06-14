package plugin

import (
	"fmt"
	"github.com/galenliu/gateway/configs"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
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
	wsChan chan *Connection
	locker *sync.Mutex
}

func NewIpcServer() *IpcServer {
	ipc := &IpcServer{
		addr:   "localhost:" + strconv.Itoa(configs.GetIpcPort()),
		wsChan: make(chan *Connection, 2),
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
	err := c.ws.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		c.connected = false
	}
}

func (c *Connection) readMessage() (data []byte, err error) {

	_, data, err = c.ws.ReadMessage()
	if err != nil {
		logging.Error("connection read message err:", err.Error())
		c.connected = false
	}
	return
}

func (server *IpcServer) handle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrade.Upgrade(w, r, nil)
	logging.Debug("accept new connection")
	if conn == nil {
		return
	}
	//升级协议时可能发生的错误
	if err != nil {
		logging.Error("ipc server upgrade failed,err: ", err.Error())
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
	logging.Info("plugin server execute on port: %s", server.addr)
	err := http.ListenAndServe(server.addr, nil)
	logging.Info(fmt.Sprintf("ipc server listening addr: %s", server.addr))
	if err != nil {
		logging.Error("ipc server fail,err: %s", err.Error())
	}
}

func (server *IpcServer) close() {
	close(server.wsChan)
}
