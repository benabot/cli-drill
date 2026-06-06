# cli-drill — Cahier des charges

## 1. Vision

`cli-drill` est une application CLI/TUI d’entraînement qui transforme un repo
dotfiles existant en parcours d’apprentissage chapitré.

L’objectif n’est pas seulement de documenter des commandes, mais d’aider
l’utilisateur à les mémoriser par la pratique.

Le principe produit est :

```text
scan dotfiles -> annuaire typé -> chapitres YAML éditables -> entraînement
```

Le projet est :

```text
perso-first, generic-ready
```

Il doit d’abord fonctionner avec le repo dotfiles de Benoît, puis être conçu
pour être diffusable à d’autres utilisateurs ayant leur propre repo dotfiles.

## 2. Nom

```text
cli-drill
```

## 3. Public cible

### Public initial

Benoît, développeur / power user macOS, utilisant :

- ZSH ;
- Ghostty ;
- dotfiles ;
- alias et fonctions shell ;
- outils CLI modernes ;
- Micro ;
- Yazi ;
- Atuin ;
- zoxide ;
- fzf ;
- fd ;
- rg ;
- bat ;
- Glow ;
- eza ;
- lazygit ;
- agents IA comme Codex, Claude Code et Gemini CLI.

### Public futur

Power users, développeurs, administrateurs systèmes ou utilisateurs Linux/macOS
qui maintiennent un repo dotfiles et veulent apprendre leur propre environnement.

## 4. Emplacement du projet

Le projet vit ici :

```text
/Users/benoitabot/Sites/cli-drill
```

Structure racine attendue :

```text
/Users/benoitabot/Sites/cli-drill/
├── README.md
├── AGENTS.md
├── docs/
├── app/
└── mcp/
```

Rôle des éléments :

```text
README.md   Présentation publique du projet.
AGENTS.md   Consignes agnostiques pour Codex, Claude Code, Gemini CLI, etc.
docs/       Spécification, architecture, décisions, formats, roadmap.
app/        Code Go de l’application.
mcp/        Réservé pour futures intégrations MCP / LLM tools.
```

Le code Go ne doit pas vivre directement à la racine. Il doit vivre dans `app/`.

## 5. Stack technique

Langage :

```text
Go
```

Bibliothèques prévues :

```text
Cobra       CLI
Bubble Tea  TUI
Bubbles     composants TUI
Lip Gloss   style TUI
TOML        configuration utilisateur
YAML        chapitres éditables
JSON        progression locale
```

Ne pas utiliser :

- Python ;
- Swift ;
- Electron ;
- base de données ;
- télémétrie ;
- réseau ;
- IA dans le MVP.

Ollama pourra être envisagé plus tard pour aider à générer des exercices, mais
pas dans le MVP.

## 6. Objectifs MVP

Le MVP doit permettre :

1. d’initialiser une configuration ;
2. de scanner un repo dotfiles ;
3. de créer un annuaire typé ;
4. de générer des chapitres YAML éditables ;
5. de lancer un entraînement ;
6. de consulter la progression ;
7. de consulter l’annuaire ;
8. de fonctionner en CLI et en TUI.

## 7. Shells supportés

MVP :

```text
ZSH uniquement
```

Architecture future :

```text
ZSH
Bash
Fish
```

Le code doit prévoir une extension future sans créer une abstraction excessive.

## 8. Sources de données

### Sources dotfiles

Exemples de fichiers à exploiter dans le cas de Benoît :

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

### Sources pédagogiques

Fichiers pédagogiques pouvant servir d’inspiration :

```text
ZSH — Mode d’emploi des fonctions et alias.md
Tutoriel — outils CLI essentiels.md
Raccourcis shortcut terminal bash ZSH.md
```

L’application doit toutefois rester générique. Les chemins doivent être
configurables.

## 9. Sécurité

Règles strictes :

- ne jamais lire `~/.config/secrets` ;
- ne jamais lire `~/.ssh` ;
- ne jamais lire `~/.gnupg` ;
- ne jamais lire `~/.config/gh/hosts.yml` ;
- ne jamais lire `~/.config/zed/settings.json` ;
- ne jamais afficher de variables d’environnement sensibles ;
- ne jamais exécuter d’alias ;
- ne jamais exécuter de fonctions shell ;
- ne jamais lancer Docker ;
- ne jamais lancer Colima ;
- ne jamais lancer Ollama ;
- ne jamais lancer n8n ;
- ne jamais modifier le repo dotfiles de l’utilisateur ;
- ne jamais modifier `.zshrc` ;
- ne jamais commit/push sans validation humaine.

Le scan doit être statique par défaut.

Un scan runtime pourra être envisagé plus tard, mais il devra être explicite,
opt-in et sécurisé.

## 10. Annuaire typé

`cli-drill` doit générer un annuaire typé.

Types minimum :

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

Exemples :

```text
Ctrl+A       shortcut   Aller au début de la ligne
cgs          alias      git status --short
y            function   Ouvrir Yazi et revenir dans le dernier dossier
rg           tool       Chercher dans le contenu des fichiers
mdpreview    workflow   MarkEdit + Marked 2
Ctrl+R       binding    Atuin search
```

Commandes prévues :

```bash
cli-drill directory
cli-drill directory --type alias
cli-drill directory --type tool
cli-drill search atuin
cli-drill show cgs
```

À terme, l’annuaire pourra être enrichi par consultation de manuels, pages `man`,
Dash, Context7 ou MCP. Pas dans le MVP.

## 11. Entraînement

Modes d’exercice MVP :

```text
free-answer
multiple-choice
scenario
simple-shell-sim
```

Le shell simulé est indispensable, mais doit rester simple au MVP.

Il ne doit jamais exécuter de vraies commandes. Il compare seulement une réponse
saisie avec des réponses attendues.

Exemples :

```text
Question libre :
Tu veux aller au début de la ligne. Quel raccourci ?
Réponse : Ctrl+A

Choix multiple :
Quel outil sert à chercher dans le contenu des fichiers ?
A. fd
B. rg
C. eza
D. bat

Scénario :
Tu veux retrouver une commande déjà exécutée concernant gitleaks.
Quel outil utilises-tu ?

Shell simulé :
$ _
Objectif : lister les fichiers Markdown et en choisir un avec fzf.
Réponse attendue : fd -e md | fzf
```

Le moteur de correction doit accepter plusieurs variantes valides.

## 12. Chapitres MVP

Créer des chapitres séparés.

Ne pas tout mélanger.

Chapitres initiaux :

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

Outils prioritaires à couvrir s’ils sont détectés ou documentés :

```text
Atuin
Yazi
zoxide
fzf
fd
rg
bat
Glow
eza
micro
lazygit
lazydocker
git-delta
sd
tokei
xh
doggo
gum
viddy
gitleaks
hadolint
vale
markdownlint-cli2
swiftlint
swiftformat
xcodegen
```

La liste n’est pas fermée. Le scan doit pouvoir détecter d’autres outils utiles.

## 13. Génération

Approche retenue :

```text
scan automatique -> catalogue brut -> génération de chapitres YAML éditables
```

L’utilisateur doit pouvoir modifier les chapitres générés.

Ne pas enfermer l’utilisateur dans une génération opaque.

## 14. Configuration

Format :

```text
TOML
```

Emplacement utilisateur :

```text
~/.config/cli-drill/config.toml
```

Prévoir les équivalents macOS/Linux/Windows.

Exemple :

```toml
dotfiles_path = "~/dotfiles"
shell = "zsh"

[paths]
aliases = ["zsh/modules/aliases.zsh"]
functions = ["zsh/modules/functions.zsh"]
docs = [
  "README.md",
  "docs/tools-inventory.md",
  "docs/cli-tools-usage.md",
  "micro/DOC.md"
]

[security]
exclude = [
  "~/.config/secrets",
  "~/.ssh",
  "~/.gnupg",
  "~/.config/gh/hosts.yml",
  "~/.config/zed/settings.json"
]
```

## 15. Chapitres

Format :

```text
YAML
```

Emplacement utilisateur futur :

```text
~/.config/cli-drill/chapters/*.yaml
```

Emplacement projet pendant développement :

```text
app/data/chapters/*.yaml
```

Exemple :

```yaml
id: 01-raccourcis-terminal
title: Raccourcis terminal
description: Mémoriser les raccourcis de navigation et édition de ligne.
items:
  - id: ctrl-a
    type: shortcut
    exercise_type: free-answer
    prompt: Aller au début de la ligne
    answer:
      primary: Ctrl+A
      accepted:
        - "^A"
        - "control-a"
        - "ctrl a"
    explanation: Ctrl+A place le curseur au début de la ligne.
```

## 16. Progression

Stocker la progression hors dotfiles.

Emplacement souhaité :

```text
~/.local/share/cli-drill/progress.json
```

Prévoir équivalents macOS/Linux/Windows.

Contenu minimum :

```json
{
  "version": 1,
  "completed_exercises": [],
  "chapter_scores": {},
  "last_session": null,
  "streak": 0
}
```

## 17. Parsing MVP

### ZSH aliases

Détecter :

```zsh
alias cgs='git status --short'
alias lg='lazygit'
```

Ignorer :

```zsh
# alias old='...'
```

### ZSH functions

Détecter les noms :

```zsh
y() {
  ...
}

function foo() {
  ...
}
```

Ne jamais exécuter les fonctions.

### Markdown

Parser :

- headings ;
- tableaux Markdown ;
- blocs de code ;
- listes simples.

Objectif : extraire des éléments exploitables, pas comprendre parfaitement tout
le Markdown.

### Config outils

Parser prudemment :

```text
micro/bindings.json
yazi/keymap.toml
lazygit/config.yml
glow/glow.yml
```

Si le parsing est fragile, importer comme référence textuelle au lieu de casser.

## 18. Commandes CLI prévues

```bash
cli-drill init
cli-drill scan
cli-drill generate
cli-drill chapters
cli-drill train
cli-drill train 01-raccourcis-terminal
cli-drill directory
cli-drill directory --type alias
cli-drill search <query>
cli-drill show <entry>
cli-drill stats
cli-drill reset
```

## 19. TUI MVP

Écran principal :

```text
cli-drill

1. Continuer
2. Choisir un chapitre
3. Annuaire
4. Scan dotfiles
5. Stats
6. Quitter
```

Écran chapitre :

```text
Chapitre : Alias ZSH
Progression : 4/20

Question :
Quel alias affiche git status --short ?

Réponse :
> _
```

Feedback :

```text
Correct.
cgs = git status --short
```

ou :

```text
Pas encore.
Réponse attendue : cgs
```

## 20. Architecture applicative

Structure cible :

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

Cette architecture peut être ajustée, mais tout changement doit être justifié.

## 21. Distribution

Objectif futur :

```bash
go install github.com/benabot/cli-drill@latest
```

Plus tard :

```text
GitHub Releases avec binaires
Homebrew tap
```

Ne pas mettre en place Homebrew tap dans le MVP.

## 22. Qualité attendue

- code Go simple ;
- architecture lisible ;
- pas de sur-abstraction ;
- tests unitaires sur les parsers ;
- tests unitaires sur le matching de réponses ;
- README clair ;
- docs courtes mais utiles ;
- aucune dépendance inutile ;
- pas de télémétrie ;
- pas de réseau ;
- pas d’IA.

## 23. Contraintes agents

Les agents doivent :

- lire `README.md`, `AGENTS.md` et `docs/SPEC.md` avant modification ;
- proposer un plan avant les gros changements ;
- ne pas créer de complexité inutile ;
- ne pas toucher aux secrets ;
- ne pas scanner les dossiers interdits ;
- ne pas commit/push sans validation ;
- garder le MVP simple.

## 24. Non-objectifs MVP

Ne pas implémenter dans le MVP :

- IA/Ollama ;
- MCP ;
- synchronisation ;
- compte utilisateur ;
- télémétrie ;
- Homebrew tap ;
- génération parfaite depuis n’importe quel dotfiles ;
- shell Bash/Fish ;
- vrai terminal interactif exécutant des commandes.

## 25. Première étape attendue

Avant de coder :

1. inspecter le dossier projet ;
2. lire `README.md`, `AGENTS.md`, `docs/SPEC.md` ;
3. proposer une architecture ;
4. proposer les dépendances Go ;
5. proposer le modèle de données ;
6. proposer le plan MVP ;
7. proposer les commandes de test ;
8. attendre validation.
