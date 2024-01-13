# ==================================================================================== #
# ENVIRONMENT VARIABLES
# ==================================================================================== #

PROJECT := $(shell basename $(shell pwd))

DOCKER := docker
DOCKER_COMPOSE := ${DOCKER} compose

.DEFAULT_GOAL := help

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...

## lint: run linter
.PHONY: lint
lint:
	golangci-lint run -E gofumpt ./...


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out


# ==================================================================================== #
# DOCKER
# ==================================================================================== #

## run/docker: run the cmd/api application in docker
.PHONY: run/docker/app
run/docker/app: mode=prod
run/docker/app: config_file=./configs/.env.$(mode)
run/docker/app:
	CONFIG_FILE=${config_file} \
		${DOCKER_COMPOSE} -p ${PROJECT} -f docker-compose.yml up --build -d 

## stop/docker: stop the cmd/api application in docker
.PHONY: stop/docker/app
stop/docker/app: mode=prod
stop/docker/app: config_file=./configs/.env.$(mode)
stop/docker/app:
	CONFIG_FILE=${config_file} \
		${DOCKER_COMPOSE} -p ${PROJECT} -f docker-compose.yml down

## run/docker/infra: run the db, inmemory db and s3 storage 
.PHONY: run/docker/infra
run/docker/infra:
	${DOCKER_COMPOSE} -p ${PROJECT} -f infra/docker-compose.yml up -d

## stop/docker/infra: stop the db, inmemory db and s3 storage 
.PHONY: stop/docker/infra
stop/docker/infra:
	${DOCKER_COMPOSE} -p ${PROJECT} -f infra/docker-compose.yml down

## log/docker/app: show logs the cmd/api application in docker
.PHONY: log/docker/app
log/docker/app:
	${DOCKER} logs --follow ${PROJECT}-app-1

## build/app: build the cmd/api-server
.PHONY: build/app
build/app: path=/tmp/${PROJECT}
build/app:
	go build -v -o ${path}/api-server ./cmd/api-server

## run/app: run the cmd/api-server
.PHONY: run/app
run/app: path=/tmp/${PROJECT}
run/app: build/app
	${path}/api-server


# ==================================================================================== #
# SQL MIGRATIONS
# ==================================================================================== #

## migrations/new name=$1: create a new database migration
.PHONY: migrations/new
migrations/new:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest create -seq -ext=.sql -dir=./assets/migrations ${name}

## migrations/up: apply all up database migrations
.PHONY: migrations/up
migrations/up:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="postgres://${DB_DSN}" up

## migrations/down: apply all down database migrations
.PHONY: migrations/down
migrations/down:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="postgres://${DB_DSN}" down

## migrations/drop: drop all database migrations
.PHONY: migrations/drop
migrations/drop:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="postgres://${DB_DSN}" drop

## migrations/goto version=$1: migrate to a specific version number
.PHONY: migrations/goto
migrations/goto:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="postgres://${DB_DSN}" goto ${version}

## migrations/force version=$1: force database migration
.PHONY: migrations/force
migrations/force:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="postgres://${DB_DSN}" force ${version}

## migrations/version: print the current in-use migration version
.PHONY: migrations/version
migrations/version:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="postgres://${DB_DSN}" version

