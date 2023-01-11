package middlewares

import (
	engine "github.com/JoanGTSQ/api"
	"github.com/gin-gonic/gin"
	"neft.web/auth"
	"neft.web/client"
)

var Maitenance bool

func RequireAuth() gin.HandlerFunc {
	return func(context *gin.Context) {

		context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		context.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		context.Writer.Header().Set("Access-Control-Allow-Headers", "neftAuth, XMLHttpRequest, Content-Type, Content-Length")
		context.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		tokenString := context.GetHeader("neftAuth")

		if Maitenance == true {
			context.JSON(503, gin.H{"maitenance": true})
			context.Abort()
			return
		}

		if tokenString == "" {
			context.JSON(401, gin.H{"error": engine.ERR_JWT_TOKEN_REQUIRED.Error()})
			context.Abort()
			return
		}

		if err := auth.ValidateToken(tokenString); err != nil {
			context.JSON(401, gin.H{"error": err.Error()})
			context.Abort()
			return
		}

		// claim, err := auth.ReturnClaims(tokenString)
		// if err != nil {
		// 	context.JSON(401, gin.H{"error": err.Error()})
		// 	context.Abort()
		// 	return
		// }
		//  if claim.Context.User.Banned {
		//    context.JSON(401, gin.H{"error": "you are banned from this api"})
		// context.Abort()
		// return
		//  }

		context.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		client.ResponseMap = make(map[string]interface{})

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
