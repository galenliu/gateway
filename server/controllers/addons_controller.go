package controllers

import (
	"github.com/galenliu/gateway/pkg/database"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/gofiber/fiber/v2"
	json "github.com/json-iterator/go"
	"net/http"
)

type AddonHandler interface {
	GetInstallAddons() []byte
	EnableAddon(addonId string) error
	DisableAddon(addonId string) error
	InstallAddonFromUrl(id, url, checksum string, enabled bool) error
	UnloadAddon(id string) error
	UninstallAddon(id string, disabled bool) error
	AddonEnabled(id string) bool
	LoadAddon(id string) error
}

type AddonController struct {
	handler AddonHandler
	logger  logging.Logger
}

func NewAddonController(addonHandler AddonHandler, log logging.Logger) *AddonController {
	a := &AddonController{}
	a.handler = addonHandler
	a.logger = log
	return a
}

//  GET /addons
func (addon *AddonController) handlerGetAddons(c *fiber.Ctx) error {
	data := addon.handler.GetInstallAddons()
	return c.Status(fiber.StatusOK).Send(data)
}

// PUT /addon/:id
func (addon *AddonController) handlerSetAddon(c *fiber.Ctx) error {
	addonId := c.Params("addonId")
	enabled := json.Get(c.Body(), "enabled").ToBool()
	var err error
	if enabled {
		err = addon.handler.EnableAddon(addonId)
	} else {
		err = addon.handler.DisableAddon(addonId)
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
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Bad Request"})
	}
	e := addon.handler.InstallAddonFromUrl(id, url, checksum, true)
	if e != nil {
		logging.Error(e.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": e.Error()})
	}
	key := "addons." + id
	setting, ee := db.GetSetting(key)
	if ee != nil {
		logging.Error("install add-on err : %s", ee.Error())
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
	e := addon.handler.InstallAddonFromUrl(id, url, checksum, true)
	if e != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": e.Error()})

	}
	key := "addons." + id
	setting, ee := db.GetSetting(key)
	if ee != nil {
		logging.Error(ee.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": ee.Error()})
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

	config, err := db.GetSetting(key)
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
	err := db.SetSetting(key, config)
	if err != nil {
		logging.Error(err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to set options for add-on: "+addonId)
	}
	err = addon.handler.UnloadAddon(addonId)
	if addon.handler.AddonEnabled(addonId) {
		err := addon.handler.LoadAddon(addonId)
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
	err := addon.handler.UninstallAddon(addonId, true)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
