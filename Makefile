PKGS ?= $(shell go list ./...)
.PHONY: all services clean

migrate_up:
	migrate -database ${POSTGRESQL_URL} -path db/migrations up

migrate_down:
	migrate -database ${POSTGRESQL_URL} -path db/migrations down

ensure-dependencies:
	go mod tidy

services:
	docker-compose up -d psql

create_migration:
	migrate create -ext sql -dir ./db/migrations $(NAME)

