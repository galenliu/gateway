package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func RootHandle() gin.HandlerFunc {
	var handle = func(c *gin.Context) {
		var apt = strings.Split(c.Request.Header.Get("Accept"), ",")
		if IsContain(apt, "text/html") || IsContain(apt, "*/*") {
			c.HTML(http.StatusOK, "index.html", nil)
			c.Abort()
			return
		}
		c.Next()
	}
	return handle
}

func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}
