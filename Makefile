
.DEFAULT_GOAL := all

PROJECT=gotube


.PHONY: all
all: ci run


.PHONY: test
test:
	go test -v -race -timeout=5m ./...


.PHONY: lint
lint:
	golangci-lint run ./...


.PHONY: ci
ci: test lint


.PHONY: run
run:
	docker compose -p ${PROJECT} -f ./build/docker-compose.yaml up -d --build


.PHONY: run-web
run-web:
	docker compose -p ${PROJECT} -f ./build/docker-compose.yaml up web -d --build


.PHONY: run-app
run-app:
	docker compose -p ${PROJECT} -f ./build/docker-compose.yaml up app -d --build


.PHONY: run-infra
run-infra:
	docker compose -p ${PROJECT} -f ./build/docker-compose.yaml up postgres mongo s3 redis -d --build


.PHONY: stop
stop:
	docker compose -p ${PROJECT} -f ./build/docker-compose.yaml down


# TODO: Add script for migrations

.PHONY: migrate
migrate: MIGRATE_ACTION=up
migrate: MIGRATE_PATH=./assets/migrations/postgres
migrate: MIGRATE_DB="postgres://admin:123456789@localhost:5432/gotubedb?sslmode=disable"
migrate:
	migrate -path ${MIGRATE_PATH} -database ${MIGRATE_DB} ${MIGRATE_ACTION}

