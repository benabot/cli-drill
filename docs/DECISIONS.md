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

## 2026-08-08 — TUI Text Input Shortcuts

TUI footers are selected by interaction context. While typing a free-text
answer, global single-key shortcuts are disabled so every printable key reaches
the text input. `Enter` submits, `Esc` returns to the chapter detail, and
`Ctrl+C` quits the TUI.

List-based screens hide the built-in Bubbles help when cli-drill renders its
centralized contextual footer. Home, Chapters and Directory therefore show one
footer source without changing their existing navigation shortcuts.

## 2026-08-09 — TUI Chapter Completion

Text chapter completion uses a session-only summary. It shows the score and
available missed answers without changing persistent scoring. Retry restarts
the chapter and is offered only when the completed session has misses. `g`
returns Home from that summary; a session without misses returns Home directly.

<a id="2026-11-07--cli-drill-v2--corpus-integre-et-retrait-des-dotfiles"></a>

## 2026-11-07 — CLI-Drill V2 : corpus intégré et retrait des dotfiles

**Statut : décision produit validée ; migration fonctionnelle non réalisée.**
La présente livraison modifie la conception et le backlog, sans changement Go,
YAML pédagogique, dépendance, README public ou donnée utilisateur.

### Contexte et motivation

La V1 transforme des dotfiles en catalogue puis en leçons. Ce pipeline ajoute
configuration, parsing, filtrage de chemins, choix de sources et génération,
sans garantir la qualité pédagogique : un nom de fonction ou une mention
d'outil ne suffisent pas à produire une question discriminante.
Il lie aussi les leçons à des alias et workflows personnels.

L'objectif devient **CLI-Drill — Learn the terminal by doing.** La valeur
principale est l'entraînement : saisir des commandes, pratiquer les raccourcis,
comprendre ses erreurs et mémoriser par des sessions courtes répétables.
Un corpus intégré écrit et vérifié permet de concentrer l'effort sur cette valeur,
de commencer sans configuration et de fonctionner sans outils tiers installés.

### Décision

- Conserver Go, Cobra, Bubble Tea, Bubbles, Lip Gloss, l'identité graphique,
  les chapitres YAML intégrés, les exercices et la progression JSON locale.
- Maintenir comme activités centrales la saisie de commandes simulées sans
  exécution et la capture réelle des raccourcis compatibles avec le terminal.
- Organiser le corpus en Unix & Shell Fundamentals, Keyboard Shortcuts et
  Modern CLI & macOS, avec contexte de portabilité et préconditions explicites.
- Retirer scan des dotfiles, détection d'outils installés, parsing utilisateur,
  génération, catalogue, TOML dotfiles et leurs commandes/écrans/options.
- Conserver `chapters`, `train`, `stats`, `reset --yes` et `--progress`.
- Utiliser exclusivement le corpus intégré ; arrêter le chargement automatique
  des anciens YAML utilisateur, sans ajout d'import personnalisé dans ce MVP.
- Conserver les moteurs textuel et clavier dédié, les packages nécessaires et
  les composants réutilisables ; aucun nouveau framework ou moteur universel.

### Données et compatibilité

Préserver intégralement les anciennes données utilisateur : progression JSON,
chapitres YAML et configuration TOML. Le retrait de leur usage n'autorise aucun
nettoyage, déplacement, conversion ou réécriture automatique. Le premier
lancement V2 et les consultations ne déclenchent aucune migration écrite.

La progression garde son chemin, son schéma version 1, ses champs et ses clés
`chapter-id/item-id`. Les entraînements peuvent l'enrichir, sans purger les
anciens résultats. Conserver un identifiant si l'apprentissage reste identique ;
utiliser une nouvelle clé pour un contenu différent. Les indicateurs du corpus
actif filtrent les identifiants présents sans modifier le stockage historique.

Une progression invalide ou non prise en charge ne doit jamais être remplacée
silencieusement par un état vide. Prévoir sauvegarde atomique et erreurs visibles.
`reset --yes` demeure une suppression explicitement demandée, jamais une étape
implicite de mise à niveau.

### Terminal et sécurité

L'introspection abandonnée désigne le scan des dotfiles et la détection des
outils installés. Les vérifications techniques de TUI, TTY, dimensions,
couleurs et capture clavier restent nécessaires. Conserver notamment les usages
appropriés de `x/term`, `TERM` et `NO_COLOR` sans inventorier l'environnement
personnel ni afficher de variables sensibles.

Le premier chapitre cible sera textuel, sans configuration et sans capture
brute imposée. Vérifier les capacités avant activation de la TUI ou de la capture,
traiter leurs erreurs et guider vers une CLI textuelle sans couleur lorsque
nécessaire. Ne pas promettre la capture de toute combinaison macOS et ne pas
créditer une capture indisponible.

Les autres interdictions de sécurité demeurent : pas de lecture de secrets,
exécution des commandes entraînées, modification des dotfiles, télémétrie,
réseau applicatif nécessaire, IA ou MCP.

### Compromis et décisions remplacées

Accepter la perte de personnalisation automatique et la rupture des commandes
retirées ; conserver les données ne signifie pas conserver ces fonctionnalités.
Le corpus demande une maintenance éditoriale et doit distinguer portable,
BSD/macOS, GNU/Linux et outils tiers. Viser cinq à dix minutes par chapitre,
sans minuteur, aléatoire, préférences multiples ou reprise persistante au curseur.

Cette décision remplace, pour les travaux actifs, le pipeline dotfiles, le
catalogue, le TOML et les orientations de génération des décisions du
6 juin 2026, dont Catalog Cleanup P1, CLI UX P2.1 et la référence aux éléments
personnels dans Default Training Chapters P2.2. Editable Chapters reste valable
pour l'édition des sources intégrées, sans génération ni surcharge utilisateur.
La restriction ZSH comme format de scan disparaît ; ZSH sous macOS reste un
contexte pédagogique de référence, sans architecture multi-shell à anticiper.
La persistance textuelle TUI attendue en V2 remplace la limite historique de
bilan de session seul décrite le 9 août 2026. Les décisions de module racine,
fonctionnement local, saisie contextuelle et moteur clavier dédié restent utiles.
Aucun texte historique ci-dessus n'est supprimé ni réécrit.

La [spécification V2](SPEC.md), l'[architecture cible](ARCHITECTURE.md),
l'[état actuel](../PROJECT_STATE.md) et le [backlog](../TODO.md) distinguent
la décision de sa réalisation. Le README sera actualisé après migration.


<a id="2026-11-08--p20--evolution-pedagogique"></a>

## 2026-11-08 — P2.0 : évolution de la feuille de route pédagogique

**Statut : décisions produit fournies après validation de l'audit pédagogique
des 68 exercices ; cadrage documentaire approuvé le 8 novembre 2026,
implémentation P2 non réalisée.**
La demande du 8 novembre 2026 constitue la référence de cette validation ;
cette mission ne refait ni l'audit pédagogique ni une revue exhaustive de P1.
La base vérifiée est `dev` à `ef9b0a5`, avec P1 terminée et validée.
La validation de finalisation du 8 novembre 2026 autorise le commit documentaire,
l'intégration par fast-forward dans `dev` et la publication de `origin/dev`, en
conservant `docs/p2-roadmap`. Elle ne vaut pas validation fonctionnelle P2 et
n'autorise aucune implémentation P2.1.

### Ce qui change par rapport à la décision du 7 novembre

- Le périmètre explicite est commandes Unix, shells, raccourcis du terminal
  et outils CLI. La famille Modern CLI & macOS ne couvre pas les logiciels
  graphiques pour leur seule présence sur macOS ; leurs questions héritées
  devront sortir du corpus actif. Le cas `Marked 2` reporté par P1.3 est traité
  éditorialement, sans assouplissement global des réponses de workflow.
- La personnalisation manuelle devient un objectif P2 : socle Unix toujours
  disponible, modules tiers activables individuellement, même sans installation.
  Aucun scan, lecture de dotfiles, ancien scanner ou catalogue ne revient.
- Le corpus court est un premier lot, suivi d'une cible indicative de 100 à
  120 exercices au total dans les trois familles. Commandes saisies, scénarios
  autonomes et captures compatibles priment sur la répétition de noms d'outils.
- Beginner, Intermediate et Advanced qualifient chaque exercice. Le niveau
  d'une séance peut être choisi au lancement, notamment depuis une fiche chapitre.
  Les YAML historiques sans difficulté restent lisibles sans conversion.
- L'ordre aléatoire devient prévu pour chaque nouvelle séance ordinaire et
  les révisions. L'apprentissage initial suit un ordre guidé lorsqu'il est
  nécessaire. Filtrer avant de mélanger, conserver l'ordre pendant la séance
  et présenter une fois chaque exercice éligible dans le parcours initial.
  Les retries volontaires restent de nouvelles tentatives légitimes selon P1.2.

Ces choix remplacent les exclusions de l'aléatoire et des préférences
pédagogiques minimales du 7 novembre. Ils ne rétablissent ni la personnalisation
automatique ni les autres préférences différées. Les décisions historiques
ci-dessus sont conservées comme historique, y compris leurs anciens libellés P2 V1.

### Choix de conception minimaux retenus pour le cadrage

La demande fixe les objectifs ; les précisions suivantes constituent le contrat
minimal retenu dans P2.0, sans prétendre décrire le logiciel actuel :

- YAML additif : `difficulty` et `tools` sur l'exercice, `family` et
  `learning_order` sur le chapitre. Niveau absent : Beginner en mémoire ;
  ordre absent : guidé par prudence. Les sources révisées déclareront ces valeurs.
  Les préconditions de portabilité restent textuelles, sans nouveau schéma complexe.
- Par défaut, aucun module tiers activé et niveau Beginner. Le filtre de niveau
  est exact ; un exercice multi-outils exige tous ses modules. Aucun résultat
  n'entraîne une activation ou un changement automatique de niveau.
- Préférences locales facultatives dans un futur `preferences.json`, distinct
  du JSON de progression : version propre, outils activés et niveau par défaut.
  Écriture uniquement sur enregistrement explicite, défauts en mémoire sinon,
  erreurs visibles et aucune réparation automatique d'un fichier invalide.
  Les chemins et cas d'erreur sont définis dans la spécification.
- Ordre et options ponctuelles de séance en mémoire ; aucun marqueur de première
  séance, graine persistante ou curseur. Sur un chapitre guidé, proposer Apprendre
  par défaut et laisser choisir Réviser explicitement. Deux permutations
  successives peuvent être identiques ; le mélange n'est pas une promesse d'alternance.

### Garanties maintenues et limites

Le JSON de progression reste en version 1, mêmes champs, chemins, clés et
compteurs, avec `progress.RecordAnswer` et les protections de stockage P1.
Ni module, ni difficulté, ni ordre n'est stocké dans la progression. Les acquis
sont projetés sur le corpus sélectionné sans effacer les résultats historiques.
Déplacer un exercice vers une autre famille n'autorise pas à changer sa clé
`chapter-id/item-id`. Un nouvel apprentissage reçoit une nouvelle clé.
Les anciens JSON, YAML et TOML ne sont ni purgés, ni déplacés, ni convertis
par migration ; les YAML utilisateur ne sont toujours pas chargés automatiquement.

Conserver les deux moteurs, les six packages, Go/Cobra/Bubble Tea, l'identité
visuelle, le fonctionnement local hors ligne, les contrôles terminal et
l'interdiction d'exécuter les réponses. La capture indisponible ne donne aucun
acquis fictif. Aucun moteur adaptatif, répétition espacée, base de données,
framework de configuration ou nouvelle dépendance spéculative n'est prévu.
Les limites P1 de concurrence et de durabilité restent inchangées.

La [spécification](SPEC.md) est la référence des contrats ; le
[backlog](../TODO.md#p2--feuille-de-route-pedagogique) définit l'ordre des
sous-phases et leurs critères. L'[architecture](ARCHITECTURE.md) précise leur
répartition minimale ; l'[état courant](../PROJECT_STATE.md) distingue cadrage
P2.0 et fonctions réellement disponibles. Aucun Go, YAML, README ou fichier
utilisateur/configuration n'est modifié par cette mission.
