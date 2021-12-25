package controllers

import (
	"context"
	"fmt"
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/bus/topic"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/plugin"
	"github.com/galenliu/gateway/server/middleware"
	"github.com/galenliu/gateway/server/models"
	things "github.com/galenliu/gateway/server/models/container"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
	"net/http"
	"time"
)

type controllerBus interface {
	Pub(topic topic.Topic, args ...any)
	Sub(topic topic.Topic, fun any) func()
}

type Storage interface {
	models.UsersStore
	things.ThingsStorage
	models.SettingsStore
	models.JsonwebtokenStore
	models.AddonStore
}

type Config struct {
	HttpAddr  string
	HttpsAddr string
	AddonUrls []string
}

// Container  Things
type Container interface {
	GetThing(id string) *things.Thing
	GetThings() []*things.Thing
	GetMapOfThings() map[string]*things.Thing
	CreateThing(data []byte) (*things.Thing, error)
	RemoveThing(id string)
	UpdateThing(data []byte) error
}

type Models struct {
	ThingsModel  *Container
	UsersModel   *models.Users
	SettingModel *models.Settings
}

type Router struct {
	*fiber.App
	ctx    context.Context
	logger logging.Logger
	config Config
}

func NewRouter(ctx context.Context, config Config, manager *plugin.Manager, store Storage, bus *bus.Bus, log logging.Logger) *Router {

	//router init
	app := Router{}
	app.logger = log
	app.config = config
	app.ctx = ctx
	app.App = fiber.New()
	app.Use(recover.New())
	app.Use(cors.New(cors.ConfigDefault))
	app.Use(compress.New())

	//models init
	settingModel := models.NewSettingsModel(config.AddonUrls, store, log)
	containerModel := things.NewThingsContainerModel(manager, store, bus, log)
	jwtMiddleware := middleware.NewJWTMiddleware(store, log)
	auth := jwtMiddleware.Auth
	usersModel := models.NewUsersModel(store, log)
	addonModel := models.NewAddonsModel(store, log)
	jsonwebtokenModel := models.NewJsonwebtokenModel(settingModel, store, log)

	actionsModel := models.NewActionsModel(manager, containerModel, bus, log)
	//serviceModel := models.NewServicesModel(deviceManager)
	//newThingsModel := models.NewNewThingsModel(deviceManager, log)

	// Color values
	const (
		cBlack   = "\u001b[90m"
		cRed     = "\u001b[91m"
		cGreen   = "\u001b[92m"
		cYellow  = "\u001b[93m"
		cBlue    = "\u001b[94m"
		cMagenta = "\u001b[95m"
		cCyan    = "\u001b[96m"
		cWhite   = "\u001b[97m"
		cReset   = "\u001b[0m"
	)
	//logger
	app.Use(func(c *fiber.Ctx) error {
		return logger.New(logger.Config{
			Format: fmt.Sprintf("%s| %s  %s|${status} %s| -${latency} %s| ${method} %s| ${path}\n",
				cBlue, c.IP(), cRed, cMagenta, cCyan, cGreen),
			Output: log,
		})(c)
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
		if c.Path() != "/" || c.Accepts("html") == "" {
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
		return c.Next()
	})

	//ping controller
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusNoContent)
	})

	//root actions
	app.Use(rootHandler)

	//app.Get("/", controllers.RootHandle())
	app.Static("/index.htm", "")

	{
		loginController := NewLoginController(usersModel, jsonwebtokenModel, log)
		app.Post(constant.LoginPath, loginController.handleLogin)
	}

	// Users Controller
	{
		usersController := NewUsersController(usersModel, jsonwebtokenModel, log)
		usersGroup := app.Group(constant.UsersPath)
		usersGroup.Get("/count", usersController.getCount)
		usersGroup.Post("/", usersController.createUser)
	}

	actionsController := NewActionsController(actionsModel, containerModel, manager, bus, log)
	//actions Controller
	{
		actionsGroup := app.Group(constant.ActionsPath)
		actionsGroup.Post("/", actionsController.handleCreateAction)
		actionsGroup.Get("/", actionsController.handleGetActions)
		actionsGroup.Delete("/:actionName/:actionId", actionsController.handleDeleteAction)
		actionsGroup.Delete("/:thingId/:actionName/:actionId", actionsController.handleDeleteAction)
	}

	//Things Controller
	{
		thingsController := NewThingsControllerFunc(manager, containerModel, log)
		thingsGroup := app.Group(constant.ThingsPath)

		//set a properties of a thing.
		thingsGroup.Put("/:thingId/properties/*", thingsController.handleSetProperty)
		thingsGroup.Get("/:thingId/properties/*", thingsController.handleGetPropertyValue)

		//Handle creating a new thing.
		thingsGroup.Post("/", thingsController.handleCreateThing)

		thingsGroup.Get("/:thingId", thingsController.handleGetThing)
		thingsGroup.Patch("/:thingId", thingsController.handleUpdateThing)
		thingsGroup.Get("/", thingsController.handleGetThings)

		thingsGroup.Get("/:thingId", websocket.New(handleWebsocket(containerModel, bus, log)))
		thingsGroup.Get("/", websocket.New(handleWebsocket(containerModel, bus, log)))

		//Get the properties of a thing
		thingsGroup.Get("/:thingId"+constant.PropertiesPath, thingsController.handleGetProperties)

		// Modify a ThingInfo.
		thingsGroup.Put("/:thingId", thingsController.handleSetThing)
		thingsGroup.Patch("/", auth, thingsController.handlePatchThings)
		thingsGroup.Patch("/:thingId", thingsController.handlePatchThing)
		thingsGroup.Delete("/:thingId", thingsController.handleDeleteThing)

		thingsGroup.Get("/:thingId"+"/"+constant.ActionsPath, actionsController.handleGetActions)
		thingsGroup.Post("/:thingId"+"/"+constant.ActionsPath, actionsController.handleCreateAction)
	}

	//NewThing Controller
	{
		newThingsController := NewNewThingsController(log)
		newThingsGroup := app.Group(constant.NewThingsPath)
		newThingsGroup.Use("/", func(c *fiber.Ctx) error {
			if websocket.IsWebSocketUpgrade(c) {
				c.Locals("websocket", true)
				return c.Next()
			}
			return fiber.ErrUpgradeRequired
		})
		newThingsGroup.Get("/", websocket.New(newThingsController.handleNewThingsWebsocket(manager, containerModel, bus)))
	}

	//Addons Controller
	{
		addonController := NewAddonController(manager, addonModel, log)
		addonGroup := app.Group(constant.AddonsPath)
		addonGroup.Get("/", addonController.handlerGetInstalledAddons)
		addonGroup.Delete("/:addonId", addonController.handlerDeleteAddon)
		addonGroup.Get("/:addonId/license", addonController.handlerGetLicense)
		addonGroup.Post("/", addonController.handlerInstallAddon)
		addonGroup.Put("/:addonId", addonController.handlerSetAddon)
		addonGroup.Patch("/:addonId", addonController.handlerUpdateAddon)
		addonGroup.Get("/:addonId/config", addonController.handlerGetAddonConfig)
		addonGroup.Put("/:addonId/config", addonController.handlerSetAddonConfig)
	}

	//Settings Controller
	{
		settingsGroup := app.Group(constant.SettingsPath)
		settingsController := NewSettingController(settingModel, log)
		settingsGroup.Get("/addonsInfo", settingsController.handleGetAddonsInfo)
	}

	app.Start()
	return &app
}

func (app *Router) Start() {
	go func() {
		c, cancelFunc := context.WithCancel(app.ctx)
		select {
		case <-c.Done():
			cancelFunc()
			_ = app.Shutdown()
		default:
			err := app.Listen(app.config.HttpAddr)
			if err != nil {
				app.logger.Errorf("http server err:%s", err.Error())
				cancelFunc()
				return
			}
		}
		cancelFunc()
	}()
	time.Sleep(1 * time.Millisecond)

	go func() {
		c, cancelFunc := context.WithCancel(app.ctx)
		select {
		case <-c.Done():
			cancelFunc()
			_ = app.Shutdown()
		default:
			err := app.Listen(app.config.HttpsAddr)
			if err != nil {
				app.logger.Errorf("https server err:%s", err.Error())
				cancelFunc()
				return
			}
		}
		cancelFunc()
	}()

}
