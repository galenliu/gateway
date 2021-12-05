package controllers

import (
	"context"
	"fmt"
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/container"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/server/models"
	"github.com/gofiber/fiber/v2"
	"time"
)

type ActionsController struct {
	logger         logging.Logger
	thingContainer *container.ThingsContainer
	actions        *models.ActionsModel
	manager        models.ActionsManager
	bus            *bus.Bus
}

func NewActionsController(model *models.ActionsModel, thing *container.ThingsContainer, manager models.ActionsManager, bus *bus.Bus, log logging.Logger) *ActionsController {
	return &ActionsController{
		logger:         log,
		manager:        manager,
		thingContainer: thing,
		actions:        model,
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
	var actionModel *models.Action
	if thingId != "" {
		thing = a.thingContainer.GetThing(thingId)
		if thing == nil {
			err := fmt.Errorf("thing does not exist: %s", thingId)
			a.logger.Error(err.Error())
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		actionModel = models.NewActionModel(actionName, input, a.bus, a.logger, thing)
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*1)
		err := a.manager.RequestAction(ctx, thing.GetId(), actionModel.GetName(), input)
		cancelFunc()
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("create action: %s failed. err: %s", actionName, err.Error()))
		}
	}
	if thing != nil || actionModel == nil {
		actionModel = models.NewActionModel(actionName, input, a.bus, a.logger)
	}
	err = a.actions.Add(actionModel)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("create action: %s failed，err:%s", actionName, err.Error()))
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"actionName": actionModel.GetDescription()})
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
		return c.SendStatus(fiber.StatusOK)
	}
}

func (a *ActionsController) handleDeleteAction(c *fiber.Ctx) error {

	actionId := c.Params("actionId")
	actionName := c.Params("actionName")
	thingId := c.Params("thingId")

	if thingId != "" {
		err := a.manager.RemoveAction(thingId, actionId, actionName)
		if err != nil {
			a.logger.Error("delete action failed err: %s", actionName)
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}
	err := a.actions.Remove(actionId)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
