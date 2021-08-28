package controllers

import (
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/server/models"
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
	LoadAddon(id string) error
	UninstallAddon(id string, disabled bool) error
	GetAddonLicense(addonId string) (string, error)
	AddonEnabled(addonId string) bool
}

type AddonController struct {
	handler       AddonHandler
	settingsStore models.SettingsStore
	logger        logging.Logger
}

func NewAddonController(addonHandler AddonHandler, settingsStore models.SettingsStore, log logging.Logger) *AddonController {
	a := &AddonController{}
	a.settingsStore = settingsStore
	a.handler = addonHandler
	a.logger = log
	return a
}

//  GET /addons
func (addon *AddonController) handlerGetInstallAddons(c *fiber.Ctx) error {
	data := addon.handler.GetInstallAddons()
	return c.Status(fiber.StatusOK).Send(data)
}

// Get addonId/license"
func (addon *AddonController) handlerGetLicense(c *fiber.Ctx) error {
	addonId := c.Params("addonId")
	data, err := addon.handler.GetAddonLicense(addonId)
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
		err = addon.handler.EnableAddon(addonId)
	} else {
		err = addon.handler.DisableAddon(addonId)
	}
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(enabled)
}

//GET /addon/:addonId/config
func (addon *AddonController) handlerGetAddonConfig(c *fiber.Ctx) error {
	var addonId = c.Params("addonId")
	var key = "addons.config." + addonId
	if addonId == "" {
		return fiber.NewError(fiber.StatusInternalServerError, "addonId failed")
	}

	config, err := addon.settingsStore.GetSetting(key)
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
	var key = "addons.config." + addonId
	config := json.Get(c.Body(), "config").ToString()
	if config == "" {
		return fiber.NewError(fiber.StatusBadRequest, "config empty")
	}
	err := addon.settingsStore.SetSetting(key, config)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to set config for add-on: "+addonId)
	}
	err = addon.handler.UnloadAddon(addonId)
	if addon.handler.AddonEnabled(addonId) {
		err := addon.handler.LoadAddon(addonId)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to restart add-on: "+addonId)
		}
	}
	err = addon.handler.LoadAddon(addonId)
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
	e := addon.handler.InstallAddonFromUrl(id, url, checksum, true)
	if e != nil {

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": e.Error()})
	}
	key := "addons." + id
	setting, ee := addon.settingsStore.GetSetting(key)
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
	e := addon.handler.InstallAddonFromUrl(id, url, checksum, true)
	if e != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": e.Error()})
	}
	key := "addons." + id
	setting, ee := addon.settingsStore.GetSetting(key)
	if ee != nil {

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": ee.Error()})
	}
	return c.Status(fiber.StatusOK).SendString(setting)

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
