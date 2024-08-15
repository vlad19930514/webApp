package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/vlad19930514/webApp/internal/app/domain"
)

// UserService is a user service
type UserService struct {
	repo UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo UserRepository) UserService {
	return UserService{
		repo: repo,
	}
}

func (s UserService) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	return s.repo.CreateUser(ctx, user)
}
func (s UserService) GetUser(ctx context.Context, id uuid.UUID) (domain.User, error) {
	return s.repo.GetUser(ctx, id)
}
func (s UserService) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {
	return s.repo.UpdateUser(ctx, user)
}
