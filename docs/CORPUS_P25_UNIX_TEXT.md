# P2.5 — Recherche et traitement de texte Unix

Le chapitre `14-recherche-texte-unix` est ajouté localement sur
`content/p25-unix-text`, depuis
`b3f255646b2520dfdff00b1020e3f88e7f10a0d8`. Il relève de la famille
`unix-shell` et ne requiert aucun module tiers. Apprendre suit l'ordre
éditorial du niveau choisi et Réviser mélange les exercices éligibles.
Chaque exercice fournit son propre contexte.

| Domaine | Beginner | Intermediate | Advanced | Total |
| --- | ---: | ---: | ---: | ---: |
| `grep` | 5 | 5 | 2 | 12 |
| `find` | 3 | 5 | 4 | 12 |
| `sed` | 3 | 5 | 1 | 9 |
| `awk` | 3 | 5 | 2 | 10 |
| Composition | 0 | 2 | 3 | 5 |
| **Total** | **14** | **22** | **12** | **48** |

Le corpus local compte 240 exercices dans 13 chapitres, dont 192 exercices
publiés. Les 48 ajouts comprennent 44 commandes simulées, trois réponses
courtes et un QCM. Les diagnostics textuels concernent les statuts de `grep`,
les dates et le suivi des liens de `find`, ainsi que le suffixe de `sed -i`.

## Choix éditoriaux

- `00` introduit les chemins, motifs et guillemets. Ici, l'exercice de
  protection du motif fait choisir quel programme l'interprète.
- `06` conserve ripgrep, fd et fzf. Les exercices Unix sur les numéros,
  le texte littéral et les types installent les gestes du socle sans module.
- `11` conserve les flux, les descripteurs, les statuts de pipeline et ses
  filtres génériques. Ici, `grep -c` produit des comptes par fichier et
  `awk` filtre ou agrège à partir des champs. Le scénario NUL traite le
  passage de chemins en arguments.
- `12` enseigne les propriétés et opérations sur les fichiers. Ici,
  la taille, les dates et les liens servent à sélectionner des entrées.
- Les commandes `sed` et `awk` sont des transformations ponctuelles. Les
  actions `find` et `xargs` des exercices inspectent les fichiers. Le QCM
  sur `sed -i` compare l'interprétation de la syntaxe et explique comment
  prévisualiser sur la sortie standard.
- La séance Intermediate expose les 22 exercices du niveau. Elle peut
  dépasser la durée indicative de cinq à dix minutes. La durée humaine
  reste à mesurer pendant la recette.

## Portabilité

- Le contexte `grep -A` est enseigné pour GNU grep et macOS. Le scénario
  `-r` face à un lien interne est explicitement GNU. La question sur les
  statuts suit la garantie POSIX 0, 1, supérieur à 1, hors option `-q`.
  [GNU grep, Invoking grep](https://www.gnu.org/software/grep/manual/grep.html)
  · [Apple grep](https://github.com/apple-oss-distributions/text_cmds/blob/main/grep/grep.1)
  · [POSIX grep](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/grep.html)
- `-maxdepth` relève de GNU et BSD/macOS. `-newer` fournit un témoin
  portable et `-mtime 1` reste un diagnostic GNU. Sur macOS 27.2.1, la
  fixture native a classé 12 h, 36 h et 60 h dans les tranches 0, 1 et 2,
  en décalage avec l'arrondi décrit par le manuel Apple.
  [GNU Findutils, Age Ranges](https://www.gnu.org/software/findutils/manual/html_node/find_html/Age-Ranges.html)
  · [Apple find](https://github.com/apple-oss-distributions/shell_cmds/blob/main/find/find.1)
- `find -print0`, `xargs -0` et `sed -E` sont normalisés par POSIX.1-2024.
  Pour la chaîne NUL, la consigne exige explicitement ces capacités.
  Leur présence effective dépend des utilitaires du système.
  [POSIX find](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/find.html)
  · [POSIX xargs](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/xargs.html)
  · [POSIX sed](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/sed.html)
- `sed -i` est hors POSIX. Le diagnostic compare le suffixe facultatif
  joint sous GNU et le suffixe demandé en argument par le sed de macOS.
  [GNU sed, Command-Line Options](https://www.gnu.org/software/sed/manual/sed.html)
  · [Apple sed](https://github.com/apple-oss-distributions/text_cmds/blob/main/sed/sed.1)

## Références par exercice

La réponse indiquée est la principale. Les variantes explicites figurent dans
le [YAML du chapitre](../data/chapters/14-recherche-texte-unix.yaml). Les codes
ci-dessous renvoient aux manuels, avec l'option ou la section dans chaque
entrée.

| Code | Manuel |
| --- | --- |
| PG | [POSIX grep](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/grep.html) |
| PF | [POSIX find](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/find.html) |
| PD | [POSIX sed](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/sed.html) |
| PA | [POSIX awk](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/awk.html) |
| PS | [POSIX Shell Command Language](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html) |
| BRE | [POSIX Regular Expressions](https://pubs.opengroup.org/onlinepubs/9799919799/basedefs/V1_chap09.html) |
| PX | [POSIX xargs](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/xargs.html) |
| PC | [POSIX cksum](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/cksum.html) |
| PL | [POSIX file](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/file.html) |
| PO | [POSIX sort](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/sort.html) |
| PU | [POSIX uniq](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/uniq.html) |
| GG | [GNU Grep](https://www.gnu.org/software/grep/manual/grep.html) |
| GF | [GNU Findutils](https://www.gnu.org/software/findutils/manual/html_mono/find.html) |
| GD | [GNU sed](https://www.gnu.org/software/sed/manual/sed.html) |
| GA | [GNU Awk](https://www.gnu.org/software/gawk/manual/gawk.html) |
| AG | [Apple grep](https://github.com/apple-oss-distributions/text_cmds/blob/main/grep/grep.1) |
| AF | [Apple find](https://github.com/apple-oss-distributions/shell_cmds/blob/main/find/find.1) |
| AD | [Apple sed](https://github.com/apple-oss-distributions/text_cmds/blob/main/sed/sed.1) |

### Beginner

- `grep-single-log` — Chercher dans un fichier.
  Réponse `grep timeout api.log`.
  Références PG: DESCRIPTION, OPERANDS.
- `grep-two-deploy-logs` — Chercher dans plusieurs fichiers et repérer leur provenance.
  Réponse `grep rollback lundi.log mardi.log`.
  Références PG: STDOUT, plusieurs fichiers.
- `grep-ignore-case` — Ignorer la casse.
  Réponse `grep -i failed deploiement.log`.
  Références PG: OPTIONS -i.
- `grep-invert-healthchecks` — Exclure les lignes qui correspondent.
  Réponse `grep -v /health requetes.log`.
  Références PG: OPTIONS -v.
- `grep-count-per-file` — Compter les lignes retenues par fichier.
  Réponse `grep -c RETRY avant.log apres.log`.
  Références PG: OPTIONS -c, STDOUT.
- `find-root-inventory` — Choisir le point de départ du parcours.
  Réponse `find modules -print`.
  Références PF: DESCRIPTION, opérande path et -print.
- `find-name-config-pattern` — Appliquer un motif au nom d’une entrée.
  Réponse `find conf -name '*.ini'`.
  Références PF: -name, PS: Quoting.
- `find-regular-files` — Sélectionner un type d’entrée.
  Réponse `find depot -type f`.
  Références PF: OPTIONS, -type.
- `sed-first-change-per-line` — Remplacer la première occurrence sur chaque ligne.
  Réponse `sed 's/TEMP/PRET/' messages.txt`.
  Références PD: Editing Commands, s, GD: The s Command.
- `sed-show-line-number` — Afficher une ligne adressée par son numéro.
  Réponse `sed -n '17p' deploiement.log`.
  Références PD: OPTIONS -n, Addresses, commande p.
- `sed-remove-header` — Supprimer la première ligne de la sortie.
  Réponse `sed '1d' import.txt`.
  Références PD: Addresses, commande d.
- `awk-first-whitespace-field` — Extraire le premier champ séparé par des blancs.
  Réponse `awk '{ print $1 }' machines.txt`.
  Références PA: variables de champ, GA: Fields.
- `awk-delimited-reordered-fields` — Choisir un séparateur d’entrée et réordonner deux champs.
  Réponse `awk -F ':' '{ print $3, $1 }' services.txt`.
  Références PA: OPTIONS -F, instruction print, GA: Field Separators.
- `awk-last-variable-field` — Extraire le dernier champ d’une ligne de largeur variable.
  Réponse `awk '{ print $NF }' mesures.txt`.
  Références PA: variables de champ, NF, GA: Fields.

### Intermediate

- `grep-fixed-release` — Traiter un caractère spécial comme du texte.
  Réponse `grep -F 'v1.2' versions.log`.
  Références PG: OPTIONS -F.
- `grep-line-number-config` — Localiser les lignes à corriger.
  Réponse `grep -n keepalive serveur.conf`.
  Références PG: OPTIONS -n, STDOUT.
- `grep-bre-numbered-tag` — Ancrer un motif et reconnaître une classe en BRE.
  Réponse `grep '^T[0-9]:' tickets.log`.
  Références PG: DESCRIPTION, BRE: 9.3.
- `grep-ere-service-alternation` — Regrouper deux variantes dans une ERE.
  Réponse `grep -nE '^ERROR (api|worker):' application.log`.
  Références PG: OPTIONS -E et -n, BRE: 9.4.
- `grep-context-after-failure` — Lire les lignes qui suivent une correspondance.
  Réponse `grep -A 2 PANIC panne.log`.
  Références GG: Context Line Control, AG: -A.
- `find-depth-direct-children` — Limiter la descente à un niveau.
  Réponse `find exports -maxdepth 1 -type f`.
  Références GF: Directories (-maxdepth, -depth), AF: -maxdepth.
- `find-size-byte-threshold` — Comparer une taille en octets.
  Réponse `find captures -type f -size +10000c`.
  Références PF: -size, préfixes numériques.
- `find-newer-than-deploy` — Comparer deux dates de modification.
  Réponse `find config -type f -newer livraison.ok`.
  Références PF: -newer, GF: Comparing Timestamps.
- `find-group-name-alternatives` — Appliquer le même test de type à deux motifs.
  Réponse `find fixtures -type f \( -name '*.json' -o -name '*.yaml' \)`.
  Références PF: opérateurs, PS: Quoting.
- `find-negate-file-suffix` — Exclure des noms sans modifier le parcours.
  Réponse `find archives -type f '!' -name '*.bak'`.
  Références PF: -name, opérateur !.
- `sed-select-marked-section` — Sélectionner un intervalle entre deux motifs.
  Réponse `sed -n '/^BEGIN$/,/^END$/p' inventaire.txt`.
  Références PD: Addresses à deux bornes, commande p.
- `sed-replace-all-on-line` — Remplacer plusieurs occurrences d’une même ligne.
  Réponse `sed 's/DEBUG/INFO/g' options.txt`.
  Références PD: commande s et drapeau g, GD: The s Command.
- `sed-path-delimiter` — Choisir un délimiteur de substitution adapté.
  Réponse `sed 's#/srv/old/#/srv/new/#g' chemins.txt`.
  Références PD: délimiteur de la commande s, GD: The s Command.
- `sed-wrap-matched-token` — Réutiliser le texte correspondant dans le remplacement.
  Réponse `sed -E 's/INC-[0-9]{3}/[&]/g' incidents.txt`.
  Références PD: OPTIONS -E et commande s, BRE: intervalles.
- `sed-reorder-captured-fields` — Réordonner deux groupes capturés.
  Réponse `sed -E 's/^host=([^ ]+) port=([0-9]+)$/\2:\1/' points.txt`.
  Références PD: OPTIONS -E, références arrière du remplacement.
- `awk-pattern-on-service-field` — Limiter un motif à un champ précis.
  Réponse `awk -F ':' '$1 ~ /^web-/ { print }' traces.txt`.
  Références PA: expressions ~, action par défaut.
- `awk-numeric-threshold` — Filtrer des lignes sur un seuil numérique.
  Réponse `awk '$2 > 250 { print $1 }' latences.txt`.
  Références PA: comparaisons numériques et instruction print.
- `awk-skip-single-header` — Sauter l’en-tête par le rang de l’enregistrement.
  Réponse `awk 'NR > 1 { print $2 }' disque.txt`.
  Références PA: variables NR et FNR.
- `awk-count-selected-records` — Compter les enregistrements retenus et afficher zéro.
  Réponse `awk '$2 == "DONE" { n++ } END { print n + 0 }' jobs.txt`.
  Références PA: END, valeurs non initialisées, GA: Using BEGIN/END.
- `awk-set-output-separator` — Distinguer le séparateur d’entrée de celui de sortie.
  Réponse `awk -v OFS='\t' '{ print $1, $2 }' routes.txt`.
  Références PA: OPTIONS -v, FS, OFS et print, GA: Output Separators.
- `compose-grep-sed-normalize-log` — Filtrer un journal et normaliser un préfixe.
  Réponse `grep '^ERROR:' acces.log | sed 's/^ERROR:/ALERTE:/'`.
  Références PG: DESCRIPTION, PD: STDIN et commande s.
- `compose-awk-count-statuses` — Regrouper et compter une colonne extraite.
  Réponse `awk '{ print $2 }' http.txt | sort | uniq -c`.
  Références PA: champs, PO: tri, PU: -c.

### Advanced

- `grep-gnu-recursion-links` — Choisir une récursion sans suivre les liens internes.
  Réponse `grep -r 'connection refused' services/`.
  Références GG: File and Directory Selection (-r, -R).
- `grep-exit-status` — Distinguer absence de correspondance et erreur.
  Réponse `0, 1, supérieur à 1`.
  Références PG: EXIT STATUS, GG: Exit Status.
- `find-gnu-mtime-buckets` — Interpréter les tranches de 24 heures.
  Réponse `veille.log`.
  Références GF: Age Ranges (-mtime).
- `find-prune-vendor-tree` — Exclure un sous-arbre avant d’y descendre.
  Réponse `find . -path ./vendor -prune -o -type f -name '*.conf' -print`.
  Références PF: -path, -prune et opérateur -o, GF: Directories.
- `find-link-type-with-follow` — Distinguer le lien de la cible pour un test de type.
  Réponse `Seulement avec -L.`.
  Références PF: OPTIONS -L, -type.
- `find-exec-checksum-batch` — Transmettre des chemins intacts par lots.
  Réponse `find pieces -type f -name '*.bin' -exec cksum {} +`.
  Références PF: -exec … {} +, PC: OPERANDS.
- `sed-in-place-portability` — Interpréter le suffixe de sauvegarde de sed -i.
  Réponse `GNU le lit comme commande. macOS le lit comme suffixe.`.
  Références GD: Command-Line Options (-i), AD: -i, PD: STDOUT.
- `awk-sum-category-duration` — Totaliser une mesure pour une catégorie.
  Réponse `awk -F ':' '$2 == "ERROR" { total += $3 } END { print total + 0 }' services.txt`.
  Références PA: opérateur +=, champ et END.
- `awk-format-latency-report` — Fixer l’alignement et la précision d’un relevé.
  Réponse `awk '{ printf "%-12s %7.2f\n", $1, $2 }' timings.txt`.
  Références PA: instruction printf, GA: Using printf Statements.
- `compose-find-grep-always-path` — Chercher dans des fichiers choisis en gardant leur provenance.
  Réponse `find conf -type f -name '*.conf' -exec grep -n -e RETRY vide.txt {} +`.
  Références PF: -exec … {} +, PG: OPTIONS -n, STDOUT.
- `compose-null-delimited-paths` — Transporter des chemins atypiques par séparateurs NUL.
  Réponse `find copies -type f -name '*.txt' -print0 | xargs -0 file`.
  Références PF: -print0, PX: -0, PL: OPERANDS, GF: Safe File Name Handling.
- `compose-sed-awk-block-total` — Totaliser une mesure dans un bloc de rapport.
  Réponse `sed -n '/^BEGIN$/,/^END$/p' suivi.txt | awk -F ':' '$2 == "ERROR" { n += $3 } END { print n + 0 }'`.
  Références PD: Addresses et -n, PA: conditions et END.

## Vérifications

- Les tests ciblés ont échoué avant l'ajout du YAML sur l'absence du
  chapitre. Ils vérifient ensuite les 48 réponses principales, un
  distracteur plausible par exercice, les choix du QCM et l'ordre guidé
  ou mélangé à chaque niveau. Le test global contrôle aussi chaque
  variante déclarée.
- Des fixtures temporaires sur macOS 27.2.1 et GNU grep 3.12 ont vérifié
  les sorties des 44 commandes principales et de leurs 82 variantes.
  Les cas couvrent notamment l'alternance ERE, les seuils stricts, le
  total nul, les liens, les espaces, les tirets et les sauts de ligne
  dans les chemins. Le comportement GNU de `find -mtime` et la syntaxe
  GNU de `sed -i` sont établis à partir de leurs manuels.
- La comparaison Git porte sur les douze YAML publiés et les phases
  P2.6 et P3. Les 192 exercices publiés sont identiques à la base.
- La CLI a corrigé les 22 réponses Intermediate dans l'ordre guidé,
  chargé la révision au même niveau et conservé les clés et compteurs
  historiques d'une progression JSON v1 synthétique. Une commande de
  contrôle saisie comme réponse incorrecte n'a créé aucun fichier.
- La liste `chapters` contient 13 lignes en sortie redirigée. Dans un
  pseudo-terminal, le tableau aligné occupe 82 colonnes avec le titre
  complet du chapitre. Sur 80 colonnes, des lignes peuvent se replier.
- `go test -count=1 ./...`, `go vet ./...`, le contrôle du formatage
  Go et `git diff --check` réussissent. La recette humaine et les essais
  sur des systèmes Linux et BSD distincts restent ouverts.
