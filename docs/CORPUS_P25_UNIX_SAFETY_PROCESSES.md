# P2.5 — Sécurité des fichiers et processus

Le premier lot de l'[audit du socle](CURRICULUM.md#audit-des-quatre-chapitres-fondamentaux)
ajoute 13 exercices sur `content/p25-unix-safety-processes`, depuis
`0178ab90edc468a0458e0f8bd94484f9adc5c906`. Le référentiel conserve ses 36
propositions. Les 23 autres restent ouvertes. Les ACL et `lsof` restent
facultatifs.

| Chapitre | Ajouts B / I / A | Total du chapitre B / I / A | Exercices |
| --- | --- | --- | ---: |
| `12-fichiers-permissions` | 1 / 5 / 1 | 11 / 15 / 11 | 37 |
| `13-processus-environnement` | 1 / 4 / 1 | 11 / 14 / 11 | 36 |

Le corpus local compte **253 exercices dans 13 chapitres**. Les deux chapitres
conservent leurs séances mélangées et leurs filtres de niveau.

| Audit | Identifiant | Niveau et réponse principale | Référence |
| --- | --- | --- | --- |
| 12.1 | `remove-tree-with-confirmation` | I · `rm -ri ./cache-sortie` | [POSIX rm, `-r` et `-i`](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/rm.html) |
| 12.2 | `copy-tree-keep-symlinks` | I · `cp -RP config config-copie` | [POSIX cp, `-R` et `-P`](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/cp.html) · [GNU cp, liens](https://www.gnu.org/software/coreutils/manual/html_node/cp-invocation.html) |
| 12.3 | `chmod-exact-config-mode` | I · `chmod 640 service.conf` | [POSIX chmod, modes numérique et symbolique](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/chmod.html) |
| 12.4 | `mkdir-private-mode` | B · `mkdir -m 700 prive` | [POSIX mkdir, `-m`](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/mkdir.html) |
| 12.5 | `du-directory-kib` | I · `du -sk archives` | [POSIX du, `-s` et `-k`](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/du.html) |
| 12.6 | `df-target-filesystem-kib` | I · `df -Pk /mnt/backup/exports` | [POSIX df, chemin, `-P` et `-k`](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/df.html) |
| 12.7 | `move-cross-filesystem-interrupted` | A · Une copie et une suppression sont deux étapes distinctes. | [POSIX mv, cas `EXDEV`](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/mv.html) · [GNU mv](https://www.gnu.org/software/coreutils/manual/html_node/mv-invocation.html) |
| 13.1 | `run-executable-relative-path` | B · `./releve` | [POSIX Shell, § 2.9.1.4](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html) |
| 13.2 | `path-prepend-for-one-call` | I · `PATH="/opt/test/bin:$PATH" analyse` | [POSIX Shell, § 2.9.1](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html) · [POSIX env](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/env.html) |
| 13.3 | `find-process-by-exact-name` | I · `pgrep -x agentweb` | [procps pgrep, `-x`](https://man7.org/linux/man-pages/man1/pgrep.1.html) · [FreeBSD pgrep, `-x`](https://man.freebsd.org/cgi/man.cgi?query=pgrep&sektion=1) |
| 13.4 | `read-job-control-stop-state` | I · PID `4812` | [procps ps, états `T` et `S`](https://man7.org/linux/man-pages/man1/ps.1.html#PROCESS_STATE_CODES) |
| 13.5 | `diagnose-command-launch-status` | I · `127 126` | [POSIX Shell, § 2.8.2](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html) |
| 13.6 | `force-stop-verified-pid` | A · `kill -s KILL 4821` | [POSIX kill, forme `-s` et forme courte XSI](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/kill.html) · [POSIX signal.h, `SIGKILL`](https://pubs.opengroup.org/onlinepubs/9799919799/basedefs/signal.h.html) |

La suppression cible un seul arbre vérifié et demande confirmation pendant
le parcours. La copie impose `-R -P` pour préserver les liens sources sur
les implémentations POSIX, dont le choix par défaut peut varier. Le cas
`mv` porte sur la reprise après interruption, sans prédire un état unique.
Le cas `KILL` impose un PID revérifié et une demande `TERM` préalable.

## Essais et limites

- Les tests ont d'abord échoué sur les ajouts absents. Ils vérifient les
  réponses, variantes et refus plausibles, les filtres et le mélange. Les
  tests d'ordre guidé et de révision des autres chapitres passent aussi.
- Sur macOS 27.2.1, 32 formes de commande ont été vérifiées nativement sur
  des fixtures isolées. Les confirmations positives et négatives de `rm`
  ont porté sur un arbre temporaire. Les essais couvrent les liens, les
  modes, les unités disque, les statuts et un `PATH` limité à l'appel.
  Deux processus de test ont permis de vérifier la recherche par nom.
- Sur une image HFS+ temporaire, le renommage entre systèmes de fichiers
  a renvoyé `EXDEV`. Ensuite `mv` a transféré les octets et retiré la source.
  Le cas interrompu est fondé sur les manuels. Les trois formes de `kill`
  ont été vérifiées par capture d'arguments, sans émission de signal.
- La CLI a corrigé 15/15 réponses Intermediate en `12` et 14/14 en `13`
  avec une progression JSON v1 synthétique. Les acquis précédents et les
  fichiers témoins sont restés intacts.
- `go test -count=1 ./...`, `go vet ./...`, le formatage Go et
  `git diff --check` passent. Les 240 exercices historiques ont été
  comparés à la base. Les comportements GNU et FreeBSD non essayés
  localement sont documentés par leurs manuels.
