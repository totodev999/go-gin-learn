package services

import (
	"context"
	"errors"
	"flea-market/models"
	"flea-market/utils"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type IAuthRepository interface {
	CreateUser(ctx context.Context, user models.User) error
	FindUser(ctx context.Context, email string) (*models.User, error)
}

type AuthService struct {
	repository IAuthRepository
}

func (s *AuthService) Signup(ctx context.Context, email string, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return utils.NewUnknownError("bcrypt.GenerateFromPassword failed", err)
	}
	user := models.User{Email: email, Password: string(hashed)}
	return s.repository.CreateUser(ctx, user)

}

func (s *AuthService) Login(ctx context.Context, email string, password string) (*string, error) {
	user, err := s.repository.FindUser(ctx, email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, utils.NewUnauthorized("Invalid email or password", err)
	}

	token, err := CreateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *AuthService) GetUserFromToken(token string) (*models.User, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, utils.NewUnauthorized("JWT is invalid", fmt.Errorf("unexpected signing method %v", t.Header["alg"]))
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		return nil, utils.NewUnauthorized("failed to parse token", err)
	}

	var user *models.User
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		exp, ok := claims["exp"].(float64)
		if !ok {
			return nil, utils.NewUnauthorized("token is invalid", errors.New("can't get exp from the token"))
		}

		if float64(time.Now().Unix()) > exp {
			return nil, utils.NewUnauthorized("token is expired", jwt.ErrTokenExpired)
		}

		// this function is called by middleware, therefore context is not passed, so create context here.
		ctx := context.TODO()
		user, err = s.repository.FindUser(ctx, claims["email"].(string))
		if err != nil {
			return nil, err
		}
	}

	return user, nil

}

func NewAuthService(repository IAuthRepository) *AuthService {
	return &AuthService{repository: repository}
}

func CreateToken(userId uint, email string) (*string, error) {
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		return nil, utils.NewUnknownError("Internal Error", errors.New("SECRET_KEY is not set"))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   userId,
		"email": email,
		"exp":   time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, utils.NewUnknownError("Creating JWT failed", err)
	}

	return &tokenString, nil
}
