package controllers

import (
	things "github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/addon/devices"
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/bus/topic"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/gofiber/websocket/v2"
	"sync"
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
	locker         *sync.Mutex
	ws             *websocket.Conn
	foundThing     chan string
	closeChan      chan struct{}
	logger         logging.Logger
	manager        manager
	thingContainer thingContainer
}

func NewNewThingsController(manager manager, thingContainer thingContainer, log logging.Logger) *NewThingsController {
	c := &NewThingsController{}
	c.logger = log
	c.manager = manager
	c.thingContainer = thingContainer
	c.locker = new(sync.Mutex)
	c.closeChan = make(chan struct{})
	c.foundThing = make(chan string)
	return c
}

func (c *NewThingsController) handleNewThingsWebsocket() func(conn *websocket.Conn) {
	return func(conn *websocket.Conn) {
		locker := new(sync.Mutex)
		addThing := func(msg topic.DeviceAddedMessage) {
			if conn != nil && locker != nil {
				locker.Lock()
				defer locker.Unlock()
				if t := c.thingContainer.GetThing(msg.DeviceId); t != nil {
					return
				}
				err := conn.WriteJSON(things.AsWebOfThing(msg.Device))
				if err != nil {
					return
				}
			}

		}
		addonDevices := c.manager.GetMapOfDevices()
		unSub := c.manager.Subscribe(topic.DeviceAdded, addThing)
		defer func() {
			//removeFunc()
			locker = nil
			unSub()
			_ = conn.Close()
		}()
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
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}
}
