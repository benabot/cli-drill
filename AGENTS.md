# AGENTS.md

Consignes communes pour les agents travaillant sur ce dépôt.

## Lecture obligatoire

Avant toute modification, lire :

- `README.md`
- `AGENTS.md`
- `PROJECT_STATE.md`
- `TODO.md`
- `docs/SPEC.md`

Si présents, lire aussi :

- `docs/ARCHITECTURE.md`
- `docs/DECISIONS.md`
- `docs/ROADMAP.md`
- `.codex/skills/cli-drill-spec/SKILL.md`


## Règles

- Ne pas lire de secrets.

- Ne pas scanner `~/.config/secrets`, `~/.ssh`, `~/.gnupg`,

  `~/.config/gh/hosts.yml` ou `~/.config/zed/settings.json`.

- Ne pas exécuter d’alias ou fonctions shell utilisateur.

- Ne pas lancer Docker, Colima, Ollama ou n8n.

- Ne pas modifier un repo dotfiles utilisateur.

- Ne pas modifier `.zshrc`.

- Ne pas commit/push sans validation explicite.

- Garder le MVP simple.

## Projet

Le code Go vit dans `app/`.

La documentation vit dans `docs/`.

Le dossier `mcp/` est réservé pour plus tard.

## Project state and TODO

`PROJECT_STATE.md` is versioned and must describe the current state of the
project.

Update it when a durable project decision changes.

`TODO.md` is versioned and contains the human-readable backlog.

Do not use `docs/llm/scratch/`, `docs/llm/runs/`, `docs/llm/transcripts/` or
`docs/llm/context/` as durable project memory. These directories are local and
ignored by Git.

Durable decisions must be summarized into:

- `PROJECT_STATE.md`
- `TODO.md`
- `docs/DECISIONS.md`
- `docs/ARCHITECTURE.md`
