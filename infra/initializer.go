package infra

import (
	"log"

	"github.com/joho/godotenv"
)

func Initializer() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found; relying on environment variables")
	}
}
