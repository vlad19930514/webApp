package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/vlad19930514/webApp/internal/app/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) (domain.User, error)
}
