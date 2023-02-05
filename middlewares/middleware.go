package middlewares

import (
	engine "github.com/JoanGTSQ/api"
	"github.com/gin-gonic/gin"
)

var Maitenance bool

func CORSMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {

		engine.Gin.Println("NEW REQUEST | ", context.ClientIP(), " | ", context.Request.Method, " | ", context.Request.URL)

		context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		context.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		context.Writer.Header().Set("Access-Control-Allow-Headers", "neftAuth, XMLHttpRequest, Content-Type, Content-Length")
		context.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if context.Request.Method == "OPTIONS" {
			context.AbortWithStatus(204)
			return
		}
		if Maitenance == true {
			context.AbortWithStatus(503)
			return
		}
		context.Next()
	}
}
