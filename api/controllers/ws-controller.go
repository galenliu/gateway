package controllers

import (
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/gofiber/websocket/v2"
)

func handleWebsocket(model Container, bus bus.Bus, log logging.Logger) func(conn *websocket.Conn) {
	handler := func(c *websocket.Conn) {
		thingId, _ := c.Locals("thingId").(string)
		clint := NewWsClint(c, thingId, bus, model, log)
		defer clint.close()
		clint.handle()
	}
	return handler
}
