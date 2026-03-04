.PHONY: build run clean

build:
	cd backend && go build -o bin/api ./cmd/api

run:
	cd backend && go run ./cmd/api

clean:
	rm -rf backend/bin/

DB_URL ?= postgres://granttool:localdevpassword@localhost:5432/granttool?sslmode=disable

migrate-up:
	goose -dir backend/migrations postgres "$(DB_URL)" up

migrate-down:
	goose -dir backend/migrations postgres "$(DB_URL)" down

migrate-create:
	goose -dir backend/migrations create $(name) sql

generate:
	cd backend && sqlc generate