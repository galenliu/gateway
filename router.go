package gateway

import (
	"fmt"
	"gateway/controllers"
	"gateway/plugin"
	"github.com/gin-gonic/gin"

	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

func InitRouter(staticPath, template, upload, logDir string, manager *plugin.AddonsManager) error {

	//init thingsController
	thingsController := controllers.NewThingsController(manager)

	//thingsCtr := controllers.NewThingsController()
	var router = gin.Default()

	gin.SetMode(gin.DebugMode)

	//日志写入文件
	f, _ := os.Create(path.Join(logDir, "gin.log"))
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	//router.Static("/static","/index.html")
	//router.Static("/static", "./static")
	router.LoadHTMLGlob("./static/templates/*")
	//router.LoadHTMLGlob(template)

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

	// 静态文件
	router.Use(func(c *gin.Context) {
		if c.Request.RequestURI == "/" && strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
			c.HTML(http.StatusOK, "index.tmpl", "")
		}
	})

	//router handle api
	router.Use(func(c *gin.Context) {
		if (!strings.Contains(c.Request.Header.Get("Accept"), "text/html") &&
			strings.Contains(c.Request.Header.Get("Accept"), "application/json")) ||
			c.IsWebsocket() {

			//c.Request.URL.Path = ApiPrefix + c.Request.URL.String()
			//c.Request.RequestURI = c.Request.URL.String()
			r, _ := http.NewRequest("GET", ApiPrefix+c.Request.URL.String(), nil)
			c.Request = r
			c.Next()

		}
	})

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
		thingsGroup.GET(":thingId/properties", thingsController.HandleGetProperties)
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
		addonGroup.GET("/", nil)
	}

	_ = router.Run(":8080")
	return nil
}
