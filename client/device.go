 package client

import (
	engine "github.com/JoanGTSQ/api"
	"github.com/gin-gonic/gin"
	"neft.web/models"

	"net"
	"net/http"
)

type Devices struct {
	db models.DeviceDB
}

func NewDevices(db models.DeviceDB) *Devices {
	return &Devices{
		db: db,
	}
}

// define struct for mac address json
type MacAddress struct {
	Id string
}

func (db *Devices) RetrieveByMac() gin.HandlerFunc {
	return func(context *gin.Context) {
		var mac MacAddress
		ifas, err := net.Interfaces()
		if err != nil {
			engine.Warning.Println(err)
			ResponseMap["data"] = gin.H{"error": err.Error()}
			ResponseMap["message"] = "failed"
			response = engine.Response{
				ResponseCode: http.StatusInternalServerError,
				Context:      context,
				Response:     ResponseMap,
			}
			response.SendAnswer()
			return
		}
		for _, ifa := range ifas {
			a := ifa.HardwareAddr.String()
			if a != "" {
				mac = MacAddress{Id: a}
				break
			}
		}
		_, err = db.db.ByMac(mac.Id)
		switch err {
		case engine.ERR_NOT_FOUND:
			ResponseMap["data"] = gin.H{"error": engine.ERR_NOT_ENOUGH_PERMISSIONS}
			ResponseMap["message"] = "failed"
			response = engine.Response{
				ResponseCode: http.StatusForbidden,
				Context:      context,
				Response:     ResponseMap,
			}
			response.SendAnswer()
			return
		case nil:

		default:
			engine.Warning.Println(err)
			ResponseMap["data"] = gin.H{"error": err.Error()}
			ResponseMap["message"] = "failed"
			response = engine.Response{
				ResponseCode: http.StatusBadRequest,
				Context:      context,
				Response:     ResponseMap,
			}
			response.SendAnswer()
			return
		}
		context.Next()
	}
}
func (db *Devices) ValidToken() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("neftAuth")

		if tokenString != "devtest.dev" {
			engine.Debug.Println("TRY TO ACCES BETA WITHOUT PERMISSIONS")
      err := engine.ERR_NOT_ENOUGH_PERMISSIONS
			engine.Warning.Println(err)
			ResponseMap["data"] = gin.H{"error": err.Error()}
			ResponseMap["message"] = "failed"
			response = engine.Response{
				ResponseCode: http.StatusBadRequest,
				Context:      context,
				Response:     ResponseMap,
			}
			response.SendAnswer()
			return
		}
		context.Next()
	}
}
