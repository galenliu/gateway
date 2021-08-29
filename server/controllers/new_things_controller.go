package controllers

import (
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/server/models"
	"github.com/gofiber/websocket/v2"
	"sync"
)

type NewThingsController struct {
	locker     *sync.Mutex
	ws         *websocket.Conn
	foundThing chan string
	closeChan  chan struct{}
	logger     logging.Logger
	model      *models.NewThingsModel
}

func NewNewThingsController(log logging.Logger) *NewThingsController {
	c := &NewThingsController{}
	c.logger = log
	c.locker = new(sync.Mutex)
	c.closeChan = make(chan struct{})
	c.foundThing = make(chan string)
	return c
}

func (c *NewThingsController) handleNewThingsWebsocket(thingsModel models.Container) func(conn *websocket.Conn) {
	return func(conn *websocket.Conn) {
		addonDevices := c.model.Manager.GetDevicesBytes()
		savedThings := thingsModel.GetMapThings()
		for id, dev := range addonDevices {
			_, ok := savedThings[id]
			if !ok {
				dev, err := models.NewThingFromString(string(dev))
				if err == nil {
					err := conn.WriteJSON(dev)
					if err != nil {
						return
					}
				}
			}
		}
	}
}
