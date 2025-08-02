// test_helpers.go（ファイル分けてもOK, main_test.go内に直書きもOK）

package mocks

import (
	"flea-market/models"

	"golang.org/x/net/context"
)

type MockItemRepository struct {
	FindAllFunc  func(ctx context.Context) (*[]models.Item, error)
	FindByIdFunc func(ctx context.Context, itemId uint, userId uint) (*models.Item, error)
	CreateFunc   func(ctx context.Context, newItem models.Item) (*models.Item, error)
	UpdateFunc   func(ctx context.Context, updateItem models.Item) (*models.Item, error)
	DeleteFunc   func(ctx context.Context, itemId uint, userId uint) error
}

func (m *MockItemRepository) FindAll(ctx context.Context) (*[]models.Item, error) {
	return m.FindAllFunc(ctx)
}
func (m *MockItemRepository) FindById(ctx context.Context, itemId uint, userId uint) (*models.Item, error) {
	return m.FindByIdFunc(ctx, itemId, userId)
}
func (m *MockItemRepository) Create(ctx context.Context, newItem models.Item) (*models.Item, error) {
	return m.CreateFunc(ctx, newItem)
}
func (m *MockItemRepository) Update(ctx context.Context, updateItem models.Item) (*models.Item, error) {
	return m.UpdateFunc(ctx, updateItem)
}
func (m *MockItemRepository) Delete(ctx context.Context, itemId uint, userId uint) error {
	return m.DeleteFunc(ctx, itemId, userId)
}
