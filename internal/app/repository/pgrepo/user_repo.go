package pgrepo

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/vlad19930514/webApp/internal/app/domain"
	"github.com/vlad19930514/webApp/internal/pkg/pg"
)

type UserRepo struct {
	db *pg.DB
}

func NewUserRepo(db *pg.DB) (*UserRepo, error) {

	// Добавление расширения uuid-ossp
	err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	if err != nil {
		return nil, fmt.Errorf("failed to create uuid-ossp extension: %w", err)
	}

	// Миграция схемы для структуры User
	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate user table: %w", err)
	}
	return &UserRepo{
		db: db,
	}, nil
}

func (r UserRepo) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	result := r.db.WithContext(ctx).Create(&user)
	if result.Error != nil {
		return domain.User{}, fmt.Errorf("failed to create domain user: %w", result.Error)
	}
	return user, nil
}
func (r UserRepo) GetUser(ctx context.Context, id uuid.UUID) (domain.User, error) {
	dbUser := domain.User{
		Id: id,
	}
	result := r.db.WithContext(ctx).Take(&dbUser)
	if result.Error != nil {
		return domain.User{}, fmt.Errorf("failed to get a user: %w", result.Error)
	}
	return dbUser, nil

}
func (r UserRepo) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {

	result := r.db.WithContext(ctx).Save(&user)
	if result.Error != nil {
		return domain.User{}, fmt.Errorf("failed to get a user: %w", result.Error)
	}
	return user, nil

}
