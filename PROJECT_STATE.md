# Project State

## Current phase

MVP foundation implemented and validated.

`cli-drill` is being prepared as a Go CLI/TUI application that turns a
dotfiles repository into a typed command directory and a chapter-based training
tool.

## Current decisions

- Project name: `cli-drill`.
- Language: Go.
- Go module lives at the repository root.
- Executable entrypoint lives in `cmd/cli-drill/`.
- Internal Go code lives in `internal/`.
- Default chapters and embeddable data live in `data/`.
- Test fixtures live in `testdata/`.
- Project documentation lives in `docs/`.
- Agent instructions live in `AGENTS.md`.
- Current project state lives in `PROJECT_STATE.md`.
- Human backlog lives in `TODO.md`.
- `app/` is not used for the MVP.
- MVP shell support: ZSH only.
- Future shell support: Bash and Fish.
- Config format: TOML.
- Chapter format: YAML.
- Progress format: JSON.
- No AI in the MVP.
- No MCP implementation in the MVP.
- No telemetry.
- No network access.
- No execution of user aliases or shell functions.
- Static scan by default.
- Target install command:
  `go install github.com/benabot/cli-drill/cmd/cli-drill@latest`.

## Current implementation

- Root Go module declared in `go.mod`.
- Cobra CLI entrypoint lives in `cmd/cli-drill/`.
- Internal packages exist for app orchestration, catalog, chapters, config,
  detection, exercises, Markdown parsing, progress, ZSH parsing, TUI and XDG
  paths.
- Default editable chapters live in `data/chapters/` and are embedded through
  `data/defaults.go`.
- Unit tests cover ZSH parsing, answer matching, chapter loading, config
  loading, safe path filtering, catalog deduplication, scanner noise reduction
  and progress storage.
- The TUI includes a main menu, chapter picker, training screen, directory
  browser and stats summary.
- CLI commands implemented:
  `init`, `scan`, `generate`, `chapters`, `train`, `directory`, `search`,
  `show`, `stats`, `reset`.
- CLI UX P2.1 includes `scan --summary`, `scan --type <type>`,
  `init --print`, clearer empty-result messages, typed tabular directory/search
  output and detailed `show` output.
- The catalog deduplicates entries by `(type, id)` and merges tags/sources.
- Markdown headings are filtered before becoming concepts to avoid generic
  exercise noise.
- Absolute configured scan paths outside `dotfiles_path` are rejected by
  default.
- `directory`, `search` and `show` use the scan catalog when a config is
  provided or present; otherwise they use embedded/default chapters.
- Without config, `show` can fall back to default scan paths when an entry is
  not found in embedded chapters.

## Validation status

The current implementation was validated with:

```bash
gofmt -w cmd internal data
go mod tidy
go test ./...
go run ./cmd/cli-drill --help
```

## Current repository layout

```text
README.md
AGENTS.md
PROJECT_STATE.md
TODO.md
go.mod
go.sum
cmd/
internal/
data/
docs/
mcp/
testdata/
.codex/
```

`go.sum` is generated.

## Next step

Review the MVP behavior manually, then decide whether to refine CLI/TUI UX or
start using a real dotfiles config for scan/generation testing.
