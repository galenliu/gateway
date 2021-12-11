package ipc

import (
	"encoding/json"
	"fmt"
	"github.com/fasthttp/websocket"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/util"
)

type connection struct {
	*websocket.Conn
	logger     logging.Logger
	registered bool
	pluginId   string
}

func newConnection(conn *websocket.Conn, log logging.Logger) *connection {
	return &connection{
		Conn:       conn,
		logger:     log,
		registered: false,
		pluginId:   "",
	}
}

func (c *connection) WriteMessage(mt messages.MessageType, data interface{}) error {
	message := BaseMessage{
		MessageType: mt,
		Data:        data,
	}
	err := c.Conn.WriteJSON(message)
	if err != nil {
		return err
	}
	if !c.registered {
		c.logger.Debugf("plugin register message:%s \t\n", util.JsonIndent(message))
	} else {
		c.logger.Debugf("write %s message:%s \t\n", c.getPluginId(), util.JsonIndent(message))
	}
	return nil
}

func (c *connection) ReadMessage() (messages.MessageType, interface{}, error) {
	_, data, err := c.Conn.ReadMessage()
	if err != nil {
		return 0, nil, fmt.Errorf("read message error: %v", err.Error())
	}
	if data == nil {
		return 0, nil, fmt.Errorf("invalid data")
	}
	//var msg rpc.BaseMessage
	var m BaseMessage
	err = json.Unmarshal(data, &m)
	if err != nil {
		return messages.MessageType_MashalERROR, nil, fmt.Errorf("marshal err: %s", err.Error())
	}
	if c.registered {
		c.logger.Debugf("rev %s message :%s \t\n", c.getPluginId(), util.JsonIndent(m))
	}
	return m.MessageType, m.Data, nil
}

func (c *connection) getPluginId() string {
	return c.pluginId
}

func (c *connection) Register(pluginId string) {
	c.registered = true
	c.pluginId = pluginId
}

type BaseMessage struct {
	MessageType messages.MessageType `json:"messageType"`
	Data        interface{}          `json:"data"`
}
