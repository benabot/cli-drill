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

## 2026-06-06 — CLI UX P2.1

The CLI stays minimal and file-based. `scan` supports summary and type-filtered
views without storing scan state. `directory`, `search` and `show` keep using
the active catalog source selected by config presence.

When no config exists, `directory` and `search` use embedded chapters. `show`
may fall back to the default scan paths for direct lookup, which keeps audit
commands like `show tool-rg` useful without introducing scan storage.

List-style output is stable and tabular: `type`, `id`, `name`, `summary`.
Detailed output from `show` includes source and tags when present.

`init --print` prints the default TOML config without writing to disk.

## 2026-06-06 — Default Training Chapters P2.2

Default chapters are short, topic-separated YAML drills intended for daily
terminal practice. They favor realistic macOS/ZSH power-user workflows while
remaining static and non-executing.

Exercises may reference detected aliases and functions when they are present in
the catalog. Generic concepts stay explicit as concepts or workflows rather
than pretending to be user-specific shell definitions.

The embedded chapter set now covers terminal shortcuts, navigation, aliases,
functions, daily tools, search, reading/preview, Micro, Markdown and dotfiles
workflows with a mix of `free-answer`, `multiple-choice`, `scenario` and
`simple-shell-sim` exercises.

## 2026-06-06 — Raw Shortcut Training P2.4/P2.4b/P2.4c

Shortcut drills use raw key capture through `key-sequence`; text answers remain
for other exercise types.

Only `key-sequence` exercises enter raw terminal mode. The training loop
restores the terminal state with `defer` after each captured key, maps control
bytes to canonical `Ctrl+<letter>` notation and treats `Esc` or `Ctrl+C` as a
clean session exit.

`key-sequence` uses raw key capture with a lightweight command bar, not a full
TUI. After each captured answer, `Enter` advances, `r` retries the same
exercise, and `Esc` or `Ctrl+C` quits cleanly.

The command bar is intentionally contextual. Correct answers only offer next or
quit. Missed answers offer next, retry or solution. `key-sequence` uses a
lightweight drill loop with session review for missed shortcuts; the review
list is in-memory only and does not add persistent storage.
