# cli-drill

`cli-drill` est une application CLI/TUI d’entraînement pour apprendre son shell,
ses dotfiles, ses alias, ses fonctions, ses raccourcis, ses outils CLI et ses
workflows.

Le projet est perso-first, generic-ready : il doit d’abord fonctionner avec un
repo dotfiles personnel, puis être généralisable.

## Structure

```text
README.md          Présentation publique du projet
AGENTS.md          Consignes communes pour agents de code
PROJECT_STATE.md   État courant du projet
TODO.md            Backlog humain et priorités
docs/              Spécification et documentation projet
app/               Code Go de l’application
mcp/               Réservé pour futures intégrations MCP
```

Le code applicatif vit dans `app/`.

La documentation stable vit dans `docs/`.

Les notes temporaires d’agents peuvent vivre dans `docs/llm/`, mais seuls les
fichiers stables doivent être versionnés.

Les sous-dossiers suivants ne sont pas versionnés :

```text
docs/llm/scratch/
docs/llm/runs/
docs/llm/transcripts/
docs/llm/context/
```

## MVP

* Go
* CLI + TUI
* ZSH au départ
* annuaire typé
* chapitres YAML éditables
* progression locale
* aucun secret lu
* aucune commande utilisateur exécutée
