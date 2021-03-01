package controllers

import (
	"gateway/pkg/bus"
	"gateway/pkg/log"
	"gateway/pkg/util"
	"gateway/server/models"
	"gateway/server/models/thing"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type NewThingsController struct {
	locker                 *sync.Mutex
	container              *models.Things
	ws                     *websocket.Conn
	closeChan              chan struct{}
}

func NewNewThingsController(things *models.Things) *NewThingsController {
	controller := &NewThingsController{container: models.NewThings()}
	controller.locker = new(sync.Mutex)
	controller.closeChan = make(chan struct{})
	_ = bus.Subscribe(util.ThingAdded, controller.handleNewThing)
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

	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.String(http.StatusBadGateway, err.Error())
		return
	}

	log.Debug("new thing websocket...add:"+conn.RemoteAddr().String())
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
			log.Debug("new thing websocket rev data:", data)
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
		log.Error(err.Error())
		return
	}
}

func (controller *NewThingsController) Close() {
	controller.locker.Lock()
	defer controller.locker.Unlock()
	_ = bus.Unsubscribe(util.ThingAdded, controller.handleNewThing)
	controller.closeChan <- struct{}{}
	_ = controller.ws.Close()
}
