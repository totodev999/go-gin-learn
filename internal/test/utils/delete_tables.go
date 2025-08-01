package test_utils

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

func DeleteTables(db *gorm.DB) {
	var tables []string
	db.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Scan(&tables)

	fmt.Printf("Delete table %v", tables)

	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE;", table)
		if err := db.Exec(query).Error; err != nil {
			log.Fatalf("Failed to delete from %s: %v\n", table, err)
		}
	}

}
