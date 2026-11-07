# Backlog — CLI-Drill V2

**Mis à jour le 7 novembre 2026.** P0.1 intégrée dans `dev` (`b82dbaf`).
P0.2 réalisée sur `refactor/remove-dotfiles-surface`, en attente de validation
et d'intégration ; migration fonctionnelle V2 partielle.
Le code actuel reste décrit dans [PROJECT_STATE.md](PROJECT_STATE.md).
La [spécification](docs/SPEC.md) fixe le périmètre et les critères d'acceptation ;
la [décision V2](docs/DECISIONS.md#2026-11-07--cli-drill-v2--corpus-integre-et-retrait-des-dotfiles)
explique les abandons. Les cases cochées ne désignent que des travaux réalisés.

## P0 — Simplification du périmètre

- [x] Formaliser la vision V2 et la distinguer de l'implémentation dans les documents.
- [x] Charger exclusivement les chapitres intégrés pour CLI/TUI et leur liste,
  sans charger les anciens YAML utilisateur ; conserver intégralement leurs fichiers.
- [x] Retirer la lecture et les diagnostics du TOML dotfiles de la surface CLI/TUI
  (P0.2), sans toucher aux anciens fichiers utilisateur.
- [x] Retirer `init`, `scan`, `generate`, `directory`, `search`, `show`,
  `--config` et les écrans Configuration, Scan et Annuaire.
- [ ] Découpler les métadonnées de `chapter` du catalogue ; retirer génération,
  conversion vers l'annuaire, packages abandonnés et dépendance TOML.
- [x] Conserver les vérifications techniques TTY, dimensions, couleurs et capture
  clavier ; ne pas les confondre avec la détection d'outils ou de dotfiles.
- [x] Tester sur des données synthétiques que le changement de chargement ne
  supprime, ne déplace, ne réécrit ni ne purge les anciennes données JSON/YAML/TOML.
  Le premier lancement et les consultations ne doivent provoquer aucune migration écrite.

## P1 — Fiabilité des entraînements

- [ ] Sauvegarder atomiquement la progression ; conserver les champs et résultats
  historiques, refuser l'écrasement silencieux d'un JSON invalide ou non pris en charge.
- [ ] Partager la règle d'enregistrement CLI/TUI ; persister chaque réponse
  réellement soumise, sans doublon lors du feedback, d'un retour ou d'une sortie.
- [ ] Ne pas créer de tentative pour l'aide, une solution affichée, la navigation
  ou un EOF vide ; tester retries, EOF après réponse et interruption.
- [ ] Afficher les choix des QCM dans les deux interfaces et valider leur cohérence.
- [ ] Corriger les commandes sans effacer casse, options ou espaces significatifs ;
  garder les variantes explicites et la normalisation des raccourcis.
- [x] Donner priorité à la saisie des réponses et des filtres avant les raccourcis globaux (navigation P0.2).
- [ ] Vérifier capacités, capture, restauration du terminal et retour TUI, y compris
  en cas d'erreur ; ne jamais créditer une touche non réellement capturée.
- [ ] Calculer les acquis du corpus actif par identifiants présents, sans purger
  le JSON des anciens exercices. Garder `streak` hors des indicateurs actifs.
- [x] Démarrer directement un chapitre textuel depuis Start training, sans configuration
  ni activation implicite de capture brute (P0.2).
- [ ] Guider vers la CLI textuelle si la TUI est indisponible ; compléter les diagnostics
  de capacités terminal sans changer les données utilisateur.

## P2 — Corpus pédagogique et expérience utilisateur

- [ ] Fournir un corpus court jouable dans chaque famille : Unix & Shell Fundamentals,
  Keyboard Shortcuts, Modern CLI & macOS. Viser cinq à dix minutes par chapitre.
- [ ] Réviser chaque exercice conservé : consigne discriminante, correction,
  explication, contexte portable/BSD/macOS/GNU/Linux/tiers vérifié et absence
  de dépendance à une installation personnelle.
- [ ] Réviser la première leçon textuelle pour le corpus V2 ; garder saisie de commandes
  et vraie capture comme activités centrales, avec QCM complémentaires.
- [x] Adapter le menu cible en réutilisant navigateur, fiche et statistiques (P0.2) ;
  Keyboard shortcuts est une vue filtrée, pas un moteur supplémentaire.
- [ ] Conserver l'identité graphique, le mode sans couleur et le retour à l'accueil
  après capture ; ne pas ajouter de préférences de navigation.
- [ ] Relire et essayer les leçons après de vraies sessions avant d'élargir le corpus.
- [ ] Mettre à jour le README public seulement après la migration fonctionnelle.

### Évaluation du corpus existant à réaliser

| Chapitre actuel | Orientation V2, sans modification YAML dans la mission documentaire |
| --- | --- |
| `01-raccourcis-terminal` | Conserver les exercices de capture pertinents ; préciser shell, mode d'édition et limites du terminal. |
| `02-navigation-shell` | Conserver les notions générales et outils tiers utiles ; retirer les variantes d'alias personnels. |
| `03-alias-zsh`, `04-fonctions-zsh` | Retirer les leçons dépendant des noms et fonctions personnels ; ne pas les convertir automatiquement en questions générales. |
| `05-outils-quotidiens`, `06-recherche-fichiers-contenu` | Conserver et vérifier les usages généraux ; distinguer outils tiers et commandes portables. |
| `07-lecture-preview` | Garder les pratiques générales ; réévaluer les associations d'applications et fonctions personnelles. |
| `08-micro`, `09-markdown-glow` | Évaluer comme contenus tiers secondaires ; retirer les alias et fonctions dits détectés. |
| `10-workflows-dotfiles` | Retirer les workflows personnels ; extraire seulement les exercices généraux utiles, par exemple Git. |

Ne conserver un couple d'identifiants que si l'apprentissage reste identique.
Un nouveau concept reçoit de nouveaux identifiants ; les anciens résultats
restent stockés. L'encodage minimal des familles/contextes sera documenté dans
la PR de corpus, sans nouvelle infrastructure de contenu.

## Plus tard — Sous réserve d'un besoin observé

- Distribution par releases et éventuellement Homebrew, après stabilisation.
- Enrichissement mesuré des exercices et sujets avancés, après validation du corpus court.
- Sessions aléatoires, préférences, import personnalisé ou reprise au curseur :
  différés, sans tâches d'implémentation actives ni nouvelle configuration prévue.

Aucune roadmap de moteurs Bash/Fish, d'IA, MCP, compte, synchronisation,
scan runtime ou capture globale macOS n'est entretenue pour le nouveau MVP.

<a id="branches-atomiques-proposees"></a>

## Branches atomiques proposées

P0.1 est intégrée. P0.2 est réalisée sur `refactor/remove-dotfiles-surface`,
à valider avant commit, merge ou push. Les unités suivantes restent planifiées.
Chaque nouvelle branche partira de `dev` à jour,
portera un objectif cohérent et une PR vers `dev`, jamais directement `main`.
Les unités de retrait se suivent dans l'ordre indiqué ; chaque PR doit compiler
et conserver les commandes d'entraînement encore disponibles.

| Priorité / branche | Périmètre et acceptation | Validation, risque et dépendance |
| --- | --- | --- |
| P0 — `refactor/embedded-only` (intégrée) | Corpus intégré seul pour l'entraînement et sa liste ; YAML non chargés et données héritées préservées. | Tests et essais isolés CLI/TUI réussis ; risque de changement du corpus visible caractérisé ; surface TOML retirée par P0.2. |
| P0 — `refactor/remove-dotfiles-surface` (réalisée, à valider) | Retirer commandes, option et écrans dotfiles en conservant l'entraînement et les contrôles terminal. | Aide CLI et navigation ; rupture intentionnelle des anciennes commandes ; après corpus intégré seul. |
| P0 — `refactor/remove-dotfiles-internals` | Découpler `chapter`, retirer packages abandonnés et TOML, sans toucher aux données utilisateur. | Compilation, imports, tests, dépendances ; risque de couplages résiduels ; après retrait de la surface. |
| P1 — `fix/progress-storage` | Sauvegarde atomique, historique préservé et erreurs visibles. | JSON version 1, fichiers invalides, échec de remplacement ; risque de perte de données à prévenir ; après P0. |
| P1 — `fix/training-recording` | Même comptage et persistance CLI/TUI, sans tentative fictive ou double compte. | Réouverture, sortie, retries et EOF ; risque de sémantique des tentatives ; après stockage fiable. |
| P1 — `fix/exercise-correction` | QCM affichés et correction adaptée aux commandes et raccourcis. | Casse, espaces significatifs, variantes et choix ; risque de refus de faux positifs auparavant acceptés ; après P0. |
| P1 — `fix/terminal-training` | Vérifier les capacités, traiter les erreurs, préserver la capture dédiée et l'accès textuel sans configuration. | Terminal compatible/limité, sans couleur, restauration et retour TUI ; risque d'état terminal ; après enregistrement partagé. |
| P2 — `content/v2-foundations` | Réviser le corpus et ses trois familles ; fournir le premier chapitre textuel et préserver les identifiants pertinents. | Validation YAML, références officielles et essais pédagogiques ; risque de réattribution d'acquis ; après P1. |
| P2 — `refactor/v2-navigation` (à réévaluer après corpus) | Le menu et la vue clavier sont livrés par P0.2 ; limiter les futurs ajustements aux besoins démontrés du corpus révisé. | Navigation par familles si le corpus le requiert ; ne pas refaire l'accueil ni changer l'identité graphique. |
| P2 — `docs/readme-v2` | Décrire le comportement V2 effectivement livré. | Exemples, liens et comparaison avec le binaire ; risque de promesse prématurée ; après validation fonctionnelle V2. |

Les futures PR de code exécuteront les tests pertinents, `go test ./...`,
`go vet ./...`, le contrôle de formatage et l'aide CLI concernée. Les tests
utilisateur emploieront des fixtures synthétiques isolées, jamais les vraies
configurations. La livraison documentaire précédente n'a pas modifié Go.
P0.1 valide le chargement intégré et P0.2 la surface utilisateur ; aucune des deux
ne valide la V2 complète ni le retrait des packages internes.

## Archive V1 — Réalisations et abandons

La V1 a livré le module Go, Cobra, les dix chapitres intégrés, la TUI,
les moteurs textuel/clavier, le scan ZSH/Markdown, le catalogue, la génération,
le TOML et le stockage JSON. Leur présence actuelle est décrite dans
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
