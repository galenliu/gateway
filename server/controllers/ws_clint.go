package controllers

import (
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/container"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/server/models/model"
	"github.com/gofiber/websocket/v2"
)

type wsClint struct {
	ws            *websocket.Conn
	container     model.Container
	thingId       string
	thingCleanups map[string]func()
	logger        logging.Logger
}

func NewWsClint(ws *websocket.Conn, thingId string, container model.Container, log logging.Logger) *wsClint {
	c := &wsClint{}
	c.ws = ws
	c.container = container
	c.logger = log
	c.thingId = thingId
	return c
}

func (c *wsClint) handle() {

	if c.thingId == "" {
		things := c.container.GetThings()
		for _, t := range things {
			c.addThing(t)
		}
	}
	for {
		mt, data, err := c.ws.ReadMessage()
		if mt == websocket.CloseMessage {
			c.logger.Info("websocket %s close message from ws :", c.ws.LocalAddr())
			return
		}
		if err != nil {
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

func (c *wsClint) onThingAdded() {

}

func (c *wsClint) close() {
	err := c.ws.Close()
	if err != nil {
		c.logger.Infof("%s close", c.ws.LocalAddr().String())
		return
	}
}

func (c *wsClint) addThing(t *container.Thing) {
	var thingCleanup []func()

	thingCleanup = append(thingCleanup, t.AddSubscription(constant.CONNECTED, func(b bool) {
		err := c.ws.WriteJSON(map[string]interface{}{
			"id":          t.GetId(),
			"messageType": constant.CONNECTED,
			"data":        b,
		})
		if err != nil {
		}
		c.logger.Error("websocket send connected message err : %s", err.Error())
	}))

	thingCleanup = append(thingCleanup, t.AddSubscription(constant.ThingRemoved, func() {
		f, ok := c.thingCleanups[t.GetId()]
		if ok {
			f()
		}
		if c.thingId == "" {
			_ = c.ws.Close()
		} else {
			err := c.ws.WriteJSON(map[string]interface{}{
				"id":          t.GetId(),
				"messageType": constant.ThingRemoved,
				"data":        struct{}{},
			})
			if err != nil {
			}
			c.logger.Error("websocket send ThingRemoved message err : %s", err.Error())
		}
	}))

	thingCleanup = append(thingCleanup, t.AddSubscription(constant.ThingModified, func() {
		err := c.ws.WriteJSON(map[string]interface{}{
			"id":          t.GetId(),
			"messageType": constant.ThingModified,
			"data":        struct{}{},
		})
		if err != nil {
			c.logger.Error("websocket send ThingModified message err : %s", err.Error())
		}
	}))

	thingCleanup = append(thingCleanup, t.AddSubscription(constant.EVENT, func() {
		err := c.ws.WriteJSON(map[string]interface{}{
			"id":          t.GetId(),
			"messageType": constant.ThingModified,
			"data":        struct{}{},
		})
		if err != nil {
			c.logger.Error("websocket send ThingModified message err : %s", err.Error())
		}
	}))

	thingCleanups := func() {
		for _, f := range thingCleanup {
			f()
		}
	}
	c.thingCleanups[t.GetId()] = thingCleanups
}
