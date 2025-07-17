package infra

import (
	"log"

	"github.com/joho/godotenv"
)

func Initializer() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
