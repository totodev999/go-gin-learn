package main

import (
	"flea-market/infra"
	"fmt"
	"log"
)

func main() {
	infra.Initializer()
	db := infra.SetupDB()

	var tables []string
	db.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Scan(&tables)

	fmt.Printf("Delete table %v", tables)

	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE;", table)
		if err := db.Exec(query).Error; err != nil {
			log.Fatalf("Failed to delete from %s: %v", table, err)
		}
	}

}
