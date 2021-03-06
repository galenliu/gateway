package controllers

import (
	"gateway/log"
	"gateway/pkg/bus"
	"gateway/pkg/util"
	"gateway/server/models"
	"gateway/server/models/thing"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type NewThingsController struct {
	locker    *sync.Mutex
	container *models.Things
	wss       map[*websocket.Conn]bool
	closeChan chan struct{}
}

func NewNewThingsController(things *models.Things) *NewThingsController {
	controller := &NewThingsController{container: models.NewThings()}
	controller.locker = new(sync.Mutex)
	controller.wss = make(map[*websocket.Conn]bool, 0)
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

	log.Debug("new thing websocket...add:" + conn.RemoteAddr().String())

	newThings := controller.container.GetNewThings()

	for _, t := range newThings {
		err := conn.WriteJSON(t)
		if err != nil {
			return
		}
	}
	controller.wss[conn] = true
	//go handlerPing(controller, conn)
}

func (controller *NewThingsController) handleNewThing(thing *thing.Thing) {
	controller.locker.Lock()
	defer controller.locker.Unlock()
	for ws, _ := range controller.wss {
		err := ws.WriteJSON(thing)
		if err != nil {
			ws.PingHandler()
			_ = ws.Close()
			delete(controller.wss, ws)
		}
	}
}

func (controller *NewThingsController) checkErr(err error) {
	if err != nil {
		log.Error(err.Error())
		return
	}
}

//func handlerPing(c *NewThingsController, conn *websocket.Conn) {
//	var pingTime = 60 * time.Second
//	var pongTime = 180 * time.Second
//	pingTicker := time.NewTicker(pingTime)
//	pongTicker := time.NewTicker(pongTime)
//	defer pingTicker.Stop()
//	defer pongTicker.Stop()
//	for {
//		select {
//		case <-pingTicker.C:
//			_ = conn.WriteMessage(websocket.PingMessage, []byte{})
//		}
//	}
//}
