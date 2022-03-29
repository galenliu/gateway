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
	"net/url"
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
	tc.logger.Debugf("Post /things:\t\n %s", c.Body())
	thing, err := tc.model.CreateThing(c.Body())
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusOK).SendString(util.JsonIndent(thing))
}

// DELETE /things/:thingId
func (tc *thingsController) handleDeleteThing(c *fiber.Ctx) error {
	thingId, err := url.QueryUnescape(c.Params("thingId"))
	if err != nil {
		return util.BadRequestError("Invalid thing id")
	}
	err = tc.model.RemoveThing(thingId)
	if err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

//GET /things/:thingId
func (tc *thingsController) handleGetThing(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	thingId, err := url.QueryUnescape(c.Params("thingId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest)
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

//PUT /:thing/properties/:propertyName
func (tc *thingsController) handleSetProperty(c *fiber.Ctx) error {
	thingId, err := url.QueryUnescape(c.Params("thingId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest)
	}
	propName := c.Params("*")
	if thingId == "" || propName == "" {
		return fiber.NewError(fiber.StatusBadRequest, "invalid params")
	}
	var value any
	err = json.Unmarshal(c.Body(), &value)
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
	id, err := url.QueryUnescape(c.Params("thingId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest)
	}
	propName := c.Params("*")
	v, err := tc.manager.GetPropertyValue(id, propName)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(map[string]any{propName: v})
}

func (tc *thingsController) handleGetProperties(c *fiber.Ctx) error {
	id, err := url.QueryUnescape(c.Params("thingId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest)
	}
	m, err := tc.manager.GetPropertiesValue(id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(m)
}

// Put /things/:thingId
func (tc *thingsController) handleSetThing(c *fiber.Ctx) error {
	thingId, _ := url.QueryUnescape(c.Params("thingId"))
	thing := tc.model.GetThing(thingId)
	if thing == nil {
		return util.NotFoundError("thing %s not found", thingId)
	}
	title := strings.Trim(json.Get(c.Body(), "title").ToString(), " ")
	if len(title) == 0 || title == "" {
		return util.BadRequestError("Invalid title")
	}
	selectedCapability := strings.Trim(json.Get(c.Body(), "selectedCapability").ToString(), " ")
	if selectedCapability != "" {
		thing.SetSelectedCapability(selectedCapability)
	}
	thing.SetTitle(title)
	return c.SendStatus(fiber.StatusOK)
}

// Patch thing
func (tc *thingsController) handleUpdateThing(c *fiber.Ctx) error {

	ctx, cancelFunc := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancelFunc()
	thingId, err := url.QueryUnescape(c.Params("thingId"))
	if err != nil {
		return util.BadRequestError("Invalid thing id")
	}
	username, err := url.QueryUnescape(json.Get(c.Body(), "username").ToString())
	if err != nil {
		return util.BadRequestError("Invalid username")
	}
	password, err := url.QueryUnescape(json.Get(c.Body(), "password").ToString())
	if err != nil {
		return util.BadRequestError("Invalid password")
	}
	pin := json.Get(c.Body(), "pin").ToString()
	if thingId != "" && pin != "" {
		device, err := tc.manager.SetPIN(ctx, thingId, pin)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(device)
	}
	if thingId != "" && username != "" && password != "" {
		device, err := tc.manager.SetCredentials(ctx, thingId, username, password)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(device)
	}
	return fiber.NewError(fiber.StatusBadRequest)
}
