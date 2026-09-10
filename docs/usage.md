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

In the interactive chooser, press `e` to open the highlighted profile's first
configured env file in `$EDITOR`. After the editor closes, the chooser remains
open.

Pass arguments through to Claude:

```sh
cer work -- --help
```

Use an explicit config:

```sh
cer --config examples/config.toml personal
```
