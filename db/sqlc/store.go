package db

import "github.com/jackc/pgx/v5/pgxpool"

type Store interface { //TODO изменить на userStore
	Querier
}
type SQLStore struct {
	*Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}
