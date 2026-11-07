# État du projet — CLI-Drill

**Constat du 7 novembre 2026, base de code `c66f329` sur `dev`.**
La vision V2 est décidée ; sa migration fonctionnelle n'est pas réalisée.
Cette livraison est documentaire et ne change ni le binaire, ni les chapitres,
ni les dépendances, ni les données utilisateur.

## Implémentation actuelle

Le MVP CLI/TUI existant reste orienté dotfiles :

- Module Go à la racine, point d'entrée `cmd/cli-drill/` et packages dans `internal/`.
- Dix chapitres YAML intégrés et 68 exercices, dont huit en capture réelle.
- Chargement des YAML utilisateur à la place du corpus intégré dès que le
  dossier de chapitres utilisateur contient des fichiers `*.yaml`.
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

Ces limites ne sont pas corrigées par la présente livraison documentaire.
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

1. Retirer la surcharge automatique du corpus, les commandes/écrans/options
   dotfiles puis les packages et la dépendance TOML devenus inutilisés.
2. Fiabiliser le stockage, l'enregistrement partagé des résultats, les QCM,
   la correction des commandes et les événements clavier.
3. Réviser les contenus personnels, fournir un premier chapitre textuel et
   organiser le corpus en trois familles, sans réutiliser des clés pour un
   apprentissage différent.
4. Adapter la navigation en réutilisant les composants existants et vérifier
   les limites du terminal ; conserver le moteur de capture dédié.
5. Mettre à jour le README après migration et validation fonctionnelles.

L'[architecture cible](docs/ARCHITECTURE.md) conserve six packages internes.
Aucun retrait de code, nouvelle navigation ni nouvelle leçon n'est livré ici.

## Validation et limites de preuve

L'audit préalable, exécuté sur la même base de code, a validé `go test ./...`,
`go test -count=1 -cover ./...`, `go vet ./...`, l'aide CLI, le contrôle de
formatage Go et `go mod verify`, avec Go 1.27.3 sur macOS arm64.
Des essais de TUI en pseudo-terminal ont confirmé les écarts de progression,
de QCM et de filtre ; ce constat ne valide pas tous les terminaux réels.

Ces résultats ne sont pas une validation fonctionnelle V2 et ne garantissent
pas la compatibilité avec Go 1.23. La mission documentaire vérifie son diff,
les liens internes, la cohérence et le périmètre de fichiers ; elle ne relance
pas automatiquement les tests Go et n'exécute aucune migration utilisateur.
