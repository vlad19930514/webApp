package PG

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgres struct {
	Db *pgxpool.Pool
}

var (
	pgInstance *postgres
	pgOnce     sync.Once
)

func NewPG(ctx context.Context, connString string) (*postgres, error) {

	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, connString)
		if err != nil {
			err = fmt.Errorf("unable to create connection pool: %w", err)
			log.Fatalf("%v", err)
		}

		pgInstance = &postgres{db}
	})
	if pgInstance == nil {
		return nil, errors.New("connection pool not initialized")
	}
	return pgInstance, nil
}
func (pg *postgres) Ping(ctx context.Context) error {
	if pg.Db == nil {
		return errors.New("connection pool is nil")
	}
	return pg.Db.Ping(ctx)
}
func (pg *postgres) Close() {
	if pg.Db != nil {
		pg.Db.Close()
	}
}
