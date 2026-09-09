# claude-env-router

`cer` runs Claude with environment variables from a named profile.

```sh
cer
cer --config examples/config.toml personal
cer work -- --help
```

See [docs/config.md](docs/config.md) and [docs/usage.md](docs/usage.md).

## Development

Run the full test suite, including the CLI integration test:

```sh
go test ./...
```

Pull requests and pushes run formatting, `go vet`, unit tests, and integration
tests through [GitHub Actions](.github/workflows/ci.yml).
