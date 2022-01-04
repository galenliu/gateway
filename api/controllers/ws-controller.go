package controllers

import (
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/gofiber/websocket/v2"
)

const (
	MessageTypeSetProperty          = "setProperty"
	MessageTypeRequestAction        = "requestAction"
	MessageTypeAddEventSubscription = "addEventSubscription"
	MessageTypePropertyStatus       = "propertyStatus"
	MessageTypeActionStatus         = "actionStatus"
	MessageTypeEvent                = "events"
)

func handleWebsocket(model Container, bus controllerBus, log logging.Logger) func(conn *websocket.Conn) {
	handler := func(c *websocket.Conn) {
		thingId, _ := c.Locals("thingId").(string)
		clint := NewWsClint(c, bus, thingId, model, log)
		defer clint.close()
		clint.handle()
	}
	return handler
}