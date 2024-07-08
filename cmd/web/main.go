package main

import (
	"context"
	"log"

	"github.com/vlad19930514/webApp/api"
	"github.com/vlad19930514/webApp/util"

	PG "github.com/vlad19930514/webApp/db"
	db "github.com/vlad19930514/webApp/db/sqlc"
)

//mockgen -destination db/mock/store.go github.com/vlad19930514/webApp/db/sqlc Store
/* func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
} */

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	pgConnection, err := PG.NewPG(context.Background(), config.DBSource)
	if err != nil {
		log.Fatalf("Error creating connection pool: %v", err)
	}

	store := db.NewStore(pgConnection.Db)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatalf("Cannot start server: %v", err)
	}
	/* 	r := gin.Default()
	   	r.GET("/ping", func(c *gin.Context) {
	   		c.JSON(http.StatusOK, gin.H{
	   			"message": "pong",
	   		})
	   	})
	   	r.Run(os.Getenv("SERVER_PORT")) */
	defer pgConnection.Close()

	/*
		 	mux := http.NewServeMux()
			mux.HandleFunc("/", helloWorldPage)

			mux.HandleFunc("/user", getUser)

			log.Println("Starting server on :8080")
			err = http.ListenAndServe(":8080", mux)
			if err != nil {
				log.Fatalf("HTTP server error: %v", err)
			}
	*/
}
