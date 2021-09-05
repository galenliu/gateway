package controllers

import (
	"github.com/galenliu/gateway/server/models"
	"github.com/gofiber/fiber/v2"
	json "github.com/json-iterator/go"
)

type ServiceManager interface {
	EnableService(id string) error
	DisableService(id string) error
	LoadService(id string) error
}

type ServicesController struct {
	model   *models.Services
	manager ServiceManager
}

func NewServicesController(model *models.Services, s ServiceManager) *ServicesController {
	c := &ServicesController{}
	c.manager = s
	c.model = model
	return c
}

func (s ServicesController) handleGetServices(c *fiber.Ctx) error {
	return nil
}

func (s *ServicesController) handleSetService(c *fiber.Ctx) error {
	id := c.FormValue("serviceId")
	if id == "" {
		return fiber.NewError(fiber.StatusInternalServerError, "serviceId failed")
	}
	enabled := json.Get(c.Body(), "enabled").ToBool()
	var err error = nil
	if enabled {
		err = s.manager.DisableService(id)
	}
	err = s.manager.DisableService(id)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}

func (s *ServicesController) handleGetServiceConfig(c *fiber.Ctx) error {
	id := c.FormValue("serviceId", "")
	if id == "" {
		return fiber.NewError(fiber.StatusInternalServerError, "serviceId failed")
	}
	config, err := s.model.LoadConfig(id)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusCreated).JSON(config)
}

func (s *ServicesController) handleSetServiceConfig(c *fiber.Ctx) error {
	id := c.FormValue("serviceId", "")
	if id == "" {
		return fiber.NewError(fiber.StatusInternalServerError, "serviceId failed")
	}
	str := json.Get(c.Body(), "config").ToString()
	info, err := s.model.StoreConfig(id, str)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if info.Enabled {
		err := s.manager.LoadService(id)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}
	return c.SendStatus(fiber.StatusCreated)
}
