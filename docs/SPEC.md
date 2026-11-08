# CLI-Drill V2 — Spécification fonctionnelle

**Vision décidée le 7 novembre 2026 ; migration fonctionnelle partielle.**
Ce document décrit la cible V2 et ses critères, pas un état de livraison complet.
L'[état du projet](../PROJECT_STATE.md) distingue les travaux intégrés et ceux
qui attendent leur validation ; le [README](../README.md) conserve la présentation
historique du MVP antérieur, avec un avertissement de périmètre.
La [décision V2](DECISIONS.md#2026-11-07--cli-drill-v2--corpus-integre-et-retrait-des-dotfiles)
explique le changement de direction ; le [backlog](../TODO.md) organise sa réalisation.

## 1. Vision et public

**CLI-Drill — Learn the terminal by doing.**

CLI-Drill est une application Go CLI/TUI locale pour apprendre et mémoriser
les commandes Unix, les outils CLI utiles sous macOS et les raccourcis clavier.
Elle s'adresse aux personnes qui découvrent le terminal et aux développeurs ou
power users qui souhaitent consolider leurs pratiques.

L'apprentissage repose sur la mémorisation active : exercices courts,
progressifs, répétables, écrits et vérifiés. Saisir une commande et appuyer sur
un vrai raccourci sont des activités centrales. Les questions de connaissance
servent ces pratiques ; elles ne constituent pas à elles seules le produit.

Une session vise cinq à dix minutes par sa longueur éditoriale. Le nouveau MVP
n'impose ni minuteur, ni sélection aléatoire, ni préférence de durée.

## 2. Périmètre décidé

### Conserver

- Go, Cobra, Bubble Tea, Bubbles, Lip Gloss et l'identité graphique actuelle.
- Corpus YAML intégré au binaire, lisible et éditable dans les sources du projet.
- Moteur de correction et simulateur de commandes qui compare du texte.
- Capture réelle des raccourcis compatibles avec le terminal utilisé.
- Progression JSON persistante en local et tests unitaires.
- Commandes `chapters`, `train [chapter-id]`, `stats` et `reset --yes`.
- Option `--progress` et fonctionnement hors ligne à l'exécution.

### Retirer lors de la migration

- Scan des dotfiles et détection des outils installés.
- Parsing des alias et fonctions utilisateur ; génération automatique de leçons.
- Catalogue des éléments détectés et configuration TOML dédiée aux dotfiles.
- Commandes `init`, `scan`, `generate`, `directory`, `search` et `show`,
  ainsi que l'option `--config`.
- Écrans Configuration, Scan et Annuaire, et leurs états associés.
- Chargement automatique des chapitres YAML utilisateur.

Ces suppressions sont des ruptures CLI intentionnelles, à annoncer au moment
de la livraison fonctionnelle. Aucun catalogue de remplacement n'est prévu.
Les anciens fichiers utilisateur restent intégralement conservés ; voir la
[politique de progression et de migration](#7-progression-et-donnees-heritees).

L'abandon de l'introspection vise les dotfiles et les outils installés. Il ne
supprime pas les vérifications techniques du terminal nécessaires à la TUI,
aux couleurs, au redimensionnement ou à la capture clavier.

## 3. Organisation pédagogique

### A. Unix & Shell Fundamentals

Navigation, fichiers et répertoires, recherche, pipes, redirections,
permissions, processus, environnement shell, texte, archives et Git essentiel.
Les exercices de commande fournissent un contexte simulé suffisant : objectif,
noms de fichiers et résultat recherché. Ils ne dépendent pas du disque réel.

### B. Keyboard Shortcuts

Navigation dans la ligne, édition, historique et contrôle des processus.
Chaque leçon précise si le raccourci relève du shell et de son mode d'édition,
de l'émulateur de terminal ou de macOS. ZSH sous macOS sert de contexte shell
de référence ; aucun moteur de détection ou d'adaptation aux shells n'est requis.

Conserver la capture réelle pour les touches prises en charge. Les combinaisons
interceptées par macOS, indisponibles dans le terminal ou réservées à la sortie
ne sont pas promises comme capturables. Par exemple, le rôle pédagogique de
`Ctrl+C` doit être distingué de son rôle de sortie dans le moteur actuel.
Une explication ou question textuelle ne vaut pas validation d'une capture réelle.

### C. Modern CLI & macOS

Commandes macOS et outils tiers utiles : `rg`, `fd`, `fzf`, `bat`, `jq`,
`zoxide`, et autres outils retenus pour leur valeur pédagogique.
Leur installation n'est nécessaire ni pour charger une leçon ni pour la suivre.

### Règles éditoriales communes

- Identifier le contexte pertinent : Unix portable, BSD/macOS, GNU/Linux,
  outil tiers, shell/mode d'édition ou émulateur de terminal.
- Ne pas présenter une option GNU comme universellement disponible sous macOS.
- Vérifier les syntaxes et raccourcis dans leurs références officielles lors
  de la rédaction du corpus ; noter les préconditions dans la leçon.
- Poser une question discriminante, expliquer la solution et ses limites.
- Inclure de la saisie de commandes dans les chapitres Unix et outils qui s'y
  prêtent ; réserver les QCM à un rôle complémentaire.
- Préserver les identifiants des exercices dont le sens reste identique.

Ces familles organisent les contenus ; elles ne justifient pas trois moteurs.
Le nouveau MVP exige un corpus court jouable dans chaque famille, sans couvrir
immédiatement tout ce programme. Le format YAML actuel reste la base ; aucun
nouveau schéma complexe n'est nécessaire à cette phase documentaire.

## 4. Parcours utilisateur

### Premier lancement et premier entraînement

Après installation, `cli-drill` donne immédiatement accès à l'entraînement,
sans fichier de configuration, scan ni installation d'outils pédagogiques.
`Start training` et `train` sans identifiant choisissent le premier chapitre
intégré non vide dans l'ordre pédagogique. Ce premier chapitre est textuel,
afin de ne pas imposer la capture brute avant vérification de ses capacités.

Les capacités nécessaires à l'interface choisie doivent être vérifiées avant
son activation. Si la TUI est indisponible, fournir un diagnostic clair et une
commande directe vers un chapitre textuel. La CLI textuelle reste utilisable
sans couleur, ligne par ligne, avec une entrée exploitable. En l'absence d'une
telle entrée, signaler la limite sans créer de tentative fictive.

### Navigation TUI cible

| Action | Comportement |
| --- | --- |
| Start training | Démarrer le premier chapitre textuel disponible selon l'ordre pédagogique. |
| Browse chapters | Choisir un chapitre parmi les trois familles et consulter sa fiche. |
| Keyboard shortcuts | Réutiliser le navigateur, limité aux chapitres de raccourcis. |
| My progress | Consulter la progression locale avec les composants existants. |
| Quit | Quitter proprement. |

Les cinq entrées n'imposent pas cinq parcours ou composants indépendants.
Conserver la fiche chapitre et le passage au mode clavier dédié, puis le retour
automatique à l'accueil TUI. Ne pas ajouter de cible de retour configurable.
Conserver la palette, la bannière et le style existants.

### Répondre, comprendre et poursuivre

Présenter le contexte et la consigne, recevoir la réponse, afficher une
correction avec explication, puis permettre de poursuivre. Une action de
navigation ne devient jamais une réponse. Les touches imprimables doivent
atteindre la réponse ou le filtre actif avant les raccourcis globaux.

Le bilan distingue le résultat de la session et les compteurs historiques.
La révision des erreurs reste en mémoire pour la session ; aucun nouveau
stockage de session n'est requis. Une relance de chapitre repart au début,
tout en conservant les acquis et tentatives précédents. La reprise exacte au
curseur interrompu est différée.

## 5. Formats d'exercices

| Type YAML existant | Usage V2 |
| --- | --- |
| `simple-shell-sim` | Saisir une vraie syntaxe de commande dans un contexte simulé ; aucune exécution. |
| `key-sequence` | Appuyer réellement sur les touches compatibles avec la capture disponible. |
| `scenario` | Résoudre une situation concrète, notamment par une commande lorsque pertinent. |
| `free-answer` | Réponse textuelle courte et active. |
| `multiple-choice` | Choix affichés explicitement ; réponse textuelle correspondant à une solution acceptée. |

Conserver `id`, `title`, `description`, `items`, `exercise_type`, `prompt`,
`choices`, `answer.primary`, `answer.accepted` et `explanation` comme base
du format existant. Le retrait du catalogue doit découpler les métadonnées
héritées `type` et `source`, sans imposer une conversion des fichiers utilisateur.
Les décisions d'encodage minimal des familles et contextes seront prises dans
la PR de corpus ; elles ne conditionnent pas la présente livraison documentaire.

Les chapitres intégrés doivent être non vides, leurs identifiants de chapitre
uniques, et leurs identifiants d'exercice uniques dans un chapitre. Les types
doivent être reconnus et les choix cohérents avec la réponse attendue.

## 6. Correction et tentatives

- Comparer à une réponse principale et à des variantes explicitement admises.
- Pour les commandes, préserver la casse des noms, les options et les espaces
  significatifs, notamment entre guillemets. Ne pas tout convertir en minuscules.
- Pour les réponses conceptuelles, appliquer une tolérance limitée et explicable.
- Normaliser les notations des raccourcis sans confondre une lettre seule avec
  sa combinaison de contrôle.
- Ne pas chercher une équivalence shell générale, analyser le disque ou lancer
  une commande pour vérifier une réponse.
- Une réponse soumise, correcte ou incorrecte, compte comme une tentative.
  Une nouvelle réponse après retry en compte une autre.
- L'aide, l'affichage d'une solution, la navigation et un EOF sans réponse ne
  comptent pas. Un EOF après une réponse textuelle non vide peut soumettre
  cette réponse une seule fois.
- Les moteurs CLI et TUI utilisent les mêmes règles d'enregistrement. Un retour
  à l'écran, une validation du feedback ou une sortie ne doit pas doubler le compte.

<a id="7-progression-et-donnees-heritees"></a>

## 7. Progression et données héritées

La progression reste en JSON version 1 à son emplacement actuel :
`~/.local/share/cli-drill/progress.json`, ou sous `XDG_DATA_HOME` lorsqu'il
est défini, ou au chemin explicite `--progress`.

Conserver les champs existants : `version`, `completed_exercises`,
`chapter_scores`, `last_session` et `streak`. Les clés d'acquisition restent
`chapter-id/item-id` ; les scores contiennent les tentatives et réponses correctes
cumulées. Un exercice réussi plusieurs fois reste un acquis unique.
`last_session` reste une date et ne promet pas une reprise exacte.
`streak`, actuellement non mis à jour, ne devient pas un indicateur actif V2.

La migration ne réinitialise, ne purge, ne renumérote, ne déplace et ne réécrit
automatiquement aucune ancienne donnée utilisateur : progression JSON,
chapitres YAML et configuration TOML sont intégralement préservés.
Les anciens YAML et TOML ne sont plus chargés automatiquement dans la cible V2.
Aucun import personnalisé ou nettoyage des anciens dossiers n'est prévu.

L'entraînement peut enrichir normalement la progression, en conservant les
champs et entrées historiques, y compris ceux des exercices retirés.
Les acquis affichés pour le corpus actif comptent uniquement les identifiants
encore présents. Ne pas attribuer un ancien score à un nouveau concept par
réutilisation d'un identifiant. Un contenu substantiellement différent reçoit
une nouvelle clé ; les anciens résultats restent stockés.

La sauvegarde doit être atomique et ses erreurs visibles. Un fichier invalide
ou de version non prise en charge ne doit pas être remplacé par un état vide.
Une consultation ou un premier lancement ne déclenche pas de migration écrite.
`reset --yes` demeure la seule suppression explicitement demandée par
l'utilisateur ; une mise à niveau ne doit jamais l'appeler implicitement.

## 8. Terminal, sécurité et fonctionnement local

- Aucun appel réseau nécessaire au fonctionnement, compte, synchronisation,
  télémétrie, IA ou MCP. Les téléchargements de compilation ne sont pas du trafic
  applicatif et ne sont pas requis après installation.
- Aucune lecture des dotfiles pour créer ou adapter le corpus ; aucune détection
  d'installation des outils enseignés.
- Aucune exécution d'alias, de fonction utilisateur ou de commande saisie.
  Aucune modification de dotfiles, de `.zshrc` ou de données héritées par migration.
- Ne jamais lire les chemins sensibles interdits dans les
  [consignes du dépôt](../AGENTS.md#regles-de-securite).
- Conserver les contrôles techniques nécessaires : terminal attaché en entrée
  et sortie, dimensions, disponibilité du mode brut, restauration de l'état du
  terminal, possibilité d'afficher la TUI et sélection du rendu couleur.
- Respecter les indications techniques telles que `TERM` et `NO_COLOR`, sans
  collecter ni afficher de variables sensibles ou inventorier l'environnement.
- Ne pas supposer une capacité ANSI ou de capture à partir du seul système
  macOS. Prévoir un rendu textuel sans couleur et des diagnostics exploitables.
- N'entrer en mode brut que pour la capture qui le nécessite ; restaurer le
  terminal après réponse, interruption ou erreur traitée.
- Si la capture est indisponible, expliquer la limite et permettre de choisir
  un chapitre textuel ; ne pas créditer un raccourci comme capturé.

## 9. Non-objectifs du nouveau MVP

Analyse ou personnalisation d'environnement shell ; scan runtime ; génération
automatique ; catalogue d'outils ; shell réellement exécutant ; détection de
shells ; IA, MCP, compte, synchronisation, télémétrie ; base de données ;
plugins ; import personnalisé ; préférences multiples ; sessions aléatoires ;
reprise persistante au curseur ; capture globale des raccourcis macOS.

Les releases et Homebrew peuvent être traités ultérieurement. Aucun framework
ou nouveau moteur n'est requis pour compenser les suppressions.

## 10. Critères d'acceptation fonctionnels V2

Ces critères devront être vérifiés lors des futures PR ; ils ne décrivent pas
des validations acquises par la présente modification documentaire.

1. Premier entraînement textuel accessible sans configuration, scan ou outil
   tiers installé, avec interface adaptée aux capacités terminal vérifiées.
2. Corpus intégré court et jouable dans chacune des trois familles ; consignes,
   réponses et portabilité vérifiées ; aucune dépendance à des alias personnels.
3. Commandes `chapters`, `train`, `stats`, `reset --yes` disponibles ; commandes,
   option et écrans dotfiles retirés ; anciens fichiers utilisateur intacts.
4. Saisie de commandes et capture réelle compatibles testées ; les limites sont
   expliquées, aucune commande pédagogique n'est exécutée.
5. Choix des QCM visibles, correction respectant les commandes et feedback utile.
6. Progression CLI/TUI enregistrée et relue après fermeture, sans tentative
   fictive, double compte ou perte des résultats historiques.
7. Lecture des données version 1 préservée ; sauvegarde sûre et erreurs visibles ;
   aucune suppression implicite des anciennes données.
8. Navigation, saisie des filtres, sortie et restauration du terminal vérifiées
   sur macOS, en couleur et sans couleur selon les capacités disponibles.
9. Tests adaptés, `go test ./...` et `go vet ./...` réussis ; aucune nouvelle
   dépendance injustifiée ni fonctionnalité de réseau applicatif.

Voir l'[architecture cible](ARCHITECTURE.md) pour les frontières techniques et
le [backlog](../TODO.md) pour les branches atomiques proposées.
