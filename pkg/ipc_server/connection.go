package ipc_server

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Connection struct {
	*websocket.Conn
	locker *sync.Mutex
}

func NewConn(conn *websocket.Conn) *Connection {
	c := &Connection{}
	c.Conn = conn
	return c
}

func (c *Connection) Send(data []byte) {
	c.locker.Lock()
	defer c.locker.Unlock()
	err := c.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {

	}
}

func (c *Connection) Read() (data []byte, err error) {
	_, data, err = c.ReadMessage()
	return
}
