package controllers

import (
	"context"
	"fmt"
	"gateway/config"
	"gateway/server/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
)

type Config struct {
	HttpPort    int
	HttpsPort   int
	StaticDir   string
	TemplateDir string
	UploadDir   string
	LogDir      string
	Ctx         context.Context
}

type WebApp struct {
	*gin.Engine
	config Config
	ctx    context.Context
	things *models.Things
}

func NewWebAPP(conf Config) *WebApp {
	app := WebApp{}
	app.things = models.NewThings()
	app.config = conf
	app.Engine = CollectRoute(&app)
	return &app
}

func CollectRoute(app *WebApp) *gin.Engine {

	//init thingsController

	//thingsCtr := controllers.NewThingsControllerFunc()
	var router = gin.Default()

	gin.SetMode(gin.DebugMode)

	//解决跨域问题 仅测试
	if gin.Mode() == gin.DebugMode {
		router.Use(cors.Default())
	}

	//日志写入文件
	f, _ := os.Create(path.Join(app.config.LogDir, "web.log"))
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	//html template
	//router.LoadHTMLGlob(path.Join(conf.TemplateDir,"index.html"))
	//router.Static("/_assets",conf.StaticDir)
	//router.StaticFS("/_assets",http.Dir(path.Join(conf.StaticDir,"_assets")))
	//router.Static("/server",conf.StaticDir)

	router.POST("/upload", func(c *gin.Context) {
		// 单文件
		// Multipart form
		form, _ := c.MultipartForm()
		files := form.File["upload[]"]

		for _, file := range files {
			// Upload the file to specific dst.
			_ = c.SaveUploadedFile(file, app.config.UploadDir)
		}
		c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})

	//ping controller
	router.GET("/ping", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	// root controller
	//router.GET("/", controllers.RootHandle())
	router.GET("/index.html", RootHandle())

	//Things Controller
	thingsGroup := router.Group(models.ThingsPath)

	{
		thingsController := NewThingsControllerFunc()

		//set a properties of a thing.
		thingsGroup.PUT("/:thingId/properties/:propertyName", thingsController.handleSetProperty)

		//Handle creating a new thing.
		thingsGroup.POST("/", thingsController.handleCreateThing)

		thingsGroup.GET("/:thingId", thingsController.handleGetThing)

		//get the properties of a thing
		thingsGroup.GET("/:thingId/properties", thingsController.handleGetProperties)

		//get a properties of a thing
		thingsGroup.GET("/:thingId/properties/:propertyName", thingsController.handleGetProperty)

		// Modify a ThingInfo.
		thingsGroup.PUT("/:thingId", thingsController.handleSetThing)

		thingsGroup.PATCH("/", thingsController.handlePatchThings)
		thingsGroup.PATCH("/:thingId", thingsController.handlePatchThing)

		thingsGroup.DELETE("/:thingId", thingsController.handleDeleteThing)
		thingsGroup.GET("/", thingsController.handleGetThings)
	}

	newThingsGroup := router.Group(models.NewThingsPath)
	{
		newThingsController := NewNewThingsController(app.things)
		newThingsGroup.GET("", newThingsController.HandleWebsocket)
	}

	//Addons Controller
	addonGroup := router.Group(models.AddonsPath)
	{
		addonController := NewAddonController()
		addonGroup.GET("/", addonController.handlerGetAddons)
		addonGroup.POST("/", addonController.handlerInstallAddon)
		addonGroup.PUT("/:addon_id", addonController.handlerSetAddon)
		addonGroup.GET("/:addon_id/config", addonController.HandlerGetAddonConfig)
		addonGroup.PUT("/:addon_id/config", addonController.HandlerSetAddonConfig)
	}

	//settings Controller
	debugGroup := router.Group(models.SettingsPath)
	{
		settingsController := NewSettingController()
		debugGroup.GET("/addonsInfo", settingsController.handleGetAddonsInfo)
	}

	//actions Controller
	actionsGroup := router.Group(models.ActionsPath)
	{
		actionsController := NewActionsController()
		actionsGroup.POST("/", actionsController.HandleActions)
		actionsGroup.DELETE("/:actionName/:actionId", actionsController.HandleDeleteAction)
	}

	return router
}

func (app *WebApp) Start() error {
	httpPort := ":" + strconv.Itoa(app.config.HttpPort)
	err := app.Run(httpPort)
	if err != nil {
		return err
	}
	return nil
}

func NewDefaultWebConfig(ctx context.Context) Config {
	conf := Config{
		HttpPort:    config.Conf.Ports["http"],
		HttpsPort:   config.Conf.Ports["https"],
		StaticDir:   "./dist",
		TemplateDir: "./dist",
		UploadDir:   config.Conf.UploadDir,
		LogDir:      config.Conf.LogDir,
		Ctx:         ctx,
	}
	return conf
}
