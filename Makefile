# Include .env file if it exists
ifneq (,$(wildcard .env))
    include .env
    export
endif

postgres:
	docker run --name url-service-db -p $(DB_PORT):5432 -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d postgres:17-alpine

createdb:
	docker exec -it url-service-db createdb --username=$(DB_USER) --owner=$(DB_USER) $(DB_NAME)

dropdb:
	docker exec -it url-service-db dropdb --username=$(DB_USER) $(DB_NAME)

migrateup:
	migrate -path internal/storage/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path internal/storage/migration -database "$(DB_URL)" -verbose down

build-pg:
	go run cmd/url-shortener/main.go -storage=postgres

build-mem:
	go run cmd/url-shortener/main.go -storage=memory

test-coverage:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown print-config build-pg build-mem test-coverage