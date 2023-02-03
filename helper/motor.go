package helper

import (
	"encoding/json"
	"net/http"
	"reflect"
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

		engine.Debug.Println("new message to proccess")

		err = ReadMessage(ws, &c.IncomingMessage)

		if err != nil {
			engine.Warning.Println(err)
			c.LastMessage.Command = "invalid_message"
			c.LastMessage.Data["error"] = "server could not read correctly the message received, please try again"
			c.SendMessage()
			return
		}
		c.LastMessage = c.IncomingMessage

		switch c.IncomingMessage.Command {
		case "login":
			c.Login()
			break
		case "sign_up":
			c.SignUp()
			break
		case "whoami":
			if !reflect.DeepEqual(c.User, models.User{}) {
				c.LastMessage.Data["user"] = &c.User
			} else {
				c.LastMessage.Data["error"] = "you are not logged"
			}
			c.SendMessage()
			break
		case "validate_token":
			c.ValidateToken()
			break
		case "count_client":
			c.LastMessage.Data["clients_in_pool"] = len(controller.Hub)
			c.SendMessage()
		case "message":
			c.ValidateAndExecute(c.MessageController)
			break
		case "get_post_from_id":
			c.ValidateAndExecute(c.GetPost)
		case "update_user":
			c.ValidateAndExecute(c.UpdateUser)
			break
		case "delete_user":
			c.ValidateAndExecute(c.DeleteUser)
			break
		case "retrieve_user":
			c.ValidateAndExecute(c.RetrieveUser)
			break
		case "init_user_reset":
			c.ValidateAndExecute(c.InitUserReset)
			break
		case "complete_user_reset":
			c.ValidateAndExecute(c.CompleteReset)
			break
		case "follow_user":
			c.ValidateAndExecute(c.FollowUser)
			break
		case "unfollow_user":
			c.ValidateAndExecute(c.UnfollowUser)
			break
		case "quit":
			delete(controller.Hub, c.UUID)
			err := ws.Close()
			if err != nil {
				engine.Debug.Println(err)
				break
			}
			return
		default:
			c.LastMessage.Data["error"] = "command invalid, please try again"
			c.SendMessage()
		}
		//c.IncomingMessage = controller.Message{}
		data := make(map[string]interface{})
		cM := controller.Message{
			Data: data,
		}
		c.IncomingMessage = cM
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
	go newClient.CheckToken()
	return newClient, nil
}

func ReadMessage(ws *websocket.Conn, dest interface{}) error {
	engine.Debug.Println("New Incoming message")
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
