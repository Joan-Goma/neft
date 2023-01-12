package main

import (
	"flag"
	"runtime"

	engine "github.com/JoanGTSQ/api"
	"github.com/gin-gonic/gin"
	"neft.web/client"
	"neft.web/helper"
	"neft.web/models"
)

var (
	isProd  bool
	port    string
	debug   bool
	debugDB bool
	ssl     bool
	route   string
)

const version = "CERBERUS 2.0"

func init() {
	flag.BoolVar(&isProd, "isProd", false, "This will ensure all pro vars are enabled")
	flag.BoolVar(&debug, "debug", false, "This will export all stats to file log.log")
	flag.BoolVar(&debugDB, "debugDB", false, "This will enable logs of db")
	flag.StringVar(&port, "port", ":8080", "This will set the port of use")
	flag.StringVar(&route, "route", "log.txt", "This will create the log file in the desired route")
	gin.SetMode(gin.ReleaseMode)
	flag.Parse()
	engine.InitLog(debug, route, version)
	helper.InitConfig()
}

func main() {

	if err := helper.InitDB(debugDB); err != nil {
		engine.Error.Fatalln("Can not connect to DB: ", err)
	}

	defer func(Services *models.Services) {
		err := Services.Close()
		if err != nil {
			engine.Error.Fatalln("Error deferring db close", err)
		}
	}(client.Services)

	// Auto generate new tables or modifications in every start | Use DestructiveReset() to delete all data

	if err := client.Services.AutoMigrate(); err != nil {
		engine.Error.Fatalln("Can not AutoMigrate the database", err)
	}

	// Retrieve controllers struct
	engine.Debug.Println("Creating all services handlers")

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

	if !debug {
		// Start server
		engine.Info.Println("Running SSL server on port :443")
		if err := r.RunTLS(":443", cer, key); err != nil {
			engine.Error.Fatalln("Error starting web server", err)
		}
	} else {
		engine.Info.Println("Running non SSL server on port", port)
		if err := r.Run(port); err != nil {
			engine.Error.Fatalln("Error starting web server", err)
		}
	}

}
