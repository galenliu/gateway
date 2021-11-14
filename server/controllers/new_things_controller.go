package controllers

import (
	"github.com/galenliu/gateway/pkg/addon"
	"github.com/galenliu/gateway/pkg/container"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/gofiber/websocket/v2"
	"sync"
)

type deviceManager interface {
	GetMapOfDevices() map[string]*addon.Device
}
type thingContainer interface {
	GetMapOfThings() map[string]*container.Thing
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

		addonDevices := m.GetMapOfDevices()
		removeFunc := bus.AddDeviceAddedSubscription(func(device *addon.Device) {
			_ = c.handlerNewDevice(conn, device)
		})
		defer func() {
			removeFunc()
			_ = conn.Close()
		}()
		savedThings := t.GetMapOfThings()
		for id, dev := range addonDevices {
			_, ok := savedThings[id]
			if !ok {
				err := c.handlerNewDevice(conn, dev)
				if err != nil {
					return
				}
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

func (c *NewThingsController) handlerNewDevice(conn *websocket.Conn, device *addon.Device) error {
	return conn.WriteJSON(container.AsWebOfThing(device))
}
