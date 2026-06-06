# Skill — cli-drill specification guardian

## Purpose

Use this skill when working on the `cli-drill` project.

`cli-drill` is a Go CLI/TUI application that turns a dotfiles repository into a
typed command directory and a chapter-based training app.

The skill exists to keep implementation aligned with the product specification
and to prevent scope creep.

## Project location

```text
/Users/benoitabot/Sites/cli-drill
```

Expected root structure:

```text
README.md
AGENTS.md
docs/
app/
mcp/
```

Application code lives in:

```text
app/
```

Do not put Go application code directly at the repository root.

## Mandatory reading

Before making changes, read:

```text
README.md
AGENTS.md
docs/SPEC.md
```

If present, also read:

```text
docs/ARCHITECTURE.md
docs/DECISIONS.md
docs/ROADMAP.md
docs/CONFIG_FORMAT.md
docs/CHAPTER_FORMAT.md
```

## Core product rule

The product flow is:

```text
scan dotfiles -> typed directory -> editable YAML chapters -> training
```

Do not turn the MVP into a general automation framework.

Do not add AI, MCP, network access, telemetry or sync to the MVP.

## Technical stack

Use Go.

Expected stack:

```text
Cobra       CLI
Bubble Tea  TUI
Bubbles     TUI components
Lip Gloss   TUI styling
TOML        configuration
YAML        chapters
JSON        local progress
```

Do not use Python, Swift or Electron.

## MVP shell support

MVP supports:

```text
ZSH only
```

Architecture may prepare future support for Bash and Fish, but avoid excessive
abstraction.

## Security rules

Never read:

```text
~/.config/secrets
~/.ssh
~/.gnupg
~/.config/gh/hosts.yml
~/.config/zed/settings.json
```

Never:

- print environment variables containing secrets;
- execute user aliases;
- execute user functions;
- launch Docker;
- launch Colima;
- launch Ollama;
- launch n8n;
- modify the user's dotfiles repository;
- modify `.zshrc`;
- commit or push without explicit validation.

Scan statically by default.

Runtime scanning must be explicit, opt-in and safe.

## Data sources

For Benoît's personal setup, useful sources include:

```text
~/dotfiles/README.md
~/dotfiles/zsh/.zshrc
~/dotfiles/zsh/README.md
~/dotfiles/zsh/modules/aliases.zsh
~/dotfiles/zsh/modules/functions.zsh
~/dotfiles/zsh/modules/tools.zsh
~/dotfiles/docs/tools-inventory.md
~/dotfiles/docs/cli-tools-usage.md
~/dotfiles/micro/DOC.md
~/dotfiles/micro/settings.json
~/dotfiles/micro/bindings.json
~/dotfiles/glow/glow.yml
~/dotfiles/yazi/yazi.toml
~/dotfiles/yazi/keymap.toml
~/dotfiles/lazygit/config.yml
```

These paths must be configurable. Do not hard-code them as universal truth.

## Directory entries

The app builds a typed directory.

Minimum entry types:

```text
shortcut
alias
function
tool
workflow
concept
binding
chapter
```

Examples:

```text
Ctrl+A       shortcut   move to beginning of line
cgs          alias      git status --short
y            function   open Yazi and cd to selected directory
rg           tool       search file contents
mdpreview    workflow   MarkEdit + Marked 2 preview workflow
```

## Training exercise types

MVP exercise types:

```text
free-answer
multiple-choice
scenario
simple-shell-sim
```

The simple shell simulator must not execute commands. It only compares user
input with accepted answers.

## Chapter rules

Chapters are editable YAML files.

Initial MVP chapters:

```text
01-raccourcis-terminal
02-navigation-shell
03-alias-zsh
04-fonctions-zsh
05-outils-quotidiens
06-recherche-fichiers-contenu
07-lecture-preview
08-micro
09-markdown-glow
10-workflows-dotfiles
```

Keep chapters separated. Do not mix all tools into one giant lesson.

## Config and storage

User config:

```text
~/.config/cli-drill/config.toml
```

Progress:

```text
~/.local/share/cli-drill/progress.json
```

Support macOS/Linux/Windows equivalents through proper path handling.

## Implementation guidance

Keep the code simple.

Prefer clear packages over over-engineered abstractions.

Expected application structure:

```text
app/
├── go.mod
├── cmd/
│   └── cli-drill/
│       └── main.go
├── internal/
│   ├── app/
│   ├── catalog/
│   ├── chapter/
│   ├── config/
│   ├── detect/
│   ├── exercise/
│   ├── markdown/
│   ├── progress/
│   ├── shell/
│   │   └── zsh/
│   ├── tui/
│   └── xdg/
├── data/
│   └── chapters/
└── testdata/
```

Adjust only if there is a clear reason.

## Testing expectations

Add tests for:

- ZSH alias parsing;
- ZSH function name detection;
- answer matching;
- YAML chapter loading;
- TOML config loading;
- safe path filtering.

Useful commands:

```bash
cd /Users/benoitabot/Sites/cli-drill/app
go test ./...
go test ./internal/...
go run ./cmd/cli-drill --help
```

## Change policy

Before large changes:

1. summarize the intent;
2. list touched files;
3. explain risks;
4. propose tests.

After changes:

1. run tests;
2. summarize results;
3. do not commit or push unless explicitly asked.

## Scope control

Do not implement in the MVP:

- Ollama integration;
- AI chapter generation;
- MCP tools;
- Homebrew tap;
- GitHub releases;
- telemetry;
- account system;
- sync;
- real command execution in shell simulator;
- Bash/Fish support.

## When in doubt

Favor:

```text
small
clear
testable
safe
editable
```

over:

```text
automatic
clever
abstract
magical
large
```
