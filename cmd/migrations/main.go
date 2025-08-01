package main

import (
	"flea-market/infra"
	"flea-market/models"
	"log"
)

func main() {
	infra.Initializer()
	db := infra.SetupDB()

	if err := db.AutoMigrate(&models.User{}, &models.Item{}); err != nil {
		log.Fatalln("Failed to migrate database")
	}
}
