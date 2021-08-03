package controllers

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/server"
	"github.com/galenliu/gateway/server/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/gofiber/websocket/v2"
	"net/http"
	"time"
)

type Options struct {
	HttpAddr  string
	HttpsAddr string
}
type Models struct {
	ThingsModel  *models.Container
	UsersModel   *models.Users
	SettingModel *models.Settings
}

type Router struct {
	*fiber.App
	logger  logging.Logger
	options Options
}

func Setup(options Options, addonManager server.AddonManager, store server.Store, log logging.Logger) *Router {

	//router init
	app := Router{}
	app.logger = log
	app.options = options
	app.App = fiber.New()
	app.Use(recover.New())

	//controller init

	//logger
	app.Use(logger.New())

	auth := jwtware.New(jwtware.Config{
		SigningMethod:  "ES256",
		SigningKey: nil,
	})

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

	settingModel := models.NewSettingsModel(store, log)
	usersModel := models.NewUsersModel(store, log)
	jsonwebtokenModel := models.NewJsonwebtokenModel(settingModel, store, log)
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
		thingsController := NewThingsController(models.NewThingsContainer(store, log), nil, log)
		thingsGroup := app.Group(constant.ThingsPath)

		//set a properties of a thing.
		thingsGroup.Put("/:thingId/properties/*", thingsController.handleSetProperty)
		thingsGroup.Get("/:thingId/properties/*", thingsController.handleGetPropertyValue)

		//Handle creating a new thing.
		thingsGroup.Post("/",auth, thingsController.handleCreateThing)

		thingsGroup.Get("/:thingId", thingsController.handleGetThing)
		thingsGroup.Get("/", thingsController.handleGetThings)

		thingsGroup.Get("/:thingId", websocket.New(handleWebsocket))
		thingsGroup.Get("/", websocket.New(handleWebsocket))

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

	//NewThingFromString Controller
	{
		newThingsGroup := app.Group(constant.NewThingsPath)
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

	//Addons Controller
	{
		addonController := NewAddonController(addonManager, log)
		addonGroup := app.Group(constant.AddonsPath)
		addonGroup.Get("/", addonController.handlerGetAddons)
		addonGroup.Get("/:addonId/license", addonController.handlerGetLicense)
		addonGroup.Post("/", addonController.handlerInstallAddon)
		addonGroup.Put("/:addonId", addonController.handlerSetAddon)
		addonGroup.Patch("/:addonId", addonController.handlerUpdateAddon)
		addonGroup.Get("/:addonId/options", addonController.handlerGetAddonConfig)
		addonGroup.Put("/:addonId/options", addonController.handlerSetAddonConfig)
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
