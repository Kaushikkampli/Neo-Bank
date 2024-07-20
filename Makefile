postgres:
	docker run --name postgres12 --network neobank-network -p 3308:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=postgres -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root neo_bank

dropdb:
	docker exec -it postgres12 dropdb neo_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:postgres@localhost:3308/neo_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:postgres@localhost:3308/neo_bank?sslmode=disable" -verbose down

migrateup1:
	migrate -path db/migration -database "postgresql://root:postgres@localhost:3308/neo_bank?sslmode=disable" -verbose up 1

migratedown1:
	migrate -path db/migration -database "postgresql://root:postgres@localhost:3308/neo_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/kaushikkampli/neobank/db/sqlc Store

docker-create-network:
	docker network create neobank-network

docker-run:
	docker run --name neobank --network neobank-network -p 8080:8080 -e DB_SOURCE="postgresql://root:postgres@postgres12:5432/neo_bank?sslmode=disable" neobank:v2


.PHONY: postgres12 createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc server mock