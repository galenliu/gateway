package controllers

import (
	"gateway/pkg/database"
	"gateway/pkg/log"
	"gateway/plugin"
	"github.com/gofiber/fiber/v2"
	json "github.com/json-iterator/go"
	"net/http"
)

type AddonController struct {
}

func NewAddonController() *AddonController {
	return &AddonController{}
}

//  GET /addons
func (addon *AddonController) handlerGetAddons(c *fiber.Ctx) error {
	data, err := plugin.GetInstallAddons()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Send(data)
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
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusOK)
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
		log.Error("install add-on err :", e.Error())
		return fiber.NewError(http.StatusInternalServerError, "install addon err:  %v", e.Error())
	}
	key := "addons." + id
	setting, ee := database.GetSetting(key)
	if ee != nil {
		log.Error("install add-on err : %v", ee.Error())
		return fiber.NewError(http.StatusInternalServerError, "install addon err: "+ee.Error())
	}
	return c.SendString(setting)

}

//GET /addon/:addonId/config
func (addon *AddonController) handlerGetAddonConfig(c *fiber.Ctx) error {
	return nil

}

func (addon *AddonController) handlerSetAddonConfig(c *fiber.Ctx) error {
	return nil
}
