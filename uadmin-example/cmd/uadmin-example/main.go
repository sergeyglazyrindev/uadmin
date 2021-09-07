package main

import (
	"github.com/sergeyglazyrindev/uadmin"
	"github.com/sergeyglazyrindev/uadmin_example/blueprint/example"
	"os"
)


func main() {
	environment := os.Getenv("environment")
	if environment == "" {
		environment = "dev"
	}
	app1 := uadmin.NewApp(environment, true)
	app1.BlueprintRegistry.Register(example.ConcreteBlueprint)
	app1.Initialize()
	app1.InitializeRouter()
	app1.ExecuteCommand()
}

