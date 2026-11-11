# P2.5 — Navigation et flux

Le deuxième lot de l'[audit du socle](CURRICULUM.md#audit-des-quatre-chapitres-fondamentaux)
ajoute neuf exercices sur `content/p25-unix-basics-flows`, depuis
`c9ad59cf39d14f7ee4439a3fb23f1956b250c5be`. Le corpus local compte
**262 exercices dans 13 chapitres**. Les 14 propositions utiles et les pistes
facultatives restent dans le référentiel.

| Chapitre | Ajouts B / I / A | Total du chapitre B / I / A | Exercices |
| --- | --- | --- | ---: |
| `00-premiers-pas` | 1 / 2 / 0 | 11 / 12 / 10 | 33 |
| `11-flux-shell` | 1 / 3 / 2 | 11 / 13 / 12 | 36 |

| Audit | Identifiant | Niveau et réponse principale | Référence |
| --- | --- | --- | --- |
| 00.1 | `ls-long-by-modification` | B · `ls -lt livraisons` | [POSIX ls, `-l` et `-t`](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/ls.html) · [GNU ls, tri](https://www.gnu.org/s/coreutils/manual/html_node/Sorting-the-output.html) |
| 00.2 | `man-external-printf-section` | I · `man 1 printf` | [man-db, sections](https://man7.org/linux/man-pages/man1/man.1.html#DESCRIPTION) · [FreeBSD man, synopsis](https://man.freebsd.org/cgi/man.cgi?query=man&sektion=1) |
| 00.3 | `pwd-physical-directory` | I · `pwd -P` | [POSIX pwd, `-P`](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/pwd.html) |
| 11.1 | `tail-follow-journal` | B · `tail -f service.log` | [POSIX tail, `-f` et nombre de lignes par défaut](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/tail.html) |
| 11.2 | `redirect-missing-parent` | I · Le shell échoue avant de lancer mesurer. | [POSIX Shell, § 2.7 et § 2.8.1](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html) |
| 11.3 | `noclobber-single-export` | I · `(set -C; exporter > rapport.txt)` | [POSIX Shell, § 2.7.2, § 2.9.4 et `set -C`](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html) |
| 11.4 | `redirect-stderr-before-stdout` | A · stderr garde la destination initiale de stdout. | [POSIX Shell, § 2.7 et § 2.7.6](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html) · [Bash, redirections](https://www.gnu.org/s/bash/manual/html_node/Redirections.html) |
| 11.5 | `show-last-exit-status` | I · `printf '%s\n' "$?"` | [POSIX Shell, § 2.5.2](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html) |
| 11.6 | `and-or-alert-after-publish` | A · alerter s’exécute après l’échec de publier. | [POSIX Shell, § 2.9.3](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html) · [Bash, listes](https://www.gnu.org/s/bash/manual/html_node/Lists.html) |

La sélection explicite de la section distingue le nouvel exercice de
`man-ls`. Le suivi concerne un journal qui reste en place. La commande
proposée limite au sous-shell la protection du fichier ordinaire.
Les deux diagnostics Advanced font interpréter l'ordre des
redirections et la suite d'un échec dans une liste interactive.

## Vérifications et limites

- Les tests ont d'abord détecté l'absence des neuf exercices. Les tests
  de correction couvrent ensuite les réponses, leurs variantes déclarées
  et les erreurs plausibles. Les tests de séance vérifient les filtres,
  l'ordre guidé et la révision mélangée.
- Sur macOS 27.2.1, les 19 formes de commande acceptées ont été vérifiées
  sur des fixtures isolées. Les essais ont utilisé des dates distinctes,
  deux sections du manuel, un lien symbolique et un journal auquel une
  ligne a été ajoutée pendant le suivi. Les comportements GNU et FreeBSD
  ont été vérifiés dans leurs manuels, sans essai sur ces systèmes.
- Les fixtures sous sh, Bash et Zsh ont vérifié l'échec avant lancement
  quand le dossier parent manque, la protection du fichier existant, la
  portée du réglage, le routage de stderr, le statut précédent et l'alerte
  après l'échec de la publication.
- La CLI texte a obtenu 12/12 en apprentissage guidé dans `00` et 13/13
  en révision mélangée dans `11`, au niveau Intermediate. Le rapport
  synthétique est resté intact et l'exécutable témoin n'a pas été lancé.
- `go test -count=1 ./...`, `go vet ./...`, le formatage Go et
  `git diff --check` passent. Les 253 exercices antérieurs ont été comparés
  à la base. Les tests de progression utilisent des données JSON v1
  synthétiques et vérifient la conservation des acquis.
