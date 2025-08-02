package repositories

import (
	"context"
	"errors"
	"flea-market/models"
	"flea-market/utils"
	"fmt"

	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(ctx context.Context, user models.User) error {
	result := r.db.Create(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return utils.NewDuplicateKeyError(fmt.Sprintf("Duplicated key %s", user.Email), result.Error)
		}
		return utils.NewDBError("create user failed", result.Error)
	}
	return nil
}

func (r *AuthRepository) FindUser(ctx context.Context, email string) (*models.User, error) {

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
