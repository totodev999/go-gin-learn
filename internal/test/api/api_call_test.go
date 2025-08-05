package api_test

import (
	"encoding/json"
	"flea-market/internal/app"
	"flea-market/repositories"
	"net/http"
	"net/http/httptest"
	"testing"

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

func TestGetAllPosts_External(t *testing.T) {

	// url would be like 127.0.0.1:35889, port differs every time starting mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// returning data only when a request is GET "/post".
		if r.URL.Path == "/posts" && r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(dummyPosts)
			return
		}
		http.NotFound(w, r)
	}))
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
	assert.Equal(t, *dummyPosts[0].Title, *res["data"][0].Title)
	assert.Equal(t, *dummyPosts[0].UserId, *res["data"][0].UserId)
}
