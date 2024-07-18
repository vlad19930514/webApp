package main

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/vlad19930514/webApp/api"
	pg "github.com/vlad19930514/webApp/internal/pkg/pg"
	"github.com/vlad19930514/webApp/util"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	db "github.com/vlad19930514/webApp/db/sqlc"
)

/* func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
} */

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().AnErr("cannot load config", err)
	}

	pgConnection, err := pg.Dial(context.Background(), config.DBSource)

	if err != nil {
		log.Fatal().AnErr("Error creating connection pool: %v", err)
	}
	defer pgConnection.Close()

	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(pgConnection.Pool) //TODO передавать pgconnection
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal().AnErr("Cannot start server: %v", err) //TODO вынести ошибки в run
	}

}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrate up")
	}

	log.Info().Msg("db migrated successfully")
}
