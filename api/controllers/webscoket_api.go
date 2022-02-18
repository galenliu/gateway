package controllers

import (
	things "github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/gofiber/websocket/v2"
)

type MessageType = string

func handleWebsocket(model things.Container, log logging.Logger) func(conn *websocket.Conn) {
	handler := func(c *websocket.Conn) {
		log.Infof("websocket : %s", c.RemoteAddr())
		thingId, _ := c.Locals("thingId").(string)
		clint := NewWsClint(c, thingId, model, log)
		defer clint.close()
		clint.handle()
	}
	return handler
}
