package main

import (
	"context"
	"log"

	"github.com/go-api/internal/bootstrap"
	"github.com/go-api/internal/bootstrap/seed"
)

func main() {

	app, err := bootstrap.New()
	if err != nil {
		log.Fatal(err)
	}

	defer app.DB.Close()

	if err := seed.Run(context.Background(), app.Container, app.Config); err != nil {
		log.Fatal(err)
	}

	log.Println("✅ Database seeded successfully.")
}
