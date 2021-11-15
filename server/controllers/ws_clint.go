package controllers

import (
	"github.com/galenliu/gateway/pkg/addon"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/container"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/server/models/model"
	"github.com/gofiber/websocket/v2"
)

type wsClint struct {
	ws                   *websocket.Conn
	bus                  controllerBus
	container            model.Container
	thingId              string
	thingCleanups        map[string]func()
	logger               logging.Logger
	subscribedEventNames map[string]bool
}

func NewWsClint(ws *websocket.Conn, bus controllerBus, thingId string, container model.Container, log logging.Logger) *wsClint {
	c := &wsClint{}
	c.bus = bus
	c.subscribedEventNames = make(map[string]bool)
	c.thingCleanups = make(map[string]func())
	c.ws = ws
	c.container = container
	c.logger = log
	c.thingId = thingId
	return c
}

func (c *wsClint) handle() {
	var unsubscribe func()
	if c.thingId == "" {
		things := c.container.GetThings()
		unsubscribe = c.bus.AddThingAddedSubscription(c.onThingAdded)
		for _, t := range things {
			c.addThing(t)
		}
	} else {
		t := c.container.GetThing(c.thingId)
		if t == nil {
			return
		}
		c.addThing(t)
	}

	for {
		mt, data, err := c.ws.ReadMessage()
		if mt == websocket.CloseMessage {
			c.logger.Info("websocket %s close message from ws :", c.ws.LocalAddr())
			return
		}
		if err != nil {
			if unsubscribe != nil {
				unsubscribe()
			}
			return
		}
		go c.handleMessage(data)
	}
}

func (c *wsClint) handleMessage(data []byte) {

}

func (c *wsClint) sendMessage(messageType string, data map[string]interface{}) {
	err := c.ws.WriteJSON(struct {
		MessageType string                 `json:"messageType"`
		Data        map[string]interface{} `json:"data"`
	}{
		MessageType: messageType,
		Data:        data,
	})
	if err != nil {
		c.logger.Error("send message err: %s", err)
	}
}

func (c *wsClint) onThingAdded(t *container.Thing) {
	err := c.ws.WriteJSON(map[string]interface{}{
		"id":          t.GetId(),
		"messageType": constant.ThingAdded,
		"data":        struct{}{},
	})
	c.addThing(t)
	if err != nil {
		c.logger.Error("websocket send %s message err : %s", constant.Connected, err.Error())
	}
}

func (c *wsClint) close() {
	err := c.ws.Close()
	if err != nil {
		c.logger.Infof("%s close", c.ws.LocalAddr().String())
		return
	}
}

func (c *wsClint) addThing(t *container.Thing) {

	onConnected := func(connected bool) {
		err := c.ws.WriteJSON(map[string]interface{}{
			"id":          t.GetId(),
			"messageType": constant.Connected,
			"data":        connected,
		})
		if err != nil {
			c.logger.Error("websocket send %s message err : %s", constant.Connected, err.Error())
		}
	}
	removeConnectedFunc := c.bus.AddConnectedSubscription(t.GetId(), onConnected)

	onThingRemoved := func() {
		f, ok := c.thingCleanups[t.GetId()]
		if ok {
			f()
		}
		err := c.ws.WriteJSON(map[string]interface{}{
			"id":          t.GetId(),
			"messageType": constant.ThingRemoved,
			"data":        struct{}{},
		})
		if err != nil {
		}
		c.logger.Error("websocket send %s message err : %s", constant.ThingRemoved, err.Error())
	}
	removeRemovedFunc := c.bus.AddRemovedSubscription(t.GetId(), onThingRemoved)

	onThingModified := func() {
		err := c.ws.WriteJSON(map[string]interface{}{
			"id":          t.GetId(),
			"messageType": constant.ThingModified,
			"data":        struct{}{},
		})
		if err != nil {
		}
		c.logger.Error("websocket send %s message err : %s", constant.ThingRemoved, err.Error())
	}
	removeModifiedFunc := c.bus.AddModifiedSubscription(t.GetId(), onThingModified)

	onEvent := func(event *addon.Event) {
		err := c.ws.WriteJSON(map[string]interface{}{
			"id":          t.GetId(),
			"messageType": constant.Event,
			"data": struct {
				Name  string
				Event *addon.Event
			}{Name: event.GetName(), Event: event},
		})
		if err != nil {
		}
		c.logger.Error("websocket send %s message err : %s", constant.Event, err.Error())
	}
	removeThingEventStatusFunc := c.bus.AddThingEventSubscription(onEvent)

	onPropertyChanged := func(property *addon.Property) {
		err := c.ws.WriteJSON(map[string]interface{}{
			"id":          t.GetId(),
			"messageType": constant.Event,
			"data":        map[string]interface{}{property.GetName(): property.GetValue()},
		})
		if err != nil {
		}
		c.logger.Error("websocket send %s message err : %s", constant.Event, err.Error())
	}
	removePropertyChangedFunc := c.bus.AddPropertyChangedSubscription(t.GetId(), onPropertyChanged)

	onActionStatus := func(action *addon.Action) {
		err := c.ws.WriteJSON(map[string]interface{}{
			"id":          t.GetId(),
			"messageType": constant.Event,
			"data":        map[string]interface{}{action.GetName(): action.GetDescription()},
		})
		if err != nil {
		}
		c.logger.Error("websocket send %s message err : %s", constant.Event, err.Error())
	}
	removeActionStatusFunc := c.bus.AddActionStatusSubscription(onActionStatus)

	thingCleanup := func() {
		removeConnectedFunc()
		removeRemovedFunc()
		removeModifiedFunc()
		removePropertyChangedFunc()
		removeActionStatusFunc()
		removeThingEventStatusFunc()
	}
	c.thingCleanups[t.GetId()] = thingCleanup
}
