package helper

import (
	"bufio"
	"fmt"
	"os"

	"neft.web/controller"

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
	router := gin.New()
	router.Use(middlewares.CORSMiddleware())
	beta := router.Group("/beta")
	{
		beta.GET("/websocket", ControlWebsocket)
	}
	//Controllers for websocket commands
	engine.Debug.Println("Generating controllers for commands of websocket")
	GenerateController("core.whoami", controller.CoreFuncs{}.WhoAmI)
	GenerateController("core.quit", controller.CoreFuncs{}.Quit)
	GenerateController("core.token.validate", controller.CoreFuncs{}.ValidateToken)

	GenerateController("auth.login", controller.AuthFuncs{}.Login)
	GenerateController("auth.signup", controller.AuthFuncs{}.SignUp)

	GenerateController("user.retrieve.single", controller.UserFuncs{}.RetrieveUser)
	GenerateController("user.retrieve.all", controller.UserFuncs{}.RetrieveAllUser)
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
