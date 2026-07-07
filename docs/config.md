# Config

`cer` reads config from the first available source:

1. `--config path`
2. `CER_CONFIG`
3. `cer.toml` in the current directory or a parent directory
4. `~/.config/cer/config.toml`

Supported config shape:

```toml
binary = "claude"

[profiles.personal]
env_files = ["personal.env"]
args = []

[profiles.personal.env]
ANTHROPIC_API_KEY = "sk-ant-example"
```

Environment files use `KEY=value` lines. Later files override earlier values, and inline profile env values override env files.
