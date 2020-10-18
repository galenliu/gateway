package gateway

import (
	"fmt"
	"gateway/controllers"
	"gateway/plugin"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

func CollectRoute(staticPath, templates, upload, logDir string, manager *plugin.AddonsManager, _log *zap.Logger) error {

	//init thingsController
	thingsController := controllers.NewThingsController(manager, _log)
	addonController := controllers.NewAddonController(manager, _log)

	//thingsCtr := controllers.NewThingsController()
	var router = gin.Default()

	gin.SetMode(gin.DebugMode)

	//日志写入文件
	f, _ := os.Create(path.Join(logDir, "gin.log"))
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	//html template
	router.LoadHTMLGlob(templates)

	//images,js
	router.StaticFS("/images", http.Dir(staticPath+"/images"))
	router.StaticFS("/js", http.Dir(staticPath+"/js"))
	router.StaticFS("/css", http.Dir(staticPath+"/css"))
	router.StaticFS("/fonts", http.Dir(staticPath+"/fonts"))

	router.Static("/app.webmanifest", staticPath+"/app.webmanifest")

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
			log.Println(file.Filename)

			// Upload the file to specific dst.
			_ = c.SaveUploadedFile(file, upload)
		}
		c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})

	// root controller
	router.GET("/", controllers.RootHandle())

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

	_ = router.Run(":8080")
	return nil
}
