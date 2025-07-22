package infra

import (
	"os"
	"testing"
)

func TestSetupDB_PanicOnError(t *testing.T) {
	originalEnv := os.Getenv("ENV")
	defer os.Setenv("ENV", originalEnv) // ←必ず元に戻す

	os.Setenv("ENV", "production")
	// 必ず失敗するような値を入れる
	os.Setenv("ENV", "production") // testじゃないのでPostgresパス
	os.Setenv("DB_HOST", "invalid_host")
	os.Setenv("DB_USER", "invalid_user")
	os.Setenv("DB_PASSWORD", "invalid_password")
	os.Setenv("DB_NAME", "invalid_db")
	os.Setenv("DB_PORT", "1234")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic but did not panic")
		}
	}()

	SetupDB()
}
