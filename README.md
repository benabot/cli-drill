# cli-drill

`cli-drill` est une application CLI/TUI d'entraînement pour apprendre son shell,
ses dotfiles, ses alias, ses fonctions, ses raccourcis, ses outils CLI et ses
workflows.

Le projet est **perso-first, generic-ready** : il doit d'abord fonctionner avec
un repo dotfiles personnel, puis être généralisable à d'autres utilisateurs.

## Structure

```text
README.md          Présentation publique du projet
AGENTS.md          Consignes communes pour agents de code
PROJECT_STATE.md   État courant du projet
TODO.md            Backlog humain et priorités
go.mod             Module Go à la racine du repo
cmd/               Entrées exécutables Go
internal/          Code interne de l'application
data/              Chapitres par défaut et données embarquables
testdata/          Fixtures de tests
docs/              Spécification et documentation projet
mcp/               Réservé pour futures intégrations MCP
```

Le module Go vit à la racine du repo.

L'entrée principale du binaire vit dans `cmd/cli-drill/`.

Installation locale depuis le repo :

```bash
go install ./cmd/cli-drill
```

Installation future depuis GitHub, une fois le repo public et accessible :

```bash
go install github.com/benabot/cli-drill/cmd/cli-drill@latest
```

Cette commande fonctionnera quand le dépôt GitHub public sera disponible.

## Commandes utilisateur

Lancer la TUI :

```bash
cli-drill
```

Lister les chapitres :

```bash
cli-drill chapters
```

Lancer un chapitre :

```bash
cli-drill train 01-raccourcis-terminal
```

Scanner les dotfiles :

```bash
cli-drill scan --summary
```

Voir l'annuaire :

```bash
cli-drill directory
cli-drill directory --type tool
```

Chercher une entrée :

```bash
cli-drill search rg
```

Voir une entrée :

```bash
cli-drill show tool-rg
```

Créer ou prévisualiser la config :

```bash
cli-drill init --print
cli-drill init
```

La documentation stable vit dans `docs/`.

Les notes temporaires d'agents peuvent vivre dans `docs/llm/`, mais seuls les
fichiers stables doivent être versionnés.

Les sous-dossiers suivants ne sont pas versionnés :

```text
docs/llm/scratch/
docs/llm/runs/
docs/llm/transcripts/
docs/llm/context/
```

## MVP

- Go
- CLI + TUI
- ZSH au départ
- annuaire typé
- chapitres YAML éditables
- progression locale
- aucun secret lu
- aucune commande utilisateur exécutée
- pas d'IA
- pas de MCP
- pas de télémétrie

## Documentation

Lire en priorité :

- `docs/SPEC.md`
- `PROJECT_STATE.md`
- `TODO.md`
- `AGENTS.md`

## Développement local

Après initialisation du module Go :

```bash
go test ./...
go run ./cmd/cli-drill --help
```
