package services

import (
	"context"
	"flea-market/dto"
	"flea-market/models"
)

type IItemRepository interface {
	FindAll(ctx context.Context) (*[]models.Item, error)
	FindById(ctx context.Context, itemId uint, userId uint) (*models.Item, error)
	Create(ctx context.Context, newItem models.Item) (*models.Item, error)
	Update(ctx context.Context, updateItem models.Item) (*models.Item, error)
	Delete(ctx context.Context, itemId uint, userId uint) error
}

type ItemService struct {
	repository IItemRepository
}

func NewItemService(repository IItemRepository) *ItemService {
	return &ItemService{repository: repository}
}

func (s *ItemService) FindAll(ctx context.Context) (*[]models.Item, error) {
	return s.repository.FindAll(ctx)
}

func (s *ItemService) FindById(ctx context.Context, itemId uint, userId uint) (*models.Item, error) {
	return s.repository.FindById(ctx, itemId, userId)
}

func (s *ItemService) Create(ctx context.Context, createItemInput dto.CreateItemInput, userId uint) (*models.Item, error) {
	newItem := models.Item{
		Name:        createItemInput.Name,
		Price:       createItemInput.Price,
		Description: createItemInput.Description,
		SoldOut:     false,
		UserID:      userId,
	}

	return s.repository.Create(ctx, newItem)

}

func (s *ItemService) Update(ctx context.Context, itemId uint, updateItemInput dto.UpdateItemInput, userId uint) (*models.Item, error) {
	targetItem, err := s.repository.FindById(ctx, itemId, userId)
	if err != nil {
		return nil, err
	}

	if updateItemInput.Name != nil {
		targetItem.Name = *updateItemInput.Name
	}
	if updateItemInput.Price != nil {
		targetItem.Price = *updateItemInput.Price
	}
	if updateItemInput.Description != nil {
		targetItem.Description = *updateItemInput.Description
	}
	if updateItemInput.SoldOut != nil {
		targetItem.SoldOut = *updateItemInput.SoldOut
	}

	return s.repository.Update(ctx, *targetItem)

}

func (s *ItemService) Delete(ctx context.Context, itemId uint, userId uint) error {
	return s.repository.Delete(ctx, itemId, userId)
}
