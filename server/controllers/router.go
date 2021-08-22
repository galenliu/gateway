package controllers

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/server/middleware"
	"github.com/galenliu/gateway/server/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
	"net/http"
	"time"
)

type Storage interface {
	models.UsersStore
	models.ThingsStorage
	models.SettingsStore
	models.JsonwebtokenStore
}

type Config struct {
	HttpAddr  string
	HttpsAddr string
}
type Models struct {
	ThingsModel  *models.Container
	UsersModel   *models.Users
	SettingModel *models.Settings
}

type AddonManagerHandler interface {
	AddonHandler
	ThingsHandler
	NewThingAddonHandler
}

type Router struct {
	*fiber.App
	logger  logging.Logger
	options Config
}

func Setup(config Config, addonManager AddonManagerHandler, store Storage, log logging.Logger) *Router {

	//router init
	app := Router{}
	app.logger = log
	app.options = config
	app.App = fiber.New()
	app.Use(recover.New())
	app.Use(logger.New(logger.ConfigDefault))
	app.Use(cors.New(cors.ConfigDefault))

	//models init
	settingModel := models.NewSettingsModel(store, log)
	usersModel := models.NewUsersModel(store, log)
	jsonwebtokenModel := models.NewJsonwebtokenModel(settingModel, store, log)
	thingsModel := models.NewThingsContainer(store, log)

	//logger
	app.Use(logger.New())

	auth := middleware.NewJWTWare(store)

	app.Use(func(c *fiber.Ctx) error {
		if c.Protocol() == "https" {
			c.Response().Header.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		c.Response().Header.Set("Content-Security-Policy", "")
		return c.Next()
	})

	staticHandler := func(c *fiber.Ctx) error {
		return nil
	}
	app.Use(func(c *fiber.Ctx) error {
		if c.Path() == "/" && c.Accepts("html") == "" {
			return c.Next()
		}
		return staticHandler(c)
	})

	app.Use(func(c *fiber.Ctx) error {
		c.Response().Header.Set("Vary", "Accept")
		c.Response().Header.Set("Access-Control-Allow-Origin", "*")
		c.Response().Header.Set("Access-Control-Allow-Headers",
			"Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Response().Header.Set("Access-Control-Allow-Methods", "GET,HEAD,PUT,PATCH,POST,DELETE")
		return nil
	})

	//ping controller
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusNoContent)
	})

	//root handler
	app.Use(rootHandler)

	//app.Get("/", controllers.RootHandle())
	app.Static("/index.htm", "")

	{
		loginController := NewLoginController(usersModel, jsonwebtokenModel, log)
		app.Post(constant.LoginPath, loginController.handleLogin)

	}

	// Users Controller
	{
		usersController := NewUsersController(usersModel, log)
		usersGroup := app.Group(constant.UsersPath)
		usersGroup.Get("/count", usersController.getCount)
		usersGroup.Post("/", usersController.createUser)
	}

	//Things Controller
	{
		thingsController := NewThingsController(thingsModel, nil, log)
		thingsGroup := app.Group(constant.ThingsPath)

		//set a properties of a thing.
		thingsGroup.Put("/:thingId/properties/*", thingsController.handleSetProperty)
		thingsGroup.Get("/:thingId/properties/*", thingsController.handleGetPropertyValue)

		//Handle creating a new thing.
		thingsGroup.Post("/", auth, thingsController.handleCreateThing)

		thingsGroup.Get("/:thingId", thingsController.handleGetThing)
		thingsGroup.Get("/", thingsController.handleGetThings)

		thingsGroup.Get("/:thingId", websocket.New(handleWebsocket(thingsModel, log)))
		thingsGroup.Get("/", websocket.New(handleWebsocket(thingsModel, log)))

		//Get the properties of a thing
		thingsGroup.Get(constant.ThingIdParam+constant.PropertiesPath, thingsController.handleGetProperties)

		// Modify a ThingInfo.
		thingsGroup.Put("/:thingId", thingsController.handleSetThing)
		thingsGroup.Patch("/", thingsController.handlePatchThings)
		thingsGroup.Patch("/:thingId", thingsController.handlePatchThing)
		thingsGroup.Delete("/:thingId", thingsController.handleDeleteThing)

		actionsController := NewActionsController(log)
		thingsGroup.Get(constant.ThingIdParam+constant.ActionsPath, actionsController.handleGetActions)
		thingsGroup.Post(constant.ThingIdParam+constant.ActionsPath, actionsController.handleAction)
	}

	//NewThing Controller
	{
		newThingsController := NewNEWThingsController(log)
		newThingsGroup := app.Group(constant.NewThingsPath)
		newThingsGroup.Use("/", func(c *fiber.Ctx) error {
			if websocket.IsWebSocketUpgrade(c) {
				c.Locals("websocket", true)
				return c.Next()
			}
			return fiber.ErrUpgradeRequired
		})

		newThingsGroup.Get("/", websocket.New(newThingsController.handleNewThingsWebsocket(thingsModel, addonManager)))
	}

	//Addons Controller
	{
		addonController := NewAddonController(addonManager, store, log)
		addonGroup := app.Group(constant.AddonsPath)
		addonGroup.Get("/", addonController.handlerGetAddons)
		addonGroup.Get("/:addonId/license", addonController.handlerGetLicense)
		addonGroup.Post("/", addonController.handlerInstallAddon)
		addonGroup.Put("/:addonId", addonController.handlerSetAddon)
		addonGroup.Patch("/:addonId", addonController.handlerUpdateAddon)
		addonGroup.Get("/:addonId/config", addonController.handlerGetAddonConfig)
		addonGroup.Put("/:addonId/config", addonController.handlerSetAddonConfig)
	}

	//settings Controller
	{
		settingsController := NewSettingController(log)
		debugGroup := app.Group(constant.SettingsPath)
		debugGroup.Get("/addonsInfo", settingsController.handleGetAddonsInfo)
	}

	//Actions Controller
	{
		actionsController := NewActionsController(log)
		actionsGroup := app.Group(constant.ActionsPath)
		actionsGroup.Post("/", actionsController.handleAction)
		actionsGroup.Get("/", actionsController.handleGetActions)
		actionsGroup.Delete("/:actionName/:actionId", actionsController.handleDeleteAction)
	}

	return &app
}

func (app *Router) Start() error {

	var errs []error
	go func() {

		err := app.Listen(app.options.HttpAddr)
		if err != nil {
			errs[0] = fmt.Errorf("http server err:%s", err.Error())

			return
		}
	}()

	go func() {

		err1 := app.Listen(app.options.HttpsAddr)
		if err1 != nil {
			errs[1] = fmt.Errorf("https server err:%s", err1.Error())

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

func (app *Router) Stop() error {
	return app.App.Shutdown()
}
