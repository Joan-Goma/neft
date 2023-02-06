package controller

import engine "github.com/JoanGTSQ/api"

type CoreFuncs struct {
}

func (core CoreFuncs) WhoAmI(client *Client) {
	client.LastMessage.Data["client"] = client
	client.SendMessage()
}

func (core CoreFuncs) Quit(client *Client) {
	client.LastMessage.Data["message"] = "closing connection, goodbye"
	err := client.WS.Close()
	if err != nil {
		engine.Warning.Println("error trying to close connection", err.Error())
		return
	}
}

func (core CoreFuncs) ValidateToken(client *Client) {
	client.MessageReader <- client.IncomingMessage
}
