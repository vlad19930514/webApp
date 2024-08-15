package main

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/vlad19930514/webApp/internal/app/repository/pgrepo"
	"github.com/vlad19930514/webApp/internal/app/services"
	"github.com/vlad19930514/webApp/internal/app/transport/httpserver"
	pg "github.com/vlad19930514/webApp/internal/pkg/pg"
	"github.com/vlad19930514/webApp/util"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Caller().Logger()
	err := runApplication()
	if err != nil {
		log.Fatal().Err(err)
	}
}

func runApplication() error {
	config, err := util.LoadConfig(".")
	if err != nil {
		return fmt.Errorf("cannot load config: %w", err)
	}

	pgDB, err := pg.Dial(config.DSN)

	if err != nil {
		return fmt.Errorf("error creating connection pool: %w", err)
	}

	// create repositories
	userRepo, err := pgrepo.NewUserRepo(pgDB)
	if err != nil {
		return fmt.Errorf("failed to create userRepo: %w", err)
	}
	// create services
	userService := services.NewUserService(userRepo)

	server := httpserver.NewHttpServer(userService)

	err = server.Start(config.ServerAddress)
	if err != nil {
		return fmt.Errorf("cannot start server: %w", err)
	}
	return nil
}
