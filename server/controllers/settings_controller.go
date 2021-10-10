package controllers

import (
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/server/models"
	"github.com/gofiber/fiber/v2"
)

type SettingsController struct {
	model *models.Settings
}

func NewSettingController(model *models.Settings, log logging.Logger) *SettingsController {
	s := &SettingsController{}
	s.model = model
	return s
}

// GET  /settings/addonInfo
func (s *SettingsController) handleGetAddonsInfo(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(s.model.GetAddonInfo())
}
