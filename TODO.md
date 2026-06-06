# TODO

## P0 — Initialization

- [x] Initialize Git repository.
- [ ] Run `codex init`.
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

## P3 — TUI

- [x] Add Bubble Tea main menu.
- [x] Add chapter picker.
- [x] Add training screen.
- [x] Add directory browser.
- [x] Add stats screen.

## P4 — Later

- [ ] Bash support.
- [ ] Fish support.
- [ ] Ollama-assisted chapter generation.
- [ ] MCP integrations.
- [ ] GitHub Releases.
- [ ] Homebrew tap.
