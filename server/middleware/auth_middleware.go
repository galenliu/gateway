package middleware

//import (
//	"github.com/gin-gonic/gin"
//	"net/http"
//	"strings"
//)
//
//func AuthMiddleware() gin.HandlerFunc {
//	return func(ctx *gin.Context) {
//		//获取 Authorization header
//		tokenString := ctx.GetHeader("Authorization")
//		if tokenString == "" || strings.HasPrefix(tokenString, "Bearer") {
//			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
//			ctx.Abort()
//			return
//		}
//		tokenString = tokenString[7:]
//	}
//}
