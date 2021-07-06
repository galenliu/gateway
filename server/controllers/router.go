package controllers

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/util"
	"github.com/galenliu/gateway/server/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
	"net/http"
	"time"
)

type Options struct {
	HttpAddr  string
	HttpsAddr string
}
type Models struct {
	ThingsModel  *models.Things
	UsersModel   *models.Users
	SettingModel *models.Settings
}

type Router struct {
	*fiber.App
	logger             logging.Logger
	options            Options
	HttpRunning        bool
	HttpsRunning       bool
	thingsController   *thingsController
	usersController    *userController
	settingsController *SettingsController
	addonController    *AddonController
}

func Setup(options Options, models Models, log logging.Logger) *Router {

	//router init
	app := Router{}
	app.logger = log
	app.options = options
	app.App = fiber.New()
	app.Use(recover.New())
	app.HttpRunning = false
	app.HttpsRunning = false

	//controller init
	app.thingsController = NewThingsController(models.ThingsModel, log)
	app.usersController = NewUsersController(models.UsersModel, log)
	app.settingsController = NewSettingController(log)
	app.addonController = NewAddonController(log)

	//logger
	app.Use(logger.New())

	//ping controller
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusNoContent)
	})

	//root handler
	app.Use(rootHandler)

	//app.Get("/", controllers.RootHandle())
	app.Static("/index.htm", "")

	actionsController := NewActionsController()

	thingId := "/:thingId"

	//Things Controller
	{
		thingsGroup := app.Group(util.ThingsPath)
		//set a properties of a thing.
		thingsGroup.Put("/:thingId/properties/*", app.thingsController.handleSetProperty)
		thingsGroup.Get("/:thingId/properties/*", app.thingsController.handleGetPropertyValue)

		//Handle creating a new thing.
		thingsGroup.Post("/", app.thingsController.handleCreateThing)

		thingsGroup.Get("/:thingId", app.thingsController.handleGetThing)
		thingsGroup.Get("/", app.thingsController.handleGetThings)

		thingsGroup.Get("/:thingId", websocket.New(handleWebsocket))
		thingsGroup.Get("/", websocket.New(handleWebsocket))

		//Get the properties of a thing
		thingsGroup.Get(thingId+util.PropertiesPath, app.thingsController.handleGetProperties)

		thingsGroup.Get(thingId+util.ActionsPath, actionsController.handleGetActions)
		thingsGroup.Post(thingId+util.ActionsPath, actionsController.handleAction)

		// Modify a ThingInfo.
		thingsGroup.Put("/:thingId", app.thingsController.handleSetThing)

		thingsGroup.Patch("/", app.thingsController.handlePatchThings)
		thingsGroup.Patch("/:thingId", app.thingsController.handlePatchThing)

		thingsGroup.Delete("/:thingId", app.thingsController.handleDeleteThing)

	}

	//NewThingFromString Controller
	{
		newThingsGroup := app.Group(util.NewThingsPath)
		newThingsGroup.Use("/", func(c *fiber.Ctx) error {
			if websocket.IsWebSocketUpgrade(c) {
				c.Locals("websocket", true)
				return c.Next()
			}
			return fiber.ErrUpgradeRequired
		})

		newThingsGroup.Get("/", websocket.New(func(conn *websocket.Conn) {
			handleNewThingsWebsocket(conn)
		}))
	}

	// Users Controller
	{
		usersGroup := app.Group(util.UsersPath)

		usersGroup.Get("/count", app.usersController.getCount)
		usersGroup.Post("/", app.usersController.createUser)
	}

	//Addons Controller
	{
		addonGroup := app.Group(util.AddonsPath)
		addonGroup.Get("/", app.addonController.handlerGetAddons)
		addonGroup.Post("/", app.addonController.handlerInstallAddon)
		addonGroup.Put("/:addonId", app.addonController.handlerSetAddon)
		addonGroup.Patch("/:addonId", app.addonController.handlerUpdateAddon)
		addonGroup.Get("/:addonId/options", app.addonController.handlerGetAddonConfig)
		addonGroup.Put("/:addonId/options", app.addonController.handlerSetAddonConfig)
	}

	//settings Controller
	{
		debugGroup := app.Group(util.SettingsPath)
		debugGroup.Get("/addonsInfo", app.settingsController.handleGetAddonsInfo)
	}

	//Actions Controller
	{
		actionsGroup := app.Group(util.ActionsPath)
		actionsGroup.Post("/", actionsController.handleAction)
		actionsGroup.Get("/", actionsController.handleGetActions)
		actionsGroup.Delete("/:actionName/:actionId", actionsController.handleDeleteAction)
	}

	return &app
}

func (web *Router) Start() error {

	var errs []error
	go func() {
		web.HttpRunning = true
		err := web.Listen(web.options.HttpAddr)
		if err != nil {
			errs[0] = fmt.Errorf("http server err:%s", err.Error())
			web.HttpRunning = false
			return
		}
	}()

	go func() {
		web.HttpRunning = true
		err1 := web.Listen(web.options.HttpsAddr)
		if err1 != nil {
			errs[1] = fmt.Errorf("https server err:%s", err1.Error())
			web.HttpsRunning = false
			return
		}
	}()
	time.Sleep(1 * time.Second)

	if errs[0] != nil && errs[1] != nil {
		return fmt.Errorf("http serve err:%s ,https serve err:%s", errs[0].Error(), errs[1].Error())
	}
	if errs[0] != nil {
		return fmt.Errorf("http serve err:%s ", errs[0].Error())
	}
	if errs[1] != nil {
		return fmt.Errorf("https serve err:%s ", errs[1].Error())
	}
	return nil
}

func (web *Router) Stop() error {
	return web.App.Shutdown()
}
