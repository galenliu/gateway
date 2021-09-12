package ipc_server

import (
	"github.com/galenliu/gateway-grpc"
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

func (c *Clint) Send(message *gateway_grpc.BaseMessage) error {
	return c.Conn.WriteJSON(message)
}

func (c *Clint) Read() (message *gateway_grpc.BaseMessage, err error) {
	err = c.ReadJSON(message)
	if err != nil {
		return nil, err
	}
	return
}
