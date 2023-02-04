package controller

import (
	engine "github.com/JoanGTSQ/api"
	"neft.web/auth"
)

type AuthFuncs struct{}

func (a AuthFuncs) Login(client *Client) {
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
	user.Password = ""
	client.User = user
	client.LastMessage.Data["user"] = client.User
	client.LastMessage.Data["message"] = "login succesful"
	client.LastMessage.Data["token"] = tokenString
	client.SendMessage()
}

func (a AuthFuncs) SignUp(client *Client) {
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
	user.Password = ""
	client.User = user
	client.LastMessage.Data["user"] = client.User
	client.LastMessage.Data["message"] = "User created successfully"
	client.SendMessage()
}
