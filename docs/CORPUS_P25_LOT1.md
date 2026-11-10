# P2.5, lot 1 — Fondamentaux Unix et Shell

Branche `content/v2-expansion-unix-01`, base `dev` à `f096eec`.
Le lot porte le corpus de 39 à 57 exercices et de neuf à douze chapitres.
La cible indicative de 100 à 120 exercices reste ouverte.

## Nouveaux chapitres

| Chapitre | Ordre | Beginner | Intermediate | Advanced |
| --- | --- | ---: | ---: | ---: |
| [Flux du shell](../data/chapters/11-flux-shell.yaml) | Guidé, révision mélangée | 2 | 3 | 1 |
| [Fichiers et permissions](../data/chapters/12-fichiers-permissions.yaml) | Indépendant | 5 | 1 | 0 |
| [Processus et environnement](../data/chapters/13-processus-environnement.yaml) | Indépendant | 2 | 4 | 0 |
| **Ajouts** | | **9** | **8** | **1** |
| **Corpus complet** | | **43** | **12** | **2** |

Le corpus complet compte 29 exercices Unix/Shell, 8 de raccourcis et 20
d’outils CLI. Les 18 ajouts sont des commandes simulées de la famille
`unix-shell` avec `tools: []`. Le filtre de niveau exact répartit les
exercices d’un même chapitre en séances distinctes. La durée dépend du niveau
et des reprises. Les cinq à dix minutes restent un objectif à vérifier en usage.

## Relecture des 18 ajouts

Chaque ligne ci-dessous indique le résultat attendu et le point vérifié avant
la saisie. Les titres et descriptions des chapitres annoncent le thème. Les
consignes donnent les noms, chemins, valeurs ou PID nécessaires à leur propre
scénario. Aucune réponse attendue n’est affichée avant la saisie.

| Chapitre | Identifiant | Niveau | Compétence et point de relecture |
| --- | --- | --- | --- |
| `11-flux-shell` | `redirect-output` | Beginner | Rediriger et remplacer. La destination est hors du dossier listé. |
| `11-flux-shell` | `append-output` | Beginner | Ajouter sans écraser. Source et destination sont distinctes. |
| `11-flux-shell` | `redirect-input` | Intermediate | Lire par l’entrée standard pour obtenir le nombre seul. |
| `11-flux-shell` | `count-matching-lines` | Intermediate | Enchaîner recherche et comptage. Compter les lignes, même si un motif apparaît plusieurs fois. |
| `11-flux-shell` | `redirect-errors` | Intermediate | Rediriger uniquement les diagnostics. Le résultat reste à l’écran. |
| `11-flux-shell` | `merge-output-errors` | Advanced | Réunir les deux flux avec une seule ouverture. L’ordre des redirections change le résultat. |
| `12-fichiers-permissions` | `copy-spaced-path` | Beginner | Copier dans un dossier existant en préservant le nom avec espace. |
| `12-fichiers-permissions` | `rename-file` | Beginner | Renommer sur place. Le nouveau nom est libre. |
| `12-fichiers-permissions` | `remove-file` | Beginner | Retirer un fichier temporaire précisément nommé. |
| `12-fichiers-permissions` | `remove-empty-directory` | Beginner | Refuser implicitement toute suppression récursive. |
| `12-fichiers-permissions` | `list-permissions` | Beginner | Choisir le format détaillé et lire les droits, le propriétaire et le groupe. |
| `12-fichiers-permissions` | `owner-execute` | Intermediate | Ajouter un seul droit. L’utilisateur possède le script du scénario. |
| `13-processus-environnement` | `show-environment-variable` | Beginner | Consulter une seule valeur. Le scénario précise qu’elle est exportée. |
| `13-processus-environnement` | `set-shell-variable` | Beginner | Définir une nouvelle variable qui reste dans le shell courant. |
| `13-processus-environnement` | `inspect-pid` | Intermediate | Sélectionner un seul processus par son PID simulé. |
| `13-processus-environnement` | `request-termination` | Intermediate | Choisir explicitement le signal d’arrêt demandé et cibler le PID fourni. |
| `13-processus-environnement` | `export-variable` | Intermediate | Définir et transmettre une valeur aux futurs processus enfants. |
| `13-processus-environnement` | `one-command-environment` | Intermediate | Modifier la locale d’un seul programme en préservant le shell courant. |

## Contrôle des 39 exercices du corpus P2.2

25 exercices reformulés et 14 conservés. La colonne « Réponse révélée »
compare la présentation avant et après relecture. Les QCM affichent leurs choix
sans désigner le bon dans la consigne. Les titres et descriptions de certains
chapitres ont été reformulés pour retirer les noms d’outils qui donnaient
une partie de la réponse. Les [clés du corpus P2.2](CORPUS_P2.md#contrat-des-identifiants-et-acquis),
les compétences et les réponses acceptées restent les mêmes.

| Chapitre | Identifiant | Statut | Motif | Réponse révélée avant → après |
| --- | --- | --- | --- | --- |
| `00-premiers-pas` | `pwd` | Reformulé | Consigne directe, sans rappel du simulateur. | Non → non |
| `00-premiers-pas` | `ls` | Reformulé | Retrait du nom de commande et de l’inventaire superflu. | Oui → non |
| `00-premiers-pas` | `ls-hidden` | Reformulé | Retrait du nom de commande. Les deux entrées spéciales distinguent le résultat attendu. | Oui → non |
| `00-premiers-pas` | `cd-docs` | Reformulé | Départ et cible relative suffisent. | Oui → non |
| `00-premiers-pas` | `cd-parent` | Reformulé | Retrait de la commande et du composant spécial. | Oui → non |
| `00-premiers-pas` | `cd-absolute` | Reformulé | Trajet explicite et contrainte du chemin absolu. | Oui → non |
| `00-premiers-pas` | `mkdir-essais` | Reformulé | Retrait de la commande et du préambule. | Oui → non |
| `00-premiers-pas` | `cat-spaces` | Reformulé | Retrait de la commande et de l’indice d’échappement. | Oui → non |
| `01-raccourcis-terminal` | `ctrl-a` | Conservé | Contexte ZLE Emacs et cible début de ligne explicites. | Non → non |
| `01-raccourcis-terminal` | `ctrl-e` | Conservé | Contexte ZLE Emacs et cible fin de ligne explicites. | Non → non |
| `01-raccourcis-terminal` | `ctrl-u-whole-line` | Conservé | La portée sur toute la ligne distingue ce widget. | Non → non |
| `01-raccourcis-terminal` | `ctrl-k` | Conservé | Délimite le texte du curseur à la fin de ligne. | Non → non |
| `01-raccourcis-terminal` | `ctrl-w` | Conservé | Désigne le mot situé avant le curseur. | Non → non |
| `01-raccourcis-terminal` | `ctrl-y` | Conservé | Le tampon ZLE est précisé. | Non → non |
| `01-raccourcis-terminal` | `ctrl-r` | Conservé | Donne le sens de la recherche dans l’historique. | Non → non |
| `01-raccourcis-terminal` | `ctrl-l` | Conservé | Décrit le rafraîchissement et le retour à la ligne en cours. | Non → non |
| `02-navigation-shell` | `zoxide-docs-command` | Reformulé | Intégration standard et correspondance unique, sans donner la fonction. | Oui → non |
| `02-navigation-shell` | `fzf-stdin-command` | Reformulé | Décrit l’entrée et le mode de sélection. | Oui → non |
| `02-navigation-shell` | `fd-find-dir` | Conservé | Le QCM distingue recherche par nom et recherche par contenu. | Non hors choix → non |
| `02-navigation-shell` | `eza-tree-command` | Reformulé | Décrit la vue en arbre sans nommer l’outil ni ses options. | Oui → non |
| `02-navigation-shell` | `yazi-command` | Reformulé | Identifie le gestionnaire par sa description et demande le lancement direct. | Oui → non |
| `02-navigation-shell` | `fd-fzf-dir` | Reformulé | Exprime le filtre de type et la sélection sans syntaxe. | Oui → non |
| `05-outils-quotidiens` | `atuin-search-command` | Reformulé | Garde l’historique indexé, le motif et la limite. | Oui → non |
| `05-outils-quotidiens` | `lazygit` | Conservé | Consigne conservée. Retrait du nom de l’outil dans la description du chapitre. | Oui hors choix → non |
| `06-recherche-fichiers-contenu` | `rg-content` | Conservé | QCM sur la recherche dans les contenus. | Non hors choix → non |
| `06-recherche-fichiers-contenu` | `fd-notes-command` | Reformulé | Conserve le motif, demande une recherche par nom. | Oui → non |
| `06-recherche-fichiers-contenu` | `fd-fzf-md` | Reformulé | Conserve l’extension et la limite des noms, retire les commandes. | Oui → non |
| `06-recherche-fichiers-contenu` | `rg-fzf-todo` | Reformulé | Garde le motif en majuscules et la sélection d’une ligne. | Oui → non |
| `06-recherche-fichiers-contenu` | `rg-fzf-docs` | Reformulé | Donne le dossier et le motif, sans imposer la syntaxe. | Oui → non |
| `07-lecture-preview` | `bat-readme` | Reformulé | Retire le nom du titre, de la description et de la consigne. La coloration précise l’objectif. | Oui → non |
| `08-micro` | `micro-open-file` | Reformulé | Demande le lancement direct de l’éditeur sans donner son exécutable. | Oui → non |
| `08-micro` | `micro-save` | Conservé | Le fichier nommé et modifié suffit à situer l’enregistrement. | Non → non |
| `08-micro` | `micro-search` | Conservé | QCM sur l’ouverture de la recherche avec les réglages standard. | Non hors choix → non |
| `08-micro` | `micro-quit` | Conservé | Un seul fichier ouvert et enregistré précise l’effet attendu. | Non → non |
| `09-markdown-glow` | `glow-readme` | Reformulé | Le rendu Markdown identifie le besoin sans donner la commande. | Oui → non |
| `09-markdown-glow` | `markdown-select-preview` | Reformulé | Conserve sélection et contraintes de chemins, retire les commandes et l’indice de substitution. | Oui → non |
| `10-workflows-dotfiles` | `dotfiles-status-shell` | Reformulé | Le format à deux colonnes remplace le nom de l’option. | Oui → non |
| `10-workflows-dotfiles` | `git-diff-stat` | Reformulé | Demande un graphique des changements par fichier sans donner l’option. | Oui → non |
| `10-workflows-dotfiles` | `git-diff-name-only` | Reformulé | Demande uniquement les noms des fichiers suivis modifiés hors index. | Oui → non |

## Références techniques

- Shell et ordre des redirections. Context7 a servi à retrouver la documentation
  des redirections. [Grammaire ZSH](https://zsh.sourceforge.io/Doc/Release/Shell-Grammar.html),
  [redirections ZSH](https://zsh.sourceforge.io/Doc/Release/Redirection.html),
  [variables ZSH](https://zsh.sourceforge.io/Doc/Release/Parameters.html) et
  [commandes internes ZSH](https://zsh.sourceforge.io/Doc/Release/Shell-Builtin-Commands.html).
- macOS. Pages de manuel locales `ls(1)`, `cat(1)`, `grep(1)`, `wc(1)`,
  `cp(1)`, `mv(1)`, `rm(1)`, `rmdir(1)`, `chmod(1)`, `ps(1)`,
  `kill(1)`, `printenv(1)`, `env(1)` et `sort(1)). Sources Apple
  [file_cmds](https://github.com/apple-oss-distributions/file_cmds),
  [text_cmds](https://github.com/apple-oss-distributions/text_cmds) et
  [ps(1)](https://github.com/apple-oss-distributions/adv_cmds/blob/main/ps/ps.1).
- GNU/Linux. `printenv` est commun à BSD et GNU, hors norme POSIX. Les colonnes
  par défaut de `ps` peuvent varier. [Manuel GNU Coreutils](https://raw.githubusercontent.com/coreutils/coreutils/master/doc/coreutils.texi),
  [GNU Grep](https://www.gnu.org/software/grep/manual/grep.html) et
  [procps-ng ps(1)](https://gitlab.com/procps-ng/procps/-/raw/master/man/ps.1).
- Outils du corpus d’origine. [Références P2.2](CORPUS_P2.md#références-éditoriales-vérifiées)
  et [raccourcis standard de Micro](https://github.com/micro-editor/micro/blob/master/runtime/help/defaultkeys.md).

## Validation

Les tests chargent le corpus embarqué, contrôlent les métadonnées, les clés
et le rendu des lignes. Les nouvelles séances sont vérifiées avec les trois
niveaux et aucun module actif, en ordre guidé et avec permutation contrôlée.
Les tests de correction incluent les variantes déclarées et des contre-exemples
sur les flux, les droits et la portée des variables. Les tests CLI soumettent
les réponses du corpus textuel sur des données isolées. Des fichiers sentinelles
vérifient qu’une saisie de redirection, acceptée ou refusée, n’écrit rien.
L’historique de progression v1 et les préférences v1 sont contrôlés sur fixtures.

Essais sans couleur sous pseudo-terminal à 80 × 24, sur macOS arm64. En CLI,
le chapitre Flux du shell a été joué en mode Apprendre aux trois niveaux.
Les scores sont de 1/2 en Beginner, 3/3 en Intermediate et 1/1 en Advanced.
En TUI, les deux exercices Beginner ont suivi l’ordre éditorial, l’exercice
Advanced a affiché son explication et les quatre exercices Intermediate de
Processus et environnement ont été présentés en ordre mélangé, une fois chacun.
`stats` et **My progress** ont distingué les acquis Beginner actifs de
l’historique. Les JSON synthétiques et les fichiers sentinelles ont été relus.

## Recette humaine et correction T07

Le 10 novembre 2026, la recette humaine finale sous Ghostty sur macOS a
validé T01 à T08. T07 est clôturé : la liste TUI et la commande `chapters`
sont validées. Le lot 1 et ses 57 exercices sont prêts pour commit sur
`content/v2-expansion-unix-01` ; ils attendent leur intégration et publication.

La deuxième recette Ghostty a validé l’entraînement CLI, la navigation TUI
et le fonctionnement sans couleur. Elle a demandé une liste TUI plus
compacte et une reprise de la présentation et des pluriels de `chapters`.

La liste TUI tient sur deux lignes par chapitre, avec le titre, l’identifiant
et le nombre d’exercices éligibles sur le total. La description complète
reste lisible sur la fiche et utilisable dans la recherche. `chapters` affiche
un en-tête et des colonnes alignées en terminal, avec « 1 exercice ».
En redirection, la sortie conserve les trois champs tabulés sans en-tête,
avec « 1 exercise » au singulier.

Les essais sous pseudo-terminal à 80 × 24 ont vérifié le défilement des
douze chapitres, le repère de sélection, la recherche par description et
la fiche complète. Les deux sorties de `chapters` ont été contrôlées.
Depuis la racine du dépôt,
`go run ./cmd/cli-drill` lance la TUI et `go run ./cmd/cli-drill chapters`
affiche la liste compilée depuis le checkout courant.
