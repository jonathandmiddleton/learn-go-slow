package main

import (
	"github.com/and-cru/go-service/api/app"
	"github.com/and-cru/go-service/api/config"
)

func main() {
	config := config.GetConfig()

	app := &app.App{}
	app.Initialize(config)
	app.Run(":3000")
}
