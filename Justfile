# Display the list of available commands
help:
    @just --list --list-heading $'\033[33mCommands\033[0m:\n'

isCI := env_var_or_default("CI","false")
buildprops := if isCI == "true" { "env GOOS=linux GOARCH=arm64"  } else { "" }



@_lint_requirements:
    which golangci-lint 2> /dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1


# Run linting check
lint: _lint_requirements
    golangci-lint run --fix

fmt:
	gci write .
	gofumpt -l -w .

# Generate swagger documentation
swag:
   swag init -g cmd/server/main.go  --parseInternal true

servers := "main & gbp"

test mode="all":
    @just _test_{{mode}}

@_test_all: _test_unit _test_it _test_e2e

@_test_unit:
    go test -v  -count=1 ./... -tags=unit

@_test_it:
   docker compose up -d mongo mongo-express
   SERVER_PORT=8050 \
   SERVER_HOST=localhost \
   GIN_MODE=debug \
   DB_HOST=localhost \
   DB_PORT=27017 \
   DB_NAME=oukoud \
   DB_USER=root \
   DB_PASS=root \
   DB_SSL=false \
   DB_CLUSTERED=false \
   ENV=dev \
   LOG_LEVEL=debug \
   JWT_KEY="bd32c589-f9cd-479b-9500-a0392b13252b" \
   FIREBASE_PROJECT_ID=oukoud-dev \
   FIREBASE_FILE="firebase.json" \
   bash -c 'go test -v -count=1 -tags=integration ./tests/it/...'
   docker compose down

@_test_e2e:
   docker compose up -d mongodbtest -d -p 27018:27017  mongo-express
   SERVER_PORT=8051 \
   SERVER_HOST=localhost \
   GIN_MODE=debug \
   DB_HOST=localhost \
   DB_PORT=27017 \
   DB_NAME=oukoud \
   DB_USER=root \
   DB_PASS=root \
   DB_SSL=false \
   DB_CLUSTERED=false \
   ENV=dev \
   LOG_LEVEL=debug \
   FIREBASE_PROJECT_ID=oukoud-dev \
   FIREBASE_FILE="firebase.json" \
   JWT_KEY="bd32c589-f9cd-479b-9500-a0392b13252b" \
   bash -c 'go test -v -count=1 -tags=integration ./tests/e2e/...'
   docker compose down

build what="all":
    @just _build_{{what}}

@_build_server:
    {{buildprops}} go build -o ./bin/app ./cmd/server/main.go

@_build_migrate:
    {{buildprops}} go build -o ./bin/migrate -tags migration ./cmd/migrate/main.go

@_build_all: _build_server _build_migrate


# Run migrations
migrate:
    cp .env ./bin
    go run -tags migration ./cmd/migrate/main.go
