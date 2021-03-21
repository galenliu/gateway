package controllers

import (
	"gateway/pkg/bus"
	"gateway/pkg/log"
	"gateway/pkg/util"
	"gateway/server/models"
	"gateway/server/models/thing"
	"github.com/gofiber/websocket/v2"
	"sync"
)

type NewThingsController struct {
	locker    *sync.Mutex
	container *models.Things
	ws        *websocket.Conn
	closeChan chan struct{}
}

func NewNewThingsController(ws *websocket.Conn) *NewThingsController {
	controller := &NewThingsController{}
	controller.locker = new(sync.Mutex)
	controller.closeChan = make(chan struct{})
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

	for {
		select {
		case <-controller.closeChan:
			log.Info("new things websocket disconnection")
			return
		default:
			_, message, err := controller.ws.ReadMessage()
			if err != nil {
				log.Error("read:", err)
				return
			}
			log.Info("recv: %s", message)
		}

	}
}

func (controller *NewThingsController) handleNewThing(thing *thing.Thing) {
	controller.locker.Lock()
	defer controller.locker.Unlock()
	err := controller.ws.WriteJSON(thing)
	if err != nil {
		controller.closeChan <- struct{}{}
	}
}

func handleNewThingsWebsocket(conn *websocket.Conn) {
	if !conn.Locals("websocket").(bool) {
		return
	}
	controller := NewNewThingsController(conn)
	controller.handlerConnection()

}
