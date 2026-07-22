package main

import (
	"log"

	"github.com/go-api/internal/bootstrap"
	"github.com/go-api/internal/interfaces/http/routes"
	"go.uber.org/zap"
)

func main() {

	app, err := bootstrap.New()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := app.Close(); err != nil {
			app.Logger.Error("failed to shutdown application", zap.Error(err))
		}
	}()

	router := routes.SetupRouter(app.Logger, app.Container)
	router.Static("/storage", "./storage")

	app.Logger.Info("application started")

	if err := router.Run(":" + app.Config.App.Port); err != nil {
		log.Fatal(err)
	}
}
