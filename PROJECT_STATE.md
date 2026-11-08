# État du projet — CLI-Drill

**Constat du 7 novembre 2026 : P0.1 à P0.3 et P1.1 terminées, intégrées localement dans `dev`.**
P1.1 comprend le stockage initial (`25c344b`) et les corrections de
`fix/progress-storage-hardening`, acceptées après revue. La migration V2
complète reste à réaliser. P1.2 est réalisée sur `fix/training-recording`,
depuis `dev` à `1f61cb9`, en attente de validation humaine et d'intégration.

## Implémentation actuelle

L'interface présente CLI-Drill comme un outil d'entraînement local au terminal.

- Module Go à la racine ; point d'entrée `cmd/cli-drill/`.
- Dix chapitres YAML intégrés, 68 exercices, contenus et identifiants inchangés.
- Corpus exclusivement intégré depuis P0.1 ; aucun chargement automatique des
  anciens YAML utilisateur, même lorsqu'ils sont invalides.
- Commandes publiques : `chapters`, `train [chapter-id]`, `stats`, `reset --yes`,
  avec `--progress`, l'aide et la complétion standard Cobra.
- Anciennes commandes `init`, `scan`, `generate`, `directory`, `search`, `show`
  et option `--config` retirées ; elles produisent désormais une erreur.
- Accueil TUI : Start training, Browse chapters, Keyboard shortcuts, My progress, Quit.
- Bannière ASCII, palette, composants de liste et de saisie et pieds de page conservés.
- Deux moteurs d'entraînement conservés : texte dans la CLI/TUI et capture réelle
  dédiée. La TUI revient automatiquement à l'accueil après le mode de capture.
- Progression JSON version 1 et chemin inchangés ; stockage atomique P1.1
  réutilisé par P1.2 pour sauvegarder chaque réponse soumise.

### P0.1 — Corpus intégré exclusif, intégrée dans `dev`

Le chargeur partagé `loadChapters` lit uniquement le système de fichiers
embarqué transmis par `data.Chapters()`. Les anciens YAML ne remplacent ni ne
complètent le corpus et ne sont ni déplacés, ni convertis, ni supprimés.
Un chapitre présent seulement dans le dossier utilisateur est indisponible.

### P0.2 — Surface dotfiles retirée, intégrée dans `dev`

Le câblage Cobra, les aides et l'interface TUI ne lisent plus de configuration
TOML, n'exposent plus d'annuaire, de scan ou de génération et ne présentent plus
les messages d'installation des dotfiles. Les diagnostics et chemins associés
ont été retirés de l'interface, sans supprimer les fichiers utilisateur.

Start training choisit le premier chapitre textuel non vide du corpus actuel,
soit `02-navigation-shell`, et ouvre directement l'exercice. Si aucun chapitre
textuel n'est disponible, il ouvre le navigateur sans forcer la capture brute.
Le chapitre sélectionné reste cohérent avec la liste lors des retours arrière.
Le comportement de `train` sans identifiant reste inchangé : premier chapitre
intégré, actuellement le chapitre de raccourcis.

Keyboard shortcuts réutilise le navigateur en ne montrant que les chapitres
contenant un exercice `key-sequence`. La sélection conserve les indices du
corpus complet et délègue au moteur existant. Les familles pédagogiques n'ont
pas encore de métadonnées nouvelles ; les raccourcis seulement textuels des
chapitres mixtes, comme Micro, ne sont pas ajoutés à cette vue de capture.
Browse chapters restaure la liste complète.

Les filtres reçoivent leurs lettres et le retour arrière avant les raccourcis
globaux. Entrée applique le filtre, puis permet de choisir ; Échap annule la
saisie ou efface le filtre appliqué avant de quitter le navigateur. Les pieds
de page reflètent ces contextes. Les dimensions de liste réservent l'espace du
titre, du pied de page et du padding pour éviter un écran tronqué.

My progress réutilise les statistiques existantes, sans compteur d'annuaire.
Aucun schéma, clé, compteur ou mécanisme de sauvegarde de progression n'est modifié.

### P0.3 — Internes hérités retirés, intégrée dans `dev`

L'architecture conserve six packages internes : `app`, `chapter`, `exercise`,
`progress`, `tui` et `xdg`, conformément à l'[architecture cible](docs/ARCHITECTURE.md).
Les cinq packages `detect`, `catalog`, `config`, `markdown` et `shell/zsh`,
leurs tests et leurs fixtures exclusivement internes ont été supprimés.
Aucun composant actif ne dépend plus du scanner, du parseur utilisateur ou du catalogue.

`chapter` conserve le chargement YAML (`Load`, `LoadFS`) et sa validation.
La génération, la conversion vers le catalogue et les lecteurs/écrivains de
répertoires externes ont été retirés. `Item.Type` est une chaîne conservant
les valeurs pédagogiques existantes ; la provenance `Source`, absente du corpus
intégré et utilisée seulement par le catalogue, a été retirée du modèle Go.
Aucun YAML intégré ni identifiant n'est modifié.

`xdg` ne résout plus que le chemin de progression, toujours
`$XDG_DATA_HOME/cli-drill/progress.json` ou, par défaut,
`~/.local/share/cli-drill/progress.json`. Les anciens chemins TOML/YAML ne sont
plus résolus par la production. Les tests construisent seuls leurs fixtures héritées.
La dépendance `BurntSushi/toml` a été retirée ; aucune version conservée ne change.

Dans `app` et `tui`, seuls des helpers et une option interne sans consommateur
réel ont été retirés ; un helper réservé aux tests a été déplacé dans les tests.
Le rendu, les parcours actifs, les moteurs, les contrôles techniques du terminal
et le stockage JSON restent inchangés. Aucun fichier utilisateur n'est lu,
déplacé, converti ou supprimé par cette refactorisation.

### P1.1 — Stockage de progression terminé, intégré localement dans `dev`

Sur la base P0.3, `Save` tronquait directement le fichier cible ; `Load`
acceptait `null` et des versions non prises en charge. Ces défauts sont corrigés.
Le stockage sérialise le même JSON version 1, indenté et terminé par un saut de
ligne, puis écrit un fichier temporaire dans le répertoire de destination.
Les permissions POSIX, la synchronisation du fichier (`Sync`) et sa fermeture
sont traitées avant le remplacement par `Rename`, atomique sur les Unix ciblés.
Un échec avant ce remplacement laisse le dernier fichier valide intact ; les
temporaires sont nettoyés sur les sorties d'erreur traitées. Si le système
refuse aussi leur suppression, cette erreur est retournée avec l'erreur initiale.

Un fichier absent produit l'état vide historique sans écriture. Les champs
optionnels absents gardent les valeurs par défaut existantes, dont la version 1
si elle est omise et une table de scores vide si elle est absente ou `null`.
Un JSON invalide, un objet racine `null`, une version autre que 1 ou une erreur de
lecture produisent une erreur contextualisée. `Save` refuse aussi de remplacer
un fichier existant illisible, invalide ou non pris en charge, même si l'appelant
lui fournit un état vide. Les identifiants et résultats des exercices retirés
restent présents après sauvegarde et rechargement.

La revue de `25c344b` a exposé des `null` structurels acceptés silencieusement,
le remplacement des fichiers en lecture seule et une couverture insuffisante
de `Sync`/`Close`. P1.1b rejette une version explicitement `null`, les scores
`null`, leurs compteurs `null`, un `streak` `null` et les acquisitions `null`
dans une liste. Les champs omis, les scores partiels et les conteneurs v1
`completed_exercises`/`chapter_scores` à `null` restent admis ; `last_session`
reste nullable. La validation typée existante continue à refuser les mauvais
types. Aucun schéma ni état historique valide n'est converti par migration.

Les nouveaux JSON ont le mode `0600` ; les bits POSIX d'un fichier existant sont
conservés. Avant remplacement, P1.1b ouvre la cible existante en écriture sans
troncature ni création, puis la referme ; un refus ou une erreur de fermeture
arrête la sauvegarde. Les destinations non régulières, notamment les liens
symboliques, sont refusées explicitement pour éviter de remplacer le lien à la place de sa
cible. Aucun schéma, chemin, score ou déclenchement de sauvegarde CLI/TUI ne
change. `reset --yes` conserve sa suppression explicitement demandée.

Garanties limitées aux fichiers locaux sous macOS et Linux, avec un seul
écrivain. Le remplacement change l'inode : les ACL, attributs étendus et autres
métadonnées spécifiques ne sont pas recopiés. Pas de verrouillage concurrent,
de synchronisation du répertoire ou de garantie contre une coupure électrique.
Un arrêt brutal empêchant les nettoyages peut laisser un temporaire voisin,
qui n'est jamais chargé comme progression ni purgé automatiquement.

L'injection de `Sync` et `Close` est privée et passée par appel au helper de
remplacement, sans état global mutable ni abstraction générale de fichiers.
Le nettoyage ferme toujours le vrai descripteur, même après une erreur injectée.

### P1.2 — Enregistrement des entraînements, réalisé à valider

L'audit sur `1f61cb9` a reproduit trois défauts : `tui.answerOrAdvance`
évaluait les réponses et le bilan sans mettre à jour la progression ;
`app.runTraining` sauvegardait seulement à la fin et perdait les réponses déjà
validées en cas d'erreur de lecture ; un EOF textuel vide créait une tentative
incorrecte. Les parcours clavier de Browse chapters et Keyboard shortcuts
délèguent au moteur CLI et héritaient de sa sauvegarde tardive.

`progress.RecordAnswer` partage uniquement l'enregistrement : charger l'état
v1 courant, appeler `MarkCompleted`, puis `Save`. Cette fonction concrète est
commune aux moteurs texte et clavier CLI et au texte Bubble Tea. Elle ne renvoie
un état utilisable qu'après sauvegarde réussie. Le stockage P1.1, ses protections,
ses chemins et son contrat à un seul écrivain sont conservés.

Règles conservées : une soumission correcte ou incorrecte ajoute une tentative,
une réponse correcte ajoute un succès et l'acquis `chapter-id/item-id` reste
unique. Les retries et nouvelles séances ajoutent des tentatives légitimes ;
les scores et identifiants retirés restent stockés. Le bilan de séance en mémoire
reste distinct des compteurs cumulés. Aide, solution, feedback et navigation ne
sont pas enregistrés. Un texte non vide suivi d'EOF est soumis une fois ; un
EOF vide abandonne sans écriture.

Choix explicité pour les validations vides : Entrée seule, espaces et fins de
ligne ne soumettent pas de réponse, dans les deux interfaces. Cela évite les
tentatives fictives sur la question suivante lors d'Entrées répétées. Les octets
de contrôle historiquement acceptés comme réponses textuelles, dont `Ctrl+L`,
ne sont pas effacés par ce filtre. La normalisation et les réponses admises
restent inchangées.

La sauvegarde a lieu avant le feedback et avant la lecture suivante. Esc et
Ctrl+C abandonnent la saisie courante ; les réponses déjà sauvegardées restent
disponibles après retour au chapitre, menu, sortie et redémarrage. La CLI
textuelle conserve sa lecture par ligne : Esc peut nécessiter Entrée avant
d'être traité ; Ctrl+C natif peut interrompre le processus. Aucun passage
supplémentaire en mode brut ni gestionnaire de signal n'est ajouté.

La CLI remonte toute erreur de lecture ou sauvegarde. La TUI affiche le
diagnostic dans les vues existantes ; en cas d'échec, elle garde la réponse
éditable et n'avance ni le bilan ni la progression affichée. Entrée permet
de réessayer après correction de la cause, sans double comptage. Quitter à ce
stade ne sauvegarde pas la réponse en échec : aucun journal ni rejeu implicite
n'est ajouté. Une consultation ou un abandon sans réponse ne crée pas de JSON.

## Limites encore présentes

Les contenus YAML conservent des alias, fonctions et formulations personnelles
héritées ; leur révision pédagogique est une étape distincte. Le README public
n'est pas modifié et attend son alignement après migration fonctionnelle.

Défauts connus encore au backlog :

- Choix des QCM non affichés dans les entraînements.
- Normalisation trop permissive de la casse et de certains espaces de commandes.
- Parcours Start training/Keyboard shortcuts redondants et footer compact tronqué
  au redimensionnement ; panneau d'aide complet inchangé, travaux P1.4 différés.
- `streak` non mis à jour, `last_session` sans curseur de reprise.
- Contrôles et diagnostics supplémentaires des terminaux indisponibles à traiter
  en P1 ; la capture n'est pas promise pour toutes les combinaisons macOS.

## Vision décidée et travaux suivants

**CLI-Drill — Learn the terminal by doing.** La [spécification](docs/SPEC.md)
organise un corpus vérifié en Unix & Shell Fundamentals, Keyboard Shortcuts,
Modern CLI & macOS, avec contexte portable/BSD/macOS/GNU/Linux/tiers explicite.
Les sessions visent cinq à dix minutes, sans détection d'outils installés.

La [décision V2](docs/DECISIONS.md#2026-11-07--cli-drill-v2--corpus-integre-et-retrait-des-dotfiles)
conserve intégralement les anciennes données JSON/YAML/TOML. Leur abandon
fonctionnel n'autorise aucune purge, réécriture ou migration automatique.
Les vérifications techniques du terminal restent requises ; les deux moteurs
et le fonctionnement local et hors ligne restent centraux.

P1.2 attend sa validation avant intégration. P1.3/P1.4, la correction des réponses,
les choix QCM, les parcours redondants et le footer restent hors de cette branche.
Les travaux suivants porteront sur ces défauts, puis sur le corpus pédagogique.
Voir le [backlog](TODO.md).

## Validation et limites de preuve

### P0.3 — Historique

Les contrôles de référence sur la base P0.2 ont réussi avant modification.
Après suppression : `go test -count=1 ./...`, `go vet ./...`,
`go run ./cmd/cli-drill --help`, `gofmt -w cmd internal data`, `go mod tidy`
et `git diff --check`. La cartographie des imports confirme l'absence des
packages retirés et du code de génération. Le corpus, `exercise`, `progress`
et les fonctions actives de rendu/capture ont été comparés à la base Git.

Deux tests de caractérisation ajoutés passent avant et après le découplage :
chargement des dix chapitres/68 exercices dans le même ordre avec leurs
métadonnées pédagogiques, et chemin JSON historique par défaut, avec surcharge
XDG et sans répertoire personnel disponible. Les tests de préservation des
anciennes données et de la surface V2 restent en place ; leurs chemins de
fixtures ne dépendent plus de helpers de production supprimés. Le test du pied
de page clavier vise désormais la méthode réellement utilisée par le rendu.

Le binaire a aussi été essayé dans un pseudo-terminal isolé à 80×24 : premier
lancement sans configuration, entraînement textuel immédiat, menu paginé,
progression, capture réelle de `Ctrl+A`, sortie `Esc`, retour à l'accueil et
restauration des modes du terminal. Les commandes `chapters`, `train`, `stats`
et `reset --yes` ont été exécutées sur des fichiers temporaires distincts.
L'historique JSON version 1 est chargé et conservé lors de l'ajout de résultats.
Les anciens TOML/YAML synthétiques, même invalides, sont ignorés et leurs octets,
ainsi que ceux du JSON historique de référence, restent identiques après les essais.
Aucune vraie donnée utilisateur n'est lue ou modifiée pour ces validations.

La comparaison des attributs du pseudo-terminal exclut seulement `PENDIN`,
[drapeau d'état macOS](https://github.com/apple-oss-distributions/xnu/blob/main/bsd/sys/termios.h)
lié à l'entrée en attente. Ce même écart est reproduit avec le binaire de la base
`14b12bd` ; tous les autres attributs sont restaurés. La capture n'a pas été
modifiée pour masquer cet écart dans le code applicatif.

Ces validations ne couvrent ni la V2 complète, ni tous les émulateurs et
raccourcis. La compatibilité avec le minimum Go 1.23 n'est pas vérifiée ;
les contrôles de développement utilisent Go 1.27.3 sur macOS arm64.

### P1.1 — Validations initiales, antérieures à la revue

Tests de référence sur P0.3 puis contrôles après modification :
`go test -count=1 ./...`, `go vet ./...`, formatage Go et `git diff --check`.
Les tests ajoutés ont d'abord exposé l'écrasement des JSON invalides, l'acceptation
de versions non prises en charge, la troncature de l'ancien inode et la création
en `0644`, puis réussi après correction.

La suite couvre absence de fichier, création des répertoires, aller-retour du
JSON historique et de ses cinq champs, valeurs par défaut, erreurs de lecture,
sérialisation et écriture, permissions et refus des liens symboliques.
Des erreurs provoquées après une écriture partielle et sur un fichier fermé
vérifient l'intégrité de la destination et le nettoyage du temporaire ; un
renommage réellement refusé vers un répertoire teste aussi ce nettoyage.
Un refus de nettoyage par permissions est signalé avec l'erreur d'écriture.
Les tests de permissions Unix sont exécutés sur macOS avec un utilisateur non
root. Les tests de stockage sont compilés pour Linux amd64, sans exécution Linux.

Le binaire a été essayé avec un HOME et un XDG isolés : aide, `chapters`,
`stats`, six réponses de `train`, conservation des résultats retirés et du mode
`0640`, nouvelle progression en `0600`, refus des JSON invalides et reset avec
confirmation. Le lancement et la sortie de la TUI n'écrasent pas un JSON invalide.
Aucun consommateur, YAML, moteur clavier, règle pédagogique ou dépendance ne
change. Toutes les données de test sont synthétiques et temporaires.

### P1.1b — Corrections après revue, acceptées et intégrées localement

Les tests ajoutés avant correction ont reproduit l'écrasement des `null`
invalides et des fichiers `0400`/`0444`. Les tests applicatifs passent par
`Load`, les commandes `stats`/`train` et le refus direct de `Save`, avec contrôle
des octets d'origine après chaque action. Cinq fixtures positives couvrent les
omissions historiques, conteneurs `null`, scores partiels et un JSON v1
complet enrichi par un entraînement, sans perte des résultats retirés.

Les tests de remplacement atteignent séparément les erreurs injectées de `Sync`
et de `Close`, vérifient l'intégrité de la destination et le nettoyage, ainsi
que l'ordre nominal écriture/synchronisation/fermeture avant renommage.
Des suppressions de temporaires réellement refusées par les permissions
vérifient aussi la remontée conjointe de l'erreur initiale et du nettoyage.
Les injections sont exercées en parallèle sur des fichiers isolés.

Contrôles réussis sur macOS arm64, Go 1.27.3, utilisateur non root (UID 501) :
`go test -count=1 ./...`, `go vet ./...`, formatage et `git diff --check` ;
`go test -race -count=1 ./internal/progress ./internal/app` réussit également.
Les tests de permissions ne sont pas sautés. L'aide CLI et le refus de `reset`
sans confirmation puis `reset --yes` sont vérifiés sur une fixture temporaire.

Cinq mutations par overlays temporaires font échouer les nouveaux tests :
suppression de la validation structurelle, du contrôle d'écriture, de `Sync`,
de `Close`, et fermeture déplacée après renommage. Les sources du checkout ne
sont pas modifiées par ces essais. Aucune donnée utilisateur réelle n'est utilisée.
Ni Linux, ni Go 1.23, ni une nouvelle session terminal interactive ne sont
validés par P1.1b ; les limites de durabilité et de concurrence restent celles
décrites plus haut. La revue des corrections est acceptée et leur intégration
locale est autorisée. La finalisation de P1.1 comprend une dernière exécution
des tests complets, de `go vet ./...` et de `git diff --check`, puis un commit
des corrections et un fast-forward exclusivement vers `dev`. Aucun push ni
travail P1.2 ne fait partie de cette finalisation.

### P1.2 — Validation de l'enregistrement, à valider humainement

Avant modification : arbre propre, deux commits P1.1 présents, fetch de `origin/dev`
et égalité avec `dev` à `1f61cb9`, tests complets sans cache et vet réussis.
L'index QMD ne renvoyait pas les règles de tentative ; les documents originaux
du checkout ont servi de référence. Les défauts ont été reproduits avant correction.

Les nouveaux tests couvrent les quatre types textuels par un même parcours et
le moteur clavier séparément : réponses correctes/incorrectes, enregistrement
avant la lecture suivante, sessions complètes/partielles, EOF vide et non vide,
Esc/Ctrl+C, aides et solutions, retries et révision, feedback et Entrées répétées,
nouvelles séances et rechargement. Ils passent par les commandes Cobra ou les
événements Bubble Tea, avec lecture du JSON réel après soumission. Les résultats
historiques retirés, les cinq champs v1, `stats` et `reset --yes` sont contrôlés.
Les erreurs de stockage sont exercées sur des liens synthétiques refusés par
P1.1, et la TUI est testée après correction de la cause et nouvelle validation.

Les tests échouent avec les sources de `1f61cb9` en overlay temporaire. Deux
mutations supplémentaires sont détectées : suppression de la sauvegarde commune
et suppression de la garde de feedback TUI. Aucun fichier utilisateur réel
n'est utilisé, et aucun code d'exercice n'est exécuté.

Contrôles finaux réussis sur macOS arm64 avec Go 1.27.3 : formatage des seuls
fichiers Go modifiés, `go test -count=1 ./...`, `go vet ./...`, `git diff --check`
et `go test -race -count=1 ./internal/app ./internal/tui ./internal/progress`.

Le binaire est aussi essayé dans un pseudo-terminal macOS isolé : Start training,
Browse chapters textuel, Browse chapters clavier, Keyboard shortcuts et retour
à l'accueil, erreur de sauvegarde TUI visible, interruption CLI par SIGINT et
abandon CLI par Esc suivi d'Entrée. Les sept scénarios vérifient le JSON et la
sortie du processus. Cela ne valide pas un émulateur macOS particulier, les
raccourcis interceptés par le système ni tous les états du terminal.

Les garanties contre coupure électrique, les écritures concurrentes et les
arrêts avant déclenchement de la sauvegarde restent exclues. Aucun curseur de
reprise ni stockage des réponses en échec n'est ajouté. Linux et Go 1.23 ne sont
pas validés par ces essais ; P1.3/P1.4 restent hors périmètre.

Pour une vérification complémentaire dans un terminal natif, depuis la racine
du dépôt (procédure fournie, non exécutée dans un émulateur natif) :

```sh
recording_dir="$(mktemp -d)"
go run ./cmd/cli-drill --progress "$recording_dir/progress.json"
go run ./cmd/cli-drill --progress "$recording_dir/progress.json" stats
```

Dans la TUI : Start training, saisir `zoxide`, Entrée, Esc, `g`, puis `q`.
Les statistiques doivent annoncer un exercice acquis et un chapitre pratiqué.
Relancer avec le même `--progress`, vérifier My progress, puis pratiquer le même
chapitre via Browse chapters : une nouvelle réponse ajoute une tentative,
sans dupliquer l'acquis. Keyboard shortcuts, `Ctrl+A`, puis Esc doit enregistrer
la réponse avant le retour à l'accueil. Tous ces essais restent dans le dossier
temporaire explicite.
