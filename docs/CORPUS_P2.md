# Premier corpus V2 — P2.2

Corpus préparé sur `content/v2-first-corpus` depuis `ab2096c` et validé le
9 novembre 2026 après recette manuelle dans Terminal.app.
[Détail de la validation](../PROJECT_STATE.md) · [Enrichissement P2.5, lot 1](CORPUS_P25_LOT1.md)

Le corpus actif passe de 10 chapitres / 68 exercices à 9 chapitres / 39 exercices :
31 situations issues du corpus historique sont révisées et huit exercices Unix
sont créés. Le volume indicatif de 32 exercices révisés n’est pas un quota.
Les chapitres `03-alias-zsh` et `04-fonctions-zsh` sortent du corpus actif.
Les autres gardent leur identifiant même lorsque leur titre change.

`00-premiers-pas` est le premier fichier embarqué : CLI `train` sans identifiant
et TUI Start training ouvrent ce chapitre textuel, guidé, Beginner, sans outils
requis. Chaque situation fournit un point de départ indépendant ; aucune réponse
ne modifie le disque. Les outils tiers ne sont jamais supposés installés.
Les métadonnées préparent les futures sélections : aucun filtre, préférence ou
mélange n’est activé par P2.2.

## Répartition du corpus

| Chapitre / titre affiché | Famille | Exercices | Ordre éditorial |
| --- | --- | ---: | --- |
| `00-premiers-pas` — Premiers pas Unix | `unix-shell` | 8 | `guided` |
| `01-raccourcis-terminal` — Raccourcis ZSH — mode Emacs | `keyboard-shortcuts` | 8 | `independent` |
| `02-navigation-shell` — Navigation avec les outils CLI | `modern-cli-macos` | 6 | `independent` |
| `05-outils-quotidiens` — Historique et interface Git | `modern-cli-macos` | 2 | `independent` |
| `06-recherche-fichiers-contenu` — Rechercher et sélectionner | `modern-cli-macos` | 5 | `independent` |
| `07-lecture-preview` — Lire un fichier avec Bat | `modern-cli-macos` | 1 | `independent` |
| `08-micro` — Éditer dans Micro | `modern-cli-macos` | 4 | `independent` |
| `09-markdown-glow` — Lire du Markdown avec Glow | `modern-cli-macos` | 2 | `independent` |
| `10-workflows-dotfiles` — Git essentiel | `unix-shell` | 3 | `independent` |

À la livraison de P2.2, l’ordre `guided` ou `independent` exprimait une
intention éditoriale. Le moteur parcourait le YAML dans l’ordre.
Les consignes sont repliées à 70 caractères et les explications à 60, pour
préserver la lisibilité à 80 colonnes avec les préfixes du rendu existant.
Les durées de cinq à dix minutes restent un objectif à confirmer par l’usage ;
les chapitres très courts Bat et Historique/Git sont volontairement des mini-leçons,
sans ajout de remplissage pour atteindre une durée annoncée.

| Répartition | Effectifs |
| --- | --- |
| Familles | Unix/Shell : 11 ; raccourcis : 8 ; outils CLI : 20 |
| Difficultés | Beginner : 34 ; Intermediate : 4 ; Advanced : 1 |
| Formats | Commandes simulées : 25 ; captures réelles : 8 ; QCM : 4 ; réponses textuelles de raccourcis Micro : 2 |

Git est classé dans Unix/Shell pour son apprentissage général, avec `tools: [git]`.
Le socle Unix de huit exercices a `tools: []`. Les widgets ZLE standards n’exigent
aucun module tiers. Les choix erronés d’un QCM ne sont pas des modules requis :
seuls les outils effectivement enseignés figurent dans `tools`.

## Contrat des identifiants et acquis

Les clés restent `chapter-id/item-id`. Vingt clés sont conservées pour une
compétence identique, 48 sortent du corpus actif et 19 sont nouvelles. Aucune
clé retirée n’est recyclée. Un titre ou une famille ne sert pas de clé.
La table couvre chaque ancien exercice et chaque création : les identifiants
ci-dessous sont relatifs au chapitre de la première colonne.

| Chapitre | Clés conservées | Clés retirées | Nouvelles clés |
| --- | --- | --- | --- |
| `00-premiers-pas` | — | — | `pwd`, `ls`, `ls-hidden`, `cd-docs`, `cd-parent`, `cd-absolute`, `mkdir-essais`, `cat-spaces` |
| `01-raccourcis-terminal` | `ctrl-a`, `ctrl-e`, `ctrl-k`, `ctrl-w`, `ctrl-y`, `ctrl-r`, `ctrl-l` | `ctrl-u` | `ctrl-u-whole-line` |
| `02-navigation-shell` | `fd-find-dir`, `fd-fzf-dir` | `zoxide-jump`, `fzf-pick`, `eza-tree`, `yazi-open` | `zoxide-docs-command`, `fzf-stdin-command`, `eza-tree-command`, `yazi-command` |
| `03-alias-zsh` | — | `cgs`, `cgd`, `cgdn`, `lg`, `dtf`, `mc`, `ezac` | — |
| `04-fonctions-zsh` | — | `y`, `mdpreview`, `funclist`, `clip`, `slug`, `port`, `zedstdin` | — |
| `05-outils-quotidiens` | `lazygit` | `atuin`, `yazi`, `zoxide`, `bat`, `eza`, `glow` | `atuin-search-command` |
| `06-recherche-fichiers-contenu` | `rg-content`, `fd-fzf-md`, `rg-fzf-todo` | `fd-files`, `fzf-filter`, `bat-preview`, `rg-bat` | `fd-notes-command`, `rg-fzf-docs` |
| `07-lecture-preview` | `bat-readme` | `bat-file`, `glow-md`, `markedit-app`, `marked2-preview`, `mdpreview-function` | — |
| `08-micro` | `micro-open-file`, `micro-save`, `micro-search`, `micro-quit` | `micro-open`, `mc-alias` | — |
| `09-markdown-glow` | `glow-readme` | `glow-preview`, `markdownlint`, `vale`, `mdpreview`, `markdown-pipeline`, `markdown-generate-preview` | `markdown-select-preview` |
| `10-workflows-dotfiles` | `dotfiles-status-shell` | `static-scan`, `git-status-short`, `check-dotfiles`, `update-dotfiles`, `generate-tools-inventory`, `gitleaks-before-commit` | `git-diff-stat`, `git-diff-name-only` |

Les retraits sans remplacement concernent alias/fonctions personnels, anciennes
opérations de scan/génération et inventaire, applications graphiques et questions
de reconnaissance redondantes. Les exercices Git concrets et les commandes
Bat/Glow sont conservés plutôt que les anciennes questions de noms répétées.
`Marked 2` et MarkEdit sont retirés ; le correcteur P1 n’est pas assoupli.

Les remplacements de compétence justifiant une nouvelle clé sont explicites :

| Ancienne clé | Nouvelle clé | Motif |
| --- | --- | --- |
| `01-raccourcis-terminal/ctrl-u` | `01-raccourcis-terminal/ctrl-u-whole-line` | Le binding Emacs ZLE est kill-whole-line, pas backward-kill-line. |
| `02-navigation-shell/zoxide-jump` | `02-navigation-shell/zoxide-docs-command` | Saisir z docs avec intégration et chemin connu, au lieu de nommer un outil. |
| `02-navigation-shell/fzf-pick` | `02-navigation-shell/fzf-stdin-command` | Saisir la commande de sélection sur une entrée standard explicitement fournie. |
| `02-navigation-shell/eza-tree` | `02-navigation-shell/eza-tree-command` | Construire une commande d’arborescence, sans alias ni police à icônes. |
| `02-navigation-shell/yazi-open` | `02-navigation-shell/yazi-command` | Lancer yazi directement ; y reste un wrapper facultatif, hors réponse cible. |
| `05-outils-quotidiens/atuin` | `05-outils-quotidiens/atuin-search-command` | Construire une recherche bornée dans un historique fourni, au lieu de nommer Atuin/Ctrl+R. |
| `06-recherche-fichiers-contenu/fd-files` | `06-recherche-fichiers-contenu/fd-notes-command` | Saisir un motif de recherche de noms. |
| `06-recherche-fichiers-contenu/rg-bat` | `06-recherche-fichiers-contenu/rg-fzf-docs` | Supprimer la promesse d’ouverture Bat : le pipeline sélectionne une ligne et ne lance jamais Bat. |
| `09-markdown-glow/markdown-pipeline` | `09-markdown-glow/markdown-select-preview` | Évaluer explicitement une substitution protégée pour les espaces, avec sélection non vide et limites de noms déclarées. |
| `03-alias-zsh/cgd` | `10-workflows-dotfiles/git-diff-stat` | Commande Git portable dans un dépôt, au lieu de rappeler un alias personnel. |
| `03-alias-zsh/cgdn` | `10-workflows-dotfiles/git-diff-name-only` | Commande Git complète et distinction arbre de travail/index. |

Les suppressions des variantes personnelles `mc`, `lg` ou `cgs` des réponses
conservées ne changent pas la compétence cible Micro/Git ; aucune migration
ou nouvelle attribution d’acquis n’en découle. Les variantes de syntaxe utiles
restent explicites. L’équivalence générale entre commandes n’est pas déduite.

Aucun JSON utilisateur n’est déplacé, purgé ou réécrit par migration. Les acquis
retirés et scores cumulés restent stockés en v1 ; un entraînement ajoute seulement
ses réponses selon P1.2. Les nouvelles clés n’héritent d’aucun ancien succès.
Les dix YAML de la base sont conservés octet pour octet dans
`testdata/legacy-chapters/` pour les tests de compatibilité et de correction P1.
Ils ne sont pas embarqués, importés ou exposés à l’entraînement.

## Références éditoriales vérifiées

Consultées le 8 novembre 2026. Context7 n’a pas fourni de correspondance officielle
pertinente pour la recherche ZLE et la requête regroupant les outils ; les sources
primaires ci-dessous ont donc servi de référence. Les préconditions sont dans les
consignes et explications, sans champ YAML supplémentaire.

- Unix/macOS : [ls](https://github.com/apple-oss-distributions/file_cmds/blob/main/ls/ls.1),
  [mkdir](https://github.com/apple-oss-distributions/file_cmds/blob/main/mkdir/mkdir.1),
  [cat](https://github.com/apple-oss-distributions/text_cmds/blob/main/cat/cat.1),
  [cd et pwd dans ZSH](https://zsh.sourceforge.io/Doc/Release/Shell-Builtin-Commands.html).
  Le chapitre d’introduction évite les options longues spécifiques à GNU.
- [Widgets ZLE](https://zsh.sourceforge.io/Doc/Release/Zsh-Line-Editor.html) et
  [table des bindings standards](https://github.com/zsh-users/zsh/blob/master/Src/Zle/zle_bindings.c).
  La capture d’une touche ne démontre pas son effet dans un shell.
- [fd](https://github.com/sharkdp/fd#how-to-use),
  [ripgrep](https://github.com/BurntSushi/ripgrep/blob/master/GUIDE.md),
  [fzf](https://github.com/junegunn/fzf#usage) : noms/contenus, exclusions et entrée/sortie standard.
- [eza](https://github.com/eza-community/eza/blob/main/man/eza.1.md),
  [Zoxide](https://github.com/ajeetdsouza/zoxide#readme),
  [Yazi](https://yazi-rs.github.io/docs/quick-start/) : --tree, intégration z et wrapper y facultatif.
- [Bat](https://github.com/sharkdp/bat#how-to-use),
  [Glow](https://github.com/charmbracelet/glow#readme) et
  [substitution ZSH](https://zsh.sourceforge.io/Doc/Release/Expansion.html#Command-Substitution).
  La sélection Glow suppose un chemin choisi non vide, sans retour à la ligne
  ni tiret initial ; elle ne gère pas l’annulation.
- [Micro](https://github.com/zyedidia/micro/blob/master/runtime/help/defaultkeys.md),
  [Atuin search](https://docs.atuin.sh/latest/reference/search/),
  [Lazygit](https://github.com/jesseduffield/lazygit#readme),
  [git status](https://git-scm.com/docs/git-status) et [git diff](https://git-scm.com/docs/git-diff).

Les références documentaires vérifient la syntaxe et le sens ; elles ne prouvent
ni l’installation des outils ni une durée de séance réelle. Les réponses sont
comparées comme texte ; aucun de ces outils n’est lancé par l’exercice.
