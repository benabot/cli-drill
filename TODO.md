# TODO

## P0 — Initialization

- [x] Initialize Git repository.
- [X] Run `codex init`.
- [x] Ask Codex to inspect the project and propose a plan.
- [x] Validate standard Go architecture.
- [x] Validate CLI command list.
- [x] Validate TUI MVP.
- [x] Validate config and chapter formats.
- [x] Remove unused `app/` directory if still empty.
- [x] Run `gofmt -w cmd internal data` once Go tooling is available.
- [x] Run `go mod tidy` once Go tooling is available.
- [x] Run `go test ./...` once Go tooling is available.
- [x] Run `go run ./cmd/cli-drill --help` once Go tooling is available.

## P1 — MVP foundation

- [x] Create Go module at repository root.
- [x] Add Cobra CLI skeleton in `cmd/cli-drill/`.
- [x] Add config loader.
- [x] Add YAML chapter loader.
- [x] Add progress JSON storage.
- [x] Add answer matching engine.
- [x] Add ZSH alias parser.
- [x] Add ZSH function name parser.
- [x] Add basic Markdown parser.
- [x] Add catalog model.
- [x] Deduplicate catalog entries by `(type, id)`.
- [x] Merge duplicate catalog sources and tags.
- [x] Reduce noisy Markdown heading concepts.
- [x] Reject absolute scan paths outside `dotfiles_path`.
- [x] Clarify `directory`, `search` and `show` catalog source behavior.

## P2 — Training

- [x] Add `free-answer` exercises.
- [x] Add `multiple-choice` exercises.
- [x] Add `scenario` exercises.
- [x] Add simple non-executing shell simulator.
- [x] Add `scan --summary`.
- [x] Add `scan --type <type>`.
- [x] Improve `directory --type <type>` validation and output.
- [x] Improve `search <query>` output and empty-result message.
- [x] Improve `show <entry>` detail output and ambiguity handling.
- [x] Add `init --print`.
- [x] Add concise config-missing notices.
- [x] Enrich default YAML chapters with about 5 to 10 exercises per chapter.
- [x] Cover required terminal shortcuts, navigation tools, detected aliases and
  detected ZSH functions in default chapters.
- [x] Add realistic non-executing shell-sim and scenario drills for search,
  Markdown preview and dotfiles workflows.
- [x] Add `key-sequence` shortcut drills with raw key capture for terminal
  control-key practice.
- [x] Add a lightweight command bar to `key-sequence` training for help, retry,
  next and quit actions.
- [x] Add an in-memory session review loop for missed `key-sequence` shortcut
  drills.
- [ ] Review chapter wording after several real training sessions.
- [ ] Add more advanced drills only after the MVP exercises feel stable.

## P3 — TUI MVP

- [x] Add Bubble Tea main menu.
- [x] Add chapter picker.
- [x] Add training screen.
- [x] Add directory browser.
- [x] Add stats screen.

## P3.1 — TUI onboarding and navigation

- [ ] Add a home screen that introduces cli-drill.
  - Explain what the app does.
  - Show whether a dotfiles repo is configured.
  - Offer quick access to configuration, scan, chapters, directory and stats.
- [ ] Add a first-run configuration flow.
  - Detect whether a dotfiles repo is configured.
  - Let the user point cli-drill to a dotfiles path.
  - Keep the flow safe and non-destructive.
- [ ] Improve TUI chapter context.
  - Always show the current chapter title.
  - Show current exercise position inside the chapter.
  - Show whether the user is in training, review, directory or stats.
- [ ] Improve TUI navigation.
  - Add a clear way to go back.
  - Add a clear way to return to the chapter list.
  - Add a clear way to return to the home screen.
  - Keep keyboard hints visible.
- [ ] Improve TUI chapter browsing.
  - Navigate previous/next chapter.
  - Preview chapter description and exercise count.
  - Start or resume a chapter from the picker.

## P4 — CLI training UX polish

- [x] Polish key-sequence training UX.
  - Added minimal redraw for key-sequence training.
  - Reduced repeated command bars.
  - Added stable footer.
  - Added inline `h help` / `h hide help`.
  - Kept Bubble Tea out of the training flow for now.
- [x] Improve visual polish of key-sequence training screens.
  - Added restrained colors.
  - Improved separators.
  - Improved spacing and visual hierarchy.
  - Kept output readable in plain terminals.
  - Kept ANSI styling disabled in tests / non-interactive output.
- [ ] P4.2 — Final summary polish.
  - Improve final chapter summary layout.
  - Clarify score, missed items and review action.
  - Keep summary readable without color.
  - Avoid changing raw key capture or training logic.
- [ ] Consider Bubble Tea for training mode only if the CLI renderer becomes too limited.

## P5 — Distribution and release

- [ ] Add a Makefile or Taskfile for common commands.
- [ ] Add GitHub repository metadata.
- [ ] Add GitHub Releases.
- [ ] Add release builds for macOS/Linux/Windows.
- [ ] Add Homebrew tap.

## P6 — Later

- [ ] Bash support.
- [ ] Fish support.
- [ ] Ollama-assisted chapter generation.
- [ ] MCP integrations.
- [ ] Review chapter wording after several real training sessions.
- [ ] Add more advanced drills only after the MVP exercises feel stable.
