# Project State

## Current phase

Project initialization.

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

`go.mod` and `go.sum` may not exist yet during initialization.

## Next step

Ask Codex to read:

- `README.md`
- `AGENTS.md`
- `PROJECT_STATE.md`
- `TODO.md`
- `docs/SPEC.md`
- `.codex/skills/cli-drill-spec/SKILL.md`

Codex must update its plan according to the standard Go layout before writing
application code.
