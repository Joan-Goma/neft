package client

import (
	"net/http"

	engine "github.com/JoanGTSQ/api"
	"github.com/gin-gonic/gin"
	"neft.web/auth"
	"neft.web/models"
)

type Controllers struct {
	Users   *Users
	Posts   *Posts
	Devices *Devices
}

var (
	UsersAuth   *Users
	Services    *models.Services
	ResponseMap = make(map[string]interface{})
	response    = engine.Response{}
)

func handleError(err error, c *gin.Context) {
	var responseCode int
	switch err.Error() {
	case "EOF":
		ResponseMap["data"] = gin.H{"error": engine.ERR_INVALID_JSON.Error()}
		responseCode = http.StatusBadRequest
	case "record not found":
		ResponseMap["data"] = gin.H{"error": engine.ERR_NOT_FOUND.Error()}
		responseCode = http.StatusBadRequest
	default:
		ResponseMap["data"] = gin.H{"error": err.Error()}
		responseCode = http.StatusInternalServerError
	}

	ResponseMap["message"] = "Failed"
	response = engine.Response{
		ResponseCode: responseCode,
		Context:      c,
		Response:     ResponseMap,
	}
	response.SendAnswer()
}

// getUserFromToken take the context header "neftAuth" and search the user by the email
func getUserFromContext(context *gin.Context) (*models.User, error) {
	// Obtain data from JWT token
	tokenNeft := context.GetHeader("neftAuth")
	return GetUserFromToken(tokenNeft)
}

func GetUserFromToken(token string) (*models.User, error) {
	claims, err := auth.ReturnClaims(token)
	if err != nil {
		engine.Warning.Println(err)
		return nil, err
	}

	// Search the user from the claims by Email
	user := &models.User{
		Email: claims.Context.User.Email,
	}

	err = user.ByEmail()

	if err != nil {
		return nil, err
	}

	return user, nil
}
