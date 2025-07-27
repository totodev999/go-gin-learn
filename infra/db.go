package infra

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupDB() *gorm.DB {
	env := os.Getenv("ENV")
	dns := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var (
		db  *gorm.DB
		err error
	)

	if env == "test" {
		db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		fmt.Println("using sqlite")
	} else {
		db, err = gorm.Open(postgres.Open(dns), &gorm.Config{})
		fmt.Println("using postgres")
	}

	if err != nil {
		errMsg := fmt.Sprintf("failed to connect database %v\n", err)
		panic(errMsg)
	}

	return db
}
