package helper

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	engine "github.com/JoanGTSQ/api"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"neft.web/controller"
	"neft.web/models"
)

var upgrader = websocket.Upgrader{

	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ControlWebsocket(context *gin.Context) {

	//upgrade get request to websocket protocol
	ws, err := upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		engine.Warning.Println(err)
		return
	}
	err = ws.SetReadDeadline(time.Now().Add(45 * time.Minute))
	if err != nil {
		engine.Warning.Println("Could not set dead time out for the new client", err)
		return
	}
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			engine.Warning.Println("Can not close this connection")
			return
		}
	}(ws)

	engine.Debug.Println("New client connected!")
	c, err := GenerateClient(ws, ws.RemoteAddr().String())

	if err != nil {
		err := ws.Close()
		if err != nil {
			engine.Warning.Println("Can not close this connection")
			return
		}
		engine.Warning.Println("error generating new client", err)
	}

	for {

		err = ReadMessage(ws, &c.IncomingMessage)

		if err != nil {
			engine.Warning.Println(err)
			c.LastMessage.Command = "invalid_message"
			c.LastMessage.Data["error"] = "server could not read correctly the message received, please try again"
			c.SendMessage()
			return
		}

		c.LastMessage = c.IncomingMessage

		//Execute the command readed
		c.ExecuteCommand(c.IncomingMessage.Command)

		// Reset the data of the incomming message
		data := make(map[string]interface{})
		cM := controller.Message{
			Data: data,
		}

		c.IncomingMessage = cM
		engine.Debug.Println("command processed")
	}
}

func GenerateClient(ws *websocket.Conn, addr string) (controller.Client, error) {
	u := uuid.NewV4()
	if controller.Hub[u] != nil {
		u = uuid.NewV4()
	}
	mTemplate := make(map[string]interface{})
	message := controller.Message{Data: mTemplate}
	newClient := controller.Client{
		UUID:            u,
		Addr:            addr,
		WS:              ws,
		LastMessage:     message,
		IncomingMessage: message,
		User:            models.User{},
		Sync:            &sync.Mutex{},
	}
	newClient.RegisterToPool()
	go newClient.StartMessageServer()
	go newClient.StartValidator()
	return newClient, nil
}

func ReadMessage(ws *websocket.Conn, dest interface{}) error {
	//	engine.Debug.Println("New Incoming message")
	//engine.Debug.Printf("Client, %d sent another request %d", message.RequestID, message.RequestID)
	//Read Message from client
	_, m, err := ws.ReadMessage()
	if err != nil {
		engine.Warning.Println(err)
		return err
	}

	err = json.Unmarshal(m, &dest)
	if err != nil {
		engine.Warning.Println(err)
		return err
	}
	return nil
}
