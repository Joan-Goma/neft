package main

import (
	"flag"
	"runtime"

	"neft.web/models"

	engine "github.com/JoanGTSQ/api"
	"github.com/gin-gonic/gin"
	"neft.web/helper"
)

var (
	port    string
	debug   bool
	debugDB bool
	route   string
)

const version = "CERBERUS DEV 3.2.0"

func init() {
	flag.BoolVar(&debug, "debug", false, "Enable debug prints")
	flag.BoolVar(&debugDB, "debugDB", false, "Enable logs of database")
	flag.StringVar(&port, "port", ":8080", "Set the port of the web server")
	flag.StringVar(&route, "route", "log.txt", "Set the log route and the file name")
	gin.SetMode(gin.ReleaseMode)
	flag.Parse()
	engine.InitLog(debug, route, version)
	helper.InitConfig()
}

func main() {

	if err := helper.InitDB(debugDB); err != nil {
		engine.Error.Fatalln("Can not connect to DB: ", err)
	}

	/*defer func(Services *models.Services) {
		err := Services.Close()
		if err != nil {
			engine.Error.Fatalln("Error deferring db close", err)
		}
	}(client.Services)*/

	// Auto generate new tables or modifications in every start | Use DestructiveReset() to delete all data

	if err := models.AutoMigrate(); err != nil {
		engine.Error.Fatalln("Can not AutoMigrate the database", err)
	}

	// Retrieve controllers struct

	// Generate Router
	r := helper.InitRouter()

	// Start printing stats
	go engine.PrintStats()
	go helper.ReadInput(debug)

	var route string

	if runtime.GOOS == "linux" {
		engine.Info.Println("Welcome Linux user:", runtime.GOOS, runtime.GOARCH)
		route = "/srv/"
	} else {
		engine.Info.Println("Welcome Windows user:", runtime.GOOS, runtime.GOARCH)
	}

	var cer = route + "cer.cer"
	var key = route + "key.key"

	if port == ":443" {
		engine.Info.Println("Running SSL server on port :443")
		if err := r.RunTLS(":443", cer, key); err != nil {
			engine.Error.Fatalln("Error starting web server", err)
		}
	}

	engine.Info.Println("Running non SSL server on port", port)
	if err := r.Run(port); err != nil {
		engine.Error.Fatalln("Error starting web server", err)
	}

}
