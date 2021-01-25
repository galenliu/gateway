package controllers

import (
	"fmt"
	"gateway/pkg/bus"
	"gateway/pkg/util"
	"gateway/server/models"
	"gateway/server/models/thing"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type NewThingsController struct {
	locker                 *sync.Mutex
	container              *models.Things
	removeSubscriptionFunc func()
	ws                     *websocket.Conn
	closeChan              chan struct{}
}

func NewNewThingsController(things *models.Things) *NewThingsController {
	controller := &NewThingsController{container: models.NewThings()}
	controller.locker = new(sync.Mutex)
	controller.closeChan = make(chan struct{})
	bus.Subscribe(util.ThingAdded, controller.handleNewThing)
	return controller
}

func (controller *NewThingsController) HandleGetThing(c *gin.Context) {

	c.JSON(http.StatusOK, controller.container.GetThings())
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (controller *NewThingsController) HandleWebsocket(c *gin.Context) {

	if !c.IsWebsocket() {
		c.Next()
	}
	fmt.Printf("is websocket %v", c.IsWebsocket())
	fmt.Printf("c.gin %v", c.Request)

	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.String(http.StatusBadGateway, err.Error())
		return
	}
	controller.ws = conn

	newThings := controller.container.GetNewThings()

	for _, t := range newThings {
		controller.handleNewThing(t)
	}

	for {
		select {
		case <-controller.closeChan:
			return

		default:
			_, data, e := conn.ReadMessage()
			controller.checkErr(e)
			log.Print(data)
		}

	}
}

func (controller *NewThingsController) handleNewThing(thing *thing.Thing) {
	controller.locker.Lock()
	defer controller.locker.Unlock()
	if controller.ws != nil {
		err := controller.ws.WriteJSON(thing)
		controller.checkErr(err)
	}
}

func (controller *NewThingsController) checkErr(err error) {
	if err != nil {
		log.Print(err.Error())
		controller.closeChan <- struct{}{}
		return
	}

}

func (controller *NewThingsController) Close() {
	controller.locker.Lock()
	defer controller.locker.Unlock()
	bus.Unsubscribe(util.ThingAdded, controller.handleNewThing)
	_ = controller.ws.Close()
}
