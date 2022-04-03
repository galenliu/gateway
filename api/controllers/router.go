package controllers

import (
	"github.com/galenliu/gateway/api/middleware"
	"github.com/galenliu/gateway/api/models"
	things "github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/rules_engine"
	"github.com/galenliu/gateway/plugin"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
	"net/http"
	"time"
)

const InstallAddonTimeOut = 60 * time.Second

type Storage interface {
	models.UsersStore
	things.ThingsStorage
	models.SettingsStore
	models.JsonwebtokenStore
	models.AddonStore
	rules_engine.RuleDB
}

type Models struct {
	ThingsModel  *things.Container
	UsersModel   *models.Users
	SettingModel *models.Settings
}

type Router struct {
	*fiber.App
	addonUrls []string
	logger    logging.Logger
}

func NewRouter(addonUrls []string, manager *plugin.Manager, store Storage, log logging.Logger) *Router {

	//router init
	app := Router{}
	app.addonUrls = addonUrls
	app.logger = log
	app.App = fiber.New()
	app.Use(recover.New())
	app.Use(cors.New(cors.ConfigDefault))
	app.Use(compress.New())

	//models init
	settingModel := models.NewSettingsModel(app.addonUrls, store, log)
	containerModel := things.NewThingsContainerModel(manager, store, log)
	manager.SetThingsContainer(containerModel)
	jwtMiddleware := middleware.NewJWTMiddleware(store, log)
	auth := jwtMiddleware.Auth
	usersModel := models.NewUsersModel(store, log)
	addonModel := models.NewAddonsModel(store, log)
	jsonwebtokenModel := models.NewJsonwebtokenModel(settingModel, store, log)

	actionsModel := models.NewActionsModel(manager, containerModel, log)

	// recover middleware
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	//logger middleware
	app.Use(logger.New())

	//root handler
	app.Use(rootHandler)

	//ping controller
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusNoContent)
	})

	// login
	{
		loginController := NewLoginController(usersModel, jsonwebtokenModel, log)
		app.Post(constant.LoginPath, loginController.handleLogin)
	}

	// Users
	{
		usersController := NewUsersController(usersModel, jsonwebtokenModel, log)
		usersGroup := app.Group(constant.UsersPath)
		usersGroup.Get("/count", usersController.getCount)
		usersGroup.Post("/", usersController.createUser)
	}

	//actions EventBus
	actionsController := NewActionsController(actionsModel, containerModel, manager, log)
	{
		actionsGroup := app.Group(constant.ActionsPath)
		actionsGroup.Post("/", actionsController.handleCreateAction)
		actionsGroup.Get("/", actionsController.handleGetActions)
		actionsGroup.Delete("/:actionName/:actionId", actionsController.handleDeleteAction)
		actionsGroup.Delete("/:thingId/:actionName/:actionId", actionsController.handleDeleteAction)
	}

	//Things
	{
		thingsController := NewThingsControllerFunc(manager, containerModel, log)
		thingsGroup := app.Group(constant.ThingsPath)

		thingsGroup.Put("/:thingId/properties/*", thingsController.handleSetProperty)
		thingsGroup.Get("/:thingId/properties/*", thingsController.handleGetPropertyValue)

		//Handle creating a new thing.
		thingsGroup.Post("/", thingsController.handleCreateThing)

		thingsGroup.Get("/:thingId", thingsController.handleGetThing)
		thingsGroup.Patch("/:thingId", thingsController.handleUpdateThing)
		thingsGroup.Get("/", thingsController.handleGetThings)

		//thingsGroup.Get("/:thingId", websocket.New(handleWebsocket(containerModel, log)))
		thingsGroup.Get("/", websocket.New(handleWebsocket(containerModel, log)))

		//Get the properties of a thing
		thingsGroup.Get("/:thingId"+constant.PropertiesPath, thingsController.handleGetProperties)

		// Modify a Thing
		thingsGroup.Put("/:thingId", thingsController.handleSetThing)
		thingsGroup.Patch("/", auth, thingsController.handlePatchThings)
		thingsGroup.Patch("/:thingId", thingsController.handlePatchThing)
		thingsGroup.Delete("/:thingId", thingsController.handleDeleteThing)

		thingsGroup.Get("/:thingId"+"/"+constant.ActionsPath, actionsController.handleGetActions)
		thingsGroup.Post("/:thingId"+"/"+constant.ActionsPath, actionsController.handleCreateAction)
	}

	//New Things
	{
		newThingsController := NewNewThingsController(manager, containerModel, log)
		newThingsGroup := app.Group(constant.NewThingsPath)
		newThingsGroup.Use("/", func(c *fiber.Ctx) error {
			if websocket.IsWebSocketUpgrade(c) {
				return c.Next()
			}
			return fiber.ErrUpgradeRequired
		})
		newThingsGroup.Get("/", websocket.New(newThingsController.handleNewThingsWebsocket()))
	}

	//Addons
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

	{
		eventsController := NewEventsController()
		eventsGroup := app.Group(constant.EventsPath)
		eventsGroup.Get("/", eventsController.handleGetEvents)
		eventsGroup.Get("/:eventName", eventsController.handlerGetEvent)
	}

	//Settings
	{
		settingsGroup := app.Group(constant.SettingsPath)
		settingsController := NewSettingController(settingModel, log)
		settingsGroup.Get("/addonsInfo", settingsController.handleGetAddonsInfo)
	}

	//rules
	{
		rulesGroup := app.Group(constant.RulesPath)
		rulesController := NewRulesController(store, containerModel)
		rulesGroup.Get("/", rulesController.handleGetRules)
		rulesGroup.Get("/:id", rulesController.handleGetRule)
		rulesGroup.Delete("/:id", rulesController.handlerDeleteRule)
		rulesGroup.Put("/:id", rulesController.handlerUpdateRule)
		rulesGroup.Post("/", rulesController.handleCreateRule)
	}

	//groups
	{
		group := app.Group(constant.GroupsPath)
		groupsController := NewGroupsController()
		group.Get("/", groupsController.handleGetGroups)
		group.Get("/:id", groupsController.handleGetGroup)
		group.Post("/", groupsController.handlerCreateGroup)
		group.Delete("/:id", groupsController.handlerDeleteGroup)
		group.Put("/:id", groupsController.handlerUpdateGroup)
	}

	return &app
}
