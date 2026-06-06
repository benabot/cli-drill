# Decisions

## 2026-06-06 — Root Go Module

The Go module lives at the repository root. The executable entrypoint is
`cmd/cli-drill/`, internal packages live under `internal/`, default chapters
live under `data/chapters/`, and fixtures live under `testdata/`.

This supports the target installation command:

```bash
go install github.com/benabot/cli-drill/cmd/cli-drill@latest
```

## 2026-06-06 — Static MVP

The MVP is intentionally static and local:

- ZSH only.
- No AI.
- No MCP implementation.
- No telemetry.
- No network access.
- No execution of user aliases, functions or shell commands.

The shell simulator compares accepted answers only.

## 2026-06-06 — Editable Chapters

Generated and default lessons use YAML chapters. Generation should remain
transparent and editable rather than opaque or magical.

## 2026-06-06 — Catalog Cleanup P1

Catalog entries are deduplicated by `(type, id)`. Duplicate entries merge
sources and tags while keeping the first command/name as the canonical display
value.

Markdown headings are filtered before becoming concepts. Single-word or generic
headings are not enough to create training material.

Absolute configured scan paths are rejected by default when they point outside
`dotfiles_path`. A future explicit allowlist can relax this without changing the
default safe behavior.
