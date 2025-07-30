package services

import (
	test_utils "flea-market/internal/test/utils"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	test_utils.ReadEnv()

	code := m.Run()

	os.Exit(code)
}

func TestCreateToken(t *testing.T) {

	userId := uint(123)
	email := "test@example.com"

	tokenStr, err := CreateToken(userId, email)
	assert.NoError(t, err)
	assert.NotNil(t, tokenStr)

	// トークンのパース
	token, err := jwt.Parse(*tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	assert.NoError(t, err)
	assert.True(t, token.Valid)

	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, float64(userId), claims["sub"])
	assert.Equal(t, email, claims["email"])

	// exp(有効期限)が将来になっていることを確認
	exp := int64(claims["exp"].(float64))
	assert.Greater(t, exp, time.Now().Unix())
}

func TestCreateToken_NoSecret(t *testing.T) {
	// SECRET_KEY未設定時
	oldSecret := os.Getenv("SECRET_KEY")
	os.Setenv("SECRET_KEY", "")
	defer os.Setenv("SECRET_KEY", oldSecret)

	userId := uint(123)
	email := "test@example.com"

	tokenStr, err := CreateToken(userId, email)
	assert.Error(t, err)
	assert.Nil(t, tokenStr)
}
