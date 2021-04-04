// @Title  things_controllers
// @Description  app router
// @Author  liuguilin
// @update  liuguilin

package controllers

import (
	"encoding/json"
	"fmt"
	"gateway/pkg/log"
	AddonManager "gateway/plugin"
	"gateway/server/models"
	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/websocket/v2"
	"github.com/tidwall/gjson"
	"net/http"
	"strings"
)

type ThingsController struct {
	Container *models.Things
}

func NewThingsControllerFunc() *ThingsController {
	tc := &ThingsController{}
	tc.Container = models.NewThings()
	return tc
}

// POST /things
func (tc *ThingsController) handleCreateThing(c *fiber.Ctx) error {

	log.Debug("Post /thing,Body: \t\n %s", c.Body())

	id := gjson.GetBytes(c.Body(), "id").String()
	if len(id) < 1 {
		return fiber.NewError(http.StatusBadRequest, "bad request")

	}
	t := tc.Container.GetThing(id)
	if t != nil {
		return fiber.NewError(http.StatusBadRequest, "thing already added")
	}

	des, e := tc.Container.CreateThing(id, c.Body())
	if e != nil {
		return fiber.NewError(http.StatusInternalServerError, fmt.Sprintf("create thing(%s) err: %v", id, e.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(des)
}

// DELETE /things/:thingId
func (tc *ThingsController) handleDeleteThing(c *fiber.Ctx) error {
	thingId := c.Params("thingId")
	_ = AddonManager.RemoveDevice(thingId)
	err := tc.Container.RemoveThing(thingId)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	log.Info(fmt.Sprintf("Successfully deleted %v from database", thingId))
	return c.SendStatus(http.StatusNoContent)
}

//GET /things/:thingId
func (tc *ThingsController) handleGetThing(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("thingId", c.Params("thingId"))
		return c.Next()
	}

	thingId := c.Params("thingId")
	if thingId == "" {
		return fiber.NewError(http.StatusBadRequest, "thing id invalid")

	}
	t := tc.Container.GetThing(thingId)
	if t == nil {
		return fiber.NewError(http.StatusBadRequest, "thing not found")

	}
	return c.Status(fiber.StatusOK).JSON(t)
}

//GET /things
func (tc *ThingsController) handleGetThings(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	ts := tc.Container.GetListThings()
	data, err := json.MarshalIndent(ts, "", " ")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusOK).SendString(string(data))
}

//patch things
func (tc *ThingsController) handlePatchThings(c *fiber.Ctx) error {
	log.Info("container controller handle patch container")
	return nil
}

//PATCH /things/:thingId
func (tc *ThingsController) handlePatchThing(c *fiber.Ctx) error {
	log.Info("container controller handle patch thing")
	return nil
}

//PUT /:thing/properties/:propertyName
func (tc *ThingsController) handleSetProperty(c *fiber.Ctx) error {

	thingId := c.Params("thingId")
	propName := c.Params("*")
	if thingId == "" || propName == "" {
		return fiber.NewError(fiber.StatusBadRequest, "invalid params")
	}
	value := c.Body()
	prop, e := tc.Container.SetThingProperty(thingId, propName, value)
	if e != nil {
		log.Error("Failed set thing(%s) property:(%s) value:(%s),err:(%s)", thingId, propName, value, e.Error())
		return fiber.NewError(fiber.StatusGatewayTimeout, e.Error())
	}
	newValue := gjson.GetBytes(prop, "value").String()
	data := map[string]interface{}{propName: newValue}
	return c.Status(fiber.StatusOK).JSON(data)
}

func (tc *ThingsController) handleGetProperty(c *fiber.Ctx) error {
	id := c.Params("thingId")
	propName := c.Params("*")
	v, err := AddonManager.GetPropertyValue(id, propName)
	var result = make(map[string]interface{})
	if err != nil {
		log.Info("get property err: %s", err.Error())
	}
	result[propName] = v
	return c.Status(fiber.StatusOK).JSON(result)
}

func (tc *ThingsController) handleGetProperties(c *fiber.Ctx) error {
	id := c.Params("thingId")
	th := tc.Container.GetThing(id)
	if th == nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	var result = make(map[string]interface{})
	for propName, _ := range th.Properties {
		v, err := AddonManager.GetPropertyValue(id, propName)
		if err != nil {
			log.Info("get property err: %s", err.Error())
		}
		result[propName] = v
	}
	return c.Status(fiber.StatusOK).JSON(result)
}

func (tc *ThingsController) handleSetThing(c *fiber.Ctx) error {
	thingId := c.Params("thingId")
	thing := tc.Container.GetThing(thingId)
	if thing == nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to retrieve thing(%s)", thingId)
	}

	title := strings.Trim(gjson.GetBytes(c.Body(), "title").String(), " ")
	if len(title) == 0 || title == "" {
		return fiber.NewError(http.StatusInternalServerError, "Invalid title")
	}

	selectedCapability := strings.Trim(gjson.GetBytes(c.Body(), "selectedCapability").String(), " ")
	if selectedCapability != "" {
		thing.SelectedCapability = selectedCapability
	}
	t := thing.SetTitle(title)
	return c.Status(fiber.StatusOK).SendString(t)
}
