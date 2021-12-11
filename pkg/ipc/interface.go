package ipc

import (
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
)

type Connection interface {
	WriteMessage(message messages.MessageType, data interface{}) error
	ReadMessage() (messages.MessageType, interface{}, error)
	Register(pluginId string)
}

type PluginHandler interface {
	OnMsg(mt messages.MessageType, data interface{})
}

type PluginServer interface {
	RegisterPlugin(clint Connection) (PluginHandler, error)
}
