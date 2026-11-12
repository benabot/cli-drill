# CLI-Drill — Recette finale P2.6

11 novembre 2026 · base `dev` à `1f53b29eece127505d3648ed0fc02143e79bffe5`.
Les lots 1 à 3 sont intégrés. La recette physique et la validation du mainteneur restent à faire.

## Matrice des preuves

| Réf. | Type et résultat | Périmètre vérifié |
| --- | --- | --- |
| A01 | Automatisé, réussi | `go test -count=1 ./...`, `go test -race -count=1 ./...` et `go vet ./...` sur macOS 27.2.1 arm64, Go 1.27.4. Contrôles `gofmt -l` et `git diff --check`. |
| A02 | Automatisé, réussi | [Moteur](../internal/chapter/visit_test.go), [CLI](../internal/app/short_sessions_test.go) et [TUI](../internal/tui/visit_test.go) : lot de six par défaut, reliquat, ordre guidé en Apprendre, mélange des chapitres indépendants et de la révision, filtres, suite sans doublon, erreur non bloquante, reprise ciblée, parcours intégral et filtre acquis. |
| A03 | Automatisé, réussi | [Progression](../internal/app/progress_test.go) et [enregistrement TUI](../internal/tui/recording_test.go) : relecture, erreur de sauvegarde, reprise depuis les acquis, historique retiré du corpus conservé, une seule acquisition par clé. JSON v1 et ses cinq champs conservés (`version`, `completed_exercises`, `chapter_scores`, `last_session`, `streak`). `stats` et `My progress` comptent les chapitres avec une tentative. |
| A04 | Automatisé, réussi | Git disponible sans fichier. Un choix explicite désactivé et les anciennes préférences sans Git restent tels quels. Les autres modules restent optionnels. [Tests des préférences](../internal/app/short_sessions_test.go) |
| A05 | Automatisé, réussi | 446 identifiants composés uniques dans 13 chapitres, dont 149 exercices débutants, 169 intermédiaires et 128 avancés. Le diff de `data/` entre `ef5d6ba` et la base de recette est vide. |
| P01 | Pseudo-terminal macOS, réussi | Consigne, correction `compose-sed-awk-block-total` et explication complètes consultables avec PgUp/PgDn à 80×24, 60×20 et 40×12. Redimensionnement en cours d'exercice et pendant la correction sans changer d'exercice ni enregistrer de tentative supplémentaire. QCM, navigation, recherche et aide à 40×12. |
| P02 | Pseudo-terminal et pipes, réussi | `chapters` tient sur 75, 37 et 37 colonnes aux trois dimensions, avec identifiants et titres. En pipe, les trois champs tabulés et les 13 entrées sont conservés, sans ANSI. Les aides générale, `train` et `prefs` répondent. Diagnostic sans TTY, repli `train` sous `TERM=dumb`, capture sans effacement sous `TERM=dumb`, aucune couleur SGR avec `NO_COLOR` même sous `CLICOLOR_FORCE=1`. Diagnostic et commande d'activation avec Git désactivé. |
| P03 | Pseudo-terminal macOS, réussi | Lot TUI de six à 40×12, reprise d'une erreur et lot suivant. Le JSON synthétique comptait sept tentatives et six acquis au point contrôlé. Retour à l'accueil après sortie de la capture lancée depuis la TUI. |
| P04 | Pseudo-terminal macOS, réussi | [Script de recette](../testdata/capture_pty.py) : restauration termios sur SIGINT, SIGTERM, SIGHUP, SIGQUIT et SIGTSTP, réponse déjà enregistrée conservée ; Échap seul et Ctrl+C sortent sans tentative ; flèche, F1 et collage encadré ignorés en capture ; libellés `DEL (Retour arrière)`, `Ctrl+_` et `Tab (Ctrl+I)`. Parcours intégral mixte à 11/11 avec un score, interruption à 5/6 sans menu de reprise, lot court mixte avec reprise à 1/1 et passage aux QCM. |
| P05 | Correctif local, tests Go et pseudo-terminal macOS réussis | Les quatre flèches suivies d'Entrée laissent la question active. Échap seul quitte avec le bilan et la progression enregistrés. Vérification sur clavier physique à faire. |
| C01 | Compilation croisée, réussie | `CGO_ENABLED=0 go build ./...` pour Linux amd64/arm64, FreeBSD amd64, OpenBSD amd64, NetBSD amd64 et DragonFly BSD amd64. |
| D01 | Documentation et inspection des tables de touches | En édition Emacs, Ctrl+U correspond à `unix-line-discard` dans Bash et à `kill-whole-line` dans Zsh. Ctrl+_ correspond à `undo`. Associations vérifiées dans les tables de touches sous pseudo-terminal avec Bash 3.2.57 et 5.3.20, et Zsh 5.9, sans configuration personnelle. Elles décrivent le comportement du shell lorsque c'est lui qui reçoit la touche. [GNU Bash, manuel Readline](https://www.gnu.org/software/bash/manual/bash.html) · [Zsh Line Editor](https://zsh.sourceforge.io/Doc/Release/Zsh-Line-Editor.html) |
| H01–H11 | Recette physique, à faire | Aucun clavier physique, émulateur en usage manuel ni thème clair/sombre testé pour cette recette. Linux et BSD restent à tester en exécution native. Go 1.23 n'a pas été exécuté. |

Les vérifications applicatives ont utilisé des fichiers temporaires, des réponses synthétiques et des variables XDG isolées. Les commandes Go de recette ont utilisé `GOTOOLCHAIN=local GOPROXY=off` avec les dépendances déjà en cache.

Pour rejouer P04 et P05, compiler le binaire temporaire comme indiqué dans la
checklist, puis exécuter `python3 testdata/capture_pty.py "$P26_BIN"` depuis
le dépôt. Le script crée ses propres fichiers de progression temporaires.

### P05 — Sortie involontaire en CLI textuelle

1. Lancer `"$P26_BIN" --progress "$P26_DIR/fleche.json" train` et valider `pwd`.
2. À la deuxième question, appuyer sur la flèche gauche puis sur Entrée. Le pseudo-terminal a envoyé les octets `1b 5b 44 0d`.
3. Sur la base `fe12e52`, le lecteur prenait la présence d'Échap pour un abandon. La CLI affichait `Session interrompue.` et un bilan `1/1`. Le JSON conservait une tentative correcte et le seul acquis `00-premiers-pas/pwd`.

Le correctif local sur `fix/p26-text-input-escape` ignore toute ligne
contenant une séquence préfixée par Échap. Échap isolé, y compris après un
texte inachevé, et Ctrl+C gardent leur effet d'abandon. Dans le menu du bilan,
une ligne ignorée laisse le menu ouvert et Entrée seule conserve son action.
Si une séquence et du texte partagent une ligne, il faut ressaisir la réponse.
La CLI textuelle reste en lecture par ligne. [Lecteur CLI](../internal/app/root.go#L477)

Le test de régression a échoué sur la base publiée avec un bilan de `1/1`
et passe après correction. Les tests Go vérifient aussi le menu du bilan,
Ctrl+C et la transition du texte à la capture dans un parcours mixte.
Sous pseudo-terminal macOS arm64, les quatre flèches, F1 et une ligne vide
ont laissé la deuxième question active. Après `ls` et Échap suivi
d'Entrée, le bilan était `2/2` et le JSON v1 comptait deux tentatives
correctes. [Tests de régression](../internal/app/recording_test.go) ·
[Essai sous pseudo-terminal](../testdata/capture_pty.py)

Pour vérifier sur clavier physique, répéter les étapes 1 et 2, puis essayer
les trois autres flèches suivies d'Entrée. La deuxième question doit rester
active. Répondre `ls`, puis quitter avec Échap et Entrée à la troisième
question. Attendre le bilan `2/2` et noter R, E ou NT avec l'émulateur et
la disposition du clavier.

### Limites à consigner pendant la recette physique

- La capture compare les octets transmis. Une combinaison peut être interceptée par le système ou l'émulateur, ou dépendre de la disposition et du mode d'édition. Le délai entre deux octets d'une séquence Échap est de 50 ms. Une touche transmise comme `0x08` peut être affichée `Ctrl+H` ; `0x7f` est libellé `DEL (Retour arrière)`.
- Les essais de restauration sous pseudo-terminal concernent les cinq signaux de la matrice. SIGTSTP quitte l'application. SIGKILL ne peut pas être intercepté.
- Un fichier de progression suppose un seul écrivain. Deux exécutions simultanées peuvent perdre des mises à jour. Apprendre recharge les acquis. Le curseur, le tirage et les erreurs de la visite restent en mémoire.

## Taille et durée des séances

Avec les préférences initiales (Git activé), les 39 combinaisons chapitre/niveau donnent 30 filtres accessibles et 9 filtres vides. Le découpage de chaque filtre complet en lots successifs donne 63 lots : 4 de taille 1, 5 de taille 2, 2 de taille 3, 7 de taille 4, 7 de taille 5 et 38 de taille 6. Le comptage porte sur des filtres complets, donc sur Réviser ou une première visite Apprendre sans acquis. Les modules supplémentaires changent la population éligible.

| Filtre, préférences initiales | Exercices éligibles | Lots courts successifs |
| --- | ---: | --- |
| `00-premiers-pas`, débutant | 12 | 6 + 6 |
| `01-raccourcis-terminal`, intermédiaire | 11 | 6 + 5 en Apprendre |
| `10-workflows-dotfiles`, débutant | 10 | 6 + 4 |
| `14-recherche-texte-unix`, intermédiaire | 22 | 6 + 6 + 6 + 4 |

Le chronomètre démarre à la première consigne et s'arrête au premier bilan. Inclure la lecture des corrections et noter à part les interruptions, la reprise et la suite. La cible reste `5 à 10 minutes par séance`. Aucun temps de réponse humain n'a encore été mesuré. Si un niveau dépasse régulièrement dix minutes, relever les exercices qui prennent le plus de temps, puis comparer une taille de lot réduite pour ce niveau et une taille réglable. Ces options demandent un arbitrage produit.

## Checklist physique sur Mac

Depuis la racine du dépôt, compiler et ouvrir un sous-shell temporaire. Toutes les commandes de la checklist utilisent le binaire temporaire et les préférences XDG isolées. Chaque scénario interactif a son propre fichier JSON, sauf quand la relance doit réutiliser le même état.

~~~sh
p26_dir="$(mktemp -d)"
mkdir -p "$p26_dir/data" "$p26_dir/config" "$p26_dir/zsh"
go build -o "$p26_dir/cli-drill" ./cmd/cli-drill &&
  env XDG_DATA_HOME="$p26_dir/data" XDG_CONFIG_HOME="$p26_dir/config" \
    ZDOTDIR="$p26_dir/zsh" HISTFILE=/dev/null \
    P26_DIR="$p26_dir" P26_BIN="$p26_dir/cli-drill" /bin/zsh -f
~~~

Avant les essais, noter macOS (`sw_vers`), l'émulateur et sa version, le clavier et sa disposition, le shell et son mode d'édition, `TERM`, les dimensions (`stty size` renvoie lignes puis colonnes) et le thème. Utiliser une fenêtre temporaire et redimensionner avec la souris. Sur un clavier compact Mac, utiliser `Fn+↑` et `Fn+↓` pour PgUp/PgDn si l'émulateur les transmet. Dans le navigateur, pour ouvrir un chapitre par identifiant : `Browse chapters → /`, taper l'identifiant, puis `Entrée` pour valider le filtre et `Entrée` pour ouvrir la fiche.

Pour chaque ligne, noter `R` (réussi), `E` (échec) ou `NT` (non testé), avec le temps ou le résultat observé si pertinent.

| ID | Manipulation | Résultat attendu | R / E / NT |
| --- | --- | --- | --- |
| H01 | Lancer `"$P26_BIN" --progress "$P26_DIR/premier.json"`. Ouvrir `Browse chapters`, chercher `00-premiers-pas`, ouvrir la fiche, basculer `d`, `m` et `a` puis appuyer sur Échap pour revenir. Effacer la recherche, tester `h`, une recherche sans résultat et le retour à l'accueil. Quitter avec `q`. | La recherche est annoncée, le filtre vide est explicite, l'aide et le retour restent accessibles, `premier.json` est absent. | |
| H02 | Lancer `"$P26_BIN" --progress "$P26_DIR/chrono.json" train 00-premiers-pas --difficulty beginner`. Résoudre les six exercices sans aide externe et chronométrer jusqu'au bilan. Quitter avec `q`, Entrée. | Six exercices, bilan des bonnes réponses, erreurs et acquis. Temps et difficulté ressentie notés. | |
| H03 | Relancer la même commande avec `--mode review` et chronométrer dans les mêmes conditions. Relever les exercices rencontrés. | Six exercices différents dans le lot. Une acquisition antérieure peut réapparaître ; un tirage isolé peut ne pas en contenir. Temps noté. | |
| H04 | Lancer `"$P26_BIN" --progress "$P26_DIR/reprise.json" train`. Soumettre successivement `pwd`, `ls`, `ls -a`, `cd docs`, `faux`, `cd /atelier`, chaque fois avec Entrée. Au bilan, `r` puis Entrée, répondre `cd ..`. Taper `c`, Entrée, répondre `mkdir essais`, puis interrompre avec Ctrl+C à la question suivante. Relancer la même commande. | Lot initial `5/6`, reprise `1/1` sur la seule erreur, suite sans répétition. Après relance, `1/5` ; le JSON synthétique conserve huit tentatives, sept bonnes réponses et sept clés acquises uniques. | |
| H05 | Lancer `"$P26_BIN" --progress "$P26_DIR/integral-cli.json" train 00-premiers-pas --all`. Dans une autre exécution, lancer `"$P26_BIN" --progress "$P26_DIR/integral-tui.json"`, ouvrir la fiche `00-premiers-pas`, taper `a` puis Entrée. Dans chaque interface, répondre `pwd` et interrompre avant de répondre à la suivante. Examiner `"$P26_BIN" --progress "$P26_DIR/integral-cli.json" stats` et `My progress` avec `integral-tui.json`. | `1/12` au lancement dans les deux interfaces, une seule bonne réponse conservée de chaque côté, mêmes totaux. | |
| H06 | Lancer `"$P26_BIN" --progress "$P26_DIR/long.json"`, ouvrir `14-recherche-texte-unix`, appuyer deux fois sur `d` depuis Débutant, puis `a` et Entrée. Soumettre `faux` et passer la correction avec Entrée pour les 11 premiers exercices. Au 12e, lire la consigne sur le bloc `BEGIN` à `END` puis soumettre `faux`. Redimensionner à 80×24, 60×20 et 40×12 et parcourir la correction avec PgUp/PgDn. | Même exercice, consigne entière, solution jusqu'à `END { print n + 0 }`, explication entière et aides accessibles. Douze tentatives au total malgré la navigation et les redimensionnements. | |
| H07 | À 40×12, lancer `"$P26_BIN" --progress "$P26_DIR/qcm.json"`, ouvrir `01-raccourcis-terminal`, appuyer deux fois sur `d` depuis Débutant, puis `a` et Entrée. Appuyer sur ↓ pour choisir **Reprendre…**, valider, paginer le feedback et revenir par Échap. | Choix et feedback consultables, contexte encore accessible. Le JSON ne contient qu'une tentative. | |
| H08 | Lancer `"$P26_BIN" --progress "$P26_DIR/touches.json" train 01-raccourcis-terminal`. À la première invite, presser ← et →, puis Ctrl+B, Ctrl+F, Retour arrière et Tab dans l'ordre des consignes, en passant chaque feedback avec Entrée. À l'invite suivante, presser Échap seul. Pour le parcours mixte, lancer `"$P26_BIN" --progress "$P26_DIR/mixte.json" train 01-raccourcis-terminal --difficulty intermediate`. Suivre Ctrl+D, Ctrl+T, Ctrl+P, Ctrl+N, Ctrl+_ et Ctrl+G avec Entrée après chaque feedback. Si Ctrl+_ est indisponible, noter NT et saisir `x` pour continuer. Appuyer sur `c` au bilan, répondre `git switch \|sauvegarde` puis quitter par Échap et Entrée à la question suivante. | Les flèches et les sorties n'ajoutent aucune tentative. Retour arrière en DEL et Tab ont leurs libellés. Noter la combinaison physique et le libellé observé pour Ctrl+_. Six captures, puis un lot de cinq QCM. Le bilan d'interruption ne compte que les réponses soumises. | |
| H09 | Exécuter la procédure de signal ci-dessous. Refaire au besoin avec Ctrl+C à l'invite de capture en mode brut. | Le réglage terminal initial est rétabli ; saisie, retour arrière et écho fonctionnent. Aucune tentative sans réponse. | |
| H10 | Ouvrir la TUI dans des fenêtres claire et sombre. Relancer avec `NO_COLOR=1 "$P26_BIN" --progress "$P26_DIR/theme.json"`. Puis exécuter `TERM=dumb "$P26_BIN"` et `"$P26_BIN" </dev/null`. Comparer `"$P26_BIN" chapters` aux trois dimensions puis `"$P26_BIN" chapters \| cat`. | Texte, état actif et feedback lisibles dans les deux thèmes ; sans couleur avec `NO_COLOR` ; diagnostic et commande `train` dans les deux modes dégradés. Sortie `chapters` redirigée tabulée et sans ANSI. | |
| H11 | Exécuter `"$P26_BIN" prefs` (Git actif, autres modules facultatifs), `"$P26_BIN" prefs --disable git`, puis `"$P26_BIN" --progress "$P26_DIR/git.json" train 10-workflows-dotfiles`. Ouvrir aussi la fiche Git en TUI. Relancer `prefs` puis `"$P26_BIN" prefs --enable git` et retenter le chapitre. | Le diagnostic nomme Git et indique son activation ; Git reste désactivé jusqu'à l'action explicite. Après activation le chapitre démarre, sans installer d'outil. | |

### H09 — Signal envoyé pendant la capture

Dans le sous-shell isolé, lancer :

~~~sh
p26_stty="$(stty -g)"
printf 'Fichier PID à consulter depuis une autre fenêtre : %s\n' "$P26_DIR/capture.pid"
sh -c 'printf "%s\n" "$$" > "$1"; shift; exec "$@"' sh \
  "$P26_DIR/capture.pid" "$P26_BIN" --progress "$P26_DIR/signal.json" \
  train 01-raccourcis-terminal
~~~

Attendre `Waiting for key`. Depuis une deuxième fenêtre, envoyer `kill -TERM "$(cat '/chemin/affiché/capture.pid')"`, en remplaçant le chemin par celui affiché, pendant que la capture est encore active. Dans la première fenêtre, vérifier :

~~~sh
p26_status=$?
printf 'Code de sortie : %s\n' "$p26_status"
test "$(stty -g)" = "$p26_stty" && printf 'Réglages du terminal restaurés\n'
test ! -e "$P26_DIR/signal.json" && printf 'Aucune tentative enregistrée\n'
~~~

Attendu : code `143`. Si la comparaison échoue, noter l'écart et, si la saisie est perturbée, restaurer les réglages de cette fenêtre avec `stty "$p26_stty"`. Après la recette, taper `exit` pour quitter le sous-shell. Le dossier temporaire peut être conservé pour joindre les résultats.

Sur un système Linux ou BSD disponible, compiler localement le binaire et
reprendre les mêmes scénarios avec un sous-shell sans configuration personnelle
disponible sur la machine. Pour un premier passage natif, prioriser H01, H06,
H08, H09 et H10. Noter le système et la version avec `uname -sr` ainsi que
le terminal, le shell et la disposition réellement utilisés.

## Critères de clôture

1. Les preuves automatisées et sous pseudo-terminal restent valides sur la révision candidate. Le corpus et le JSON v1 satisfont les contrôles A03 et A05.
2. Les scénarios H01 à H11 sont renseignés sur le Mac du mainteneur, avec l'environnement et les octets ou libellés observés pour les touches litigieuses. Tout échec ou NT reçoit une décision explicite.
3. Deux durées humaines sont relevées, avec le niveau et le contexte. Si l'objectif de 5 à 10 minutes est dépassé, le mainteneur décide de l'accepter ou de faire un essai complémentaire avant d'ajuster les lots.
4. Le résultat physique de P05 est consigné sur la révision candidate. Les plateformes sans recette native restent signalées.
5. Le mainteneur donne explicitement sa validation de P2.6.
