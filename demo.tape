Output demo.gif

Set Shell "bash"
Set FontSize 14
Set FontFamily "JetBrainsMono Nerd Font"
Set Width 1200
Set Height 700
Set Theme "ayu"
Set TypingSpeed 75ms
Set Padding 20
Set Framerate 60
Set Margin 20
Set MarginFill "#1A1F29"
Set BorderRadius 10
Set WindowBar Colorful

Env TERM "xterm-256color"
Env COLORTERM "truecolor"
Env CLICOLOR_FORCE "1"

Type "grimoire examples/"
Enter
Sleep 1s

# Browse the docker cheatsheet
Sleep 1.5s
Type "j"
Sleep 500ms
Type "j"
Sleep 500ms
Type "j"
Sleep 500ms

# Navigate commands
Type "n"
Sleep 800ms
Type "n"
Sleep 800ms
Type "n"
Sleep 800ms

# Switch to vim tab
Type "L"
Sleep 1.5s

# Scroll down
Type "j"
Sleep 500ms
Type "j"
Sleep 800ms

# Search
Type "/"
Sleep 300ms
Type "undo"
Sleep 800ms
Enter
Sleep 1.5s

# Clear search
Escape
Sleep 500ms

# Switch back to docker tab and show link picker
Type "H"
Sleep 1s
Type "G"
Sleep 800ms
Type "u"
Sleep 2s
Escape
Sleep 500ms

# Quit
Type "q"
Sleep 500ms
