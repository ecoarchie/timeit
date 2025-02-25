include .env
export $(shell sed 's/=.*//' .env)

run:
	go run cmd/app/main.go

test:
	go test -v ./...

migrate-up:
	goose -dir ./migrations postgres postgresql://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DB) up

migrate-down:
	goose -dir ./migrations postgres postgresql://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DB) down