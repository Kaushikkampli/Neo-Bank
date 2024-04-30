postgres:
	docker run --name postgres12 -p 3308:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=postgres -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root neo_bank

dropdb:
	docker exec -it postgres12 dropdb neo_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:postgres@localhost:3308/neo_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:postgres@localhost:3308/neo_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/kaushikkampli/neobank/db/sqlc Store

.PHONY: postgres12 createdb dropdb migrateup migratedown sqlc server mock