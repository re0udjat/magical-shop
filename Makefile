# Default target
.DEFAULT_GOAL := help

# Colors
YELLOW := \033[33m
GREEN  := \033[32m
RESET  := \033[0m

# ==================================================================================== #
# HELPERS
# ==================================================================================== 

.PHONY: help
help:
	@echo ""
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Available targets:"
	@echo ""
	@sed -n 's/^##//p' ${MAKEFILE_LIST} \
	| sort \
	| (command -v column >/dev/null \
		&& column -t -s ':' \
		|| awk -F ':' '{ printf "%-20s %s\n", $$1, $$2 }') \
	| sed -e 's/^/ /'

# Create the new confirm target
.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N]' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

.PHONY: run/api
run/api:
	go run ./cmd/api -db-dsn=${DB_DSN} -jwt-secret=${JWT_SECRET}

.PHONY: db/psql
db/psql:
	psql ${DB_DSN}

.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

.PHONY: db/migrations/up
db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${DB_DSN} up

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

.PHONY: tidy
tidy:
	@echo 'Formatting .go files...'
	go fmt ./...
	@echo 'Tidying module dependencies...'
	go mod tidy
	@echo 'Verifying and vendoring module dependencies...'
	go mod verify
	go mod vendor

.PHONY: audit
audit:
	@echo 'Checking module dependencies'
	go mod tidy -diff
	go mod verify
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

# ==================================================================================== #
# BUILD
# ==================================================================================== #

.PHONY: build/api
build/api:
	@echo 'Building API...'
	go build -ldflags='-s' -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/api ./cmd/api
