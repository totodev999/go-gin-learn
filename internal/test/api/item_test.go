package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"flea-market/controllers"
	"flea-market/dto"
	"flea-market/infra"
	"flea-market/internal/app"
	"flea-market/internal/mocks"
	test "flea-market/internal/test/utils"
	"flea-market/middlewares"
	"flea-market/models"
	"flea-market/repositories"
	"flea-market/services"
	"flea-market/utils"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working dir: %v", err)
	}

	root := test.FindProjectRoot(wd)
	envPath := filepath.Join(root, ".env.test")

	if err := godotenv.Load(envPath); err != nil {
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
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Item{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.User{})

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
	// This is necessary for SQLite. Because every time this function is called, new SQLite instance is created.
	db.Exec("DROP TABLE IF EXISTS users;")
	db.Exec("DROP TABLE IF EXISTS items;")

	db.AutoMigrate(&models.User{}, &models.Item{})
	setupTestData(db)
	router := app.NewRouter(db)
	return router
}

func setUpRouterWithItemRepo(itemRepo repositories.IItemRepository) *gin.Engine {
	itemService := services.NewItemService(itemRepo)
	itemController := controllers.NewItemController(itemService)

	router := gin.New()
	router.Use(middlewares.APIErrorHandler())
	router.Use(gin.Recovery())
	itemRouter := router.Group("/items")
	itemRouter.GET("", itemController.FindAll)

	return router
}

type ItemForTest struct {
	Name        string
	Price       uint
	Description string
	SoldOut     bool
	UserID      uint
}

func toTestItem(item models.Item) ItemForTest {
	return ItemForTest{
		Name:        item.Name,
		Price:       item.Price,
		Description: item.Description,
		SoldOut:     item.SoldOut,
		UserID:      item.UserID,
	}
}

func toTestItems(items []models.Item) []ItemForTest {
	var res []ItemForTest
	for _, item := range items {
		res = append(res, toTestItem(item))
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

func TestFindById(t *testing.T) {
	router := setup()

	token, _ := services.CreateToken(1, userData[0].Email)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/items/1", nil)
	req.Header.Set("Authorization", "Bearer "+*token)

	router.ServeHTTP(w, req)

	var res map[string]models.Item
	json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, uint(1), res["data"].ID)

	assert.Equal(t, toTestItem(itemData[0]), toTestItem(res["data"]))
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
	exp := models.Item{
		Name:        createItemInput.Name,
		Price:       createItemInput.Price,
		Description: createItemInput.Description,
		SoldOut:     false,
		UserID:      1,
	}
	assert.Equal(t, toTestItem(exp), toTestItem(res["data"]))

}

func TestUpdate(t *testing.T) {
	router := setup()

	token, err := services.CreateToken(1, userData[0].Email)
	assert.Equal(t, nil, err)
	reqBody := `{"name":"12","price":999,"description":"updated","soldOut":true}`

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("PUT", "/items/1", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Authorization", "Bearer "+*token)

	router.ServeHTTP(w, req)

	var res map[string]models.Item
	json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, uint(1), res["data"].ID)

	var expexted models.Item
	if err := json.Unmarshal([]byte(reqBody), &expexted); err != nil {
		t.Fatalf("json.Unmarshal failed %v\n", err)
	}

	exp := toTestItem(expexted)
	got := toTestItem(res["data"])
	exp.UserID = 0
	got.UserID = 0

	assert.Equal(t, exp, got)

}

func Test_Delete(t *testing.T) {
	router := setup()

	token, err := services.CreateToken(1, userData[0].Email)
	assert.Equal(t, nil, err)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("DELETE", "/items/1", nil)
	req.Header.Set("Authorization", "Bearer "+*token)

	router.ServeHTTP(w, req)

	var res map[string]models.Item
	json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(t, http.StatusOK, w.Code)

	// Check if the record is deleted.
	w2 := httptest.NewRecorder()
	reqGet, _ := http.NewRequest("GET", "/items/1", nil)
	reqGet.Header.Set("Authorization", "Bearer "+*token)
	router.ServeHTTP(w2, reqGet)

	assert.Equal(t, http.StatusNotFound, w2.Code)

}

func Test_FindById_Wrong_ID(t *testing.T) {
	router := setup()

	token, err := services.CreateToken(1, userData[0].Email)
	assert.Equal(t, nil, err)

	cases := []struct {
		name       string
		param      string
		wantStatus int
	}{
		{
			name:       "param is string",
			param:      "id",
			wantStatus: http.StatusBadRequest,
		},
		// Gin will treat /items and /items/ differently. /items/ will be redirected to /items
		{
			name:       "param is missing",
			param:      "",
			wantStatus: http.StatusMovedPermanently,
		},
		{
			name:       "param is number but no data found",
			param:      "9999999",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/items/%v", tc.param), nil)
			req.Header.Set("Authorization", "Bearer "+*token)

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)
		})
	}

}

func Test_Delete_Wrong_ID(t *testing.T) {
	router := setup()

	token, err := services.CreateToken(1, userData[0].Email)
	assert.Equal(t, nil, err)

	cases := []struct {
		name       string
		param      string
		wantStatus int
	}{
		{
			name:       "param is string",
			param:      "id",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "param is missing",
			param:      "",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "param is number but no data found",
			param:      "9999999",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/items/"+tc.param, nil)
			req.Header.Set("Authorization", "Bearer "+*token)

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)
		})
	}

}

func Test_Create_Wrong_Input(t *testing.T) {
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
			body:       `{"name":"12",price":"文字列","description":"test"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "name is less than 2",
			body:       `{"name":"a","price":200,"description":""}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing price",
			body:       `{"name":"12","description":"test"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "extra unknown field",
			body:       `{"name":"12","price":200,"description":"test","unknown":"field"}`,
			wantStatus: http.StatusCreated,
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

func Test_Update_Wrong_Input(t *testing.T) {
	router := setup()

	token, err := services.CreateToken(1, userData[0].Email)
	assert.Equal(t, nil, err)

	cases := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "Name is less than 2",
			body:       `{"Name":"1","description":""}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "price is string",
			body:       `{"price":"文字列","description":"test"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "price is less than 1",
			body:       `{"price":0,"description":"test"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "price is greater than 9999999",
			body:       `{"price":1000000,"description":"test"}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/items/1", strings.NewReader(tc.body))
			req.Header.Set("Authorization", "Bearer "+*token)

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)
		})
	}

}

// auth_middleware will check the Authorization token.
func Test_Unauthorized(t *testing.T) {
	router := setup()
	createItemInput := dto.CreateItemInput{
		Name:        "test item New",
		Price:       2010,
		Description: "This is test",
	}
	createItemInputBytes, _ := json.Marshal(createItemInput)

	cases := []struct {
		path       string
		method     string
		body       string
		wantStatus int
	}{
		{
			path:       "/items",
			method:     "POST",
			body:       string(createItemInputBytes),
			wantStatus: http.StatusUnauthorized,
		},
		// No authorization is required for FindAll
		{
			path:       "/items",
			method:     "GET",
			body:       "",
			wantStatus: http.StatusOK,
		},
		{
			path:       "/items/1",
			method:     "GET",
			body:       "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			path:       "/items/1",
			method:     "DELETE",
			body:       "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			path:       "/items/1",
			method:     "PUT",
			body:       string(createItemInputBytes),
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range cases {
		t.Run(tc.path+":"+tc.method, func(t *testing.T) {
			w := httptest.NewRecorder()

			req, _ := http.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.wantStatus, w.Code)
		})
	}

}

func TestItems_FindAll_DBError_Non_CustomError(t *testing.T) {
	mockRepo := &mocks.MockItemRepository{
		FindAllFunc: func() (*[]models.Item, error) {
			return nil, errors.New("mock db error")
		},
	}

	router := setUpRouterWithItemRepo(mockRepo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/items", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Internal server error")
}

func TestItems_FindAll_DBError_CustomError(t *testing.T) {
	mockRepo := &mocks.MockItemRepository{
		FindAllFunc: func() (*[]models.Item, error) {
			return nil, utils.NewDBError("Mock", errors.New("Mock"))
		},
	}

	router := setUpRouterWithItemRepo(mockRepo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/items", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), utils.DBError)
}

func Test_Unauthorized_Invalid_Token(t *testing.T) {
	router := setup()
	req, _ := http.NewRequest("GET", "/items/1", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func Test_Forbidden_Access_OtherUserItem(t *testing.T) {
	router := setup()

	token, _ := services.CreateToken(2, userData[1].Email)
	req, _ := http.NewRequest("DELETE", "/items/1", nil)
	req.Header.Set("Authorization", "Bearer "+*token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.True(t, w.Code == http.StatusNotFound)
}

func Test_Forbidden_Update_OtherUserItem(t *testing.T) {
	router := setup()
	token, _ := services.CreateToken(2, userData[1].Email)
	reqBody := `{"name":"test update","price":111,"description":"try update"}`
	req, _ := http.NewRequest("PUT", "/items/1", strings.NewReader(reqBody))
	req.Header.Set("Authorization", "Bearer "+*token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func Test_Update_DeletedItem(t *testing.T) {
	router := setup()
	token, _ := services.CreateToken(1, userData[0].Email)
	// まず削除
	reqDel, _ := http.NewRequest("DELETE", "/items/1", nil)
	reqDel.Header.Set("Authorization", "Bearer "+*token)
	wDel := httptest.NewRecorder()
	router.ServeHTTP(wDel, reqDel)

	// 削除済みIDで更新
	reqUpd, _ := http.NewRequest("PUT", "/items/1", strings.NewReader(`{"name":"test","price":123,"description":""}`))
	reqUpd.Header.Set("Authorization", "Bearer "+*token)
	wUpd := httptest.NewRecorder()
	router.ServeHTTP(wUpd, reqUpd)
	assert.Equal(t, http.StatusNotFound, wUpd.Code)
}

func Test_AllEndpoints_Concurrent(t *testing.T) {
	if !test.RaceEnabled {
		t.Skip("skip: race detector not enabled")
	}
	router := setup()
	token, err := services.CreateToken(1, userData[0].Email)
	assert.NoError(t, err)

	var wg sync.WaitGroup
	threadNum := 10

	// --- Createを並列実行
	for i := 0; i < threadNum; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			body := dto.CreateItemInput{
				Name:        fmt.Sprintf("RaceItem%d", i),
				Price:       uint(2000 + i),
				Description: fmt.Sprintf("desc%d", i),
			}
			bodyBytes, _ := json.Marshal(body)
			req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Authorization", "Bearer "+*token)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusCreated, w.Code)
		}(i)
	}

	// --- Updateを並列実行（ID=1固定と仮定）
	for i := 0; i < threadNum; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			name := fmt.Sprintf("UpdatedByThread%d", i)
			price := uint(3000 + i)
			description := fmt.Sprintf("desc_update%d", i)
			soldOut := i%2 == 0
			body := dto.UpdateItemInput{
				Name:        &name,
				Price:       &price,
				Description: &description,
				SoldOut:     &soldOut,
			}
			bodyBytes, _ := json.Marshal(body)
			req, _ := http.NewRequest("PUT", "/items/1", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Authorization", "Bearer "+*token)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			// 既に消されてる場合も考慮して2xx/404でも許可
			assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusNotFound)
		}(i)
	}

	// --- Deleteを並列実行（ID=1固定と仮定）
	for i := 0; i < threadNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req, _ := http.NewRequest("DELETE", "/items/1", nil)
			req.Header.Set("Authorization", "Bearer "+*token)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			// 既に消されてる場合も考慮して2xx/404でも許可
			assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusNotFound)
		}()
	}

	// --- FindAllを並列実行
	for i := 0; i < threadNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req, _ := http.NewRequest("GET", "/items", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}()
	}

	// --- FindById(ID=1固定)を並列実行
	for i := 0; i < threadNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req, _ := http.NewRequest("GET", "/items/1", nil)
			req.Header.Set("Authorization", "Bearer "+*token)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			// 更新/削除が同時なのでOK/404どちらでも
			assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusNotFound)
		}()
	}

	// 競合を高めるため少しsleep
	time.Sleep(50 * time.Millisecond)

	wg.Wait()

	// 最終的に残ってるitemsを確認
	req, _ := http.NewRequest("GET", "/items", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var res map[string][]models.Item
	json.Unmarshal(w.Body.Bytes(), &res)
	t.Logf("final items: %+v", res["data"])
}
