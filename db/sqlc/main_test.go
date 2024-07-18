package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/vlad19930514/webApp/internal/pkg/pg"
	"github.com/vlad19930514/webApp/util"
)

var testQueries *Queries

func TestMain(m *testing.M) {

	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("cannot load config", err)
	}
	pgConnection, err := pg.Dial(context.Background(), config.DBSource)
	if err != nil {
		log.Fatalf("Error creating connection pool:%v", err)
	}
	testQueries = New(pgConnection)
	os.Exit(m.Run())
}
