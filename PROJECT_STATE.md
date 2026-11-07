# État du projet — CLI-Drill

**Constat du 7 novembre 2026 : P0.1 à P0.3 et P1.1 terminées, intégrées localement dans `dev`.**
P1.1 comprend le stockage initial (`25c344b`) et les corrections de
`fix/progress-storage-hardening`, acceptées après revue. La migration V2
complète reste à réaliser ; P1.2 n'est pas commencée.

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
- Progression JSON version 1, chemin et déclenchements de sauvegarde inchangés ;
  écriture atomique et erreurs de stockage renforcées par P1.1.

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

## Limites encore présentes

Les contenus YAML conservent des alias, fonctions et formulations personnelles
héritées ; leur révision pédagogique est une étape distincte. Le README public
n'est pas modifié et attend son alignement après migration fonctionnelle.

Défauts connus encore au backlog :

- Bilan textuel TUI en mémoire sans persistance des réponses.
- Choix des QCM non affichés dans les entraînements.
- Normalisation trop permissive de la casse et de certains espaces de commandes.
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

P1.1 terminée, les prochaines étapes portent sur les autres travaux
P1, puis le corpus pédagogique et les ajustements de navigation nécessaires.
`fix/training-recording` n'est pas commencé : ni comptage des réponses, ni
calcul des scores, ni persistance textuelle TUI ne sont modifiés ici.
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
