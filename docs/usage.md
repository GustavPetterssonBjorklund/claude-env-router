# Usage

Run `cer` without arguments to choose or create an environment profile interactively.
New profiles create an env file and open it in `$EDITOR`:

```sh
cer
```

Run Claude directly with a named environment profile:

```sh
cer personal
```

Pass arguments through to Claude:

```sh
cer work -- --help
```

Use an explicit config:

```sh
cer --config examples/config.toml personal
```
