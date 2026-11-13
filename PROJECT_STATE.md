# État du projet — CLI-Drill

**Au 12 novembre 2026, P2.6 est terminée, validée sur macOS avec Ghostty et
publiée sur `dev` à `91efab0`. Le corpus compte 446 exercices dans
13 chapitres.** [Rapport de recette et réserves](docs/P26_FINAL_ACCEPTANCE.md)

P3 est le prochain chantier fonctionnel. La
[décision P3.0](docs/DECISIONS.md#p30-politique-linguistique-et-lancement)
fixe l'anglais par défaut et approuve les objectifs du lancement public
bilingue. Les travaux P3.1 à P6 restent à réaliser.

P2.4 est publiée à `f096eec` après huit tests manuels
réussis dans Terminal.app.
P1.1 comprend le stockage initial (`25c344b`) et les corrections de
`fix/progress-storage-hardening`, acceptées après revue. La migration V2
complète reste à réaliser. P1.2 comprend l'implémentation `e7b4170` et les tests
complémentaires `a60c819`, sur `fix/training-recording`, depuis `dev` à `1f61cb9`.
Son intégration et sa publication sont réalisées dans `dev` à `be2bf2b`.
P1.3 est terminée et validée sur `fix/exercise-correction` depuis cette base.
Son intégration et sa publication sont réalisées dans `dev` à `db40d00`.
P1.4 est intégrée et publiée à `ef9b0a5`, après validation humaine et recette
réelle dans Terminal.app sous macOS. Le cadrage P2.0 décrit ci-dessous est
approuvé le 8 novembre 2026 pour commit, intégration par fast-forward dans `dev`
et publication de `origin/dev`. Cette finalisation est réalisée à `91f13c6`.
P2.1 part de cette base sur `feat/pedagogical-metadata`. La validation humaine
du 8 novembre 2026 autorise son commit, son intégration par fast-forward dans
`dev` et sa publication sur `origin/dev`, en conservant la branche de travail.
P2.1 est intégrée et publiée à `ab2096c`. P2.2 est intégrée et publiée à
`3ecdf9c`, après recette manuelle du 9 novembre 2026. P2.3 part de cette base
sur `feat/cli-preferences`. Après six tests manuels validés le 9 novembre
2026, elle est intégrée et publiée sur `dev` à `1e0c396`.

## Implémentation actuelle

L'interface présente CLI-Drill comme un outil d'entraînement local au terminal.

- Module Go à la racine ; point d'entrée `cmd/cli-drill/`.
- Le corpus publié depuis `ef5d6ba` compte
  446 exercices dans 13 chapitres, avec 149 Beginner, 169 Intermediate et
  128 Advanced. [Rapports par jalon](docs/CURRICULUM.md#enrichissement-par-lots)
- La base publiée à `d8b383a` inclut les neuf exercices de navigation et
  flux et porte le corpus à 262 exercices.
  [Rapport du lot](docs/CORPUS_P25_UNIX_BASICS_FLOWS.md)
- La base publiée `dev` à `c9ad59c` compte 253 exercices. Elle inclut
  les 13 ajouts en fichiers, permissions, processus et environnement.
  [Rapport du lot](docs/CORPUS_P25_UNIX_SAFETY_PROCESSES.md)
- La base publiée `dev` à `0178ab9` compte 240 exercices dans 13 chapitres.
  Elle inclut les 48 exercices sans module tiers du chapitre 14.
  [Rapport du chapitre 14](docs/CORPUS_P25_UNIX_TEXT.md)
- La base `dev` à `b3f2556` compte 192 exercices dans 12 chapitres. Elle inclut
  les 31 exercices de `01` intégrés à `a6af7b4`.
  [Rapport de la vague 3](docs/CORPUS_P25_WAVE3.md)
- La vague 2 publiée à `fd37fe4` comptait 169
  exercices. La vague 1 en comptait 121 à `f97718a` et le lot 2 en comptait
  75 à `4320650`. [Rapport de la vague 2](docs/CORPUS_P25_WAVE2.md)
  Les rapports des [lots 1](docs/CORPUS_P25_LOT1.md) et
  [2](docs/CORPUS_P25_LOT2.md) recensent les ajouts et validations.
  Les dix chapitres du corpus de référence P2.1 restent des fixtures de test.
- Le démarrage direct CLI et TUI choisit le premier chapitre textuel éligible.
  Avec les préférences initiales, `00-premiers-pas` propose un premier lot
  de six parmi ses douze exercices Unix textuels de niveau Beginner.
- Les choix QCM sont désormais affichés aussi en CLI texte ; la saisie reste un
  libellé ou une variante acceptée, sans nouveau système de correction.
- Corpus exclusivement intégré depuis P0.1 ; aucun chargement automatique des
  anciens YAML utilisateur, même lorsqu'ils sont invalides.
- Commandes publiques `chapters`, `train [chapter-id]`, `stats`, `reset --yes`
  et `prefs`, avec `--progress`, l'aide et la complétion standard Cobra.
  `train` accepte les choix ponctuels `--difficulty`, `--mode` et `--all`.
- Anciennes commandes `init`, `scan`, `generate`, `directory`, `search`, `show`
  et option `--config` retirées ; elles produisent désormais une erreur.
- Accueil TUI : Start training, Browse chapters, Keyboard shortcuts, My progress,
  Preferences et Quit.
- Bannière ASCII, palette, composants de liste et de saisie et pieds de page conservés.
- Deux moteurs d'entraînement conservés : texte dans la CLI/TUI et capture réelle
  dédiée. La TUI revient automatiquement à l'accueil après le mode de capture.
- Progression JSON version 1 et chemin inchangés ; stockage atomique P1.1
  réutilisé par P1.2 pour sauvegarder chaque réponse soumise.

### P0.1 — Corpus intégré exclusif, intégrée dans `dev`

Le chargeur partagé `loadChapters` lit uniquement le système de fichiers
embarqué transmis par `data.Chapters()`. Les anciens YAML ne remplacent ni ne
complètent le corpus et ne sont ni déplacés, ni convertis, ni supprimés.
Un chapitre présent seulement dans le dossier utilisateur est indisponible.

### P0.2 — Surface dotfiles retirée, intégrée dans `dev`

Le câblage Cobra, les aides et l'interface TUI ne lisent plus de configuration
TOML, n'exposent plus d'annuaire, de scan ou de génération et ne présentent plus
les messages d'installation des dotfiles. Les diagnostics et chemins associés
ont été retirés de l'interface, sans supprimer les fichiers utilisateur.

Lors de P0.2, Start training choisissait le premier chapitre textuel non vide,
soit `02-navigation-shell`, et ouvrait directement l'exercice. Si aucun chapitre
textuel n'est disponible, il ouvre le navigateur sans forcer la capture brute.
Le chapitre sélectionné reste cohérent avec la liste lors des retours arrière.
`train` sans identifiant choisissait le premier chapitre intégré, alors celui
des raccourcis. Depuis P2.2, ces deux accès ouvrent le nouveau
premier chapitre `00-premiers-pas`, sans changer leur code de sélection.

Depuis P1.4, Keyboard shortcuts lance directement le moteur existant lorsqu'un
seul chapitre contient un exercice `key-sequence` (le corpus actuel). Si plusieurs
chapitres de capture existent, le navigateur filtré est conservé ; s'il n'y en a
aucun, sa liste reste vide. La sélection conserve les indices du corpus complet.
À l'époque de P0.2/P1.4, les familles pédagogiques n'étaient pas encore
annotées ; P2.1 a depuis ajouté ces métadonnées. Les raccourcis uniquement
textuels des chapitres mixtes, comme Micro, ne sont pas ajoutés à cette vue
de capture. Browse chapters restaure la liste complète.

Les filtres reçoivent leurs lettres et le retour arrière avant les raccourcis
globaux. Entrée applique le filtre, puis permet de choisir ; Échap annule la
saisie ou efface le filtre appliqué avant de quitter le navigateur. Les pieds
de page reflètent ces contextes. Les dimensions de liste réservent l'espace du
titre, du pied de page et du padding pour éviter un écran tronqué.

My progress réutilise les statistiques existantes, sans compteur d'annuaire.
Aucun schéma, clé, compteur ou mécanisme de sauvegarde de progression n'est modifié.

### P0.3 — Internes hérités retirés, intégrée dans `dev`

L'architecture conserve six packages internes : `app`, `chapter`, `exercise`,
`progress`, `tui` et `xdg`, conformément à l'[architecture cible](docs/ARCHITECTURE.md).
Les cinq packages `detect`, `catalog`, `config`, `markdown` et `shell/zsh`,
leurs tests et leurs fixtures exclusivement internes ont été supprimés.
Aucun composant actif ne dépend plus du scanner, du parseur utilisateur ou du catalogue.

`chapter` conserve le chargement YAML (`Load`, `LoadFS`) et sa validation.
La génération, la conversion vers le catalogue et les lecteurs/écrivains de
répertoires externes ont été retirés. `Item.Type` est une chaîne conservant
les valeurs pédagogiques existantes ; la provenance `Source`, absente du corpus
intégré et utilisée seulement par le catalogue, a été retirée du modèle Go.
Aucun YAML intégré ni identifiant n'est modifié.

`xdg` ne résout plus que le chemin de progression, toujours
`$XDG_DATA_HOME/cli-drill/progress.json` ou, par défaut,
`~/.local/share/cli-drill/progress.json`. Les anciens chemins TOML/YAML ne sont
plus résolus par la production. Les tests construisent seuls leurs fixtures héritées.
La dépendance `BurntSushi/toml` a été retirée ; aucune version conservée ne change.

Dans `app` et `tui`, seuls des helpers et une option interne sans consommateur
réel ont été retirés ; un helper réservé aux tests a été déplacé dans les tests.
Le rendu, les parcours actifs, les moteurs, les contrôles techniques du terminal
et le stockage JSON restent inchangés. Aucun fichier utilisateur n'est lu,
déplacé, converti ou supprimé par cette refactorisation.

### P1.1 — Stockage de progression terminé, intégré localement dans `dev`

Sur la base P0.3, `Save` tronquait directement le fichier cible ; `Load`
acceptait `null` et des versions non prises en charge. Ces défauts sont corrigés.
Le stockage sérialise le même JSON version 1, indenté et terminé par un saut de
ligne, puis écrit un fichier temporaire dans le répertoire de destination.
Les permissions POSIX, la synchronisation du fichier (`Sync`) et sa fermeture
sont traitées avant le remplacement par `Rename`, atomique sur les Unix ciblés.
Un échec avant ce remplacement laisse le dernier fichier valide intact ; les
temporaires sont nettoyés sur les sorties d'erreur traitées. Si le système
refuse aussi leur suppression, cette erreur est retournée avec l'erreur initiale.

Un fichier absent produit l'état vide historique sans écriture. Les champs
optionnels absents gardent les valeurs par défaut existantes, dont la version 1
si elle est omise et une table de scores vide si elle est absente ou `null`.
Un JSON invalide, un objet racine `null`, une version autre que 1 ou une erreur de
lecture produisent une erreur contextualisée. `Save` refuse aussi de remplacer
un fichier existant illisible, invalide ou non pris en charge, même si l'appelant
lui fournit un état vide. Les identifiants et résultats des exercices retirés
restent présents après sauvegarde et rechargement.

La revue de `25c344b` a exposé des `null` structurels acceptés silencieusement,
le remplacement des fichiers en lecture seule et une couverture insuffisante
de `Sync`/`Close`. P1.1b rejette une version explicitement `null`, les scores
`null`, leurs compteurs `null`, un `streak` `null` et les acquisitions `null`
dans une liste. Les champs omis, les scores partiels et les conteneurs v1
`completed_exercises`/`chapter_scores` à `null` restent admis ; `last_session`
reste nullable. La validation typée existante continue à refuser les mauvais
types. Aucun schéma ni état historique valide n'est converti par migration.

Les nouveaux JSON ont le mode `0600` ; les bits POSIX d'un fichier existant sont
conservés. Avant remplacement, P1.1b ouvre la cible existante en écriture sans
troncature ni création, puis la referme ; un refus ou une erreur de fermeture
arrête la sauvegarde. Les destinations non régulières, notamment les liens
symboliques, sont refusées explicitement pour éviter de remplacer le lien à la place de sa
cible. Aucun schéma, chemin, score ou déclenchement de sauvegarde CLI/TUI ne
change. `reset --yes` conserve sa suppression explicitement demandée.

Garanties limitées aux fichiers locaux sous macOS et Linux, avec un seul
écrivain. Le remplacement change l'inode : les ACL, attributs étendus et autres
métadonnées spécifiques ne sont pas recopiés. Pas de verrouillage concurrent,
de synchronisation du répertoire ou de garantie contre une coupure électrique.
Un arrêt brutal empêchant les nettoyages peut laisser un temporaire voisin,
qui n'est jamais chargé comme progression ni purgé automatiquement.

L'injection de `Sync` et `Close` est privée et passée par appel au helper de
remplacement, sans état global mutable ni abstraction générale de fichiers.
Le nettoyage ferme toujours le vrai descripteur, même après une erreur injectée.

### P1.2 — Enregistrement des entraînements terminé et validé

L'audit sur `1f61cb9` a reproduit trois défauts : `tui.answerOrAdvance`
évaluait les réponses et le bilan sans mettre à jour la progression ;
`app.runTraining` sauvegardait seulement à la fin et perdait les réponses déjà
validées en cas d'erreur de lecture ; un EOF textuel vide créait une tentative
incorrecte. Les parcours clavier de Browse chapters et Keyboard shortcuts
délèguent au moteur CLI et héritaient de sa sauvegarde tardive.

`progress.RecordAnswer` partage uniquement l'enregistrement : charger l'état
v1 courant, appeler `MarkCompleted`, puis `Save`. Cette fonction concrète est
commune aux moteurs texte et clavier CLI et au texte Bubble Tea. Elle ne renvoie
un état utilisable qu'après sauvegarde réussie. Le stockage P1.1, ses protections,
ses chemins et son contrat à un seul écrivain sont conservés.

Règles conservées : une soumission correcte ou incorrecte ajoute une tentative,
une réponse correcte ajoute un succès et l'acquis `chapter-id/item-id` reste
unique. Les retries et nouvelles séances ajoutent des tentatives légitimes ;
les scores et identifiants retirés restent stockés. Le bilan de séance en mémoire
reste distinct des compteurs cumulés. Aide, solution, feedback et navigation ne
sont pas enregistrés. Un texte non vide suivi d'EOF est soumis une fois ; un
EOF vide abandonne sans écriture.

Choix explicité pour les validations vides : Entrée seule, espaces et fins de
ligne ne soumettent pas de réponse, dans les deux interfaces. Cela évite les
tentatives fictives sur la question suivante lors d'Entrées répétées. Les octets
de contrôle historiquement acceptés comme réponses textuelles, dont `Ctrl+L`,
ne sont pas effacés par ce filtre. P1.2 n'a pas modifié la normalisation
ni les réponses admises ; la correction est traitée séparément par P1.3.

La sauvegarde a lieu avant le feedback et avant la lecture suivante. Esc et
Ctrl+C abandonnent la saisie courante ; les réponses déjà sauvegardées restent
disponibles après retour au chapitre, menu, sortie et redémarrage. La CLI
textuelle conserve sa lecture par ligne : Esc peut nécessiter Entrée avant
d'être traité ; Ctrl+C natif peut interrompre le processus. Aucun passage
supplémentaire en mode brut ni gestionnaire de signal n'est ajouté.

La CLI remonte toute erreur de lecture ou sauvegarde. La TUI affiche le
diagnostic dans les vues existantes ; en cas d'échec, elle garde la réponse
éditable et n'avance ni le bilan ni la progression affichée. Entrée permet
de réessayer après correction de la cause, sans double comptage. Quitter à ce
stade ne sauvegarde pas la réponse en échec : aucun journal ni rejeu implicite
n'est ajouté. Une consultation ou un abandon sans réponse ne crée pas de JSON.

### P1.3 — Correction terminée et validée

L'audit sur `be2bf2b` a reproduit des faux positifs dans le comparateur commun :
`bat readme.md` pour `bat README.md`, `ls -a` pour `ls -A`, motifs de recherche
ou commandes de casse différente, espaces comprimés dans les guillemets,
pipes cités/échappés et séparateurs Unicode ou sauts de ligne confondus.
Les réponses principales, variantes YAML et bonnes/mauvaises options des QCM
étaient déjà reconnues ; aucun faux négatif du corpus n'a été confirmé.

`chapter.Item.MatchAnswer` sélectionne la règle pure de `exercise` pour les
deux moteurs CLI et la TUI, à partir des métadonnées existantes :

- `simple-shell-sim` conserve la casse. Seuls les espaces ASCII/tabulations
  entre mots simples et autour de pipes simples sont normalisés. Les mots
  de cette tolérance se limitent aux lettres/chiffres ASCII et `_-./:=`.
  Guillemets, échappements, expansions, redirections, `||` et autres syntaxes
  sortent de cette tolérance et sont comparés littéralement.
- `key-sequence` et les réponses textuelles de type `shortcut` conservent les
  octets de contrôle et les notations `Ctrl+A`, `^A`, `control-a`, `ctrl a`.
  Les lettres seules restent distinctes de Ctrl et respectent leur casse ;
  les séparateurs malformés tels que `ctrl++a` ne sont plus admis implicitement.
- Les noms pédagogiques `tool`/`concept` gardent une tolérance de casse limitée
  aux noms nus ASCII (lettres, chiffres, tiret, souligné). Les phrases, chemins
  et commandes ne bénéficient pas de cette tolérance. Une variante de contrôle
  explicitement admise, comme `Ctrl+R` pour Atuin, conserve ses notations clavier.
- Les autres réponses, dont les alias, fonctions et workflows textuels,
  sont littérales. Aucune équivalence supplémentaire n'est inventée ;
  les variantes déclarées, dont `MarkEdit`, restent acceptées.

Les espaces, tabulations de bord et fins de ligne de saisie restent tolérés,
sans supprimer un caractère final échappé. Le traitement des octets et notations
clavier conserve ses règles antérieures, notamment pour `Ctrl+I`.
Une réponse vide ne peut pas correspondre à une
variante vide. Les 68 exercices et toutes leurs variantes sont testés sans
modifier le YAML, les identifiants ou le contenu pédagogique. Des variantes
sémantiquement équivalentes mais non déclarées peuvent rester refusées,
notamment les changements de présentation des commandes complexes ou de
casse des libellés de workflow (`Marked 2` reste une formulation littérale).
Le cas `marked 2` refusé pour le nom d'application `Marked 2` est reporté
explicitement à P2 ; aucune réponse YAML n'est ajoutée dans P1.3.

Les tests ciblent casse, espaces significatifs, opérateurs, redirections,
guillemets, échappements, options, raccourcis et les parcours réels CLI/TUI.
Ils vérifient les scores et acquis exacts sur une progression v1 synthétique
avec historique, le rechargement et l'absence de double compte du feedback.
Sept mutations ciblées ont été détectées avec des overlays temporaires Go.
Le stockage et `progress.RecordAnswer`, la capture, la navigation, le schéma
JSON v1 et les dix YAML intégrés ne sont pas modifiés ; aucune dépendance
ni version supplémentaire n'est introduite.

Après revue, les deux anomalies importantes sont corrigées : les collages
multilignes sont refusés avant leur conversion en espaces par Bubbles ;
un diagnostic apparaît dans le champ, la saisie existante est conservée et
Entrée ne soumet rien tant qu'elle n'est pas corrigée. Les collages du terminal
et de `Ctrl+V` suivent ce contrôle, avec un message Bubble Tea privé pour le
presse-papiers. La bibliothèque `atotto/clipboard`, déjà présente, devient
une dépendance directe sans changement de version ni ajout de bibliothèque.
Les tabulations de bord non échappées sont acceptées dans les réponses
textuelles, y compris les variantes shell citées. Les tests ciblés couvrent
le refus sans écriture du JSON initial, la reprise après édition, les deux
interfaces et la préservation des tabulations significatives et de `Ctrl+I`.

Validation exécutée : `go test -count=1 ./...`, `go vet ./...`, `-race -count=1`
sur `app`, `tui`, `progress`, `chapter` et `exercise`, formatage Go et diff.
Un essai du binaire en CLI puis en TUI sous pseudo-terminal sur le chapitre
`07-lecture-preview` donne les mêmes six tentatives, cinq succès et acquis
exacts ; `bat readme.md` est refusé, la sortie TUI normale fonctionne.
Après corrections de revue, un nouvel essai sous pseudo-terminal confirme
l'acceptation d'une tabulation de bord, le refus multiligne avec diagnostic
et JSON inchangé, puis la reprise avec une seule sauvegarde et la sortie normale.
Ces essais utilisent exclusivement des progressions temporaires synthétiques ;
ils ne constituent pas une validation de tous les terminaux physiques.

### P1.4 — QCM TUI et accès direct aux raccourcis, terminée et validée

Les cinq exercices `multiple-choice` intégrés affichent tous leurs choix dans
l'écran d'entraînement existant. Le premier est sélectionné à l'ouverture ;
le marqueur `>` reste visible sans couleur. Flèches haut/bas et j/k déplacent
la sélection, bornée aux extrémités. Entrée soumet le libellé choisi à
`chapter.Item.MatchAnswer`, puis `progress.RecordAnswer`, sans modifier ces
mécanismes P1.1–P1.3. Le feedback affiche la correction et l'explication ;
Entrée poursuit normalement. Échap revient à la fiche et Ctrl+C quitte, sans
enregistrer le choix non soumis. La sélection repart du premier choix à chaque
exercice ou relance. Les collages ne déplacent pas la sélection.

Keyboard shortcuts évite le sélecteur et la fiche intermédiaires lorsqu'un
seul chapitre de capture existe, en réutilisant la délégation déjà employée
par la fiche chapitre. Start training garde le démarrage textuel existant ;
Browse chapters conserve l'exploration et les fiches. Le moteur de capture et
le retour automatique à l'accueil sont conservés.

Tests ajoutés : chaque proposition des cinq QCM, bonne/mauvaise réponse,
explication, navigation, bornes, feedback sans double compte, interruptions
avant/après réponse, erreur de sauvegarde et nouvelle tentative, transitions
vers les autres types textuels et relance. Les parcours clavier avec zéro,
un ou plusieurs chapitres restent couverts. Toutes les progressions sont temporaires.
`go test -count=1 ./...`, `go vet ./...`, contrôle du formatage Go et
`git diff --check` réussissent.

Essais du binaire sous pseudo-terminal à 100×35, avec et sans couleur : QCM
correct/incorrect, flèches et j/k, explication, continuation, Ctrl+C sans tentative
fictive, scores relus, lancement direct de la capture, Ctrl+A réellement reçu et
Échap vers l'accueil. Les attributs terminal sont restaurés hors du bit d'état
macOS `PENDIN` déjà documenté ci-dessous. Ces essais ne valident pas Terminal.app.

La recette Terminal.app rapportée par l'utilisateur avant P1.4 a validé bannière,
palette, navigation, aide complète, footer compact aux dimensions testées,
redimensionnement, thèmes, progression et capture. Aucun travail de layout n'est
engagé ici.

La recette réelle des nouveaux comportements dans Terminal.app sous macOS est
validée par l'utilisateur le 8 novembre 2026 : les huit tests manuels ont réussi,
couvrant l'affichage des QCM, la navigation entre propositions, les bonnes et
mauvaises réponses, l'enregistrement unique, l'interruption sans comptage,
l'accès direct à Keyboard shortcuts, la capture des touches et la restauration
du terminal. Aucun défaut bloquant ou important n'a été constaté. Cette validation
clôture P1.4 et la phase P1 (stockage, enregistrement, correction et TUI).

À la clôture de P1.4, le corpus YAML, les dépendances, le README et les
mécanismes de correction et de stockage étaient inchangés. L'affichage des
choix en CLI texte a ensuite été livré en P2.2.

## Limites encore présentes

Le corpus actif de P2.2 retire alias/fonctions personnels, scan/génération,
applications graphiques et questions redondantes. Les anciennes leçons conservées
comme fixtures ne sont jamais chargées en production. La recette manuelle valide la lisibilité et la qualité pédagogique de ce lot.
La recette P2.6 fournit deux mesures de durée. La durée sur d'autres
parcours reste à mesurer en usage. Le README actuel a déjà été simplifié
et présente l'entraînement de `dev`. La refonte publique EN/FR relève de P6.

Avant P2.1 (base `91f13c6`), `chapter.Item` n'avait pas de difficulté ni de liste
de modules et les chapitres n'avaient pas de famille ou d'ordre déclaré.
P2.1 ajoute ces métadonnées, P2.3 enregistre les préférences et P2.4 utilise
les deux pour préparer les séances.

Défauts connus encore au backlog :

- `streak` non mis à jour, `last_session` sans curseur de reprise.
- Les diagnostics et les contrôles de capacités du terminal gardent des cas
  différés au backlog. La capture dépend des touches transmises par l'émulateur
  et le système.

## Vision décidée et travaux suivants

**CLI-Drill — Learn the terminal by doing.** La [spécification](docs/SPEC.md)
organise un corpus vérifié en Unix & Shell Fundamentals, Keyboard Shortcuts,
Modern CLI & macOS, avec contexte portable/BSD/macOS/GNU/Linux/tiers explicite.
Les sessions visent cinq à dix minutes, sans détection d'outils installés.
Le périmètre P2 enseigne les commandes Unix, shells, raccourcis du terminal et
outils CLI ; la famille macOS ne justifie pas de leçons sur les logiciels graphiques.

La [décision V2](docs/DECISIONS.md#2026-11-07--cli-drill-v2--corpus-integre-et-retrait-des-dotfiles)
conserve intégralement les anciennes données JSON/YAML/TOML. Leur abandon
fonctionnel n'autorise aucune purge, réécriture ou migration automatique.
Les vérifications techniques du terminal restent requises ; les deux moteurs
et le fonctionnement local et hors ligne restent centraux.

P1 est terminée et validée, P1.4 comprise après recette réelle dans Terminal.app.
Les compléments différés restent ouverts dans le backlog ; ils ne sont pas
déclarés réalisés par la clôture de P1.

### P2.0 — Évolution de la feuille de route, terminée et validée

La demande du 8 novembre 2026 indique la validation de l'audit pédagogique des
68 exercices. Cette mission en formalise les nouvelles décisions produit, sans
refaire cet audit ni une revue exhaustive de P1. Les 10 YAML/68 exercices
restent inchangés ; aucune nouvelle qualité pédagogique ou capture n'est validée.

La [décision P2.0](docs/DECISIONS.md#2026-11-08--p20--evolution-pedagogique)
complète la décision V2 : sélection manuelle des modules CLI tiers, socle Unix
toujours disponible, difficulté Beginner/Intermediate/Advanced par exercice,
séances mélangées avec apprentissage guidé lorsque nécessaire, et enrichissement
progressif. Sa cible initiale de 100 à 120 exercices est remplacée par le
[nouvel objectif de P2.5](docs/CURRICULUM.md).
Le choix ne dépend pas d'une installation réelle, d'un scan ou des dotfiles.

La [spécification](docs/SPEC.md) définit les métadonnées additives et leur
compatibilité historique, les préférences minimales séparées du JSON v1, les
filtres et l'ordre stable pendant une séance. L'[architecture](docs/ARCHITECTURE.md)
les répartit dans les packages existants ; aucun nouveau moteur ou framework
n'est prévu. Les identifiants pertinents et les garanties d'enregistrement,
de correction et de stockage P1 sont maintenus.

Le [backlog P2](TODO.md#p2--feuille-de-route-pedagogique) ordonne métadonnées,
premier corpus révisé, préférences, séances, enrichissement et recette/documentation,
avec critères d'acceptation par sous-phase. Ce cadrage P2.0 ne réalisait aucune
des fonctions P2.1 à P2.6. P2.1 et P2.2 sont depuis validées et décrites
ci-dessous. P2.0 ne modifie ni Go, ni YAML, ni
dépendance, ni README et ne crée aucun fichier de configuration. La validation
de finalisation autorise uniquement le commit et la publication du cadrage P2.0,
sans commencer P2.1.

### P2.1 — Métadonnées pédagogiques, terminée et validée

Sur `feat/pedagogical-metadata`, depuis `dev` à `91f13c6`, `chapter.Chapter`
expose `Family` et `LearningOrder`, et `chapter.Item` expose `Difficulty` et
`Tools`. Les types chaîne et constantes `Difficulty`, `Family`, `LearningOrder`
restent dans `internal/chapter`. Aucun nouveau package ni dépendance n'est ajouté.
Le contrat et ses valeurs restent ceux de la [spécification](docs/SPEC.md#5-formats-dexercices).

Au chargement YAML, une difficulté absente devient `beginner` et un ordre absent
`guided`, uniquement en mémoire. Une famille absente reste vide, sans inférence ;
une liste d'outils absente ou vide n'ajoute aucun module requis. Les valeurs
de difficulté, famille et ordre explicitement vides, nulles, inconnues ou de
mauvais type sont refusées, de même que les listes d'outils mal typées et les
identifiants d'outils vides ou dupliqués. Les alias et fusions YAML sont résolus
avant les contrôles de type, sans assouplissement des valeurs. Les champs hérités
non reconnus restent ignorés, sans imposer de conversion aux anciens YAML.

`Load`/`LoadFS` et `Validate` diagnostiquent les types d'exercice inconnus,
identifiants de chapitre/exercice dupliqués et QCM sans choix cohérents avec une
réponse acceptée. Les erreurs identifient le champ et, selon la couche concernée,
le chapitre, l'exercice et le fichier. Les 10 YAML intégrés et leurs 68 exercices
restent inchangés sur disque, avec leurs identifiants et leur ordre historiques.
Le niveau Beginner chargé par défaut ne constitue pas une réévaluation éditoriale.

Les tests ajoutés couvrent défauts historiques, trois difficultés, trois familles,
deux ordres, outils multiples, champs hérités, valeurs invalides, alias/fusions,
doublons, QCM et la conservation de tous les champs historiques des 68 exercices
intégrés, comparés à un modèle miroir sans les nouveaux hooks de décodage.
Un test CLI compare un chapitre historique à sa version annotée Advanced avec un module fictif non installé :
mêmes sorties et règles de correction, mêmes acquis et compteurs v1, historique
retiré conservé, consultation sans écriture de progression ni configuration créée.
Les autres tests de correction et d'enregistrement CLI/TUI restent inchangés.

Contrôles exécutés sur macOS arm64, Go 1.27.3 : tests ciblés, puis
`go test -count=1 ./...`, `go vet ./...`, contrôle du formatage Go et
`git diff --check`. Les tests de refus des métadonnées invalides ont d'abord
échoué sur le chargeur antérieur, puis réussi après implémentation.
Ces résultats n'établissent pas une nouvelle recette terminal physique,
Linux ou Go 1.23. La validation humaine de P2.1 est accordée le 8 novembre 2026 ;
la finalisation exécute de nouveau les tests complets, vet, formatage et contrôles
de diff avant le commit, puis vérifie les tests après fast-forward.
La validation de P2.2 figure dans la section suivante.

Aucun code de `exercise`, `progress`, `tui`, de capture ou de navigation n'est
modifié. `chapter.Item.MatchAnswer`, `progress.RecordAnswer`, le stockage JSON v1
et ses protections restent inchangés. Les outils sont des références éditoriales,
sans détection d'installation, scan ou lecture des dotfiles. Les contextes de
plateforme/shell restent dans les champs textuels prévus par P2.0. Aucun filtre,
préférence utilisateur, mélange, nouveau contenu, YAML utilisateur importé ou
fichier de configuration n'est ajouté ; ces fonctions restent dans les
sous-phases suivantes.

### P2.2 — Premier corpus V2, terminé et validé

Sur `content/v2-first-corpus`, depuis `dev` à `ab2096c`, le corpus actif passe de
68 à 39 exercices, répartis en 9 chapitres : huit nouveaux fondamentaux et
31 situations historiques révisées, sans quota de remplissage. Les chapitres
personnels `03-alias-zsh` et `04-fonctions-zsh` sont retirés. Le chapitre Git
conserve son identifiant `10-workflows-dotfiles` pour ne pas casser ses acquis.
Les [décisions éditoriales et la table complète des clés](docs/CORPUS_P2.md)
font référence : 20 clés conservées, 48 retirées du corpus actif, 19 nouvelles.
Aucun résultat historique n'est déplacé, supprimé ou réattribué.

`00-premiers-pas` est le premier chapitre intégré : `Start training` et `train`
sans identifiant y accèdent immédiatement, sans configuration ni mode brut.
Ses huit commandes simulées sont `pwd`, `ls`, `ls -a`, `cd docs`, `cd ..`,
`cd /atelier`, `mkdir essais` et la lecture d'un nom avec espaces par `cat`.
À la livraison de P2.2, le premier prompt précisait l'absence d'exécution.
Chaque exercice fournit son propre contexte. Les variantes utiles sont
déclarées, sans équivalence shell générale.

Les trois familles sont jouables : Unix/Shell 11 exercices, raccourcis 8,
outils CLI 20. Niveaux : 34 Beginner, 4 Intermediate, 1 Advanced.
Formats : 25 commandes simulées, 8 captures, 4 QCM, 2 raccourcis Micro textuels.
Les outils tiers sont tous déclarés ; Git est un module de la famille Unix/Shell.
À la livraison de P2.2, la sélection des outils, les préférences et le mélange
restaient futurs. P2.3 ajoute les préférences, puis P2.4 filtre et ordonne les
séances.
Les consignes et explications comportent des retours à la ligne éditoriaux pour
éviter leur troncature dans la TUI à 80 colonnes, sans changement du rendu Go.
Les mini-leçons Bat et Historique/Git sont volontairement courtes ; une durée
uniforme de cinq à dix minutes n'est pas démontrée et reste à éprouver par l'usage.

Les raccourcis ciblent les bindings ZSH/ZLE Emacs de référence ; Ctrl+U reçoit
une nouvelle clé pour kill-whole-line. La capture reconnaît une touche, sans
exécuter son widget shell. Les wrappers et intégrations Zoxide, Yazi et Atuin
sont explicités, sans supposer leur présence. Le remplacement de `rg-bat`
sélectionne une ligne et ne prétend plus lancer Bat. Les références officielles
et les limites des pipelines figurent dans le document de corpus.

Le seul changement de production Go ajoute les choix YAML avant la saisie QCM
CLI. Il conserve les réponses principales et variantes, `MatchAnswer`,
`progress.RecordAnswer`, l'absence de tentative pour entrée/EOF vide et le
schéma JSON v1. La TUI, la capture et le stockage ne sont pas modifiés.
Les anciennes données utilisateur ne sont ni lues pour ces essais ni migrées.

Tests ajoutés : corpus complet, métadonnées explicites, politique des clés,
parcours par défaut CLI/TUI, soumission de toutes les réponses textuelles du
corpus courant avec historique retiré, non-exécution des commandes, QCM CLI.
Les tests QCM TUI couvrent toutes les propositions des quatre QCM courants.
Les tests historiques P1/P2.1 restent exercés sur les dix fixtures identiques à
la base, en particulier les 68 exercices sans métadonnées.

Sur macOS arm64 / Go 1.27.3 : `go test -count=1 ./...`, `go vet ./...`,
formatage Go et diff contrôlés. Les nouveaux tests du premier corpus et de
l'affichage QCM ont d'abord échoué sur la base puis réussi après modification.
Essais du binaire sous pseudo-terminal à 80 colonnes, sans couleur : parcours
Unix 8/8 avec nom de fichier contenant des espaces, QCM TUI incorrect puis
correct après relance, choix visibles et explication, Ctrl+A réellement capturé,
Échap vers l'accueil puis sortie normale. Les fichiers JSON synthétiques sont
relus pour les scores et identifiants exacts. La troncature initiale des longues
consignes a été constatée, puis corrigée par les retours à la ligne éditoriaux.

Le 9 novembre 2026, l'utilisateur a validé huit essais manuels dans Terminal.app
sans défaut bloquant. La recette couvre le premier lancement du chapitre Unix,
les huit exercices fondamentaux, la correction des commandes, les QCM en CLI,
la capture, les outils modernes dont `z docs`, la persistance après Ctrl+C et
la relecture dans un nouveau processus, ainsi que la lisibilité et la qualité
pédagogique. Les tests automatisés vérifient aussi la conservation des résultats
historiques, y compris ceux des exercices retirés, au format JSON v1.

La recette porte sur Terminal.app sous macOS. La durée d'apprentissage, Linux
et Go 1.23 restent à valider.

### P2.3 — Sélection manuelle et préférences minimales, terminée et validée

Sur `feat/cli-preferences`, depuis `dev` à `3ecdf9c` :

```sh
cli-drill prefs                           # afficher le niveau et les modules
cli-drill prefs --enable fd                # activer fd
cli-drill prefs --disable fd               # désactiver fd
cli-drill prefs --difficulty intermediate  # beginner | intermediate | advanced
```

Une seule option de modification par invocation. Une modification effective
est enregistrée immédiatement en CLI ; une valeur identique ne réécrit rien.
Dans la TUI, ouvrir **Preferences** depuis l'accueil (ou avec `p`). Les flèches
ou `j`/`k` déplacent le curseur, Entrée bascule un module ou fait défiler
les trois niveaux, `s` sauvegarde le brouillon et Échap revient à l'accueil en
abandonnant les changements non sauvegardés. Un échec laisse le brouillon
modifiable et affiche une erreur.

Les modules sont dédupliqués et triés à partir de `Item.Tools` ; aucun outil
n'a besoin d'être installé. Les identifiants retirés sont signalés, conservés
dans le JSON et ignorés pour l'éligibilité. `chapter.Item.EligibleForTools`
vérifie que tous les outils requis sont présents dans le corpus et activés ;
les exercices sans outil sont toujours éligibles selon ce critère.
P2.3 conserve la composition des séances de P2.2. P2.4 applique les
préférences aux séances et aux acquis affichés.

Le [contrat JSON v1](docs/SPEC.md#préférences-pédagogiques-p2-distinctes-de-la-progression)
est stocké dans `$XDG_DATA_HOME/cli-drill/preferences.json`, avec repli dans
`~/.local/share/cli-drill/preferences.json`. Sans fichier, le défaut est
`beginner` sans module actif, en mémoire uniquement. Un fichier invalide,
illisible ou d'une version inconnue produit un diagnostic ; la consultation
affiche les défauts et la sauvegarde est refusée. L'écriture explicite remplace
le fichier atomiquement avec des permissions privées (`0600` sur Unix
avec un umask standard). La relecture avant sauvegarde refuse un état différent
depuis l'ouverture de l'écran ; le stockage reste prévu pour un seul écrivain.

`--progress` ne déplace pas les préférences et ne peut viser le fichier réservé.
`reset --yes` continue de supprimer uniquement la progression. Les tests sur
fixtures contrôlent les octets de la progression v1 et des anciens YAML/TOML
lors des changements de préférences. Aucun identifiant, exercice, correcteur
ou moteur de capture n'a changé.

Sur macOS arm64 / Go 1.27.3 : `go test -count=1 ./...`, `go vet ./...`,
formatage Go et diff vérifiés. Essais du binaire sous pseudo-terminal, avec des données
temporaires et sans couleur : consultation sans création, accès par le menu
et `p`, brouillon annulé, sauvegarde et relecture CLI/TUI dans des processus
distincts, conservation d'un module retiré et d'une version inconnue, erreur
d'écriture visible, et socle accessible.

Recette réelle du 9 novembre 2026 dans Terminal.app : les six essais manuels
sont réussis, sans anomalie signalée. Ils couvrent les défauts sans création
de fichier, l'activation CLI de `fd` et `fzf`, le niveau `intermediate`, la
relecture après relancement, l'annulation d'un brouillon TUI par Échap, la
sauvegarde TUI par `s`, la désactivation CLI et la navigation. Les préférences
sont vérifiées dans un répertoire `XDG_DATA_HOME` temporaire, indépendant de
la progression habituelle. Après cette recette, P2.3 est publiée sur `dev`
à `1e0c396`.

### P2.4 — Séances filtrées, guidées et mélangées

Implémentée sur `refactor/session-order` depuis `1e0c396`, puis intégrée dans
`dev` et publiée sur `origin/dev` à `f096eec` après huit tests manuels réussis
dans Terminal.app.

`chapter.SessionFilter` applique d'abord les modules requis, tous activés
manuellement, puis le niveau exact. Il fournit les décomptes et prépare la copie
du chapitre consommée par les moteurs existants. Apprendre suit l'ordre éditorial
des exercices éligibles. Réviser et les séances indépendantes mélangent la copie
une fois au lancement. La TUI remet cette copie au moteur de capture. L'ordre
reste stable pendant la séance, y compris en cas de feedback, d'aide, de
redimensionnement ou de nouvelle validation après un échec de sauvegarde en TUI.
Le retry de la TUI réutilise la liste. La révision des erreurs en capture reste
une passe séparée.

```sh
cli-drill train
cli-drill train 00-premiers-pas --mode review
cli-drill train 02-navigation-shell --difficulty intermediate
```

Le niveau par défaut vient de `preferences.json` avec repli sur Beginner.
`--difficulty beginner|intermediate|advanced` remplace ce choix pour la séance.
Les chapitres guidés acceptent `--mode learn|review` et démarrent en Apprendre
par défaut. Dans la TUI, **Browse chapters** ouvre la fiche. `d` change le niveau,
`m` bascule le mode pour un chapitre guidé et Entrée lance la sélection éligible.
La fiche indique les effectifs par niveau et la cause d'une sélection vide.
Les choix ponctuels disparaissent à la sortie de la séance.

**Start training** et `train` sans identifiant cherchent le premier chapitre
textuel éligible. **Keyboard shortcuts** lance directement la capture pour un
seul chapitre éligible, ouvre le navigateur pour plusieurs et indique la raison
si la sélection est vide. Une capture exige une entrée et une sortie TTY. Après
une sortie normale de capture, l'accueil TUI reprend. Une sélection vide
n'enregistre ni réponse, ni préférence. Le socle reste consultable via
**Browse chapters**.

`stats`, **My progress** et les fiches séparent les acquis actifs des compteurs
historiques. La projection joint les clés du corpus éligible et les clés acquises.
Le score historique de la fiche cumule toutes les tentatives du chapitre.
Les réponses de séance conservent leurs clés `chapter-id/item-id` et enrichissent
le JSON v1 avec `progress.RecordAnswer`.

Les tests automatisés couvrent filtres, permutation, lancement, stabilité et
historique. Des essais macOS arm64 sous pseudo-terminal ont exercé le socle
guidé, une séance indépendante avec `fd` et `fzf`, les niveaux Intermediate et
Advanced, une sélection vide, les choix TUI, la capture avec et sans couleur,
la restauration du terminal et la relecture de données JSON v1 synthétiques.
Linux et la version minimale Go 1.23 restent à valider.
Un échec de sauvegarde en CLI ou en capture arrête la séance. La TUI textuelle
conserve la réponse pour une nouvelle validation.

### P2.5 — Fondamentaux Unix/Shell, lot 1 validé

Sur `content/v2-expansion-unix-01`, depuis `f096eec`, trois chapitres ajoutent
18 commandes simulées sans module tiers. Les 39 exercices existants ont été
relus avec conservation de leurs clés et réponses. Le corpus totalise 57
exercices. La recette humaine finale sous Ghostty a validé T01 à T08 le
10 novembre 2026, y compris T07. Le lot est intégré et publié sur `dev` à
`b809151`. Le [rapport du lot](docs/CORPUS_P25_LOT1.md) détaille ses validations.

### P2.5 — Outils CLI modernes, lot 2 validé et publié

Sur `content/v2-expansion-cli-02`, depuis `b809151`, 18 exercices enrichissent
six chapitres existants. Le corpus atteint 75 exercices dans 12 chapitres.
La recette humaine finale sous Ghostty sur macOS a validé les parcours
fonctionnels et la présentation le 10 novembre 2026. Le lot est intégré et
publié sur `dev` à `4320650`.
[Rapport du lot](docs/CORPUS_P25_LOT2.md)

Le premier référentiel prévoyait 360 exercices pour douze chapitres. Le
cadrage du 11 novembre 2026 conserve ce repère historique et fait dépendre
le volume et les niveaux de la couverture des compétences.
[Référentiel](docs/CURRICULUM.md)

### P2.5 — Repères, fichiers et droits, vague 1 publiée

Depuis `d945059`, la branche locale `content/p25-wave1-paths-files` ajoute
46 exercices sans module tiers : 22 dans `00-premiers-pas` (Learn guidé,
Review mélangé) et 24 dans `12-fichiers-permissions` (mélangé). Les 75
exercices historiques et leurs clés sont conservés. Le
[rapport de vague](docs/CORPUS_P25_WAVE1.md) documente la couverture, les
références, les tests et les essais isolés. La vague 1 est publiée sur `dev`
à `f97718a`.

### P2.5 — Flux, processus et environnement, vague 2 publiée

Depuis `f97718a`, la branche locale `content/p25-wave2-flows-processes` ajoute
24 exercices dans `11-flux-shell` (Learn guidé, Review mélangé) et 24 dans
`13-processus-environnement` (mélangé), sans module tiers. Les 121 exercices
de la base publiée restent identiques. Le
[rapport de vague](docs/CORPUS_P25_WAVE2.md) relie les 48 ajouts au référentiel
et détaille les tests. La vague 2 est publiée sur `dev` à `fd37fe4`.

### P2.5 — Édition de ligne dans Bash et Zsh, vague 3 révisée

Depuis `fd37fe4`, la branche locale `content/p25-wave3-line-editing` ajoute
23 exercices à `01-raccourcis-terminal`. Les gestes communs suivent la
référence Bash/Readline. Les différences utiles sont précisées pour Zsh/ZLE.
Les huit exercices historiques gardent leur identifiant et leur sens.

Apprendre suit l'ordre éditorial dans le niveau choisi et Réviser mélange
les exercices éligibles. En Intermediate, le navigateur propose une séance
mixte avec six captures et cinq QCM. En Advanced, le navigateur propose
huit diagnostics textuels. La recette physique sur des terminaux
représentatifs relève de P2.6.
[Rapport et limites de capture](docs/CORPUS_P25_WAVE3.md)

### P2.5 — Recherche et traitement de texte Unix, chapitre 14 intégré

Intégrée à `0178ab9`, la branche `content/p25-unix-text` ajoute 48 exercices
sur `grep`, `find`, `sed`, `awk` et leurs combinaisons. Apprendre suit
l'ordre éditorial du niveau et Réviser mélange la sélection. Les essais
sur fixtures couvrent les réponses de commande, dont les noms atypiques.
La durée de la séance Intermediate, qui contient 22 exercices, reste à
mesurer lors de la recette humaine. [Rapport](docs/CORPUS_P25_UNIX_TEXT.md)

### P2.5 — Sécurité des fichiers et processus, premier lot intégré

Le lot intégré à `c9ad59c` ajoute les sept indispensables de `12` et les six
de `13` prévus par l'audit. Les deux chapitres restent en ordre mélangé.
Les essais CLI sur une progression synthétique ont conservé les acquis JSON v1.
[Rapport du lot](docs/CORPUS_P25_UNIX_SAFETY_PROCESSES.md)

### P2.5 — Navigation et flux, deuxième lot intégré

Le lot intégré à `d8b383a` ajoute les trois indispensables de `00` et les six de `11`.
Les deux chapitres proposent un ordre guidé et une révision mélangée.
[Rapport du lot](docs/CORPUS_P25_UNIX_BASICS_FLOWS.md)

### P2.5 — Finalisation Unix, Git et outils CLI, corpus publié

Les 184 ajouts couvrent les 14 propositions utiles du socle et les sept
chapitres Git et outils CLI. `07`, `08` et `09` proposent des séances guidées
de less, Nano et Markdown, puis les compléments activables Bat, Micro et
Glow. Git suit aussi une progression guidée. La révision mélange les
exercices éligibles. Le corpus de 446 exercices est publié sur `dev`
depuis `ef5d6ba`. [Rapports par jalon](docs/CURRICULUM.md#enrichissement-par-lots)

### P2.6 — Lot 2, lisibilité et navigation

Le lot 2 a été publié dans `dev` à `0b8bddb`, après le lot 1 à `fb29a22`.
La TUI affiche le contenu à la largeur disponible et garde les
commandes de navigation ainsi que le champ de réponse textuelle accessibles.
`PgUp` et `PgDn` permettent de lire le contenu qui dépasse la hauteur.
Les flèches continuent de sélectionner les choix d'un QCM. La correction
garde la consigne et le redimensionnement conserve l'exercice et sa réponse.
Le navigateur annonce la recherche avec `/`. Les fiches et la CLI nomment
les groupes de modules à activer pour les chapitres bloqués par les filtres.

`chapters` utilise des totaux compacts ou présente l'identifiant et le titre
sur des lignes distinctes selon la largeur mesurée. En redirection, les
trois champs tabulés conservent leur format. Le lancement TUI demande un TTY
en entrée et en sortie et indique `cli-drill train` sans TTY ou avec
`TERM=dumb`. La capture avec `TERM=dumb` affiche ses écrans à la suite
sans séquence d'effacement. Avec `NO_COLOR`, la TUI interactive conserve
les contrôles de positionnement du curseur.

Sur macOS arm64, les tests Go, le détecteur de courses et `go vet` ont
réussi. Sous pseudo-terminal, le binaire de la base `fb29a22` coupe la
solution au milieu du test `awk` à 80×24. Le lot 2 rend sa fin et
l'explication consultables aux trois dimensions de recette, sans ajouter
de tentative au défilement ou au redimensionnement. Les essais couvrent
aussi un démarrage à 40×12, la recherche, le QCM, le mode sans couleur,
`TERM=dumb`, les flux redirigés et les régressions de capture du lot 1.

### P2.6 — Lot 3, séances courtes et progression

Le lot 3 est publié dans `dev` à `1f53b29`. La CLI,
la TUI et le parcours clavier utilisent le même ordre préparé, les mêmes
lots de six et le même décompte. `c` poursuit la visite et `r` reprend les
seules erreurs de la dernière passe. Le lot initial garde son score.
Le parcours intégral s'active avec `--all` ou avec `a` dans la fiche TUI.
Apprendre s'appuie sur les acquis JSON v1. Réviser tire dans tout le filtre.
Le curseur et les erreurs de la visite restent en mémoire.

Git est actif lorsqu'aucun fichier de préférences n'existe. Une préférence
enregistrée garde ses modules, y compris si Git est absent ou désactivé.
`stats` et My progress comptent les chapitres qui ont au moins une tentative.
Le remplacement atomique du JSON v1 suppose un seul écrivain. Des exécutions
simultanées peuvent perdre des mises à jour.

Sur macOS arm64 avec Go 1.27.4, les tests ciblés et complets, le détecteur
de courses et `go vet` ont réussi. Sous pseudo-terminal, la TUI à 40×12
a enchaîné un lot, une reprise et le lot suivant. Le parcours mixte a
enchaîné six captures, la reprise d'une erreur et les QCM. Les scénarios
de restauration du terminal du lot 1 et de lisibilité aux trois dimensions
du lot 2 ont réussi. Les durées et les résultats sur clavier physique
figurent dans le rapport du lot 4.

### P2.6 — Lot 4, recette finale validée sur macOS

P05 a été publié à `3c4d63d`. La clôture de P2.6 est publiée à `91efab0`.
Le mainteneur a terminé H01 à H11 sous macOS
avec Ghostty et confirmé après correction les flèches gauche et droite,
Ctrl+B, Ctrl+F et le séparateur de la CLI textuelle. Les séances Apprendre
et Réviser ont duré 1 min 20 s (6/6) et 1 min 58 s (4/6).

Le compteur TUI `Exercice 1/12` est confirmé sous pseudo-terminal. Son
affichage en recette physique reste non confirmé. Le terminal a été
restauré physiquement après SIGTERM. Le statut `143` a été vérifié sous
pseudo-terminal et n'a pas été relevé dans Ghostty. Linux et BSD ont été
compilés en croisé pendant le lot 4. Leur recette native reste à faire.
[Rapport de recette](docs/P26_FINAL_ACCEPTANCE.md)

## Préparation du lancement public

Le corpus est en français et la CLI/TUI mêle des libellés français et
anglais. Le logiciel n'a pas encore de réglage de langue. La
[spécification P3](docs/SPEC.md#politique-linguistique-p3) fixe le
comportement futur, avec choix EN/FR persistant, préférences existantes
préservées et progression JSON v1 commune.

Le retrait fonctionnel du scan et de la génération appartient aux travaux
P0 intégrés. Le nettoyage éditorial P4 reste à mener. Les distributions
officielles, le nouveau README public en deux fichiers et le Wiki bilingue
relèvent des phases P5 et P6. Le
[backlog de lancement](TODO.md#lancement-public-p3-p6) précise les travaux
parallélisables et les conditions de sortie. La recette native Linux et BSD
reste à faire. Les compilations croisées P2.6 sont consignées dans le rapport.

CLI-Drill est en pré-release, sans utilisateurs externes. La
[décision du 13 novembre](docs/DECISIONS.md#p30-compatibilite-des-preferences-en-pre-release)
exclut la garantie de rétrogradation. La mise à niveau doit préserver les
préférences valides et le schéma JSON v1 de progression.

Restent à instruire le format des ressources, le choix du schéma de
préférences et sa prise en charge à la mise à niveau, la licence, les
versions minimales et cibles de support, la forme du tap Homebrew, les
contrôles bloquants de release et les modalités de lancement, avec bêta
technique ou sortie publique unique.

## Validation et limites de preuve

### P0.3 — Historique

Les contrôles de référence sur la base P0.2 ont réussi avant modification.
Après suppression : `go test -count=1 ./...`, `go vet ./...`,
`go run ./cmd/cli-drill --help`, `gofmt -w cmd internal data`, `go mod tidy`
et `git diff --check`. La cartographie des imports confirme l'absence des
packages retirés et du code de génération. Le corpus, `exercise`, `progress`
et les fonctions actives de rendu/capture ont été comparés à la base Git.

Deux tests de caractérisation ajoutés passent avant et après le découplage :
chargement des dix chapitres/68 exercices dans le même ordre avec leurs
métadonnées pédagogiques, et chemin JSON historique par défaut, avec surcharge
XDG et sans répertoire personnel disponible. Les tests de préservation des
anciennes données et de la surface V2 restent en place ; leurs chemins de
fixtures ne dépendent plus de helpers de production supprimés. Le test du pied
de page clavier vise désormais la méthode réellement utilisée par le rendu.

Le binaire a aussi été essayé dans un pseudo-terminal isolé à 80×24 : premier
lancement sans configuration, entraînement textuel immédiat, menu paginé,
progression, capture réelle de `Ctrl+A`, sortie `Esc`, retour à l'accueil et
restauration des modes du terminal. Les commandes `chapters`, `train`, `stats`
et `reset --yes` ont été exécutées sur des fichiers temporaires distincts.
L'historique JSON version 1 est chargé et conservé lors de l'ajout de résultats.
Les anciens TOML/YAML synthétiques, même invalides, sont ignorés et leurs octets,
ainsi que ceux du JSON historique de référence, restent identiques après les essais.
Aucune vraie donnée utilisateur n'est lue ou modifiée pour ces validations.

La comparaison des attributs du pseudo-terminal exclut seulement `PENDIN`,
[drapeau d'état macOS](https://github.com/apple-oss-distributions/xnu/blob/main/bsd/sys/termios.h)
lié à l'entrée en attente. Ce même écart est reproduit avec le binaire de la base
`14b12bd` ; tous les autres attributs sont restaurés. La capture n'a pas été
modifiée pour masquer cet écart dans le code applicatif.

Ces validations ne couvrent ni la V2 complète, ni tous les émulateurs et
raccourcis. La compatibilité avec le minimum Go 1.23 n'est pas vérifiée ;
les contrôles de développement utilisent Go 1.27.3 sur macOS arm64.

### P1.1 — Validations initiales, antérieures à la revue

Tests de référence sur P0.3 puis contrôles après modification :
`go test -count=1 ./...`, `go vet ./...`, formatage Go et `git diff --check`.
Les tests ajoutés ont d'abord exposé l'écrasement des JSON invalides, l'acceptation
de versions non prises en charge, la troncature de l'ancien inode et la création
en `0644`, puis réussi après correction.

La suite couvre absence de fichier, création des répertoires, aller-retour du
JSON historique et de ses cinq champs, valeurs par défaut, erreurs de lecture,
sérialisation et écriture, permissions et refus des liens symboliques.
Des erreurs provoquées après une écriture partielle et sur un fichier fermé
vérifient l'intégrité de la destination et le nettoyage du temporaire ; un
renommage réellement refusé vers un répertoire teste aussi ce nettoyage.
Un refus de nettoyage par permissions est signalé avec l'erreur d'écriture.
Les tests de permissions Unix sont exécutés sur macOS avec un utilisateur non
root. Les tests de stockage sont compilés pour Linux amd64, sans exécution Linux.

Le binaire a été essayé avec un HOME et un XDG isolés : aide, `chapters`,
`stats`, six réponses de `train`, conservation des résultats retirés et du mode
`0640`, nouvelle progression en `0600`, refus des JSON invalides et reset avec
confirmation. Le lancement et la sortie de la TUI n'écrasent pas un JSON invalide.
Aucun consommateur, YAML, moteur clavier, règle pédagogique ou dépendance ne
change. Toutes les données de test sont synthétiques et temporaires.

### P1.1b — Corrections après revue, acceptées et intégrées localement

Les tests ajoutés avant correction ont reproduit l'écrasement des `null`
invalides et des fichiers `0400`/`0444`. Les tests applicatifs passent par
`Load`, les commandes `stats`/`train` et le refus direct de `Save`, avec contrôle
des octets d'origine après chaque action. Cinq fixtures positives couvrent les
omissions historiques, conteneurs `null`, scores partiels et un JSON v1
complet enrichi par un entraînement, sans perte des résultats retirés.

Les tests de remplacement atteignent séparément les erreurs injectées de `Sync`
et de `Close`, vérifient l'intégrité de la destination et le nettoyage, ainsi
que l'ordre nominal écriture/synchronisation/fermeture avant renommage.
Des suppressions de temporaires réellement refusées par les permissions
vérifient aussi la remontée conjointe de l'erreur initiale et du nettoyage.
Les injections sont exercées en parallèle sur des fichiers isolés.

Contrôles réussis sur macOS arm64, Go 1.27.3, utilisateur non root (UID 501) :
`go test -count=1 ./...`, `go vet ./...`, formatage et `git diff --check` ;
`go test -race -count=1 ./internal/progress ./internal/app` réussit également.
Les tests de permissions ne sont pas sautés. L'aide CLI et le refus de `reset`
sans confirmation puis `reset --yes` sont vérifiés sur une fixture temporaire.

Cinq mutations par overlays temporaires font échouer les nouveaux tests :
suppression de la validation structurelle, du contrôle d'écriture, de `Sync`,
de `Close`, et fermeture déplacée après renommage. Les sources du checkout ne
sont pas modifiées par ces essais. Aucune donnée utilisateur réelle n'est utilisée.
Ni Linux, ni Go 1.23, ni une nouvelle session terminal interactive ne sont
validés par P1.1b ; les limites de durabilité et de concurrence restent celles
décrites plus haut. La revue des corrections est acceptée et leur intégration
locale est autorisée. La finalisation de P1.1 comprend une dernière exécution
des tests complets, de `go vet ./...` et de `git diff --check`, puis un commit
des corrections et un fast-forward exclusivement vers `dev`. Aucun push ni
travail P1.2 ne fait partie de cette finalisation.

### P1.2 — Validation de l'enregistrement acceptée

Avant modification : arbre propre, deux commits P1.1 présents, fetch de `origin/dev`
et égalité avec `dev` à `1f61cb9`, tests complets sans cache et vet réussis.
L'index QMD ne renvoyait pas les règles de tentative ; les documents originaux
du checkout ont servi de référence. Les défauts ont été reproduits avant correction.

Les nouveaux tests couvrent les quatre types textuels par un même parcours et
le moteur clavier séparément : réponses correctes/incorrectes, enregistrement
avant la lecture suivante, sessions complètes/partielles, EOF vide et non vide,
Esc/Ctrl+C, aides et solutions, retries et révision, feedback et Entrées répétées,
nouvelles séances et rechargement. Ils passent par les commandes Cobra ou les
événements Bubble Tea, avec lecture du JSON réel après soumission. Les résultats
historiques retirés, les cinq champs v1, `stats` et `reset --yes` sont contrôlés.
Les erreurs de stockage sont exercées sur des liens synthétiques refusés par
P1.1, et la TUI est testée après correction de la cause et nouvelle validation.

Les tests échouent avec les sources de `1f61cb9` en overlay temporaire. Deux
mutations supplémentaires sont détectées : suppression de la sauvegarde commune
et suppression de la garde de feedback TUI. Aucun fichier utilisateur réel
n'est utilisé, et aucun code d'exercice n'est exécuté.

Contrôles finaux réussis sur macOS arm64 avec Go 1.27.3 : formatage des seuls
fichiers Go modifiés, `go test -count=1 ./...`, `go vet ./...`, `git diff --check`
et `go test -race -count=1 ./internal/app ./internal/tui ./internal/progress`.

La revue finale n'a identifié aucun défaut bloquant ou important. Les deux
lacunes mineures de tests sont corrigées dans `a60c819` : identifiants acquis
exacts, commande produisant réellement `tea.QuitMsg` et sauvegarde échouée suivie
d'une interruption sans progression artificielle ni modification du JSON.
Trois mutations passaient les anciennes assertions et échouent avec les tests
renforcés : identifiants permutés, commande non nulle ne quittant pas et progression
artificielle après erreur. La validation humaine finale est accordée le
8 novembre 2026 ; les tests complets sans cache, vet, race et contrôles de diff
et formatage ont de nouveau réussi avant le commit complémentaire.

Le binaire est aussi essayé dans un pseudo-terminal macOS isolé : Start training,
Browse chapters textuel, Browse chapters clavier, Keyboard shortcuts et retour
à l'accueil, erreur de sauvegarde TUI visible, interruption CLI par SIGINT et
abandon CLI par Esc suivi d'Entrée. Les sept scénarios vérifient le JSON et la
sortie du processus. Cela ne valide pas un émulateur macOS particulier, les
raccourcis interceptés par le système ni tous les états du terminal.

Les garanties contre coupure électrique, les écritures concurrentes et les
arrêts avant déclenchement de la sauvegarde restent exclues. Aucun curseur de
reprise ni stockage des réponses en échec n'est ajouté. Linux et Go 1.23 ne sont
pas validés par ces essais ; P1.3/P1.4 restent hors périmètre.

Pour une vérification complémentaire dans un terminal natif, depuis la racine
du dépôt (procédure fournie, non exécutée dans un émulateur natif) :

```sh
recording_dir="$(mktemp -d)"
go run ./cmd/cli-drill --progress "$recording_dir/progress.json"
go run ./cmd/cli-drill --progress "$recording_dir/progress.json" stats
```

Dans la TUI : Start training, saisir `zoxide`, Entrée, Esc, `g`, puis `q`.
Les statistiques doivent annoncer un exercice acquis et un chapitre pratiqué.
Relancer avec le même `--progress`, vérifier My progress, puis pratiquer le même
chapitre via Browse chapters : une nouvelle réponse ajoute une tentative,
sans dupliquer l'acquis. Keyboard shortcuts, `Ctrl+A`, puis Esc doit enregistrer
la réponse avant le retour à l'accueil. Tous ces essais restent dans le dossier
temporaire explicite.
