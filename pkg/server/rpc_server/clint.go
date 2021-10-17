package rpc_server

import (
	rpc "github.com/galenliu/gateway-grpc"
)

type Stream interface {
	Send(message *rpc.BaseMessage) error
	Recv() (message *rpc.BaseMessage, err error)
}

type Clint struct {
	pluginId string
	stream   Stream
}

func NewClint(pluginId string, s Stream) *Clint {
	c := &Clint{}
	c.pluginId = pluginId
	c.stream = s
	return c
}

func (c *Clint) Send(message *rpc.BaseMessage) error {
	return c.stream.Send(message)
}

func (c *Clint) Read() (message *rpc.BaseMessage, err error) {
	return c.stream.Recv()
}


