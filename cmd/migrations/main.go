package main

import (
	"flea-market/infra"
	"flea-market/models"
)

func main() {
	infra.Initializer()
	db := infra.SetupDB()

	if err := db.AutoMigrate(&models.Item{}, &models.User{}); err != nil {
		panic("Failed to migrate database")
	}
}
