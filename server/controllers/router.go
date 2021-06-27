package controllers

import (
	"github.com/galenliu/gateway/configs"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/util"
	"github.com/galenliu/gateway/server"
	"github.com/galenliu/gateway/server/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
	"net/http"
	"strconv"
)

type Router struct {
	*fiber.App
	logger           logging.Logger
	thingsController *thingsController
	userController   userController
}

func NewAPP(thingsModel thingsModel, log logging.Logger) *Router {

	//init
	app := Router{}

	app.logger = log

	app.App = fiber.New()
	app.Use(recover.New())
	app.thingsController = NewThingsController(thingsModel, log)

	//app.Use("/", filesystem.New(filesystem.Options{
	//	Root: http.FS(server.File),
	//}))

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
	thingsController := NewThingsController()
	usersController := NewUsersController()
	addonController := NewAddonController()

	Tid := "/:thingId"

	//Things Controller
	{
		thingsGroup := app.Group(util.ThingsPath)
		//set a properties of a thing.
		thingsGroup.Put("/:thingId/properties/*", thingsController.handleSetProperty)
		thingsGroup.Get("/:thingId/properties/*", thingsController.handleGetPropertyValue)

		//Handle creating a new thing.
		thingsGroup.Post("/", thingsController.handleCreateThing)

		thingsGroup.Get("/:thingId", thingsController.handleGetThing)
		thingsGroup.Get("/", thingsController.handleGetThings)

		thingsGroup.Get("/:thingId", websocket.New(handleWebsocket))
		thingsGroup.Get("/", websocket.New(handleWebsocket))

		//Get the properties of a thing
		thingsGroup.Get(Tid+util.PropertiesPath, thingsController.handleGetProperties)

		thingsGroup.Get(Tid+util.ActionsPath, actionsController.handleGetActions)
		thingsGroup.Post(Tid+util.ActionsPath, actionsController.handleAction)

		// Modify a ThingInfo.
		thingsGroup.Put("/:thingId", thingsController.handleSetThing)

		thingsGroup.Patch("/", thingsController.handlePatchThings)
		thingsGroup.Patch("/:thingId", thingsController.handlePatchThing)

		thingsGroup.Delete("/:thingId", thingsController.handleDeleteThing)

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

	{
		usersGroup := app.Group(util.UsersPath)

		usersGroup.Get("/count", usersController.getCount)
		usersGroup.Post("/", usersController.createUser)
	}

	{ //Addons Controller
		addonGroup := app.Group(util.AddonsPath)
		addonGroup.Get("/", addonController.handlerGetAddons)
		addonGroup.Post("/", addonController.handlerInstallAddon)
		addonGroup.Put("/:addonId", addonController.handlerSetAddon)
		addonGroup.Patch("/:addonId", addonController.handlerUpdateAddon)
		addonGroup.Get("/:addonId/options", addonController.handlerGetAddonConfig)
		addonGroup.Put("/:addonId/options", addonController.handlerSetAddonConfig)
	}

	{ //settings Controller
		debugGroup := app.Group(util.SettingsPath)
		settingsController := NewSettingController()
		debugGroup.Get("/addonsInfo", settingsController.handleGetAddonsInfo)
	}

	{ //actions Controller
		actionsGroup := app.Group(util.ActionsPath)

		actionsGroup.Post("/", actionsController.handleAction)
		actionsGroup.Get("/", actionsController.handleGetActions)
		actionsGroup.Delete("/:actionName/:actionId", actionsController.handleDeleteAction)
	}

	return app
}

func (web *server.WebServe) Start() error {
	httpPort := ":" + strconv.Itoa(web.options.HttpPort)
	var err error
	go func() {
		err = web.Listen(httpPort)
		if err != nil {
			logging.Error("web server err:%s", err.Error())
		}
	}()
	if err != nil {
		event_bus.Publish(util.WebServerStarted)
	}
	return err
}

func (web *server.WebServe) Stop() {
	err := web.App.Shutdown()
	if err != nil {
		logging.Error(err.Error())
	}
	event_bus.Publish(util.WebServerStopped)
}

func NewDefaultWebConfig() server.Options {
	conf := server.Options{
		HttpPort:    configs.GetPorts().HTTP,
		HttpsPort:   configs.GetPorts().HTTPS,
		StaticDir:   "./dist",
		TemplateDir: "./dist",
		UploadDir:   configs.GetUploadDir(),
		LogDir:      configs.GetLogDir(),
	}
	return conf
}
