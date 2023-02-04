package helper

import (
	engine "github.com/JoanGTSQ/api"
	"neft.web/controller"
)

func GenerateController(handlerName string, funcType controller.ClientCommandExecution) {
	engine.Debug.Printf("New handler registered %s", handlerName)
	controller.MapFuncs[handlerName] = funcType
}
