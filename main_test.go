package main

import (
	"encoding/json"
	"free-market/infra"
	"free-market/models"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load(".env.test"); err != nil {
		log.Fatal("Error loading .env.test")
	}

	code := m.Run()

	os.Exit(code)
}

var itemData = []models.Item{
	{Name: "test1", Price: 100, Description: "", SoldOut: false, UserID: 1},
	{Name: "test2", Price: 200, Description: "テスト2", SoldOut: true, UserID: 1},
	{Name: "test3", Price: 300, Description: "テスト3", SoldOut: false, UserID: 2},
}

var userData = []models.User{
	{Email: "test1@test.com", Password: "testpass"},
	{Email: "test2@test.com", Password: "testpass"},
}

func setupTestData(db *gorm.DB) {
	items := itemData

	users := userData

	for _, user := range users {
		db.Create(&user)
	}

	for _, item := range items {
		db.Create(&item)
	}

}
func setup() *gin.Engine {
	db := infra.SetupDB()
	db.AutoMigrate(&models.User{}, &models.Item{})

	setupTestData(db)

	router := setUpRouter(db)

	return router
}

type ItemForTest struct {
	Name        string
	Price       uint
	Description string
	SoldOut     bool
	UserID      uint
}

func toTestItems(items []models.Item) []ItemForTest {
	var res []ItemForTest
	for _, item := range items {
		res = append(res, ItemForTest{
			Name:        item.Name,
			Price:       item.Price,
			Description: item.Description,
			SoldOut:     item.SoldOut,
			UserID:      item.UserID,
		})
	}
	return res
}

func TestFindAll(t *testing.T) {
	router := setup()

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/items", nil)

	router.ServeHTTP(w, req)

	var res map[string][]models.Item
	json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 3, len(res["data"]))

	assert.ElementsMatch(t, toTestItems(itemData), toTestItems(res["data"]))

	// Reflect Ver little bit too complicated

	// for index, v := range res["data"] {
	// 	expectType := reflect.TypeOf(v)
	// 	expectVal := reflect.ValueOf(v)

	// 	for i := 0; i < expectType.NumField(); i++ {
	// 		field := expectType.Field(i).Name
	// 		value := expectVal.FieldByName(field).Interface()
	// 		assert.Equal(t, value, reflect.ValueOf(res["data"][index]).FieldByName(field).Interface())
	// 	}
	// }

}
