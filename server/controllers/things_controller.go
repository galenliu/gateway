package controllers

import (
	"github.com/galenliu/gateway/pkg/addon"
	"github.com/galenliu/gateway/pkg/container"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	json "github.com/json-iterator/go"
	"net/http"
	"strings"
)

type ThingsManager interface {
	SetPropertyValue(thingId, propertyName string, value interface{}) (interface{}, error)
	GetPropertyValue(thingId, propertyName string) (interface{}, error)
	GetPropertiesValue(thingId string) (map[string]interface{}, error)
	GetMapOfDevices() map[string]*addon.Device
}

type thingsController struct {
	model   *container.ThingsContainer
	logger  logging.Logger
	manager ThingsManager
}

func NewThingsControllerFunc(manager ThingsManager, model *container.ThingsContainer, log logging.Logger) *thingsController {
	tc := &thingsController{}
	tc.manager = manager
	tc.model = model
	tc.logger = log
	return tc
}

// POST /things
func (tc *thingsController) handleCreateThing(c *fiber.Ctx) error {
	tc.logger.Infof("Post /thing,Body: %s", c.Body())
	des, err := tc.model.CreateThing(c.Body())
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(des)
}

// DELETE /things/:thingId
func (tc *thingsController) handleDeleteThing(c *fiber.Ctx) error {
	thingId := c.Params("thingId")
	err := tc.model.RemoveThing(thingId)
	if err != nil {
		return err
	}
	tc.logger.Infof("Successfully deleted %v from database", thingId)
	return c.SendStatus(http.StatusNoContent)
}

//GET /things/:thingId
func (tc *thingsController) handleGetThing(c *fiber.Ctx) error {
	thingId := c.Params("thingId")
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("thingId", thingId)
		return c.Next()
	}
	t := tc.model.GetThing(thingId)
	if t == nil {
		return fiber.NewError(http.StatusBadRequest, "thing not found")
	}
	return c.Status(fiber.StatusOK).JSON(t)
}

//GET /things
func (tc *thingsController) handleGetThings(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("thingId", "")
		return c.Next()
	}
	ts := tc.model.GetThings()
	return c.Status(fiber.StatusOK).JSON(ts)
}

//patch things
func (tc *thingsController) handlePatchThings(c *fiber.Ctx) error {
	tc.logger.Info("container controller handle patch container")
	return nil
}

//PATCH /things/:thingId
func (tc *thingsController) handlePatchThing(c *fiber.Ctx) error {
	tc.logger.Info("container controller handle patch thing")
	return nil
}

//PUT /:thing/a/:propertyName
func (tc *thingsController) handleSetProperty(c *fiber.Ctx) error {
	thingId := c.Params("thingId")
	propName := c.Params("*")
	if thingId == "" || propName == "" {
		return fiber.NewError(fiber.StatusBadRequest, "invalid params")
	}
	value := c.Body()
	v, err := tc.model.SetThingProperty(thingId, propName, value)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(v)
}

func (tc *thingsController) handleGetPropertyValue(c *fiber.Ctx) error {
	id := c.Params("thingId")
	propName := c.Params("*")
	v, err := tc.manager.GetPropertyValue(id, propName)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(map[string]interface{}{propName: v})
}

func (tc *thingsController) handleGetProperties(c *fiber.Ctx) error {
	id := c.Params("thingId")
	m, err := tc.manager.GetPropertiesValue(id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(m)
}

func (tc *thingsController) handleSetThing(c *fiber.Ctx) error {
	thingId := c.Params("thingId")
	thing := tc.model.GetThing(thingId)
	if thing == nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to retrieve thing(%s)", thingId)
	}
	title := strings.Trim(json.Get(c.Body(), "title").ToString(), " ")
	if len(title) == 0 || title == "" {
		return fiber.NewError(http.StatusInternalServerError, "Invalid title")
	}
	thing.SetTitle(title)
	selectedCapability := strings.Trim(json.Get(c.Body(), "selectedCapability").ToString(), " ")
	if selectedCapability != "" {
		thing.SetSelectedCapability(selectedCapability)
	}
	return c.Status(fiber.StatusOK).SendString("")
}
