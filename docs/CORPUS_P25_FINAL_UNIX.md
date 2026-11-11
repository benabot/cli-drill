# P2.5 — Fondamentaux Unix

Jalon 1 de la [campagne de finalisation](CURRICULUM.md#enrichissement-par-lots).
Base publiée `d8b383a`. Branche locale `content/p25-unix-git-modern-expansion`.

**262 → 276 exercices**, soit **14 ajouts**
(1 / 9 / 4 Beginner / Intermediate / Advanced).

| Chapitre | Ajouts B / I / A |
| --- | --- |
| [`00-premiers-pas`](../data/chapters/00-premiers-pas.yaml) | 1 / 2 / 0 |
| [`11-flux-shell`](../data/chapters/11-flux-shell.yaml) | 0 / 2 / 1 |
| [`12-fichiers-permissions`](../data/chapters/12-fichiers-permissions.yaml) | 0 / 3 / 1 |
| [`13-processus-environnement`](../data/chapters/13-processus-environnement.yaml) | 0 / 2 / 2 |

Les 14 propositions utiles de l’audit sont couvertes. Les 22
indispensables sont déjà publiés et les pistes facultatives restent dans
le référentiel. [Audit du socle](CURRICULUM.md#audit-des-quatre-chapitres-fondamentaux)

- `00.4–00.6` abordent la recherche dans l’index des manuels, le chemin
  lexical d’un dossier absent et la résolution physique d’un fichier.
- `11.7–11.9` font écrire sur stderr, suivre une rotation et relier un
  SIGPIPE au statut du pipeline.
- `12.8–12.11` couvrent le dossier temporaire, la comparaison d’inodes,
  l’écriture dans le dossier de destination et le renommage protégé.
- `13.7–13.10` couvrent les PID des jobs, le retrait d’une variable pour
  un seul appel, le signal zéro et la collecte détachée.

Sur macOS 27.2.1, les fixtures ont vérifié la rotation effective de
`tail -F`, le chemin lexical, les liens, les inodes, `mktemp`,
`mv -n`, les descripteurs, le proxy absent, les jobs, le signal zéro,
les flux de `nohup` et le statut Bash avec `pipefail`. Le `realpath`
natif testé ne connaît pas `-e`. L’exercice annonce un exécutable
compatible POSIX.1-2024 et la réponse est vérifiée dans cette norme.
La recherche de pages locales a été essayée avec un index de deux pages
synthétiques. Le refus de renommage est couvert par la norme et la
simulation. Les comportements GNU et FreeBSD sont
documentaires. [GNU tail](https://www.gnu.org/s/coreutils/manual/html_node/tail-invocation.html)
· [FreeBSD tail](https://man.freebsd.org/cgi/man.cgi?query=tail&sektion=1)
· [FreeBSD env](https://man.freebsd.org/cgi/man.cgi?query=env&sektion=1)
· [Source du manuel GNU Coreutils](https://github.com/coreutils/coreutils/blob/master/doc/coreutils.texi)
· [Manuel Bash chez son mainteneur, Pipelines](https://tiswww.case.edu/php/chet/bash/bashref.html#Pipelines)

## Vérifications automatisées

`go test -count=1 ./...`, `go vet ./...`, contrôle `gofmt` et
`git diff --check` réussis. Les cas ciblés couvrent réponse, erreur
plausible, niveau et modules. Le contrôle du corpus couvre les variantes
déclarées. La suite de séances vérifie le filtrage et la conservation
de la progression JSON v1 avec des fixtures.

## Références par exercice

### `00`

- `man-keyword-checksum`. [POSIX/man](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/man.html)
- `dirname-export-path`. [POSIX/dirname](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/dirname.html)
- `realpath-existing-export`. [POSIX/realpath](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/realpath.html)

### `11`

- `printf-diagnostic-stderr`. [POSIX/redirection](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html)
- `tail-follow-rotated-log`. [GNU/tail](https://www.gnu.org/s/coreutils/manual/html_node/tail-invocation.html)
- `pipefail-head-sigpipe`. [BASH/pipelines](https://www.gnu.org/software/bash/manual/html_node/Pipelines.html)

### `12`

- `mktemp-work-directory`. [GNU/mktemp](https://www.gnu.org/s/coreutils/manual/html_node/mktemp-invocation.html)
- `ls-compare-file-inodes`. [POSIX/ls](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/ls.html)
- `rename-destination-write`. [POSIX/rename](https://pubs.opengroup.org/onlinepubs/9799919799/functions/rename.html)
- `mv-macos-no-clobber`. [APPLE/mv](https://github.com/apple-oss-distributions/file_cmds/blob/main/mv/mv.1)

### `13`

- `jobs-with-process-ids`. [POSIX/jobs](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/jobs.html)
- `env-client-without-proxy`. [GNU/env](https://www.gnu.org/s/coreutils/manual/html_node/env-invocation.html)
- `kill-zero-observation`. [POSIX/kill](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/kill.html)
- `nohup-background-collect`. [POSIX/nohup](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/nohup.html)
