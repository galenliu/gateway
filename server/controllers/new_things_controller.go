package controllers

import (
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/util"
	AddonManager "github.com/galenliu/gateway/plugin"
	"github.com/galenliu/gateway/server/models"
	"github.com/gofiber/websocket/v2"
	"sync"
)

type NewThingsController1 struct {
	locker     *sync.Mutex
	container  *models.ThingsModel
	ws         *websocket.Conn
	foundThing chan string
	closeChan  chan struct{}
}

func NewNewThingsController(ws *websocket.Conn) *NewThingsController1 {
	controller := &NewThingsController1{}
	controller.locker = new(sync.Mutex)
	controller.closeChan = make(chan struct{})
	controller.foundThing = make(chan string)
	controller.ws = ws
	return controller
}

func (controller *NewThingsController1) handlerConnection() {

	newThings := controller.container.GetNewThings()
	for _, t := range newThings {
		err := controller.ws.WriteJSON(t)
		if err != nil {
			logging.Error("web socket err: %s", err.Error())
			return
		}
	}
	AddonManager.Subscribe(util.ThingAdded, controller.handleNewThing)
	defer func() {
		err := controller.ws.Close()
		if err != nil {
			logging.Error(err.Error())
		}
		AddonManager.Unsubscribe(util.ThingAdded, controller.handleNewThing)
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
			logging.Info("new things websocket disconnection")
			return
		case s := <-controller.foundThing:
			thing := models.NewThingFromString(s)
			if thing != nil {
				err := controller.ws.WriteJSON(thing)
				if err != nil {
					controller.closeChan <- struct{}{}
				}
			}

		}

	}
}

func (controller *NewThingsController1) handleNewThing(data []byte) {
	controller.foundThing <- string(data)
}

func handleNewThingsWebsocket(conn *websocket.Conn) {
	if !conn.Locals("websocket").(bool) {
		return
	}
	controller := NewNewThingsController(conn)
	controller.handlerConnection()

}
