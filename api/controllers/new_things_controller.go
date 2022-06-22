package controllers

import (
	things "github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/addon/devices"
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/bus/topic"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/galenliu/gateway/pkg/util"
	"github.com/gofiber/websocket/v2"
	"sync"
	"time"
)

type manager interface {
	GetMapOfDevices() map[string]*devices.Device
	bus.ThingsBus
}
type thingContainer interface {
	GetMapOfThings() map[string]*things.Thing
	GetThing(id string) *things.Thing
}

type NewThingsController struct {
	locker     sync.Mutex
	ws         *websocket.Conn
	foundThing chan string
	closeChan  chan struct{}

	manager        manager
	thingContainer thingContainer
}

func NewNewThingsController(manager manager, thingContainer thingContainer) *NewThingsController {
	c := &NewThingsController{}

	c.manager = manager
	c.thingContainer = thingContainer
	c.locker = sync.Mutex{}
	c.closeChan = make(chan struct{})
	c.foundThing = make(chan string)
	return c
}

func (c *NewThingsController) handleNewThingsWebsocket() func(conn *websocket.Conn) {

	return func(conn *websocket.Conn) {

		addThing := func(msg topic.DeviceAddedMessage) {
			if conn != nil {
				if t := c.thingContainer.GetThing(msg.DeviceId); t != nil {
					return
				}
				c.locker.Lock()
				defer c.locker.Unlock()
				data := util.JsonIndent(things.AsWebOfThing(msg.Device))
				log.Infof("New thing: %s \t\n", data)
				err := conn.WriteMessage(websocket.TextMessage, []byte(data))
				if err != nil {
					log.Infof("new thing websocket err:%s\t\n", err.Error())
					return
				}
			}
		}

		_ = c.manager.Subscribe(topic.DeviceAdded, addThing)
		defer func() {
			c.manager.Unsubscribe(topic.DeviceAdded, addThing)
			_ = conn.Close()
		}()
		go func() {
			addonDevices := c.manager.GetMapOfDevices()
			savedThings := c.thingContainer.GetMapOfThings()
			for id, dev := range addonDevices {
				_, ok := savedThings[id]
				if !ok {
					addThing(topic.DeviceAddedMessage{
						DeviceId: dev.GetId(),
						Device:   *dev,
					})
				}
			}
		}()
		for {
			err := conn.SetReadDeadline(time.Now().Add(30 * time.Second))
			if err != nil {
				return
			}
			_, _, err = conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}
}
