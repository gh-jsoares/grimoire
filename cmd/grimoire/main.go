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
		libraryPath string
		tabSelect   string
		sectionSel  string
		noIcons     bool
		plain       bool
		showVersion bool
	)

	flag.StringVar(&libraryPath, "library", "", "path to grimoire library directory")
	flag.StringVar(&tabSelect, "tab", "", "select a document tab by name or alias")
	flag.StringVar(&sectionSel, "section", "", "select a section by id")
	flag.BoolVar(&noIcons, "no-icons", false, "disable icons")
	flag.BoolVar(&plain, "plain", false, "plain mode (minimal styling)")
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.Parse()

	if showVersion {
		fmt.Printf("grimoire %s (%s, %s)\n", version, commit, date)
		os.Exit(0)
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

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "grimoire: "+format+"\n", args...)
	os.Exit(1)
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
