package controllers

import (
	"gateway/config"
	"gateway/pkg/log"
	"gateway/pkg/util"
	//"gateway/server"
	"gateway/server/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
	"net/http"
	"strconv"
)

type Config struct {
	HttpPort    int
	HttpsPort   int
	StaticDir   string
	TemplateDir string
	UploadDir   string
	LogDir      string
}

type Web struct {
	*fiber.App
	config Config
	things *models.Things
}

func NewWebAPP(conf Config, opts ...webOption) *Web {
	options := defaultConfig
	for _, o := range opts {
		o(&options)
	}
	web := Web{}
	web.things = models.NewThings()
	web.config = conf
	web.App = CollectRoute(conf)

	return &web
}

func CollectRoute(conf Config) *fiber.App {

	//init
	var app = fiber.New()
	app.Use(recover.New())

	//app.Use("/", filesystem.New(filesystem.Config{
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

	//Things Controller
	{
		thingsGroup := app.Group(util.ThingsPath)

		thingsController := NewThingsControllerFunc()
		//set a properties of a thing.
		thingsGroup.Put("/:thingId/properties/*", thingsController.handleSetProperty)
		thingsGroup.Get("/:thingId/properties/*", thingsController.handleGetProperty)

		//Handle creating a new thing.
		thingsGroup.Post("/", thingsController.handleCreateThing)

		thingsGroup.Get("/:thingId", thingsController.handleGetThing)
		thingsGroup.Get("/", thingsController.handleGetThings)

		thingsGroup.Get("/:thingId", websocket.New(handleWebsocket))
		thingsGroup.Get("/", websocket.New(handleWebsocket))

		//Get the properties of a thing
		thingsGroup.Get("/:thingId/properties", thingsController.handleGetProperties)

		// Modify a ThingInfo.
		thingsGroup.Put("/:thingId", thingsController.handleSetThing)

		thingsGroup.Patch("/", thingsController.handlePatchThings)
		thingsGroup.Patch("/:thingId", thingsController.handlePatchThing)

		thingsGroup.Delete("/:thingId", thingsController.handleDeleteThing)

	}

	//NewThing Controller
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

	{ //Addons Controller
		addonGroup := app.Group(util.AddonsPath)
		addonController := NewAddonController()
		addonGroup.Get("/", addonController.handlerGetAddons)
		addonGroup.Post("/", addonController.handlerInstallAddon)
		addonGroup.Put("/:addonId", addonController.handlerSetAddon)
		addonGroup.Patch("/:addonId", addonController.handlerUpdateAddon)
		addonGroup.Get("/:addonId/config", addonController.handlerGetAddonConfig)
		addonGroup.Put("/:addonId/config", addonController.handlerSetAddonConfig)
	}

	{ //settings Controller
		debugGroup := app.Group(util.SettingsPath)
		settingsController := NewSettingController()
		debugGroup.Get("/addonsInfo", settingsController.handleGetAddonsInfo)
	}

	{ //actions Controller
		actionsGroup := app.Group(util.ActionsPath)
		actionsController := NewActionsController()
		actionsGroup.Post("/", actionsController.handleActions)
		actionsGroup.Delete("/:actionName/:actionId", actionsController.handleDeleteAction)
	}

	return app
}

func (web *Web) Start() error {
	httpPort := ":" + strconv.Itoa(web.config.HttpPort)
	err := web.Listen(httpPort)
	if err != nil {
		log.Error("web err:%s", err.Error())
		return err
	}
	return nil
}

func (web *Web) Close() error {
	return web.App.Shutdown()
}

func NewDefaultWebConfig() Config {
	conf := Config{
		HttpPort:    config.GetPorts().HTTP,
		HttpsPort:   config.GetPorts().HTTPS,
		StaticDir:   "./dist",
		TemplateDir: "./dist",
		UploadDir:   config.GetUploadDir(),
		LogDir:      config.GetLogDir(),
	}
	return conf
}

var defaultConfig = Config{
	HttpPort:    config.GetPorts().HTTP,
	HttpsPort:   config.GetPorts().HTTPS,
	StaticDir:   "./dist",
	TemplateDir: "./dist",
	UploadDir:   config.GetUploadDir(),
	LogDir:      config.GetLogDir(),
}

type webOption func(*Config)

func WithHttpPort(p int) webOption {
	return func(c *Config) {
		c.HttpPort = p
	}
}

func WithHttpsPort(p int) webOption {
	return func(c *Config) {
		c.HttpsPort = p
	}
}

func WithStaticDir(s string) webOption {
	return func(c *Config) {
		c.StaticDir = s
	}
}

func WithTemplateDir(t string) webOption {
	return func(c *Config) {
		c.StaticDir = t
	}
}
