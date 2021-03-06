package controllers

import (
	"fmt"
	"gateway/log"
	"gateway/plugin"
	"gateway/server/models"
	"gateway/server/models/thing"
	"github.com/gin-gonic/gin"
	json "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
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
	controller.Actions.Add(action)

	var actionDesc string
	actionDesc, err = action.GetDescription()
	if err != nil {
		c.String(http.StatusBadGateway, err.Error())
		return
	}
	c.String(http.StatusOK, actionDesc)
}

func (controller *ActionsController) HandleDeleteAction(c *gin.Context) {
	log.Debug(fmt.Sprintf("Handler: %s %s", c.Request.Method, c.Request.URL.String()))
	actionId := c.Param("actionId")
	actionName := c.Param("actionName")
	thingId := c.Param("thingId")

	if thingId != "" {
		err := plugin.RemoveAction(thingId, actionId, actionName)
		if err != nil {
			log.Error(fmt.Sprintf("Removing acotion actionId: %s faild,err: %v", actionId, err))
			c.String(400, err.Error())
			return
		}
	}
	err := controller.Actions.Remove(actionId)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

}
