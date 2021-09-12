package server

import "github.com/galenliu/gateway-grpc"

type Clint interface {
	Send(message *rpc.BaseMessage) error
	Read() (*rpc.BaseMessage, error)
}

type PluginHandler interface {
	OnMsg(mt rpc.MessageType, data []byte) error
}

type PluginServer interface {
	RegisterPlugin(pluginId string, clint Clint) PluginHandler
}
