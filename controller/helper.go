package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"time"

	engine "github.com/JoanGTSQ/api"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"neft.web/auth"
	"neft.web/models"
)

type Client struct {
	UUID            uuid.UUID       `json:"UUID,omitempty"`
	Addr            string          `json:"-"`
	User            models.User     `json:"user,omitempty"`
	Sync            *sync.Mutex     `json:"-"`
	WS              *websocket.Conn `json:"-"`
	LastMessage     Message         `json:"-"`
	IncomingMessage Message         `json:"-"`
	Token           string          `json:"token,omitempty"`
	UserFuncs       UserFuncs       `json:"-"`
}

type Message struct {
	RequestID int64                  `json:"request_id,omitempty"`
	Command   string                 `json:"command"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

type clientCommandExecution func()

type newClientCommandExecution func(c *Client)

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func (client *Client) ValidateAndExecute(functionToExecute clientCommandExecution) {

	if client.User.Banned {
		client.LastMessage.Data["error"] = "You are banned"
		client.SendMessage()
		return
	} else if client.Token == "" {
		client.LastMessage.Data["error"] = "please log in first"
		client.SendMessage()
		return
	}
	engine.Debug.Println("new command executed:", getFunctionName(functionToExecute))
	functionToExecute()
}

func (client *Client) ValidateAndExecuteNew(functionToExecute newClientCommandExecution) {

	functionToExecute(client)
}

// RegisterToPool Add this client to the general pool
func (client *Client) RegisterToPool() {
	Hub[client.UUID] = client
}

// CheckClientIsSync Check if the user of client is not null
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

// StartMessageServer This loop will update the client messages every time someone sends
func (client *Client) StartMessageServer() {
	for {
		select {
		case m := <-Lobby:
			message := Message{}
			message.Data["message"] = message
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

// MessageController Control the desired message to send all users or single user
func (client *Client) MessageController() {
	if !client.CheckClientIsSync() {
		client.LastMessage.Data["error"] = "Please sync your client before sending messages"
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
		messagee.Receiver = uuidReceiver
		Lobby <- messagee
		client.LastMessage.Data["message"] = "Message succesful"
		client.SendMessage()
	}
}

// GetUserFromRequest Return a user and error from the message request
func (client *Client) GetUserFromRequest() (models.User, error) {
	var user models.User
	// Convert map to json string
	jsonStr, err := json.Marshal(client.IncomingMessage.Data["user"])
	if err != nil {
		engine.Debug.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return models.User{}, err
	}

	// Obtain the body in the request and parse to the user
	if err := json.Unmarshal(jsonStr, &user); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = engine.ERR_INVALID_JSON
		client.SendMessage()
		return models.User{}, err
	}
	return user, nil
}

// GetPostFromRequest Return a post and error from the message request
func (client *Client) GetPostFromRequest() (models.Post, error) {
	var post models.Post
	// Convert map to json string
	jsonStr, err := json.Marshal(client.IncomingMessage.Data["post"])
	if err != nil {
		engine.Debug.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return models.Post{}, err
	}

	// Obtain the body in the request and parse to the user
	if err := json.Unmarshal(jsonStr, &post); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = engine.ERR_INVALID_JSON
		client.SendMessage()
		return models.Post{}, err
	}
	return post, nil
}

// GetInterfaceFromMap Search from the message request and save it into dest
func (client *Client) GetInterfaceFromMap(position string, dest interface{}) error {
	// Convert map to json string
	jsonStr, err := json.Marshal(client.IncomingMessage.Data[position])
	if err != nil {
		engine.Debug.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return err
	}
	// Obtain the body in the request and parse to the user
	if err := json.Unmarshal(jsonStr, dest); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = engine.ERR_INVALID_JSON
		client.SendMessage()
		return err
	}
	return nil
}

func (client *Client) ApplyTemporalBan() {
	client.User.Banned = true
}

func (client *Client) ValidateToken() {
	var token string

	if err := client.GetInterfaceFromMap("token", &token); err != nil {
		client.LastMessage.Command = "invalid_token"
		client.LastMessage.Data["error"] = "please try again"
		return
	}

	if err := auth.ValidateToken(token); err != nil {
		client.Sync.Lock()
		client.LastMessage.Command = "invalid_token"
		client.LastMessage.Data["error"] = "the token was invalid, please verify and connect again"
		client.SendMessage()
		client.ApplyTemporalBan()
		client.Sync.Unlock()
		return
	}
	client.Token = token
	client.TokenToUser()
}
func (client *Client) TokenToUser() {
	claims, err := auth.ReturnClaims(client.Token)
	if err != nil {
		client.LastMessage.Command = "invalid_token"
		client.LastMessage.Data["error"] = "please try again"
		return
	}
	engine.Debug.Println(claims.Context.User)
	client.User = claims.Context.User
}
func (client *Client) CheckToken() {
	////m := rand.Intn(20)
	ticker := time.NewTicker(1 * time.Minute)
	done := make(chan bool)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			//client.LastMessage.Command = "temporal_login"
			//client.SendMessage()
		}
	}
}
