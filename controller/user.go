package controller

import (
	"fmt"

	"neft.web/auth"

	engine "github.com/JoanGTSQ/api"
	"neft.web/models"
)

type UserFuncs struct {
}

func (u UserFuncs) Login(client *Client) {

	user, err := client.GetUserFromRequest()
	if err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}
	err = user.Authenticate()
	if err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	// Generate  JWT Token
	tokenString, err := auth.GenerateJWT(user)
	if err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	client.User = user
	client.LastMessage.Data["message"] = "login succesful"
	client.LastMessage.Data["token"] = tokenString
	client.SendMessage()
}

// UpdateUser Get the user from the json request, compare ID and update it
func (u UserFuncs) UpdateUser(client *Client) {

	newUser, err := client.GetUserFromRequest()
	if err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	if client.User.ID != newUser.ID {
		engine.Warning.Println(engine.ERR_NOT_SAME_USER, newUser.ID, client.User.ID)
		client.LastMessage.Data["error"] = engine.ERR_NOT_SAME_USER.Error()
		client.SendMessage()
		return
	}
	// Try to update the user
	if err := newUser.Update(); err != nil {
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		engine.Warning.Println(err)
		return
	}
	client.User = newUser
	client.LastMessage.Data["message"] = "User updated successfully"
	client.SendMessage()
}

// DeleteUser Obtain user data, search by ID and delete it, return code http.StatusOK
func (u UserFuncs) DeleteUser(client *Client) {
	var user models.User

	user, err := client.GetUserFromRequest()
	if err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	if user.ID != client.User.ID {
		engine.Warning.Println("Someone is trying to delete a user without rights")
		client.LastMessage.Data["error"] = "you are trying to delete a user without rights"
		client.SendMessage()
		return
	}
	// Try to delete the user
	if err := user.Delete(); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}
	client.User = models.User{}

	// Close connection with status http.StatusOK (resource deleted)
	client.LastMessage.Data["message"] = "user deleted"
	client.SendMessage()
}

// RetrieveUser Obtain the user from the json request and search it by ID
func (u UserFuncs) RetrieveUser(client *Client) {

	user, err := client.GetUserFromRequest()
	if err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	err = user.ByID()
	if err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	client.LastMessage.Data["user"] = user
	client.SendMessage()
}

// SignUp Register a new user
func (u UserFuncs) SignUp(client *Client) {
	user, err := client.GetUserFromRequest()
	if err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}
	// Create user with the data received
	if err := user.Create(); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	client.User = user

	client.LastMessage.Data["message"] = "User created successfully"
	client.SendMessage()
}

// InitUserReset Initiate the process of restore a password
func (u UserFuncs) InitUserReset(client *Client) {

	user, err := client.GetUserFromRequest()
	if err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	// Initiate reset with the data received
	token, err := user.InitiateReset()
	if err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	client.LastMessage.Data["token"] = token
	client.SendMessage()
}

// CompletePasswdReset This struct is used to unmarshall the token and password to complete the reset
type CompletePasswdReset struct {
	Token    string
	Password string
}

// CompleteReset Use this controller to complete the password reset with token
func (u UserFuncs) CompleteReset(client *Client) {
	user, err := client.GetUserFromRequest()
	if err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}
	form := &CompletePasswdReset{}
	err = client.GetInterfaceFromMap("token", form)
	if err != nil {
		engine.Debug.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	err = user.CompleteReset(form.Token, user.Password)
	if err != nil {
		engine.Debug.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	client.LastMessage.Data["message"] = "User password reset successfully"
	client.SendMessage()

}

// FollowUser Use this controller to follow between users
func (u UserFuncs) FollowUser(client *Client) {

	userToFollow := &models.User{}

	err := client.GetInterfaceFromMap("user_to_follow", userToFollow)
	if err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}
	if err = client.User.Follow(userToFollow.ID); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}
	client.LastMessage.Data["message"] = fmt.Sprintf("User ID %d followed successfully", userToFollow.ID)
	client.SendMessage()
}

// UnfollowUser Use this controller to unfollow between users
func (u UserFuncs) UnfollowUser(client *Client) {

	userToUnfollow := &models.User{}
	err := client.GetInterfaceFromMap("user_to_follow", userToUnfollow)
	if err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}
	if err = client.User.Unfollow(userToUnfollow.ID); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}
	client.LastMessage.Data["message"] = fmt.Sprintf("User ID %d unfollowed successfully", userToUnfollow.ID)
	client.SendMessage()
}
