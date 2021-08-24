package controllers

import (
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/gofiber/fiber/v2"
)

type ActionsController struct {
	logger logging.Logger
}

func NewActionsController(log logging.Logger) *ActionsController {
	return &ActionsController{

		logger: log,
	}
}

func (controller *ActionsController) handleAction(c *fiber.Ctx) error {

	// POST /actions
	//var thingId = c.Params("thingId")
	//
	//action := models.NewAction(c.Body(), thingId)
	//if thingId != "" {
	//	t := models.NewThingsOnce().GetThing(thingId)
	//	if t != nil {
	//		err := plugin.RequestAction(thingId, action.ID, action.Name, action.Input)
	//		if err != nil {
	//			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	//		}
	//
	//	} else {
	//		return c.Status(fiber.StatusBadRequest).SendString("thing id invalid")
	//	}
	//}
	//err := controller.Actions.Add(action)
	//if err != nil {
	//	return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	//}
	//
	//data, err := json.MarshalIndent(action, "", " ")
	//if err != nil {
	//	return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	//}
	return c.Status(fiber.StatusCreated).Send(nil)
}

func (controller *ActionsController) handleGetActions(c *fiber.Ctx) error {

	// POST /actions
	var thingId = c.Params("thingId", "")
	//var actionName = c.Params("actionName", "")

	if thingId != "" {
		//actions := controller.Actions.GetAction(thingId, actionName)
		return c.Status(fiber.StatusOK).JSON("")
	} else {
		//actions := controller.Actions.GetGatewayActions(actionName)
		return c.Status(fiber.StatusOK).JSON("")
	}
}

func (controller *ActionsController) handleDeleteAction(c *fiber.Ctx) error {
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
	//err := controller.Actions.Remove(actionId)
	//if err != nil {
	//	return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	//}
	return c.SendStatus(fiber.StatusOK)
}
