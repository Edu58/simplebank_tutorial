postgres:
	docker run --name postgres12 -p 5432:5433 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:16.2-alpine

createdb:
	 docker exec -it postgres createdb --username=postgres --owner=postgres simple_bank

dropdb:
	 docker exec -it postgres dropdb --username=postgres simple_bank

migrateup:
	migrate -path db/migration -database "postgres://postgres:postgres@localhost:5433/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgres://postgres:postgres@localhost:5433/simple_bank?sslmode=disable" -verbose down

migrate-rollback:
	migrate -path db/migration -database "postgres://postgres:postgres@localhost:5433/simple_bank?sslmode=disable" -verbose down 1


sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen --package mockdb --destination db/mock/store.go github.com/techschool/simplebank/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock
