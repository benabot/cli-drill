# P2.5, vague 1 — Repères, fichiers et droits

10 novembre 2026. Branche locale `content/p25-wave1-paths-files`, créée à
partir de `dev` propre à `d945059`, également pointe de `origin/dev` et
de la branche distante lors de la vérification initiale.

La vague couvre intégralement les deux chapitres du [référentiel](CURRICULUM.md).
Les 46 ajouts portent le corpus de **75 à 121 exercices dans 12 chapitres**.
Les 75 exercices historiques, leurs réponses et leurs clés restent identiques.

| Chapitre | Avant B / I / A | Ajouts B / I / A | Après B / I / A | Total | Séance |
| --- | ---: | ---: | ---: | ---: | --- |
| [`00-premiers-pas`](../data/chapters/00-premiers-pas.yaml) | 8 / 0 / 0 | 2 / 10 / 10 | 10 / 10 / 10 | 30 | Guidée en Learn ; mélangée en Review |
| [`12-fichiers-permissions`](../data/chapters/12-fichiers-permissions.yaml) | 5 / 1 / 0 | 5 / 9 / 10 | 10 / 10 / 10 | 30 | Indépendante, mélangée |
| **Corpus entier** | **44 / 27 / 4** | **7 / 19 / 20** | **51 / 46 / 24** | **121** | |

Les tableaux suivent l’ordre de la matrice approuvée. B = Beginner, I =
Intermediate, A = Advanced. Les codes SH, FS, BASH et ZSH renvoient aux
[familles de références du référentiel](CURRICULUM.md#références-par-compétence).

## `00-premiers-pas`

| Identifiant | Niveau | Compétence du référentiel | Réponse principale | Point vérifié | Référence officielle et section |
| --- | --- | --- | --- | --- | --- |
| `cd-home` | B | SH/tilde · Retour au dossier personnel | `cd ~` | Le tilde non protégé devient le chemin de `HOME` avant le déplacement. | [Bash — Tilde Expansion](https://www.gnu.org/software/bash/manual/html_node/Tilde-Expansion.html) ; [Bash — Bourne Shell Builtins, `cd`](https://www.gnu.org/software/bash/manual/html_node/Bourne-Shell-Builtins.html#index-cd) |
| `man-ls` | B | FS/man · Page locale d’une commande externe | `man ls` | La section 1 distingue les commandes des autres entrées du manuel. | [POSIX `man` — SYNOPSIS et DESCRIPTION](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/man.html) ; [macOS `man(1)` — sections du manuel](https://keith.github.io/xcode-man-pages/man.1.html) |
| `cd-previous` | I | SH/cd · Dossier précédent et non parent | `cd -` | `OLDPWD` garde le dernier dossier quitté ; `..` part du dossier courant. | [Bash — Bourne Shell Builtins, `cd`](https://www.gnu.org/software/bash/manual/html_node/Bourne-Shell-Builtins.html#index-cd) |
| `ls-other-directory` | I | FS/ls · Inspection sans déplacement | `ls /atelier/livraisons` | Un opérande absolu suffit ; le shell garde son répertoire courant. | [POSIX `ls` — DESCRIPTION, `-d`, format long](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/ls.html) |
| `mkdir-parents` | I | FS/mkdir · Création des parents manquants | `mkdir -p livraisons/juin/brut` | `-p` crée les trois composants manquants dans l’ordre. | [POSIX `mkdir` — option `-p`](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/mkdir.html) |
| `cat-hidden-path` | I | FS/cat · Accès direct à un chemin caché | `cat /atelier/.etat` | Le point fait partie du nom ; l’accès direct ne dépend d’aucun listage. | [POSIX `cat` — SYNOPSIS et DESCRIPTION](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/cat.html) ; [POSIX `ls` — DESCRIPTION, `-d`, format long](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/ls.html) |
| `ls-directory-entry` | I | FS/ls · Entrée d’un dossier plutôt que son contenu | `ls -d sorties` | `-d` affiche l’entrée elle-même. | [POSIX `ls` — DESCRIPTION, `-d`, format long](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/ls.html) |
| `glob-print-names` | I | SH/glob · Expansion avant exécution | `printf '%s\n' scan-*.png` | Le shell développe l’astérisque ; le format imprime ensuite chaque argument. | [Bash — Filename Expansion et Pattern Matching](https://www.gnu.org/software/bash/manual/html_node/Filename-Expansion.html) ; [Bash — Bash Builtin Commands, `printf`](https://www.gnu.org/software/bash/manual/html_node/Bash-Builtins.html#index-printf) |
| `glob-bracket-list` | I | SH/glob · Une position, plusieurs caractères précis | `plan-[135].svg` | Les crochets correspondent à une position et excluent les chiffres 2 et 4. | [Bash — Pattern Matching, `[…]`](https://www.gnu.org/software/bash/manual/html_node/Pattern-Matching.html) |
| `quote-literal-question` | I | SH/quoting · Point d’interrogation dans un nom | `cat 'lot?.csv'` | Les guillemets neutralisent `?` sans changer le nom transmis. | [Bash — Quoting](https://www.gnu.org/software/bash/manual/html_node/Quoting.html) ; [Bash — Pattern Matching, `[…]`](https://www.gnu.org/software/bash/manual/html_node/Pattern-Matching.html) |
| `literal-tilde-directory` | I | SH/tilde · Protéger le nom littéral `~` | `cat '~/journal.txt'` | Le tilde cité est le nom du sous-dossier ; le tilde nu désigne `HOME`. | [Bash — Tilde Expansion](https://www.gnu.org/software/bash/manual/html_node/Tilde-Expansion.html) ; [Bash — Quoting](https://www.gnu.org/software/bash/manual/html_node/Quoting.html) |
| `command-origin` | I | SH/resolution · Interroger la résolution du shell | `command -v archiver` | `command -v` connaît les fonctions du shell, en plus des exécutables. | [Bash — Bash Builtins, `command`](https://www.gnu.org/software/bash/manual/html_node/Bash-Builtins.html#index-command) ; [Bash — Bash Builtin Commands, `type`](https://www.gnu.org/software/bash/manual/html_node/Bash-Builtins.html#index-type) |
| `cdpath-priority` | A | SH/cd · Priorité de CDPATH sur un homonyme local | `/partage/images` | Le premier dossier correspondant de `CDPATH` l’emporte sur `.`. | [Bash — Bourne Shell Builtins, `cd`](https://www.gnu.org/software/bash/manual/html_node/Bourne-Shell-Builtins.html#index-cd) |
| `cd-symlink-parent` | A | SH/cd · Parent logique et parent physique | `/atelier /srv` | `..` est traité avant ou après la résolution du lien selon le mode. | [Bash — Bourne Shell Builtins, `cd`](https://www.gnu.org/software/bash/manual/html_node/Bourne-Shell-Builtins.html#index-cd) |
| `glob-quoted-prefix` | A | SH/quoting · Valeur protégée et motif final actif | `printf '%s\n' "$lot"-*.csv` | Seule la valeur est protégée ; le motif final reste actif. | [Bash — Double Quotes et Filename Expansion](https://www.gnu.org/software/bash/manual/html_node/Double-Quotes.html) ; [Bash — Filename Expansion et Pattern Matching](https://www.gnu.org/software/bash/manual/html_node/Filename-Expansion.html) ; [Bash — Bash Builtin Commands, `printf`](https://www.gnu.org/software/bash/manual/html_node/Bash-Builtins.html#index-printf) ; [Bash — Shell Parameter Expansion, `${parameter}`](https://www.gnu.org/software/bash/manual/html_node/Shell-Parameter-Expansion.html) |
| `glob-no-match-shells` | A | BASH/glob + ZSH/glob · Comparer une expansion sans résultat | `Bash transmet le motif ; Zsh refuse la ligne.` | Le motif ne disparaît pas dans Bash ; Zsh arrête avant l’appel. | [Bash — Filename Expansion et Pattern Matching](https://www.gnu.org/software/bash/manual/html_node/Filename-Expansion.html) ; [Zsh — Options, `NOMATCH` et `NULL_GLOB`](https://zsh.sourceforge.io/Doc/Release/Options.html#Expansion-and-Globbing) |
| `glob-leading-dot` | A | SH/glob · Point initial et suffixe identique | `menu.json` | Le suffixe commun ne suffit pas si le point initial n’est pas admis. | [Bash — Filename Expansion et Pattern Matching](https://www.gnu.org/software/bash/manual/html_node/Filename-Expansion.html) |
| `builtin-shadowed` | A | SH/resolution · Contourner une fonction homonyme | `builtin printf '%s\n' prêt` | `builtin` sélectionne la commande intégrée ; `command` saute la fonction. | [Bash — Bash Builtins, `builtin`](https://www.gnu.org/software/bash/manual/html_node/Bash-Builtins.html#index-builtin) ; [Bash — Command Search and Execution](https://www.gnu.org/software/bash/manual/html_node/Command-Search-and-Execution.html) ; [Bash — Bash Builtins, `command`](https://www.gnu.org/software/bash/manual/html_node/Bash-Builtins.html#index-command) |
| `subshell-directory` | A | SH/subshell · Isoler un changement de dossier | `(cd exports && pwd)` | Les parenthèses isolent le changement de dossier ; `&&` protège l’étape suivante. | [Bash — Command Grouping, parenthèses](https://www.gnu.org/software/bash/manual/html_node/Command-Grouping.html) ; [Bash — Bourne Shell Builtins, `cd`](https://www.gnu.org/software/bash/manual/html_node/Bourne-Shell-Builtins.html#index-cd) |
| `relative-from-peer` | A | SH/cd · Recalculer un chemin depuis une autre branche | `../../plans/est/coupe.svg` | Deux remontées rejoignent l’ancêtre commun avant la descente. | [Bash — Bourne Shell Builtins, `cd`](https://www.gnu.org/software/bash/manual/html_node/Bourne-Shell-Builtins.html#index-cd) ; [POSIX — Pathname Resolution, chapitre 4](https://pubs.opengroup.org/onlinepubs/9799919799/basedefs/V1_chap04.html#tag_04_16) |
| `dash-leading-name` | A | FS/options · Séparer opérande et options | `cat -- -notes.txt` | `--` termine les options ; `./` supprime l’ambiguïté autrement. | [POSIX `cat` — SYNOPSIS et DESCRIPTION](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/cat.html) ; [POSIX — Utility Syntax Guidelines, règle 10](https://pubs.opengroup.org/onlinepubs/9799919799/basedefs/V1_chap12.html#tag_12_02) ; [macOS `cat(1)` — utilitaire natif](https://keith.github.io/xcode-man-pages/cat.1.html) |
| `case-sensitive-volume` | A | FS/ls · Comparer deux modes de résolution APFS | `Seulement sur le volume insensible à la casse.` | Les modes APFS diffèrent sur la comparaison du nom, sans changer sa graphie. | [Apple — Formats de système de fichiers, APFS sensible à la casse](https://support.apple.com/guide/disk-utility/file-system-formats-dsku19ed921c/mac) ; [POSIX `ls` — DESCRIPTION, `-d`, format long](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/ls.html) |

## `12-fichiers-permissions`

| Identifiant | Niveau | Compétence du référentiel | Réponse principale | Point vérifié | Référence officielle et section |
| --- | --- | --- | --- | --- | --- |
| `touch-new-file` | B | FS/touch · Création d’un fichier ordinaire absent | `touch accueil.txt` | L’utilitaire crée le nom absent, sans avoir à écrire de contenu. | [POSIX `touch` — DESCRIPTION et option `-c`](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/touch.html) |
| `copy-backup-name` | B | FS/cp · Copier en choisissant un nouveau nom | `cp devis.txt devis.bak` | Le second opérande est le nom libre de la nouvelle copie. | [POSIX `cp` — DESCRIPTION et options `-R`, `-i`, `-p`](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/cp.html) |
| `read-owner-column` | B | FS/ls · Lire le propriétaire dans une sortie longue | `nora` | Après le nombre de liens vient le propriétaire, puis le groupe. | [POSIX `ls` — format long, colonne propriétaire et type](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/ls.html) |
| `chmod-group-no-write` | B | FS/chmod · Retirer un seul droit au groupe | `chmod g-w maquette.txt` | `g-w` conserve le droit d’écriture du propriétaire. | [POSIX `chmod` — modes symbolique et absolu](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/chmod.html) |
| `recognize-link-entry` | B | FS/ls · Reconnaître le type d’une entrée | `edition` | Le caractère de type `l` identifie l’entrée avant la flèche. | [POSIX `ls` — format long, colonne propriétaire et type](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/ls.html) |
| `copy-directory-tree` | I | FS/cp · Copie récursive des descendants | `cp -R modeles livraison` | `-R` descend dans les sous-dossiers ; la cible libre prend le nom choisi. | [POSIX `cp` — DESCRIPTION et options `-R`, `-i`, `-p`](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/cp.html) ; [GNU Coreutils — `cp` invocation](https://www.gnu.org/software/coreutils/manual/html_node/cp-invocation.html) |
| `copy-confirm-overwrite` | I | FS/cp · Demander avant un écrasement | `cp -i index.txt depot/index.txt` | `-i` laisse décider au moment de la copie ; `-n` ne demande rien. | [POSIX `cp` — DESCRIPTION et options `-R`, `-i`, `-p`](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/cp.html) ; [GNU Coreutils — `cp` invocation](https://www.gnu.org/software/coreutils/manual/html_node/cp-invocation.html) |
| `move-many-sources` | I | FS/mv · Déplacer une liste précise de sources | `mv rouge.txt bleu.txt collecte/` | Le dernier opérande est le dossier commun ; pas de motif qui inclurait `vert.txt`. | [POSIX `mv` — deux formes du SYNOPSIS, DESCRIPTION](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/mv.html) |
| `copy-preserve-times` | I | FS/cp · Préserver accès et modification sur la copie | `cp -p source.txt copie.txt` | `-p` préserve notamment accès et modification, pas une date de création. | [POSIX `cp` — DESCRIPTION et options `-R`, `-i`, `-p`](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/cp.html) ; [macOS `cp(1)` — option `-p`](https://keith.github.io/xcode-man-pages/cp.1.html) |
| `link-relative-target` | I | FS/ln · Écrire une cible relative au dossier du lien | `ln -s ../source/fiche.txt vues/fiche` | La cible est relative au répertoire du lien, pas au répertoire courant. | [POSIX `ln` — liens symboliques et option `-s`](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/ln.html) ; [POSIX — §4.16 Pathname Resolution](https://pubs.opengroup.org/onlinepubs/9799919799/basedefs/V1_chap04.html#tag_04_16) |
| `readlink-stored-target` | I | FS/readlink · Lire la cible stockée sans la résoudre | `readlink raccourci` | L’outil affiche la chaîne brute, même si la cible manque. | [macOS `stat(1)` — mode `readlink`](https://keith.github.io/xcode-man-pages/stat.1.html) |
| `chmod-two-classes` | I | FS/chmod · Modifier deux classes avec des clauses ciblées | `chmod g+r,o-w bilan.txt` | Deux clauses ciblent deux classes et préservent les autres bits. | [POSIX `chmod` — modes symbolique et absolu](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/chmod.html) |
| `interpret-octal-mode` | I | FS/chmod · Décoder trois classes en octal | `rw-r-x--x` | Chaque chiffre est une somme de bits : lecture = 4, écriture = 2, exécution = 1. | [POSIX `chmod` — modes symbolique et absolu](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/chmod.html) |
| `missing-directory-search` | I | FS/permissions · Traverser un parent malgré un fichier lisible | `Exécution` | Le droit de lire le fichier ne remplace pas la traversée du parent. | [POSIX — §4.7 File Access Permissions](https://pubs.opengroup.org/onlinepubs/9799919799/basedefs/V1_chap04.html#tag_04_07) ; [POSIX — §4.16 Pathname Resolution](https://pubs.opengroup.org/onlinepubs/9799919799/basedefs/V1_chap04.html#tag_04_16) |
| `remove-symlink-only` | A | FS/rm + FS/ln · Retirer le nom du lien, pas celui de la cible | `Seul acces disparaît.` | Le retrait du lien ne touche pas à l’entrée de la cible. | [POSIX `unlink()` — DESCRIPTION et ERRORS](https://pubs.opengroup.org/onlinepubs/9799919799/functions/unlink.html) ; [POSIX `ln` — liens symboliques et option `-s`](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/ln.html) ; [POSIX `rm` — DESCRIPTION et `unlink`](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/rm.html) |
| `umask-file-and-directory` | A | FS/umask · Appliquer un masque de bits à deux modes demandés | `640 750` | Le complément du masque filtre le mode demandé, sans soustraction. | [POSIX `umask()` — DESCRIPTION](https://pubs.opengroup.org/onlinepubs/9799919799/functions/umask.html) ; [Linux `open(2)` — `O_CREAT`, mode et groupe](https://man7.org/linux/man-pages/man2/open.2.html) ; [Linux `umask(2)` — masque et ACL par défaut](https://man7.org/linux/man-pages/man2/umask.2.html) |
| `unlink-vs-write` | A | FS/rm + FS/permissions · Distinguer édition et retrait d’entrée | `Supprimer le nom seulement.` | Le parent autorise le retrait d’entrée ; le fichier interdit l’écriture. | [POSIX `unlink()` — DESCRIPTION et ERRORS](https://pubs.opengroup.org/onlinepubs/9799919799/functions/unlink.html) ; [POSIX — §4.7 File Access Permissions](https://pubs.opengroup.org/onlinepubs/9799919799/basedefs/V1_chap04.html#tag_04_07) ; [Linux `open(2)` — `O_CREAT`, mode et groupe](https://man7.org/linux/man-pages/man2/open.2.html) ; [POSIX `rm` — DESCRIPTION et `unlink`](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/rm.html) |
| `sticky-shared-directory` | A | FS/permissions · Appliquer la règle d’effacement en dossier partagé | `Mon fichier seulement.` | Dans ce dossier, seul le propriétaire de l’entrée a le droit de la retirer. | [POSIX — §4.5 Directory Protection](https://pubs.opengroup.org/onlinepubs/9799919799/basedefs/V1_chap04.html#tag_04_05) ; [POSIX `unlink()` — DESCRIPTION et ERRORS](https://pubs.opengroup.org/onlinepubs/9799919799/functions/unlink.html) ; [Linux `inode(7)` — mode de fichier, `S_ISGID` et `S_ISVTX`](https://man7.org/linux/man-pages/man7/inode.7.html) |
| `setgid-inherit-group` | A | FS/permissions · Hériter du groupe du répertoire Linux | `atelier` | Le groupe du dossier l’emporte sur le groupe principal du créateur. | [Linux `inode(7)` — mode de fichier, `S_ISGID` et `S_ISVTX`](https://man7.org/linux/man-pages/man7/inode.7.html) ; [Linux `open(2)` — `O_CREAT`, mode et groupe](https://man7.org/linux/man-pages/man2/open.2.html) ; [GNU Coreutils — Directory Setuid and Setgid](https://www.gnu.org/software/coreutils/manual/html_node/Directory-Setuid-and-Setgid.html) |
| `move-relative-symlink` | A | FS/ln + FS/mv · Recalculer une cible après déplacement du lien | `/atelier/archives/source/plan.txt` | La chaîne stockée ne change pas, mais son point de départ change. | [POSIX `mv` — deux formes du SYNOPSIS, DESCRIPTION](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/mv.html) ; [POSIX `ln` — liens symboliques et option `-s`](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/ln.html) ; [POSIX — §4.16 Pathname Resolution](https://pubs.opengroup.org/onlinepubs/9799919799/basedefs/V1_chap04.html#tag_04_16) ; [POSIX `rename()` — remplacement de destination](https://pubs.opengroup.org/onlinepubs/9799919799/functions/rename.html) |
| `hardlink-after-replace` | A | FS/ln · Suivre l’inode ou le nom après remplacement | `temoin seulement.` | Le lien physique garde l’inode d’origine ; le lien symbolique suit le nom remplacé. | [POSIX `ln` — liens symboliques et option `-s`](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/ln.html) ; [POSIX `rename()` — remplacement de destination](https://pubs.opengroup.org/onlinepubs/9799919799/functions/rename.html) ; [POSIX — §4.16 Pathname Resolution](https://pubs.opengroup.org/onlinepubs/9799919799/basedefs/V1_chap04.html#tag_04_16) |
| `chmod-conditional-search` | A | FS/chmod · Ajouter la traversée avec exécution conditionnelle | `chmod -R g+X depot` | `X` ajoute le bit aux dossiers et aux fichiers déjà exécutables seulement. | [POSIX `chmod` — modes symbolique et absolu](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/chmod.html) ; [GNU Coreutils — Conditional Executability](https://www.gnu.org/software/coreutils/manual/html_node/Conditional-Executability.html) |
| `copy-exact-destination-gnu` | A | FS/cp · Empêcher l’interprétation d’une destination comme dossier | `cp -T plan.txt cible` | `-T` interdit de réinterpréter la cible comme un dossier qui recevrait la copie. | [GNU Coreutils — Target directory, `-T`](https://www.gnu.org/software/coreutils/manual/html_node/Target-directory.html) ; [GNU Coreutils — `cp` invocation](https://www.gnu.org/software/coreutils/manual/html_node/cp-invocation.html) |
| `stat-mode-macos` | A | FS/stat · Adapter le relevé du mode octal à macOS | `stat -f %OLp fiche.txt` | `%OLp` prend les bits de mode en octal ; GNU utilise un autre sélecteur. | [macOS `stat(1)` — `-f`, format octal `O`, partie basse `Lp`](https://keith.github.io/xcode-man-pages/stat.1.html#Formats) ; [GNU Coreutils — `stat` invocation, `-c`, `%a`](https://www.gnu.org/software/coreutils/manual/html_node/stat-invocation.html) |

## Relecture éditoriale et portabilité

- Les **46 ajouts** ont été relus individuellement puis comparés aux **14
  exercices existants** des deux chapitres. Chaque consigne donne une situation
  indépendante, une action ou un résultat précis, sans commande à recopier
  ni indication du bon choix. Les QCM proposent des résultats concurrents.
  Les explications et solutions sont réservées au feedback.
- Dans `00`, la navigation simple précède la résolution de chemins et
  l’expansion. Le listage du contenu et l’affichage de l’entrée d’un dossier
  évaluent deux décisions différentes. Les deux exercices de lecture avec
  caractères spéciaux séparent le regroupement des mots (espaces, exercice
  historique) de la neutralisation d’un motif (point d’interrogation).
  Le tilde sert d’abord à rentrer chez soi, puis devient un caractère à protéger.
- Dans `12`, les exercices existants de copie, de renommage, de retrait et
  de lecture de métadonnées servent de base. Nouveau nom, récursion,
  confirmation, conservation des dates et interprétation de la destination
  demandent des décisions distinctes. Le changement historique d’un droit
  propriétaire reste ; les nouveaux modes symboliques travaillent le retrait,
  le ciblage de plusieurs classes et la condition sur le type ou les bits.
- Les 20 cas Advanced demandent de changer d’hypothèse ou de combiner des
  règles : ordre des expansions, résolution d’un nom ou d’un parent, portée
  d’un déplacement, comparaison de plateformes, droits sur le parent plutôt
  que sur le fichier, identité d’un inode ou ancrage d’un lien relatif. La
  longueur de la saisie ne détermine pas le niveau. Les diagnostics utilisent
  le texte libre ou un QCM ; les gestes restent des commandes simulées.
- Bash et Zsh sont nommés lorsque le shell change le résultat ; les options
  pertinentes sont données dans les diagnostics. APFS est comparé selon le
  mode du volume. Les scénarios de droits distinguent les classes, les
  droits du parent et les restrictions utiles. L’héritage setgid et le sticky
  bit sont cadrés sous Linux. La copie avec `-T` vise GNU Coreutils ; le
  relevé `stat` demande le format BSD de macOS et explique le format GNU.
- [Introduction au shell du MIT](https://missing.csail.mit.edu/2020/course-shell/)
  et [The Unix Shell de Software Carpentry](https://swcarpentry.github.io/shell-novice/)
  ont servi à vérifier l’enchaînement navigation → fichiers. Les décisions
  techniques reposent sur POSIX, les manuels Bash, Zsh et GNU Coreutils,
  les pages macOS, Apple et les Linux man-pages cités ci-dessus. Context7
  a complété la consultation des manuels Bash et GNU Coreutils.
- Aucune compétence de cette vague n’est différée ou substituée. Les cinq
  autres vagues restent ouvertes ; la matrice initiale est conservée.

## Vérification

- Les tests de contrat ont d’abord échoué sur le corpus initial, puis réussi
  avec les 46 ajouts. `go test -count=1 ./...` et `go vet ./...` : succès.
  Aucun fichier Go modifié ne nécessite de formatage ; contrôle du diff réussi.
- Comparaison à `d945059` : les 75 exercices historiques sont identiques
  octet pour octet. Les 46 clés ajoutées n’apparaissent pas dans l’historique
  local du corpus ni dans les fixtures héritées. Les 121 consignes sont
  distinctes après normalisation des espaces.
- Chargement et métadonnées explicites, 121 exercices dans 12 chapitres, niveaux
  et formats valides, modules vides pour les deux chapitres. Le correcteur
  accepte les 46 réponses principales et toutes les variantes déclarées,
  rejette une erreur plausible par ajout et vérifie chaque distracteur des
  nouveaux QCM. L’acceptation reste limitée aux variantes déclarées et à
  la normalisation du moteur existant.
- Les 28 commandes principales et leurs variantes passent l’analyse
  syntaxique Bash sans exécution. Les autres ajouts sont 12 réponses libres
  et 6 QCM. Les 46 lignes du rapport ont été rapprochées des réponses YAML.
- Les tests de séance contrôlent les dix identifiants de chaque niveau, l’ordre
  éditorial de `00` en Learn et la permutation complète en Review ainsi que
  pour `12`. Les tests CLI/TUI préservent les acquisitions, les scores v1 et
  les préférences synthétiques. Tous les QCM embarqués sont parcourus dans
  le modèle TUI, dont les six nouveaux, avec contrôle du rendu à 80 × 24.
- Le binaire a confirmé le total via `chapters`, l’aide `train --help` et
  le chargement des deux chapitres. Sous pseudo-terminal à 80 × 24, un
  scénario Advanced a vérifié l’affichage avant réponse, le feedback,
  la sortie par Ctrl+C et l’enregistrement de la clé en JSON v1.
- Sur fixtures temporaires, Bash 5.3.20 et Zsh 5.9 ont confirmé les expansions,
  la résolution et la portée du changement de dossier. Les utilitaires
  macOS ont confirmé la copie récursive, la confirmation, les dates,
  les liens, les modes symboliques, le calcul du masque et le format de
  `stat`. Ces essais indépendants ne passent aucune commande saisie
  à un interpréteur depuis CLI-Drill.

Le moteur, les dépendances, la CLI/TUI, le schéma de progression JSON v1 et les
préférences v1 ne sont pas modifiés. Aucun vrai fichier utilisateur n’a servi
aux validations. La recette humaine d’ensemble sous Ghostty reste en P2.6.
Les comportements propres à Linux et la comparaison de deux volumes APFS
s’appuient sur les manuels ; ils n’ont pas fait l’objet d’essais sur ces
plateformes ou volumes. La branche reste locale, avec des changements non
indexés ; ni commit, ni intégration, ni publication.
