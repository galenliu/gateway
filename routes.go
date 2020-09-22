package gateway

import "github.com/gin-gonic/gin"

func CollectRoute(r *gin.Engine) *gin.Engine {

	apiV1 := r.Group("api/v1")
	apiV1.GET("/")
	return r
}
