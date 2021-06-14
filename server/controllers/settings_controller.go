package controllers

import (
	"github.com/galenliu/gateway"
	"github.com/galenliu/gateway/configs"
	"github.com/galenliu/gateway/pkg/util"
	"github.com/gofiber/fiber/v2"
)

type addonInfo struct {
	Urls          []string `json:"urls"`
	Architecture  string   `json:"architecture"`
	Version       string   `json:"version"`
	NodeVersion   string   `json:"nodeVersion"`
	PythonVersion []string `json:"pythonVersion"`
}

type SettingsController struct {
}

func NewSettingController() *SettingsController {
	return &SettingsController{}
}

func (settings *SettingsController) handleGetAddonsInfo(c *fiber.Ctx) error {
	var addonInfo = addonInfo{
		Urls:          configs.GetAddonListUrls(),
		Architecture:  util.GetArch(),
		Version:       gateway.Version,
		NodeVersion:   util.GetNodeVersion(),
		PythonVersion: util.GetPythonVersion(),
	}
	return c.Status(fiber.StatusOK).JSON(addonInfo)
}
