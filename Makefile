.PHONY: build run clean

build:
	cd backend && go build -o bin/api ./cmd/api

run:
	cd backend && go run ./cmd/api

clean:
	rm -rf backend/bin/