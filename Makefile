# including env variables from .envrc
include .envrc

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

confirm:
	@echo -n 'Are you sure? [y/n] ' && read ans && [ $${ans:-n} = y ]

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	@go run ./cmd/api -db-dsn=${GOBLOG_DSN}

.PHONY: build/api
build/api:
	@echo "Building binaries..."
	GOOS=windows GOARCH=amd64 go build -ldflags='-s' -o ./bin/api.exe ./cmd/api
	GOOS=linux   GOARCH=amd64 go build -ldflags='-s' -o ./bin/api     ./cmd/api

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${GOBLOG_DSN}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path=./migrations -database=${GOBLOG_DSN} up

## db/migrations/down: apply down migrations
.PHONY: db/migrations/down
db/migrations/down: confirm
	@echo 'Applying down migration by ${num}...'
	migrate -path=./migrations -database=${GOBLOG_DSN} down ${num}

.PHOY: tidy
tidy:
	@echo 'Tidying module dependencies...'
	go mod tidy
	@echo 'Verifying and vendoring module dependencies...'
	go mod verify
	go mod vendor
	@echo 'Formatting .go files...'
	go fmt ./...

