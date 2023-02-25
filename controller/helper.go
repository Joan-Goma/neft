package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"math/rand"
	"reflect"
	"strings"
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
	MessageReader   chan Message    `json:"-"`
	Token           string          `json:"token,omitempty"`
}

type Message struct {
	RequestID int64                  `json:"request_id,omitempty"`
	Command   string                 `json:"command,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

var (
	MapFuncs = make(map[string]ClientCommandExecution)
)

func (client *Client) ExecuteCommand(commandName string) {

	if MapFuncs[commandName] == nil || client.IncomingMessage.RequestID == 0 {
		client.LastMessage.Data["error"] = "command invalid, please try again"
		client.SendMessage()
		return
	}
	if !strings.Contains(commandName, "auth") && !strings.Contains(commandName, "core") {
		validateAndExecute(MapFuncs[commandName], client)
		return
	}
	MapFuncs[commandName](client)
}

type ClientCommandExecution func(c *Client)

func validateAndExecute(functionToExecute ClientCommandExecution, client *Client) {

	if client.User.Banned {
		client.LastMessage.Data["error"] = "You are banned"
		client.SendMessage()
		return
	} else if reflect.DeepEqual(client.User, models.User{}) {
		client.LastMessage.Data["error"] = "please log in first"
		client.SendMessage()
		return
	}

	functionToExecute(client)
}

func GenerateClient(ws *websocket.Conn, addr string) *Client {
	u := uuid.NewV4()

	mTemplate := make(map[string]interface{})
	mmTemplate := make(map[string]interface{})
	mssg := Message{Data: mmTemplate}
	message := Message{Data: mTemplate}
	client := &Client{
		UUID:            u,
		Addr:            addr,
		WS:              ws,
		LastMessage:     mssg,
		IncomingMessage: message,
		User:            models.User{},
		Sync:            &sync.Mutex{},
	}
	client.RegisterToPool()
	go client.StartMessageServer()
	go client.StartValidator()
	return client
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
			if m.Receiver == uuid.FromStringOrNil("0") {
				m.Type = "global"
				for _, client := range Hub {
					client.LastMessage.Command = "global_incoming_message"
					engine.Debug.Printf("client from hub %d, sender %d", client.User.ID, m.Sender.ID)
					if !reflect.DeepEqual(client.User, m.Sender) {
						client.LastMessage.Data["message"] = m
						client.SendMessage()
						m.RegisterMessage()
					}
				}
			} else {
				m.Type = "private"
				Hub[m.Receiver].LastMessage.Command = "private_incoming_message"
				Hub[m.Receiver].LastMessage.Data["message"] = m
				Hub[m.Receiver].SendMessage()
				engine.Debug.Println("New private message")
				m.RegisterMessage()
			}
		}
	}
}

// GetInterfaceFromMap Search from the message request and save it into dest
func (client *Client) GetInterfaceFromMap(position string, dest interface{}) error {

	if client.IncomingMessage.Data[position] == nil {
		return errors.New("could not find the object please try again")
	}
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

func (client *Client) ValidateToken() error {
	var token string

	if err := client.GetInterfaceFromMap("token", &token); err != nil {
		err = errors.New("could not load the token, please try again")
		return err
	}

	if err := auth.ValidateToken(token); err != nil {
		err = errors.New("token not valid, please try again")
		return err
	}
	client.Token = token
	client.TokenToUser()
	return nil
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

func (client *Client) StartValidator() {
	r := time.Now().UnixNano()
	rand.Seed(r)
	m := rand.Intn(320)
	request := 01000 + m
	mTemplate := make(map[string]interface{})
	mssg := Message{
		RequestID: int64(request),
		Command:   "temporal_validator",
		Data:      mTemplate,
	}

	ticker := time.NewTicker(120 * time.Minute)
	done := make(chan bool)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			engine.Debug.Println("starting to ban")
			client.LastMessage = mssg
			client.SendMessage()
			client.CompleteValidator(mssg.RequestID)
		}
	}
}

func (client *Client) CompleteValidator(requestID int64) {
	engine.Debug.Println("Starting validation....")
	tries := 0
	ticker := time.NewTicker(5 * time.Minute)
	client.LastMessage.Command = "core.token.validator"
	for {
		select {
		case message := <-client.MessageReader:
			if message.RequestID == requestID {
				err := client.ValidateToken()
				if err != nil {
					tries++
					if tries >= 3 {
						client.ApplyTemporalBan()
						client.LastMessage.Data["error"] = "you will be banned"
						client.SendMessage()
						err := client.WS.Close()
						if err != nil {
							return
						}
						return
					}
				}
				tries = 0
				return
			}
		case <-ticker.C:
			tries++
			client.LastMessage.Data["error"] = "you didn't send any message"
			client.SendMessage()
			if tries >= 3 {
				client.ApplyTemporalBan()
				client.LastMessage.Data["error"] = "you will be banned"
				client.SendMessage()
				err := client.WS.Close()
				if err != nil {
					return
				}
				return
			}
		}
	}
}
