set dotenv-load := false

gocache := env_var_or_default("GOCACHE", "/tmp/cer-go-build")
cerconfig := env_var_or_default("CER_CONFIG", ".run/config.toml")

ensure-cerconfig-dir:
    @case "{{cerconfig}}" in \
        .run/*) mkdir -p -- "$$(dirname -- "{{cerconfig}}")" ;; \
    esac
default:
    @just --list

test: ensure-cerconfig-dir
    GOCACHE={{gocache}} go test ./...

run *args: ensure-cerconfig-dir
    @if [ -z "{{args}}" ]; then \
        CER_CONFIG={{cerconfig}} GOCACHE={{gocache}} go run ./cmd/cer; \
    else \
        set -- {{args}}; \
        if [ "${1:-}" = "--" ]; then shift; fi; \
        CER_CONFIG={{cerconfig}} GOCACHE={{gocache}} go run ./cmd/cer "$@"; \
    fi

help: ensure-cerconfig-dir
    CER_CONFIG={{cerconfig}} GOCACHE={{gocache}} go run ./cmd/cer --help
