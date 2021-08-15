package ipc_server

import (
	"github.com/galenliu/gateway/pkg/rpc"
	"github.com/gorilla/websocket"
	"sync"
)

type Clint struct {
	*websocket.Conn
	locker   *sync.Mutex
	pluginId string
}

func NewClint(pluginId string, conn *websocket.Conn) *Clint {
	c := Clint{}
	c.Conn = conn
	return &c
}

func (c *Clint) Send(message *rpc.BaseMessage) error {
	return c.Conn.WriteJSON(message)
}

func (c *Clint) Read() (message *rpc.BaseMessage, err error) {
	err = c.ReadJSON(message)
	if err != nil {
		return nil, err
	}
	return
}