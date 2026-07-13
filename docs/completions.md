# Shell Completions

Grimoire supports completions for bash, zsh, and fish via the `--completion` flag.

## Quick Setup

Add to your shell config:

```sh
# bash (~/.bashrc)
eval "$(grimoire --completion bash)"

# zsh (~/.zshrc, before compinit)
eval "$(grimoire --completion zsh)"

# fish (~/.config/fish/config.fish)
grimoire --completion fish | source
```

## Static Generation (recommended)

Evaluating completions on every shell startup adds latency. For faster shells, generate a static file and source it instead.

### zsh

```sh
# Generate once (or regenerate after upgrading grimoire)
grimoire --completion zsh > "${XDG_CACHE_HOME:-$HOME/.cache}/zsh/completions/_grimoire"
```

Make sure the completions directory is in your `fpath`:

```sh
fpath=("${XDG_CACHE_HOME:-$HOME/.cache}/zsh/completions" $fpath)
autoload -Uz compinit && compinit
```

To regenerate daily in the background:

```sh
_comp_dir="${XDG_CACHE_HOME:-$HOME/.cache}/zsh/completions"
if [[ ! -f "$_comp_dir/_grimoire" || ! $(find "$_comp_dir/_grimoire" -newermt "24 hours ago" -print) ]]; then
  grimoire --completion zsh >| "$_comp_dir/_grimoire" 2>/dev/null &|
fi
```

### bash

```sh
# Generate once
grimoire --completion bash > "${XDG_DATA_HOME:-$HOME/.local/share}/bash-completion/completions/grimoire"
```

bash-completion will source files from this directory automatically.

### fish

```sh
# Generate once
grimoire --completion fish > "${XDG_CONFIG_HOME:-$HOME/.config}/fish/completions/grimoire.fish"
```

Fish automatically loads files from `~/.config/fish/completions/`.
