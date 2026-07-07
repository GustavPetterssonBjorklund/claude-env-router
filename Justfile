set dotenv-load := false

gocache := env_var_or_default("GOCACHE", "/tmp/cer-go-build")
cerconfig := env_var_or_default("CER_CONFIG", "examples/config.toml")

default:
    @just --list

test:
    GOCACHE={{gocache}} go test ./...

run *args:
    @if [ -z "{{args}}" ]; then \
        GOCACHE={{gocache}} go run ./cmd/cer --help; \
    else \
        set -- {{args}}; \
        if [ "${1:-}" = "--" ]; then shift; fi; \
        CER_CONFIG={{cerconfig}} GOCACHE={{gocache}} go run ./cmd/cer "$@"; \
    fi

help:
    GOCACHE={{gocache}} go run ./cmd/cer --help
