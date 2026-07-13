package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gh-jsoares/grimoire/internal/app"
	"github.com/gh-jsoares/grimoire/internal/config"
	"github.com/gh-jsoares/grimoire/internal/library"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	var (
		libraryPath    string
		tabSelect      string
		sectionSel     string
		noIcons        bool
		plain          bool
		noColor        bool
		showVersion    bool
		showCompletion string
	)

	flag.StringVar(&libraryPath, "library", "", "path to grimoire library directory")
	flag.StringVar(&tabSelect, "tab", "", "select a document tab by name or alias")
	flag.StringVar(&sectionSel, "section", "", "select a section by id")
	flag.BoolVar(&noIcons, "no-icons", false, "disable icons")
	flag.BoolVar(&plain, "plain", false, "plain mode (minimal styling)")
	flag.BoolVar(&noColor, "no-color", false, "disable color output")
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.StringVar(&showCompletion, "completion", "", "output shell completion (bash|zsh|fish)")
	flag.Usage = usage
	flag.Parse()

	if showVersion {
		fmt.Printf("grimoire %s (%s, %s)\n", version, commit, date)
		os.Exit(0)
	}

	if showCompletion != "" {
		printCompletion(showCompletion)
		os.Exit(0)
	}

	if noColor {
		_ = os.Setenv("NO_COLOR", "1")
	}

	arg := flag.Arg(0)

	// Determine mode: single file, directory, or default library
	var lib *library.Library
	singleDoc := false

	switch {
	case arg != "" && strings.HasSuffix(arg, ".grim"):
		// Single file mode
		path, err := filepath.Abs(arg)
		if err != nil {
			fatal("invalid path: %v", err)
		}
		lib, err = library.LoadSingle(path)
		if err != nil {
			fatal("%v", err)
		}
		singleDoc = true

	case arg != "" && isDir(arg):
		// Directory as library
		path, err := filepath.Abs(arg)
		if err != nil {
			fatal("invalid path: %v", err)
		}
		var loadErr error
		lib, loadErr = library.Load(path)
		if loadErr != nil {
			fatal("%v", loadErr)
		}

	default:
		// Default library
		dir := libraryPath
		if dir == "" {
			dir = config.ResolveLibraryPath()
		}
		if dir == "" {
			fatal("cannot determine library path (set $GRIMOIRE_HOME or create ~/.config/grimoire)")
		}

		var loadErr error
		lib, loadErr = library.Load(dir)
		if loadErr != nil {
			fatal("%v", loadErr)
		}
	}

	// Print warnings for any parsing errors
	if len(lib.Errors) > 0 {
		for _, err := range lib.Errors {
			fmt.Fprintf(os.Stderr, "warning: %v\n", err)
		}
	}

	if len(lib.Documents) == 0 {
		fatal("no valid documents found")
	}

	// Resolve initial document if arg is a name/alias
	if arg != "" && !strings.HasSuffix(arg, ".grim") && !isDir(arg) {
		doc, err := library.Resolve(arg, lib)
		if err != nil {
			fatal("%v", err)
		}
		// Find index
		for i := range lib.Documents {
			if lib.Documents[i].Path == doc.Path {
				tabSelect = lib.Documents[i].Title
				break
			}
		}
	}

	cfg := app.Config{
		SingleDoc:   singleDoc,
		NoIcons:     noIcons,
		Plain:       plain,
		InitTab:     tabSelect,
		InitSection: sectionSel,
	}

	model := app.New(lib, cfg)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fatal("runtime error: %v", err)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: grimoire [options] [file|name|directory]

A terminal cheatsheet viewer. Renders .grim files as styled reference cards.

Arguments:
  file          Open a .grim file directly
  name          Open a document by name or alias
  directory     Open a directory as a library

Options:
`)
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, `
Shell completion:
  grimoire --completion bash >> ~/.bashrc
  grimoire --completion zsh >> ~/.zshrc
  grimoire --completion fish > ~/.config/fish/completions/grimoire.fish
`)
}

func printCompletion(shell string) {
	switch shell {
	case "bash":
		fmt.Print(bashCompletion)
	case "zsh":
		fmt.Print(zshCompletion)
	case "fish":
		fmt.Print(fishCompletion)
	default:
		fatal("unknown shell %q (use bash, zsh, or fish)", shell)
	}
}

const bashCompletion = `_grimoire() {
    local cur="${COMP_WORDS[COMP_CWORD]}"
    local prev="${COMP_WORDS[COMP_CWORD-1]}"
    case "$prev" in
        --tab|--library|--section|--completion)
            return 0
            ;;
    esac
    if [[ "$cur" == -* ]]; then
        COMPREPLY=($(compgen -W "--tab --library --section --no-icons --plain --no-color --version --completion" -- "$cur"))
    else
        COMPREPLY=($(compgen -f -X '!*.grim' -- "$cur"))
    fi
}
complete -o default -F _grimoire grimoire
`

const zshCompletion = `#compdef grimoire

_grimoire() {
    _arguments \
        '--tab[select a document tab by name or alias]:tab name:' \
        '--library[path to grimoire library directory]:directory:_directories' \
        '--section[select a section by id]:section id:' \
        '--no-icons[disable icons]' \
        '--plain[plain mode (minimal styling)]' \
        '--no-color[disable color output]' \
        '--version[show version]' \
        '--completion[output shell completion]:shell:(bash zsh fish)' \
        '*:file:_files -g "*.grim"'
}

_grimoire "$@"
`

const fishCompletion = `complete -c grimoire -l tab -d 'Select a document tab by name or alias'
complete -c grimoire -l library -d 'Path to grimoire library directory' -r -F
complete -c grimoire -l section -d 'Select a section by id'
complete -c grimoire -l no-icons -d 'Disable icons'
complete -c grimoire -l plain -d 'Plain mode (minimal styling)'
complete -c grimoire -l no-color -d 'Disable color output'
complete -c grimoire -l version -d 'Show version'
complete -c grimoire -l completion -d 'Output shell completion' -r -f -a 'bash zsh fish'
`

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "grimoire: "+format+"\n", args...)
	os.Exit(1)
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
