package pg

import (
	"context"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func Dial(ctx context.Context, dsn string) (*DB, error) {

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		err = fmt.Errorf("unable to create connection pool: %w", err)
		return nil, err
	}

	return &DB{db}, nil
}
