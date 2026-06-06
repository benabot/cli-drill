# TODO

## P0 — Initialization

- [x] Initialize Git repository.
- [ ] Run `codex init`.
- [ ] Ask Codex to inspect the project and propose a plan.
- [ ] Validate standard Go architecture.
- [ ] Validate CLI command list.
- [ ] Validate TUI MVP.
- [ ] Validate config and chapter formats.
- [ ] Remove unused `app/` directory if still empty.

## P1 — MVP foundation

- [ ] Create Go module at repository root.
- [ ] Add Cobra CLI skeleton in `cmd/cli-drill/`.
- [ ] Add config loader.
- [ ] Add YAML chapter loader.
- [ ] Add progress JSON storage.
- [ ] Add answer matching engine.
- [ ] Add ZSH alias parser.
- [ ] Add ZSH function name parser.
- [ ] Add basic Markdown parser.
- [ ] Add catalog model.

## P2 — Training

- [ ] Add `free-answer` exercises.
- [ ] Add `multiple-choice` exercises.
- [ ] Add `scenario` exercises.
- [ ] Add simple non-executing shell simulator.

## P3 — TUI

- [ ] Add Bubble Tea main menu.
- [ ] Add chapter picker.
- [ ] Add training screen.
- [ ] Add directory browser.
- [ ] Add stats screen.

## P4 — Later

- [ ] Bash support.
- [ ] Fish support.
- [ ] Ollama-assisted chapter generation.
- [ ] MCP integrations.
- [ ] GitHub Releases.
- [ ] Homebrew tap.
