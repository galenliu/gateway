package controllers

import (
	"fmt"
	"gateway/server/models"
	"gateway/server/models/thing"
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

	// POST /actions

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	log.Debug(fmt.Sprintf("Handler: %s %s body: %s", c.Request.Method, c.Request.URL.String(), data))

	var action *thing.Action
	var thingId = c.Param("thingId")

	var actionInfo map[string]struct {
		Input map[string]interface{} `json:"input"`
	}
	err = json.Unmarshal(data, &actionInfo)
	if err != nil {
		return
	}

	var actionName string
	var input map[string]interface{}
	for k, v := range actionInfo {
		actionName = k
		input = v.Input
	}

	if thingId != "" {
		action = thing.NewThingAction(thingId, actionName, input)
	} else {
		action = thing.NewAction(actionName, input)
	}
	controller.Actions.AddAction(action)

	var actionDesc []byte
	actionDesc, err = action.GetDescription()
	if err != nil {
		c.String(http.StatusBadGateway, err.Error())
		return
	}
	c.String(http.StatusOK, string(actionDesc))
}

func (controller *ActionsController) HandleDeleteAction(c *gin.Context) {
	log.Debug(fmt.Sprintf("Handler: %s %s", c.Request.Method, c.Request.URL.String()))
	actionId := c.Param("actionId")
	id, err := strconv.Atoi(actionId)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("actions id  must int,err:%v", err.Error()))
	}
	err = controller.Actions.RemoveAction(uint(id))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

}
