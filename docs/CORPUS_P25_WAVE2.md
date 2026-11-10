# P2.5, vague 2 — Flux, processus et environnement

10 novembre 2026. Branche locale `content/p25-wave2-flows-processes`, créée
depuis `dev` propre à `f97718a`, également pointe de `origin/dev` et de la
branche distante au départ. Les 48 ajouts couvrent les compétences de la
[vague 2 du référentiel](CURRICULUM.md#enrichissement-par-lots).

| Chapitre | Avant B / I / A | Ajouts B / I / A | Après B / I / A | Total | Séance |
| --- | ---: | ---: | ---: | ---: | --- |
| [`11-flux-shell`](../data/chapters/11-flux-shell.yaml) | 2 / 3 / 1 | 8 / 7 / 9 | 10 / 10 / 10 | 30 | Learn guidé, Review mélangé |
| [`13-processus-environnement`](../data/chapters/13-processus-environnement.yaml) | 2 / 4 / 0 | 8 / 6 / 10 | 10 / 10 / 10 | 30 | Indépendant, mélangé |
| **Corpus entier** | **51 / 46 / 24** | **16 / 13 / 19** | **67 / 59 / 43** | **169** | **12 chapitres** |

B = Beginner, I = Intermediate, A = Advanced. Les réponses ci-dessous sont
les formes principales ; les variantes acceptées et les explications
complètes sont dans les YAML. Aucun ajout ne requiert de module tiers.

## `11-flux-shell`

| Identifiant | Niveau | Compétence du référentiel | Réponse principale | Source officielle et section |
| --- | --- | --- | --- | --- |
| `stdin-report` | B | Lire stdin depuis un fichier | `cat < accueil.txt` | [POSIX §2.7.1 — Redirecting Input](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html#tag_19_07_01) ; [POSIX cat — DESCRIPTION et STDIN](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/cat.html) |
| `preview-first-lines` | B | Afficher les premières lignes | `head -n 3 consignes.txt` | [POSIX head — DESCRIPTION et -n](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/head.html) |
| `preview-last-lines` | B | Afficher les dernières lignes | `tail -n 2 journal.txt` | [POSIX tail — DESCRIPTION et -n](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/tail.html) |
| `count-file-words` | B | Compter les mots | `wc -w annonce.txt` | [POSIX wc — options -l et -w](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/wc.html) |
| `pipe-produced-words` | B | Relier un producteur à un lecteur | `emet \| wc -w` | [POSIX §2.9.2 — Pipelines](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html#tag_19_09_02) ; [POSIX wc — options -l et -w](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/wc.html) |
| `sort-full-lines` | B | Trier des lignes | `sort villes.txt` | [POSIX sort — DESCRIPTION, -k, -n, -t, LC_COLLATE](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/sort.html) |
| `dedupe-adjacent-lines` | B | Réduire des doublons adjacents | `uniq passages.txt` | [POSIX uniq — DESCRIPTION et -c](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/uniq.html) |
| `tee-new-copy` | B | Afficher tout en conservant une copie | `emet \| tee copie.txt` | [POSIX §2.9.2 — Pipelines](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html#tag_19_09_02) ; [POSIX tee — DESCRIPTION et -a](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/tee.html) |
| `sort-dedupe-scattered` | I | Trier puis dédupliquer des entrées non adjacentes | `sort couleurs.txt \| uniq` | [POSIX sort — DESCRIPTION, -k, -n, -t, LC_COLLATE](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/sort.html) ; [POSIX uniq — DESCRIPTION et -c](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/uniq.html) ; [POSIX §2.9.2 — Pipelines](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html#tag_19_09_02) |
| `count-sorted-runs` | I | Compter les occurrences par ligne distincte | `uniq -c choix.txt` | [POSIX uniq — DESCRIPTION et -c](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/uniq.html) |
| `sort-numeric-field` | I | Trier une colonne numérique | `sort -t : -k 2,2n scores.txt` | [POSIX sort — DESCRIPTION, -k, -n, -t, LC_COLLATE](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/sort.html) |
| `cut-delimited-field` | I | Extraire une colonne avec un délimiteur annoncé | `cut -d ':' -f 2 tarifs.txt` | [POSIX cut — options -d et -f](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/cut.html) |
| `tee-append-copy` | I | Ajouter une copie avec tee sans écraser | `mesurer \| tee -a mesures.log` | [POSIX tee — DESCRIPTION et -a](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/tee.html) ; [POSIX §2.9.2 — Pipelines](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html#tag_19_09_02) |
| `split-output-diagnostics` | I | Séparer les sorties normale et d’erreur | `controle > releve.log 2> incidents.log` | [POSIX Shell Command Language §2.7 — Redirection](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html#tag_19_07) ; [POSIX §2.7.2 — Redirecting Output](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html#tag_19_07_02) |
| `heredoc-literal-start` | I | Fournir un bloc littéral via un document incorporé | `<<'FIN'` | [POSIX §2.7.4 — Here-Document](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html#tag_19_07_04) ; [Bash — Redirections, Here Documents](https://www.gnu.org/software/bash/manual/html_node/Redirections.html#Here-Documents) |
| `pipefail-rightmost-failure` | A | Identifier le maillon en échec qui fixe le statut | `2` | [Bash — Pipelines, option pipefail](https://www.gnu.org/software/bash/manual/html_node/Pipelines.html) ; [Zsh — Options, PIPE_FAIL](https://zsh.sourceforge.io/Doc/Release/Options.html#index-PIPE_005fFAIL) |
| `bash-pipeline-status-vector` | A | Distinguer statut global et statuts individuels | `printf '%s\n' "${PIPESTATUS[@]}"` | [Bash Variables — PIPESTATUS](https://www.gnu.org/software/bash/manual/html_node/Bash-Variables.html#index-PIPESTATUS) ; [Bash — Pipelines, option pipefail](https://www.gnu.org/software/bash/manual/html_node/Pipelines.html) ; [Zsh — Parameters Set By The Shell, pipestatus](https://zsh.sourceforge.io/Doc/Release/Parameters.html#index-pipestatus) |
| `nul-delimited-sort` | A | Préserver les octets d’un flux de noms contenant des sauts de ligne | `sort -z noms.bin` | [GNU Coreutils — sort invocation, -z/--zero-terminated](https://www.gnu.org/software/coreutils/manual/html_node/sort-invocation.html) ; [POSIX §2.9.2 — Pipelines](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html#tag_19_09_02) |
| `sort-before-replacement` | A | Éviter la troncature si l’entrée et la sortie désignent le même fichier | `sort notes.txt > tri.tmp && mv tri.tmp notes.txt` | [POSIX §2.7.2 — Redirecting Output](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html#tag_19_07_02) ; [POSIX §2.9.3 — Lists](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html#tag_19_09_03) ; [POSIX mv — DESCRIPTION et SOURCE/DEST](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/mv.html) |
| `exec-persistent-stdout` | A | Distinguer redirection persistante et locale | `exec > trace.log` | [POSIX exec — DESCRIPTION](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/exec.html) ; [POSIX §2.7.2 — Redirecting Output](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html#tag_19_07_02) |
| `redirect-brace-group` | A | Comparer la redirection d’un groupe et d’un seul maillon | `{ preparer \| publier; } 2> cycle.log` | [POSIX §2.9.4 — Compound Commands](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html#tag_19_09_04) ; [POSIX §2.9.2 — Pipelines](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html#tag_19_09_02) ; [POSIX Shell Command Language §2.7 — Redirection](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html#tag_19_07) |
| `heredoc-expanded-body` | A | Expliquer l’expansion d’un document incorporé | `$secteur -> nord` | [POSIX §2.7.4 — Here-Document](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html#tag_19_07_04) ; [Bash — Redirections, Here Documents](https://www.gnu.org/software/bash/manual/html_node/Redirections.html#Here-Documents) |
| `sort-locale-precedence` | A | Diagnostiquer un tri différent sous deux locales | `LC_ALL prend le dessus sur LC_COLLATE.` | [POSIX sort — DESCRIPTION, -k, -n, -t, LC_COLLATE](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/sort.html) ; [POSIX Base Definitions §8.2 — Internationalization Variables](https://pubs.opengroup.org/onlinepubs/9799919799/basedefs/V1_chap08.html#tag_08_02) |
| `newline-vs-record-count` | A | Distinguer fins de ligne et enregistrements | `2 3` | [POSIX wc — options -l et -w](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/wc.html) ; [GNU Coreutils — wc invocation, -l/--lines](https://www.gnu.org/software/coreutils/manual/html_node/wc-invocation.html) |

## `13-processus-environnement`

| Identifiant | Niveau | Compétence du référentiel | Réponse principale | Source officielle et section |
| --- | --- | --- | --- | --- |
| `ps-terminal-processes` | B | Lister les processus de la session | `ps -T` | [macOS ps(1) — -T, -p, -o et champ state](https://keith.github.io/xcode-man-pages/ps.1.html) |
| `read-process-id` | B | Reconnaître le PID dans une sortie fournie | `3082` | [POSIX ps — DESCRIPTION, -p et -o](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/ps.html) ; [macOS ps(1) — -T, -p, -o et champ state](https://keith.github.io/xcode-man-pages/ps.1.html) |
| `list-exported-environment` | B | Consulter l’environnement exporté | `env` | [POSIX env — DESCRIPTION et -i](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/env.html) ; [GNU Coreutils — env invocation, -i/--ignore-environment](https://www.gnu.org/software/coreutils/manual/html_node/env-invocation.html) ; [macOS printenv(1) — DESCRIPTION](https://keith.github.io/xcode-man-pages/printenv.1.html) |
| `unset-shell-variable` | B | Retirer une variable du shell | `unset ATELIER_CACHE` | [POSIX unset — DESCRIPTION et -v](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/unset.html) |
| `list-shell-jobs` | B | Afficher les jobs du shell | `jobs` | [Bash — Job Control Builtins, jobs, fg, bg, wait](https://www.gnu.org/software/bash/manual/html_node/Job-Control-Builtins.html) ; [Zsh — Jobs & Signals, Jobs](https://zsh.sourceforge.io/Doc/Release/Jobs-_0026-Signals.html) |
| `launch-background-job` | B | Démarrer une commande fictive en arrière-plan | `rendu apercu &` | [Bash — Lists of Commands, opérateur &](https://www.gnu.org/software/bash/manual/html_node/Lists.html) ; [Zsh — Jobs & Signals, Jobs](https://zsh.sourceforge.io/Doc/Release/Jobs-_0026-Signals.html) |
| `foreground-known-job` | B | Ramener un job connu au premier plan | `fg %2` | [Bash — Job Control Builtins, jobs, fg, bg, wait](https://www.gnu.org/software/bash/manual/html_node/Job-Control-Builtins.html) ; [Zsh — Jobs & Signals, Jobs](https://zsh.sourceforge.io/Doc/Release/Jobs-_0026-Signals.html) |
| `empty-exported-variable` | B | Distinguer variable vide et variable absente | `Définie et vide.` | [POSIX env — DESCRIPTION et -i](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/env.html) ; [Bash — Shell Parameter Expansion, défini et non nul](https://www.gnu.org/software/bash/manual/html_node/Shell-Parameter-Expansion.html) |
| `ps-pid-parent-columns` | I | Afficher PID et PPID dans des colonnes choisies | `ps -p 5810 -o pid,ppid` | [POSIX ps — DESCRIPTION, -p et -o](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/ps.html) ; [macOS ps(1) — -T, -p, -o et champ state](https://keith.github.io/xcode-man-pages/ps.1.html) |
| `last-background-pid` | I | Récupérer le PID du dernier lancement en arrière-plan | `printf '%s\n' "$!"` | [Bash — Special Parameters, !](https://www.gnu.org/software/bash/manual/html_node/Special-Parameters.html) ; [Zsh — Parameters Set By The Shell, !](https://zsh.sourceforge.io/Doc/Release/Parameters.html) |
| `wait-specific-child` | I | Attendre la fin d’un processus enfant précis | `wait 5814` | [POSIX wait — DESCRIPTION et EXIT STATUS](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/wait.html) ; [Bash — Job Control Builtins, jobs, fg, bg, wait](https://www.gnu.org/software/bash/manual/html_node/Job-Control-Builtins.html) |
| `resume-suspended-in-background` | I | Reprendre un job suspendu en arrière-plan | `bg %4` | [Bash — Job Control Builtins, jobs, fg, bg, wait](https://www.gnu.org/software/bash/manual/html_node/Job-Control-Builtins.html) ; [Zsh — Jobs & Signals, Jobs](https://zsh.sourceforge.io/Doc/Release/Jobs-_0026-Signals.html) |
| `wait-child-exit-status` | I | Lire la signification d’un statut renvoyé par wait | `Le code de sortie de l’enfant.` | [POSIX wait — DESCRIPTION et EXIT STATUS](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/wait.html) ; [Bash — Job Control Builtins, jobs, fg, bg, wait](https://www.gnu.org/software/bash/manual/html_node/Job-Control-Builtins.html) |
| `launch-with-clean-environment` | I | Lancer une commande avec seulement les variables utiles | `env -i LC_ALL=C MODE=essai /atelier/bin/analyse` | [POSIX env — DESCRIPTION et -i](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/env.html) ; [GNU Coreutils — env invocation, -i/--ignore-environment](https://www.gnu.org/software/coreutils/manual/html_node/env-invocation.html) |
| `identify-job-process-group` | A | Distinguer numéro de job, PID et groupe de processus | `Le groupe de processus 4101.` | [Bash — Job Control Basics, groupe et SIGTTIN](https://www.gnu.org/software/bash/manual/html_node/Job-Control-Basics.html) ; [Zsh — Jobs & Signals, Jobs](https://zsh.sourceforge.io/Doc/Release/Jobs-_0026-Signals.html) ; [POSIX ps — DESCRIPTION, -p et -o](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/ps.html) |
| `signal-service-cleanup` | A | Choisir un signal en tenant compte du nettoyage | `kill -s USR1 5174` | [POSIX kill — -s et noms symboliques](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/kill.html) ; [POSIX Base Definitions — signal.h, noms et actions par défaut](https://pubs.opengroup.org/onlinepubs/9799919799/basedefs/signal.h.html) ; [Linux signal(7) — dispositions et signaux standard](https://man7.org/linux/man-pages/man7/signal.7.html) |
| `background-terminal-read-stop` | A | Expliquer la suspension d’un job qui lit son terminal | `SIGTTIN` | [Bash — Job Control Basics, groupe et SIGTTIN](https://www.gnu.org/software/bash/manual/html_node/Job-Control-Basics.html) ; [Zsh — Jobs & Signals, Jobs](https://zsh.sourceforge.io/Doc/Release/Jobs-_0026-Signals.html) ; [POSIX Base Definitions — signal.h, noms et actions par défaut](https://pubs.opengroup.org/onlinepubs/9799919799/basedefs/signal.h.html) |
| `path-change-clears-cache` | A | Diagnostiquer une commande absente après modification de PATH | `Bash a vidé son cache et le dossier a quitté PATH.` | [Bash — Command Search and Execution, PATH](https://www.gnu.org/software/bash/manual/html_node/Command-Search-and-Execution.html) ; [Bash — Bourne Shell Builtins, hash, export, dot et trap](https://www.gnu.org/software/bash/manual/html_node/Bourne-Shell-Builtins.html) |
| `export-after-child-launch` | A | Distinguer export ultérieur et environnement déjà reçu | `avant apres` | [Bash — Command Execution Environment, export et sous-shell](https://www.gnu.org/software/bash/manual/html_node/Command-Execution-Environment.html) ; [Bash — Bourne Shell Builtins, hash, export, dot et trap](https://www.gnu.org/software/bash/manual/html_node/Bourne-Shell-Builtins.html) |
| `subshell-variable-isolation` | A | Expliquer pourquoi un sous-shell ne modifie pas le parent | `Le groupe modifie une copie de la variable.` | [Bash — Command Grouping, parenthèses](https://www.gnu.org/software/bash/manual/html_node/Command-Grouping.html) ; [Bash — Command Execution Environment, export et sous-shell](https://www.gnu.org/software/bash/manual/html_node/Command-Execution-Environment.html) |
| `load-script-current-shell` | A | Distinguer lancement et chargement dans le shell courant | `. ./profil.sh` | [Bash — Bourne Shell Builtins, hash, export, dot et trap](https://www.gnu.org/software/bash/manual/html_node/Bourne-Shell-Builtins.html) ; [Bash — Bash Builtin Commands, source](https://www.gnu.org/software/bash/manual/html_node/Bash-Builtins.html#index-source) ; [Bash — Command Execution Environment, export et sous-shell](https://www.gnu.org/software/bash/manual/html_node/Command-Execution-Environment.html) |
| `cwd-with-empty-environment` | A | Distinguer répertoire courant hérité et variable d’environnement | `/atelier/sorties` | [POSIX exec — fonctions, attributs hérités](https://pubs.opengroup.org/onlinepubs/9799919799/functions/exec.html) ; [POSIX pwd — -P et DESCRIPTION](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/pwd.html) ; [GNU Coreutils — env invocation, -i/--ignore-environment](https://www.gnu.org/software/coreutils/manual/html_node/env-invocation.html) |
| `read-zombie-process-state` | A | Interpréter un processus présent qui ne peut plus s’exécuter | `Récupérer son statut de fin.` | [macOS ps(1) — -T, -p, -o et champ state](https://keith.github.io/xcode-man-pages/ps.1.html) ; [POSIX wait — conservation du statut](https://pubs.opengroup.org/onlinepubs/9799919799/functions/wait.html) |
| `trap-scope-subshell` | A | Déduire la portée d’un gestionnaire de signal | `Le sous-shell s’arrête sans lancer le gestionnaire.` | [POSIX trap — DESCRIPTION, réinitialisation dans le sous-shell](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/trap.html) ; [Bash — Command Execution Environment, export et sous-shell](https://www.gnu.org/software/bash/manual/html_node/Command-Execution-Environment.html) ; [Bash — Signals, commandes séparées et traps](https://www.gnu.org/software/bash/manual/html_node/Signals.html) |

## Choix éditoriaux et portabilité

- Les **60 exercices des deux chapitres** ont été relus ensemble. Chaque
  consigne donne une situation autonome, sans demander d’exécuter un geste
  précédent ni afficher la commande attendue. Les exemples dans les
  diagnostics sont les données à interpréter ; les QCM proposent des
  résultats concurrents, sans désigner le bon avant la réponse.
- En `11`, lire un flux, brancher deux commandes et dupliquer une sortie
  précèdent la composition. Réduire une série adjacente, rapprocher des
  répétitions par tri, puis compter chaque série sont trois décisions
  différentes. Les cas Advanced font raisonner sur l’identité du fichier
  ouvert, la portée d’une redirection, la délimitation des enregistrements
  ou le calcul du statut. La redirection des deux maillons porte sur
  stderr ; le tube relie toujours la sortie standard de l’amont à
  l’entrée du second maillon.
- Le document incorporé Intermediate demande le seul fragment à saisir :
  l’opérateur et son délimiteur protégé. L’Advanced demande de prévoir le
  texte transmis quand un dollar est protégé et l’autre développé. La
  saisie reste sur une ligne ; aucune nouvelle mécanique n’est nécessaire.
- En `13`, le numéro de job, le PID, le PGID et le PPID ont des rôles
  distincts. Le scénario USR1 repose sur un gestionnaire propre au service
  fictif ; la demande TERM du corpus historique reste une compétence
  différente. Les autres cas Advanced confrontent la durée de vie d’un
  processus, l’héritage initial et la portée d’un gestionnaire de signal.
- Les commandes portables ont été vérifiées dans POSIX. Les écarts sont
  cadrés par les manuels Bash et Zsh (`pipefail`/`PIPE_FAIL`,
  `PIPESTATUS`/`pipestatus`), GNU Coreutils (`sort -z`) et les pages
  macOS/BSD (`ps -T` et l’état `Z`). Context7 a complété la consultation
  des manuels Bash et GNU. Aucune compétence de la vague n’est différée.

## Vérifications

- Les tests de contrat ont d’abord échoué avec 121 exercices, puis ont
  accepté les 48 ajouts. `go test -count=1 ./...`, `go vet ./...`,
  le contrôle du formatage Go et celui du diff passent.
- Les deux YAML conservent leurs entrées historiques octet pour octet ;
  les dix autres YAML sont identiques à `f97718a`. Les 48 identifiants
  sont nouveaux, y compris par rapport aux fixtures et à l’historique
  local du corpus. Le chargement valide les métadonnées, les formats et
  le décompte exact. Les consignes sont distinctes.
- Les 48 réponses principales et les variantes déclarées passent le
  correcteur. Le test ciblé rejette une erreur plausible par ajout et
  contrôle tous les distracteurs des huit nouveaux QCM. L’analyse
  syntaxique Bash contrôle les 32 commandes simulées et leurs variantes
  sans les exécuter ; les autres ajouts sont huit réponses libres et huit
  QCM. Les 48 lignes du rapport ont été comparées aux réponses du YAML.
- Les tests de séance vérifient les dix clés de chaque niveau, l’ordre
  éditorial de `11` en Learn, le mélange en Review et celui de `13`.
  Les tests CLI vérifient la conservation des clés, des scores et des
  préférences sur données synthétiques ; les QCM sont parcourus dans le
  modèle TUI. Le binaire confirme `chapters`, l’aide de `train`, la
  correction et l’enregistrement de deux clés v1 dans une séance guidée.
- Les fixtures isolées ont confirmé le tri par clé et le tri NUL avec
  les utilitaires macOS, la distinction entre LF et enregistrements, les
  redirections, les expansions, la portée de l’environnement et les
  statuts de pipeline sous Bash 5.3.20 et Zsh 5.9. Des processus
  créés pour les essais ont confirmé le retour de `wait`, l’état zombie,
  l’exécution d’un gestionnaire USR1 et la réinitialisation d’un piège dans
  un sous-shell. L’essai de contrôle des jobs sous pseudo-terminal n’a pas
  fourni de validation complète de la fermeture ; les commandes de jobs
  s’appuient sur les manuels Bash et Zsh.

Le moteur, la TUI, les préférences et la progression JSON v1 sont
inchangés. Aucune réponse saisie dans CLI-Drill n’est exécutée. La variante
GNU de `sort -z` est vérifiée dans le manuel ; l’essai de tri NUL a utilisé
l’utilitaire BSD de macOS. La recette humaine globale sous Ghostty reste
prévue en P2.6. Quatre vagues restent ouvertes. Le diff est local et non
indexé ; aucun commit, merge, push ni PR.
