package controllers

import (
	"fmt"
	things "github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/addon/actions"
	"github.com/galenliu/gateway/pkg/addon/events"
	"github.com/galenliu/gateway/pkg/bus/topic"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/gofiber/websocket/v2"
)

type wsClint struct {
	ws                   *websocket.Conn
	container            things.Container
	thingId              string
	thingCleanups        map[string]func()
	logger               logging.Logger
	subscribedEventNames map[string]bool
}

func NewWsClint(ws *websocket.Conn, thingId string, container things.Container, log logging.Logger) *wsClint {
	c := &wsClint{}
	c.subscribedEventNames = make(map[string]bool)
	c.thingCleanups = make(map[string]func())
	c.ws = ws
	c.container = container
	c.logger = log
	c.thingId = thingId
	return c
}

func (c *wsClint) handle() error {

	onThingAdded := func(message topic.ThingAddedMessage) {
		c.logger.Infof("OnThingsAdded message: %v", message)
		err := c.ws.WriteJSON(map[string]any{
			"id":          message.ThingId,
			"messageType": constant.ThingAdded,
			"data":        struct{}{},
		})
		if err != nil {
			c.logger.Error("websocket send %s message err : %s", constant.ThingAdded, err.Error())
		}
		thing := c.container.GetThing(message.ThingId)
		if thing != nil {
			c.addThing(thing)
		}
	}

	if c.thingId == "" {
		ts := c.container.GetThings()
		f := c.container.Subscribe(topic.ThingAdded, onThingAdded)
		c.thingCleanups[""] = f
		for _, t := range ts {
			c.addThing(t)
		}
	} else {
		t := c.container.GetThing(c.thingId)
		if t == nil {
			return fmt.Errorf("thing:%s not found in container", c.thingId)
		}
		c.addThing(t)
	}

	for {
		mt, data, err := c.ws.ReadMessage()
		if mt == websocket.CloseMessage || err != nil {
			c.close()
			return fmt.Errorf("websocket  close message")
		}
		go c.handleMessage(data)
	}
}

func (c *wsClint) handleMessage(data []byte) {
	c.logger.Infof("websocket read message: %s", data)
}

func (c *wsClint) close() {
	if c.thingCleanups != nil && len(c.thingCleanups) > 0 {
		for _, f := range c.thingCleanups {
			f()
		}
	}
	if c.ws != nil {
		err := c.ws.Close()
		if err != nil {
			c.logger.Infof("%s close", c.ws.LocalAddr().String())
		}
	}
	c.ws = nil
	c.thingCleanups = nil
}

func (c *wsClint) addThing(t *things.Thing) {

	write := func(id string, messageType string, data any) {
		if c.ws != nil {
			err := c.ws.WriteJSON(struct {
				Id          string `json:"id"`
				MessageType string `json:"messageType"`
				Data        any    `json:"data,omitempty"`
			}{
				id,
				messageType,
				data,
			})
			if err != nil {
				c.logger.Error("websocket error: messageType:%s  err : %s", messageType, err.Error())
				return
			}
		}
	}

	onThingConnected := func(message topic.ThingConnectedMessage) {
		if message.ThingId != t.GetId() {
			return
		}
		write(t.GetId(), constant.Connected, message.Connected)
	}

	removeConnectedFunc := t.AddConnectedSubscription(onThingConnected)
	onThingRemoved := func(message topic.ThingRemovedMessage) {
		c.logger.Infof("onThingRemoved message: %s", message)
		if message.ThingId != t.GetId() {
			return
		}
		f, ok := c.thingCleanups[t.GetId()]
		if ok {
			f()
		}
		delete(c.thingCleanups, t.GetId())
		write(t.GetId(), constant.ThingRemoved, nil)
	}
	removeRemovedFunc := c.container.Subscribe(topic.ThingRemoved, onThingRemoved)

	onThingModified := func(message topic.ThingModifyMessage) {
		if message.ThingId != t.GetId() {
			return
		}
		write(t.GetId(), constant.ThingModified, nil)
	}
	removeModifiedFunc := c.container.Subscribe(topic.ThingModify, onThingModified)

	onEvent := func(e *events.Event) {
		write(t.GetId(), constant.Event, struct {
			Name  string        `json:"name"`
			Event *events.Event `json:"events"`
		}{Name: e.GetName(), Event: e})
	}
	removeThingEventStatusFunc := c.container.Subscribe(topic.ThingEvent, onEvent)

	onPropertyChanged := func(message topic.ThingPropertyChangedMessage) {
		if message.ThingId != t.GetId() {
			return
		}
		write(t.GetId(), constant.PropertyStatus, map[string]any{message.PropertyName: message.Value})
	}
	removePropertyChangedFunc := c.container.Subscribe(topic.ThingPropertyChanged, onPropertyChanged)

	onActionStatus := func(thingId string, action *actions.ActionDescription) {
		if t.GetId() != thingId {
			return
		}
		write(t.GetId(), constant.ActionStatus, map[string]any{action.GetName(): action.GetDescription()})
	}
	removeActionStatusFunc := c.container.Subscribe(topic.ThingActionStatus, onActionStatus)

	thingCleanup := func() {
		removeConnectedFunc()
		removeRemovedFunc()
		removeModifiedFunc()
		removePropertyChangedFunc()
		removeActionStatusFunc()
		removeThingEventStatusFunc()
	}

	var data = make(map[string]any)
	for n, _ := range t.Properties {
		v, err := t.GetPropertyValue(n)
		if err == nil {
			write(t.GetId(), constant.PropertyStatus, map[string]any{n: v})
			c.logger.Infof("data message: %v", data)
			continue
		}
		c.logger.Errorf(err.Error())
	}

	c.thingCleanups[t.GetId()] = thingCleanup
}
