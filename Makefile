postgres:
	docker run --name some-postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres

createdb:
	docker exec -it some-postgres createdb --username=root --owner=root web_app

nodemon:
	nodemon --exec go run ./cmd/web/main.go --signal SIGTERM	

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/web_app?sslmode=disable" -verbose down
 
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/web_app?sslmode=disable" -verbose up

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/vlad19930514/webApp/db/sqlc Store

.PHONY: postgres createdb nodemon migratedown migrateup sqlc test mock