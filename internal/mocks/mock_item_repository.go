// test_helpers.go（ファイル分けてもOK, main_test.go内に直書きもOK）

package mocks

import (
	"free-market/models"
)

type MockItemRepository struct {
	FindAllFunc  func() (*[]models.Item, error)
	FindByIdFunc func(itemId uint, userId uint) (*models.Item, error)
	CreateFunc   func(newItem models.Item) (*models.Item, error)
	UpdateFunc   func(updateItem models.Item) (*models.Item, error)
	DeleteFunc   func(itemId uint, userId uint) error
}

func (m *MockItemRepository) FindAll() (*[]models.Item, error) {
	return m.FindAllFunc()
}
func (m *MockItemRepository) FindById(itemId uint, userId uint) (*models.Item, error) {
	return m.FindByIdFunc(itemId, userId)
}
func (m *MockItemRepository) Create(newItem models.Item) (*models.Item, error) {
	return m.CreateFunc(newItem)
}
func (m *MockItemRepository) Update(updateItem models.Item) (*models.Item, error) {
	return m.UpdateFunc(updateItem)
}
func (m *MockItemRepository) Delete(itemId uint, userId uint) error {
	return m.DeleteFunc(itemId, userId)
}
