package controllers

import (
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/server/models"
	"github.com/gofiber/fiber/v2"
)

type ActionsController struct {
	logger logging.Logger
	model  *models.ActionsModel
}

func NewActionsController(model *models.ActionsModel, log logging.Logger) *ActionsController {
	return &ActionsController{
		logger: log,
		model:  model,
	}
}

func (a *ActionsController) handleCreateAction(c *fiber.Ctx) error {

	var thingId = c.FormValue("thingId", "")
	var actionBody map[string]interface{}
	err := c.BodyParser(&actionBody)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if len(actionBody) > 1 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	var actionName string
	var actionParams interface{}
	for name, a := range actionBody {
		actionName = name
		actionParams = a
	}
	if actionName == "" || actionParams == nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	var m *models.Action
	if thingId != "" {

	} else {
		inputParams, ok := actionParams.(map[string]interface{})
		if !ok {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		input, ok := inputParams["input"].(map[string]interface{})
		if !ok {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		m = models.NewActionModel(actionName, input, a.logger)
	}
	err = a.model.Add(m)
	if err != nil {
		return err
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
