# Backlog — CLI-Drill V2

**Mis à jour le 10 novembre 2026.** P0.1 à P0.3 et P1.1 à P2.4 sont terminées.
P2.5 est en cours. Le stockage initial (`25c344b`) est complété par les
corrections acceptées de `fix/progress-storage-hardening` (P1.1b).
Migration fonctionnelle V2 partielle. P1.2 est terminée et validée : implémentation
`e7b4170` et complément de tests `a60c819`, depuis `dev` à `1f61cb9`.
Son intégration et sa publication sur `dev` sont réalisées à `be2bf2b`.
P1.3 est terminée, validée et intégrée dans `dev` à `db40d00`.
P1.4 est terminée et validée après recette réelle dans Terminal.app sous macOS.
P1 est clôturée et intégrée dans `dev` à `ef9b0a5`. P2.0 est terminée et
approuvée le 8 novembre 2026, intégrée et publiée à `91f13c6`.
P2.1 est terminée et validée le 8 novembre 2026 sur `feat/pedagogical-metadata`,
base `91f13c6`. La finalisation autorise commit, intégration par fast-forward
dans `dev` et publication de `origin/dev`, en conservant la branche de travail.
P2.1 est intégrée et publiée à `ab2096c`. P2.2 est intégrée et publiée à
`3ecdf9c` après huit essais manuels dans Terminal.app le 9 novembre 2026.
P2.3 est intégrée et publiée à `1e0c396` après six essais manuels le même jour.
P2.4 est intégrée et publiée sur `origin/dev` à `f096eec` après huit tests
manuels réussis dans Terminal.app.

## P0 — Simplification du périmètre

- [x] Formaliser la vision V2 et la distinguer de l'implémentation dans les documents.
- [x] Charger exclusivement les chapitres intégrés pour CLI/TUI et leur liste,
  sans charger les anciens YAML utilisateur ; conserver intégralement leurs fichiers.
- [x] Retirer la lecture et les diagnostics du TOML dotfiles de la surface CLI/TUI
  (P0.2), sans toucher aux anciens fichiers utilisateur.
- [x] Retirer `init`, `scan`, `generate`, `directory`, `search`, `show`,
  `--config` et les écrans Configuration, Scan et Annuaire.
- [x] Découpler les métadonnées de `chapter` du catalogue ; retirer génération,
  conversion vers l'annuaire, packages abandonnés et dépendance TOML (P0.3 intégrée).
  Conserver le corpus, les moteurs, les contrôles terminal et le chemin JSON existant.
- [x] Conserver les vérifications techniques TTY, dimensions, couleurs et capture
  clavier ; ne pas les confondre avec la détection d'outils ou de dotfiles.
- [x] Tester sur des données synthétiques que le changement de chargement ne
  supprime, ne déplace, ne réécrit ni ne purge les anciennes données JSON/YAML/TOML.
  Le premier lancement et les consultations ne doivent provoquer aucune migration écrite.

## P1 — Fiabilité des entraînements — terminée et validée

Clôture validée le 8 novembre 2026 après réussite des huit tests manuels P1.4
dans Terminal.app sous macOS, sans défaut bloquant ou important. Les compléments
non réalisés sont conservés explicitement dans la section différée ci-dessous.

- [x] Sauvegarder atomiquement la progression ; conserver les champs et résultats
  historiques, refuser l'écrasement silencieux d'un JSON invalide ou non pris en charge
  (P1.1 terminée, corrections P1.1b acceptées). Temporaire voisin, Sync, fermeture et renommage ;
  nouveaux fichiers privés et modes POSIX existants conservés. Un seul écrivain,
  sans garantie contre une coupure électrique ni copie des ACL/attributs étendus.
  P1.1b traite exclusivement les `null` invalides, l'ouverture des cibles en écriture
  sans troncature et les tests d'erreurs et d'ordre des opérations `Sync`/`Close`.
- [x] Partager la règle d'enregistrement CLI/TUI ; persister chaque réponse
  réellement soumise, sans doublon lors du feedback, d'un retour ou d'une sortie
  (P1.2 terminée et validée). Sauvegarde avant feedback ; erreur visible et réponse
  TUI conservée pour une nouvelle tentative de sauvegarde.
- [x] Ne pas créer de tentative pour l'aide, une solution affichée, la navigation
  ou un EOF vide ; tester retries, EOF après réponse et interruption
  (P1.2 terminée et validée). Validations sans texte ignorées ; compteurs v1 inchangés.
- [x] P1.4 : afficher et sélectionner les choix des cinq QCM dans la TUI avec
  flèches/j/k, Entrée, feedback et explication ; réutiliser correction P1.3 et
  enregistrement P1.2. Tests du corpus, interruptions, sauvegarde unique et erreurs
  réussis ; recette réelle dans Terminal.app validée.
- [x] P1.4 : lancer directement Keyboard shortcuts lorsqu'un seul chapitre de
  capture existe ; réutiliser le moteur dédié et conserver le navigateur sinon.
  Start training et Browse chapters conservent leurs comportements existants.
  Essais sous pseudo-terminal et recette réelle dans Terminal.app réussis.
- [x] Recette Terminal.app préalable rapportée par l'utilisateur : bannière,
  palette, navigation, aide complète, footer compact aux dimensions testées,
  redimensionnement, thèmes, progression et capture validés. Pas de refonte de layout.
- [x] Corriger les commandes sans effacer casse, options ou espaces significatifs ;
  garder les variantes explicites et la normalisation des raccourcis
  (P1.3 terminée et validée). Routage commun par `chapter.Item.MatchAnswer` ;
  espaces normalisés uniquement pour mots/pipes simples, autres commandes
  littérales, noms pédagogiques nus tolérants à la casse et raccourcis dédiés.
  Corpus et progression inchangés ; variantes, QCM, faux positifs, parcours
  CLI/TUI et sept mutations testés. Les équivalences complexes non déclarées
  et changements de casse des libellés de workflow ne sont pas déduits.
  Après revue : collage multiligne refusé avant Bubbles (terminal et `Ctrl+V`),
  sans tentative ni écriture ; reprise après édition et tabulations de bord
  testées dans les deux interfaces, raccourcis inchangés.
- [x] Donner priorité à la saisie des réponses et des filtres avant les raccourcis globaux (navigation P0.2).
- [x] Valider la capture réelle, les interruptions, la restauration du terminal et
  le retour TUI dans Terminal.app sous macOS. Les compléments relatifs aux terminaux
  indisponibles et aux cas d'erreur restent différés.
- [x] Démarrer directement un chapitre textuel depuis Start training, sans configuration
  ni activation implicite de capture brute (P0.2).

<a id="p2--feuille-de-route-pedagogique"></a>

## P2 — Feuille de route pédagogique

La [décision P2.0](docs/DECISIONS.md#2026-11-08--p20--evolution-pedagogique)
complète la V2. L'audit pédagogique des 68 exercices est déclaré validé par la
demande du 8 novembre 2026 : appliquer ses orientations, sans refaire cet audit
ni rouvrir P1. Les contrats YAML, filtres, ordre et préférences sont définis dans
la [spécification](docs/SPEC.md), qui fait référence pour les détails.

Ordre de réalisation **P2.0 → P2.1 → P2.2 → P2.3 → P2.4 → P2.5 → P2.6**.
Chaque sous-phase fonctionnelle fait l'objet d'une mission et d'une validation
distinctes.

### P2.0 — Cadrage documentaire — terminée et validée

Dépendance : P1 validée, base `ef9b0a5`.

- [x] Définir le périmètre terminal/CLI, les modules manuels, les trois niveaux,
  le mélange et l'exception guidée ; conserver les décisions historiques.
- [x] Décrire les métadonnées additives, les préférences séparées et les garanties P1.
- [x] Ordonner les sous-phases et préciser leur acceptation.

Acceptation : documents cohérents, liens locaux valides, diff limité aux références
Markdown autorisées et `git diff --check` réussi. Aucun Go, YAML, dépendance,
README ou fichier de configuration créé/modifié. Après validation documentaire,
la finalisation du 8 novembre 2026 autorise commit, fast-forward et publication
de ce seul cadrage ; aucun développement P2.1 n'est engagé.

### P2.1 — Métadonnées pédagogiques rétrocompatibles — terminée et validée

Dépendance : P2.0. Ne pas encore activer les filtres sur le corpus non annoté.

**Statut du 8 novembre 2026 : terminée et validée par l'utilisateur sur
`feat/pedagogical-metadata`, base `91f13c6`.**
Les champs, défauts et contrôles ci-dessous sont implémentés, testés et approuvés.
La validation de P2.2 est décrite ci-dessous. Tests ciblés et complets sans cache, vet, formatage Go
et diff contrôlés ; détails et limites dans [PROJECT_STATE.md](PROJECT_STATE.md).
Le corpus P2.1 de 10 YAML/68 exercices reste conservé comme fixtures de référence
dans P2.2. Les clés et mécanismes P1 sont préservés et aucun
filtre n'est activé.

- [x] Ajouter `difficulty` et `tools` à l'exercice, `family` et `learning_order`
  au chapitre selon SPEC ; conserver les champs historiques et les deux moteurs.
- [x] Valider les valeurs, listes de modules et identifiants sans nouveau catalogue.

Acceptation : les fixtures YAML historiques sans difficulté se chargent comme
Beginner, sans réécriture ni changement d'identifiants ; l'ordre absent reste
prudemment guidé. Valeurs explicites invalides, types inconnus, doublons de clés
et incohérences de QCM sont diagnostiqués. Les tests couvrent les trois niveaux,
les outils multiples, l'ancien format et le corpus intégré inchangé. Aucun import
automatique des YAML utilisateur, aucune modification du JSON de progression v1.

### P2.2 — Premier corpus V2 jouable et annoté — terminée et validée

Dépendance : P2.1. Appliquer l'audit validé, sans nouvelle revue exhaustive.

**Recette humaine validée le 9 novembre 2026 dans Terminal.app.** Les huit
essais ont réussi, dont le parcours Unix, les QCM CLI, la capture et la relecture
de la progression après interruption. Le corpus P2.2 comptait 39 exercices
dans neuf chapitres. Les tests sur fixtures confirment la conservation des résultats
historiques JSON v1. La durée des mini-leçons reste à confirmer par l'usage.
[Corpus et identifiants](docs/CORPUS_P2.md) · [Détail des validations](PROJECT_STATE.md)

- [x] Fournir un premier lot court dans les trois familles, avec un chapitre
  textuel Beginner du socle Unix utilisable sans outil tiers ni configuration.
- [x] Retirer du corpus actif les dépendances aux alias/fonctions personnels,
  workflows dotfiles et questions sur les logiciels graphiques hors périmètre.
  Le cas `Marked 2` est résolu par ce retrait, sans changer le correcteur P1.3.
- [x] Déclarer famille et ordre du chapitre, difficulté de chaque exercice et
  tous les modules tiers requis ; rendre chaque consigne autonome.
- [x] Documenter pour chaque clé conservée, nouvelle ou retirée le maintien ou
  changement d'apprentissage. Réorganiser les familles sans renommer les clés
  historiques pertinentes ni transférer de faux acquis.
- [x] Lever le complément différé d'affichage des choix QCM en CLI texte avant
  de déclarer ces leçons jouables dans les deux interfaces.

Acceptation : chaque famille est jouable ; les chapitres visent cinq à dix minutes.
Consignes, solutions, variantes, explications et portabilité sont vérifiées dans
les références officielles utiles. Les scénarios fournissent leur contexte sans
état réel ni réponse préalable ; aucune commande saisie n'est exécutée.
Tous les exercices tiers sont annotés avant d'activer leur sélection. Premier
lancement et consultations laissent les fixtures historiques JSON/YAML/TOML
intactes ; les entraînements enrichissent le JSON sans perdre les résultats retirés.
Essais pédagogiques CLI texte, TUI et capture compatibles réalisés,
y compris choix QCM visibles et diagnostics des captures indisponibles.

### P2.3 — Sélection manuelle et préférences minimales — terminée et validée

Dépendance : P2.2, pour ne pas exposer des exercices tiers non annotés.
Implémentation sur `feat/cli-preferences` depuis `3ecdf9c`, validée le
9 novembre 2026 dans Terminal.app et publiée sur `dev` à `1e0c396`.
Voir les [commandes et le parcours TUI](PROJECT_STATE.md#p23--sélection-manuelle-et-préférences-minimales-terminée-et-validée).

- [x] Permettre d'activer séparément les modules tiers issus du corpus et de
  choisir le niveau par défaut ; conserver le socle Unix toujours accessible.
- [x] Ajouter le stockage facultatif de préférences décrit dans SPEC, distinct
  de la progression, avec des fonctions concrètes dans les packages existants.
- [x] Raccorder les choix CLI/TUI aux composants existants, sans refonte de
  navigation ; conserver `--config` et les commandes de scan hors de l'interface.

Acceptation : sans préférences, le socle Beginner démarre sans création de
préférences et sans installation tierce. Chaque module peut être activé/désactivé
sans scan ; un exercice multi-outils exige tous ses modules. Une modification
enregistrée est relue après relance ; un choix ponctuel n'est pas persisté. Absence, JSON invalide,
version inconnue, outil retiré et échec d'écriture suivent les diagnostics et
protections de SPEC. Aucun changement d'outils/niveau ne modifie la progression
v1 ou les anciens YAML/TOML ; `--progress` et `reset --yes` gardent leur périmètre.
Le prédicat d'éligibilité par module est disponible pour la préparation commune
de séance P2.4 ; aucun filtre n'est appliqué aux séances dans P2.3.

- [x] Recette humaine du 9 novembre 2026 : six tests Terminal.app réussis
  (défauts sans fichier, activation CLI, persistance, annulation et sauvegarde
  TUI, désactivation/navigation), sans anomalie signalée.

### P2.4 — Séances filtrées, guidées et mélangées — terminée

Implémentée sur `refactor/session-order` depuis `1e0c396`, puis intégrée et
publiée sur `origin/dev` à `f096eec`. Voir le [fonctionnement livré et les
limites](PROJECT_STATE.md#p24--séances-filtrées-guidées-et-mélangées).

- [x] Permettre de choisir Beginner, Intermediate ou Advanced au lancement de la
  séance, y compris depuis une fiche chapitre ; filtrer par niveau exact.
- [x] Former une liste éligible puis mélanger une seule fois par nouvelle séance
  ordinaire/révision ; garder l'ordre en mémoire jusqu'à la sortie.
- [x] Pour un chapitre guidé, proposer Apprendre en ordre éditorial et Réviser
  en ordre aléatoire, sans marqueur persistant de première séance.
- [x] Adapter les décomptes et accès Start training / Keyboard shortcuts aux
  exercices éligibles, sans refonte visuelle ; calculer les acquis actifs par
  intersection d'identifiants, sans purger l'historique.

- [x] Valider les huit tests manuels dans Terminal.app.

Acceptation : pour zéro, un ou plusieurs exercices, chaque liste est une
permutation exacte de l'ensemble éligible, sans omission ni doublon. Deux tirages
identiques sont permis ; les tests vérifient l'invariant, pas une différence
systématique. Aide, feedback, retour de capture, redimensionnement, erreur de
sauvegarde puis nouvelle tentative ne changent pas l'ordre ou la clé enregistrée.
Les filtres et préférences restent figés pendant la séance. Une sélection vide
explique la cause sans tentative ni écriture ; le socle reste consultable.
Apprendre conserve l'ordre filtré et Réviser remélange ; retry/révision des erreurs
ajoutent uniquement les réponses réellement soumises selon P1.2. Les tests et
essais couvrent CLI, TUI et capture, avec conservation des résultats v1 retirés,
compteurs cumulés et absence de faux acquis. Le corpus source ne change pas d'ordre.

### P2.5 — Enrichissement progressif vers 100 à 120 exercices — en cours

Les parcours P2.2 à P2.4 sont validés. Le [lot 1 Unix/Shell](docs/CORPUS_P25_LOT1.md)
ajoute 18 exercices dans trois chapitres et relit les 39 exercices existants.
Le premier lot est terminé et validé après la recette finale dans Ghostty
sur macOS le 10 novembre 2026 : T01 à T08 sont réussis, dont T07. Les 57
exercices sont acceptés et le lot est prêt pour commit sur sa branche.
Son intégration et sa publication restent à réaliser.

- [ ] Enrichir les trois familles par petits lots relus et essayés, avec les
  trois niveaux progressivement représentés ; priorité aux commandes, scénarios
  concrets et captures prises en charge.
- [ ] Viser 100 à 120 exercices intégrés au total, cible indicative subordonnée
  à la qualité ; ne pas multiplier les questions de noms de logiciels.

Acceptation par lot : nouveaux identifiants pour nouveaux apprentissages,
métadonnées complètes, variantes explicites, explications et contexte vérifiés,
autonomie en ordre aléatoire et parcours guidé essayés. Décompte du corpus par
famille/niveau et équilibre des formats examinés ; aucun quota de remplissage.
Un écart à la cible est expliqué plutôt que compensé par des doublons.
Aucune installation tierce nécessaire, aucune exécution pédagogique et aucune
extension de schéma de progression, moteur adaptatif ou répétition espacée.

### P2.6 — Recette d'ensemble et documentation du logiciel livré

Dépendance : sous-phases fonctionnelles précédentes validées.

- [ ] Vérifier les parcours avec/sans modules, niveaux, guidage/révision,
  interruptions et sauvegardes, avec et sans couleur ; conserver le retour TUI.
- [ ] Actualiser README, état du projet et aides sur le comportement réellement
  livré, après validation fonctionnelle ; distinguer la cible indicative atteinte
  des lots encore prévus et conserver les limites terminal documentées.

Acceptation : recette macOS des parcours touchés et tests ciblés de compatibilité
sur fixtures synthétiques réussis ; exemples et aides cohérents avec le binaire,
liens et diff vérifiés. Les limites non testées sont nommées. Cette recette
porte sur P2 et ses régressions possibles, pas une nouvelle revue exhaustive P1.

## P3 — Internationalisation FR/EN

Phase future après P2.6. Permettre à l'utilisateur de choisir et de conserver
sa langue entre le français et l'anglais. Le français est le défaut envisagé,
à confirmer en P3.

- [ ] Localiser les menus CLI et TUI, les consignes et explications, le feedback,
  les diagnostics, l'aide et les textes de navigation.
- [ ] Persister le choix de langue en préservant la compatibilité des préférences
  existantes. Décider en P3 de la stratégie de localisation et du format de stockage.
- [ ] Garder les clés `chapter-id/item-id` indépendantes de la langue et préserver
  les acquis et scores historiques lors de tout changement de langue, sans
  suppression ni réattribution de résultat.
- [ ] Vérifier l'équivalence pédagogique des deux versions. Les commandes et
  raccourcis gardent leur signification et aucune traduction ne facilite la réponse.

## Compléments différés après clôture de P1

Ces éléments du backlog restent ouverts ; ils ne font pas partie de la livraison
P1.4 validée et ne sont pas implémentés par sa finalisation.

P2.2 valide l'affichage des choix QCM en CLI texte. P2.4 implémente la projection
des acquis du corpus actif. `streak` reste hors des indicateurs actifs.

- [ ] Compléter les vérifications des terminaux indisponibles et des cas d'erreur
  de capture/restauration ; ne jamais créditer une touche non réellement capturée.
- [ ] Guider vers la CLI textuelle si la TUI est indisponible ; compléter les diagnostics
  de capacités terminal sans changer les données utilisateur.

## Plus tard — Sous réserve d'un besoin observé

- Distribution par releases et éventuellement Homebrew, après stabilisation.
- Import personnalisé et reprise persistante au curseur : toujours différés.
- Préférences au-delà des modules CLI, du niveau par défaut et du choix de langue
  prévu en P3 : différées.

L'enrichissement avancé, le mélange des séances et les préférences pédagogiques
minimales sont désormais dans P2. Moteur adaptatif et répétition espacée restent
hors de cette phase.

Aucune roadmap de moteurs Bash/Fish, d'IA, MCP, compte, synchronisation,
scan runtime ou capture globale macOS n'est entretenue pour le nouveau MVP.

<a id="branches-atomiques-proposees"></a>

## Branches atomiques proposées

P0.1 à P0.3 et P1.1, corrections de `fix/progress-storage-hardening` comprises,
sont terminées et intégrées dans `dev`. P1.2 est terminée et validée sur
`fix/training-recording`, avec `e7b4170` et `a60c819` ; intégration et publication
sur `dev` réalisées à `be2bf2b`. P1.3 est terminée et validée sur
`fix/exercise-correction`, intégrée dans `dev` à `db40d00`. P1.4 est terminée et
validée dans Terminal.app et intégrée à `ef9b0a5` ; P1 est clôturée. P2.0 cadre
la feuille de route pédagogique. P2.1 à P2.4 sont intégrées et publiées.
Chaque nouvelle branche partira de `dev` à jour,
portera un objectif cohérent et une PR vers `dev`, jamais directement `main`.
Les sous-phases P2 suivent les dépendances ci-dessus ; chaque future PR de
code doit compiler et conserver les commandes d'entraînement disponibles.
Les lignes P0/P1 ci-dessous sont historiques, sans travail à recommencer.

| Priorité / branche | Périmètre et acceptation | Validation, risque et dépendance |
| --- | --- | --- |
| P0 — `refactor/embedded-only` (intégrée) | Corpus intégré seul pour l'entraînement et sa liste ; YAML non chargés et données héritées préservées. | Tests et essais isolés CLI/TUI réussis ; risque de changement du corpus visible caractérisé ; surface TOML retirée par P0.2. |
| P0 — `refactor/remove-dotfiles-surface` (intégrée) | Retirer commandes, option et écrans dotfiles en conservant l'entraînement et les contrôles terminal. | Aide CLI et navigation ; rupture intentionnelle des anciennes commandes ; après corpus intégré seul. |
| P0 — `refactor/remove-dotfiles-internals` (intégrée) | Six packages internes ; `chapter` sans catalogue ni génération, anciens chemins supprimés, TOML retiré ; corpus et données utilisateur préservés. | Tests sans cache, vet, aide, imports et dépendances vérifiés ; essais CLI/TUI isolés et restauration terminal réussis ; intégrée dans `dev` à `8ec5c00`. |
| P1 — `fix/progress-storage` (intégrée localement) | Sauvegarde atomique et historique v1 ; aucun changement du comptage ou des appels CLI/TUI. | `25c344b` ; trois problèmes de revue corrigés par P1.1b, P1.1 terminée. |
| P1.1b — `fix/progress-storage-hardening` (intégrée localement) | Refus des `null` invalides et des cibles interdites en écriture ; injection privée de `Sync`/`Close` par appel. | Revue acceptée ; tests applicatifs, permissions non root, erreurs et ordre de `Sync`/`Close`, nettoyage et cinq mutations vérifiés ; P1.1 terminée, aucun travail P1.2. |
| P1.2 — `fix/training-recording` (terminée et validée) | Enregistrement commun par réponse via les API P1.1 ; CLI texte/clavier et TUI texte persistés ; pas de tentative vide ni de double compte sur feedback. | Revue acceptée ; identifiants exacts, vrais messages de sortie et échec suivi d'interruption vérifiés ; tests complets, vet et race réussis ; aucun travail P1.3/P1.4. |
| P1.3 — `fix/exercise-correction` (terminée et validée) | Correction partagée des cinq types, variantes YAML conservées ; aucun changement d'affichage QCM (P1.4). | Tests de casse, espaces, syntaxe, choix, raccourcis et résultats CLI/TUI ; sept mutations détectées ; refus intentionnel de faux positifs ; base `be2bf2b`, stockage P1.1/P1.2 inchangé. |
| P1.4 — `fix/terminal-training` (terminée et validée) | QCM TUI sélectionnables ; accès direct au seul chapitre clavier via la capture existante. Aucun changement de layout, YAML, correction ou stockage. | Tests complets, vet, formatage et diff réussis ; essais sous pseudo-terminal et huit tests manuels dans Terminal.app réussis. Aucun défaut bloquant ou important ; clôture de P1 validée. |
| P2.0 — `docs/p2-roadmap` (terminée et validée) | Spécification, décisions et backlog P2, documentation seulement. | Liens, cohérence, périmètre et diff ; après P1 à `ef9b0a5`. |
| P2.1 — `feat/pedagogical-metadata` (terminée et validée) | Métadonnées additives, défauts historiques et validation approuvés. | Base `91f13c6` ; tests réussis, corpus inchangé ; finalisation par commit, fast-forward et publication autorisée. |
| P2.2 — `content/v2-first-corpus` (terminée et validée) | Premier corpus révisé et annoté avec 39 exercices et choix QCM visibles en CLI. | Base `ab2096c`. Huit essais manuels réussis dans Terminal.app le 9 novembre 2026. |
| P2.3 — `feat/cli-preferences` (intégrée et publiée) | Modules manuels et préférences séparées. | Six essais Terminal.app réussis le 9 novembre 2026. Publiée sur `dev` à `1e0c396`. |
| P2.4 — `refactor/session-order` (terminée) | Préparation commune, niveau, mélange, guidage et projection des acquis. | Huit essais Terminal.app validés. Intégrée et publiée sur `origin/dev` à `f096eec`. |
| P2.5, lot 1 — `content/v2-expansion-unix-01` (terminé et validé) | Trois chapitres Unix/Shell et révision éditoriale du corpus existant. | Base `f096eec`. T01 à T08 validés dans Ghostty. Prêt pour commit ; intégration et publication à réaliser. [Détails](docs/CORPUS_P25_LOT1.md). |
| P2.5, lots suivants — `content/v2-expansion-<lot>` | Enrichissement vers la cible indicative de 100 à 120 exercices. | Critères P2.5 et validation par lot. |
| P2.6 — `docs/readme-v2` | Présentation du comportement effectivement validé. | Critères P2.6 ; après recette P2 ; risque de promesse prématurée. |
| P3 — branche à définir | Choix de langue FR/EN et continuité des acquis. | Après P2.6. [Périmètre et garanties](#p3--internationalisation-fren). |

Les futures PR de code exécuteront les tests pertinents, `go test ./...`,
`go vet ./...`, le contrôle de formatage et l'aide CLI concernée. Les tests
utilisateur emploieront des fixtures synthétiques isolées, jamais les vraies
configurations. La livraison documentaire précédente n'a pas modifié Go.
P0.1 valide le chargement intégré, P0.2 la surface utilisateur et P0.3 le retrait
des internes hérités. P1.1, corrections P1.1b comprises, termine la fiabilisation
du stockage dans son périmètre, sans modifier les entraînements.
P1 est terminée et validée. La recette P2.6 reste à faire après l'enrichissement
du corpus en P2.5.
Les compléments différés restent explicitement ouverts.

## Archive V1 — Réalisations et abandons

La V1 a livré le module Go, Cobra, les dix chapitres intégrés, la TUI,
les moteurs textuel/clavier, le scan ZSH/Markdown, le catalogue, la génération,
le TOML et le stockage JSON. Leur état actuel est décrit dans
[l'état du projet](PROJECT_STATE.md) ; les décisions historiques restent dans
[DECISIONS.md](docs/DECISIONS.md).

Sont sortis du backlog actif : amélioration du scan, déduplication du catalogue,
génération depuis les dotfiles, préférences de configuration dotfiles et support
de nouveaux parseurs. Les corrections de sécurité propres au scanner proposées
par l'audit précédent sont remplacées par son retrait ; la sécurité des données
et du terminal reste obligatoire. L'ancienne préparation du README bilingue est
réalisée ; sa prochaine mise à jour attendra la V2. Les doublons de relecture
pédagogique, retour après capture et migration Bubble Tea ont été retirés.
Makefile/Taskfile, refonte visuelle et unification des moteurs ne constituent
pas des objectifs actifs.
