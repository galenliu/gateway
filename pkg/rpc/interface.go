package rpc

type Clint interface {
	Send(message *BaseMessage) error
	Read() (*BaseMessage, error)
}

type PluginHandler interface {
	MessageHandler(mt MessageType, data []byte) error
}

type PluginServer interface {
	RegisterPlugin(pluginId string, clint Clint) PluginHandler
}
