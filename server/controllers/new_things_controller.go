package controllers

import (
	"github.com/galenliu/gateway/pkg/addon"
	b "github.com/galenliu/gateway/pkg/bus/topic"
	"github.com/galenliu/gateway/pkg/logging"
	container2 "github.com/galenliu/gateway/server/models/container"
	"github.com/gofiber/websocket/v2"
	"sync"
)

type deviceManager interface {
	GetMapOfDevices() map[string]*addon.Device
	GetLanguage() string
}
type thingContainer interface {
	GetMapOfThings() map[string]*container2.Thing
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

func (c *NewThingsController) handleNewThingsWebsocket(m deviceManager, t thingContainer, bus controllerBus) func(conn *websocket.Conn) {
	return func(conn *websocket.Conn) {
		addThing := func(deviceId string, device *addon.Device) {
			err := conn.WriteJSON(container2.AsWebOfThing(device))
			if err != nil {
				return
			}
		}
		addonDevices := m.GetMapOfDevices()
		unSub := bus.Sub(b.DeviceAdded, addThing)
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
