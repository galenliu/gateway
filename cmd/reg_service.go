package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
	"log"
	"net/http"
	"strings"
)

var box *packr.Box

func main() {

	box = packr.New("config", "../config")
	server()

}

func server() {
	r := gin.Default()
	r.GET("/addons", addonsHandler)
	_ = r.Run(":8443")
}

func addonsHandler(c *gin.Context) {
	d, err := box.Find("addons.json")
	if err != nil {
		log.Print(err)
	}
	js := strings.ReplaceAll(string(d), "\\", "")
	c.String(http.StatusOK, js)
}
