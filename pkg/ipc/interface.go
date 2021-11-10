package ipc

import "github.com/galenliu/gateway-grpc"

type Clint interface {
	WriteMessage(message *rpc.BaseMessage) error
	ReadMessage() (*rpc.BaseMessage, error)
	SetPluginId(id string)
	GetPluginId() string
}

type PluginHandler interface {
	OnMsg(mt rpc.MessageType, data []byte)
}

type PluginServer interface {
	RegisterPlugin(clint Clint) (PluginHandler, error)
}
