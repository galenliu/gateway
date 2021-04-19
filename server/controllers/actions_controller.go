package controllers

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/galenliu/gateway/plugin"
	"github.com/galenliu/gateway/server/models"
	"github.com/gofiber/fiber/v2"
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

func (controller *ActionsController) handleActions(c *fiber.Ctx) error {

	// POST /actions
	var action *models.Action
	var thingId = c.Params("thingId")
	var actionInfo map[string]struct {
		Input map[string]interface{} `json:"input"`
	}
	err := c.BodyParser(&actionInfo)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Bad Request: "+string(c.Body()))
	}
	var actionName string
	var input map[string]interface{}
	for k, v := range actionInfo {
		actionName = k
		input = v.Input
	}

	if thingId != "" {
		action = models.NewThingAction(thingId, actionName, input)
	} else {
		action = models.NewAction(actionName, input, nil)
	}
	controller.Actions.Add(action)

	var actionDesc string
	actionDesc = action.GetDescription()
	if actionDesc == "" {
		return fiber.NewError(http.StatusBadGateway, "action GetDescription err")
	}
	return c.SendString(actionDesc)
}

func (controller *ActionsController) handleDeleteAction(c *fiber.Ctx) error {

	actionId := c.Params("actionId")
	actionName := c.Params("actionName")
	thingId := c.Params("thingId")

	if thingId != "" {
		err := plugin.RemoveAction(thingId, actionId, actionName)
		if err != nil {
			log.Error(fmt.Sprintf("Removing acotion actionId: %s faild,err: %v", actionId, err))
			return fiber.NewError(http.StatusBadGateway, err.Error())

		}
	}
	err := controller.Actions.Remove(actionId)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusOK)
}
