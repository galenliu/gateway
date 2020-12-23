package controllers

import (
	"fmt"
	"gateway/app/models"
	"gateway/pkg/log"
	"github.com/gin-gonic/gin"
	json "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"strconv"
)

type ActionsController struct {
	Actions *models.Actions
}

func NewActionsController() *ActionsController {
	return &ActionsController{
		Actions: models.NewActions(),
	}
}

func (controller *ActionsController) HandleActions(c *gin.Context) {

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	log.Debug(fmt.Sprintf("Handler: %s %s body: %s", c.Request.Method, c.Request.URL.String(), data))

	var action *models.Action
	var input json.Any

	if pair := json.Get(data, "pair").ToString(); pair != "" {
		input = json.Get(data, "pair", "input")
		if input == nil {
			c.String(http.StatusBadRequest, "action input err")
			return
		}
	}
	thingId := c.Param("thingId")

	if thingId != "" {
		action = models.NewThingAction()
	} else {
		action = models.NewAction("pair", input)
	}
	controller.Actions.AddAction(action)

	var dataDesc []byte
	dataDesc, err = action.GetDescription()
	if err != nil {
		c.String(http.StatusBadGateway, err.Error())
		return
	}
	c.String(http.StatusOK, string(dataDesc))
}

func (controller *ActionsController) HandleDeleteAction(c *gin.Context) {
	log.Debug(fmt.Sprintf("Handler: %s %s", c.Request.Method, c.Request.URL.String()))
	actionId := c.Param("actionId")
	id, err := strconv.Atoi(actionId)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("action id  must int,err:%v", err.Error()))
	}
	err = controller.Actions.RemoveAction(id)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

}
