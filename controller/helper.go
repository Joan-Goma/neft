package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	engine "github.com/JoanGTSQ/api"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	uuid "github.com/satori/go.uuid"
	"neft.web/models"
	"reflect"
	"sync"
)

type Client struct {
	UUID            uuid.UUID       `json:"UUID"`
	Addr            string          `json:"-"`
	User            models.User     `json:"user"`
	Sync            *sync.Mutex     `json:"-"`
	WS              *websocket.Conn `json:"-"`
	LastMessage     Message         `json:"-"`
	IncomingMessage IncomingMessage `json:"-"`
}

type Message struct {
	RequestID int64       `json:"request_id,omitempty"`
	Command   string      `json:"command"`
	Data      interface{} `json:"data"`
}

type IncomingMessage struct {
	RequestID int64                  `json:"request_id,omitempty"`
	Command   string                 `json:"command"`
	Data      map[string]interface{} `json:"data"`
}

//Login Authenticate the user from the request message
func (client *Client) Login() {
	var user models.User
	err := mapstructure.Decode(client.IncomingMessage.Data["user"], &user)
	if err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data = err.Error()
		client.SendMessage()
		return
	}
	if err := user.Authenticate(); err != nil {
		client.LastMessage.Data = err.Error()
		client.SendMessage()
		return
	}
	client.User = user
	client.LastMessage.Data = "login succesfull! Welcome to the hub!"
	client.SendMessage()
}

func (client *Client) AssignUserToClient(user models.User) error {
	client.User = user
	engine.Debug.Println("Client sync succesfull!")
	return nil
}

//RegisterToPool Add this client to the general pool
func (client *Client) RegisterToPool() {
	Hub[client.UUID] = client
}

//CheckClientIsSync Check if the user of client is not null
func (client *Client) CheckClientIsSync() bool {
	if reflect.DeepEqual(client.User, models.User{}) {
		return false
	}
	return true
}

func (client *Client) SendMessage() {
	reqBodyBytes := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes).Encode(client.LastMessage)
	if err != nil {
		engine.Warning.Println(err)
		return
	}
	client.Sync.Lock()
	err = client.WS.WriteMessage(1, reqBodyBytes.Bytes())
	if err != nil {
		engine.Warning.Println(err)
		return
	}
	engine.Debug.Println("New message sent")
	client.Sync.Unlock()
}

var (
	Hub   = make(map[uuid.UUID]*Client)
	Lobby = make(chan models.UserMessage)
)

//StartMessageServer This loop will update the client messages every time someone sends
func (client *Client) StartMessageServer() {
	for {
		select {
		case m := <-Lobby:
			message := Message{
				Data: m,
			}
			if m.Receiver == uuid.FromStringOrNil("0") {
				m.Type = "global"
				for _, client := range Hub {
					message.Command = "global_incoming_message"
					if reflect.DeepEqual(client.User, m.Sender) {
						client.LastMessage = message
						client.SendMessage()
						m.RegisterMessage()
					}
				}
			} else {
				m.Type = "private"
				engine.Debug.Println("New private message")
				message.Command = "private_incoming_message"
				Hub[m.Receiver].LastMessage = message
				Hub[m.Receiver].SendMessage()
				m.RegisterMessage()
			}
		}
	}
}

//MessageController Control the desired message to send all users or single user
func (client *Client) MessageController() {
	if !client.CheckClientIsSync() {
		client.LastMessage.Data = "Please sync your client before sending messages"
		client.SendMessage()
		return
	}
	messagee := models.UserMessage{
		Sender:  client.User,
		Message: fmt.Sprintf("%v", client.IncomingMessage.Data["message"]),
	}

	if client.IncomingMessage.Data["receiver"] == nil {
		if client.User.RoleID == 3 {
			messagee.Receiver = uuid.FromStringOrNil("0")
			Lobby <- messagee
			client.LastMessage.Data = "Message succesful"
			client.SendMessage()
		} else {
			client.LastMessage.Data = "You don't have enough rights to send this messages"
			client.SendMessage()
		}

	} else {
		uuidReceiver, err := uuid.FromString(fmt.Sprintf("%v", client.IncomingMessage.Data["receiver"]))
		if err != nil {
			engine.Warning.Println(err)
			client.LastMessage.Data = err.Error()
			client.SendMessage()
		}
		messagee.Receiver = uuidReceiver
		Lobby <- messagee
		client.LastMessage.Data = "Message succesful"
		client.SendMessage()
	}
}

//GetUserFromMap Return an user and error from the message request
func (client *Client) GetUserFromMap() (models.User, error) {
	var user models.User
	// Convert map to json string
	jsonStr, err := json.Marshal(client.IncomingMessage.Data["user"])
	if err != nil {
		engine.Debug.Println(err)
		client.LastMessage.Data = "Error: " + err.Error()
		client.SendMessage()
		return models.User{}, err
	}

	// Obtain the body in the request and parse to the user
	if err := json.Unmarshal(jsonStr, &user); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data = engine.ERR_INVALID_JSON
		client.SendMessage()
		return models.User{}, err
	}
	return user, nil
}

//GetInterfaceFromMap Search from the message request and save it into dest
func (client *Client) GetInterfaceFromMap(position string, dest interface{}) error {
	// Convert map to json string
	jsonStr, err := json.Marshal(client.IncomingMessage.Data[position])
	if err != nil {
		engine.Debug.Println(err)
		client.LastMessage.Data = "Error: " + err.Error()
		client.SendMessage()
		return err
	}
	// Obtain the body in the request and parse to the user
	if err := json.Unmarshal(jsonStr, dest); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data = engine.ERR_INVALID_JSON
		client.SendMessage()
		return err
	}
	return nil
}
