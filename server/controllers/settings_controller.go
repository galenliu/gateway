package controllers

import (
	"gateway/config"
	"gateway/pkg/util"

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
		Urls:          config.GetAddonListUrls(),
		Architecture:  util.GetArch(),
		Version:       util.Version,
		NodeVersion:   util.GetNodeVersion(),
		PythonVersion: util.GetPythonVersion(),
	}
	return c.JSON(addonInfo)
}
