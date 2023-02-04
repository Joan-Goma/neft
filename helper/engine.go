package helper

import (
	"bufio"
	"fmt"
	"neft.web/controller"
	"os"

	engine "github.com/JoanGTSQ/api"
	"github.com/gin-gonic/gin"
	"neft.web/middlewares"
	"neft.web/models"
)

// InitDB start a connection with the database, return error if can't connect
func InitDB(debugdb bool) error {
	// Create connection with DB
	engine.Debug.Println("Creating connection with DB")
	var err error
	err = models.NewServices(fmt.Sprintf(
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
	//controllersR := client.Controllers{
	//	Users:   client.NewUsers(client.Services.User),
	//	Posts:   client.NewPosts(client.Services.Post),
	//	Devices: client.NewDevices(client.Services.Device),
	//}
	//client.UsersAuth = controllersR.Users
	router := gin.New()
	router.Use(middlewares.CORSMiddleware())
	/*api := router.Group("/v1")
	{
		api.POST("/auth", controllersR.Users.Login)
		api.PUT("/auth", controllersR.Users.RegisterUser)
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
	*/

	beta := router.Group("/beta")
	{
		beta.GET("/websocket", ControlWebsocket)
	}

	GenerateController("core.whoami", controller.CoreFuncs{}.WhoAmI)
	GenerateController("core.quit", controller.CoreFuncs{}.Quit)

	GenerateController("auth.login", controller.AuthFuncs{}.Login)
	GenerateController("auth.signup", controller.AuthFuncs{}.SignUp)

	GenerateController("user.retrieve.single", controller.UserFuncs{}.RetrieveUser)
	GenerateController("user.retrieve.all", controller.UserFuncs{}.RetrieveUser)
	GenerateController("user.update", controller.UserFuncs{}.UpdateUser)
	GenerateController("user.password.init_reset", controller.UserFuncs{}.InitUserReset)
	GenerateController("user.password.complete_reset", controller.UserFuncs{}.CompleteReset)
	GenerateController("user.delete", controller.UserFuncs{}.DeleteUser)
	GenerateController("user.interaction.follow", controller.UserFuncs{}.FollowUser)
	GenerateController("user.interaction.unfollow", controller.UserFuncs{}.UnfollowUser)

	GenerateController("spot.create", controller.SpotFuncs{}.CreateSpot)
	GenerateController("spot.delete", controller.SpotFuncs{}.DeleteSpot)
	GenerateController("spot.update", controller.SpotFuncs{}.UpdateSpot)
	GenerateController("spot.retrieve.single", controller.SpotFuncs{}.RetrieveSingle)
	GenerateController("spot.retrieve.all", controller.SpotFuncs{}.RetrieveAllSpot)
	GenerateController("spot.interaction.like", controller.SpotFuncs{}.LikeSpot)
	GenerateController("spot.interaction.unlike", controller.SpotFuncs{}.UnlikeSpot)
	GenerateController("spot.interaction.comment", controller.SpotFuncs{}.CommentSpot)
	GenerateController("spot.interaction.uncomment", controller.SpotFuncs{}.UncommentSpot)

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
