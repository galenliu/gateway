package controllers

import "github.com/gofiber/fiber/v2"

type EventsController struct {
}

func NewEventsController() *EventsController {
	return &EventsController{}
}

func (c *EventsController) handleGetEvents(ctx *fiber.Ctx) error {
	return ctx.JSON(struct{}{})
}

func (c *EventsController) handlerGetEvent(ctx *fiber.Ctx) error {
	eventName := ctx.Params("eventName")
	return ctx.JSON(eventName)
}
