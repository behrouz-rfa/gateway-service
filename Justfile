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

build what="all":
    @just _build_{{what}}

@_build_server:
    {{buildprops}} go build -o ./bin/app ./cmd/server/main.go

@_build_migrate:
    {{buildprops}} go build -o ./bin/migrate -tags migration ./cmd/migrate/main.go

@_build_all: _build_serve
