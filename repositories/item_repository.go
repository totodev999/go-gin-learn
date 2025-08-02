package repositories

import (
	"context"
	"errors"
	"flea-market/models"
	"flea-market/utils"
	"fmt"

	"gorm.io/gorm"
)

type ItemRepository struct {
	db *gorm.DB
}

// Create implements IItemRepository.
func (r *ItemRepository) Create(ctx context.Context, newItem models.Item) (*models.Item, error) {
	result := r.db.Create(&newItem)
	if result.Error != nil {
		return nil, utils.NewDBError("Create item failed", result.Error)
	}

	return &newItem, nil
}

// 論理削除となる。物理削除の場合は.Unscoped().Delete()にする
func (r *ItemRepository) Delete(ctx context.Context, itemId uint, userId uint) error {
	deleteItem, err := r.FindById(ctx, itemId, userId)
	if err != nil {
		return utils.NewNotFoundError(
			fmt.Sprintf("Data not found itemId:%d userId:%s", itemId, fmt.Sprint(userId)),
			err,
		)
	}

	result := r.db.Delete(&deleteItem)
	if result.Error != nil {
		return utils.NewDBError("Delete from item failed", result.Error)
	}
	return nil
}

// FindAll implements IItemRepository.
func (r *ItemRepository) FindAll(ctx context.Context) (*[]models.Item, error) {
	var items []models.Item
	result := r.db.Find(&items)
	if result.Error != nil {
		return nil, utils.NewDBError("DB Error", result.Error)
	}

	return &items, nil
}

// FindById implements IItemRepository.
func (r *ItemRepository) FindById(ctx context.Context, itemId uint, userId uint) (*models.Item, error) {
	var item models.Item
	result := r.db.First(&item, "id = ? AND user_id = ?", itemId, userId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, utils.NewNotFoundError("Not Found From DB", result.Error)
		}
		return nil, utils.NewDBError("DB Error", result.Error)
	}
	return &item, nil
}

// Update implements IItemRepository.
func (r *ItemRepository) Update(ctx context.Context, updateItem models.Item) (*models.Item, error) {
	result := r.db.Save(&updateItem)
	if result.Error != nil {
		return nil, utils.NewDBError("DB Error", result.Error)
	}
	return &updateItem, nil
}

func NewItemRepository(db *gorm.DB) *ItemRepository {
	return &ItemRepository{db: db}
}
