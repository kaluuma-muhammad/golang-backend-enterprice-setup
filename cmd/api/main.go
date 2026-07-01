package main

import (
	"log"

	"github.com/go-api/internal/bootstrap"
	"github.com/go-api/internal/interfaces/http/routes"
)

func main() {

	app, err := bootstrap.New()
	if err != nil {
		log.Fatal(err)
	}

	router := routes.SetupRouter(app.Logger, app.Container)

	app.Logger.Info("application started")

	if err := router.Run(":" + app.Config.App.Port); err != nil {
		log.Fatal(err)
	}
}
