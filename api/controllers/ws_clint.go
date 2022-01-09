package controllers

import (
	things "github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/addon/actions"
	"github.com/galenliu/gateway/pkg/addon/devices"
	"github.com/galenliu/gateway/pkg/addon/events"
	"github.com/galenliu/gateway/pkg/addon/properties"
	"github.com/galenliu/gateway/pkg/bus/topic"
	"github.com/galenliu/gateway/pkg/constant"
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
		ts := c.container.GetThings()
		unsubscribe = c.bus.Sub(topic.DeviceAdded, func(deviceId string, device *devices.Device) {
			thing := things.AsWebOfThing(device)
			c.addThing(&thing)
		})
		for _, t := range ts {
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

func (c *wsClint) sendMessage(messageType string, data map[string]any) {
	err := c.ws.WriteJSON(struct {
		MessageType string         `json:"messageType"`
		Data        map[string]any `json:"data"`
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

func (c *wsClint) addThing(t *things.Thing) {

	onConnected := func(deviceId string, connected bool) {
		if deviceId != t.GetId() {
			return
		}
		err := c.ws.WriteJSON(map[string]any{
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
		err := c.ws.WriteJSON(map[string]any{
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
		err := c.ws.WriteJSON(map[string]any{
			"id":          t.GetId(),
			"messageType": constant.ThingModified,
			"data":        struct{}{},
		})
		if err != nil {
		}
		c.logger.Error("websocket send %s message err : %s", constant.ThingRemoved, err.Error())
	}
	removeModifiedFunc := c.bus.Sub(topic.ThingModify, onThingModified)

	onEvent := func(e *events.Event) {
		err := c.ws.WriteJSON(map[string]any{
			"id":          t.GetId(),
			"messageType": constant.Event,
			"data": struct {
				Name  string        `json:"name"`
				Event *events.Event `json:"events"`
			}{Name: e.GetName(), Event: e},
		})
		if err != nil {
		}
		c.logger.Error("websocket send %s message err : %s", constant.Event, err.Error())
	}
	removeThingEventStatusFunc := c.bus.Sub(topic.ThingEvent, onEvent)

	onPropertyChanged := func(thingId string, property *properties.PropertyDescription) {
		if thingId != t.GetId() {
			return
		}
		err := c.ws.WriteJSON(map[string]any{
			"id":          t.GetId(),
			"messageType": constant.PropertyChanged,
			"data":        map[string]any{*property.Name: property.Value},
		})
		if err != nil {
		}
		c.logger.Error("websocket send %s message err : %s", constant.Event, err.Error())
	}
	removePropertyChangedFunc := c.bus.Sub(topic.ThingPropertyChanged, onPropertyChanged)

	onActionStatus := func(thingId string, action *actions.ActionDescription) {
		if t.GetId() != thingId {
			return
		}
		err := c.ws.WriteJSON(map[string]any{
			"id":          t.GetId(),
			"messageType": constant.ActionStatus,
			"data":        map[string]any{action.GetName(): action.GetDescription()},
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
