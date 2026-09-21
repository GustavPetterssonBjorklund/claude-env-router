# Claude Env Router

`cer` lets you keep separate Claude Code setups—for example, personal, work, and
another API provider—and start the one you want with a single command. Store API
keys in encrypted secure storage instead of putting them in shell configuration,
`.env` files, or a repository.

Pick a profile interactively or launch one directly:

```sh
cer
cer personal
cer work -- --help
```

Each profile can also include non-secret provider settings, Claude arguments,
and optional environment files.

## Install

You will need [Claude Code](https://docs.anthropic.com/en/docs/claude-code) and Go
1.25.10 or newer. Install `cer` with:

```sh
go install github.com/GustavPetterssonBjorklund/claude-env-router/cmd/cer@latest
```

Make sure your Go bin directory is on your `PATH`, then check the installation:

```sh
cer --help
```

## Get started: store your API key securely

Create the configuration directory if needed:

```sh
mkdir -p ~/.config/cer
```

Then create `~/.config/cer/config.toml` with a profile for each setup you use.
The profile itself contains no secret values:

```toml
binary = "claude"

[profiles.personal]
args = []

[profiles.work]
args = []
```

Save the API key for each profile. `cer` prompts without echoing the value:

```sh
cer secret set personal ANTHROPIC_API_KEY
cer secret set work ANTHROPIC_API_KEY
```

The values are encrypted in `~/.config/cer/secrets.vault`; the encryption key is
stored in your operating system's credential store. The vault and its key are
both required to read a secret.

Start Claude with the profile you need:

```sh
cer personal
cer work
```

Run `cer` with no profile to choose one interactively. The picker can create and
edit profiles too; those actions use your configured `$EDITOR`, but an editor is
not needed for the secure-storage workflow above.

## Use DeepSeek with Claude Code

First, create an API key on the [DeepSeek Platform](https://platform.deepseek.com/).
Then add a profile to your `config.toml`:

```toml
[profiles.deepseek]
env_files = ["deepseek.env"]
args = []
```

Create `deepseek.env` next to the config file with DeepSeek's non-secret Claude
Code settings:

```dotenv
ANTHROPIC_BASE_URL=https://api.deepseek.com/anthropic
ANTHROPIC_MODEL=deepseek-flash
ANTHROPIC_DEFAULT_OPUS_MODEL=deepseek-flash
ANTHROPIC_DEFAULT_SONNET_MODEL=deepseek-flash
ANTHROPIC_DEFAULT_HAIKU_MODEL=deepseek-flash
CLAUDE_CODE_SUBAGENT_MODEL=deepseek-flash
CLAUDE_CODE_EFFORT_LEVEL=max
CLAUDE_CODE_AUTO_COMPACT_WINDOW=786432
```

Store the API key in secure storage, then launch the profile:

```sh
cer secret set deepseek ANTHROPIC_AUTH_TOKEN
cer deepseek
```

See DeepSeek's [official Claude Code integration
guide](https://api-docs.deepseek.com/quick_start/agent_integrations/claude_code/)
for the latest provider and model settings.

## Common commands

```sh
# Open the interactive profile picker
cer

# Start Claude with a profile
cer personal

# Pass extra arguments to Claude
cer work -- --help

# Use a config in a different location
cer --config ./examples/config.toml personal
```

In the profile picker, press `e` to edit the highlighted profile's first env
file. The picker reopens when your editor closes.

## Manage secrets

Use the vault for tokens, API keys, and other sensitive environment variables:

```sh
# Prompt securely for a value
cer secret set personal ANTHROPIC_API_KEY

# Show the secret names stored for a profile
cer secret list personal

# Remove a secret
cer secret unset personal ANTHROPIC_API_KEY
```

Vault secrets override values from environment files and inline profile
variables. If your config already contains inline secrets, move them into the
vault with:

```sh
cer secret migrate
```

## How configuration is found

`cer` uses the first config it finds:

1. A path passed with `--config`
2. The `CER_CONFIG` environment variable
3. A `cer.toml` in the current directory or one of its parents
4. `~/.config/cer/config.toml`

See the [configuration reference](docs/config.md) and [usage guide](docs/usage.md)
for more detail.

## Development

Clone the repository and run the complete test suite with:

```sh
go test ./...
```

Pull requests and pushes run formatting checks, `go vet`, unit tests, and the CLI
integration test through [GitHub Actions](.github/workflows/ci.yml).

## License

[MIT](LICENSE)
