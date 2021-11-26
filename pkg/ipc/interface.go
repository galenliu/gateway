package ipc

import (
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
)

type Clint interface {
	WriteMessage(message messages.MessageType, data interface{}) error
	ReadMessage() (messages.MessageType, interface{}, error)
	SetPluginId(id string)
	GetPluginId() string
}

type PluginHandler interface {
	OnMsg(mt messages.MessageType, data interface{})
}

type PluginServer interface {
	RegisterPlugin(clint Clint) (PluginHandler, error)
}
