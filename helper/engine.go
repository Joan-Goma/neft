package helper

import (
	"bufio"
	"fmt"
	"os"

	engine "github.com/JoanGTSQ/api"
	"github.com/gin-gonic/gin"
	"neft.web/client"
	"neft.web/middlewares"
	"neft.web/models"
)

// InitDB start a connection with the database, return error if can't connect
func InitDB(debugdb bool) error {
	// Create connection with DB
	engine.Debug.Println("Creating connection with DB")
	var err error
	client.Services, err = models.NewServices(fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		database.URL,
		5432,
		database.User,
		database.Password,
		database.Name,
		database.SslMode),
		debugdb)
	if err != nil {
		return err
	}
	return nil
}

// InitRouter Generate a router with directions and middlewares
func InitRouter() *gin.Engine {
	engine.Debug.Println("Creating gin router")
	controllersR := client.Controllers{
		Users:   client.NewUsers(client.Services.User),
		Posts:   client.NewPosts(client.Services.Post),
		Devices: client.NewDevices(client.Services.Device),
	}
	client.UsersAuth = controllersR.Users
	router := gin.New()
	router.Use(middlewares.CORSMiddleware())
	api := router.Group("/v1")
	{

		secured := api.Group("/secured").Use(middlewares.RequireAuth())
		{
			// USER

			secured.GET("/users", controllersR.Users.RetrieveAllUsers)

			secured.PUT("/post", controllersR.Posts.CreatePost)
			secured.DELETE("/post", controllersR.Posts.DeletePost)
			secured.PATCH("/post", controllersR.Posts.UpdatePost)
			secured.GET("/posts", controllersR.Posts.RetrieveAllPost)
			// secured.GET("/post/:id", controllersR.Posts.GetPost)
			secured.GET("/post/like/:id", controllersR.Posts.Like)
			secured.DELETE("/post/like/:id", controllersR.Posts.Unlike)
			secured.PUT("/post/comment/:id", controllersR.Posts.Comment)
			secured.DELETE("/post/comment/:id", controllersR.Posts.Uncomment)
		}
	}

	beta := router.Group("/beta")
	{
		beta.GET("/websocket", ControlWebsocket)
	}
	return router
}

// ReadInput read in every moment the console and change maintenance and debug mode
func ReadInput(debug bool) {
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		switch input.Text() {
		case "maitenance":
			middlewares.Maitenance = !middlewares.Maitenance
			engine.Info.Println("maintenance", middlewares.Maitenance)
		case "debug":
			debug = !debug
			engine.EnableDebug(debug)
			engine.Info.Println("debug mode", debug)
		}
	}
}
