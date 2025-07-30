package infra

import (
	"flea-market/utils"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDB() *gorm.DB {
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

	db, err = gorm.Open(postgres.Open(dns), &gorm.Config{})
	fmt.Println("using postgres")
	if err != nil {
		cstmErr := utils.NewDBError("Connecting db error", err)
		utils.Logger(cstmErr.MessageCode, nil, cstmErr)
		panic(cstmErr.Error())
	}

	return db
}
