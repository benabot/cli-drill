# État du projet — CLI-Drill

**Constat du 7 novembre 2026, base `dev` : `b82dbaf` (P0.1 intégrée).**
P0.2 est réalisée sur `refactor/remove-dotfiles-surface`, en attente de validation
et d'intégration. La migration V2 complète reste à réaliser.

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
- Progression JSON version 1 et comportement de sauvegarde inchangés.

### P0.1 — Corpus intégré exclusif, intégrée dans `dev`

Le chargeur partagé `loadChapters` lit uniquement le système de fichiers
embarqué transmis par `data.Chapters()`. Les anciens YAML ne remplacent ni ne
complètent le corpus et ne sont ni déplacés, ni convertis, ni supprimés.
Un chapitre présent seulement dans le dossier utilisateur est indisponible.

### P0.2 — Surface dotfiles retirée, à valider

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

## Internes et limites encore présents

Les packages `detect`, `catalog`, `config`, `markdown` et `shell/zsh`, leurs
modèles et tests restent présents et inchangés pour P0.3. `chapter` conserve
ses dépendances aux types du catalogue, les fonctions de génération et de
conversion ; `xdg` conserve les anciens helpers de chemins. La dépendance TOML
est encore déclarée. Ces éléments ne sont plus accessibles par la surface CLI/TUI.

Les contenus YAML conservent des alias, fonctions et formulations personnelles
héritées ; leur révision pédagogique est une étape distincte. Le README public
n'est pas modifié et attend son alignement après migration fonctionnelle.

Défauts connus encore au backlog :

- Bilan textuel TUI en mémoire sans persistance des réponses.
- Choix des QCM non affichés dans les entraînements.
- Normalisation trop permissive de la casse et de certains espaces de commandes.
- Sauvegarde JSON directe, sans remplacement atomique.
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

Le prochain chantier est `refactor/remove-dotfiles-internals` : retirer les
packages devenus inutiles, découpler les métadonnées et nettoyer les dépendances.
Ensuite viennent la fiabilité, le corpus pédagogique et les ajustements de
navigation réellement nécessaires après révision des contenus. Aucun travail
P0.3, nouvelle leçon, nouvelle dépendance ou refonte graphique n'est réalisé ici.
Voir le [backlog](TODO.md) et l'[architecture cible](docs/ARCHITECTURE.md).

## Validation et limites de preuve

Sur la base P0.1, les tests de référence ont réussi sans réutilisation des
résultats. Les nouvelles régressions P0.2 ont échoué avant les changements,
puis réussi après. Les tests de l'ancienne surface ont été retirés ; les tests
des packages conservés et des moteurs d'entraînement restent en place.

Contrôles P0.2 : `go test ./...`, `go vet ./...`, aide du binaire,
contrôle du formatage et `git diff --check`. Les régressions couvrent les
commandes/options retirées, le menu, le démarrage textuel, les indices de la
vue clavier, les filtres, les dimensions et la préservation des fixtures héritées.

Le binaire a aussi été essayé dans un pseudo-terminal isolé : premier lancement
sans configuration, navigation, entraînement CLI, historique JSON, capture de
`Ctrl+A`, sortie `Esc`, retour à l'accueil et restauration de l'état du terminal.
Les vérifications de données utilisent uniquement des fichiers synthétiques ;
les sessions d'entraînement et le test explicite de reset visent des fichiers
temporaires distincts. Aucune vraie donnée utilisateur n'est lue ou modifiée.

Ces validations ne couvrent ni la V2 complète, ni tous les émulateurs et
raccourcis. La compatibilité avec le minimum Go 1.23 n'est pas vérifiée ;
les contrôles de développement utilisent Go 1.27.3 sur macOS arm64.
