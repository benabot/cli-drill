# Project State

## Current phase

Project initialization.

`cli-drill` is being prepared as a Go CLI/TUI application that turns a
dotfiles repository into a typed command directory and a chapter-based training
tool.

## Current decisions

- Project name: `cli-drill`.
- Language: Go.
- App code lives in `app/`.
- Project documentation lives in `docs/`.
- Agent instructions live in `AGENTS.md`.
- Current project state lives in `PROJECT_STATE.md`.
- Human backlog lives in `TODO.md`.
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

## Current repository layout

```text
README.md
AGENTS.md
PROJECT_STATE.md
TODO.md
docs/
app/
mcp/
.codex/
```

## Next step

Run `codex init`, then ask Codex to read:

- `README.md`
- `AGENTS.md`
- `PROJECT_STATE.md`
- `TODO.md`
- `docs/SPEC.md`
- `.codex/skills/cli-drill-spec/SKILL.md`

Codex must propose a plan before writing application code.
