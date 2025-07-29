package main

import (
	"flea-market/internal/app"
	"log"
)

func main() {
	application := app.NewApp()
	if err := application.Run(); err != nil {
		log.Fatalf("Starting server failed: %v", err)
	}
}
