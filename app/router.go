package app

import (
	"context"
	"fmt"
	"gateway/app/controllers"
	"gateway/pkg/runtime"
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
}

type WebApp struct {
	*gin.Engine
	config Config
	ctx    context.Context
}

func NewWebAPP(conf Config) *WebApp {
	app := WebApp{}
	app.config = conf
	app.Engine = CollectRoute(app.config)
	return &app
}

func CollectRoute(conf Config) *gin.Engine {

	//init thingsController
	thingsController := controllers.NewThingsController()
	addonController := controllers.NewAddonController()
	settingsController := controllers.NewSettingController()

	//thingsCtr := controllers.NewThingsController()
	var router = gin.Default()

	gin.SetMode(gin.DebugMode)

	//日志写入文件
	f, _ := os.Create(path.Join(conf.LogDir, "web.log"))
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
			_ = c.SaveUploadedFile(file, conf.UploadDir)
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

	//Addons Controller
	addonGroup := router.Group(AddonsPath)
	{
		addonGroup.GET("/", addonController.HandlerGetAddons)
		addonGroup.PUT("/", addonController.HandlerInstallAddon)
		addonGroup.PUT("/:addon_id", addonController.HandlerSetAddon)
		addonGroup.GET("/:addon_id/config", addonController.HandlerGetAddonConfig)
		addonGroup.PUT("/:addon_id/config", addonController.HandlerSetAddonConfig)
	}

	//settings Controller
	debugGroup := router.Group(SettingsPath)
	{
		debugGroup.GET("/addons_info", settingsController.HandleGetAddonsInfo)
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

func NewDefaultWebConfig() Config {
	conf := Config{
		HttpPort:    runtime.RuntimeConf.Ports["http"],
		HttpsPort:   runtime.RuntimeConf.Ports["https"],
		StaticDir:   "./dist",
		TemplateDir: "./dist",
		UploadDir:   runtime.RuntimeConf.UploadDir,
		LogDir:      runtime.RuntimeConf.LogDir,
	}
	return conf
}
