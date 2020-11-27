package app

import (
	"fmt"
	"gateway/app/controllers"
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
}

func NewApp(config Config) *WebApp {
	app := WebApp{}
	app.config = config
	app.Engine = CollectRoute(config)
	return &app
}

func CollectRoute(conf Config) *gin.Engine {

	//init thingsController
	thingsController := controllers.NewThingsController()
	addonController := controllers.NewAddonController()

	//thingsCtr := controllers.NewThingsController()
	var router = gin.Default()

	gin.SetMode(gin.DebugMode)

	//日志写入文件
	f, _ := os.Create(path.Join(conf.LogDir, "web.log"))
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	//html template
	router.LoadHTMLGlob(path.Join(conf.TemplateDir,"index.html"))
	//router.Static("/_assets",conf.StaticDir)
	router.StaticFS("/_assets",http.Dir(path.Join(conf.StaticDir,"_assets")))
	//router.Static("/app",conf.StaticDir)


	//curl -X POST http://localhost:8080/upload \
	//-F "upload[]=@/Users/appleboy/test1.zip" \
	//-F "upload[]=@/Users/appleboy/test2.zip" \
	//-H "Content-Type: multipart/form-data"
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

	// root controller
	router.GET("/", controllers.RootHandle())
	router.GET("/index.html", controllers.RootHandle())

	//Things Controller
	thingsGroup := router.Group(ThingsPath)
	{
		//Handle creating a new thing.
		thingsGroup.POST("/", thingsController.HandleCreateThing)

		//get a list of thingsController
		thingsGroup.GET("/", thingsController.HandleGetThings)

		//get a thing
		thingsGroup.GET("/:thingId", thingsController.HandleGetThing)

		//get the properties of a thing
		thingsGroup.GET("/:thingId/properties", thingsController.HandleGetProperties)

		//get a property of a thing
		thingsGroup.GET("/:thingId/properties/:propertyName", thingsController.HandleGetProperty)

		//set a property of a thing.
		thingsGroup.PUT("/:thingId/properties/:propertyName", thingsController.HandleSetProperty)

		// Modify a Thing.
		thingsGroup.PUT("/:thingId", thingsController.HandleSetThing)

		thingsGroup.PATCH("/", thingsController.HandlePatchThings)
		thingsGroup.PATCH("/:thingId", thingsController.HandlePatchThing)

		thingsGroup.DELETE("/:thingId", thingsController.HandleDeleteThing)
	}

	//Addons Controller
	addonGroup := router.Group(AddonsPath)
	{
		addonGroup.GET("/", addonController.HandlerGetInstalledAddons)
	}

	//Addons Controller
	debugGroup := router.Group(DebugPath)
	{
		debugGroup.GET("/", nil)
	}

	//Settings Controller
	//TODO

	return router
}

func (app *WebApp) Start() {
	httpPort := ":" + strconv.Itoa(app.config.HttpPort)
	_ = app.Run(httpPort)
	_ = app.Run(httpPort)
}
