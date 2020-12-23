package app

import (
	"context"
	"fmt"
	"gateway/app/controllers"
	"gateway/app/models"
	"gateway/config"
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

	//日志写入文件
	f, _ := os.Create(path.Join(app.config.LogDir, "web.log"))
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	//html template
	//router.LoadHTMLGlob(path.Join(conf.TemplateDir,"index.html"))
	//router.Static("/_assets",conf.StaticDir)
	//router.StaticFS("/_assets",http.Dir(path.Join(conf.StaticDir,"_assets")))
	//router.Static("/app",conf.StaticDir)

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
	router.GET("/", controllers.RootHandle())
	router.GET("/index.html", controllers.RootHandle())

	//Things Controller
	thingsGroup := router.Group(ThingsPath)
	{
		thingsController := controllers.NewThingsControllerFunc()
		//Handle creating a new thing.
		thingsGroup.POST("/", thingsController.HandleCreateThing)
		thingsGroup.GET("/", thingsController.HandleGetThings)

		//get a thing
		thingsGroup.GET("/:thing_id", thingsController.HandleGetThing)

		//get the properties of a thing
		thingsGroup.GET("/:thing_id/properties", thingsController.HandleGetProperties)

		//get a property of a thing
		thingsGroup.GET("/:thing_id/properties/:propertyName", thingsController.HandleGetProperty)

		//set a property of a thing.
		thingsGroup.PUT("/:thing_id/properties/:propertyName", thingsController.HandleSetProperty)

		// Modify a Thing.
		thingsGroup.PUT("/:thing_id", thingsController.HandleSetThing)

		thingsGroup.PATCH("/", thingsController.HandlePatchThings)
		thingsGroup.PATCH("/:thing_id", thingsController.HandlePatchThing)

		thingsGroup.DELETE("/:thing_id", thingsController.HandleDeleteThing)
	}

	newThingsController := controllers.NewNewThingsController(app.things)
	newThingsGroup := router.Group(NewThingsPath)
	{
		newThingsGroup.GET("", newThingsController.HandleWebsocket)
	}


	//Addons Controller
	addonGroup := router.Group(AddonsPath)
	{
		addonController := controllers.NewAddonController()
		addonGroup.GET("/", addonController.HandlerGetAddons)
		addonGroup.PUT("/", addonController.HandlerInstallAddon)
		addonGroup.PUT("/:addon_id", addonController.HandlerSetAddon)
		addonGroup.GET("/:addon_id/config", addonController.HandlerGetAddonConfig)
		addonGroup.PUT("/:addon_id/config", addonController.HandlerSetAddonConfig)
	}

	//settings Controller
	debugGroup := router.Group(SettingsPath)
	{
		settingsController := controllers.NewSettingController()
		debugGroup.GET("/addons_info", settingsController.HandleGetAddonsInfo)
	}

	//actions Controller
	actionsGroup := router.Group(ActionsPath)
	{

		actionsController := controllers.NewActionsController()
		actionsGroup.POST("/", actionsController.HandleActions)
		actionsGroup.DELETE("/:actionName/:actionId", actionsController.HandleActions)
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
