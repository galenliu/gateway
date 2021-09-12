package rpc

import "github.com/galenliu/gateway-grpc"

type Clint interface {
	Send(message *gateway_grpc.BaseMessage) error
	Read() (*gateway_grpc.BaseMessage, error)
}

type PluginHandler interface {
	OnMsg(mt gateway_grpc.MessageType, data []byte) error
}

type PluginServer interface {
	RegisterPlugin(pluginId string, clint Clint) PluginHandler
}
