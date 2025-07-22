package main

import (
	"bytes"
	"encoding/json"
	"free-market/dto"
	"free-market/infra"
	"free-market/models"
	"free-market/services"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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

func TestCreate(t *testing.T) {
	router := setup()

	token, err := services.CreateToken(1, userData[0].Email)
	assert.Equal(t, nil, err)

	createItemInput := dto.CreateItemInput{
		Name:        "test item New",
		Price:       2010,
		Description: "This is test",
	}

	reqBody, _ := json.Marshal(createItemInput)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+*token)

	router.ServeHTTP(w, req)

	var res map[string]models.Item
	json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, uint(4), res["data"].ID)

}

func TestCreateWrongInput(t *testing.T) {
	router := setup()

	token, err := services.CreateToken(1, userData[0].Email)
	assert.Equal(t, nil, err)

	cases := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "price is string",
			body:       `{"price":"文字列","description":"test"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "description is empty",
			body:       `{"price":200,"description":""}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing price",
			body:       `{"description":"test"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "extra unknown field",
			body:       `{"price":200,"description":"test","unknown":"field"}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/items", strings.NewReader(tc.body))
			req.Header.Set("Authorization", "Bearer "+*token)

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)
		})
	}

}

func TestCreateUnauthorized(t *testing.T) {
	router := setup()

	createItemInput := dto.CreateItemInput{
		Name:        "test item New",
		Price:       2010,
		Description: "This is test",
	}

	reqBody, _ := json.Marshal(createItemInput)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer(reqBody))

	router.ServeHTTP(w, req)

	var res map[string]models.Item
	json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
