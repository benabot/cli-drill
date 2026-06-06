# Architecture

`cli-drill` is a standard Go module rooted at the repository top level.

The product flow is:

```text
scan dotfiles -> typed directory -> editable YAML chapters -> training
```

## Layout

```text
cmd/cli-drill/     Cobra executable entrypoint
internal/app/      CLI command wiring and use-case orchestration
internal/catalog/  Typed directory model, filtering and search
internal/chapter/  YAML chapter model, loading and generation
internal/config/   TOML configuration model and persistence
internal/detect/   Static scan orchestration and safe path filtering
internal/exercise/ Answer matching and exercise types
internal/markdown/ Small Markdown extractor
internal/progress/ Local JSON progress storage
internal/shell/zsh Static ZSH alias/function parsers
internal/tui/      Bubble Tea MVP screens
internal/xdg/      Config/data path resolution
data/chapters/     Default editable chapters
mcp/               Reserved for later
testdata/          Fixtures for tests
```

## Boundaries

- Parsers never execute user aliases, functions or commands.
- Detection reads only configured paths and applies the security exclude list.
- Absolute configured paths are rejected unless they stay inside
  `dotfiles_path`.
- Catalog entries are deduplicated by `(type, id)` and merge their
  sources/tags.
- Markdown headings are filtered before they become `concept` entries.
- CLI and TUI share the same internal packages.
- Default chapters are embedded so `go install` can ship a usable first run.
- User-generated chapters remain editable YAML files.
- `directory`, `search` and `show` use the scan catalog when a config is
  provided or present; otherwise they use the default chapter catalog.
- CLI list output uses a stable tabular shape: `type`, `id`, `name`,
  `summary`.
- `show` resolves exact `id` or `name`; ambiguous exact matches return a clear
  error rather than choosing silently.
- Without config, `directory` and `search` use embedded chapters; `show`
  falls back to default scan paths only when the requested entry is not found
  in embedded chapters.

## Validation

Expected validation commands:

```bash
gofmt -w cmd internal data
go mod tidy
go test ./...
go run ./cmd/cli-drill --help
```
