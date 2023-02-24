package controller

import (
	"fmt"
	engine "github.com/JoanGTSQ/api"
	uuid "github.com/satori/go.uuid"
	"neft.web/models"
)

type CoreFuncs struct {
}

func (core CoreFuncs) WhoAmI(client *Client) {
	client.LastMessage.Data["client"] = client
	client.SendMessage()
}

func (core CoreFuncs) PrintHub(client *Client) {

	client.LastMessage.Data["hub"] = Hub
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

func (core CoreFuncs) MessageController(client *Client) {
	if !client.CheckClientIsSync() {
		client.LastMessage.Data["error"] = "Please sync your client before sending messages"
		client.SendMessage()
		return
	}
	engine.Debug.Println(client.User.ID)
	message := models.UserMessage{
		Sender:  client.User,
		Message: fmt.Sprintf("%v", client.IncomingMessage.Data["message"]),
	}
	engine.Debug.Println(message.Sender.ID)
	if client.IncomingMessage.Data["receiver"] == nil {
		if client.User.RoleID == 3 {
			message.Receiver = uuid.FromStringOrNil("0")
			Lobby <- message
			client.LastMessage.Data["message"] = "Message succesful"
			client.SendMessage()
		} else {
			client.LastMessage.Data["error"] = "You don't have enough rights to send this messages"
			client.SendMessage()
		}

	} else {
		uuidReceiver, err := uuid.FromString(fmt.Sprintf("%v", client.IncomingMessage.Data["receiver"]))
		if err != nil {
			engine.Warning.Println(err)
			client.LastMessage.Data["error"] = err.Error()
			client.SendMessage()
		}
		message.Receiver = uuidReceiver
		Lobby <- message
		client.LastMessage.Data["message"] = "Message succesful"
		client.SendMessage()
	}
}
