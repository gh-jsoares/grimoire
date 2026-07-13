# tmux Integration

## Popup

Add to your `tmux.conf` to open grimoire in a popup with `?`:

```sh
bind-key ? display-popup -E -w 90% -h 85% "grimoire"
```

## Tips

- The popup dimensions (`-w 90% -h 85%`) give enough room for the 2-column layout in most terminals
- Use `-E` to close the popup when grimoire exits
- You can pass arguments: `"grimoire git"` to open a specific cheatsheet directly
