package services

import (
	"context"
	"flea-market/repositories"
)

type IAPICallRepository interface {
	GetAllPosts(ctx context.Context) (*[]repositories.Post, error)
}

type APICallService struct {
	repository IAPICallRepository
}

func NewAPICallService(repository IAPICallRepository) *APICallService {
	return &APICallService{repository: repository}
}

// In many example, context is passed as it is, not convert into pointer.
func (s *APICallService) GetAllPosts(ctx context.Context) (*[]repositories.Post, error) {
	return s.repository.GetAllPosts(ctx)
}
