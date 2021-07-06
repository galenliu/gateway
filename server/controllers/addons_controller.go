package controllers

import (
	"github.com/galenliu/gateway/pkg/database"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/plugin"
	"github.com/gofiber/fiber/v2"
	json "github.com/json-iterator/go"
	"net/http"
)

type AddonController struct {
	logger logging.Logger
}

func NewAddonController(log logging.Logger) *AddonController {
	c := &AddonController{}
	c.logger = log
	return c
}

//  GET /addons
func (addon *AddonController) handlerGetAddons(c *fiber.Ctx) error {
	data, err := plugin.GetInstallAddons()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusOK).Send(data)
}

// PUT /addon/:id
func (addon *AddonController) handlerSetAddon(c *fiber.Ctx) error {
	addonId := c.Params("addonId")
	enabled := json.Get(c.Body(), "enabled").ToBool()
	var err error
	if enabled {
		err = plugin.EnableAddon(addonId)
	} else {
		err = plugin.DisableAddon(addonId)
	}
	if err != nil {
		logging.Error(err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(enabled)
}

// Post /addons
func (addon *AddonController) handlerInstallAddon(c *fiber.Ctx) error {

	id := json.Get(c.Body(), "id").ToString()
	url := json.Get(c.Body(), "url").ToString()
	checksum := json.Get(c.Body(), "checksum").ToString()
	if id == "" || url == "" || checksum == "" {
		return fiber.NewError(http.StatusBadRequest, "Bad Request")

	}
	e := plugin.InstallAddonFromUrl(id, url, checksum, true)
	if e != nil {
		logging.Error(e.Error())
		return fiber.NewError(http.StatusInternalServerError, e.Error())
	}
	key := "addons." + id
	setting, ee := database.GetSetting(key)
	if ee != nil {
		logging.Error("install add-on err : %s", ee.Error())
		return fiber.NewError(http.StatusInternalServerError, ee.Error())
	}
	return c.Status(fiber.StatusOK).SendString(setting)

}

// Patch /:addonId
func (addon *AddonController) handlerUpdateAddon(c *fiber.Ctx) error {
	id := c.Params("addonId")
	url := json.Get(c.Body(), "url").ToString()
	checksum := json.Get(c.Body(), "checksum").ToString()
	if id == "" || url == "" || checksum == "" {
		return fiber.NewError(http.StatusBadRequest, "Bad Request")

	}
	e := plugin.InstallAddonFromUrl(id, url, checksum, true)
	if e != nil {
		logging.Error("install add-on err :%s", e.Error())
		return fiber.NewError(http.StatusInternalServerError, e.Error())
	}
	key := "addons." + id
	setting, ee := database.GetSetting(key)
	if ee != nil {
		logging.Error(ee.Error())
		return fiber.NewError(http.StatusInternalServerError, ee.Error())
	}
	return c.Status(fiber.StatusOK).SendString(setting)

}

//GET /addon/:addonId/options
func (addon *AddonController) handlerGetAddonConfig(c *fiber.Ctx) error {
	var addonId = c.Params("addonId")
	var key = "addons.options." + addonId
	if addonId == "" {
		return fiber.NewError(fiber.StatusInternalServerError, "addonId failed")
	}

	config, err := database.GetSetting(key)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if config == "" {
		return c.Status(fiber.StatusOK).SendString("{}")
	}
	return c.Status(fiber.StatusOK).SendString(config)
}

//Put /:addonId/options
func (addon *AddonController) handlerSetAddonConfig(c *fiber.Ctx) error {
	var addonId = c.Params("addonId")
	var key = "addons.options." + addonId
	config := json.Get(c.Body(), "options").ToString()
	if config == "" {
		return fiber.NewError(fiber.StatusBadRequest, "options empty")
	}
	err := database.SetSetting(key, config)
	if err != nil {
		logging.Error(err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to set options for add-on: "+addonId)
	}
	err = plugin.UnloadAddon(addonId)
	if plugin.AddonEnabled(addonId) {
		err := plugin.LoadAddon(addonId)
		if err != nil {
			logging.Error(err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to restart add-on: "+addonId)
		}
	}
	return c.Status(fiber.StatusOK).SendString(config)
}

//Delete /:addonId
func (addon *AddonController) handlerDeleteAddon(c *fiber.Ctx) error {
	var addonId = c.Params("addonId")
	err := plugin.UninstallAddon(addonId, true)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
