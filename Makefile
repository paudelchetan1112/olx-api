.PHONY: build run
build:
	@go build -buildvcs=false -o bin/api/ ./cmd/api
run:build
	@./bin/api/api
migrate-up:
	@go run ./cmd/migrate up
migrate-down:
	@go run ./cmd/migrate down