package controllers

import (
	things "github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/addon/devices"
	"github.com/galenliu/gateway/pkg/bus"
	b "github.com/galenliu/gateway/pkg/bus/topic"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/gofiber/websocket/v2"
	"sync"
)

type deviceManager interface {
	GetMapOfDevices() map[string]*devices.Device
	bus.Bus
}
type thingContainer interface {
	GetMapOfThings() map[string]*things.Thing
}

type NewThingsController struct {
	locker     *sync.Mutex
	ws         *websocket.Conn
	foundThing chan string
	closeChan  chan struct{}
	logger     logging.Logger
}

func NewNewThingsController(log logging.Logger) *NewThingsController {
	c := &NewThingsController{}
	c.logger = log
	c.locker = new(sync.Mutex)
	c.closeChan = make(chan struct{})
	c.foundThing = make(chan string)
	return c
}

func (c *NewThingsController) handleNewThingsWebsocket(m deviceManager, t thingContainer) func(conn *websocket.Conn) {
	return func(conn *websocket.Conn) {
		addThing := func(deviceId string, device *devices.Device) {
			err := conn.WriteJSON(things.AsWebOfThing(device))
			if err != nil {
				return
			}
		}
		addonDevices := m.GetMapOfDevices()
		unSub := m.Subscribe(b.DeviceAdded, addThing)
		defer func() {
			//removeFunc()
			unSub()
			_ = conn.Close()
		}()
		savedThings := t.GetMapOfThings()
		for id, dev := range addonDevices {
			_, ok := savedThings[id]
			if !ok {
				addThing(dev.GetId(), dev)
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
