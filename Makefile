DB_URL=postgres://root:secret@localhost:5432/web_app?sslmode=disable
network:
	docker network create web-network

postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root web_app

nodemon:
	nodemon --exec go run ./cmd/web/main.go --signal SIGTERM	

migratedown:
	migrate -path db/migration -database postgres://root:secret@localhost:5432/web_app?sslmode=disable -verbose down
 
migrateup:
	migrate -path db/migration -database postgres://root:secret@localhost:5432/web_app?sslmode=disable -verbose up

sqlc:
	sqlc generate

dockerdown:
	docker-compose down -v

dockerprune:
	docker system prune -a --volumes

test:
	go test -v -cover ./...

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/vlad19930514/webApp/db/sqlc Store

.PHONY: postgres createdb nodemon migratedown migrateup sqlc test mock