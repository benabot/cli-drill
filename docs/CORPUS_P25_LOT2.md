# P2.5, lot 2 — Outils CLI modernes

Branche `content/v2-expansion-cli-02`, depuis `dev` à `b809151`.
Le lot ajoute 18 exercices dans six chapitres existants. Le corpus atteint
75 exercices dans 12 chapitres. La cible indicative de 100 à 120 reste ouverte.
Le lot 2 est terminé et validé. La recette humaine finale sous Ghostty sur
macOS a réussi le 10 novembre 2026. Le lot est prêt pour commit sur sa
branche. L'intégration et la publication restent à faire.

| Répartition | Beginner | Intermediate | Advanced | Total |
| --- | ---: | ---: | ---: | ---: |
| Ajouts | 1 | 15 | 2 | 18 |
| Corpus complet | 44 | 27 | 4 | 75 |

Les nouveaux exercices sont indépendants et classés dans la famille
`modern-cli-macos`. Aucun module ne nécessite une installation pour
répondre. Les raccourcis Micro sont des réponses textuelles.

## Relecture des nouveaux exercices

Chaque consigne a été relue pour son autonomie, sa clarté, l'absence de
commande complète révélée et la pertinence de l'explication. Les titres des
chapitres et les noms d'outils précisent le contexte. La table indique les
nouvelles clés et la raison du niveau choisi.

| Chapitre et identifiant | Modules | Niveau | Compétence et justification |
| --- | --- | --- | --- |
| `02-navigation-shell/eza-tree-depth` | eza | Intermediate | Combiner arbre et profondeur. |
| `02-navigation-shell/eza-hidden-long` | eza | Intermediate | Combiner détails et fichiers cachés, sans . ni ... |
| `02-navigation-shell/zoxide-list-matches` | zoxide | Intermediate | Interroger la base sans changer de répertoire. |
| `05-outils-quotidiens/atuin-failed-search` | atuin | Intermediate | Filtrer une recherche d'historique par code de sortie. |
| `06-recherche-fichiers-contenu/rg-numbered` | rg | Intermediate | Demander les numéros de ligne même en sortie redirigée. |
| `06-recherche-fichiers-contenu/rg-matching-files` | rg | Intermediate | Obtenir les chemins des fichiers correspondants. |
| `06-recherche-fichiers-contenu/rg-fixed-dot` | rg | Intermediate | Chercher un point littéral plutôt qu'une expression régulière. |
| `06-recherche-fichiers-contenu/rg-go-glob` | rg | Intermediate | Limiter la recherche au glob d'extension protégé du shell. |
| `06-recherche-fichiers-contenu/fd-go-files` | fd | Intermediate | Combiner type fichier et extension. |
| `06-recherche-fichiers-contenu/fd-hidden-notes` | fd | Intermediate | Inclure les chemins cachés tout en conservant les exclusions. |
| `06-recherche-fichiers-contenu/rg-files-select` | rg, fzf, bat | Advanced | Choisir un fichier par son contenu et protéger le chemin avec espaces. Sélection non vide, noms sans saut de ligne ni tiret initial. |
| `07-lecture-preview/bat-number-lines` | bat | Intermediate | Choisir les numéros sans les autres décorations. |
| `07-lecture-preview/bat-line-range` | bat | Intermediate | Restreindre la lecture à une plage inclusive. |
| `07-lecture-preview/bat-json-stdin` | bat | Advanced | Interpréter un flux sans extension, forcer les couleurs ANSI en redirection et désactiver le pager. |
| `08-micro/micro-undo` | micro | Beginner | Annuler une modification avec le raccourci standard. |
| `08-micro/micro-line-numbers` | micro | Intermediate | Basculer les numéros de ligne. |
| `09-markdown-glow/glow-pager` | glow | Intermediate | Parcourir un rendu avec le pager configuré. |
| `09-markdown-glow/glow-width` | glow | Intermediate | Fixer la largeur de retour à la ligne du rendu. |

## Références officielles

- [Guide ripgrep](https://github.com/BurntSushi/ripgrep/blob/master/GUIDE.md)
  et [référence des options](https://github.com/BurntSushi/ripgrep/blob/master/crates/core/flags/defs.rs).
  Context7 a été consulté pour les filtres, les chemins et la recherche littérale.
- [fd](https://github.com/sharkdp/fd#how-to-use),
  [fzf](https://github.com/junegunn/fzf#usage) et
  [manuel fzf](https://github.com/junegunn/fzf/blob/master/man/man1/fzf.1).
- [bat](https://github.com/sharkdp/bat#how-to-use). Context7 et l'aide locale
  ont confirmé la plage de lignes, la syntaxe, le pager et la sortie couleur.
- [eza](https://github.com/eza-community/eza/blob/main/man/eza.1.md),
  [zoxide](https://github.com/ajeetdsouza/zoxide#readme) et
  [recherche Atuin](https://docs.atuin.sh/latest/reference/search/).
  Les aides locales ont précisé la sélection de la base et le code de sortie.
- [Raccourcis Micro](https://github.com/micro-editor/micro/blob/master/runtime/help/defaultkeys.md)
  et [Glow](https://github.com/charmbracelet/glow#usage).

Les chemins avec espaces sont conservés par les guillemets du scénario avancé.
Ce scénario suppose une sélection confirmée et ne traite pas l'annulation.
Les sorties colorées forcées comprennent des séquences ANSI. Les réglages des
outils réels et les programmes de pagination peuvent varier.

## Vérifications

Les tests de chargement valident 75 exercices et la conservation des clés
historiques. Les 18 nouvelles clés ont des tests de correction, de niveau et
de modules. Chaque cas vérifie l'absence d'éligibilité sans module, la présence
avec les modules requis et une permutation contrôlée du chapitre indépendant.
Le cas à trois modules est aussi vérifié en omettant chacun d'eux.
Les tests existants soumettent chaque réponse textuelle avec une progression
v1 synthétique et vérifient la conservation des résultats retirés.

`go test -count=1 ./...` et `go vet ./...` ont réussi sous macOS arm64.
Le rendu des 18 ajouts avant et après les deux verdicts est testé à 80 × 24.
Des essais CLI et TUI sous pseudo-terminal à 80 × 24 ont soumis les deux
scénarios Advanced. Une progression v1 et des préférences synthétiques ont
été relues et les commandes saisies n'ont pas été exécutées.

## Recette finale sous Ghostty

Le 10 novembre 2026, les sept tests fonctionnels ont réussi. La recette UX
a validé Ma progression, la fiche de chapitre et la lisibilité sous Ghostty.
Aucune anomalie restante n’a été signalée.

Ma progression distingue les acquis éligibles de l’historique global.
La fiche sépare les informations du chapitre, les choix de séance et la
progression.

Les numéros 03 et 04 correspondent à deux chapitres historiques retirés.
Les identifiants existants servent de clés aux acquis du JSON v1.

Dans tmux, quitter la capture clavier a reproduit « Session interrompue. »
au-dessus de l’accueil. Le parcours CLI écrit cette ligne avant le retour
dans Bubble Tea. Elle persiste dans la sortie du terminal.

Les deux vues ont été vérifiées à 80 × 24 sans couleur et après
redimensionnement à 64 × 24. Les tests couvrent les compteurs à vide,
l’effet du niveau et des modules, l’historique et les modes de séance.
