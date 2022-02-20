package ipc

import (
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
)

type Connection interface {
	WriteMessage(message messages.MessageType, data any) error
	ReadMessage() (*BaseMessage, error)
	Register(pluginId string)
}

type PluginHandler interface {
	OnMsg(mt messages.MessageType, data any)
}

type PluginServer interface {
	RegisterPlugin(clint Connection) (PluginHandler, error)
}
