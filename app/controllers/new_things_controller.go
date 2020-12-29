package controllers

import (
	"fmt"
	"gateway/addons"
	"gateway/app/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

type NewThingsController struct {
	things *models.Things
}

func NewNewThingsController(things *models.Things) *NewThingsController {
	return &NewThingsController{things: things}
}

func (newThings *NewThingsController) HandleGetThing(c *gin.Context) {

	newThings.things.GetNewThings()
	things := addons.GetThings()
	if things == nil {
		c.String(http.StatusNotFound, "things nil")
		return
	}
	c.JSON(http.StatusOK, things)
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (newThings *NewThingsController) HandleWebsocket(c *gin.Context) {

	if !c.IsWebsocket() {
		c.Next()
	}
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.String(http.StatusBadGateway, err.Error())
		return
	}
	newThing := newThings.things.GetNewThings()
	newThings.things.RegisterWebsocket(conn)
	for _, v := range newThing {
		err = conn.WriteJSON(v)
	}

	if err != nil {
		c.String(http.StatusBadGateway, fmt.Sprint("websocket err:", err.Error()))

	}

}
