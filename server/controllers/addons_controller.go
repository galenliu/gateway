package controllers

import (
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/server/models"
	json "github.com/json-iterator/go"
	"net/http"
)

type AddonController struct {
	model  *models.AddonsModel
	logger logging.Logger
}

func NewAddonController(m *models.AddonsModel, log logging.Logger) *AddonController {
	a := &AddonController{}
	a.model = m
	a.logger = log
	return a
}

//  GET /addons
func (addon *AddonController) handlerGetInstallAddons(c *fiber.Ctx) error {
	data := addon.model.GetInstallAddons()
	return c.Status(fiber.StatusOK).Send(data)
}

// Get addonId/license"
func (addon *AddonController) handlerGetLicense(c *fiber.Ctx) error {
	addonId := c.Params("addonId")
	data, err := addon.model.GetAddonLicense(addonId)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Type("text/plain").Status(fiber.StatusOK).SendString(data)
}

// PUT /addons/:id
func (addon *AddonController) handlerSetAddon(c *fiber.Ctx) error {
	addonId := c.FormValue("addonId")
	enabled := json.Get(c.Body(), "enabled").ToBool()
	var err error
	if enabled {
		err = addon.model.EnableAddon(addonId)
	} else {
		err = addon.model.DisableAddon(addonId)
	}
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(enabled)
}

//GET /addon/:addonId/config
func (addon *AddonController) handlerGetAddonConfig(c *fiber.Ctx) error {
	var addonId = c.Params("addonId")
	if addonId == "" {
		return fiber.NewError(fiber.StatusInternalServerError, "addonId failed")
	}
	config, err := addon.model.Store.LoadAddonConfig(addonId)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if config == "" {
		return c.Status(fiber.StatusOK).SendString("{}")
	}
	return c.Status(fiber.StatusOK).SendString(config)
}

//Put /:addonId/config
func (addon *AddonController) handlerSetAddonConfig(c *fiber.Ctx) error {
	var addonId = c.Params("addonId")
	config := json.Get(c.Body(), "config").ToString()
	if config == "" {
		return fiber.NewError(fiber.StatusBadRequest, "config empty")
	}
	err := addon.model.Store.StoreAddonsConfig(addonId, config)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to set config for add-on: "+addonId)
	}
	err = addon.model.UnloadAddon(addonId)
	if addon.model.AddonEnabled(addonId) {
		err := addon.model.LoadAddon(addonId)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to restart add-on: "+addonId)
		}
	}
	err = addon.model.LoadAddon(addonId)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to restart add-on: "+addonId)
	}
	return c.Status(fiber.StatusOK).SendString(config)
}

// Post /addons
func (addon *AddonController) handlerInstallAddon(c *fiber.Ctx) error {

	id := c.FormValue("id")
	url := c.FormValue("url")
	checksum := c.FormValue("checksum")
	if id == "" || url == "" || checksum == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Bad Request"})
	}
	e := addon.model.InstallAddonFromUrl(id, url, checksum)
	if e != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": e.Error()})
	}
	setting, ee := addon.model.Store.LoadAddonSetting(id)
	if ee != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": ee.Error()})
	}
	return c.Status(fiber.StatusOK).SendString(setting)
}

// Patch /:addonId
func (addon *AddonController) handlerUpdateAddon(c *fiber.Ctx) error {
	id := c.Params("addonId")
	url := json.Get(c.Body(), "url").ToString()
	checksum := json.Get(c.Body(), "checksum").ToString()
	if id == "" || url == "" || checksum == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Bad Request"})
	}
	e := addon.model.InstallAddonFromUrl(id, url, checksum)
	if e != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": e.Error()})
	}
	setting, ee := addon.model.Store.LoadAddonSetting(id)
	if ee != nil {

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": ee.Error()})
	}
	return c.Status(fiber.StatusOK).SendString(setting)

}

//Delete /:addonId
func (addon *AddonController) handlerDeleteAddon(c *fiber.Ctx) error {
	var addonId = c.Params("addonId")
	err := addon.model.UninstallAddon(addonId, true)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
