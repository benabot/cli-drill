# État du projet — CLI-Drill

**Constat du 7 novembre 2026, base `dev` : `563fdd3`.**
P0.1 est réalisée sur `refactor/embedded-only`, en attente de validation et
d'intégration. La migration fonctionnelle V2 complète reste à réaliser.
Ce changement porte uniquement sur la source des chapitres ; contenus, identifiants,
dépendances, moteurs d'entraînement et données utilisateur sont conservés.

## Implémentation actuelle

Le MVP CLI/TUI existant reste orienté dotfiles :

- Module Go à la racine, point d'entrée `cmd/cli-drill/` et packages dans `internal/`.
- Dix chapitres YAML intégrés et 68 exercices, dont huit en capture réelle.
- Chargement exclusif du corpus intégré pour la liste des chapitres,
  les entraînements CLI et les chapitres fournis à la TUI, même si des YAML
  utilisateur valides ou invalides existent.
- Scan statique de sources ZSH et Markdown configurées, parsing d'alias et noms
  de fonctions, catalogue dédupliqué et génération de chapitres YAML.
- Configuration TOML dotfiles ; commandes `init`, `scan`, `generate`, `chapters`,
  `train`, `directory`, `search`, `show`, `stats`, `reset`.
- Options globales `--config` et `--progress`.
- CLI d'entraînement textuel et moteur dédié de capture des touches de contrôle.
- TUI avec accueil, navigateur et fiche chapitre, entraînement textuel,
  annuaire, statistiques et écrans d'instructions Configuration / Scan.
- Passage de la TUI au moteur de capture, puis retour automatique à l'accueil.
- JSON local version 1 avec acquis uniques et tentatives cumulées par chapitre.
- Contrôles techniques existants pour la capture et les couleurs ; rendu
  CLI monochrome selon le terminal, `TERM` et `NO_COLOR`.

### Étape P0.1 — Source des chapitres

Le chargeur partagé `loadChapters` lit uniquement le système de fichiers
embarqué transmis par `data.Chapters()`. Le dossier utilisateur n'est plus
recherché ni chargé. Un chapitre présent seulement dans ce dossier n'est plus
disponible pour `train`, et un YAML utilisateur invalide ne bloque plus le lancement.

`generate` reste disponible et peut toujours écrire des YAML, dans son dossier
habituel ou via `--out` ; ces sorties ne sont plus automatiquement activées.
Les fichiers existants ne sont ni déplacés, ni convertis, ni supprimés.
`--config` ne choisit pas le corpus pédagogique. Les diagnostics TOML de la TUI
et les commandes de scan restent présents pour les prochaines étapes P0.

L'annuaire CLI sans configuration partage ce chargeur : il utilise lui aussi
le corpus intégré. Avec configuration, son catalogue reste issu du scan,
comme auparavant. Aucun lecteur ou package des fonctions dotfiles n'est supprimé ici.

Les cinq types d'exercices existent dans le modèle : `free-answer`, `scenario`,
`multiple-choice`, `simple-shell-sim`, `key-sequence`. Le simulateur compare
les réponses sans exécuter les commandes. Le scan n'exécute pas les définitions.

### Limites constatées dans l'audit de cette base

- L'entraînement textuel TUI calcule un bilan en mémoire sans enregistrer la
  progression JSON. La CLI enregistre sa progression en fin de session.
- Les choix YAML des QCM ne sont affichés ni par la CLI ni par la TUI.
- La correction textuelle normalise toutes les réponses en minuscules, y compris
  les commandes sensibles à la casse.
- Les raccourcis globaux peuvent intercepter la saisie dans les filtres des listes.
- La progression est écrite directement ; aucune sauvegarde atomique n'est implémentée.
- `streak` reste à zéro ; `last_session` est une date, sans curseur de reprise.
- Le premier chapitre actuel utilise la capture brute ; un premier parcours
  textuel indépendant de cette capacité reste à définir dans le corpus V2.
- La protection actuelle du scan ne résout pas les liens symboliques et ses
  exclusions sont configurables. Le scan sera retiré, pas développé davantage.

Ces limites ne sont pas corrigées par P0.1.
Le [README public](README.md) conserve la présentation du logiciel disponible.

## Vision V2 décidée

**CLI-Drill — Learn the terminal by doing.**

La [spécification V2](docs/SPEC.md) définit un entraînement local par la pratique :
corpus écrit et vérifié, saisie de commandes, vraie capture clavier compatible,
correction utile et progression persistante. Les sessions visent cinq à dix
minutes, sans configuration initiale ni dépendance à l'environnement personnel.

Trois familles : Unix & Shell Fundamentals, Keyboard Shortcuts, Modern CLI &
macOS. La portabilité, les variantes BSD/macOS et GNU/Linux et les préconditions
shell/terminal seront explicites. Les outils tiers restent des contenus,
sans détection ni obligation d'installation.

Le corpus cible est exclusivement intégré. La migration préservera intégralement
les anciennes données utilisateur JSON, YAML et TOML, sans réécriture, déplacement
ou suppression automatiques ; seuls les usages abandonnés seront retirés.
L'entraînement pourra enrichir la progression sans purger les résultats historiques.

L'introspection abandonnée concerne dotfiles et outils installés. Les contrôles
techniques nécessaires à la TUI, aux couleurs et à la capture restent autorisés
et requis avant d'utiliser ces capacités. Le premier chapitre cible sera textuel ;
un terminal limité recevra un diagnostic et un accès CLI textuel adapté.

La [décision du 7 novembre 2026](docs/DECISIONS.md#2026-11-07--cli-drill-v2--corpus-integre-et-retrait-des-dotfiles)
remplace les anciennes orientations dotfiles pour les travaux à venir, tout en
préservant leur historique.

## Travaux non encore réalisés

Le [backlog V2](TODO.md) décrit les futures branches vers `dev` :

1. Après validation de P0.1, retirer les commandes/écrans/options dotfiles,
   puis les packages et la dépendance TOML devenus inutilisés.
2. Fiabiliser le stockage, l'enregistrement partagé des résultats, les QCM,
   la correction des commandes et les événements clavier.
3. Réviser les contenus personnels, fournir un premier chapitre textuel et
   organiser le corpus en trois familles, sans réutiliser des clés pour un
   apprentissage différent.
4. Adapter la navigation en réutilisant les composants existants et vérifier
   les limites du terminal ; conserver le moteur de capture dédié.
5. Mettre à jour le README après migration et validation fonctionnelles.

L'[architecture cible](docs/ARCHITECTURE.md) conserve six packages internes.
P0.1 ne livre ni retrait des commandes ou packages dotfiles, ni nouvelle
navigation, ni nouvelle leçon. Le premier chapitre et le moteur de capture
restent inchangés ; le premier parcours textuel V2 appartient aux étapes suivantes.

## Validation et limites de preuve

L'audit préalable, exécuté sur la même base de code, a validé `go test ./...`,
`go test -count=1 -cover ./...`, `go vet ./...`, l'aide CLI, le contrôle de
formatage Go et `go mod verify`, avec Go 1.27.3 sur macOS arm64.
Des essais de TUI en pseudo-terminal ont confirmé les écarts de progression,
de QCM et de filtre ; ce constat ne valide pas tous les terminaux réels.

Pour P0.1, les nouvelles régressions ont échoué avant la modification du chargeur,
puis réussi après. Elles couvrent CLI et transmission à la TUI, remplacement et
ajout de YAML utilisateur, YAML invalide, refus d'un chapitre utilisateur seul,
lecture de progression version 1 et préservation des fixtures JSON/YAML/TOML.
La boucle terminal est substituée dans le test unitaire du lancement TUI ;
les contrôles interactifs ci-dessous ont été effectués séparément.

Contrôles P0.1 exécutés : `go test ./...`, `go vet ./...`,
`go run ./cmd/cli-drill --help`, contrôle du formatage et `git diff --check`.
Le binaire a aussi été essayé dans un pseudo-terminal isolé : lancement sans
configuration malgré des YAML utilisateur invalides, lecture des acquis
historiques, ouverture d'un exercice textuel, capture de `Ctrl+A`, sortie par
`Esc`, retour à l'accueil et restauration des attributs du terminal.
Un entraînement CLI textuel de six exercices a été terminé avec succès.

Seuls des fichiers synthétiques temporaires ont été utilisés ; leur JSON,
YAML et TOML hérités sont restés identiques à l'octet près. Les entraînements
de vérification ont écrit leurs résultats dans des JSON de session distincts,
également temporaires. Aucune vraie donnée utilisateur n'a été lue ou modifiée.

Ces résultats ne valident ni la V2 complète, ni tous les émulateurs/raccourcis,
ni la compatibilité avec Go 1.23. Les défauts de progression textuelle TUI,
de QCM et de correction restent dans le backlog P1.
