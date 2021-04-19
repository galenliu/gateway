package controllers

import (
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/galenliu/gateway/pkg/util"
	"github.com/galenliu/gateway/server/models"
	"github.com/gofiber/websocket/v2"
	"sync"
)

type NewThingsController struct {
	locker     *sync.Mutex
	container  *models.Things
	ws         *websocket.Conn
	foundThing chan string
	closeChan  chan struct{}
}

func NewNewThingsController(ws *websocket.Conn) *NewThingsController {
	controller := &NewThingsController{}
	controller.locker = new(sync.Mutex)
	controller.closeChan = make(chan struct{})
	controller.foundThing = make(chan string)
	controller.container = models.NewThings()
	controller.ws = ws
	return controller
}

func (controller *NewThingsController) handlerConnection() {

	newThings := controller.container.GetNewThings()
	for _, t := range newThings {
		err := controller.ws.WriteJSON(t)
		if err != nil {
			log.Error("web socket err: %s", err.Error())
			return
		}
	}
	_ = bus.Subscribe(util.ThingAdded, controller.handleNewThing)
	defer func() {
		controller.ws.Close()
		_ = bus.Unsubscribe(util.ThingAdded, controller.handleNewThing)
	}()

	go func() {
		for {
			_, _, err := controller.ws.ReadMessage()
			if err != nil {
				controller.closeChan <- struct{}{}
				return
			}
		}
	}()

	for {
		select {
		case <-controller.closeChan:
			log.Info("new things websocket disconnection")
			return
		case s := <-controller.foundThing:
			thing := models.NewThing(s)
			if thing != nil {
				err := controller.ws.WriteJSON(thing)
				if err != nil {
					controller.closeChan <- struct{}{}
				}
			}

		}

	}
}

func (controller *NewThingsController) handleNewThing(data string) {
	controller.foundThing <- data
}

func handleNewThingsWebsocket(conn *websocket.Conn) {
	if !conn.Locals("websocket").(bool) {
		return
	}
	controller := NewNewThingsController(conn)
	controller.handlerConnection()

}
