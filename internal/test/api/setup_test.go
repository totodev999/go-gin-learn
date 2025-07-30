package api_test

import (
	"flea-market/infra"
	test_utils "flea-market/internal/test/utils"
	"os"
	"testing"

	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	test_utils.ReadEnv()
	testDB = infra.SetupDB()

	code := m.Run()

	os.Exit(code)
}
