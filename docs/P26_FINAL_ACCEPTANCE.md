# CLI-Drill — Recette finale P2.6

11 novembre 2026 · base `dev` à `1f53b29eece127505d3648ed0fc02143e79bffe5`.
**P2.6 est validée sur macOS avec Ghostty.** Les onze scénarios physiques
ont été exécutés lors de la recette du 12 novembre 2026. Après correction,
le mainteneur a confirmé les flèches gauche et droite, Ctrl+B, Ctrl+F et
le séparateur de la CLI textuelle. La clôture est documentée sur
`fix/p26-physical-acceptance`, issue de `dev` à `3c4d63d8eed3fdb3b67e46e7d3fa141ddc1acd4f`.

## Matrice des preuves

| Réf. | Type et résultat | Périmètre vérifié |
| --- | --- | --- |
| A01 | Automatisé, réussi | `go test -count=1 ./...`, `go test -race -count=1 ./...` et `go vet ./...` sur macOS 27.2.1 arm64, Go 1.27.4. Contrôles `gofmt -l` et `git diff --check`. |
| A02 | Automatisé, réussi | [Moteur](../internal/chapter/visit_test.go), [CLI](../internal/app/short_sessions_test.go) et [TUI](../internal/tui/visit_test.go) : lot de six par défaut, reliquat, ordre guidé en Apprendre, mélange des chapitres indépendants et de la révision, filtres, suite sans doublon, erreur non bloquante, reprise ciblée, parcours intégral et filtre acquis. |
| A03 | Automatisé, réussi | [Progression](../internal/app/progress_test.go) et [enregistrement TUI](../internal/tui/recording_test.go) : relecture, erreur de sauvegarde, reprise depuis les acquis, historique retiré du corpus conservé, une seule acquisition par clé. JSON v1 et ses cinq champs conservés (`version`, `completed_exercises`, `chapter_scores`, `last_session`, `streak`). `stats` et `My progress` comptent les chapitres avec une tentative. |
| A04 | Automatisé, réussi | Git disponible sans fichier. Un choix explicite désactivé et les anciennes préférences sans Git restent tels quels. Les autres modules restent optionnels. [Tests des préférences](../internal/app/short_sessions_test.go) |
| A05 | Automatisé, réussi | 446 identifiants composés uniques dans 13 chapitres, dont 149 exercices débutants, 169 intermédiaires et 128 avancés. La révision candidate ajoute deux variantes de réponse aux exercices existants. |
| P01 | Pseudo-terminal macOS, réussi | Consigne, correction `compose-sed-awk-block-total` et explication complètes consultables avec PgUp/PgDn à 80×24, 60×20 et 40×12. Redimensionnement en cours d'exercice et pendant la correction sans changer d'exercice ni enregistrer de tentative supplémentaire. QCM, navigation, recherche et aide à 40×12. |
| P02 | Pseudo-terminal et pipes, réussi | `chapters` tient sur 75, 37 et 37 colonnes aux trois dimensions, avec identifiants et titres. En pipe, les trois champs tabulés et les 13 entrées sont conservés, sans ANSI. Les aides générale, `train` et `prefs` répondent. Diagnostic sans TTY, repli `train` sous `TERM=dumb`, capture sans effacement sous `TERM=dumb`, aucune couleur SGR avec `NO_COLOR` même sous `CLICOLOR_FORCE=1`. Diagnostic et commande d'activation avec Git désactivé. |
| P03 | Pseudo-terminal macOS, réussi | Lot TUI de six à 40×12, reprise d'une erreur et lot suivant. Le JSON synthétique comptait sept tentatives et six acquis au point contrôlé. Retour à l'accueil après sortie de la capture lancée depuis la TUI. |
| P04 | Pseudo-terminal macOS, réussi | Restauration termios sur SIGINT, SIGTERM, SIGHUP, SIGQUIT et SIGTSTP. Réponse déjà enregistrée conservée. Échap et Ctrl+C sortent sans tentative. Flèche haut, F1 et collage encadré ignorés. Libellés `DEL (Retour arrière)`, `Ctrl+_` et `Tab (Ctrl+I)`. Parcours mixte et reprise ciblée vérifiés. [Script de recette](../testdata/capture_pty.py) |
| P05 | Correctif publié à `3c4d63d`, tests Go et pseudo-terminal macOS réussis | Les quatre flèches suivies d'Entrée laissent la question textuelle active. Échap seul quitte avec le bilan et la progression enregistrés. La preuve consignée pour ce scénario est sous pseudo-terminal. |
| P06 | Automatisé et pseudo-terminal macOS, réussi | H08 : flèches gauche et droite CSI et SS3 acceptées dans les exercices correspondants. Les tests Go couvrent aussi l'appui étendu Kitty sans modificateur et l'absence de tentative pour les autres séquences testées. H05 : sur la base publiée et la révision candidate, la fiche affiche `Parcours intégral (12)` et la première question `Exercice 1/12`. H09 : sous pseudo-terminal, SIGTERM avant la première réponse donne `143`, restaure le terminal et laisse le fichier de progression absent. |
| A06 | Automatisé, réussi | Séparateurs dans la CLI textuelle sans couleur : 24 caractères à 24 colonnes, 40 en sortie redirigée. Le lot de six suivi de la reprise conserve ses scores `5/6` et `1/1`, avec sept tentatives au total. |
| C01 | Compilation croisée, réussie | `CGO_ENABLED=0 go build ./...` pour Linux amd64/arm64, FreeBSD amd64, OpenBSD amd64, NetBSD amd64 et DragonFly BSD amd64. La validation native Linux et BSD reste à faire. |
| D01 | Documentation et inspection des tables de touches | En édition Emacs, Ctrl+U correspond à `unix-line-discard` dans Bash et à `kill-whole-line` dans Zsh. Ctrl+_ correspond à `undo`. Associations vérifiées dans les tables de touches sous pseudo-terminal avec Bash 3.2.57 et 5.3.20, et Zsh 5.9, sans configuration personnelle. Elles décrivent le comportement du shell lorsque c'est lui qui reçoit la touche. [GNU Bash, manuel Readline](https://www.gnu.org/software/bash/manual/bash.html) · [Zsh Line Editor](https://zsh.sourceforge.io/Doc/Release/Zsh-Line-Editor.html) |
| H01–H11 | Recette physique macOS/Ghostty validée | H01 à H04, H06, H07, H08, H10 et H11 réussis. H05 et H09 comportent les réserves consignées dans la checklist. Le séparateur CLI a été confirmé physiquement après correction. Go 1.23 n'a pas été exécuté. |

Les vérifications applicatives ont utilisé des fichiers temporaires, des réponses synthétiques et des variables XDG isolées. Les commandes Go de recette ont utilisé `GOTOOLCHAIN=local GOPROXY=off` avec les dépendances déjà en cache.

Pour rejouer P04 à P06, compiler le binaire temporaire comme indiqué dans la
checklist, puis exécuter `python3 testdata/capture_pty.py "$P26_BIN"` depuis
le dépôt. Le script crée ses propres fichiers de progression temporaires.

### P05 — Sortie involontaire en CLI textuelle

1. Lancer `"$P26_BIN" --progress "$P26_DIR/fleche.json" train` et valider `pwd`.
2. À la deuxième question, appuyer sur la flèche gauche puis sur Entrée. Le pseudo-terminal a envoyé les octets `1b 5b 44 0d`.
3. Sur la base `fe12e52`, le lecteur prenait la présence d'Échap pour un abandon. La CLI affichait `Session interrompue.` et un bilan `1/1`. Le JSON conservait une tentative correcte et le seul acquis `00-premiers-pas/pwd`.

Le correctif publié à `3c4d63d` ignore toute ligne
contenant une séquence préfixée par Échap. Échap isolé, y compris après un
texte inachevé, et Ctrl+C gardent leur effet d'abandon. Dans le menu du bilan,
une ligne ignorée laisse le menu ouvert et Entrée seule conserve son action.
Si une séquence et du texte partagent une ligne, il faut ressaisir la réponse.
La CLI textuelle reste en lecture par ligne. [Lecteur CLI](../internal/app/root.go)

Le test de régression a échoué sur la base publiée avec un bilan de `1/1`
et passe après correction. Les tests Go vérifient aussi le menu du bilan,
Ctrl+C et la transition du texte à la capture dans un parcours mixte.
Sous pseudo-terminal macOS arm64, les quatre flèches, F1 et une ligne vide
ont laissé la deuxième question active. Après `ls` et Échap suivi
d'Entrée, le bilan était `2/2` et le JSON v1 comptait deux tentatives
correctes. [Tests de régression](../internal/app/recording_test.go) ·
[Essai sous pseudo-terminal](../testdata/capture_pty.py)

Pour rejouer P05 sur clavier physique, répéter les étapes 1 et 2, puis essayer
les trois autres flèches suivies d'Entrée. La deuxième question doit rester
active. Répondre `ls`, puis quitter avec Échap et Entrée à la troisième
question. Attendre le bilan `2/2` et noter R, E ou NT avec l'émulateur et
la disposition du clavier.

### H08 — Flèches gauche et droite

Sur la base publiée, le lecteur consommait les séquences CSI et SS3, puis
les écartait avant la correction. Les deux exercices n'acceptaient que les
variantes de Ctrl+B et Ctrl+F. Ces deux causes sont reproduites dans les
tests. Les octets transmis par le clavier pendant la recette Ghostty
et le mode alors actif n'ont pas été relevés. Aucun changement de délai ni
de négociation de protocole n'est nécessaire pour corriger les cas reproduits.

Le mainteneur a confirmé les flèches gauche et droite ainsi que Ctrl+B
et Ctrl+F dans Ghostty après correction. La capture reconnaît les flèches
horizontales sans modificateur.
Elle accepte ces deux variantes dans les exercices de déplacement
correspondants. Ailleurs, une flèche ignorée ne crée pas de tentative. Les
séquences modifiées et les événements de relâchement restent ignorés.
Ghostty 1.4.0 définit CSI C/D en mode normal, SS3 C/D en mode application
et une forme CSI paramétrée quand le mode Kitty signale les événements.
[Ghostty, touches en mode normal et application](https://github.com/ghostty-org/ghostty/blob/v1.4.0/src/input/function_keys.zig#L73-L76) · [Ghostty, encodage des événements Kitty](https://github.com/ghostty-org/ghostty/blob/v1.4.0/src/input/key_encode.zig#L1006-L1026)

Si une flèche échoue encore dans Ghostty, lancer ce diagnostic dans la même
fenêtre et appuyer une seule fois sur la flèche concernée. Refaire pour
l'autre flèche. Le programme affiche les octets et les intervalles de lecture en
millisecondes. Il attend au plus dix secondes. Ctrl+C quitte en mode brut.
Le terminal est restauré à la sortie, y compris sur signal intercepté.
Sous pseudo-terminal macOS, la procédure a été vérifiée avec CSI, SS3,
Ctrl+C et SIGTERM.

~~~sh
python3 -c '
import os, select, signal, sys, termios, time, tty

def stop(signum, _frame):
    raise SystemExit(128 + signum)

signals = {signal.SIGINT, signal.SIGTERM, signal.SIGHUP, signal.SIGQUIT, signal.SIGTSTP}
for sig in signals:
    signal.signal(sig, stop)

fd = sys.stdin.fileno()
if not os.isatty(fd):
    sys.exit("Lance ce diagnostic directement dans Ghostty.")
saved = termios.tcgetattr(fd)
samples = []
print("Appuie une fois sur ← ou →. Ctrl+C quitte.", flush=True)
try:
    tty.setraw(fd, termios.TCSANOW)
    deadline = time.monotonic() + 10
    while len(samples) < 32:
        timeout = max(0, deadline - time.monotonic())
        if samples:
            timeout = min(timeout, 0.3)
        if not select.select([fd], [], [], timeout)[0]:
            break
        byte = os.read(fd, 1)
        if not byte or byte == b"\x03":
            break
        samples.append((byte, time.monotonic()))
finally:
    mask = signal.pthread_sigmask(signal.SIG_BLOCK, signals)
    try:
        termios.tcsetattr(fd, termios.TCSANOW, saved)
    finally:
        signal.pthread_sigmask(signal.SIG_SETMASK, mask)

if samples:
    print("Octets :", b"".join(byte for byte, _ in samples).hex(" "))
    print("Intervalles (ms) :", [round((b[1] - a[1]) * 1000, 1) for a, b in zip(samples, samples[1:])])
else:
    print("Aucun octet relevé.")
'
~~~

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

Le chronomètre démarre à la première consigne et s'arrête au premier bilan.
La recette a relevé **1 min 20 s en Apprendre (6/6)** et **1 min 58 s en
Réviser (4/6)**. La répétition des acquis a été observée. Certaines questions
du niveau Débutant ont paru difficiles. Ces deux mesures concernent les
lots joués le 12 novembre et restent en deçà de la cible de 5 à 10 minutes.

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

La dernière colonne consigne la recette du 12 novembre et les confirmations
après correction. `R` signifie réussi.

| ID | Manipulation | Résultat attendu | Observation physique du 12 novembre |
| --- | --- | --- | --- |
| H01 | Lancer `"$P26_BIN" --progress "$P26_DIR/premier.json"`. Ouvrir `Browse chapters`, chercher `00-premiers-pas`, ouvrir la fiche, basculer `d`, `m` et `a` puis appuyer sur Échap pour revenir. Effacer la recherche, tester `h`, une recherche sans résultat et le retour à l'accueil. Quitter avec `q`. | La recherche est annoncée, le filtre vide est explicite, l'aide et le retour restent accessibles, `premier.json` est absent. | R. |
| H02 | Lancer `"$P26_BIN" --progress "$P26_DIR/chrono.json" train 00-premiers-pas --difficulty beginner`. Résoudre les six exercices sans aide externe et chronométrer jusqu'au bilan. Quitter avec `q`, Entrée. | Six exercices, bilan des bonnes réponses, erreurs et acquis. Temps et difficulté ressentie notés. | R. 1 min 20 s, 6/6. |
| H03 | Relancer la même commande avec `--mode review` et chronométrer dans les mêmes conditions. Relever les exercices rencontrés. | Six exercices différents dans le lot. Une acquisition antérieure peut réapparaître ; un tirage isolé peut ne pas en contenir. Temps noté. | R. 1 min 58 s, 4/6. Répétition des acquis observée. |
| H04 | Lancer `"$P26_BIN" --progress "$P26_DIR/reprise.json" train`. Soumettre successivement `pwd`, `ls`, `ls -a`, `cd docs`, `faux`, `cd /atelier`, chaque fois avec Entrée. Au bilan, `r` puis Entrée, répondre `cd ..`. Taper `c`, Entrée, répondre `mkdir essais`, puis interrompre avec Ctrl+C à la question suivante. Relancer la même commande. | Lot initial `5/6`, reprise `1/1` sur la seule erreur, suite sans répétition. Après relance, `1/5` ; le JSON synthétique conserve huit tentatives, sept bonnes réponses et sept clés acquises uniques. | R. |
| H05 | Lancer `"$P26_BIN" --progress "$P26_DIR/integral-cli.json" train 00-premiers-pas --all`. Dans une autre exécution, lancer `"$P26_BIN" --progress "$P26_DIR/integral-tui.json"`, ouvrir la fiche `00-premiers-pas`, taper `a` puis Entrée. Dans chaque interface, répondre `pwd` et interrompre avant de répondre à la suivante. Examiner `"$P26_BIN" --progress "$P26_DIR/integral-cli.json" stats` et `My progress` avec `integral-tui.json`. | `1/12` au lancement dans les deux interfaces, une seule bonne réponse conservée de chaque côté, mêmes totaux. | Réserve. `1/12` confirmé en CLI et sous pseudo-terminal en TUI. Le compteur TUI n'a pas été confirmé dans Ghostty. « Ma progression » affichait un acquis sur 100 exercices éligibles. |
| H06 | Lancer `"$P26_BIN" --progress "$P26_DIR/long.json"`, ouvrir `14-recherche-texte-unix`, appuyer deux fois sur `d` depuis Débutant, puis `a` et Entrée. Soumettre `faux` et passer la correction avec Entrée pour les 11 premiers exercices. Au 12e, lire la consigne sur le bloc `BEGIN` à `END` puis soumettre `faux`. Redimensionner à 80×24, 60×20 et 40×12 et parcourir la correction avec PgUp/PgDn. | Même exercice, consigne entière, solution jusqu'à `END { print n + 0 }`, explication entière et aides accessibles. Douze tentatives au total malgré la navigation et les redimensionnements. | R. Aucun problème de lisibilité à 80×24, 60×20 et 40×12. |
| H07 | À 40×12, lancer `"$P26_BIN" --progress "$P26_DIR/qcm.json"`, ouvrir `01-raccourcis-terminal`, appuyer deux fois sur `d` depuis Débutant, puis `a` et Entrée. Appuyer sur ↓ pour choisir **Reprendre…**, valider, paginer le feedback et revenir par Échap. | Choix et feedback consultables, contexte encore accessible. Le JSON ne contient qu'une tentative. | R. |
| H08 | Lancer `"$P26_BIN" --progress "$P26_DIR/touches.json" train 01-raccourcis-terminal`. Presser ← à la première invite, Entrée après la correction, puis → à la deuxième. Après Entrée, presser Échap à la troisième. Refaire avec un fichier distinct pour vérifier Ctrl+B et Ctrl+F, Retour arrière et Tab dans l'ordre des consignes. Pour le parcours mixte, lancer `"$P26_BIN" --progress "$P26_DIR/mixte.json" train 01-raccourcis-terminal --difficulty intermediate`. Suivre Ctrl+D, Ctrl+T, Ctrl+P, Ctrl+N, Ctrl+_ et Ctrl+G avec Entrée après chaque feedback. Si Ctrl+_ est indisponible, noter NT et saisir `x` pour continuer. Appuyer sur `c` au bilan, répondre `git switch \|sauvegarde` puis quitter par Échap et Entrée à la question suivante. | Les deux flèches sont acceptées sur leurs consignes respectives. Le bilan à la sortie est de 2/2, sans tentative pour Échap. Retour arrière en DEL et Tab ont leurs libellés. Noter la combinaison physique et le libellé observé pour Ctrl+_. Six captures, puis un lot de cinq QCM. Le bilan d'interruption ne compte que les réponses soumises. | R après correction. Les flèches gauche et droite et Ctrl+B/F fonctionnent dans Ghostty. Ctrl+_ et la transition capture/QCM avaient réussi lors de la première recette. |
| H09 | Exécuter la procédure de signal ci-dessous. Refaire au besoin avec Ctrl+C à l'invite de capture en mode brut. | Le réglage terminal initial est rétabli ; saisie, retour arrière et écho fonctionnent. Aucune tentative sans réponse. | Réserve. Le terminal a été restauré après SIGTERM et aucune tentative n'a été créée. Le statut `143`, vérifié sous pseudo-terminal, n'a pas été relevé dans Ghostty. |
| H10 | Ouvrir la TUI dans des fenêtres claire et sombre. Relancer avec `NO_COLOR=1 "$P26_BIN" --progress "$P26_DIR/theme.json"`. Puis exécuter `TERM=dumb "$P26_BIN"` et `"$P26_BIN" </dev/null`. Comparer `"$P26_BIN" chapters` aux trois dimensions puis `"$P26_BIN" chapters \| cat`. | Texte, état actif et feedback lisibles dans les deux thèmes ; sans couleur avec `NO_COLOR` ; diagnostic et commande `train` dans les deux modes dégradés. Sortie `chapters` redirigée tabulée et sans ANSI. | R. |
| H11 | Exécuter `"$P26_BIN" prefs` (Git actif, autres modules facultatifs), `"$P26_BIN" prefs --disable git`, puis `"$P26_BIN" --progress "$P26_DIR/git.json" train 10-workflows-dotfiles`. Ouvrir aussi la fiche Git en TUI. Relancer `prefs` puis `"$P26_BIN" prefs --enable git` et retenter le chapitre. | Le diagnostic nomme Git et indique son activation ; Git reste désactivé jusqu'à l'action explicite. Après activation le chapitre démarre, sans installer d'outil. | R. |

### H09 — Signal envoyé pendant la capture

Dans le sous-shell isolé, lancer :

~~~sh
p26_stty="$(stty -g)"
printf 'Fichier PID à consulter depuis une autre fenêtre : %s\n' "$P26_DIR/capture.pid"
sh -c 'printf "%s\n" "$$" > "$1"; shift; exec "$@"' sh \
  "$P26_DIR/capture.pid" "$P26_BIN" --progress "$P26_DIR/signal.json" \
  train 01-raccourcis-terminal
p26_status=$?
printf 'Code de sortie : %s\n' "$p26_status"
~~~

Coller le bloc entier dans la première fenêtre pour enregistrer le statut
dès la sortie de CLI-Drill. Attendre `Waiting for key`. Depuis une deuxième
fenêtre, envoyer `kill -TERM "$(cat '/chemin/affiché/capture.pid')"`, en
remplaçant le chemin par celui affiché. Dans la première fenêtre, vérifier :

~~~sh
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

Le mainteneur a validé P2.6 sur macOS avec Ghostty après exécution de H01 à
H11. Les résultats et les deux durées de séance sont consignés ci-dessus.

- H05 : `Exercice 1/12` dans la TUI est confirmé sous pseudo-terminal.
  L'affichage physique du compteur reste non confirmé.
- H09 : la restauration du terminal après SIGTERM est confirmée dans
  Ghostty. Le code `143` est vérifié sous pseudo-terminal et n'a pas été
  relevé pendant la recette physique.
- Linux et BSD : compilation croisée réussie. Validation native à faire.

La difficulté de certains exercices Beginner et l'ergonomie des QCM sont
reportées au backlog comme améliorations non bloquantes. [Backlog](../TODO.md)
