package helper

import (
	"encoding/json"
	"net/http"
	"sync"

	"reflect"

	engine "github.com/JoanGTSQ/api"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"neft.web/client"
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

	defer ws.Close()

	engine.Debug.Println("New client connected!")
	c, err := GenerateClient(ws, ws.RemoteAddr().String())

	if err != nil {
		ws.Close()
		engine.Warning.Println("error generating new client", err)
	}

	for {
		engine.Debug.Println("New Incoming message")

		engine.Info.Printf("New request ID: %d", c.LastMessage.RequestID)
		//Read Message from client
		_, message, err := ws.ReadMessage()
		if err != nil {
			engine.Warning.Println(err)
			break
		}
		engine.Debug.Println(string(message))
		err = json.Unmarshal(message, &c.IncomingMessage)
		if err != nil {
			engine.Warning.Println(err)
		}

		m, err := ReadMessage(ws)

		if err != nil {
			engine.Warning.Println(err)
			return
		}
		c.LastMessage = m
		//c.IncomingMessage = m

		if reflect.DeepEqual(c.User, models.User{}) || c.User.Banned {
			{
				switch c.IncomingMessage.Command {
				case "validate_token":
					c.ValidateToken()
					break
				case "whoami":
					c.LastMessage.Data["user"] = c
					c.SendMessage()
					break
				case "login":
					c.Login()
					break
				case "sign_up":
					c.SignUp()
					break
				case "count_client":
					c.LastMessage.Data["len"] = len(controller.Hub)
					c.SendMessage()
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
					c.LastMessage.Data["error"] = "command invalid or access denied, please try again"
					c.SendMessage()
				}
			}
		} else {
			switch c.IncomingMessage.Command {
			case "whoami":
				c.LastMessage.Data["user"] = c
				c.SendMessage()
				break
			case "validate_token":
				c.ValidateToken()
				break
			case "count_client":
				c.LastMessage.Data["message"] = len(controller.Hub)
				c.SendMessage()
			case "message":
				c.MessageController()
				break
			case "get_post_from_id":
				if !c.CheckClientIsSync() {
					c.LastMessage.Data["alert"] = "Please sync your client before request posts"
					c.SendMessage()
					break
				}
				postID := int(c.IncomingMessage.Data["postID"].(float64))
				post, err := client.GetPost(postID)
				if err != nil {
					engine.Warning.Println(err)
					c.LastMessage.Data["error"] = err.Error()
					c.SendMessage()
					break
				}
				c.LastMessage.Data["post"] = post
				c.SendMessage()
			case "logout":
				c.User = models.User{}
				break
			case "update_user":
				c.UpdateUser()
				break
			case "delete_user":
				c.DeleteUser()
				break
			case "retrieve_user":
				c.RetrieveUser()
				break
			case "init_user_reset":
				c.InitUserReset()
				break
			case "complete_user_reset":
				c.CompleteReset()
				break
			case "follow_user":
				c.FollowUser()
				break
			case "unfollow_user":
				c.UnfollowUser()
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
		}

		c.IncomingMessage = controller.Message{}
	}
}

func GenerateClient(ws *websocket.Conn, addr string) (controller.Client, error) {
	u := uuid.NewV4()
	if controller.Hub[u] != nil {
		u = uuid.NewV4()
	}
	newClient := controller.Client{
		UUID: u,
		Addr: addr,
		WS:   ws,
		User: models.User{},
		Sync: &sync.Mutex{},
	}
	newClient.RegisterToPool()
	go newClient.StartMessageServer()
	go newClient.CheckToken()
	return newClient, nil
}

func ReadMessage(ws *websocket.Conn) (controller.Message, error) {
	engine.Debug.Println("New Incoming message")
	//engine.Debug.Printf("Client, %d sent another request %d", message.RequestID, message.RequestID)
	var message controller.Message
	//Read Message from client
	_, m, err := ws.ReadMessage()
	if err != nil {
		engine.Warning.Println(err)
		return controller.Message{}, err
	}

	err = json.Unmarshal(m, &message)
	if err != nil {
		engine.Warning.Println(err)
	}
	return message, nil
}
