package controllers

import (
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/server/models"
	"github.com/gofiber/websocket/v2"
	"sync"
)

type NewThingAddonHandler interface {
	GetDevicesBytes() map[string][]byte
}

type NEWThingsController struct {
	locker     *sync.Mutex
	ws         *websocket.Conn
	foundThing chan string
	closeChan  chan struct{}
	logger     logging.Logger
}

func NewNEWThingsController(log logging.Logger) *NEWThingsController {
	c := &NEWThingsController{}
	c.logger = log
	c.locker = new(sync.Mutex)
	c.closeChan = make(chan struct{})
	c.foundThing = make(chan string)
	return c
}

func (controller *NEWThingsController) handleNewThingsWebsocket(thingsModel models.Container, addonHandler NewThingAddonHandler) func(conn *websocket.Conn) {
	return func(conn *websocket.Conn) {
		addonDevices := addonHandler.GetDevicesBytes()
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
