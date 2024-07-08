package db

import (
	"context"
	"log"
	"os"
	"testing"

	PG "github.com/vlad19930514/webApp/db"
	"github.com/vlad19930514/webApp/util"
)

var testQueries *Queries

func TestMain(m *testing.M) {

	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("cannot load config", err)
	}
	pgConnection, err := PG.NewPG(context.Background(), config.DBSource)
	if err != nil {
		log.Fatalf("Error creating connection pool:%v", err)
	}
	defer pgConnection.Close()
	testQueries = New(pgConnection.Db)
	os.Exit(m.Run())
}
