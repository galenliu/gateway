package controllers

import (
	"github.com/galenliu/gateway/server/models"
	"github.com/gofiber/fiber/v2"
)

type ServicesController struct {
	model *models.Services
}

func NewServicesController(model *models.Services) *ServicesController {
	c := &ServicesController{}
	c.model = model
	return c
}

func (s ServicesController) handleGetServices(c *fiber.Ctx) error {
	return nil
}
