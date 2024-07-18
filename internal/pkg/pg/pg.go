package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	*pgxpool.Pool
}

func Dial(ctx context.Context, connString string) (*DB, error) {
	db, err := pgxpool.New(ctx, connString)
	if err != nil {
		err = fmt.Errorf("unable to create connection pool: %w", err)
		return nil, err
	}
	return &DB{db}, nil
}

func (pg *DB) Close() { //TODO k нужно ли закрытие?
	if pg.Pool != nil {
		pg.Pool.Close()
	}
}
