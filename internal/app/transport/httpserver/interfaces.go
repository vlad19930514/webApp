package httpserver

import (
	"context"
	"github.com/google/uuid"
	"github.com/vlad19930514/webApp/internal/app/domain"
)

//go:generate mockery --output=./mocks --filename=mock_userService.go --name=IUserService  --outpkg=mocks --structname=MockIUserService
type IUserService interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)

	GetUser(ctx context.Context, id uuid.UUID) (domain.User, error)

	UpdateUser(ctx context.Context, user domain.User) (domain.User, error)
}
