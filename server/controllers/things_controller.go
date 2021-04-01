// @Title  things_controllers
// @Description  app router
// @Author  liuguilin
// @update  liuguilin

package controllers

import (
	"fmt"
	"gateway/pkg/log"
	"gateway/plugin"
	"gateway/server/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	json "github.com/json-iterator/go"
	"go.uber.org/zap"
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

	id := json.Get(c.Body(), "id").ToString()
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
	return c.SendString(des)
}

// DELETE /things/:thingId
func (tc *ThingsController) handleDeleteThing(c *fiber.Ctx) error {
	thingId := c.Params("thingId")
	_ = plugin.RemoveDevice(thingId)
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
		c.Locals("websocket", true)
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
		c.Locals("websocket", true)
		return c.Next()
	}
	log.Debug("GET things")
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
	newValue := json.Get(prop, "value").GetInterface()
	data := map[string]interface{}{propName: newValue}
	return c.Status(fiber.StatusOK).JSON(data)
}

func (tc *ThingsController) handleGetProperty(c *fiber.Ctx) error {
	thingId := c.Params("thingId")
	propName := c.Params("*")
	property := tc.Container.FindThingProperty(thingId, propName)
	if property == nil {
		return fiber.NewError(fiber.StatusBadRequest, "property no found")

	}
	data := map[string]interface{}{propName: property.Value}
	return c.JSON(data)
}

func (tc *ThingsController) handleGetProperties(c *fiber.Ctx) error {
	thingId := c.Params("thing_id")
	//thing := tc.Container.GetThing(thingId)
	var props = make(map[string]interface{})
	//for propName, _ := range thing.Properties {
	//	props[propName] = tc.Container.Manager.GetProperty(thingId, propName)
	//}
	log.Info("container handler:GetProperties", zap.String("thingId", thingId), zap.String("method", "PUT"))
	return c.JSON(props)
}

func (tc *ThingsController) handleSetThing(c *fiber.Ctx) error {
	thingId := c.Params("thingId")
	thing := tc.Container.GetThing(thingId)
	if thing == nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to retrieve thing(%s)", thingId)
	}

	title := strings.Trim(json.Get(c.Body(), "title").ToString(), " ")
	if len(title) == 0 || title == "" {
		return fiber.NewError(http.StatusInternalServerError, "Invalid title")
	}

	selectedCapability := strings.Trim(json.Get(c.Body(), "selectedCapability").ToString(), " ")
	if selectedCapability != "" {
		thing.SelectedCapability = selectedCapability
	}
	t := thing.SetTitle(title)
	return c.Status(fiber.StatusOK).SendString(t)
}
