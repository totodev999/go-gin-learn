package api_test

import (
	"encoding/json"
	"flea-market/internal/app"
	"flea-market/repositories"
	"flea-market/utils"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// 外部API用のダミーデータ
var dummyPosts = []repositories.Post{
	{
		Id:     uintPtr(1),
		Title:  strPtr("test title"),
		Body:   strPtr("test body"),
		UserId: intPtr(99),
		Dummy:  strPtr("dummy"),
	},
}

func uintPtr(i uint) *uint    { return &i }
func intPtr(i int) *int       { return &i }
func strPtr(s string) *string { return &s }

func setupAPICallTest() *gin.Engine {
	db := testDB
	router := app.NewRouter(db)
	return router
}

func setuoAPICallMockServer(statusCode int, sleepSec uint) *httptest.Server {
	var response any

	if statusCode == http.StatusOK {
		response = dummyPosts
	} else {
		response = map[string]string{
			"message": "something wrong",
		}
	}

	// url would be like 127.0.0.1:35889, port differs every time starting mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// returning data only when a request is GET "/post".
		if r.URL.Path == "/posts" && r.Method == http.MethodGet {
			if sleepSec != 0 {
				time.Sleep(time.Duration(sleepSec) * time.Second)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(statusCode)
			_ = json.NewEncoder(w).Encode(response)
			return
		}
		http.NotFound(w, r)
	}))
	return mockServer
}

func TestGetAllPosts(t *testing.T) {

	// url would be like 127.0.0.1:35889, port differs every time starting mock server
	mockServer := setuoAPICallMockServer(200, 0)
	defer mockServer.Close()

	t.Setenv("BASE_URL", mockServer.URL)

	router := setupAPICallTest()

	req := httptest.NewRequest("GET", "/external", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res map[string][]repositories.Post
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Equal(t, len(dummyPosts), len(res["data"]))
	assert.Equal(t, dummyPosts[0].Title, res["data"][0].Title)
	assert.Equal(t, dummyPosts[0].UserId, res["data"][0].UserId)
}

func TestGetAllPosts_Timeout(t *testing.T) {

	// url would be like 127.0.0.1:35889, port differs every time starting mock server
	mockServer := setuoAPICallMockServer(200, 4)
	defer mockServer.Close()

	t.Setenv("BASE_URL", mockServer.URL)

	router := setupAPICallTest()

	req := httptest.NewRequest("GET", "/external", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var res map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Equal(t, "Internal server error", res["error"])
}

func TestGetAllPosts_Request_URL_Wrong(t *testing.T) {

	// url would be like 127.0.0.1:35889, port differs every time starting mock server
	mockServer := setuoAPICallMockServer(200, 0)
	defer mockServer.Close()

	t.Setenv("BASE_URL", "127.0.0.1:0")

	router := setupAPICallTest()

	req := httptest.NewRequest("GET", "/external", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var res map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Equal(t, string(utils.ExternalAPIConnectionError), res["error"])
}

func TestGetAllPosts_Get_HTTPStatus_400(t *testing.T) {

	// url would be like 127.0.0.1:35889, port differs every time starting mock server
	mockServer := setuoAPICallMockServer(400, 0)
	defer mockServer.Close()

	t.Setenv("BASE_URL", mockServer.URL)

	router := setupAPICallTest()

	req := httptest.NewRequest("GET", "/external", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var res map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Equal(t, string(utils.ExternalAPIReturnsError), res["error"])
}
