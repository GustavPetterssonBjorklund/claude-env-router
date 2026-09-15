# Claude Env Router

Switch between Claude Code environments without repeatedly exporting variables or
editing shell configuration.

`cer` groups environment files, environment variables, and Claude arguments into
named profiles. Pick a profile interactively or launch one directly:

```sh
cer
cer personal
cer work -- --help
```

This is useful when you use different API providers, credentials, projects, or
Claude settings and want a simple way to keep them separate.

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

## Get started

Set your preferred editor and open the interactive profile picker:

```sh
export EDITOR="vim" # or nano, code, etc.
cer
```

If no config exists yet, press `n` to create a profile. `cer` creates the config
and an env file, opens the env file in your editor, and then starts Claude with
that profile.

You can also create the config yourself at `~/.config/cer/config.toml`:

```toml
binary = "claude"

[profiles.personal]
env_files = ["personal.env"]
args = []

[profiles.work]
env_files = ["work.env"]
args = []
```

Paths to env files are relative to the config file. For example,
`~/.config/cer/personal.env` might contain:

```dotenv
ANTHROPIC_API_KEY=your-api-key
```

Now choose a profile interactively with `cer`, or skip the picker by naming it:

```sh
cer personal
cer work
```

## Use DeepSeek with Claude Code

First, create an API key on the [DeepSeek Platform](https://platform.deepseek.com/).
Then add a profile to your `config.toml`:

```toml
[profiles.deepseek]
env_files = ["deepseek.env"]
args = []
```

Create `deepseek.env` next to the config file with DeepSeek's Claude Code
settings:

```dotenv
ANTHROPIC_BASE_URL=https://api.deepseek.com/anthropic
ANTHROPIC_MODEL=deepseek-flash[1m]
ANTHROPIC_DEFAULT_OPUS_MODEL=deepseek-flash[1m]
ANTHROPIC_DEFAULT_SONNET_MODEL=deepseek-flash[1m]
ANTHROPIC_DEFAULT_HAIKU_MODEL=deepseek-flash
CLAUDE_CODE_SUBAGENT_MODEL=deepseek-flash
CLAUDE_CODE_EFFORT_LEVEL=max
CLAUDE_CODE_AUTO_COMPACT_WINDOW=786432
```

Store your API key in `cer`'s encrypted vault and launch the profile:

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

## Keep secrets out of plain text

`cer` can store sensitive values in an encrypted vault. The encryption key is
kept in your operating system's keychain.

```sh
# Prompt securely for a value
cer secret set personal ANTHROPIC_API_KEY

# Show the secret names stored for a profile
cer secret list personal

# Remove a secret
cer secret unset personal ANTHROPIC_API_KEY
```

Secrets override values from env files and inline profile variables. If your
config already contains inline secrets, move them into the vault with:

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
