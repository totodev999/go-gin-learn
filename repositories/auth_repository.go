package repositories

import (
	"errors"
	"free-market/models"

	"gorm.io/gorm"
)

type IAuthRepository interface {
	CreateUser(user models.User) error
	FindUser(email string) (*models.User, error)
}

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) IAuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(user models.User) error {
	result := r.db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *AuthRepository) FindUser(email string) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, "email = ?", email)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}

	return &user, nil
}
