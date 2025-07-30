package api_test

import (
	"bytes"
	"encoding/json"
	"flea-market/dto"
	"flea-market/internal/app"
	test_utils "flea-market/internal/test/utils"
	"flea-market/models"
	"flea-market/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func setupAuthTestData(db *gorm.DB) {
	test_utils.DeleteTables(db)

	//　if you write like "db.Create(&test_utils.UserData)",
	// it may cause a problem. &test_utils.UserData means pass pointer to gorm,
	// and gorm can rewrite "&test_utils.UserData".
	// So ID and other fields can be set and cause a problem following.
	// duplicate key value violates unique constraint "users_pkey"
	for _, user := range test_utils.UserData {
		db.Create(&user)
	}

}

func setupAuthTest() *gin.Engine {
	db := testDB
	setupAuthTestData(db)
	router := app.NewRouter(db)
	return router
}

func TestSignupAndLogin(t *testing.T) {
	// 1. Request to /auth/signup
	router := setupAuthTest()

	signupW := httptest.NewRecorder()

	signUpInput := dto.SignupInput{
		Email:    "orenotest@test.com",
		Password: "nikutaberu",
	}
	reqBody, _ := json.Marshal(signUpInput)

	signupReq, _ := http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(reqBody))

	router.ServeHTTP(signupW, signupReq)

	// Check /auth/signup response and DB and hashed password is comparable
	assert.Equal(t, http.StatusCreated, signupW.Code)

	db := testDB
	var userResult models.User
	db.First(&userResult, "email = ?", signUpInput.Email)
	assert.Equal(t, userResult.Email, signUpInput.Email)
	err := bcrypt.CompareHashAndPassword([]byte(userResult.Password), []byte(signUpInput.Password))
	assert.Equal(t, err, nil)

	// 2. Request to /auth/login to check if the created user can login soon after signing u
	loginW := httptest.NewRecorder()
	loginReq, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
	router.ServeHTTP(loginW, loginReq)

	// Check /auth/signup response(token) and token can have the email used to login
	assert.Equal(t, http.StatusOK, loginW.Code)

	var res map[string]string
	json.Unmarshal(loginW.Body.Bytes(), &res)
	token, _, _ := jwt.NewParser().ParseUnverified(res["token"], jwt.MapClaims{})

	claims, _ := token.Claims.(jwt.MapClaims)

	assert.Equal(t, claims["email"].(string), signUpInput.Email)
}

func TestSignupAndLoginWithWrongInput(t *testing.T) {
	router := setupAuthTest()

	cases := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "Email is invalid",
			body:       `{"email":"test1test.com","password":"1234567"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Password is less then 8",
			body:       `{"email":"test1@test.com","password":"1234567"}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			signupW := httptest.NewRecorder()
			signupReq, _ := http.NewRequest("POST", "/auth/signup", strings.NewReader(tc.body))
			router.ServeHTTP(signupW, signupReq)
			assert.Equal(t, tc.wantStatus, signupW.Code)

			loginW := httptest.NewRecorder()
			loginReq, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(tc.body))
			router.ServeHTTP(loginW, loginReq)
			assert.Equal(t, tc.wantStatus, loginW.Code)

		})
	}

}

func TestSignupWithDuplicatedEmail(t *testing.T) {
	router := setupAuthTest()

	time.Sleep(1000 * time.Millisecond)

	input := dto.SignupInput{
		Email:    test_utils.UserData[0].Email,
		Password: "12345678",
	}

	w := httptest.NewRecorder()
	reqBody, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(reqBody))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	var res map[string]string
	json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(t, res["error"], string(utils.DuplicateKeyError))
}

func TestLoginWithWrongPassword(t *testing.T) {
	// 1. Request to /auth/signup to register User
	// At setup, users are already registered but password is not hashed.
	// So use /auth/signup to register a user in a proper way.
	router := setupAuthTest()

	signupW := httptest.NewRecorder()

	const email = "tesorenotest@test.com"

	signUpInput := dto.SignupInput{
		Email:    email,
		Password: "nikutaberu",
	}
	reqBody, _ := json.Marshal(signUpInput)

	signupReq, _ := http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(reqBody))

	router.ServeHTTP(signupW, signupReq)

	// Check /auth/signup response and DB and hashed password is comparable
	assert.Equal(t, http.StatusCreated, signupW.Code)

	loginInput := dto.SignupInput{
		Email:    email,
		Password: "12345678",
	}

	loginW := httptest.NewRecorder()
	loginReqBody, _ := json.Marshal(loginInput)
	loginReq, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(loginReqBody))
	router.ServeHTTP(loginW, loginReq)

	assert.Equal(t, http.StatusUnauthorized, loginW.Code)

	var res map[string]string
	json.Unmarshal(loginW.Body.Bytes(), &res)

	assert.Equal(t, res["error"], string(utils.UnAuthorized))
}
