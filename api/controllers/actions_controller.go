package controllers

import (
	"context"
	"fmt"
	"github.com/galenliu/gateway/api/models"
	"github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/util"
	"github.com/gofiber/fiber/v2"
	"time"
)

type ActionsController struct {
	logger         logging.Logger
	thingContainer *container.ThingsContainer
	actions        *models.ActionsModel
	manager        models.ActionsManager
	bus            *bus.ThingsBus
}

func NewActionsController(model *models.ActionsModel, thing *container.ThingsContainer, manager models.ActionsManager, log logging.Logger) *ActionsController {
	return &ActionsController{
		logger:         log,
		manager:        manager,
		thingContainer: thing,
		actions:        model,
	}
}

// Post actions
type action struct {
	Input map[string]any `json:"input,omitempty"`
}

func (a *ActionsController) handleCreateAction(c *fiber.Ctx) error {
	thingId := c.Params("thingId")
	var actionBody map[string]*action
	err := c.BodyParser(&actionBody)
	// 确保一个Action，有且只有一个Input
	if err != nil || len(actionBody) != 1 {
		err := fmt.Errorf("incorrect number of parameters. body:  %s", c.Body())
		a.logger.Error(err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	var actionName string
	var actionParams *action
	for a, params := range actionBody {
		actionName = a
		actionParams = params
	}

	if actionName == "" || actionParams == nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	var thing *container.Thing
	var actionModel *models.Action
	if thingId != "" {
		thing = a.thingContainer.GetThing(thingId)
		if thing == nil {
			return util.NotFoundError("Thing: %s do not exist", thingId)
		}
		actionModel = models.NewActionModel(actionName, actionParams.Input, a.bus, a.logger, thing)
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*3)
		err := a.manager.RequestAction(ctx, thing.GetId(), actionModel.GetName(), actionParams.Input)
		cancelFunc()
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("create actions: %s failed. err: %s", actionName, err.Error()))
		}
	}
	if thing == nil && actionModel == nil {
		actionModel = models.NewActionModel(actionName, actionParams.Input, a.bus, a.logger)
	}
	err = a.actions.Add(actionModel)
	if err != nil {
		return fiber.NewError(fiber.StatusProcessing, fmt.Sprintf("create actions: %s failed，err:%s", actionName, err.Error()))
	}
	//return c.Status(fiber.StatusCreated).JSON(map[string]any{actionName: actionModel.GetDescription()})
	m := make(map[string]models.ActionDescription)
	actionDescription := actionModel.GetDescription()
	m[actionName] = actionDescription
	return c.Status(fiber.StatusCreated).JSON(m)
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
			a.logger.Error("delete actions failed err: %s", actionName)
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}
	err := a.actions.Remove(actionId)
	if err != nil {
		a.logger.Error(err.Error())
		return fiber.NewError(fiber.StatusNoContent)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
