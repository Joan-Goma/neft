package helper

import (
	"encoding/json"
	engine "github.com/JoanGTSQ/api"
	"io/ioutil"
	"os"
    "runtime"
)

type dbConnection struct {
	URL      string `json:"dbDirection"`
	User     string `json:"dbUser"`
	Name     string `json:"dbName"`
	Password string `json:"dbPsswd"`
	SslMode  string `json:"sslMode"`
}

var database dbConnection

func InitConfig() {
	var route string
    if runtime.GOOS == "linux" {
        route = "/srv/"
    }
    // Open our jsonFile
	jsonFile, err := os.Open(route + "env.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		engine.Error.Fatalln(err)
	}
	engine.Info.Println("Successfully loaded enviroment configuration")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &database)
}
