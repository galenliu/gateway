package controllers

import (
	"context"
	"github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/addon/devices"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/util"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	json "github.com/json-iterator/go"
	"net/http"
	"strings"
	"time"
)

type ThingsManager interface {
	SetPropertyValue(ctx context.Context, thingId, propertyName string, value any) (any, error)
	GetPropertyValue(thingId, propertyName string) (any, error)
	GetPropertiesValue(thingId string) (map[string]any, error)
	GetMapOfDevices() map[string]*devices.Device
	SetPIN(ctx context.Context, thingId string, pin string) (*messages.Device, error)
	SetCredentials(ctx context.Context, thingId, username, password string) (*messages.Device, error)
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

// POST /things  create a new thing
func (tc *thingsController) handleCreateThing(c *fiber.Ctx) error {
	tc.logger.Debug("Post /things:\t\n %s", c.Body())
	thing, err := tc.model.CreateThing(c.Body())
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusOK).SendString(util.JsonIndent(thing))
}

// DELETE /things/:thingId
func (tc *thingsController) handleDeleteThing(c *fiber.Ctx) error {
	thingId := c.Params("thingId")
	if thingId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "thing id must be provided")
	}
	tc.model.RemoveThing(thingId)
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
	if len(ts) == 0 {
		return c.Status(fiber.StatusOK).JSON([]string{})
	}
	return c.Status(fiber.StatusOK).SendString(util.JsonIndent(ts))
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
	var value any
	err := json.Unmarshal(c.Body(), &value)
	if err != nil {
		return err
	}
	v, err := tc.model.SetThingPropertyValue(thingId, propName, value)
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
	return c.Status(fiber.StatusOK).JSON(map[string]any{propName: v})
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

// Patch thing
func (tc *thingsController) handleUpdateThing(c *fiber.Ctx) error {

	ctx, cancelFunc := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancelFunc()
	thingId := c.Params("thingId")
	username := json.Get(c.Body(), "username").ToString()
	password := json.Get(c.Body(), "password").ToString()
	pin := json.Get(c.Body(), "pin").ToString()

	if thingId != "" && pin != "" {
		device, err := tc.manager.SetPIN(ctx, thingId, pin)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return c.JSON(device)
	}
	if thingId != "" && username != "" && password != "" {
		device, err := tc.manager.SetCredentials(ctx, thingId, username, password)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return c.JSON(device)
	}
	return fiber.NewError(fiber.StatusBadRequest)
}
