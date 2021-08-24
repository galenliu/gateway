package controllers

import (
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/logging"
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

func NewSettingController(log logging.Logger) *SettingsController {
	return &SettingsController{}
}

func (settings *SettingsController) handleGetAddonsInfo(c *fiber.Ctx) error {
	var addonInfo = addonInfo{
		Urls:          nil,
		Architecture:  util.GetArch(),
		Version:       constant.Version,
		NodeVersion:   util.GetNodeVersion(),
		PythonVersion: util.GetPythonVersion(),
	}
	return c.Status(fiber.StatusOK).JSON(addonInfo)
}
