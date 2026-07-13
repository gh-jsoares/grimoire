Output demo.gif

Set Shell "bash"
Set FontSize 14
Set Env "TERM" "xterm-256color"
Set Env "COLORTERM" "truecolor"
Set Width 1200
Set Height 700
Set Theme "Catppuccin Mocha"

Type "./grimoire examples/"
Enter
Sleep 1s

# Browse the docker cheatsheet
Sleep 1.5s
Type "j"
Sleep 300ms
Type "j"
Sleep 300ms
Type "j"
Sleep 300ms

# Navigate commands
Type "n"
Sleep 400ms
Type "n"
Sleep 400ms
Type "n"
Sleep 400ms

# Switch to vim tab
Type "L"
Sleep 1.5s

# Scroll down
Type "j"
Sleep 200ms
Type "j"
Sleep 200ms
Type "j"
Sleep 200ms

# Search
Type "/"
Sleep 300ms
Type "undo"
Sleep 800ms
Enter
Sleep 1.5s

# Clear search
Type@50ms ""
Sleep 300ms

# Quit
Type "q"
Sleep 500ms
