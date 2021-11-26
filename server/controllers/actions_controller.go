package controllers

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/container"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/server/models"
	"github.com/gofiber/fiber/v2"
)

type ActionsController struct {
	logger         logging.Logger
	thingContainer *container.ThingsContainer
	model          *models.ActionsModel
	manager        models.ActionsManager
	bus            *bus.Bus
}

func NewActionsController(model *models.ActionsModel, thing *container.ThingsContainer, manager models.ActionsManager, bus *bus.Bus, log logging.Logger) *ActionsController {
	return &ActionsController{
		logger:         log,
		manager:        manager,
		thingContainer: thing,
		model:          model,
		bus:            bus,
	}
}

func (a *ActionsController) handleCreateAction(c *fiber.Ctx) error {

	var thingId = c.FormValue("thingId", "")
	var actionBody map[string]map[string]interface{}
	err := c.BodyParser(&actionBody)

	if err != nil || len(actionBody) != 1 {
		err := fmt.Errorf("incorrect number of parameters. body:  %s", c.Body())
		a.logger.Error(err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var actionName string
	var actionParams map[string]interface{}
	for a, params := range actionBody {
		actionName = a
		actionParams = params
	}

	i, ok := actionParams["input"]
	if actionName == "" || actionParams == nil || !ok {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	input, ok := i.(map[string]interface{})
	if !ok {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	var thing *container.Thing
	if thingId != "" {
		thing = a.thingContainer.GetThing(thingId)
		if thing == nil {
			err := fmt.Errorf("thing does not exist: %s", thingId)
			a.logger.Error(err.Error())
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
	}
	actionModel := models.NewActionModel(actionName, input, thing, a.bus, a.logger)
	if thing != nil {
		err := a.manager.RequestAction(thing.GetId(), actionModel.GetName(), input)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
	}
	err = a.model.Add(actionModel)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	return c.Status(fiber.StatusCreated).Send(nil)
}

func (a *ActionsController) handleGetActions(c *fiber.Ctx) error {

	// POST /actions
	var thingId = c.Params("thingId", "")
	//var actionName = c.Params("actionName", "")

	if thingId != "" {
		//actions := a.actions.GetAction(thingId, actionName)
		return c.Status(fiber.StatusOK).JSON("")
	} else {
		//actions := a.actions.GetGatewayActions(actionName)
		return c.Status(fiber.StatusOK).JSON("")
	}
}

func (a *ActionsController) handleDeleteAction(c *fiber.Ctx) error {
	//
	//actionId := c.Params("actionId")
	//actionName := c.Params("actionName")
	//thingId := c.Params("thingId")
	//
	//if thingId != "" {
	//	err := plugin.RemoveAction(thingId, actionId, actionName)
	//	if err != nil {
	//		logging.Error(fmt.Sprintf("Removing acotion actionId: %s faild,err: %v", actionId, err))
	//		return fiber.NewError(http.StatusBadGateway, err.Error())
	//
	//	}
	//}
	//err := a.actions.Remove(actionId)
	//if err != nil {
	//	return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	//}
	return c.SendStatus(fiber.StatusOK)
}
