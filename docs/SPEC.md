# CLI-Drill V2 — Spécification fonctionnelle

**Vision du 7 novembre 2026, complétée par P2.0 le 8 novembre 2026 ;
migration fonctionnelle partielle.**
Ce document décrit la cible V2 et ses critères, pas un état de livraison complet.
L'[état du projet](../PROJECT_STATE.md) distingue les travaux intégrés et ceux
qui attendent leur validation ; le [README](../README.md) conserve la présentation
historique du MVP antérieur, avec un avertissement de périmètre.
La [décision V2](DECISIONS.md#2026-11-07--cli-drill-v2--corpus-integre-et-retrait-des-dotfiles)
explique le changement de direction ; le [backlog](../TODO.md) organise sa réalisation.
La [décision P2](DECISIONS.md#2026-11-08--p20--evolution-pedagogique)
remplace les exclusions de l'aléatoire et des préférences pédagogiques minimales.
P2.0 a cadré la feuille de route. P2.4 est publiée à `f096eec` et P2.5 est en
cours. Le choix de langue français/anglais est prévu en [P3](../TODO.md#p3--internationalisation-fren)
et relève d'une phase distincte de P2.

## 1. Vision et public

**CLI-Drill — Learn the terminal by doing.**

CLI-Drill est une application Go CLI/TUI locale pour apprendre et mémoriser
les commandes Unix, les shells, les raccourcis clavier du terminal et les outils CLI.
Elle s'adresse aux personnes qui découvrent le terminal et aux développeurs ou
power users qui souhaitent consolider leurs pratiques.

L'apprentissage repose sur la mémorisation active : exercices courts,
progressifs, répétables, écrits et vérifiés. Saisir une commande et appuyer sur
un vrai raccourci sont des activités centrales. Les questions de connaissance
servent ces pratiques ; elles ne constituent pas à elles seules le produit.

Une session vise cinq à dix minutes par sa longueur éditoriale, sans minuteur
ni préférence de durée. P2 prévoit la sélection manuelle des modules CLI tiers,
trois niveaux et le mélange des séances, avec une exception pour l'apprentissage
guidé ; ces choix ne nécessitent ni scan ni moteur adaptatif.

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
Le nom historique de cette famille ne permet pas d'y enseigner des logiciels
graphiques simplement parce qu'ils sont installés sur macOS. Une commande macOS
peut être étudiée pour sa syntaxe CLI ; reconnaître le nom d'une application
graphique n'est pas un objectif pédagogique. Les éditeurs fonctionnant dans le
terminal peuvent constituer des modules tiers, avec leurs limites de capture.

### Modules CLI et périmètre actif

Les commandes Unix fondamentales restent toujours disponibles : aucun réglage
ne peut désactiver ce socle. Les modules CLI tiers sont activables individuellement
par choix manuel, indépendamment des logiciels réellement installés. Par défaut,
aucun module tiers n'est activé ; le premier entraînement textuel utilise le socle.
Les commandes macOS natives et les notions de shell gardent leur contexte
explicite sans être assimilées à des outils tiers installés.

Un exercice nécessitant plusieurs outils tiers n'est éligible que si **tous**
ses modules sont activés. Par exemple, un pipeline `fd` vers `fzf` exige ces
deux choix. Désactiver un module masque ses exercices pour les nouvelles séances,
sans retirer leurs acquis ou leurs scores. Le socle reste accessible dans le
navigateur ; un filtre de niveau peut néanmoins n'y trouver aucun exercice.
La liste des choix provient uniquement des métadonnées du corpus intégré :
aucun annuaire d'installation, scanner ou ancien catalogue n'est réintroduit.

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
- Rendre chaque exercice autonome : tout contexte utile figure dans sa consigne,
  sans fichier réel ni résultat d'une question précédente nécessaire.
- Préférer commandes saisies, scénarios concrets et captures réellement compatibles
  aux questions répétitives demandant uniquement le nom d'un logiciel.

Ces familles organisent les contenus ; elles ne justifient pas trois moteurs.
P2 commence par un corpus court jouable dans chaque famille, puis vise une
première cible **indicative de 100 à 120 exercices de qualité** dans le corpus
intégré total, répartis progressivement entre les trois familles. Il ne s'agit
ni de 100 exercices supplémentaires, ni d'un quota par séance ou par famille.
La sélection personnelle peut en rendre moins visibles. La qualité et les essais
des lots priment sur le nombre ; aucune duplication de questions pour atteindre
la cible. Les trois niveaux sont introduits progressivement, sans quota artificiel.

### Difficulté de l'exercice

| Niveau affiché | Valeur YAML | Intention éditoriale |
| --- | --- | --- |
| Beginner | `beginner` | Une action courante, contexte explicite, syntaxe élémentaire. |
| Intermediate | `intermediate` | Combiner options, pipes ou étapes conceptuelles dans un scénario autonome. |
| Advanced | `advanced` | Résoudre un cas plus subtil, avec contraintes et limites explicites. |

Le niveau appartient à l'exercice, indépendamment de son format ou de son
chapitre. Choisir un niveau pour une séance ou depuis une fiche chapitre retient
exactement ce niveau, sans inclure implicitement les niveaux inférieurs.
Le défaut est Beginner. Ce choix filtre les exercices ; il ne modifie ni leur
correction ni les scores et ne constitue pas une mesure automatique de maîtrise.
Un ensemble vide produit un diagnostic et permet de changer les filtres, sans
séance fictive, sauvegarde de progression ou activation d'un outil implicite.

## 4. Parcours utilisateur

### Premier lancement et premier entraînement

Après installation, `cli-drill` donne immédiatement accès à l'entraînement,
sans fichier de configuration, scan ni installation d'outils pédagogiques.
`Start training` et `train` sans identifiant choisiront le premier chapitre
textuel ayant des exercices éligibles, dans l'ordre pédagogique du corpus.
Au premier lancement, ce chapitre appartient au socle Unix, au niveau Beginner.
Ce premier chapitre est textuel,
afin de ne pas imposer la capture brute avant vérification de ses capacités.

Les capacités nécessaires à l'interface choisie doivent être vérifiées avant
son activation. Si la TUI est indisponible, fournir un diagnostic clair et une
commande directe vers un chapitre textuel. La CLI textuelle reste utilisable
sans couleur, ligne par ligne, avec une entrée exploitable. En l'absence d'une
telle entrée, signaler la limite sans créer de tentative fictive.

### Navigation TUI cible

| Action | Comportement |
| --- | --- |
| Start training | Démarrer le premier chapitre textuel éligible selon l'ordre pédagogique. |
| Browse chapters | Choisir une famille ou un chapitre, consulter sa fiche et choisir le niveau de la séance. |
| Keyboard shortcuts | Réutiliser l'accès P1.4 à la capture, appliqué aux exercices éligibles ; accès direct si un seul chapitre reste. |
| My progress | Consulter la progression locale avec les composants existants. |
| Quit | Quitter proprement. |

Les cinq entrées n'imposent pas cinq parcours ou composants indépendants.
Conserver la fiche chapitre et le passage au mode clavier dédié, puis le retour
automatique à l'accueil TUI. Ne pas ajouter de cible de retour configurable.
Conserver la palette, la bannière et le style existants.
La sélection des modules et du niveau par défaut utilise les composants existants,
sans refonte de l'accueil ni nouveau parcours de configuration obligatoire.
Une sélection de niveau depuis une fiche vaut pour le lancement concerné ;
aucune table de préférences par chapitre n'est nécessaire. La future CLI doit
permettre les mêmes choix de séance que la TUI ; sa syntaxe sera fixée dans la
sous-phase d'interface, sans rétablir `--config`.

### Répondre, comprendre et poursuivre

Présenter le contexte et la consigne, recevoir la réponse, afficher une
correction avec explication, puis permettre de poursuivre. Une action de
navigation ne devient jamais une réponse. Les touches imprimables doivent
atteindre la réponse ou le filtre actif avant les raccourcis globaux.

Le bilan distingue le résultat de la session et les compteurs historiques.
La révision des erreurs reste en mémoire pour la session ; aucun nouveau
stockage de session n'est requis. Une relance de chapitre crée une nouvelle
séance, tout en conservant les acquis et tentatives précédents. La reprise
exacte au curseur interrompu est différée.

### Ordre de séance et apprentissage guidé

Au lancement, former un instantané des exercices du chapitre : appliquer la
sélection des modules, puis le niveau, puis déterminer l'ordre. Une séance
ordinaire mélange cette liste une seule fois. La permutation reste en mémoire
et stable jusqu'à la sortie ; feedback, aide, redimensionnement, retour d'une
capture et retry ne la recalculent pas. Chaque exercice éligible apparaît une
fois dans le parcours initial, aucun autre n'y apparaît. Le corpus source n'est
pas réordonné et la clé enregistrée est toujours l'identifiant original,
jamais la position dans la permutation.

Chaque nouvelle séance de révision refait le mélange ; deux tirages peuvent
produire le même ordre, notamment avec zéro ou un exercice. Un retry volontaire
reste sur le même exercice ; la révision des erreurs est une passe distincte
sur la liste des erreurs. Ces nouvelles réponses ajoutent des tentatives selon
P1.2, sans dupliquer les acquis ni le total du parcours initial.

Lorsqu'un ordre pédagogique est indispensable, le chapitre le déclare :
**Apprendre** suit l'ordre éditorial filtré, **Réviser** mélange les mêmes
exercices lors des séances ultérieures. Le parcours guidé est proposé par défaut
pour ces chapitres ; l'utilisateur choisit explicitement la révision. Aucun
marqueur persistant de première séance, ni déduction de maîtrise depuis les
scores cumulés, n'est requis. Le guidage ordonne des exercices autonomes ; il
ne permet pas de faire dépendre une réponse du résultat précédent.
Les préférences modifiées s'appliquent au lancement suivant. Quitter une séance
abandonne son ordre en mémoire, sans perdre les réponses déjà sauvegardées.

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
P2 prévoit seulement les métadonnées additives suivantes ; aucun YAML n'est
modifié par P2.0 :

| Portée / champ prévu | Contrat minimal P2 |
| --- | --- |
| Exercice : `difficulty` | `beginner`, `intermediate` ou `advanced`. Absence historique : `beginner` en mémoire, sans réécriture. Valeur explicite vide ou inconnue : erreur contextualisée. |
| Exercice : `tools` | Liste d'identifiants stables de modules CLI tiers, par exemple `[fd, fzf]`. Absence ou liste vide : aucun module requis. Tous les modules listés doivent être activés. |
| Chapitre : `family` | `unix-shell`, `keyboard-shortcuts` ou `modern-cli-macos` pour le corpus révisé ; l'absence reste tolérée par le chargeur historique, sans famille inférée. Valeur explicite vide ou inconnue : erreur. |
| Chapitre : `learning_order` | `independent` ou `guided`. Absence historique : `guided` par prudence ; chaque chapitre révisé déclare sa valeur. Valeur explicite vide ou inconnue : erreur. |

Les identifiants de modules désignent uniquement les outils enseignés dans le
corpus, sans chemins, binaire détecté, version installée ou source dotfiles.
Chaque liste contient des identifiants non vides et sans doublons ; elle ne
constitue pas un inventaire des logiciels disponibles sur la machine.
Les contextes de portabilité et préconditions restent dans `description`,
`prompt` et `explanation` ; aucune matrice de capacités ou grammaire de scénario
n'est nécessaire. Les nouveaux exercices ont un niveau explicite ; tous les
exercices tiers révisés déclarent leurs outils avant l'activation des filtres.
Un ancien YAML sans difficulté reste lisible avec son sens et ses identifiants,
mais n'est pas automatiquement importé. Les défauts historiques permettent le
chargement ; ils ne certifient pas sa qualité éditoriale ni l'absence d'outil tiers.

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
encore présents et éligibles selon les modules et le niveau affichés ; préciser
ce périmètre dans l'interface. Les scores de chapitre v1 restent cumulés sur
toutes les tentatives, sans prétendre les ventiler par outil ou difficulté.
Le score de séance reste calculé en mémoire sur ses seules réponses.
Ne pas attribuer un ancien score à un nouveau concept par
réutilisation d'un identifiant. Un contenu substantiellement différent reçoit
une nouvelle clé ; les anciens résultats restent stockés.
Un regroupement par famille ou un changement de titre conserve le couple
`chapter-id/item-id` si l'apprentissage est inchangé ; changer seulement la
position d'un exercice ne doit jamais changer sa clé.

La sauvegarde doit être atomique et ses erreurs visibles. Un fichier invalide
ou de version non prise en charge ne doit pas être remplacé par un état vide.
Une consultation ou un premier lancement ne déclenche pas de migration écrite.
`reset --yes` demeure la seule suppression explicitement demandée par
l'utilisateur ; une mise à niveau ne doit jamais l'appeler implicitement.

### Préférences pédagogiques P2, distinctes de la progression

Contrat minimal prévu, **non implémenté par P2.0** : un fichier local
`preferences.json` séparé, sous `$XDG_DATA_HOME/cli-drill/` ou, par défaut,
`~/.local/share/cli-drill/`. Le choix de ce répertoire réutilise la résolution
de données existante ; il ne rétablit ni le TOML dotfiles ni `--config`.
`--progress` ne déplace que la progression. Aucun fichier n'est créé dans cette
mission et le futur premier lancement n'en crée pas automatiquement.

| Champ prévu | Rôle et valeur par défaut en mémoire |
| --- | --- |
| `version` | Version propre aux préférences, initialement `1` ; indépendante du JSON de progression v1. |
| `enabled_tools` | Liste d'identifiants de modules CLI activés manuellement ; défaut `[]`. |
| `default_difficulty` | `beginner`, `intermediate` ou `advanced` ; défaut `beginner`. |

Les champs facultatifs omis utilisent ces défauts ; `version`, lorsqu'un fichier
existe, doit être présente et prise en charge. Un mauvais type ou une valeur
explicite hors contrat est invalide, sans conversion silencieuse.

Ces choix d'encodage et les défauts précisent la mise en œuvre minimale de P2.0.
Ils n'ajoutent ni préférences de durée/navigation, ni profil de shell, ni réglage
persistant par chapitre. La sélection ponctuelle de séance prime sur le défaut
et reste en mémoire. Seule une modification explicitement enregistrée par
l'utilisateur écrit les préférences, avec remplacement atomique et erreur visible.
Une consultation ou un entraînement n'écrit jamais ce fichier.

Absence du fichier : utiliser les défauts sans écriture. Fichier illisible,
invalide ou de version inconnue : diagnostic, aucune réparation ni réécriture
automatique ; permettre le socle avec les défauts en mémoire, signalés comme tels,
et refuser la sauvegarde tant que le fichier ne peut être lu et validé.
Un identifiant de module devenu absent du corpus est ignoré pour la sélection,
signalé et conservé lors des sauvegardes, sans recherche d'installation.
Un échec d'enregistrement ne doit pas annoncer les nouveaux choix comme persistés.

Séparer les deux stockages : changer un outil ou un niveau ne lit ni ne réécrit
les anciens YAML/TOML et ne modifie pas le JSON de progression. Inversement,
`reset --yes` continue de viser uniquement la progression, jamais les préférences.

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
plugins ; import personnalisé ; préférences au-delà des modules CLI et du niveau
par défaut ; reprise persistante au curseur ; capture globale des raccourcis macOS ;
moteur adaptatif et répétition espacée. Enseigner des logiciels graphiques pour
leur seule présence sur macOS reste hors périmètre.

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
10. Sélection manuelle des modules tiers utilisable sans installation ; socle
    Unix toujours accessible, filtres combinés prévisibles et diagnostic si vide.
11. Trois niveaux portés par les exercices ; anciens YAML sans difficulté lus
    comme Beginner sans modification ; ni import automatique ni changement de clés.
12. Chaque séance ordinaire ou révision utilise une permutation complète des seuls
    exercices éligibles, stable jusqu'à la sortie. Apprentissage guidé conservé
    lorsque nécessaire ; retries et feedback respectent les compteurs P1.
13. Préférences séparées, facultatives, créées seulement sur enregistrement
    explicite ; erreurs visibles, aucune modification du JSON v1 par ces réglages.
14. Enrichissement par lots essayés vers 100 à 120 exercices au total, sans
    sacrifier l'autonomie, la qualité, les trois familles ou les limites de capture.

Voir l'[architecture cible](ARCHITECTURE.md) pour les frontières techniques et
le [backlog](../TODO.md) pour les branches atomiques proposées.
