package main

import (
	"flea-market/infra"
	"fmt"
	"log"
	"strings"
)

func main() {
	infra.Initializer()
	db := infra.SetupDB()

	var tables []string
	db.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Scan(&tables)

	for _, table := range tables {
		query := fmt.Sprintf(`DELETE FROM "%s";`, table)
		if err := db.Exec(query).Error; err != nil {
			log.Printf("Failed to delete from %s: %v", table, err)
		}
	}
	fmt.Printf("全テーブルのデータ削除完了 削除したテーブル：%v\n", strings.Join(tables, ", "))
}
