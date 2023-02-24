package controller

import (
	engine "github.com/JoanGTSQ/api"
	"neft.web/auth"
	"neft.web/models"
)

type AuthFunc struct{}

func (a AuthFunc) Login(client *Client) {
	user := models.User{}
	err := client.GetInterfaceFromMap("user", &user)
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
	user.Password = ""
	client.User = user
	client.LastMessage.Data["user"] = client.User
	client.LastMessage.Data["message"] = "login successful"
	client.LastMessage.Data["token"] = tokenString
	client.SendMessage()
}

func (a AuthFunc) SignUp(client *Client) {
	user := models.User{}
	err := client.GetInterfaceFromMap("user", &user)
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
	user.Password = ""
	client.User = user
	client.LastMessage.Data["user"] = client.User
	client.LastMessage.Data["message"] = "User created successfully"
	client.SendMessage()
}
