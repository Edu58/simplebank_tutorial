postgres:
	docker run --name postgres17 -p 5433:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres

createdb:
	 docker exec -it postgres17 createdb --username=postgres --owner=postgres simple_bank

dropdb:
	 docker exec -it postgres17 dropdb --username=postgres simple_bank

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

proto:
	rm -f pb/*.go
	rm -f docs/swagger/*swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=./docs/swagger --openapiv2_opt=allow_merge=true,merge_file_name=SimpleBank \
    proto/*.proto

evans :
	evans -r -p 9090 repl

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock migrate-rollback proto evans
