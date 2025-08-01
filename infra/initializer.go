package infra

import (
	"flea-market/utils"

	"github.com/joho/godotenv"
)

func Initializer() {
	if err := godotenv.Load(); err != nil {
		utils.Logger(utils.GenericMessage, nil, ".env file not found; relying on environment variables")
	}
}
