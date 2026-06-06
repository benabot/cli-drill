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
- Ne pas exécuter d'alias ou fonctions shell utilisateur.
- Ne pas lancer Docker, Colima, Ollama ou n8n.
- Ne pas modifier un repo dotfiles utilisateur.
- Ne pas modifier `.zshrc`.
- Ne pas commit/push sans validation explicite.
- Garder le MVP simple.

## Projet

Le module Go vit à la racine du repo.

Le point d'entrée du binaire vit dans :

```text
cmd/cli-drill/
```

Le code interne vit dans :

```text
internal/
```

Les chapitres par défaut et données embarquables vivent dans :

```text
data/
```

Les fixtures de tests vivent dans :

```text
testdata/
```

La documentation vit dans :

```text
docs/
```

Le dossier `mcp/` est réservé pour plus tard.

## Project state and TODO

`PROJECT_STATE.md` est versionné et doit décrire l'état courant du projet.

Le mettre à jour quand une décision durable change.

`TODO.md` est versionné et contient le backlog lisible par un humain.

Ne pas utiliser `docs/llm/scratch/`, `docs/llm/runs/`,
`docs/llm/transcripts/` ou `docs/llm/context/` comme mémoire durable du
projet. Ces dossiers sont locaux et ignorés par Git.

Les décisions durables doivent être résumées dans :

- `PROJECT_STATE.md`
- `TODO.md`
- `docs/DECISIONS.md`
- `docs/ARCHITECTURE.md`

## Politique de changement

Avant un changement important :

1. résumer l'intention ;
2. lister les fichiers prévus ;
3. expliquer les risques ;
4. proposer les tests.

Après modification :

1. lancer les tests pertinents ;
2. résumer les résultats ;
3. ne pas commit/push sans validation explicite.
