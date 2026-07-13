# Changelog

## 1.0.0 — 2026-07-13

### Added

- remove install.md, consolidate install in README with homebrew as recommended
- add shell completions, --no-color flag, clipboard error feedback, narrow terminal guard
- polish demo GIF with padding, rounded corners, nerd font, ayu theme
- add VHS demo tape, man page, GoDoc comments, and invalid test fixtures

### Documentation

- add open-source community files, tests, and examples

### Fixed

- replace retired Go Report Card badge with Go version badge
- add CLICOLOR_FORCE to bypass TTY check, use Ayu theme
- use correct VHS Env syntax for color vars
- set TERM and COLORTERM in demo.tape for color output
- bypass vhs-action, install vhs and dependencies manually
- install ffmpeg before vhs-action, add write permissions

## 0.1.0 — 2026-07-10

### Added

- skip install if already up-to-date
- initial release

### Documentation

- add vibe coded disclaimer

### Fixed

- remove homebrew tap (not created yet), fix deprecated goreleaser fields
- gitignore was excluding cmd/grimoire directory
- resolve lint errors (errcheck, unused, gosimple)

