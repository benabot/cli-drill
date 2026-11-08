# AGENTS.md

Consignes communes pour les agents travaillant sur ce dépôt.

## Lecture obligatoire

Avant toute modification, lire :

- [README.md](README.md) pour le logiciel disponible ;
- [AGENTS.md](AGENTS.md), [PROJECT_STATE.md](PROJECT_STATE.md) et [TODO.md](TODO.md) ;
- [docs/SPEC.md](docs/SPEC.md), [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)
  et [docs/DECISIONS.md](docs/DECISIONS.md) ;
- [.codex/skills/cli-drill-spec/SKILL.md](.codex/skills/cli-drill-spec/SKILL.md).

Lire les instructions plus proches et éventuels fichiers d'override applicables.
Examiner le code courant avant de décrire un comportement implémenté.

## Direction produit et transition V2

**CLI-Drill — Learn the terminal by doing.** La V2 vise l'entraînement par
corpus intégré vérifié, saisie de commandes simulées et capture réelle des
raccourcis compatibles. La [spécification](docs/SPEC.md) définit le périmètre.

La migration V2 est partiellement réalisée. Consulter [PROJECT_STATE.md](PROJECT_STATE.md)
pour les travaux intégrés et ceux qui attendent une validation humaine ; la
présente clarification documentaire ne change aucun de ces statuts. Distinguer
code existant, choix futurs et preuves réellement obtenues. Le README conserve
la présentation du MVP antérieur, signalée comme historique ; ses exemples de
scan/génération ne décrivent plus le checkout courant. Sa refonte complète reste
distincte de cette clarification.

Ne pas ajouter ni réintroduire scan des dotfiles, détection des outils installés,
parsing utilisateur, génération ou catalogue. Ne pas recréer de configuration
dotfiles, importer automatiquement d'anciens YAML ou anticiper IA/MCP/multi-shell.
Les suppressions Go se feront uniquement dans les tâches explicitement autorisées.

L'abandon de l'introspection ne concerne pas les vérifications techniques du
terminal. Conserver les contrôles nécessaires à la TUI, aux couleurs et à la
capture, et traiter les capacités indisponibles sans les supposer. Le premier
entraînement cible est textuel et immédiatement accessible sans configuration.
Conserver l'identité graphique et réutiliser les composants existants.

P2 est cadrée par la [décision du 8 novembre 2026](docs/DECISIONS.md#2026-11-08--p20--evolution-pedagogique)
et la spécification : commandes Unix, shells, raccourcis du terminal et outils CLI,
sans leçons sur des logiciels graphiques pour leur seule présence sur macOS.
La sélection manuelle des modules tiers, les niveaux par exercice et le mélange
des séances sont désormais prévus, avec apprentissage guidé si nécessaire.
Leurs préférences minimales restent distinctes de la progression JSON v1 ;
elles n'autorisent ni réintroduction du scanner/catalogue, ni import automatique,
ni implémentation hors d'une tâche explicitement consacrée à la sous-phase.

<a id="regles-de-securite"></a>

## Règles de sécurité

- Ne pas lire de secrets ni afficher de variables sensibles.
- Ne pas scanner ou lire `~/.config/secrets`, `~/.ssh`, `~/.gnupg`,
  `~/.config/gh/hosts.yml` ou `~/.config/zed/settings.json`.
- Ne pas exécuter d'alias, fonctions shell utilisateur ou commandes saisies
  pour les exercices. Le simulateur compare du texte.
- Ne pas lancer Docker, Colima, Ollama ou n8n.
- Ne pas modifier un repo dotfiles utilisateur ni `.zshrc`.
- Préserver intégralement les anciennes données utilisateur JSON, YAML et TOML :
  aucune purge, conversion, réécriture ou déplacement automatiques par migration.
  Les entraînements enrichissent la progression sans supprimer ses résultats
  historiques ; `reset --yes` ne doit jamais être une étape implicite.
- Ne pas lire ni modifier les vraies données utilisateur pour tester la migration ;
  utiliser des fixtures synthétiques isolées.
- Ne pas ajouter réseau applicatif nécessaire, télémétrie, compte, sync, IA ou MCP.
- Vérifier le terminal avant usage de la TUI ou du mode brut et restaurer son état
  après capture, interruption ou erreur traitée ; prévoir le mode sans couleur.
- Garder le MVP simple et local. Ne pas ajouter de dépendance ou abstraction
  spéculative pour remplacer les fonctionnalités retirées.

## Structure du projet

Le module Go est à la racine ; le point d'entrée est `cmd/cli-drill/`.
Le code interne vit dans `internal/`, les chapitres intégrés et l'embarquement
dans `data/`, les fixtures dans `testdata/` et la documentation dans `docs/`.
Le périmètre conserve `app`, `chapter`, `exercise`, `progress`, `tui` et `xdg`.
L'état des retraits et du stockage est décrit dans `PROJECT_STATE.md` ; ne pas
planifier à nouveau un retrait déjà réalisé ni déclarer toute la V2 terminée.
Ne pas créer de répertoires réservés à des intégrations futures.

## État, backlog et documentation

`PROJECT_STATE.md` décrit l'implémentation et la transition ; `TODO.md` porte
le backlog actif. La spécification porte le produit, l'architecture les
frontières techniques et `docs/DECISIONS.md` les décisions datées.
Préserver les anciennes décisions ; ajouter la décision qui les remplace.
Éviter les copies concurrentes ; préférer les liens vers les documents de référence.

Ne pas utiliser `docs/llm/scratch/`, `docs/llm/runs/`,
`docs/llm/transcripts/` ou `docs/llm/context/` comme mémoire durable du projet.
Ces dossiers restent locaux et ignorés par Git. Ne pas publier d'historiques
d'agents sans autorisation explicite.

## Politique de changement et Git

Avant un changement important : résumer l'intention, lister les fichiers,
expliquer les risques et proposer les validations. Préserver les changements
existants ; ne pas faire de reset destructif ou de stash automatique.

Pour les futurs changements V2 : partir de `dev` propre et actualisée par
fast-forward uniquement ; une branche, un objectif cohérent, une PR vers `dev`.
Utiliser `refactor/<scope>`, `fix/<scope>`, `content/<scope>` ou `docs/<scope>`
selon la tâche. Ne pas cibler directement `main`.

Ne pas créer de commit, push, merge ou PR sans validation explicite correspondante.
La validation d'un plan ne vaut pas autorisation d'intégrer les modifications.
Une mission documentaire n'autorise pas de refactorisation Go ni de changement
des chapitres, dépendances ou du README hors du périmètre validé.

Après modification, vérifier le résultat et indiquer limites et validations :

- Documentation seule : cohérence, liens, diff, `git diff --check` et périmètre.
- Code Go : tests pertinents, `go test ./...`, `go vet ./...`, contrôle du
  formatage et aide CLI concernée ; ne pas confondre tests et validation terminal réelle.
- Ne pas lancer un formateur réécrivant les sources ou `go mod tidy` dans une
  tâche exclusivement documentaire.
