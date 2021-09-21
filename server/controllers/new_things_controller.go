package controllers

import (
	"github.com/galenliu/gateway-addon/devices"
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/server/models"
	"github.com/galenliu/gateway/server/models/model"
	"github.com/gofiber/websocket/v2"
	"sync"
)

type busController struct {
	bus.Controller
}

func NewBusController(bus bus.Controller) *busController {
	b := busController{}
	b.Controller = bus
	return &b
}

func (b *busController) SubNewDevice(f func(device devices.Device)) {
	b.Subscribe(constant.DeviceAdded, f)
}
func (b *busController) UnsubNewDevice(f func(device devices.Device)) {
	b.Unsubscribe(constant.DeviceAdded, f)
}

type NewThingsController struct {
	locker     *sync.Mutex
	ws         *websocket.Conn
	foundThing chan string
	closeChan  chan struct{}
	logger     logging.Logger
	model      *models.NewThingsModel
	bus        *busController
}

func NewNewThingsController(model *models.NewThingsModel, bus bus.Controller, log logging.Logger) *NewThingsController {
	c := &NewThingsController{}
	c.bus = NewBusController(bus)
	c.logger = log
	c.model = model
	c.locker = new(sync.Mutex)
	c.closeChan = make(chan struct{})
	c.foundThing = make(chan string)
	return c
}

func (c *NewThingsController) handleNewThingsWebsocket(thingsModel model.Container) func(conn *websocket.Conn) {
	return func(conn *websocket.Conn) {
		devicesMap := c.model.Manager.GetDevicesMaps()
		savedThingsMap := thingsModel.GetMapThings()
		for id, dev := range devicesMap {
			_, ok := savedThingsMap[id]
			if !ok {
				err := conn.WriteJSON(dev)
				if err != nil {
					return
				}
			}
		}
		newDeviceHandler := func(device devices.Device) {
			err := conn.WriteJSON(device)
			if err != nil {
				return
			}
		}
		c.bus.SubNewDevice(newDeviceHandler)
		defer func() {
			c.bus.UnsubNewDevice(newDeviceHandler)
		}()
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}
}
