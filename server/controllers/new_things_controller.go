package controllers

import (
	"gateway/pkg/bus"
	"gateway/pkg/log"
	"gateway/pkg/util"
	"gateway/server/models"
	"gateway/server/models/thing"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"sync"
)

type NewThingsController struct {
	locker    *sync.Mutex
	container *models.Things
	ws        *websocket.Conn
	closeChan chan struct{}
}

func NewNewThingsController() *NewThingsController {
	controller := &NewThingsController{container: models.NewThings()}
	controller.locker = new(sync.Mutex)
	controller.closeChan = make(chan struct{})
	controller.container = models.NewThings()
	_ = bus.Subscribe(util.ThingAdded, controller.handleNewThing)
	return controller
}

func (controller *NewThingsController) HandleGetThing(c *fiber.Ctx) error {
	return c.JSON(controller.container.GetThings())
}

func handleNewThingsWebsocket(conn *websocket.Conn) {
	if !conn.Locals("websocket").(bool) {
		return
	}
	controller := NewNewThingsController()
	controller.ws = conn

	log.Debug("new thing websocket...add:" + conn.RemoteAddr().String())
	newThings := controller.container.GetNewThings()

	for _, t := range newThings {
		err := conn.WriteJSON(t)
		if err != nil {
			return
		}
	}
	for {
		select {
		case <-controller.closeChan:
			controller.close()
			return
		}

	}
}

func (controller *NewThingsController) handleNewThing(thing *thing.Thing) {
	controller.locker.Lock()
	defer controller.locker.Unlock()
	err := controller.ws.WriteJSON(thing)
	controller.checkErr(err)
}

func (controller *NewThingsController) checkErr(err error) {
	if err != nil {
		log.Error(err.Error())
		controller.closeChan <- struct{}{}
		return
	}
}

func (controller *NewThingsController) close() {
	_ = controller.ws.Close()
}
