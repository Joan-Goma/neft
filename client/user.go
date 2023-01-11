package client

import (
	"net/http"

	engine "github.com/JoanGTSQ/api"
	"github.com/gin-gonic/gin"
	"neft.web/models"
)

type Users struct {
	us models.UserService
}

type LoginStruct struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}


func NewUsers(us models.UserService) *Users {
	return &Users{
		us: us,
	}
}

// RetrieveAllUsers GET /users
// Return all users in a JSON
func (us *Users) RetrieveAllUsers(context *gin.Context) {
	count, err := us.us.Count()
	if err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}
	pagination, links := GeneratePaginationFromRequest(context, count)
	if (pagination == models.Pagination{}) && (links == Links{}) {
		return
	}

	// Retrieve all users data
	users, err := us.us.GetAllUsers(pagination)
	if err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}
	if users == nil {
		// TODO create handle error
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}
	// Close connection returning code 200 and JSON with all users
	ResponseMap["data"] = users
	ResponseMap["links"] = links
	response = engine.Response{
		ResponseCode: http.StatusOK,
		Context:      context,
		Response:     ResponseMap,
	}
	response.SendAnswer()
}

