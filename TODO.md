# Backlog — CLI-Drill V2

**Mis à jour le 8 novembre 2026.** P0.1 à P0.3 et P1.1 terminées, intégrées
localement dans `dev`. Le stockage initial (`25c344b`) est complété par les
corrections acceptées de `fix/progress-storage-hardening` (P1.1b).
Migration fonctionnelle V2 partielle. P1.2 est terminée et validée : implémentation
`e7b4170` et complément de tests `a60c819`, depuis `dev` à `1f61cb9`.
Son intégration et sa publication sur `dev` sont réalisées à `be2bf2b`.
P1.3 est terminée et validée sur `fix/exercise-correction` ; son intégration
par fast-forward et sa publication sur `dev` sont autorisées. P1.4 est la prochaine étape.
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
- [x] Découpler les métadonnées de `chapter` du catalogue ; retirer génération,
  conversion vers l'annuaire, packages abandonnés et dépendance TOML (P0.3 intégrée).
  Conserver le corpus, les moteurs, les contrôles terminal et le chemin JSON existant.
- [x] Conserver les vérifications techniques TTY, dimensions, couleurs et capture
  clavier ; ne pas les confondre avec la détection d'outils ou de dotfiles.
- [x] Tester sur des données synthétiques que le changement de chargement ne
  supprime, ne déplace, ne réécrit ni ne purge les anciennes données JSON/YAML/TOML.
  Le premier lancement et les consultations ne doivent provoquer aucune migration écrite.

## P1 — Fiabilité des entraînements

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
- [ ] P1.4 : afficher les choix des QCM dans les deux interfaces ; leur correction
  est testée indépendamment de l'affichage en P1.3.
- [ ] P1.4 : revoir les parcours redondants Start training/Keyboard shortcuts et
  la troncature du footer compact au redimensionnement ; conserver l'aide complète.
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
- [ ] Décider la tolérance de casse du nom d'application `Marked 2` : la saisie
  `marked 2` reste refusée par P1.3. Traiter ce point mineur en P2, sans rendre
  tous les workflows insensibles à la casse.
- [ ] Réviser la première leçon textuelle pour le corpus V2 ; garder saisie de commandes
  et vraie capture comme activités centrales, avec QCM complémentaires.
- [x] Adapter le menu cible en réutilisant navigateur, fiche et statistiques (P0.2) ;
  Keyboard shortcuts est une vue filtrée, pas un moteur supplémentaire.
- [ ] Conserver l'identité graphique, le mode sans couleur et le retour à l'accueil
  après capture ; ne pas ajouter de préférences de navigation.
- [ ] Relire et essayer les leçons après de vraies sessions avant d'élargir le corpus.
- [ ] Mettre à jour le README public seulement après la migration fonctionnelle.

### Évaluation du corpus existant à réaliser

| Chapitre actuel | Orientation V2, corpus inchangé jusqu'à la révision pédagogique |
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

P0.1 à P0.3 et P1.1, corrections de `fix/progress-storage-hardening` comprises,
sont terminées et intégrées dans `dev`. P1.2 est terminée et validée sur
`fix/training-recording`, avec `e7b4170` et `a60c819` ; intégration et publication
sur `dev` réalisées à `be2bf2b`. P1.3 est terminée et validée sur
`fix/exercise-correction` ; intégration et publication autorisées. P1.4 est
la prochaine étape ; les autres unités restent planifiées.
Chaque nouvelle branche partira de `dev` à jour,
portera un objectif cohérent et une PR vers `dev`, jamais directement `main`.
Les unités de retrait se suivent dans l'ordre indiqué ; chaque PR doit compiler
et conserver les commandes d'entraînement encore disponibles.

| Priorité / branche | Périmètre et acceptation | Validation, risque et dépendance |
| --- | --- | --- |
| P0 — `refactor/embedded-only` (intégrée) | Corpus intégré seul pour l'entraînement et sa liste ; YAML non chargés et données héritées préservées. | Tests et essais isolés CLI/TUI réussis ; risque de changement du corpus visible caractérisé ; surface TOML retirée par P0.2. |
| P0 — `refactor/remove-dotfiles-surface` (intégrée) | Retirer commandes, option et écrans dotfiles en conservant l'entraînement et les contrôles terminal. | Aide CLI et navigation ; rupture intentionnelle des anciennes commandes ; après corpus intégré seul. |
| P0 — `refactor/remove-dotfiles-internals` (intégrée) | Six packages internes ; `chapter` sans catalogue ni génération, anciens chemins supprimés, TOML retiré ; corpus et données utilisateur préservés. | Tests sans cache, vet, aide, imports et dépendances vérifiés ; essais CLI/TUI isolés et restauration terminal réussis ; intégrée dans `dev` à `8ec5c00`. |
| P1 — `fix/progress-storage` (intégrée localement) | Sauvegarde atomique et historique v1 ; aucun changement du comptage ou des appels CLI/TUI. | `25c344b` ; trois problèmes de revue corrigés par P1.1b, P1.1 terminée. |
| P1.1b — `fix/progress-storage-hardening` (intégrée localement) | Refus des `null` invalides et des cibles interdites en écriture ; injection privée de `Sync`/`Close` par appel. | Revue acceptée ; tests applicatifs, permissions non root, erreurs et ordre de `Sync`/`Close`, nettoyage et cinq mutations vérifiés ; P1.1 terminée, aucun travail P1.2. |
| P1.2 — `fix/training-recording` (terminée et validée) | Enregistrement commun par réponse via les API P1.1 ; CLI texte/clavier et TUI texte persistés ; pas de tentative vide ni de double compte sur feedback. | Revue acceptée ; identifiants exacts, vrais messages de sortie et échec suivi d'interruption vérifiés ; tests complets, vet et race réussis ; aucun travail P1.3/P1.4. |
| P1.3 — `fix/exercise-correction` (terminée et validée) | Correction partagée des cinq types, variantes YAML conservées ; aucun changement d'affichage QCM (P1.4). | Tests de casse, espaces, syntaxe, choix, raccourcis et résultats CLI/TUI ; sept mutations détectées ; refus intentionnel de faux positifs ; base `be2bf2b`, stockage P1.1/P1.2 inchangé. |
| P1 — `fix/terminal-training` | Vérifier les capacités, traiter les erreurs, préserver la capture dédiée et l'accès textuel sans configuration. | Terminal compatible/limité, sans couleur, restauration et retour TUI ; risque d'état terminal ; après enregistrement partagé. |
| P2 — `content/v2-foundations` | Réviser le corpus et ses trois familles ; fournir le premier chapitre textuel et préserver les identifiants pertinents. | Validation YAML, références officielles et essais pédagogiques ; risque de réattribution d'acquis ; après P1. |
| P2 — `refactor/v2-navigation` (à réévaluer après corpus) | Le menu et la vue clavier sont livrés par P0.2 ; limiter les futurs ajustements aux besoins démontrés du corpus révisé. | Navigation par familles si le corpus le requiert ; ne pas refaire l'accueil ni changer l'identité graphique. |
| P2 — `docs/readme-v2` | Décrire le comportement V2 effectivement livré. | Exemples, liens et comparaison avec le binaire ; risque de promesse prématurée ; après validation fonctionnelle V2. |

Les futures PR de code exécuteront les tests pertinents, `go test ./...`,
`go vet ./...`, le contrôle de formatage et l'aide CLI concernée. Les tests
utilisateur emploieront des fixtures synthétiques isolées, jamais les vraies
configurations. La livraison documentaire précédente n'a pas modifié Go.
P0.1 valide le chargement intégré, P0.2 la surface utilisateur et P0.3 le retrait
des internes hérités. P1.1, corrections P1.1b comprises, termine la fiabilisation
du stockage dans son périmètre, sans modifier les entraînements.
La V2 complète reste à valider ; les autres corrections P1 et la révision pédagogique restent
à réaliser.

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
