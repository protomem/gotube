
.DEFAULT_GOAL := all

PROJECT_NAME := gotube


.PHONY: all
all: ci run-stage


.PHONY: test
test:
	go test -v -race -timeout 5m ./...


.PHONY: lint
lint:
	golangci-lint run ./...


.PHONY: ci
ci: test lint


.PHONY: run-infra-local
run-infra-local:
	docker compose -p ${PROJECT_NAME}-local -f ./build/local/docker-compose.yaml up -d --build


.PHONY: stop-infra-local
stop-infra-local:
	docker compose -p ${PROJECT_NAME}-local -f ./build/local/docker-compose.yaml down


.PHONY: run-local
run-local:
	go run ./cmd/${PROJECT_NAME} --conf ./configs/local/app.yaml


.PHONY: run-web-local
run-web-local:
	cd ./web && npm run dev


.PHONY: run-stage
run-stage:
	docker compose -p ${PROJECT_NAME}-stage -f ./build/stage/docker-compose.yaml up -d --build


.PHONY: stop-stage
stop-stage:
	docker compose -p ${PROJECT_NAME}-stage -f ./build/stage/docker-compose.yaml down


.PHONY: migrate-up
migrate: MIGRATE_ACTION=up
migrate: MIGRATE_PATH=./assets/migrations
migrate: MIGRATE_DB="postgres://admin:123456789@localhost:5432/gotubedb?sslmode=disable"
migrate:
	migrate -path ${MIGRATE_PATH} -database ${MIGRATE_DB} ${MIGRATE_ACTION}

