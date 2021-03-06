package controllers

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/things"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/tidwall/gjson"
	"net/http"
	"strings"
)

type thingsModel interface {
	GetThing(thingId string) *things.Thing
	CreateThing(data []byte) (*things.Thing, error)
	RemoveThing(thingId string) error
	GetThings() []*things.Thing
	SetPropertyValue(thingId, propName string, value interface{}) (v interface{}, err error)
	GetPropertyValue(thingId, propName string) (interface{}, error)
	GetPropertiesValue(id string) (map[string]interface{}, error)
}

type thingsController struct {
	model  thingsModel
	logger logging.Logger
}

func NewThingsController(model thingsModel, log logging.Logger) *thingsController {
	tc := &thingsController{}
	tc.model = model
	tc.logger = log
	return tc
}

// POST /things
func (tc *thingsController) handleCreateThing(c *fiber.Ctx) error {

	tc.logger.Debug("Post /thing,Body: \t\n %s", c.Body())
	des, e := tc.model.CreateThing(c.Body())
	if e != nil {
		return fiber.NewError(http.StatusInternalServerError, fmt.Sprintf("create thing err: %v", e.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(des)
}

// DELETE /things/:thingId
func (tc *thingsController) handleDeleteThing(c *fiber.Ctx) error {
	thingId := c.Params("thingId")
	err := tc.model.RemoveThing(thingId)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	tc.logger.Info(fmt.Sprintf("Successfully deleted %v from database", thingId))
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

//PUT /:thing/properties/:propertyName
func (tc *thingsController) handleSetProperty(c *fiber.Ctx) error {

	thingId := c.Params("thingId")
	propName := c.Params("*")
	if thingId == "" || propName == "" {
		return fiber.NewError(fiber.StatusBadRequest, "invalid params")
	}
	value := c.Body()
	v, e := tc.model.SetPropertyValue(thingId, propName, value)
	if e != nil {
		logging.Error("Failed set thing(%s) property:(%s) value:(%s),err:(%s)", thingId, propName, value, e.Error())
		return fiber.NewError(fiber.StatusGatewayTimeout, e.Error())
	}
	data := map[string]interface{}{propName: v}
	return c.Status(fiber.StatusOK).JSON(data)
}

func (tc *thingsController) handleGetPropertyValue(c *fiber.Ctx) error {
	id := c.Params("thingId")
	propName := c.Params("*")
	v, err := tc.model.GetPropertyValue(id, propName)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(map[string]interface{}{propName: v})
}

func (tc *thingsController) handleGetProperties(c *fiber.Ctx) error {
	id := c.Params("thingId")
	m, err := tc.model.GetPropertiesValue(id)
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

	title := strings.Trim(gjson.GetBytes(c.Body(), "title").String(), " ")
	if len(title) == 0 || title == "" {
		return fiber.NewError(http.StatusInternalServerError, "Invalid title")
	}

	selectedCapability := strings.Trim(gjson.GetBytes(c.Body(), "selectedCapability").String(), " ")
	if selectedCapability != "" {
		thing.SelectedCapability = selectedCapability
	}
	return c.Status(fiber.StatusOK).SendString("")
}
