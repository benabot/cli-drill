# Backlog — CLI-Drill V2

**Mis à jour le 13 novembre 2026.** P0.1 à P0.3 et P1.1 à P2.4 sont terminées.
Le corpus P2.5 est publié. P2.6 est terminée et publiée sur `dev` à `91efab0`
après recette macOS avec Ghostty. La [feuille de route de lancement](#lancement-public-p3-p6)
prépare les travaux P3.1 à P6. Le stockage initial (`25c344b`) est complété par les
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

<a id="p25--enrichissement-par-chapitre-et-par-niveau--en-cours"></a>

### P2.5 — Enrichissement par chapitre et par niveau — corpus publié

Les parcours P2.2 à P2.4 sont validés. Le [lot 1 Unix/Shell](docs/CORPUS_P25_LOT1.md)
ajoute 18 exercices dans trois chapitres et relit les 39 exercices existants.
Le premier lot est terminé et validé après la recette finale dans Ghostty
sur macOS le 10 novembre 2026 : T01 à T08 sont réussis, dont T07. Les 57
exercices sont acceptés. Le lot 1 est intégré et publié sur `dev` à `b809151`.
Le lot 2 ajoute 18 exercices, soit 75 au total. Il est terminé et validé
après la recette finale sous Ghostty sur macOS le 10 novembre 2026.
Il est intégré et publié sur `dev` à `4320650`.
[Rapport du lot](docs/CORPUS_P25_LOT2.md)

- [x] Vague 1 « Repères, fichiers et droits » publiée sur `dev` à
  `f97718a` : 46 ajouts dans `00` et `12`, 121 au total.
  [Rapport de vague](docs/CORPUS_P25_WAVE1.md).
- [x] Vague 2 « Flux, processus et environnement » publiée sur `dev` à
  `fd37fe4`. Les 48 ajouts dans `11` et `13` portent le corpus à 169 exercices.
  [Rapport de vague](docs/CORPUS_P25_WAVE2.md).
- [x] Vague 3 intégrée à `dev` à `a6af7b4`. Les 23 ajouts portent `01` à
  12 Beginner, 11 Intermediate et 8 Advanced. La base compte 192 exercices.
  [Rapport de vague](docs/CORPUS_P25_WAVE3.md)
- [ ] Valider la révision pédagogique de la vague 3.
- [x] Intégrer `14-recherche-texte-unix` à `dev` à `0178ab9`. Ses 48
  exercices portent le corpus publié à 240 exercices dans 13 chapitres.
  [Rapport](docs/CORPUS_P25_UNIX_TEXT.md)
- [ ] Faire la recette du chapitre 14 et mesurer la durée d'une séance
  Intermediate de 22 exercices.
- [x] Intégrer les 13 indispensables du lot sécurité des fichiers et
  processus à `dev` à `c9ad59c`. Le corpus publié compte 253 exercices.
  [Rapport du lot](docs/CORPUS_P25_UNIX_SAFETY_PROCESSES.md)
- [x] Intégrer les neuf indispensables de navigation et flux à `dev` à
  `d8b383a`. Le corpus publié compte 262 exercices.
  [Rapport du lot](docs/CORPUS_P25_UNIX_BASICS_FLOWS.md)
- [x] Implémenter et tester les 14 propositions utiles de l'audit,
  les parcours Git et outils modernes et le recentrage de `07` à `09`.
  Le corpus de 446 exercices dans 13 chapitres est publié sur `dev` à `ef5d6ba`.
  [Rapports par jalon](docs/CURRICULUM.md#enrichissement-par-lots)
- [x] Vérifier des parcours représentatifs lors de la recette P2.6.
  Les preuves par scénario figurent dans le [rapport](docs/P26_FINAL_ACCEPTANCE.md).
  Les ACL et `lsof` restent facultatifs.

Le minimum indicatif est de 30 exercices pour un chapitre général. Le volume
et la répartition des niveaux suivent les compétences à couvrir. Le référentiel
conserve la planification initiale des lots et le repère historique de 360
exercices. [Priorités P2.5](docs/CURRICULUM.md#enrichissement-par-lots)

Acceptation par lot : nouveaux identifiants pour nouveaux apprentissages,
métadonnées complètes, variantes explicites, explications et contexte vérifiés,
autonomie en ordre aléatoire et parcours guidé essayés. Décompte du corpus par
chapitre et par niveau, équilibre des formats et progression examinés.
Une lacune est documentée et reste ouverte si la seule solution est un doublon
ou une difficulté artificielle.
Aucune installation tierce nécessaire, aucune exécution pédagogique et aucune
extension de schéma de progression, moteur adaptatif ou répétition espacée.

### P2.6 — Recette d'ensemble validée sur macOS

Clôture intégrée et publiée sur `dev` à `91efab0` après la recette macOS.

- Lot 1 (sécurité et capture) publié dans `dev` à `fb29a22`.
- Lot 2 (lisibilité et navigation) publié dans `dev` à `0b8bddb`.
  Les commandes restent visibles pendant la lecture paginée, et la CLI
  propose un parcours textuel si la TUI est indisponible.
- Lot 3 (séances courtes et progression) publié dans `dev` à `1f53b29`.
  [Contrat des lots et des reprises](docs/SPEC.md#répondre-comprendre-et-poursuivre).
- [x] Lot 4 : contrôles automatisés et essais sous pseudo-terminal sur macOS,
  compilation croisée Linux/BSD et documentation. Une sortie involontaire
  en CLI textuelle a été reproduite. [Rapport et preuves](docs/P26_FINAL_ACCEPTANCE.md).
- [x] Exécuter H01 à H11 sous macOS avec Ghostty. Apprendre a duré 1 min 20 s
  (6/6) et Réviser 1 min 58 s (4/6).
  [Observations physiques](docs/P26_FINAL_ACCEPTANCE.md#checklist-physique-sur-mac).
- [x] Valider P2.6 sur macOS. Le mainteneur a confirmé les flèches H08,
  Ctrl+B/F et le séparateur CLI après correction sur
  `fix/p26-physical-acceptance`. Les réserves H05 et H09 sont consignées
  dans le [rapport de clôture](docs/P26_FINAL_ACCEPTANCE.md#critères-de-clôture).
- La validation native sur Linux et BSD reste à faire. La compilation
  croisée a réussi pendant le lot 4.

Améliorations non bloquantes issues de la recette P2.6 :

- [ ] Revoir la difficulté des exercices Beginner jugés difficiles pendant la
  recette.
- [ ] Examiner l'ergonomie des QCM à partir des observations de recette.

<a id="lancement-public-p3-p6"></a>

## Dépendances du lancement public

Les objectifs sont fixés par la [décision P3.0](docs/DECISIONS.md#p30-politique-linguistique-et-lancement),
à partir de la [roadmap du 10 novembre](docs/CLI-DRILL_ROADMAP_LANCEMENT_GITHUB.md).
Les travaux ci-dessous restent à réaliser.

| Moment | Travaux et conditions |
| --- | --- |
| Préparation commune | Instruire la licence et les cibles de distribution, auditer les références V1, inventorier les textes et concevoir la localisation. |
| Réalisation en parallèle | P3.1 prépare P3.2 et P3.3. P4.1 stabilise les textes et titres à traduire avant P3.3. La chaîne GoReleaser, les packages, les maquettes README, la structure du Wiki et les scénarios de démo peuvent avancer pendant P3. |
| Version candidate | Valider les deux langues et l'équivalence pédagogique en P3.4, vérifier la cohérence en P4.2, tester les artefacts P5 et vérifier P6 sur les parcours et la version distribués. |
| Lancement | La release publique bilingue exige traductions et recettes validées, binaires installables et vérifiés, README fidèle, Wiki EN/FR, licence et contribution. Les annonces suivent un essai complet de découverte, d'installation et de première séance. Publication et communication sur autorisation explicite. |

Le choix entre bêta technique et sortie publique unique reste ouvert.
Une bêta apporte des retours plus tôt sur les plateformes et l'installation,
avec une diffusion et un suivi supplémentaires. Une sortie unique concentre
la préparation sur une version et reporte ces retours après sa publication.
Les critères de lancement restent les mêmes dans les deux cas.

<a id="p3--internationalisation-fren"></a>

## P3 — Internationalisation EN/FR

La [spécification](docs/SPEC.md#politique-linguistique-p3) fixe le contrat.
Anglais universel par défaut. Français et anglais sélectionnables dans la CLI
et la TUI, avec persistance du choix et conservation des réglages existants.
Même progression JSON v1 dans les deux langues.

### P3.1 — Architecture de localisation

- [ ] Inventorier les chaînes CLI/TUI, y compris aides, erreurs précoces et
  capture, ainsi que les champs du corpus. Distinguer textes affichés,
  identifiants et fragments techniques. S'appuyer sur l'audit P4.1.
- [ ] Choisir le format des ressources et les interfaces d'accès aux
  traductions. Définir la couverture vérifiable et le traitement des réponses
  conceptuelles et des QCM. Préparer un petit lot représentatif.
- [ ] Auditer le chargement, la validation et l'enregistrement des préférences.
  Comparer l'extension de `preferences.json` v1 et un passage justifié en v2.
  Décrire comment la nouvelle version lira les préférences valides et
  conservera les modules et la difficulté à la mise à niveau, selon la
  [décision du 13 novembre](docs/DECISIONS.md#p30-compatibilite-des-preferences-en-pre-release).
  Définir l'accès au choix de langue en CLI/TUI.

La conception est prête pour P3.2 et P3.3 lorsque ses choix et compromis
sont documentés et que la matrice sur fixtures synthétiques couvre l'absence
de fichier ou de champ, `en`, `fr`, les valeurs invalides, la sauvegarde
explicite et la mise à niveau. Elle vérifie la préservation des réglages
valides et l'absence d'écriture au lancement et à la consultation.

### P3.2 — CLI et TUI

- [ ] Localiser les menus, aides, navigation, réglages, feedback et diagnostics
  des deux interfaces et des deux moteurs. Rendre visible le choix de langue
  et permettre son enregistrement sans modifier les modules ni le niveau.
- [ ] Vérifier le démarrage en anglais quelles que soient les variables
  régionales, la sélection explicite, l'annulation, les erreurs de lecture
  ou d'écriture et la relance dans la langue enregistrée.

La recette vérifie que les deux interfaces répondent dans la langue choisie.
Les consultations et la lecture d'une préférence ancienne restent en lecture
seule. Une langue invalide produit un diagnostic sans réparation.

### P3.3 — Corpus pédagogique

- [ ] Après P3.1 et le nettoyage préparatoire P4.1, produire des lots par
  chapitre. Commencer par le lot représentatif, puis couvrir les 446
  exercices des 13 chapitres. Suivre les traductions par clé et par champ,
  avec une terminologie commune aux interfaces et au corpus.
- [ ] Relire les titres, descriptions, consignes, choix, corrections et
  explications. Vérifier par lot le contexte, les distracteurs, les réponses
  admises et la difficulté dans les deux langues. Protéger la syntaxe des
  commandes, les options et les raccourcis.

La revue humaine valide les formulations de chaque lot. Chaque clé identifie
le même apprentissage, avec la même bonne réponse conceptuelle et le même
niveau. Les commandes gardent leur syntaxe et les séquences capturées leur
signification.

### P3.4 — Validation bilingue

- [ ] Mettre en place des contrôles bloquants de complétude et de cohérence
  des ressources actives. Vérifier les champs attendus, la correspondance
  des clés et des réponses et l'intégrité des fragments techniques.
- [ ] Tester dans les deux langues le démarrage, l'aide, les préférences,
  Apprendre, Réviser, le feedback, la capture et les diagnostics, y compris
  les petites fenêtres. Reprendre la matrice P2.6 en recette humaine bilingue.
- [ ] Comparer le JSON v1 avant et après changement de langue sur des
  fixtures avec acquis et historique. Vérifier qu'une réussite dans une
  langue est reconnue dans l'autre sans nouvelle tentative.

La validation exige la couverture EN/FR complète des interfaces et du corpus
actif, la revue de l'équivalence pédagogique, la réussite des tests de
régression et la recette humaine bilingue.
Le passage à l'anglais par défaut et la commande effectivement livrée pour
choisir le français figurent dans les notes de version P5.

## P4 — Nettoyage des références V1

- [ ] **P4.1 — Audit et préparation avant P3.3.** Relever les textes visibles,
  aides CLI/TUI, titres et descriptions de chapitres, exemples et textes
  publics. Retirer les formulations obsolètes et harmoniser les titres
  concernés avant traduction. Le titre « Git essentiel » et la clé
  interne `10-workflows-dotfiles` illustrent le rôle distinct de ces champs.
- [ ] **P4.2 — Revue avant publication.** Aligner les aides, le README, le
  Wiki, la spécification et le logiciel sur la version candidate. Retirer
  des guides actifs les références aux commandes de scan, d'annuaire et
  de génération. Préserver décisions historiques, rapports, identifiants
  persistants et fichiers utilisateur.

La revue P4.1 vérifie les textes à traduire avant P3.3. En P4.2, les aides
et exemples doivent présenter CLI-Drill comme un entraînement au terminal
et correspondre aux parcours vérifiés.

## P5 — Distribution multiplateforme

- [ ] **P5.1 — Licence et cibles.** Choisir la licence avant diffusion
  publique officielle. Définir les versions minimales et la matrice de
  recette native. Préparer les cibles macOS ARM64/AMD64 et Linux AMD64/ARM64.
  Consigner la recette macOS P2.6 et la compilation croisée Linux/BSD
  comme preuves distinctes. Toute annonce BSD exige sa propre recette native.
- [ ] **P5.2 — Chaîne de release.** Préparer et éprouver GoReleaser et GitHub
  Actions sur une version candidate, même pendant P3. Partir d'un tag
  validé pour publier les archives sur GitHub Releases. Documenter tests Go,
  contrôles statiques, permissions minimales des workflows et accès au tap,
  artefacts, checksums SHA-256, notes de version et critères bloquants.
  Évaluer la provenance vérifiable des artefacts.
- [ ] **P5.3 — Installation.** Préparer le tap Homebrew et choisir formule
  ou Cask en vérifiant l'outillage retenu. Vérifier le chemin du module,
  le versionnement et l'installation distante avec `go install` avant
  de l'annoncer. Tester sur les cibles annoncées l'installation, la
  vérification d'intégrité, l'aide, la TUI et le premier entraînement
  hors ligne, ainsi que la mise à jour et la désinstallation.
- [ ] **P5.4 — Paquets Debian, complément.** Définir la matrice `.deb`,
  produire et tester les paquets de release installés depuis un fichier
  local. La distribution par un dépôt APT reste différée.

Le premier lancement exige une licence disponible, des artefacts et
checksums vérifiés, des notes cohérentes et une installation native
reproductible pour chaque cible annoncée. Les instructions décrivent
les canaux éprouvés.
La release publique bilingue dépend de P3.4 et de la convergence P4 à P6.

## P6 — README, Wiki et lancement GitHub

- [ ] **P6.1 — README.** Préparer `README.md` en anglais et `README.fr.md`
  en français, avec accès aux deux langues en tête. Conserver la
  [bannière](docs/assets/cli-drill-hero.webp) et la palette cyan, menthe et
  fuchsia. Associer le slogan **Make your terminal second nature.** au
  positionnement **Learn the terminal by doing.** Prendre Gum, zoxide et bat
  comme références de démonstration, de concision et de structure, avec
  les éléments graphiques du projet. Présenter les fonctionnalités vérifiées,
  l'installation validée, une première séance, les trois familles pédagogiques,
  le fonctionnement hors ligne et la capture réelle. Renvoyer vers le Wiki,
  les releases, la licence et la contribution pour les détails.
- [ ] **P6.2 — Démonstrations CleanShot X.** Préparer trois parcours sur la
  version candidate avec des données de démo. Capturer une séance générale
  pour le GIF principal du README, puis la capture réelle d'un raccourci et
  un entraînement avec réponse, explication et progression pour le Wiki.
  Montrer la fenêtre du terminal et les touches utiles. Viser 900 à
  1 100 px de large, 12 à 15 images/s et 8 à 20 s, avec un GIF principal
  sous les 5 Mo environ. Vérifier contraste, lisibilité et chargement sur
  GitHub desktop et mobile.
- [ ] **P6.3 — Wiki EN/FR.** Commencer par quelques pages utiles avec
  navigation entre langues. Couvrir installation, premiers pas, parcours,
  modes Apprendre et Réviser, types d'exercices, raccourcis, commandes CLI,
  préférences et langue, progression et statistiques, dépannage, FAQ et
  contribution. Lier les spécifications internes depuis le guide. À chaque
  release, relire les deux langues et les README sur le candidat, épingler
  les médias au tag ou commit correspondant, publier les mises à jour dans
  le dépôt Wiki distinct et vérifier les liens sur GitHub.
- [ ] **P6.4 — Présentation et retours.** Préparer About, Topics, Social
  Preview, badges vérifiables, licence et contribution. Préparer une
  présentation pour des communautés pertinentes et, si retenu, une page
  sur beabot.fr. Après autorisation et lancement, suivre téléchargements,
  trafic GitHub, issues et retours volontaires sur l'installation, les
  séances terminées et le retour à l'usage. Les étoiles restent un indicateur
  secondaire de visibilité. Le produit reste sans télémétrie.

La recette P6 suit le parcours d'un nouveau lecteur dans chaque langue.
Il doit pouvoir comprendre rapidement le produit, l'installer puis
l'essayer. Les README restent concis, les captures montrent le programme
réel, les commandes et liens sont vérifiés et le Wiki correspond à la release.

### Compléments de distribution et de documentation différés

RPM, dépôt APT signé, soumission à `homebrew/core`, autres canaux,
automatisation des captures et site documentaire GitHub Pages/MkDocs.
Les réexaminer selon les besoins observés. Les paquets `.deb` sont suivis
en P5.4.

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
| P2.5, lot 1 — `content/v2-expansion-unix-01` (intégré et publié) | Trois chapitres Unix/Shell et révision éditoriale du corpus existant. | T01 à T08 validés dans Ghostty. Publié sur `dev` à `b809151`. [Détails](docs/CORPUS_P25_LOT1.md). |
| P2.5, lot 2 — `content/v2-expansion-cli-02` (intégré et publié) | Outils CLI modernes, 18 exercices dans six chapitres. | Recette fonctionnelle et UX réussie sous Ghostty le 10 novembre 2026. Publié sur `dev` à `4320650`. [Détails](docs/CORPUS_P25_LOT2.md). |
| P2.5, vague 1 — `content/p25-wave1-paths-files` (publiée sur `dev`) | 46 exercices sur les repères, les fichiers et les droits ; 30 dans chacun des chapitres `00` et `12`. | Publiée à `f97718a` ; [contrôles et limites](docs/CORPUS_P25_WAVE1.md). |
| P2.5, vague 2 — `content/p25-wave2-flows-processes` (publiée sur `dev`) | 48 exercices sur les flux, processus et environnement. | Publiée à `fd37fe4`. [Rapport](docs/CORPUS_P25_WAVE2.md). |
| P2.5, vague 3 — `content/p25-wave3-line-editing` (intégrée sur `dev`) | 31 exercices d’édition de ligne dans Bash et Zsh. Apprentissage guidé et révision mélangée. | Intégrée à `a6af7b4`. [Contrôles et limites](docs/CORPUS_P25_WAVE3.md). |
| P2.5, chapitre 14 — `content/p25-unix-text` (intégré à `dev`) | 48 exercices sur la recherche et le traitement de texte Unix. | Intégré à `0178ab9`. [Contrôles et limites](docs/CORPUS_P25_UNIX_TEXT.md). |
| P2.5, sécurité des fichiers et processus — `content/p25-unix-safety-processes` (intégré à `dev`) | 7 exercices en `12` et 6 en `13`. | Intégré à `c9ad59c`. [Contrôles et limites](docs/CORPUS_P25_UNIX_SAFETY_PROCESSES.md). |
| P2.5, navigation et flux — `content/p25-unix-basics-flows` (intégré à `dev`) | 3 exercices en `00` et 6 en `11`. | Intégré à `d8b383a`. [Contrôles et limites](docs/CORPUS_P25_UNIX_BASICS_FLOWS.md). |
| P2.5, finalisation — `content/p25-unix-git-modern-expansion` (publiée sur `dev`) | 184 ajouts, dont 14 utiles Unix. | Publiée à `ef5d6ba`. [Contrôles par jalon](docs/CURRICULUM.md#enrichissement-par-lots). |
| P2.6 — `fix/p26-physical-acceptance` (publiée sur `dev`) | Recette validée sur macOS avec Ghostty. | Clôture à `91efab0`. [Rapport](docs/P26_FINAL_ACCEPTANCE.md) |
| P3.0 — `docs/launch-roadmap-p3-p6` | Politique linguistique et cadrage du lancement. | Base `91efab0`. [Dépendances et lots](#lancement-public-p3-p6). |

Les futures PR de code exécuteront les tests pertinents, `go test ./...`,
`go vet ./...`, le contrôle de formatage et l'aide CLI concernée. Les tests
utilisateur emploieront des fixtures synthétiques isolées, jamais les vraies
configurations. La livraison documentaire précédente n'a pas modifié Go.
P0.1 valide le chargement intégré, P0.2 la surface utilisateur et P0.3 le retrait
des internes hérités. P1.1, corrections P1.1b comprises, termine la fiabilisation
du stockage dans son périmètre, sans modifier les entraînements.
P1 est terminée et validée. Le mainteneur a validé P2.6 sur macOS.
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
et du terminal reste obligatoire. La préparation historique du README bilingue
est réalisée. La prochaine refonte publique relève de P6. Les doublons de relecture
pédagogique, retour après capture et migration Bubble Tea ont été retirés.
Makefile/Taskfile, refonte visuelle et unification des moteurs ne constituent
pas des objectifs actifs.
