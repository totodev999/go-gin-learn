package repositories

import (
	"errors"
	"fmt"
	"free-market/models"
	"free-market/utils"

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
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return utils.NewDuplicateKeyError(fmt.Sprintf("Duplicated key %s", user.Email), result.Error)
		}
		return utils.NewDBError("create user failed", result.Error)
	}
	return nil
}

func (r *AuthRepository) FindUser(email string) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, "email = ?", email)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, utils.NewNotFoundError(fmt.Sprintf("user %v not found", email), result.Error)
		}
		return nil, utils.NewDBError("Find user failed", result.Error)
	}
	return &user, nil
}
