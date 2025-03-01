include .env
export $(shell sed 's/=.*//' .env)

run: build
	@./bin/timeit

build: fmt
	@go build -o bin/timeit ./cmd/app 

fmt:
	@go fmt ./...

test:
	go test -v ./...

migrate-up:
	goose -dir ./migrations postgres postgresql://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DB) up

migrate-down:
	goose -dir ./migrations postgres postgresql://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DB) down