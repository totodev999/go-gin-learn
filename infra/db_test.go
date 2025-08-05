package infra

import (
	"testing"
)

func TestSetupDB_PanicOnError(t *testing.T) {
	t.Setenv("ENV", "production") // testじゃないのでPostgresパス
	t.Setenv("DB_HOST", "invalid_host")
	t.Setenv("DB_USER", "invalid_user")
	t.Setenv("DB_PASSWORD", "invalid_password")
	t.Setenv("DB_NAME", "invalid_db")
	t.Setenv("DB_PORT", "1234")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic but did not panic")
		}
	}()

	SetupDB()
}
