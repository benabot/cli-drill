# Project State

## Current state

cli-drill now has a functional CLI/TUI MVP.

Implemented:

- Root Go module with Cobra commands.
- Embedded YAML chapters.
- Static dotfiles scanner for ZSH aliases, functions, tools and Markdown docs.
- Typed catalog with deduplication.
- Local JSON progress storage.
- CLI training modes:
  - free-answer;
  - multiple-choice;
  - scenario;
  - simple shell simulation;
  - key-sequence with real raw key capture.
- Polished CLI key-sequence renderer:
  - real Ctrl-key capture;
  - inline `h help`;
  - stable footer;
  - restrained colors;
  - monochrome fallback;
  - session review for missed shortcuts.
- Bubble Tea TUI cockpit:
  - Home screen;
  - chapter browser;
  - chapter detail;
  - text training screen;
  - directory browser;
  - stats screen;
  - read-only configuration status;
  - read-only scan guidance.
- Key-sequence chapters are detected by the TUI and delegated to the dedicated key training mode.
- After dedicated key training, cli-drill returns automatically to the TUI Home screen.

Current product model:

- The TUI is the navigation cockpit.
- Text-based chapters can be trained inside Bubble Tea.
- Real keyboard shortcut chapters use the dedicated key training mode.
- The user should not have to copy commands manually or relaunch cli-drill to continue.

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
- Default chapters now contain about 68 short daily-training exercises across
  terminal shortcuts, shell navigation, ZSH aliases/functions, daily tools,
  search, Markdown preview, Micro and dotfiles workflows.
- Shortcut drills in `01-raccourcis-terminal` use the `key-sequence` exercise
  type to capture real control keys in raw terminal mode during CLI training.
- `key-sequence` training uses a lightweight CLI command bar for help, retry,
  next and clean quit actions; it does not use the Bubble Tea TUI.
- `key-sequence` keeps an in-memory session review list for missed shortcuts
  and can relaunch only missed items at the end of the chapter.
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
- `key-sequence` training now uses a minimal renderer with restrained ANSI
    styling when supported by the terminal.
- Styling is centralized, disabled for tests/non-interactive output, and keeps
    all semantic labels readable without color.
- The current CLI training renderer supports inline `h help`, clean redraw,
    stable footer hints, correct/incorrect feedback, solution display and
    session-only review of missed shortcuts.

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

Next active work:

1. Improve Home visual identity.
   - Add a large ASCII `cli-drill` banner.
   - Keep the Home screen readable in narrow terminals.
   - Keep configuration and chapter actions visible.

2. Improve TUI visual consistency if needed.
   - Keep footers homogeneous.
   - Keep colors restrained and accessible.
   - Avoid introducing a large UI framework change.

Deferred:

- Return to chapter detail after dedicated key training instead of Home.
- Polish final training summaries after real usage.
- Prepare public GitHub README, screenshots and GIF.
- GitHub Releases and Homebrew tap.
