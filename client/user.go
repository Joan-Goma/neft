package client

import (
	"net/http"

	engine "github.com/JoanGTSQ/api"
	"github.com/gin-gonic/gin"
	"neft.web/auth"
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

	var mu *models.MultipleUsers
	err := mu.Count()
	if err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}
	var links Links
	mu.Pagination, links = GeneratePaginationFromRequest(context, int(mu.Quantity))
	if (mu.Pagination == models.Pagination{}) && (links == Links{}) {
		return
	}

	// Retrieve all users data
	err = mu.GetAllUsers()
	if err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}
	if mu.Users == nil {
		// TODO create handle error
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}
	// Close connection returning code 200 and JSON with all users
	ResponseMap["data"] = mu.Users
	ResponseMap["links"] = links
	response = engine.Response{
		ResponseCode: http.StatusOK,
		Context:      context,
		Response:     ResponseMap,
	}
	response.SendAnswer()
}

// Login POST /auth
// Obtain login data (email,password), authenticate it and return jwt token in header
func (us *Users) Login(context *gin.Context) {
	var form LoginStruct

	// Obtain the body in the request and parse to the LoginStruct
	if err := context.ShouldBindJSON(&form); err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}
	user := &models.User{
		Email:    form.Email,
		Password: form.Password,
	}
	// Try to auth with the inserted data and return an error or a user
	err := user.Authenticate()
	if err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}

	// Generate  JWT Token
	tokenString, err := auth.GenerateJWT(user)
	if err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}
	// return token in 200
	ResponseMap["data"] = gin.H{"token": tokenString}
	response = engine.Response{
		ResponseCode: http.StatusOK,
		Context:      context,
		Response:     ResponseMap,
	}
	response.SendAnswer()
}

// RegisterUser PUT /auth
// Obtain user data, register it in the database and return a JWT TOKEN and 201 code
func (us *Users) RegisterUser(context *gin.Context) {
	var user models.User

	// Obtain the body in the request and parse to the user
	if err := context.ShouldBindJSON(&user); err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}

	// Create user with the data received
	if err := user.Create(); err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}

	// Generate  JWT Token
	tokenString, err := auth.GenerateJWT(&user)
	if err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}

	// return token in 200
	ResponseMap["data"] = gin.H{"token": tokenString}
	response = engine.Response{
		ResponseCode: http.StatusOK,
		Context:      context,
		Response:     ResponseMap,
	}
	response.SendAnswer()
}
