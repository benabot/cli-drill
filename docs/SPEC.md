# CLI-Drill V2 — Spécification fonctionnelle

**Vision du 7 novembre 2026, complétée par P2.0 le 8 novembre et P3.0
les 12 et 13 novembre 2026.**
Ce document décrit la cible V2 et ses critères.
L'[état du projet](../PROJECT_STATE.md) distingue les travaux intégrés et ceux
qui attendent leur validation. Le [README](../README.md) présente le logiciel
disponible sur `dev`.
La [décision V2](DECISIONS.md#2026-11-07--cli-drill-v2--corpus-integre-et-retrait-des-dotfiles)
explique le changement de direction ; le [backlog](../TODO.md) organise sa réalisation.
La [décision P2](DECISIONS.md#2026-11-08--p20--evolution-pedagogique)
remplace les exclusions de l'aléatoire et des préférences pédagogiques minimales.
La base `dev` à `91efab0` intègre le corpus de P2.5 et la clôture de P2.6.
La [politique linguistique P3](#politique-linguistique-p3) est décidée.
La localisation reste à réaliser.

## 1. Vision et public

**CLI-Drill — Learn the terminal by doing.**

CLI-Drill est une application Go CLI/TUI locale qui enseigne la pratique du
terminal, des commandes indispensables jusqu'à leur composition et au
diagnostic de situations concrètes. Elle s'adresse aux débutants, aux
développeurs, aux administrateurs et aux utilisateurs avancés. Les systèmes
cibles du produit sont macOS, Linux et les systèmes Unix/BSD. Les plateformes
annoncées pour la première distribution publique seront celles validées
nativement en [P5](../TODO.md#p5--distribution-multiplateforme).

L'apprentissage repose sur la mémorisation active : exercices courts,
progressifs, répétables, écrits et vérifiés. Saisir une commande et appuyer sur
un vrai raccourci sont des activités centrales. Les questions de connaissance
servent ces pratiques ; elles ne constituent pas à elles seules le produit.

Un lot propose jusqu'à six exercices après sélection du niveau et des modules.
Les cinq à dix minutes restent un repère à vérifier en recette humaine.
Le parcours intégral est disponible sur demande.

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

Le socle couvre les commandes Unix/POSIX et leurs options utiles pour la
navigation, les fichiers, les permissions, les liens, les processus,
l'environnement et `PATH`. Il aborde aussi les pipes, les redirections,
les archives et Git essentiel. La recherche et le traitement de texte avec
`grep`, `find`, `sed` et `awk` partent de tâches courantes au terminal.
Les exercices de commande fournissent un contexte simulé suffisant : objectif,
noms de fichiers et résultat recherché. Ils ne dépendent pas du disque réel.
Une construction shell courte peut servir un usage interactif précis.
[Priorités du socle Unix](CURRICULUM.md#approfondissements-du-socle-unix)

### B. Keyboard Shortcuts

L'édition interactive couvre les déplacements, la correction, l'historique,
la complétion et les incidents qui interrompent la saisie. Bash/Readline sert
de référence. Les différences Zsh/ZLE viennent après les gestes communs et
sont enseignées quand elles modifient le résultat. Les modes Emacs et Vi
sont précisés quand leurs liaisons expliquent le comportement observé.
Les processus, jobs et signaux relèvent des fondamentaux Unix.

Chaque exercice précise le shell, le mode d'édition et les réglages qui
déterminent sa réponse. Les essais distinguent le comportement de l'éditeur,
les octets transmis par le terminal et ceux que CLI-Drill reconnaît. Les
séquences compatibles sont évaluées par capture réelle. Les séquences
interceptées, indisponibles ou réservées à la sortie peuvent faire l'objet
d'un diagnostic textuel, évalué comme tel. La validation choisit des
terminaux représentatifs des systèmes cibles.

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
par choix manuel, indépendamment des logiciels réellement installés. Sans
fichier de préférences, Git essentiel est actif. Les autres modules restent
facultatifs et le premier entraînement textuel utilise le socle.
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

- Indiquer le contexte pertinent, qu'il relève de POSIX, de GNU, de BSD/macOS,
  du shell, de son mode d'édition, du terminal ou d'un outil tiers. Préciser
  une version minimale ou un réglage lorsque la réponse en dépend.
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

P2.5 enrichit le corpus par lots. **30 exercices par chapitre** est un
minimum indicatif pour les périmètres généraux. Le corpus peut dépasser ce
repère selon les compétences à couvrir. Un périmètre spécialisé peut rester
en dessous lorsque ses usages utiles sont couverts et que la justification
figure dans le référentiel. Les niveaux suivent la progression des compétences.
Advanced privilégie les diagnostics, la composition, la portabilité et la
robustesse des commandes dans des situations concrètes. Écarter les détails
ésotériques et les répétitions qui servent uniquement à atteindre un nombre.

**360 exercices** reste le repère de planification historique des douze
chapitres. Les objectifs sont réévalués selon les compétences utiles à couvrir.
La création d'un chapitre requiert une validation explicite de son périmètre.
[Référentiel pédagogique](CURRICULUM.md)

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
aucune table de préférences par chapitre n'est nécessaire. Dans la fiche,
`d` change le niveau, `m` bascule entre Apprendre et Réviser et `a` choisit
le format court ou intégral. Les équivalents CLI sont `--difficulty`,
`--mode learn|review` et `--all`.

### Répondre, comprendre et poursuivre

Présenter le contexte et la consigne, recevoir la réponse, afficher une
correction avec explication, puis permettre de poursuivre. Une action de
navigation ne devient jamais une réponse. Les touches imprimables doivent
atteindre la réponse ou le filtre actif avant les raccourcis globaux.

Chaque bilan affiche les exercices réalisés, les bonnes réponses, les erreurs
et les acquis cumulés de l'historique. `c` continue s'il reste des exercices
dans la visite, `r` reprend les seules erreurs et `q` termine. Entrée continue
si possible, sinon ouvre la reprise quand il reste des erreurs. Dans la CLI
textuelle, valider la lettre par Entrée. Le menu en mode capture lit la touche
directement. Dans la TUI, Échap revient à la fiche et `g` à l'accueil.

La reprise affiche son propre score et conserve celui du lot initial. Une
nouvelle reprise porte sur les erreurs restantes de la dernière passe. Le
réessai immédiat déjà disponible en capture ajoute une tentative et son
résultat apparaît séparément au bilan. Toute bonne réponse peut ajouter un
acquis unique. L'interruption conserve les réponses sauvegardées.

### Ordre de séance et apprentissage guidé

Au lancement, appliquer le niveau, les modules et le périmètre de l'entrée
choisie. **Apprendre**, le mode par défaut, propose les exercices non acquis.
Il conserve l'ordre éditorial pour les chapitres guidés et mélange les chapitres
indépendants. Une fois le filtre acquis, le logiciel propose **Réviser**.
**Réviser** mélange tous les exercices éligibles, y compris ceux déjà réussis.
Le parcours intégral explicite (`--all` ou `a`) propose tous les exercices
éligibles en un seul lot, selon l'ordre du mode choisi.

La liste est fixée une fois par visite. **Continuer** prend la suite des
exercices sans rejouer ceux déjà présentés. Une erreur n'empêche pas d'avancer.
La reprise des erreurs utilise une liste distincte en mémoire. Le retour
après interruption recalcule un lot d'apprentissage à partir des acquis
JSON v1. Un exercice raté peut donc revenir à la prochaine visite.
Le curseur, l'ordre du tirage et les erreurs de la visite ne sont pas
persistés. Le JSON v1 ne permet pas de classer les erreurs anciennes par
exercice. La révision ne les privilégie pas.

Les modifications de préférences s'appliquent au lancement suivant.
Chaque réponse conserve sa clé d'origine `chapter-id/item-id`.

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
Utiliser un seul processus écrivain par fichier. Deux entraînements
simultanés peuvent perdre des mises à jour malgré le remplacement atomique.
Une consultation ou un premier lancement ne déclenche pas de migration écrite.
`reset --yes` demeure la seule suppression explicitement demandée par
l'utilisateur ; une mise à niveau ne doit jamais l'appeler implicitement.

### Préférences pédagogiques P2, distinctes de la progression

Le fichier local
`preferences.json` séparé, sous `$XDG_DATA_HOME/cli-drill/` ou, par défaut,
`~/.local/share/cli-drill/`. Le choix de ce répertoire réutilise la résolution
de données existante ; il ne rétablit ni le TOML dotfiles ni `--config`.
`--progress` ne déplace que la progression. Le premier lancement
n'en crée pas automatiquement.

| Champ prévu | Rôle et valeur par défaut en mémoire |
| --- | --- |
| `version` | Version propre aux préférences, initialement `1` ; indépendante du JSON de progression v1. |
| `enabled_tools` | Liste des modules actifs. Sans fichier, défaut `["git"]`. Dans un fichier existant, une liste vide ou omise reste vide. |
| `default_difficulty` | `beginner`, `intermediate` ou `advanced` ; défaut `beginner`. |

Le niveau omis prend la valeur `beginner`. Dans un fichier existant,
`version` est obligatoire. Un mauvais type ou une valeur explicite hors
contrat produit un diagnostic.

Ces choix d'encodage et les défauts précisent la mise en œuvre minimale de P2.0.
Ils n'ajoutent ni préférences de durée/navigation, ni profil de shell, ni réglage
persistant par chapitre. La sélection ponctuelle de séance prime sur le défaut
et reste en mémoire. Seule une modification explicitement enregistrée par
l'utilisateur écrit les préférences, avec remplacement atomique et erreur visible.
Une consultation ou un entraînement n'écrit jamais ce fichier.

Absence du fichier : utiliser les défauts sans écriture. Fichier illisible,
invalide ou de version inconnue : diagnostic et accès au socle en niveau
Beginner, sans module actif. La sauvegarde reste bloquée tant que le fichier
ne peut être lu et validé.
Un identifiant de module devenu absent du corpus est ignoré pour la sélection,
signalé et conservé lors des sauvegardes, sans recherche d'installation.
Un échec d'enregistrement ne doit pas annoncer les nouveaux choix comme persistés.

Séparer les deux stockages : changer un outil ou un niveau ne lit ni ne réécrit
les anciens YAML/TOML et ne modifie pas le JSON de progression. Inversement,
`reset --yes` continue de viser uniquement la progression, jamais les préférences.

<a id="politique-linguistique-p3"></a>

### Politique linguistique P3

Sur la base P2.6, le corpus et une partie des textes d'interface sont en
français. D'autres libellés sont en anglais. Les préférences v1 couvrent
les modules et la difficulté. Le choix de langue reste à implémenter.

La décision du 12 novembre 2026 fixe l'anglais (`en`) comme défaut universel,
y compris sur une installation existante sans langue enregistrée.
Le français (`fr`) se choisit explicitement dans les préférences CLI et TUI,
comme l'anglais. Le système, le shell et les variables régionales
n'influencent pas ce choix. Le premier entraînement reste accessible sans
sélection obligatoire de langue.

| Préférences | Comportement attendu |
| --- | --- |
| Fichier absent | Anglais en mémoire, sans création de fichier. |
| Fichier valide sans `language` | Anglais en mémoire, sans modification du fichier. |
| `language: "en"` | Anglais. |
| `language: "fr"` | Français. |
| Valeur explicite invalide | Diagnostic visible, sans réparation ni écrasement automatique. |

La consultation des préférences reste en lecture seule. Un choix explicite
enregistré est conservé après relance, avec les modules activés et la
difficulté par défaut. La progression conserve son schéma JSON v1 et ses
clés `chapter-id/item-id`, indépendamment de la langue. Un acquis obtenu en
français est reconnu en anglais et réciproquement. Le changement de langue
ne modifie ni les acquis, ni les scores, ni les tentatives, ni l'historique.
Aucune migration automatique de la progression.

La version publique couvre dans les deux langues les menus, aides,
indications de navigation, réglages, messages, feedback et diagnostics
de la CLI et de la TUI, y compris pendant la capture. Elle couvre les
titres et descriptions des chapitres, les consignes des 446 exercices,
les choix QCM, les explications, les corrections et la documentation
utilisateur. Les contenus localisés se rattachent au même corpus et aux
mêmes identifiants. Les commandes Unix, options, syntaxes techniques et
séquences de touches gardent leur sens et leurs caractères significatifs.
Les réponses conceptuelles et choix traduits doivent conserver les mêmes
bonnes réponses et la même difficulté, sans ambiguïté supplémentaire ni
indice involontaire.

La [décision du 13 novembre](DECISIONS.md#p30-compatibilite-des-preferences-en-pre-release)
fixe le contrat de mise à niveau en pré-release. Le nouveau binaire doit
lire les fichiers de préférences valides existants et conserver les modules
activés et la difficulté par défaut. Le lancement et la consultation
n'entraînent ni migration ni réécriture. L'enregistrement explicite peut
mettre à jour le fichier selon le schéma retenu, avec les protections
décrites ci-dessus. Un fichier invalide reste intact et l'erreur est visible.
Aucune garantie de rétrogradation vers un ancien binaire n'est exigée.

P3.1 devra concevoir les ressources de traduction, leurs interfaces d'accès
et la syntaxe CLI. Le champ logique `language` sera stocké dans
`preferences.json`. Le choix entre étendre la version 1 de ce schéma et
passer à une version 2 techniquement justifiée reste ouvert. L'étude porte
sur la lecture et l'enregistrement lors d'une mise à niveau, sur fixtures
synthétiques.
Le [backlog P3](../TODO.md#p3--internationalisation-fren) précise les lots
et la recette bilingue.

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

- Analyse ou personnalisation de l'environnement shell, scan runtime,
  génération automatique, catalogue d'outils, exécution et détection de shells.
- IA, MCP, compte, synchronisation, télémétrie, base de données et plugins.
- Import personnalisé, reprise persistante au curseur, capture globale des
  raccourcis macOS, moteur adaptatif et répétition espacée.
- Préférences au-delà des modules CLI, du niveau par défaut et du choix
  de langue P3.
- Leçons sur des logiciels graphiques pour leur seule présence sur macOS.

La programmation shell, l'écriture de scripts complets, les fonctions shell
avancées et le débogage de scripts sont hors périmètre pédagogique.

La distribution publique relève de [P5](../TODO.md#p5--distribution-multiplateforme).
Aucun framework ou nouveau moteur n'est requis pour compenser les suppressions.

## 10. Critères d'acceptation fonctionnels V2

Ces critères s'appliquent aux futures validations. Le rapport de recette P2.6
consigne les preuves déjà obtenues.

1. Premier entraînement textuel accessible sans configuration, scan ou outil
   tiers installé, avec interface adaptée aux capacités terminal vérifiées.
2. Corpus intégré jouable dans les trois familles. Consignes, réponses et
   portabilité vérifiées dans les contextes annoncés, sans dépendance aux
   alias personnels.
3. Commandes `chapters`, `train`, `stats`, `reset --yes` disponibles ; commandes,
   option et écrans dotfiles retirés ; anciens fichiers utilisateur intacts.
4. Saisie de commandes et capture réelle compatibles testées ; les limites sont
   expliquées, aucune commande pédagogique n'est exécutée.
5. Choix des QCM visibles, correction respectant les commandes et feedback utile.
6. Progression CLI/TUI enregistrée et relue après fermeture, sans tentative
   fictive, double compte ou perte des résultats historiques.
7. Lecture des données version 1 préservée ; sauvegarde sûre et erreurs visibles ;
   aucune suppression implicite des anciennes données.
8. Navigation, filtres, sortie et restauration du terminal vérifiés sur une
   matrice représentative des systèmes cibles, en couleur et sans couleur
   selon les capacités disponibles. Consigner les systèmes, shells, modes
   d'édition et terminaux effectivement essayés.
9. Tests adaptés, `go test ./...` et `go vet ./...` réussis ; aucune nouvelle
   dépendance injustifiée ni fonctionnalité de réseau applicatif.
10. Sélection manuelle des modules tiers utilisable sans installation ; socle
    Unix toujours accessible, filtres combinés prévisibles et diagnostic si vide.
11. Trois niveaux portés par les exercices ; anciens YAML sans difficulté lus
    comme Beginner sans modification ; ni import automatique ni changement de clés.
12. Les lots courts contiennent au plus six exercices du filtre. Apprendre
    privilégie les non-acquis dans l'ordre prévu, Réviser inclut tous les
    éligibles, et Continuer ne rejoue aucun exercice présenté pendant la
    visite. Le parcours intégral et la reprise ciblée respectent les compteurs P1.
13. Préférences séparées, facultatives, créées seulement sur enregistrement
    explicite ; erreurs visibles, aucune modification du JSON v1 par ces réglages.
14. Couverture des compétences et progression vérifiées par lot. Justifier la
    profondeur du socle et les périmètres spécialisés sous le minimum
    indicatif. Contrôler l'autonomie des exercices et les formats de
    réponse, y compris la distinction entre diagnostic textuel et capture.
15. Pour la version publique bilingue, premier lancement en anglais, choix
    EN/FR persistant et vérifié dans les deux interfaces, préférences anciennes
    intactes en lecture et progression préservée à chaque changement.
    Des contrôles bloquent les traductions manquantes ou incohérentes et le
    mélange accidentel des langues. Une recette humaine dans les deux langues
    vérifie les parcours et l'équivalence pédagogique.

Le [rapport de recette P2.6](P26_FINAL_ACCEPTANCE.md) consigne les preuves
par environnement et la checklist physique. La durée cible de 5 à 10 minutes
reste à vérifier par chronométrage en usage réel.

Voir l'[architecture cible](ARCHITECTURE.md) pour les frontières techniques et
le [backlog](../TODO.md) pour les branches atomiques proposées.
