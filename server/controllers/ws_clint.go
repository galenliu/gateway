package controllers

import (
	"github.com/galenliu/gateway/pkg/addon"
	"github.com/galenliu/gateway/pkg/bus/topic"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/container"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/gofiber/websocket/v2"
)

type wsClint struct {
	ws                   *websocket.Conn
	bus                  controllerBus
	container            Container
	thingId              string
	thingCleanups        map[string]func()
	logger               logging.Logger
	subscribedEventNames map[string]bool
}

func NewWsClint(ws *websocket.Conn, bus controllerBus, thingId string, container Container, log logging.Logger) *wsClint {
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
		unsubscribe = c.bus.Sub(topic.DeviceAdded, func(deviceId string, device *addon.Device) {
			thing := container.AsWebOfThing(device)
			c.addThing(&thing)
		})
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
	c.logger.Infof("websocket read message: %s", data)
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

func (c *wsClint) close() {
	err := c.ws.Close()
	if err != nil {
		c.logger.Infof("%s close", c.ws.LocalAddr().String())
		return
	}
}

func (c *wsClint) addThing(t *container.Thing) {

	onConnected := func(deviceId string, connected bool) {
		if deviceId != t.GetId() {
			return
		}
		err := c.ws.WriteJSON(map[string]interface{}{
			"id":          t.GetId(),
			"messageType": constant.Connected,
			"data":        connected,
		})
		if err != nil {
			c.logger.Error("websocket send %s message err : %s", constant.Connected, err.Error())
		}
	}
	removeConnectedFunc := c.bus.Sub(topic.ThingConnected, onConnected)

	onThingRemoved := func(thingId string) {
		if thingId != t.GetId() {
			return
		}
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
	removeRemovedFunc := c.bus.Sub(topic.ThingRemoved, onThingRemoved)

	onThingModified := func(thingId string) {
		if thingId != t.GetId() {
			return
		}
		err := c.ws.WriteJSON(map[string]interface{}{
			"id":          t.GetId(),
			"messageType": constant.ThingModified,
			"data":        struct{}{},
		})
		if err != nil {
		}
		c.logger.Error("websocket send %s message err : %s", constant.ThingRemoved, err.Error())
	}
	removeModifiedFunc := c.bus.Sub(topic.ThingModify, onThingModified)

	onEvent := func(event *addon.Event) {
		err := c.ws.WriteJSON(map[string]interface{}{
			"id":          t.GetId(),
			"messageType": constant.Event,
			"data": struct {
				Name  string       `json:"name"`
				Event *addon.Event `json:"event"`
			}{Name: event.GetName(), Event: event},
		})
		if err != nil {
		}
		c.logger.Error("websocket send %s message err : %s", constant.Event, err.Error())
	}
	removeThingEventStatusFunc := c.bus.Sub(topic.ThingEvent, onEvent)

	onPropertyChanged := func(thingId string, property *addon.PropertyDescription) {
		if thingId != t.GetId() {
			return
		}
		err := c.ws.WriteJSON(map[string]interface{}{
			"id":          t.GetId(),
			"messageType": constant.PropertyChanged,
			"data":        map[string]interface{}{property.Name: property.Value},
		})
		if err != nil {
		}
		c.logger.Error("websocket send %s message err : %s", constant.Event, err.Error())
	}
	removePropertyChangedFunc := c.bus.Sub(topic.ThingPropertyChanged, onPropertyChanged)

	onActionStatus := func(thingId string, action *addon.ActionDescription) {
		if t.GetId() != thingId {
			return
		}
		err := c.ws.WriteJSON(map[string]interface{}{
			"id":          t.GetId(),
			"messageType": constant.ActionStatus,
			"data":        map[string]interface{}{action.GetName(): action.GetDescription()},
		})
		if err != nil {
		}
		c.logger.Error("websocket send %s message err : %s", constant.Event, err.Error())
	}
	removeActionStatusFunc := c.bus.Sub(topic.ThingActionStatus, onActionStatus)

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
