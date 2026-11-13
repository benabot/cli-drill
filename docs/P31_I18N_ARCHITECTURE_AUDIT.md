# P3.1 — Audit de l’architecture d’internationalisation EN/FR

Date : 13 novembre 2026. Statut : **audit terminé ; recommandations à valider, aucune implémentation**.

La [section 14](#p31-comparaison-experimentale) actualise les recommandations de
l’audit initial à partir des décisions du mainteneur et des deux prototypes.

Base examinée : `dev`, commit `eeab86818d5e91260479fc8fefc2ea4c96648753`
(`docs: define P3-P6 launch roadmap and language policy`), après intégration de P3.0.
Son parent est la clôture P2.6, `91efab04e19f1d1bebaef715d248e8cba4412312`.
L’arbre était propre avant création de `docs/p31-i18n-architecture`.

## 1. Résumé exécutif

**Faits vérifiés.** Le produit n’a ni contexte linguistique ni ressources bilingues.
Les textes anglais et français sont répartis dans les interfaces, les erreurs et
les YAML. Le lecteur de préférences v1 refuse les champs inconnus, donc refuse
actuellement `language`, même avec `en` ou `fr`. Les QCM corrigent le **libellé**
sélectionné. Certaines transitions TUI dépendent également de textes affichés.
Une traduction des seules consignes ne suffirait pas.

**Recommandation d’ensemble, non décision acquise :**

1. Étendre les préférences v1 par un champ facultatif `language`, sans migration
   sur lecture ; conserver les protections actuelles d’écriture.
2. Injecter un petit catalogue Go typé et immuable dans les interfaces existantes,
   construit à partir des préférences avant les aides et diagnostics Cobra.
3. Conserver un corpus canonique de métadonnées et de réponses techniques ; placer
   les textes EN/FR dans des ressources symétriques par chapitre, embarquées.
4. Donner aux choix QCM une identité stable ; distinguer réponses techniques
   partagées et réponses textuelles localisées, sans changer les identités
   d’exercices ni le schéma de progression.
5. Valider structure, complétude et correction dans les deux langues, puis faire
   relire les formulations et réaliser la recette macOS/Ghostty.

Ces choix minimisent les changements du moteur et de la persistance. Ils demandent
néanmoins une adaptation explicite du modèle QCM et du chargement du corpus.
Ni le format proposé ni la syntaxe CLI illustrative ci-dessous ne sont disponibles.

## 2. Références et état réel avant P3.1

Les contrats applicables viennent de [DECISIONS](DECISIONS.md), décisions des
12 et 13 novembre 2026, de [SPEC](SPEC.md), d’[ARCHITECTURE](ARCHITECTURE.md),
de [PROJECT_STATE](../PROJECT_STATE.md), du [TODO](../TODO.md) et du
[CURRICULUM](CURRICULUM.md). [AGENTS](../AGENTS.md), le
[skill cli-drill-spec](../.codex/skills/cli-drill-spec/SKILL.md) et le
[README actuel](../README.md) ont été consultés. Les sources du checkout priment
sur l’ancien résultat TODO renvoyé par la recherche QMD lors de cet audit.

La [roadmap de lancement](CLI-DRILL_ROADMAP_LANCEMENT_GITHUB.md) reste la source
du cadrage public. Ses sections 2–6 fondent le positionnement, le README EN/FR,
les démonstrations CleanShot X et le Wiki ; les sections 8–12 fondent les
contrôles de distribution, les dépendances et les critères de publication.
Les décisions P3.0 remplacent ses anciennes questions sur la langue de l’interface
et du Wiki. Elles conservent P4 avant la traduction des appellations obsolètes,
la préparation parallèle de P5/P6 et la convergence avant lancement public.
Le présent rapport ne choisit ni licence, ni stratégie bêta/sortie unique et
n’engage aucune diffusion. La documentation et les captures devront correspondre
à la version réellement validée, notamment pour le futur changement de langue.

### Contrats déjà décidés, non soumis à un nouvel arbitrage

- `en` par défaut universel, y compris pour un ancien fichier valide sans langue ;
  `fr` sélectionnable explicitement en CLI et TUI, puis persistant.
- Aucune détection de locale système, aucun assistant obligatoire, aucun fichier
  supplémentaire uniquement pour la langue.
- Absence de fichier ou de champ : anglais en mémoire, sans écriture. Valeur
  explicite invalide : diagnostic visible, aucune réparation automatique.
- Pré-release sans utilisateurs externes : aucune garantie de rétrogradation.
  Une v2 des préférences est autorisée si justifiée, mais n’est pas imposée.
- Préserver modules et difficulté valides ; ne pas réécrire au lancement ou à la
  consultation ; refuser l’écrasement silencieux d’un fichier invalide.
- Progression JSON v1, clés `chapter-id/item-id`, acquis, scores, tentatives et
  historiques strictement conservés ; aucune migration de progression.
- Traduction complète des contenus actifs des 446 exercices avant release
  publique bilingue, avec commandes, options et raccourcis techniques inchangés.

### Mesures sur le corpus chargé par le code actuel

| Type d’exercice | Nombre vérifié |
| --- | ---: |
| `simple-shell-sim` | 270 |
| `multiple-choice` | 99 |
| `free-answer` | 59 |
| `key-sequence` | 18 |
| `scenario` | 0 |
| **Total, 13 chapitres** | **446** |

Le type `scenario` est accepté par le chargeur, sans occurrence actuelle.
P2.6 a sa recette macOS/Ghostty documentée ; cet audit n’en refait pas la recette.
Les compilations croisées historiques ne constituent pas une validation native
Linux/BSD, et les tests de cet audit ne valident pas une future interface bilingue.

## 3. Inventaire des points d’entrée et des chaînes

Cet inventaire est organisé par surface et fonction, pas par nombre de littéraux :
un littéral technique, une identité et un message utilisateur n’ont pas le même rôle.
Les chaînes des dépendances effectivement affichées font partie du périmètre.

| Surface | Sources et points d’entrée examinés | Textes à localiser / éléments à préserver |
| --- | --- | --- |
| Démarrage, sortie d’erreur | [main.go](../cmd/cli-drill/main.go), `main` ; [root.go](../internal/app/root.go), `NewRootCommand`, `commandOptions` | Diagnostic sans TTY, erreurs d’exécution, préfixes ; conserver codes de sortie et causes techniques. |
| Commandes, aides, arguments | `NewRootCommand`, `newChaptersCommand`, `newTrainCommand`, `newStatsCommand`, `newResetCommand` dans `root.go` | `Short`, `Long`, descriptions d’options, usage, arguments invalides, commande inconnue, aide/completion Cobra ; noms `train`, `prefs`, options et valeurs techniques stables. |
| Liste des chapitres | [chapters_output.go](../internal/app/chapters_output.go), `writeChapters`, `writeChapterRows`, `writeCompactChapters` | En-têtes, unités, titres ; conserver IDs, structure des sorties compactes et règles TTY/redirection. Éviter les pluriels par simple ajout de `s`. |
| Préférences CLI | [preferences_command.go](../internal/app/preferences_command.go), `newPreferencesCommand` | Réglages, modules activés/désactivés, avertissement de module retiré, aide, sauvegarde, erreurs ; conserver identifiants des modules et valeurs `beginner`/`intermediate`/`advanced` en arguments. |
| Stockage des préférences | [preferences.go](../internal/app/preferences.go), `loadPreferences`, `readPreferences`, `validate`, `savePreferences`, `writePreferences` | Schéma, champ inconnu, fichier invalide/illisible, conflit d’écriture. Ne pas confondre diagnostic et permission de remplacer le fichier. |
| Accueil et navigation TUI | [tui.go](../internal/tui/tui.go), `homeActions`, `newModel`, `Update`, `View` | Actions, descriptions, titres, aide clavier, pied de page, retour/quitter, filtre de liste. Les actions doivent avoir un ID indépendant de leur titre. |
| Fiches et parcours | `tui.go` ; [session.go](../internal/tui/session.go), `sessionLevelLabel` | Titres/descriptions du corpus, niveaux affichés, familles, progression, sélection des modes, raisons d’absence d’exercices. |
| Réglages TUI | [preferences.go](../internal/tui/preferences.go), `newPreferencesScreen`, méthodes de `preferencesScreen` ; `tuiOptions` dans `app/root.go` | Brouillon, modules, niveau, enregistrer/annuler, erreurs et confirmations ; étendre le modèle de préférences transmis aux callbacks. |
| Entraînement CLI | [training_text.go](../internal/app/training_text.go), `textTrainingView` ; [training_session.go](../internal/app/training_session.go), `readBatchAction` | Question, saisie, feedback, correction, bilan, reprise/réessai/interruption. Garder les touches `c/r/q` et traduire leur description, sauf décision ultérieure distincte. |
| Entraînement TUI | `tui.go`, `answerOrAdvance`, rendu des exercices et bilans ; [readable.go](../internal/tui/readable.go) | Choix, réponse, solution, explication, aide au défilement, erreurs de copie ; état correct/incorrect indépendant du texte affiché. |
| Capture réelle | `app/root.go`, `runKeySequenceExercise` et rendus associés ; [capture_input.go](../internal/app/capture_input.go), `capturedKeyLabel` | Consigne, touche reçue, essai, résultat, solution. Séparer libellés explicatifs et notations reconnues. |
| Terminal indisponible, signaux | [capture_terminal_unix.go](../internal/app/capture_terminal_unix.go), `captureSignalError` ; [capture_terminal_other.go](../internal/app/capture_terminal_other.go) | Activation/restauration, interruption, capture indisponible ; préserver octets, protocoles, signaux et codes de sortie. |
| Sélection pédagogique | [chapter/session.go](../internal/chapter/session.go), `SessionFilter`, raisons de sélection vide ; [visit.go](../internal/chapter/visit.go) | Raisons et descriptions, actuellement aussi produites en anglais dans le domaine ; préférer raisons structurées rendues dans la couche utilisateur. |
| Corpus et validation | [chapter.go](../internal/chapter/chapter.go), `Load`, `LoadFS`, `Chapter.Validate`, `Item.MatchAnswer` ; [metadata.go](../internal/chapter/metadata.go) | Textes pédagogiques et diagnostics de contenu ; conserver `ID`, type, difficulté, outils, ordre et règles de correction. |
| Progression et statistiques | `newStatsCommand`, écrans TUI ; [progress.go](../internal/progress/progress.go), `Load`, `RecordAnswer`, sauvegarde | Libellés, compteurs, erreurs JSON/IO ; aucune langue ni nouveau champ dans la progression. |
| Chemins et embarquement | [xdg.go](../internal/xdg/xdg.go) ; [defaults.go](../data/defaults.go), `Chapters` | Chemins, noms de fichiers et identifiants restent techniques. Adapter plus tard l’embarquement des ressources sans accès réseau. |

### Dépendances de comportement aux textes : faits vérifiés

- Dans `tui.go`, l’action de l’accueil est dispatchée à partir de `selected.title`.
  Traduire le titre sans ajouter une identité d’action casserait la navigation.
- Le rendu de correction compare le feedback au texte `Réponse incorrecte`.
  Remplacer ce couplage par un état de résultat, puis traduire son affichage.
- La TUI transmet `item.Choices[index]` à `Item.MatchAnswer`. Cette chaîne est à la
  fois un libellé de choix et une réponse pédagogique.
- Les catalogues de modules, les filtres de niveau et les touches de capture
  utilisent des valeurs techniques : ne pas les traduire pour localiser un écran.

### Textes provenant des bibliothèques

Versions examinées dans le module/cache local : Cobra `v1.9.1`, Bubbles `v0.21.0`,
Bubble Tea `v1.3.4`. Le catalogue applicatif ne couvrira pas automatiquement ces textes.

L’accueil masque filtre/aide/barre d’état Bubbles ; la liste des chapitres masque
son titre, son aide et sa barre d’état mais garde le filtrage. `list.New` initialise
`FilterInput.Prompt` à `Filter: ` ; le corps vide contient `No …` dans cette version.
Les descriptions de raccourcis par défaut sont également anglaises. Adapter les
propriétés exposées et rendre localement les états vides non paramétrables, sans
forker Bubbles. Une barre actuellement masquée n’est pas une fuite visible avérée.
Tester les états filtre actif/sans résultat, pas seulement une liste remplie.

## 4. Préférences : comportement vérifié et expérimentations

### Lecture, validation et sauvegarde actuelles

`preferences` contient `Version`, `EnabledTools`, `DefaultDifficulty`.
`readPreferences` utilise `Lstat` et exige un fichier régulier. L’absence renvoie
les valeurs en mémoire, niveau débutant et module `git` activé, sans créer le fichier.
Dans un fichier existant, une liste de modules absente vaut liste vide.

Le JSON est décodé en `map[string]json.RawMessage`. `version` est obligatoire,
entier non nul, égal à 1. La liste blanche n’accepte que `version`,
`enabled_tools`, `default_difficulty`. Les modules doivent être des chaînes non
vides et uniques ; un identifiant de module retiré du corpus reste conservé.
Une difficulté absente prend la valeur débutant ; une difficulté invalide est refusée.
Ce lecteur n’est pas un décodeur permissif des nouveaux champs.

`loadPreferences` renvoie des valeurs de secours **et une erreur** en cas d’échec.
La commande `prefs` en consultation affiche un avertissement sur stderr et les
valeurs de secours, puis termine actuellement avec le code 0. Une demande de
modification refuse alors de sauvegarder. Le secours après erreur n’active pas
automatiquement `git`, contrairement à l’absence normale de fichier.

`writePreferences` valide l’état demandé, relit le fichier, refuse un fichier
invalide/illisible ou un état changé depuis la lecture (`samePreferences`), vérifie
le droit d’écriture de la destination puis effectue un remplacement atomique via
`replacePreferencesFile` : temporaire, écriture, `Sync`, fermeture, renommage.
Le code suppose un seul écrivain ; ce contrôle n’est pas un verrou multi-processus.
Les collisions entre préférences et chemin de progression sont déjà contrôlées
dans l’application et doivent rester couvertes.

La CLI évite l’écriture si `samePreferences` ne détecte aucun changement.
La TUI distingue brouillon, état sauvegardé et annulation ; `tuiOptions` conserve
le snapshot initial et l’erreur de lecture dans les callbacks de chargement/sauvegarde.
L’ajout de langue doit toucher ces comparaisons, pas seulement la structure JSON.

Les fichiers sont sous le répertoire de données XDG, via `xdg.PreferencesFile`
et `xdg.ProgressFile`, **pas** sous le répertoire de configuration XDG.

### Protocole et résultats réellement exécutés

Binaire temporaire construit depuis le checkout, fixtures sous
`/tmp/cli-drill-p31.gFkMq2`, variables XDG isolées pour chaque cas. Aucun accès aux
préférences ou à la progression réelles. Tests de permissions exécutés comme
utilisateur non root sur macOS ARM64. Comparaison des octets avant/après ; contrôle
de l’absence de création de progression pour chaque commande de préférences.

Fixture de départ : `version: 1`, modules `fd` et `retired-cli`, difficulté
`advanced`. Les fichiers sont du JSON réel ; les valeurs ci-dessous sont abrégées.

| Cas | Résultat du binaire actuel | Attente pour la future localisation |
| --- | --- | --- |
| 1. Aucun fichier | Consultation réussie, aucune création. | `en` en mémoire, mêmes défauts pédagogiques. |
| 2. v1 valide sans langue | `advanced` et `fd` conservés ; avertissement pour `retired-cli`, octets inchangés. | Ajouter `en` en mémoire seulement. |
| 3. v1 avec `language: "en"` | Champ inconnu ; consultation avec avertissement, code 0 ; modification explicite refusée, fichier intact. | Accepter `en`. |
| 4. v1 avec `language: "fr"` | Même refus du champ ; aucune écriture. | Accepter `fr`. |
| 5. v1 avec `language: "de"` | Même refus du champ, pas encore une validation de langue ; aucune écriture. | Diagnostic spécifique de valeur invalide, aucun remplacement. |
| 6. JSON tronqué | Diagnostic JSON en consultation ; modification explicite refusée ; octets inchangés. | Conserver cette protection. |
| 7. Version 99 | Diagnostic de version non prise en charge ; modification refusée ; octets inchangés. | Conserver le refus d’une version inconnue. |
| 8a. Fichier illisible, mode `0000` | Avertissement de permission en lecture ; fichier préservé. | Diagnostic visible ; sauvegarde également refusée. |
| 8b. Fichier lisible non inscriptible, mode `0400` | `prefs --enable fzf` échoue, code 1 ; octets inchangés. | Conserver le refus malgré la possibilité technique d’un renommage dans le parent. |
| 9. Sauvegarde explicite après lecture v1 | `prefs --enable fzf` réussit et écrit un JSON v1. | Écrire seulement à la suite d’une modification demandée. |
| 10. Conservation des réglages | Après le cas 9 : `fd`, `retired-cli`, `fzf`, difficulté `advanced`. | Même conservation lors d’un choix explicite de langue. |

Les cas EN/FR n’ont pas testé une implémentation nouvelle : ils démontrent son
absence actuelle. L’écriture d’un fichier illisible est caractérisée par le code
de refus ; le cas `0000` a été exécuté en lecture seulement. Le futur comportement
de la langue, les erreurs de disque plein et les pannes de renommage restent à tester.

## 5. Comparaison des schémas de préférences v1/v2

| Critère | A — v1 avec `language` facultatif | B — v2 avec lecture des anciens v1 |
| --- | --- | --- |
| Simplicité | Une version reconnue, un champ autorisé et sa validation. | Deux versions à lire et règles explicites pour chacune. |
| Maintenance | Le changement de contrat est additif ; documenter que les anciens lecteurs stricts le refusent. | Discriminateur de format explicite, utile si une véritable restructuration survient. |
| Défaut EN et conservation | Compléter en mémoire un v1 valide sans toucher modules/niveau. | Lire un v1 en mémoire avec les mêmes garanties ; conserver le format source jusqu’à une écriture autorisée. |
| Écriture | Ajouter la langue lors d’un choix explicite ; aucune conversion nécessaire. | Sérialiser en v2 lors d’une modification explicite ; préciser le comportement pour un simple changement de module. |
| Validation | Liste blanche et valeurs `en`/`fr` ; absence distincte de `null`, vide ou type incorrect. | Même validation, plus combinaisons version/champ et conversion v1 vers v2. |
| Fichiers invalides | Protections actuelles conservées. | Aucune migration ne doit devenir un mécanisme de réparation. |
| Tests | Matrice absence/valeur/erreur et conservation. | Même matrice, plus lecture des deux formats et sauvegarde différée en v2. |
| Ancien binaire | Il refuse le champ nouveau. | Il refuse la version nouvelle. |

**Recommandation : A.** Le seul ajout de langue ne justifie pas ici une deuxième
version et un chemin de conversion. L’incompatibilité avec les anciens lecteurs
n’est pas un motif de complexification : elle est acceptée dans le périmètre
pré-release du 13 novembre. Informer de cette limite sans promettre de rétrogradation.
Si le mainteneur préfère une v2 pour identifier explicitement le nouveau format,
B reste viable, mais son bénéfice est documentaire plus que fonctionnel à ce stade.

### Mécanisme minimal recommandé

- Lecture : valider le document entier ; distinguer champ absent et valeur
  explicitement invalide. Rejeter `null`, `""`, `"EN"`, nombre et langue inconnue.
  Ne pas convertir ou corriger ces valeurs silencieusement.
- Garder séparés **la langue effective** (`en` si absente) et **le choix enregistré**
  (facultatif). Une représentation optionnelle dans `preferences` suffit ; elle
  ne nécessite ni fichier annexe ni version de progression.
- Un choix explicite `en` doit pouvoir être enregistré même si la valeur effective
  était déjà `en`. Sinon l’optimisation « inchangé » de la CLI/TUI perdrait ce choix.
- Lorsqu’un autre réglage change, préserver l’absence de langue si aucun choix de
  langue n’a été fait ; conserver tous les modules valides, y compris retirés du corpus.
- Sauvegarder via le mécanisme existant, en intégrant le champ optionnel aux
  comparaisons de snapshot/conflit et aux brouillons TUI. Aucun appel d’écriture
  depuis le chargement, l’aide, les statistiques ou le démarrage.
- Fichier invalide : diagnostic en anglais, car aucune préférence linguistique
  valide n’est disponible ; ne pas exploiter partiellement son champ `language`.
  Proposition de continuité avec le comportement actuel : consultation possible
  avec défauts annoncés, modification refusée. Le code de sortie de consultation
  n’est pas redéfini par cet audit.

## 6. Architecture de traduction de la CLI/TUI

### Options comparées

| Approche | Atouts | Coûts / risques |
| --- | --- | --- |
| Catalogue Go intégré et typé | Références contrôlées par le compilateur ; paramètres typés ; aucune nouvelle dépendance ni lecture de ressources à l’exécution. | Contributions rédactionnelles dans du Go ; le compilateur seul ne détecte pas une valeur vide. |
| Ressources JSON/YAML intégrées, IDs de messages | Traductions faciles à relire séparément ; format familier aux contributeurs du corpus. | Validation de clés/placeholders indispensable ; erreurs reportées au chargement/test ; générateur ou wrappers si l’on veut un typage équivalent. |

**Recommandation : catalogue Go typé pour l’interface**, ressources YAML pour la
pédagogie. Deux langues et des écrans déjà codés en Go ne justifient ni framework,
ni générateur, ni système de sélection automatique de locale. Le coût d’un second
format est limité : les catalogues UI et le contenu pédagogique ont des usages différents.

Un petit package partagé proposé `internal/i18n` porterait seulement le type de
langue et le catalogue UI. Il n’existe pas aujourd’hui. Sa justification est la
consommation par `app` et `tui`, sans dépendance inverse ni recours à une variable
globale mutable. Le stockage reste dans `app`, le corpus dans `chapter`, la
correction dans `chapter`/`exercise`, la progression dans `progress`.

Utiliser des noms sémantiques stables, par exemple `PreferencesSaved`,
`TrainingIncorrect`, `ChapterCount`, et des méthodes à arguments typés comme
`Summary(correct, total int)`. Une propriété traduite n’est jamais un identifiant
d’action ou un code d’erreur. Ne pas construire une phrase en concaténant plusieurs
fragments traduits ; gérer 0/1/plusieurs et l’ordre des paramètres au niveau du message.
Pas de `map[string]any` ni de propagation de dizaines de chaînes aux fonctions.

Prévoir des tests qui parcourent tous les messages des deux catalogues, interdisent
les champs vides, exercent les paramètres et détectent les résidus de formatage
(`%!…`). Des méthodes typées limitent les erreurs d’interpolation ; des tests
restent requis pour vérifier que chaque valeur utile apparaît dans le résultat.

### Transmission du contexte et démarrage précoce

**Fait :** les préférences sont chargées à plusieurs moments actuellement.
La racine vérifie le TTY avant la lecture TUI ; `chapters` ne charge pas les
préférences ; `stats` charge la progression avant celles-ci ; `reset` n’a pas besoin
des préférences aujourd’hui. Localiser seulement `tuiOptions` laisserait des sorties
anglaises et des erreurs précoces hors du contexte choisi.

**Recommandation :** résoudre le chemin puis lire les préférences une fois au
démarrage de la commande, avant la construction/localisation des aides, sans écrire.
Transmettre un contexte d’exécution contenant le catalogue et le résultat de
lecture via `app.Options`, `tui.Options`, le modèle TUI et `textTrainingView`.
Préserver l’injection de fichiers/IO dans les tests et éviter une lecture réelle
implicite dans un constructeur de commande testé. L’entrée `main` peut préparer
ce contexte ; les tests doivent pouvoir fournir leur contexte isolé.

Les erreurs du domaine devraient fournir un motif et des paramètres structurés
aux frontières concernées (champ, version, chemin, module, raison de filtre vide).
Les interfaces les rendent dans leur langue ; les causes restent exploitables
avec `errors.Is`/`errors.As`. Ne pas importer un catalogue dans chaque algorithme
ni changer les structures persistées pour traduire un diagnostic.

### Particularité Cobra : le hook pré-exécution arrive trop tard

L’ordre vérifié dans `cobra.Command.execute` v1.9.1 est : préparation des options
d’aide/version, parsing des flags, aide éventuelle, validation des arguments, puis
`PersistentPreRunE`. Un hook à ce dernier niveau ne peut donc pas localiser toutes
les erreurs ni `--help`.

Configurer le catalogue dès la construction : textes de commandes/options,
templates/fonctions d’aide et d’usage, commande d’aide, descriptions de completion,
validateurs d’arguments et traitement des erreurs de flags/résolution de commande.
La [documentation officielle Cobra v1.9.1](https://github.com/spf13/cobra/blob/v1.9.1/site/content/user_guide.md),
consultée via Context7 puis recoupée avec le module local, expose notamment
`SetHelpCommand`, `SetHelpFunc`, `SetHelpTemplate`, `SetUsageFunc`,
`SetUsageTemplate`, `SetErrPrefix` et les fonctions `Args`.
Changer le seul préfixe d’erreur ne traduit pas le message de Cobra/pflag.

Tester les diagnostics inconnus et les erreurs système : traduire le résumé et
les motifs courants, conserver chemins, tokens et codes techniques. Réafficher
systématiquement une cause anglaise complète derrière un résumé français ne
satisfait pas la couverture bilingue. Pour un échec inattendu, prévoir un message
localisé conservant des informations techniques utiles ; le format exact reste
à préciser lors du lot diagnostics, sans parser arbitrairement toutes les phrases
des dépendances ou perdre les codes de sortie.

### Sélection explicite : proposition à revoir

`cli-drill prefs --language en` / `cli-drill prefs --language fr` serait cohérent
avec `--difficulty`, `--enable`, `--disable`. **Syntaxe proposée, non disponible.**
Conserver initialement une modification par appel, suivant l’exclusivité actuelle.
Traduire les descriptions d’options, pas les noms des commandes ou les valeurs.

Dans la TUI, proposer un choix EN/FR dans l’écran existant, avec brouillon,
enregistrement et annulation. Recommandation : appliquer après sauvegarde réussie,
reconstruire les libellés et la projection pédagogique pour les prochaines séances ;
conserver l’ancienne langue si l’écriture échoue. Ne pas changer de langue au milieu
d’un lot actif. La reprise après capture doit réutiliser le contexte du lot pour
éviter des consignes et corrections de langues différentes. Ce moment d’application
reste un choix d’interaction à valider.

## 7. Corpus YAML : comparaison et exemples

### Chargeur actuel

`Chapter` contient ID, titre, description, famille, ordre et items. `Item` contient
ID, type, type d’exercice, difficulté, outils, consigne, `Choices []string`,
`Answer {Primary, Accepted}` et explication. `data.Chapters` expose actuellement
les `chapters/*.yaml` embarqués. `LoadFS` détecte les doublons de chapitres,
`Chapter.Validate` ceux d’items, les métadonnées invalides et les QCM mal formés.

L’expérience isolée avec le chargeur réel accepte un champ YAML inconnu
`unknown_translation: ignored`. Ajouter une seconde clé `title` produit bien une
erreur de doublon YAML. Les `UnmarshalYAML` de `metadata.go` rendent insuffisante
une hypothèse de validation stricte automatique de tous les nouveaux champs.
Il faudra une validation explicite des clés et de la forme des ressources,
y compris aux niveaux imbriqués, plutôt qu’ajouter des champs qui seraient ignorés.

### Options comparées

| Critère | A — textes EN/FR dans chaque YAML existant | B — ressources de langue séparées |
| --- | --- | --- |
| Organisation | 13 fichiers pédagogiques plus longs, textes proches du contenu technique. | 13 bases canoniques et, à terme, 26 ressources EN/FR par chapitre. |
| Relecture | Correspondance locale immédiate ; changements éditoriaux et techniques mélangés. | Revue par langue/chapitre facile ; nécessite de croiser avec la base. |
| Contribution | Un seul fichier pour un exercice, structure imbriquée plus dense. | Traduction sans toucher aux commandes/règles ; davantage de fichiers. |
| Validation | Présence obligatoire des deux blocs ; vérifier les champs partagés. | Jointure stricte par IDs, détection de clés orphelines/manquantes. |
| Divergence | Risque de dupliquer commandes ou métadonnées dans les blocs de langue. | Risque d’oublier une ressource après ajout/retrait dans la base. |
| Traduction progressive | Ajouter EN item par item. | Ajouter EN chapitre par chapitre, même clé que FR et base. |
| Taille du binaire | Les deux langues sont embarquées. | Même contenu utile embarqué ; séparation des fichiers sans gain de traduction intrinsèque. |

**Recommandation : B**, ressources symétriques, sans faire du français une source
de vérité cachée et de l’anglais une surcharge facultative. L’intérêt principal
est de protéger les éléments techniques et de permettre des lots de traduction
lisibles. A reste une solution raisonnable si le mainteneur privilégie fortement
la proximité de tous les textes ; elle exige les mêmes règles de correction.

### Échantillon vérifié

| Identité actuelle | Représentativité et frontière de traduction |
| --- | --- |
| `00-premiers-pas/pwd` | `simple-shell-sim` : traduire consigne/explication ; partager `pwd`. |
| `00-premiers-pas/cd-docs` | Plusieurs variantes : partager `cd docs` et `cd ./docs`, ainsi que `/atelier` et `docs` du scénario. |
| `00-premiers-pas/glob-no-match-shells` | QCM conceptuel : libellés traduits, même choix correct ; conserver les hypothèses Bash/Zsh, `nullglob`, `failglob`, `NOMATCH`, `NULL_GLOB`, `*.bak`. |
| `12-fichiers-permissions/missing-directory-search` | `free-answer` conceptuel : réponse actuelle `Exécution`, variantes `x`, `traversée`, `execution`, `exécution` ; équivalents EN à relire. |
| `01-raccourcis-terminal/ctrl-b-character-left` | `key-sequence` : partager `Ctrl+B`, `^B`, `control-b`, `ctrl b`, `←` ; traduire le contexte Readline/ZLE. |
| `00-premiers-pas/glob-bracket-list` | Contre-exemple : `type: concept`, mais réponse technique `plan-[135].svg` et ses cinq permutations. Le type ne détermine pas la traduisibilité. |

### Exemples compacts des deux structures possibles

**Illustrations seulement** : extraits incomplets, noms de champs proposés, aucune
conversion du corpus. Les traductions anglaises sont des exemples à relire.

Option A, textes au même endroit que l’exercice :

```yaml
id: cd-docs
exercise_type: simple-shell-sim
answer:
  primary: cd docs
  accepted: [cd ./docs]
text:
  en:
    prompt: From /atelier, enter its docs subdirectory using a relative path.
    explanation: cd changes the shell's working directory; docs leads to /atelier/docs.
  fr:
    prompt: Depuis /atelier, rejoins son sous-répertoire docs par un chemin relatif.
    explanation: cd change le répertoire courant du shell ; docs mène à /atelier/docs.
```

Option B recommandée, base canonique (métadonnées habituelles omises ici) :

```yaml
id: 00-premiers-pas
items:
  - id: pwd
    answer: {scope: shared, primary: pwd}
  - id: cd-docs
    answer: {scope: shared, primary: cd docs, accepted: [cd ./docs]}
  - id: glob-no-match-shells
    choices:
      - {id: both-pass}
      - {id: bash-pass-zsh-reject}
      - {id: bash-remove-zsh-pass}
    answer: {choice_id: bash-pass-zsh-reject}
```

Ressource anglaise correspondante, par exemple
`data/translations/en/00-premiers-pas.yaml` (chemin proposé) :

```yaml
chapter_id: 00-premiers-pas
title: First steps
description: Practise the fundamentals of shell navigation and expansion.
items:
  pwd:
    prompt: Print the current working directory without changing it.
    explanation: pwd prints the current path; logical and physical paths may differ for symbolic links.
  cd-docs:
    prompt: From /atelier, enter its docs subdirectory using a relative path.
    explanation: cd changes the shell's working directory; docs leads to /atelier/docs.
  glob-no-match-shells:
    prompt: >-
      No filename matches *.bak. Pass this unquoted pattern to a command with
      globbing enabled. Bash has nullglob and failglob disabled; native Zsh has
      NOMATCH enabled and NULL_GLOB disabled. What happens before the call?
    choices:
      both-pass: Bash and Zsh pass the pattern through.
      bash-pass-zsh-reject: Bash passes the pattern through; Zsh rejects the command line.
      bash-remove-zsh-pass: Bash removes the pattern; Zsh passes it through.
    explanation: Bash passes the pattern unchanged. Zsh reports no match and does not run the command.
```

La ressource `fr` a exactement les mêmes clés et porte les formulations françaises.
L’ordre des choix reste défini une seule fois dans la base. Pour un choix purement
technique (par exemple une commande), la base pourrait fournir son texte partagé :
un choix possède soit ce texte, soit un libellé dans chaque langue, jamais les deux.

Réponse conceptuelle localisée du véritable exercice `missing-directory-search` :

```yaml
# Dans la base du chapitre 12-fichiers-permissions :
id: missing-directory-search
answer: {scope: localized}
---
# Dans items.missing-directory-search de la ressource fr :
answer:
  primary: Exécution
  accepted: [x, traversée, execution, exécution]
---
# Dans le même item de la ressource en (proposition à relire) :
answer:
  primary: Execute
  accepted: [x, execute, execution, search, traversal]
```

Le contexte des propriétaires, modes `0740`/`0640`, absence d’ACL et droit de
traversée doit rester identique. Il ne faut pas élargir la question à un droit
différent pour rendre une traduction plus naturelle.

Pour `ctrl-b-character-left`, conserver dans la base :

```yaml
id: ctrl-b-character-left
exercise_type: key-sequence
answer:
  scope: shared
  primary: Ctrl+B
  accepted: [^B, control-b, ctrl b, '←']
```

Seuls `prompt` et `explanation` changent dans ses ressources. Pour
`glob-bracket-list`, même principe avec toutes les permutations techniques actuelles.

### Projection et validation recommandées

Le chargeur joint base + langue et produit une **vue de séance** des chapitres/items,
avec IDs et métadonnées inchangés, textes et réponses textuelles de la langue choisie.
Réutiliser les parcours, filtres et composants ; ne pas doubler les moteurs ni les
objets de progression. Adapter `data/defaults.go` pour embarquer les ressources
retenues ; ne dépendre ni du réseau ni du répertoire de travail de l’utilisateur.

Valider les deux langues complètes dans les tests, indépendamment de celle choisie :

- ensembles exacts de chapitres/items et langues `en`/`fr` ; clés inconnues,
  manquantes, doublons, valeurs nulles/vides interdites aux endroits requis ;
- titres/descriptions, consignes, explications et réponses localisées présents ;
- ensembles de choix identiques, IDs uniques, ordre canonique, un choix correct
  défini et présent dans les deux projections ;
- aucune réponse technique surchargée par une traduction ; absence de réponses
  localisées lorsqu’une réponse partagée est déclarée ;
- réponses/variantes de QCM résolues vers un choix existant et sans ambiguïté,
  jamais une réponse correcte qui ne corresponde à aucun choix affiché ;
- tous les exemples positifs/négatifs de correction exécutés dans les deux vues,
  et même liste de clés d’acquisition avant/après projection ;
- rejet explicite d’une ressource manquante lors de son chargement ; aucun repli
  silencieux FR au milieu d’une vue EN, ou inversement.

La structure permet de bloquer les altérations des réponses partagées. Elle ne
prouve pas automatiquement que chaque commande citée **dans la prose** est intacte.
Prévoir des fixtures de tokens/blocs techniques pour les consignes à risque et une
revue humaine côte à côte : options GNU/BSD/macOS, noms de fichiers, variables,
valeurs, guillemets, codes de sortie, portée des hypothèses et distracteurs QCM.
Les traductions de `MODE=avant`, d’un chemin français ou d’un identifiant technique
ne doivent pas changer ces données de scénario sous prétexte de passer à l’anglais.

La traduction peut progresser chapitre par chapitre sur une branche de travail.
Un état incomplet reste explicitement non publiable ; les contrôles peuvent
distinguer les fixtures de développement du corpus complet, mais le contrôle de
release doit toujours exiger 13 chapitres / 446 items EN **et** FR. Aucun indicateur
« terminé » fondé uniquement sur un nombre de fichiers ou sur leur présence.

## 8. Réponses, QCM et correcteur bilingue

### Routage actuel et invariants

`chapter.Item.MatchAnswer` est le point commun CLI/TUI : séquence clavier d’abord,
simulation shell ensuite, puis `type: shortcut`, `tool`/`concept`, enfin comparaison
littérale. Les fonctions correspondantes sont dans
[exercise.go](../internal/exercise/exercise.go) : `MatchKeySequence`,
`MatchShellAnswer`, `MatchConceptAnswer`, `MatchAnswer`.

La traduction ne doit pas modifier ce routage par accident, ni mettre toutes les
réponses en minuscules. `MatchShellAnswer` ne normalise l’espacement que pour des
mots simples et pipelines simples ; guillemets, échappements, expansions, autres
opérateurs et espaces significatifs conservent leur traitement actuel.
Les variantes acceptées explicitement restent déterminantes.

| Catégorie | Modèle recommandé | Vérification requise |
| --- | --- | --- |
| Commande shell | Une réponse partagée et mêmes variantes ; consigne localisée. | Même ensemble de cas positifs/négatifs EN/FR, casse/options/guillemets/espaces préservés. |
| Raccourci | Même réponse et mêmes octets/notations techniques ; traduction autour de la séquence. | Capture et saisie textuelle, Ctrl+I/Tab, flèches, Ctrl+C, Esc, collisions et restauration du terminal. |
| QCM | Choix à ID stable, texte localisé ou partagé, bonne réponse par ID. | Choix correct identique même si libellé et longueur changent ; aucune collision de résolution. |
| Concept en langue naturelle | Réponse principale et liste de variantes par langue. | Équivalence pédagogique et limites de tolérance explicites. |
| Concept technique / réponse mixte | Classification manuelle ; préserver la syntaxe, localiser seulement le texte réellement pédagogique. | Exemple `plan-[135].svg` partagé ; Markdown avec texte demandé potentiellement localisé. |

### QCM : retirer l’identité portée par le libellé

Aujourd’hui `Chapter.Validate` exige des choix textuels uniques et **au moins un**
choix accepté par `MatchAnswer`. Il ne garantit ni identité sémantique ni exactement
un choix correct sous toutes les normalisations. La TUI choisit un index mais
corrige son libellé ; la CLI reçoit une réponse textuelle.

Recommandation : stocker le bon `choice_id` dans la base. La TUI soumet une sélection
identifiée ; la CLI résout la saisie vers un choix de la projection courante, puis
utilise la même règle de correction. Ne pas accepter un ID interne tapé comme une
réponse correcte par un simple effet de l’implémentation.

Conserver `Item.MatchAnswer` comme entrée textuelle et lui adjoindre, si nécessaire,
une entrée de sélection typée qui rejoint le même comparateur de choix. Les deux
chemins doivent partager leurs tests ; ne pas laisser un correcteur parallèle TUI.
Pour les réponses primaires/acceptées actuelles d’un QCM, conserver toute variante
pertinente en tant qu’alias explicite du choix correspondant et vérifier qu’elle
ne résout jamais deux choix. Ne pas inventer une tolérance universelle de casse
pour des libellés qui sont des commandes techniques.

### Réponses conceptuelles et variantes

Le corpus contient réellement des réponses comme `Exécution`,
`0, 1, supérieur à 1` (`14-recherche-texte-unix/grep-exit-status`) et
`Seulement avec -L.` (`find-link-type-with-follow`). Il contient aussi des commandes
Micro, des globs et des fragments Markdown sous `type: concept`.

Expérience avec le correcteur actuel, sur réponses synthétiques : `répertoire`
exact est accepté ; `Répertoire` et `repertoire` ne le sont pas pour cette seule
réponse primaire ; `Directory` est accepté pour `directory`.
`MatchConceptAnswer` tolère en effet la casse seulement pour un nom composé des
caractères ASCII prévus. Ce constat ne prouve pas une nouvelle règle à appliquer
aux accents : il montre qu’une traduction peut changer la tolérance effective.

Recommandation initiale : conserver les règles de correction actuelles et fournir
des variantes explicites, relues, pour chaque réponse naturelle localisée.
Ne pas introduire de moteur sémantique. Si une tolérance Unicode limitée aux
réponses naturelles est souhaitée, elle doit être validée séparément et ne jamais
atteindre les commandes, chemins, motifs ou raccourcis.

Proposition : n’utiliser que les réponses de la langue choisie, plus les variantes
techniques partagées explicitement listées (par exemple `x`). L’union automatique
des réponses FR/EN augmenterait l’acceptation et les ambiguïtés sans besoin produit
établi. L’acceptation croisée des synonymes reste une question à soumettre au
mainteneur ; la même progression dans les deux langues n’impose pas cette union.

### Garantie des acquis

Le chemin actuel appelle `progress.RecordAnswer` avec `chapter.ID`, `item.ID` et
le résultat. Ni titre ni langue n’identifient un acquis. Conserver cette frontière.
Les résultats de lots gardent notamment l’item et la réponse saisie en mémoire :
figer leur langue jusqu’au bilan évite de montrer une correction étrangère à la
question posée. Changer de langue dans les préférences n’appelle aucune fonction
de sauvegarde ou de remise à zéro de progression.

Les futurs tests doivent charger une progression synthétique non vide, alterner
EN/FR et comparer les octets sans entraînement ; après une réponse, vérifier la
même clé et seulement l’évolution pédagogique normale des compteurs/historiques.
Ne pas ajouter `language` à la progression et ne pas renommer `10-workflows-dotfiles`.

## 9. Architecture recommandée et limites de responsabilité

Flux proposé, non implémenté :

```text
preferences.json --lecture seule--> app : état validé + langue effective
                                      |                  |
                                      |                  +--> catalogue UI EN/FR
                                      +--> chapter : base + ressources de langue
                                                           |
                                                   vue de séance stable
                                                           |
                                     CLI/TUI existantes + correction commune
                                                           |
                                           progression JSON v1 inchangée
```

`app` orchestre la lecture, le diagnostic initial et la sauvegarde explicite.
`tui` et les vues CLI reçoivent un contexte unique plutôt qu’un ensemble de chaînes.
`chapter` valide et projette le contenu, garde les identités, résout les choix.
`exercise` conserve les comparaisons techniques. `progress` reste indépendant
de la langue. Les DTO de préférences transmis à la TUI peuvent remplacer les
tuples actuels sans déplacer la responsabilité du fichier JSON.

## 10. Compromis et risques

| Risque | Réponse proposée / limite |
| --- | --- |
| Ancien binaire refusant les nouvelles préférences | Limite explicitement acceptée, sans garantie de rétrogradation ; aucune architecture de contournement. |
| Perte de modules/niveau lors d’une sauvegarde | Copier l’état validé, conserver les modules retirés et tester le changement de langue seul. |
| Écriture silencieuse de défauts | Séparer lecture/choix enregistré ; assertions sur octets, existence et dates de modification. |
| Égalité logique masquant un choix explicite `en` | Conserver la présence du choix enregistré et l’intégrer aux comparaisons. |
| Concurrence de fichiers | Conserver relecture et refus de conflit ; le mécanisme actuel reste mono-écrivain, pas une garantie transactionnelle générale. |
| Textes de bibliothèque non traduits | Tester aide, flags, états vides/filtre et erreurs précoces ; wrappers ciblés, pas de fork de framework par défaut. |
| Perte de sens pédagogique | Revue humaine des variantes, distracteurs, portabilité et hypothèses ; un contrôle de clés ne prouve pas l’équivalence. |
| Langues mélangées pendant un lot | Contexte figé jusqu’au bilan, puis reconstruction cohérente après sauvegarde. |
| Travail P4 perdu en traduction | Nettoyer les appellations publiques avant traduction de chaque chapitre ; IDs historiques conservés. |
| Taille et lisibilité EN/FR | Tests de rendu étroit et long, accents et cellules terminal ; aucune refonte visuelle requise. |
| Release bilingue incomplète | Validation stricte des deux catalogues et du corpus complet ; recette humaine et documentation de release cohérente. |

L’audit ne chiffre pas un budget de traduction et ne prétend pas avoir relu chaque
phrase anglaise future. Il ne démontre pas une prise en charge native supplémentaire.
Les propositions de structure ont été comparées sur des exemples ; elles n’ont
pas été implémentées ou benchmarkées.

## 11. Arbitrages à soumettre au mainteneur

1. **Schéma des préférences :** retenir l’extension v1 recommandée ou choisir v2
   pour sa séparation explicite de formats ; dans les deux cas, aucune écriture
   sur lecture et aucune garantie de rétrogradation.
2. **Catalogues UI :** valider le catalogue Go typé et son petit package partagé,
   ou préférer des ressources intégrées avec validation supplémentaire.
3. **Corpus :** valider les ressources séparées symétriques et les champs proposés,
   plutôt que les blocs de textes dans les YAML existants.
4. **QCM et réponses :** approuver les IDs de choix et la classification
   partagé/localisé ; préciser l’acceptation des réponses d’une autre langue et
   confirmer la conservation initiale des normalisations actuelles.
5. **Interaction :** valider la syntaxe candidate `prefs --language`, son exclusivité
   et l’application de la langue TUI après sauvegarde, hors lot actif.

Les noms précis de champs/méthodes seront ajustés dans les lots autorisés, à partir
de ces choix. Anglais par défaut, deux langues, progression v1, absence de migration
implicite et exclusion de la rétrogradation sont **déjà décidés**.
Les arbitrages licence, cibles natives, Wiki/release et bêta relèvent toujours de
P5/P6 ; l’audit ne les clôture pas.

## 12. Plan de première implémentation par petits lots

**Proposition de réalisation après discussion et validation ; aucun lot commencé.**
Ces lots précisent P3.1–P3.4 sans renuméroter les responsabilités P3–P6.

| Lot proposé | Périmètre et sortie reviewable | Critère de fin |
| --- | --- | --- |
| 1. Préférences et contexte | Champ retenu, lecture sans écriture, choix enregistré distinct du défaut, callbacks/DTO ; tests synthétiques. | Anciens réglages préservés, erreurs non réparées, choix `en` explicite persistable ; aucune modification de progression. |
| 2. Identités et catalogue UI | Actions/feedback séparés des textes ; catalogue typé ; contexte injecté dans app/TUI. | Navigation et correction inchangées, catalogue EN/FR complet pour le lot, absence de globale mutable. |
| 3. CLI et diagnostics | Aides/options/erreurs précoces, consultation et choix de langue, entraînement texte/capture. | Aide et erreurs avant RunE dans la langue validée ; absence/invalidité des préférences couverte. |
| 4. TUI | Accueil, chapitres, réglages, séances, bilans, statistiques, listes vides/filtre. | Parcours testés dans les deux langues, application cohérente après sauvegarde, rendu lisible. |
| 5. Modèle pédagogique pilote | Chargeur strict, IDs de QCM, réponse partagée/localisée ; cinq exemples représentatifs en fixtures, sans conversion massive. | Mêmes clés et même correction technique ; toutes les erreurs de ressource détectées. |
| 6. Corpus par chapitre | Nettoyage P4 éditorial préalable, extraction FR sans perte, traduction EN, variantes et revue humaine ; petits lots d’un chapitre ou sous-ensemble cohérent. | Checklist d’équivalence et tests du chapitre, aucun titre obsolète nouvellement traduit. |
| 7. Convergence P3.4 | Tous catalogues et 446 exercices, régressions, recette humaine EN/FR, preuves pour P5/P6. | Pas de traduction manquante ; mêmes acquis ; README/Wiki/captures préparés sur la version validée. |

Le modèle pilote peut être étudié après le choix du format, en parallèle des lots
CLI/TUI si les interfaces sont fixées. La traduction de production dépend du pilote
et du nettoyage éditorial. P5 peut préparer ses artefacts et P6 ses structures
documentaires en parallèle ; une release présentée comme bilingue exige la fin
de P3.4 et les autres critères de lancement, avec autorisation explicite de publication.

## 13. Matrice de tests et recettes

### Automatisation à ajouter lors de l’implémentation

| Domaine | Cas essentiels | Preuve attendue |
| --- | --- | --- |
| Lecture des préférences | Absence, v1 sans langue, EN/FR, chaîne vide/inconnue/majuscule, `null`, mauvais type, JSON/version/champ invalides. | Valeur/diagnostic correct ; aucun fichier créé ou modifié. |
| Environnement | `LANG`, `LC_ALL`, `LC_MESSAGES` contradictoires avec le défaut et le choix enregistré. | Langue indépendante du système. |
| Sauvegarde | Choix explicite EN depuis absence de champ, FR puis redémarrage, modules retirés, niveau avancé, annulation TUI, échec d’écriture. | Persistance seulement demandée ; autres réglages conservés ; état UI conservé si erreur. |
| Stockage | Fichier illisible/non inscriptible, répertoire/symlink, conflit depuis lecture, collision progression, échec de temporaire/renommage. | Refus explicite, pas de destruction ni progression touchée. |
| CLI/Cobra | `--help`, aide de sous-commande, commande/flag inconnu, argument manquant/en trop, exclusivité, options invalides, sans TTY, redirection, completion. | Tous messages contrôlés EN/FR, syntaxe technique et codes de sortie stables. |
| TUI | Navigation par IDs, deux catalogues, filtres sans résultat, préférence brouillon/sauvegarde/annulation, bilans et stats. | Même action malgré traduction, aucun résidu de langue hors éléments techniques. |
| Catalogue | Messages absents/vides, paramètres et pluriels 0/1/2, erreurs longues. | Échec de validation en cas d’incomplétude ; interpolations correctes. |
| YAML | Ressource/clé manquante ou inconnue, doublons, formes invalides, choix incomplets, réponse sans choix, surcharge technique interdite. | Refus déterministe mentionnant fichier, chapitre/item et champ. |
| Correction | Matrice existante shell/clé/concept en EN et FR ; QCM via sélection TUI et texte CLI, alias ambigus, réponse naturelle/mixte. | Même correction technique ; choix sémantique stable et variantes explicites. |
| Progression | Fixtures avec acquis/scores/tentatives/historiques, alternance EN/FR sans séance puis avec séance. | Octets inchangés au changement de langue seul ; même clé lors des réponses. |
| Corpus complet | Ensembles d’IDs, métadonnées, 13 chapitres et 446 exercices EN/FR, consignes/choix/explications/réponses. | Barrière de release sans fallback silencieux. |
| Affichage | Accents, textes longs, fenêtres étroites, défilement, couleurs désactivées. | Pas de troncature rendant consigne ou correction inaccessible. |

Réutiliser notamment [preferences_test.go côté app](../internal/app/preferences_test.go),
[preferences_test.go côté TUI](../internal/tui/preferences_test.go),
[choices_cli_test.go](../internal/app/choices_cli_test.go),
[choices_test.go](../internal/tui/choices_test.go),
[correction_test.go](../internal/chapter/correction_test.go),
[match_test.go](../internal/exercise/match_test.go),
[metadata_test.go](../internal/chapter/metadata_test.go) et les tests existants de
sessions, capture, rendu lisible et progression. Les tests doivent rester isolés
avec données synthétiques ; aucun entraînement ne doit exécuter une commande shell.

### Recette humaine future

Sur macOS/Ghostty : premier lancement EN sans configuration ; choix FR, sauvegarde
et redémarrage ; retour EN ; séance de commandes, QCM, concept et capture réelle ;
bilan puis révision du même acquis ; navigation, filtre vide, aide, petites fenêtres,
copie et défilement ; interruptions et retour correct au terminal. Comparer les deux
langues à partir d’une progression synthétique identique. Faire relire les consignes,
distracteurs et explications par une personne maîtrisant les deux langues.

Les validations natives des plateformes annoncées restent un travail P5 distinct.
Les GIF authentiques, le README EN/FR et le Wiki doivent suivre les commandes
effectivement implémentées ; la note de version expliquera alors l’anglais par
défaut et la commande validée pour sélectionner le français.

### Vérifications exécutées pendant cet audit

- Base Git et intégration P3.0 vérifiées avant création de branche ; aucune
  modification préexistante, aucun changement fonctionnel effectué.
- Lecture des sources, tests et contrats cités ; comparaison Cobra avec sa
  documentation versionnée et son code local ; inspection des textes Bubbles.
- `GOTOOLCHAIN=local GOPROXY=off go build -o /tmp/cli-drill-p31.gFkMq2/cli-drill ./cmd/cli-drill` : réussi.
- `GOTOOLCHAIN=local GOPROXY=off go test -count=1 ./...` : réussi pour les six
  packages internes ; `cmd/cli-drill` et `data` sans tests. Hôte macOS ARM64,
  Go `1.27.4` ; le minimum `go 1.23.0` déclaré n’a pas été testé avec cette version.
- Expériences de préférences décrites en section 4 : résultats et conservation
  vérifiés. Chargement/comptage du corpus, réponses conceptuelles synthétiques et
  clés YAML inconnues/doubles exécutés avec les packages actuels dans une copie
  temporaire issue de `git archive HEAD` ; aucun ajout de test dans le dépôt.
- Contrôle documentaire : 39 liens locaux résolus ; cohérence avec les décisions
  des 12/13 novembre relue. `git diff --check` et contrôle équivalent du fichier
  non suivi avec `git diff --no-index --check /dev/null` : sans erreur d’espacement.
  `git diff --stat` ne compte pas le fichier non suivi ; le diff sans index fournit
  son volume. Aucun fichier suivi modifié, index vide ; roadmap source inchangée.

Aucun prototype de localisation, aucune recette terminal humaine nouvelle, aucun
test natif Linux/BSD ni test d’installation de release effectué dans cette mission.
Le seul livrable du dépôt est ce rapport, laissé local et non indexé.

<a id="p31-comparaison-experimentale"></a>

## 14. Décisions et comparaison expérimentale du 13 novembre 2026

### Contrat validé par le mainteneur

La conception doit rendre l’installation et le choix de langue simples pour
l’utilisateur, limiter les changements Go nécessaires et permettre de relire
ou modifier un exercice avec peu de synchronisation. Un mécanisme
supplémentaire doit résoudre un problème observé dans le produit ou le corpus.
L’anglais sert de référence technique. Les commandes, options, identifiants,
syntaxes et valeurs Unix gardent leur forme usuelle. Le français reste naturel
et précis, y compris lorsqu’il conserve un terme technique anglais.

- L’anglais est le défaut universel. Le choix EN/FR est explicite et persistant,
  indépendamment de la langue du système. Dans la TUI, il prend effet après
  sauvegarde réussie, sans redémarrage et hors exercice en cours.
- L’orientation retenue est d’ajouter `language` comme champ facultatif de
  `preferences.json` v1. Une v2 demande un problème technique concret.
  L’absence du champ donne `en` en mémoire. Une valeur explicite invalide
  produit un diagnostic. La consultation préserve le fichier et la sauvegarde
  du choix conserve les modules et la difficulté. La progression JSON v1 et
  les clés `chapter-id/item-id` restent communes aux deux langues.
- La **politique A de correction** accepte les réponses naturelles de la
  langue sélectionnée et les variantes techniques communes. Un synonyme
  d’une autre langue doit être déclaré et justifié. Les comparateurs shell
  et raccourcis existants gardent leurs règles.
- Les choix QCM ont une identité distincte de leurs libellés. Les actions
  TUI et l’état de correction ont aussi leur propre identité.

[Décisions produit P3.1](DECISIONS.md#p31-decisions-linguistiques)

### Protocole et échantillon

Les deux chargeurs sont dans le package d’essai `p31probe`, dans une copie
temporaire de `eeab86818d5e91260479fc8fefc2ea4c96648753` créée avec
`git archive`. Les mêmes cinq exercices ont été extraits avec
`chapter.LoadFS(data.Chapters(), "*.yaml")`. Les fixtures reprennent les
textes français existants et des traductions anglaises de travail.

| Exercice | Ce que l’essai vérifie |
| --- | --- |
| `00-premiers-pas/pwd` | Commande simulée partagée. `pwd` est accepté dans les deux langues, `PWD` et `pwd -P` sont refusés pour cet énoncé. |
| `00-premiers-pas/glob-no-match-shells` | Trois libellés par langue, ordre canonique et bonne réponse par ID. Le test couvre le choix par index et la saisie du libellé. |
| `12-fichiers-permissions/missing-directory-search` | `Exécution` en FR, `execute` en EN, `x` dans les deux langues. Les réponses naturelles de l’autre langue testées sont refusées. La variante FR historique `execution` est explicitement conservée dans la fixture comme graphie sans accent. |
| `01-raccourcis-terminal/ctrl-b-character-left` | `Ctrl+B`, ses variantes historiques et l’octet de contrôle `0x02` sont acceptés dans les deux langues. |
| `00-premiers-pas/glob-bracket-list` | Motif technique `plan-[135].svg` et cinq permutations acceptées dans les deux langues. Une plage qui inclut d’autres chiffres est refusée. |

**A** décode un chapitre avec des blocs `en` et `fr` au niveau du chapitre
et de chaque item. **B** décode la même base technique et deux ressources
localisées, puis vérifie les références de chapitre et d’item avant de les
joindre. Les deux chargeurs partagent le décodage strict avec `yaml.v3`,
la projection vers `chapter.Chapter`, `Chapter.Validate` et les comparateurs
de `Item.MatchAnswer`. Pour un QCM, le code d’essai rattache le libellé ou
l’index à l’ID du choix et compare cet ID à la réponse canonique.

Les deux chargeurs valident EN et FR même lorsque seule la vue EN est demandée.
Les projections A et B sont identiques pour chaque langue. Les tests comparent
les titres, descriptions, consignes, explications, choix, identifiants, familles,
ordre pédagogique, types, niveaux et outils aux attentes ou aux sources. L’ajout
synthétique avec le module `git` couvre une liste d’outils non vide.

Les mutations testées sont refusées dans les deux formats pour une explication
FR manquante, un choix manquant ou inconnu, un ID de choix ou d’exercice
dupliqué, une clé YAML dupliquée, une bonne réponse inconnue, des libellés
ambigus après comparaison, un champ inconnu et une langue demandée inconnue.
B refuse en plus une ressource FR fournie sans document, un chapitre référencé inconnu
et une traduction d’item orpheline. Une progression v1 synthétique non vide
est chargée avec le lecteur actuel. L’acquis `00-premiers-pas/pwd` est
retrouvé dans les deux vues et les octets du fichier sont identiques après les
lectures et corrections en mémoire.

### Comparaison A/B sur le même échantillon

| Critère | Importance | A — blocs locaux dans le chapitre | B — base et traductions séparées |
| --- | --- | --- | --- |
| Simplicité pour l’utilisateur | Prioritaire | Même vue EN/FR que B. Un corpus embarqué permettrait le même parcours sans configuration. | Même vue que A et même possibilité d’embarquement. Le prototype ne mesure pas l’installation ni l’interface. |
| Simplicité d’implémentation | Prioritaire | 39 lignes Go propres au format après `gofmt`. Un document à décoder par chapitre. | 56 lignes Go propres au format. Trois décodages, jointure et contrôles de références. |
| Simplicité de maintenance | Prioritaire | L’ajout testé modifie un seul fichier. Les deux formulations et le contexte technique se relisent ensemble. | L’ajout testé modifie trois fichiers. Avec la base seule ou la base et EN, le chargement échoue jusqu’à l’ajout des traductions attendues. |
| Conservation des identifiants | Obligatoire | Les cinq clés et les métadonnées sont conservées. | Mêmes clés et métadonnées après jointure. |
| Fiabilité des corrections | Obligatoire | Les cas positifs et négatifs passent avec les comparateurs actuels et les IDs QCM. | Les mêmes cas passent avec le même code commun. |
| Validation EN/FR | Obligatoire | Blocs locaux obligatoires et contrôle des choix dans chaque langue. | Mêmes contrôles de contenu, plus la couverture des clés sur les deux ressources. |
| Lisibilité pour les contributeurs | Importante | Chapitre plus long, mais contexte et traductions visibles au même endroit. | Fichiers par langue plus courts. Une relecture bilingue et technique demande de naviguer entre les trois. |
| Nombre de fichiers et abstractions | À minimiser | 3 YAML et 157 lignes dans l’essai. À l’échelle actuelle, 13 YAML. | 9 YAML et 163 lignes dans l’essai. À l’échelle actuelle, 39 YAML. |
| Extensibilité à d’autres langues | Secondaire | Chaque langue ajoute des blocs aux chapitres et à la validation. | Chaque langue ajoute des ressources. Ce découpage aide surtout une équipe qui travaille séparément par langue. |

Les deux formats partagent 178 lignes Go d’essai. Les chiffres excluent les
tests, la génération des fixtures et l’intégration dans l’application. Ils
décrivent ces prototypes. Le coût de migration des 446 exercices et la vitesse
de relecture n’ont pas été mesurés. Les formulations anglaises demandent une
revue éditoriale avant intégration. L’essai couvre la correction d’un octet,
sans recette de capture sur terminal réel.

La suite ciblée `go test -count=1 -v ./p31probe` a validé les deux projections
et les mutations. Dans la copie temporaire, `go test -count=1 ./...` et
`go vet ./p31probe` ont réussi sur macOS ARM64 avec Go 1.27.4 et
`GOPROXY=off`. Les vérifications portent sur les fixtures et les interfaces
existantes. La future persistance du choix de langue et son application en
TUI restent à tester lors de l’intégration.

### Recommandation finale et suite

**Recommander A, sous réserve de l’approbation du mainteneur.** Les résultats
de projection, de validation et de correction sont équivalents. A réduit les
fichiers à coordonner lors d’un ajout et supprime la jointure chapitre/item.
Le compromis est la longueur des chapitres enrichis. Des blocs EN/FR adjacents
permettent de relire une modification au niveau de l’exercice. Les métadonnées,
réponses techniques et IDs de choix gardent une seule définition. Cette
recommandation actualise celle de la section 7, formulée avant les essais.

Pour la CLI et la TUI, conserver la recommandation du **catalogue Go typé**.
Un petit type partagé avec deux jeux de messages suffit. Le transmettre depuis
le démarrage aux interfaces existantes, en gardant les IDs d’action et les
résultats de correction distincts du texte. Les messages paramétrés reçoivent
des arguments typés. Les deux jeux de messages peuvent être vérifiés avec les
tests Go existants et leurs extensions ciblées, sans dépendance supplémentaire
ni état global mutable. Cette recommandation concerne l’interface. Le
prototype compare uniquement les formats du corpus.

Le mainteneur doit encore approuver **le format YAML**. La syntaxe candidate
`prefs --language en|fr` et ses règles de combinaison restent à valider.
Les noms de champs du prototype sont illustratifs. L’extension v1, la politique
de correction et le moment d’application dans la TUI ne demandent pas un nouvel
arbitrage, sauf découverte d’un obstacle technique concret pour le schéma.

Premier découpage proposé après validation humaine :

1. Ajouter le choix facultatif aux préférences v1 et vérifier la lecture seule,
   les valeurs invalides, la persistance et la conservation des réglages.
2. Relier un petit catalogue typé aux entrées CLI/TUI et séparer les actions et
   résultats des libellés. Vérifier l’application après sauvegarde et les erreurs
   précoces avec des données isolées.
3. Intégrer le chargeur du format approuvé et les IDs QCM sur un lot pilote.
   Vérifier les deux langues, les acquis et les comparateurs avant d’élargir la
   traduction par chapitre.
