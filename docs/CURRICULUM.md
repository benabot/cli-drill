# CLI-Drill — Référentiel pédagogique P2.5

**Objectif : au moins 10 Beginner, 10 Intermediate et 10 Advanced par chapitre.**
Pour les 12 chapitres actuels, le minimum est de **360 exercices intégrés**. Ce
seuil remplace la cible globale de 100 à 120. Il ne s'applique pas à une séance :
le niveau exact et les modules choisis déterminent toujours l'éligibilité.

Le point de départ est `dev` à `4320650`, après publication du
[lot 2](CORPUS_P25_LOT2.md). Le relevé porte sur les YAML intégrés de
[`data/chapters`](../data/chapters/), sans compter les fixtures. Les identifiants
03 et 04 ont été retirés. Les garder réservés et conserver les identifiants
existants évite de réattribuer des acquis à d'autres apprentissages. L'[état du projet](../PROJECT_STATE.md)
décrit le contenu livré ; le [backlog](../TODO.md#p25--enrichissement-par-chapitre-et-par-niveau--en-cours)
suit les travaux restants. Les compétences ci-dessous sont des propositions,
pas des exercices déjà écrits, testés ou publiés.

## Bilan par chapitre

B = Beginner, I = Intermediate, A = Advanced. Les écarts sont calculés
séparément par chapitre et par niveau.

| Chapitre | B | I | A | Total | B à ajouter | I à ajouter | A à ajouter |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| `00-premiers-pas` | 8 | 0 | 0 | 8 | 2 | 10 | 10 |
| `01-raccourcis-terminal` | 8 | 0 | 0 | 8 | 2 | 10 | 10 |
| `02-navigation-shell` | 5 | 4 | 0 | 9 | 5 | 6 | 10 |
| `05-outils-quotidiens` | 2 | 1 | 0 | 3 | 8 | 9 | 10 |
| `06-recherche-fichiers-contenu` | 2 | 9 | 1 | 12 | 8 | 1 | 9 |
| `07-lecture-preview` | 1 | 2 | 1 | 4 | 9 | 8 | 9 |
| `08-micro` | 5 | 1 | 0 | 6 | 5 | 9 | 10 |
| `09-markdown-glow` | 1 | 2 | 1 | 4 | 9 | 8 | 9 |
| `10-workflows-dotfiles` | 3 | 0 | 0 | 3 | 7 | 10 | 10 |
| `11-flux-shell` | 2 | 3 | 1 | 6 | 8 | 7 | 9 |
| `12-fichiers-permissions` | 5 | 1 | 0 | 6 | 5 | 9 | 10 |
| `13-processus-environnement` | 2 | 4 | 0 | 6 | 8 | 6 | 10 |
| **Total** | **44** | **27** | **4** | **75** | **76** | **93** | **116** |

Il manque au minimum **285 exercices**. La matrice réserve une compétence
distincte par exercice à concevoir. Si un sujet se révèle redondant, invérifiable
ou incompatible avec le moteur, trouver un autre objectif utile ou laisser le
manque ouvert. Une variante de commande ou un nom de fichier différent ne
constitue pas un nouvel apprentissage.

## Progression et format

| Niveau | Ce que l'exercice doit faire travailler | Exemple de progression |
| --- | --- | --- |
| Beginner | Effectuer un geste courant ou lire un résultat simple dans un contexte explicite. | Afficher la fin d'un fichier. |
| Intermediate | Assembler des gestes connus ou choisir des options selon des contraintes explicites. | Extraire une colonne en tenant compte du séparateur. |
| Advanced | Expliquer ou résoudre un cas où les hypothèses usuelles ne tiennent plus : portée, ambiguïté, données atypiques ou compatibilité. | Préserver des noms avec retours à la ligne tout au long d'un pipeline. |

Les prérequis guident la rédaction et l'ordre des lots ; ils ne bloquent pas un
exercice. Chaque énoncé fournit l'état fictif, les chemins et les préconditions
utiles, sans dévoiler la réponse. Le niveau dépend du raisonnement demandé,
indépendamment du nombre d'options ou du format. La saisie de commandes prime
lorsqu'elle permet d'évaluer le geste ; les questions de diagnostic peuvent
s'appuyer sur une sortie simulée ou un QCM discriminant.

## Matrice des compétences à ajouter

Le point-virgule sépare les compétences à transformer en exercices. Après
chacune, le code indique la [source officielle](#références-par-compétence) et
la section, l'option ou la commande à consulter.
Les extensions de périmètre sont des propositions éditoriales : elles ne
changent ni les identifiants, ni les familles, ni l'ordre actuel.

### `00-premiers-pas` — Premiers pas

**Couverture actuelle.** B : dossier courant, liste visible ou avec entrées cachées, déplacements relatif, parent et absolu, création simple, lecture d’un chemin avec espaces.

**Suite proposée.** Étendre aux repères et à l’interprétation d’une ligne de commande. Réserver les mutations, liens et droits à `12`, les flux à `11`.

| Niveau | Compétences à ajouter | Références |
| --- | --- | --- |
| Beginner (+2) | Revenir au dossier personnel (SH/tilde) ; Choisir l’aide locale d’une commande externe (FS/man). | SH · FS |
| Intermediate (+10) | Revenir au dossier précédent (SH/cd) ; Lister le contenu d’un autre dossier sans se déplacer (FS/ls) ; Créer une arborescence absente (FS/mkdir) ; Lire un chemin caché explicite (FS/cat) ; Distinguer l’entrée d’un dossier de son contenu (FS/ls) ; Développer un motif de noms sans ouvrir les fichiers (SH/glob) ; Sélectionner une famille de noms avec une classe de caractères (SH/glob) ; Protéger un nom contenant un caractère générique (SH/quoting) ; Distinguer le tilde développé du tilde littéral (SH/tilde) ; Lire la résolution d’un nom de commande (SH/resolution). | SH · FS |
| Advanced (+10) | Diagnostiquer une résolution via CDPATH (SH/cd) ; Distinguer dossier logique et physique derrière un lien (SH/cd) ; Neutraliser un glob dans une valeur et en conserver un autre actif (SH/quoting) ; Expliquer le comportement d’un motif sans correspondance selon le shell (BASH/glob+ZSH/glob) ; Prévoir la sélection des fichiers cachés par un motif (SH/glob) ; Résoudre une collision entre builtin, fonction et exécutable dans un contexte fourni (SH/resolution) ; Relancer une commande depuis un sous-dossier sans déplacer le shell appelant (SH/subshell) ; Résoudre un chemin relatif à un autre point de départ (SH/cd) ; Empêcher qu’un nom commençant par un tiret soit interprété comme option (FS/options) ; Distinguer la casse d’un nom de la sensibilité du système de fichiers (FS/ls). | SH · BASH · ZSH · FS |

### `01-raccourcis-terminal` — Édition au clavier

**Couverture actuelle.** B : début et fin de ligne, suppression de la ligne ou de sa fin, coupe avant le curseur, réinsertion, recherche arrière, rafraîchissement. ZLE en mode Emacs.

**Suite proposée.** Élargir à l’édition de ligne et à ses diagnostics. Distinguer ZLE et Readline ; réserver jobs/signaux à `13`. Les compétences diagnostiques se prêtent à des réponses textuelles ; l’entrée dédiée aux raccourcis ne montre que les captures. Tester le parcours mixte avant de l’utiliser.

| Niveau | Compétences à ajouter | Références |
| --- | --- | --- |
| Beginner (+2) | Reculer d’un caractère (ZLE/movement) ; Avancer d’un caractère (ZLE/movement). | ZLE |
| Intermediate (+10) | Reculer d’un mot (ZLE/movement) ; Avancer d’un mot (ZLE/movement) ; Supprimer le caractère sous le curseur (ZLE/edit) ; Supprimer le caractère avant le curseur (ZLE/edit) ; Transposer deux caractères (ZLE/edit) ; Couper jusqu’à la fin du mot suivant (ZLE/edit) ; Annuler une modification de ligne (ZLE/edit) ; Rappeler l’entrée précédente (ZLE/history) ; Revenir à l’entrée suivante (ZLE/history) ; Compléter un nom dans un contexte sans ambiguïté (ZLE/completion). | ZLE |
| Advanced (+10) | Diagnostiquer Ctrl+U après changement de shell (ZLE/edit+BASH/readline) ; Prévoir le sens d’un déplacement sur un chemin selon le découpage des mots (ZLE/movement) ; Distinguer la fin d’entrée sur ligne vide de l’effacement d’un caractère selon l’éditeur et le contexte (ZLE/edit+BASH/readline) ; Identifier le keymap actif depuis un diagnostic fourni (ZLE/keymaps) ; Choisir entre le début de la ligne et celui du buffer quand l’entrée occupe plusieurs lignes (ZLE/movement) ; Récupérer une autre coupe selon l’ordre du kill ring (ZLE/edit) ; Reprendre une recherche incrémentale dans l’autre sens (ZLE/history) ; Expliquer pourquoi une touche est interceptée avant ZLE (ZLE/keymaps) ; Distinguer les liaisons d’insertion des liaisons de commande en mode Vi (ZLE/keymaps) ; Diagnostiquer une complétion qui propose plusieurs candidats au lieu d’insérer un nom (ZLE/completion). | ZLE · BASH |

### `02-navigation-shell` — Navigation avec les outils CLI

**Couverture actuelle.** B : saut zoxide, liste fzf sur stdin, repérage fd, arborescence eza, ouverture Yazi. I : sélection fd et fzf de dossiers, profondeur eza, détails et fichiers cachés, liste zoxide.

**Suite proposée.** Se déplacer et s’orienter parmi les dossiers. La recherche par nom ou par contenu de fichiers relève de `06`.

| Niveau | Compétences à ajouter | Références |
| --- | --- | --- |
| Beginner (+5) | Lister uniquement les dossiers visibles (EZA/filter) ; Trier les entrées par modification (EZA/sort) ; Interroger zoxide sans se déplacer (ZOX/query) ; Ouvrir un dossier précis dans Yazi (YAZI/start) ; Revenir au dossier parent dans Yazi (YAZI/navigation). | EZA · ZOX · YAZI |
| Intermediate (+6) | Comparer tailles et unités en vue détaillée (EZA/long) ; Repérer les cibles des liens dans une vue (EZA/long) ; Afficher une branche de l’arbre en ignorant des sous-dossiers bruyants (EZA/filter) ; Ajouter un dossier connu à l’index zoxide (ZOX/add) ; Choisir parmi des homonymes via l’intégration interactive zoxide (ZOX/query) ; Retrouver le chemin d’un dossier depuis Yazi sans changer le dossier du shell (YAZI/navigation). | EZA · ZOX · YAZI |
| Advanced (+10) | Distinguer rang zoxide et date de modification du dossier (ZOX/rank) ; Écarter une entrée obsolète de l’index (ZOX/remove) ; Diagnostiquer pourquoi un saut échoue si le hook shell manque (ZOX/integration) ; Préserver le dossier courant quand une sélection est annulée (FZF/exit+SH/quoting) ; Distinguer le dossier du shell de celui atteint dans Yazi au retour (YAZI/shell) ; Prévoir si un arbre suit un lien symbolique vers un dossier hors de la racine demandée (EZA/filter) ; Inspecter une branche avec un nom contenant des métacaractères (EZA/filter+SH/quoting) ; Restreindre une sélection interactive à une racine et une profondeur utiles (FD/depth+FZF/input) ; Choisir une représentation sûre pour des chemins de dossiers avec retour à la ligne (FD/print0+FZF/input) ; Diagnostiquer un chemin relatif interprété depuis la mauvaise racine (FD/paths+SH/cd). | ZOX · FZF · SH · YAZI · EZA · FD |

### `05-outils-quotidiens` — Historique et interface Git

**Couverture actuelle.** B : recherche Atuin limitée, reconnaissance de Lazygit. I : recherche Atuin par code de sortie.

**Suite proposée.** Élargir à la consultation locale de l’historique et aux décisions dans les panneaux de Lazygit. Pour les opérations Git en ligne de commande, voir `10`.

| Niveau | Compétences à ajouter | Références |
| --- | --- | --- |
| Beginner (+8) | Afficher les statistiques locales Atuin (ATU/stats) ; Afficher la date d’exécution dans les résultats de recherche Atuin (ATU/search) ; Ouvrir Lazygit sur un dépôt donné (LG/start) ; Identifier la branche courante dans son panneau (LG/branches) ; Reconnaître le panneau des fichiers modifiés (LG/files) ; Examiner le diff d’un fichier sélectionné (LG/files) ; Consulter l’aide des touches du panneau actif (LG/global) ; Lire la liste des commits depuis le panneau correspondant (LG/commits). | ATU · LG |
| Intermediate (+9) | Limiter une recherche Atuin à un répertoire (ATU/search) ; Borner la recherche à une période (ATU/search) ; Rechercher avec l’interface interactive en ayant un motif initial (ATU/search) ; Exclure un code de sortie donné (ATU/search) ; Distinguer fichiers indexés et non indexés dans la vue (LG/files) ; Sélectionner un fragment à indexer sur une fixture (LG/files) ; Relire uniquement le diff indexé avant un commit à partir d’une capture simulée (LG/files) ; Retrouver une branche locale sans changer de branche (LG/branches) ; Examiner ce qu’un stash contient sans l’appliquer (LG/stash). | ATU · LG |
| Advanced (+10) | Croiser répertoire, période et statut dans l’historique (ATU/search) ; Expliquer une recherche vide malgré une commande exécutée dans un autre contexte (ATU/search) ; Distinguer historique du shell et historique enregistré par Atuin (ATU/integration) ; Distinguer la commande enregistrée de sa réexécution (ATU/search) ; Interpréter deux états d’un même fichier dans les panneaux Git (LG/files) ; Isoler le bon fragment lorsque d’autres modifications partagent le fichier (LG/files) ; Identifier le côté pertinent dans un conflit à partir d’un diff fourni (LG/files+GIT/diff) ; Examiner l’écart local/distant avec des références déjà disponibles (LG/branches+GIT/log) ; Diagnostiquer un HEAD détaché avant une opération (LG/branches+GIT/branch) ; Retrouver le bon stash parmi plusieurs sans altérer l’arbre (LG/stash). | ATU · LG · GIT |

### `06-recherche-fichiers-contenu` — Recherche et sélection

**Couverture actuelle.** B : nom via fd, rôle de rg. I : lignes et chemins, littéral, glob, type, extension, ignores, trois pipes de sélection. A : ouvrir avec bat la sélection issue d’une recherche de contenu.

**Suite proposée.** Deux pipes existants `rg-fzf-todo` et `rg-fzf-docs` se recouvrent. Garder leurs clés ; ne pas ajouter une troisième version du même geste. `02` possède la sélection de dossiers, `09` le rendu Markdown.

| Niveau | Compétences à ajouter | Références |
| --- | --- | --- |
| Beginner (+8) | Chercher un motif dans un fichier nommé (RG/basics) ; Limiter la recherche de contenu à un dossier donné (RG/basics) ; Ignorer la casse d’un motif (RG/case) ; Compter les lignes correspondantes par fichier (RG/count) ; Lister les lignes sans correspondance (RG/invert) ; Reconnaître la ligne de contexte autour d’un résultat (RG/context) ; Choisir un motif ancré sur le début du nom (FD/pattern) ; Filtrer une liste fournie sans modifier ses éléments (FZF/filter). | RG · FD · FZF |
| Intermediate (+1) | Faire primer une exclusion explicite sur une sélection de fichiers (RG/glob+FD/exclude). | RG · FD |
| Advanced (+9) | Inclure seulement un sous-ensemble de fichiers ignorés (RG/ignore) ; Rendre visibles des résultats dans un dépôt imbriqué selon les règles d’ignore (RG/ignore) ; Préserver des chemins arbitraires entre producteur et sélecteur avec des séparateurs NUL (FD/print0+FZF/input) ; Empêcher l’action aval sur sélection vide (FZF/exit+SH/quoting) ; Traiter plusieurs sélections sans les scinder sur les espaces (FZF/multi+SH/quoting) ; Diagnostiquer une recherche qui ignore un fichier traité comme binaire (RG/binary) ; Interpréter le statut d’échec distinct d’une absence de résultat (RG/exit) ; Vérifier une expression qui doit couvrir plusieurs lignes (RG/multiline) ; Exploiter un format structuré lorsque chemins et lignes contiennent des délimiteurs ambigus (RG/json). | RG · FD · FZF · SH |

### `07-lecture-preview` — Lecture et inspection de texte

**Couverture actuelle.** B : lire un fichier avec bat. I : numéros de ligne, plage. A : aperçu JSON sur stdin, couleur et pager explicités.

**Suite proposée.** Élargir de la simple lecture à l’inspection d’un texte. Garder bat dans les exercices du chapitre : la famille actuelle requiert au moins un module tiers. `11` possède la mécanique des flux ; `09`, la mise en forme Markdown.

| Niveau | Compétences à ajouter | Références |
| --- | --- | --- |
| Beginner (+9) | Lire plusieurs fichiers dans l’ordre fourni (BAT/usage) ; Lire depuis l’entrée standard (BAT/usage) ; Éviter le pager pour une lecture continue (BAT/paging) ; Afficher sans décor (BAT/style) ; Révéler les caractères invisibles (BAT/nonprint) ; Indiquer la syntaxe d’un fichier sans extension (BAT/syntax) ; Lister les syntaxes reconnues (BAT/syntax) ; Lister les thèmes installés (BAT/themes) ; Ouvrir un fichier dont le nom débute par un tiret (BAT/usage+FS/options). | BAT · FS |
| Intermediate (+8) | Associer durablement un motif de nom à un langage via une configuration fournie (BAT/syntax) ; Indiquer des lignes à mettre en évidence (BAT/highlight) ; Comparer largeur du terminal et gestion des retours visuels (BAT/wrap) ; Composer un extrait de fin de fichier en annonçant la perte de numérotation source (BAT/usage+FS/tail) ; Distinguer une tabulation affichée d’espaces réels (BAT/nonprint) ; Afficher un extrait avec son chemin (BAT/style) ; Comparer style réduit et style complet (BAT/style) ; Examiner les résultats de recherche tout en conservant leur provenance (BAT/integration+RG/basics). | BAT · FS · RG |
| Advanced (+9) | Produire une sortie exploitable par un outil aval sans codes de couleur (BAT/color+SH/pipelines) ; Diagnostiquer un pager imposé par l’environnement (BAT/paging) ; Examiner un fichier sans extension dont le contenu trompe la détection (BAT/syntax) ; Conserver les numéros fournis par la recherche sans renuméroter ses lignes à l’affichage (BAT/integration+RG/basics) ; Localiser un caractère de contrôle dans un extrait fourni (BAT/nonprint) ; Diagnostiquer un décalage entre largeur visuelle et nombre d’octets (BAT/wrap) ; Préserver des chemins provenant d’une recherche à délimiteurs NUL (BAT/integration+RG/print0) ; Distinguer inspection d’un fichier non textuel et tentative de rendu (BAT/syntax) ; Inspecter depuis une ligne donnée jusqu’à la fin d’un flux de longueur inconnue (BAT/range+SH/pipelines). | BAT · SH · RG |

### `08-micro` — Édition de texte

**Couverture actuelle.** B : ouvrir, enregistrer, chercher, quitter, annuler. I : bascule des numéros de ligne.

**Suite proposée.** Élargir à l’édition de plusieurs buffers, aux remplacements et à la préservation du format. Les commandes de l’invite Micro et les touches décrites en texte ne sont pas des captures de raccourcis.

| Niveau | Compétences à ajouter | Références |
| --- | --- | --- |
| Beginner (+5) | Refaire une modification annulée (MIC/keys) ; Sélectionner toute la zone d’édition (MIC/keys) ; Couper une sélection (MIC/keys) ; Copier sans effacer (MIC/keys) ; Ouvrir l’aide de l’éditeur (MIC/commands). | MIC |
| Intermediate (+9) | Aller à une ligne donnée (MIC/commands) ; Retrouver l’occurrence suivante (MIC/keys) ; Remplacer une occurrence (MIC/commands) ; Remplacer dans tout le buffer (MIC/commands) ; Enregistrer sous un autre nom (MIC/commands) ; Ouvrir une seconde vue (MIC/commands) ; Passer au buffer voisin (MIC/keys) ; Changer la largeur des tabulations (MIC/options) ; Régler le retour visuel à la ligne (MIC/options). | MIC |
| Advanced (+10) | Restreindre un remplacement aux débuts de lignes via une expression (MIC/commands) ; Traiter un point littéral dans une recherche par motif (MIC/commands) ; Distinguer indentation affichée et tabulations écrites (MIC/options) ; Préserver les fins de ligne d’un fichier partagé (MIC/options) ; Diagnostiquer pourquoi une syntaxe est mal reconnue sans renommer le fichier (MIC/options) ; Récupérer le contenu d’un buffer lorsque la destination n’est pas accessible (MIC/commands) ; Différencier réglage local au buffer et réglage global dans un contexte fourni (MIC/options) ; Repérer quelle vue partage le même buffer lors d’une modification (MIC/commands) ; Préserver le texte en quittant une vue contenant des modifications (MIC/commands) ; Diagnostiquer une touche indisponible ou remappée à partir de la configuration fournie (MIC/keys+MIC/options). | MIC |

### `09-markdown-glow` — Lecture de Markdown

**Couverture actuelle.** B : ouvrir un fichier avec Glow. I : pager, largeur. A : sélectionner un fichier Markdown via fd/fzf puis le rendre, avec protection des espaces.

**Suite proposée.** Élargir à la relecture locale du Markdown : repérer et réparer un rendu inattendu à partir de courts extraits, sans ajouter de moteur d’édition. Associer les diagnostics aux comportements réellement observés de Glow ; ne présumer ni la prise en charge complète de GFM ni un export HTML.

| Niveau | Compétences à ajouter | Références |
| --- | --- | --- |
| Beginner (+9) | Lire un extrait sur stdin (GLOW/stdin) ; Consulter un fichier depuis l’interface locale (GLOW/tui) ; Choisir un style de rendu intégré (GLOW/styles) ; Distinguer un titre d’une ligne de texte (MD/headings) ; Reconnaître un lien dont le libellé diffère de l’adresse (MD/links) ; Distinguer texte et code en ligne (MD/code) ; Lire une liste ordonnée (MD/lists) ; Repérer un bloc de citation (MD/quotes) ; Reconnaître un bloc de code clôturé (MD/fences). | GLOW · MD |
| Intermediate (+8) | Distinguer l’entrée standard d’un argument interprété comme chemin (GLOW/stdin+SH/quoting) ; Identifier pourquoi des lignes indentées ont été rendues comme code (MD/code) ; Reconnaître pourquoi deux paragraphes ont fusionné (MD/paragraphs) ; Préserver l’imbrication voulue d’une liste (MD/lists) ; Rendre visible un fragment Markdown comme code (MD/fences) ; Distinguer une URL d’un chemin local qui contient des espaces (MD/links+SH/quoting) ; Repérer une table rendue comme texte dans un scénario GFM vérifié (MD/tables+GLOW/cli) ; Relire en deux temps la source et le rendu en mobilisant bat et Glow (BAT/usage+GLOW/cli). | GLOW · SH · MD · BAT |
| Advanced (+9) | Diagnostiquer une clôture de code absorbant la suite du document (MD/fences) ; Protéger une clôture citée à l’intérieur d’un bloc (MD/fences) ; Conserver le bon niveau d’imbrication d’un bloc au sein d’une liste (MD/lists) ; Échapper un séparateur dans une cellule GFM vérifiée (MD/tables+GLOW/cli) ; Diagnostiquer un lien relatif selon l’emplacement de la source, sans suivre le lien (MD/links+GLOW/cli) ; Distinguer le retour visuel de la coupure Markdown explicite (MD/breaks+GLOW/wrap) ; Diagnostiquer un style peu lisible selon le profil terminal (GLOW/styles) ; Prévisualiser une portion de document sans couper la structure qui l’entoure (MD/blocks+GLOW/stdin) ; Distinguer la représentation du fichier de ce que le pager peut masquer (GLOW/pager+BAT/usage). | MD · GLOW · BAT |

### `10-workflows-dotfiles` — Git essentiel

**Couverture actuelle.** B : état court, résumé des changements et liste des fichiers modifiés.

**Suite proposée.** Élargir à l’index, à l’historique et aux branches. Garder l’identifiant historique ; aucun accès aux dotfiles ou à un dépôt utilisateur. `05` traite la lecture des mêmes états dans Lazygit.

| Niveau | Compétences à ajouter | Références |
| --- | --- | --- |
| Beginner (+7) | Voir la branche active (GIT/branch) ; Lire le dernier commit (GIT/show) ; Consulter l’historique condensé (GIT/log) ; Afficher le diff complet d’un seul fichier (GIT/diff) ; Indexer un chemin déterminé (GIT/add) ; Créer une branche sans y basculer (GIT/branch) ; Lire le contenu d’un fichier à un commit donné (GIT/show). | GIT |
| Intermediate (+10) | Distinguer le diff de l’index de celui du dossier de travail (GIT/diff) ; Désindexer un fichier en conservant le travail (GIT/restore) ; Indexer une partie d’un fichier (GIT/add) ; Trouver la règle responsable d’un fichier ignoré (GIT/check-ignore) ; Suivre l’historique d’un fichier renommé (GIT/log) ; Lire les commits qui n’existent que sur une branche (GIT/log) ; Comparer les deux côtés d’une divergence (GIT/log) ; Examiner les fichiers d’un commit précis (GIT/show) ; Rechercher dans une révision plutôt que dans le dossier courant (GIT/grep) ; Comparer le contenu local et une référence distante déjà connue (GIT/diff+PROGIT/remotes). | GIT · PROGIT |
| Advanced (+10) | Expliquer les deux statuts du même fichier (GIT/status) ; Comparer branches depuis leur ancêtre commun (GIT/merge-base) ; Retrouver l’identifiant d’un commit avec un reflog simulé (GIT/reflog) ; Distinguer HEAD détaché et branche sans commit (GIT/status) ; Diagnostiquer une branche sans suivi distant (GIT/branch) ; Retrouver la bonne entrée d’un stash sans l’appliquer (GIT/stash) ; Examiner les différentes versions d’un fichier non fusionné (GIT/revisions) ; Distinguer l’effet d’une exclusion Git sur un fichier déjà suivi (GIT/check-ignore) ; Identifier le premier changement suspect par bissection sur un graphe fourni (GIT/bisect) ; Séparer la recherche d’un texte supprimé de celle d’un motif dans un diff (GIT/diff). | GIT |

### `11-flux-shell` — Flux et composition

**Couverture actuelle.** B : écrire puis ajouter en fin de fichier. I : compter sur stdin, filtrer et compter, rediriger stderr. A : ordre de fusion des descripteurs.

**Suite proposée.** Couvrir les flux et la composition portable. `07` traite la présentation d’un texte ; `13` les processus et l’environnement.

| Niveau | Compétences à ajouter | Références |
| --- | --- | --- |
| Beginner (+8) | Lire stdin depuis un fichier (SH/redirection) ; Relier un producteur à un lecteur (SH/pipelines) ; Afficher les premières lignes (FS/head) ; Afficher les dernières lignes (FS/tail) ; Compter les mots (FS/wc) ; Trier des lignes (FS/sort) ; Réduire des doublons adjacents (FS/uniq) ; Afficher tout en conservant une copie (FS/tee). | SH · FS |
| Intermediate (+7) | Trier puis dédupliquer des entrées non adjacentes (FS/sort+FS/uniq) ; Compter les occurrences par ligne distincte (FS/uniq) ; Trier une colonne numérique (FS/sort) ; Extraire une colonne avec un délimiteur annoncé (FS/cut) ; Ajouter une copie avec tee sans écraser (FS/tee) ; Séparer les sorties normale et d’erreur dans deux destinations (SH/redirection) ; Fournir un bloc littéral via un document incorporé (SH/heredoc). | FS · SH |
| Advanced (+9) | Identifier le maillon en échec selon les règles du shell (BASH/pipelines+ZSH/options) ; Distinguer le statut d’un pipe de ses statuts individuels (BASH/variables+ZSH/parameters) ; Préserver les octets d’un flux mêlant noms et retours à la ligne (FS/options+SH/pipelines) ; Éviter la troncature quand l’entrée et la destination désignent le même fichier (SH/redirection) ; Distinguer redirection persistante et redirection locale à une commande (SH/redirection) ; Comparer redirection d’un groupe et redirection d’un seul maillon (SH/redirection) ; Expliquer l’expansion dans un document incorporé (SH/heredoc) ; Diagnostiquer un tri incohérent sous deux locales (FS/sort) ; Distinguer fins de ligne et enregistrements lors d’un comptage (FS/wc). | BASH · ZSH · FS · SH |

### `12-fichiers-permissions` — Fichiers et permissions

**Couverture actuelle.** B : copie avec espaces, renommage, suppression simple, retrait d’un dossier vide, lecture des droits. I : rendre exécutable pour le propriétaire.

**Suite proposée.** Couvrir la manipulation sûre de fichiers, les liens et les droits. Les contextes potentiellement destructifs restent fictifs et circonscrits.

| Niveau | Compétences à ajouter | Références |
| --- | --- | --- |
| Beginner (+5) | Créer un fichier vide absent (FS/touch) ; Créer une copie de sauvegarde sous un nom distinct (FS/cp) ; Identifier le propriétaire dans une liste détaillée (FS/ls) ; Retirer au groupe le droit d’écriture (FS/chmod) ; Identifier un lien symbolique dans une liste détaillée (FS/ls). | FS |
| Intermediate (+9) | Copier un dossier avec ses descendants (FS/cp) ; Confirmer un remplacement avant copie (FS/cp) ; Déplacer plusieurs sources vers une destination unique (FS/mv) ; Préserver les dates lors d’une copie (FS/cp) ; Créer un lien symbolique relatif (FS/ln) ; Identifier la cible d’un lien (FS/readlink) ; Donner des droits distincts à plusieurs classes sans réécrire les autres (FS/chmod) ; Interpréter une notation octale à trois chiffres (FS/chmod) ; Diagnostiquer un droit de traversée manquant (FS/permissions). | FS |
| Advanced (+10) | Distinguer suppression de l’entrée et suppression de la cible d’un lien (FS/rm+FS/ln) ; Déduire les permissions initiales depuis un mode et un umask fournis (FS/umask) ; Distinguer droits de suppression et droits d’écriture sur le fichier (FS/rm+FS/permissions) ; Prévoir l’effet du sticky bit dans un dossier partagé (FS/permissions) ; Interpréter l’héritage de groupe d’un dossier configuré avec setgid (FS/permissions) ; Déduire l’effet d’un déplacement de lien relatif (FS/ln+FS/mv) ; Distinguer inode partagé et lien symbolique face à un remplacement (FS/ln) ; Appliquer conditionnellement la traversée sans rendre tous les fichiers exécutables (FS/chmod) ; Diagnostiquer une destination de copie ambiguë qui existe déjà (FS/cp) ; Choisir une méthode de relevé des droits selon BSD ou GNU (FS/stat). | FS |

### `13-processus-environnement` — Processus et environnement

**Couverture actuelle.** B : consulter une variable et affecter une valeur shell. I : examiner un PID, demander l’arrêt par TERM, exporter, modifier l’environnement d’une seule commande.

**Suite proposée.** Étendre à la portée, à l’observation et au contrôle des jobs. Les processus et résultats sont fournis ou simulés ; aucune opération sur les processus de la machine.

| Niveau | Compétences à ajouter | Références |
| --- | --- | --- |
| Beginner (+8) | Lister les processus de la session (PROC/ps) ; Reconnaître le PID dans une sortie fournie (PROC/ps) ; Consulter l’environnement exporté (PROC/env) ; Retirer une variable du shell (SH/parameters) ; Afficher les jobs du shell (BASH/jobs+ZSH/jobs) ; Démarrer une commande fictive en arrière-plan (BASH/jobs+ZSH/jobs) ; Ramener un job connu au premier plan (BASH/jobs+ZSH/jobs) ; Distinguer variable vide et variable absente (SH/parameters). | PROC · SH · BASH · ZSH |
| Intermediate (+6) | Afficher PID et PPID dans des colonnes choisies (PROC/ps) ; Récupérer le PID du dernier lancement en arrière-plan (SH/parameters) ; Attendre la fin d’un processus enfant précis (PROC/wait) ; Reprendre un job suspendu en arrière-plan (BASH/jobs+ZSH/jobs) ; Lire la signification d’un statut renvoyé par wait dans le cas fourni (PROC/wait) ; Lancer une commande dans un environnement vidé en lui donnant les seules valeurs utiles (PROC/env). | PROC · SH · BASH · ZSH |
| Advanced (+10) | Distinguer un numéro de job, un PID et un groupe de processus (BASH/jobs+ZSH/jobs) ; Choisir un signal en tenant compte d’un nettoyage à effectuer (PROC/kill) ; Expliquer une suspension de job qui lit son terminal en arrière-plan (BASH/jobs+ZSH/jobs) ; Diagnostiquer une commande absente après modification de PATH (SH/resolution) ; Distinguer l’instant où une variable est exportée de celui où un enfant déjà lancé a reçu son environnement (PROC/env+SH/parameters) ; Expliquer pourquoi un sous-shell ne modifie pas le parent (SH/subshell) ; Distinguer les effets d’un lancement et d’un chargement dans le shell courant (SH/resolution) ; Distinguer un répertoire courant hérité d’une variable d’environnement (SH/subshell) ; Interpréter un processus qui existe encore mais ne peut plus exécuter d’instructions (PROC/ps) ; Déduire la portée d’un gestionnaire de signal dans un scénario shell fourni (SH/trap). | BASH · ZSH · PROC · SH |

## Références par compétence

Ces liens pointent vers les manuels ou guides officiels maintenus. Les
sections indiquées servent à choisir puis à vérifier les compétences de la
matrice. À la rédaction, relever la version effectivement documentée et
l'environnement du scénario. Pour une option récente ou un raccourci, citer
la version ou le tag dans le rapport du lot si le comportement en dépend.

| Code | Références et sections à utiliser |
| --- | --- |
| SH | [POSIX.1-2024, corrigendum IEEE](https://standards.ieee.org/ieee/1003.1-2024_Cor_1/11920/) ; [The Open Group, Shell Command Language](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html) : Quoting, Word Expansions (tilde, paramètres et motifs), Redirection (Here-Document), Pipelines, Compound Commands (sous-shell) et Command Search and Execution. Pour les builtins, voir `cd` et `trap` dans Utilities. |
| FS | [POSIX, Utilities](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/contents.html) : entrées `man`, `cd`, `ls`, `mkdir`, `cp`, `mv`, `rm`, `ln`, `chmod`, `umask`, `sort` et `tee` notamment. [GNU Coreutils](https://www.gnu.org/software/coreutils/manual/html_node/) : Common options, Output of parts of files, Sorting, Directory listing, Basic operations et File permissions. Pour les écarts macOS, vérifier les pages `man` de la version cible (`ls(1)`, `cp(1)`, `chmod(1)`, `stat(1)`). |
| PROC | POSIX Utilities : [`ps`](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/ps.html), [`kill`](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/kill.html), [`wait`](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/wait.html), [`env`](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/env.html) ; Bash et ZSH pour le contrôle des jobs et les extensions. |
| BASH | [GNU Bash Reference Manual](https://www.gnu.org/software/bash/manual/) : Shell Expansions, Pipelines, Job Control, Bourne Shell Builtins, Bash Variables et Command Line Editing / Readline (Killing and Yanking, Readline Init File). |
| ZSH | [ZSH Manual](https://zsh.sourceforge.io/Doc/Release/) : Expansion (Filename Generation), Parameters, Shell Builtin Commands, Jobs & Signals et Options. |
| ZLE | [Zsh Line Editor](https://zsh.sourceforge.io/Doc/Release/Zsh-Line-Editor.html) : Keymaps et Standard Widgets (Movement, History Control, Modifying Text, Completion). |
| GIT | [Git Reference](https://git-scm.com/docs) : `git-status`, `git-diff`, `git-add`, `git-restore`, `git-show`, `git-log`, `gitrevisions`, `git-branch`, `git-merge-base`, `git-check-ignore`, `git-grep`, `git-stash`, `git-reflog` et `git-bisect`. |
| PROGIT | [Pro Git](https://git-scm.com/book/en/v2) : Git Basics (Recording Changes, Viewing the Commit History), Git Branching (Remote Branches), Git Tools (Revision Selection, Interactive Staging, Stashing and Cleaning, Debugging with Git). |
| FZF | [fzf](https://github.com/junegunn/fzf) : Usage, Search Syntax, Environment Variables ; [manuel](https://github.com/junegunn/fzf/blob/master/man/man1/fzf.1) : options d'entrée/sortie, filtre, sélection multiple et statut de sortie. |
| FD | [fd](https://github.com/sharkdp/fd) : Usage (hidden and ignored files, recherche par extension, exclusion et exécution sur résultats) ; [manuel](https://github.com/sharkdp/fd/blob/master/doc/fd.1) : type, profondeur, chemins et séparateur NUL. |
| RG | [ripgrep, GUIDE](https://github.com/BurntSushi/ripgrep/blob/master/GUIDE.md) : Basics, Automatic filtering, Configuration file et Regex syntax ; [source des options](https://github.com/BurntSushi/ripgrep/blob/master/crates/core/flags/defs.rs) et [manuel de la version](https://github.com/BurntSushi/ripgrep/blob/master/crates/core/flags/doc/template.rg.1) : sorties, recherche multiligne et statuts. |
| EZA | [eza, manuel](https://github.com/eza-community/eza/blob/main/man/eza.1.md) : Filtering and sorting options, Long view options et Display options. |
| ZOX | [zoxide](https://github.com/ajeetdsouza/zoxide) : Examples, Shell Integration et Configuration ; sous-commandes et calcul du rang dans [le manuel de la version](https://github.com/ajeetdsouza/zoxide/tree/main/man). |
| YAZI | [Yazi, Quick Start](https://yazi-rs.github.io/docs/quick-start/) : navigation, chemins, touches par défaut et adaptation du dossier courant du shell. |
| ATU | [Atuin, documentation](https://docs.atuin.sh/) : [Search](https://docs.atuin.sh/latest/reference/search/) (période, dossier, codes de sortie), [Stats](https://docs.atuin.sh/latest/reference/stats/) et [Installation](https://docs.atuin.sh/latest/guide/installation/) (intégration shell). |
| LG | [Lazygit, Keybindings](https://github.com/jesseduffield/lazygit/blob/master/docs/keybindings/Keybindings_en.md) : Global, Files, Branches, Commits et Stash ; [README](https://github.com/jesseduffield/lazygit) : Usage et ouverture d'un dépôt. |
| BAT | [bat](https://github.com/sharkdp/bat) : Usage, Integration with other tools, Customization (syntaxes, thèmes, style, pager). Vérifier les options d'affichage dans le [manuel](https://github.com/sharkdp/bat/blob/master/assets/manual/bat.1.in). |
| MIC | Micro : [Default keys](https://github.com/micro-editor/micro/blob/master/runtime/help/defaultkeys.md), [Commands](https://github.com/micro-editor/micro/blob/master/runtime/help/commands.md) (aide, navigation, recherche/remplacement, buffers), [Options](https://github.com/micro-editor/micro/blob/master/runtime/help/options.md) (syntaxe, indentation, fins de ligne et portée des réglages). |
| GLOW | [Glow](https://github.com/charmbracelet/glow) : Usage (TUI, CLI, entrée standard), Word Wrapping et Styles. Vérifier séparément les fonctions du moteur de rendu sur la version retenue. |
| MD | [CommonMark](https://spec.commonmark.org/current/) : Blocks and Inlines, listes, clôtures de code et liens. [GFM](https://github.github.com/gfm/) : extension Tables. Les spécifications expliquent la syntaxe ; elles ne garantissent pas à elles seules le rendu de Glow. |

Pour l'ordre pédagogique, [MIT Missing Semester](https://missing.csail.mit.edu/)
donne des parcours séparés sur [le shell](https://missing.csail.mit.edu/2020/course-shell/),
[les outils du shell](https://missing.csail.mit.edu/2020/shell-tools/) et
[le contrôle de version](https://missing.csail.mit.edu/2020/version-control/).
[Software Carpentry, The Unix Shell](https://swcarpentry.github.io/shell-novice/)
enchaîne Navigation, Working with Files and Directories, puis
[Pipes and Filters](https://swcarpentry.github.io/shell-novice/04-pipefilter.html).
Ces cours guident la progression ; les situations et réponses restent
originales et la syntaxe est vérifiée dans les manuels.

### Différences à vérifier avant rédaction

- **Unix et plateformes.** Partir de POSIX et comparer au manuel macOS/BSD
  et à GNU pour les variantes. Indiquer la plateforme dans l'énoncé si l'option,
  le format ou le comportement diffère (par exemple, le relevé avec `stat`
  ou la sélection et les colonnes de `ps`). Vérifier la locale, la casse et
  les noms atypiques lorsqu'ils changent l'issue.
- **Bash et ZSH.** Nommer le shell et le mode d'édition. Vérifier le glob sans
  correspondance, le découpage des mots, les statuts des pipelines, les jobs et
  le keymap séparément. Pour une capture, confirmer que la touche est transmise
  et reconnue dans le terminal cible ; éviter les touches réservées à la sortie
  ou à l'interruption. Un diagnostic textuel ne vaut pas validation d'une capture.
- **Outils tiers.** Déclarer tous les modules réellement requis et vérifier
  options et défauts pour la version citée. Ne supposer ni alias, ni configuration
  personnelle, ni intégration shell implicite, ni service distant. Un diagnostic
  Markdown prévu dans `09` doit avoir un rapport concret avec le rendu Glow ;
  sinon il n'a pas sa place dans ce chapitre.

## Enrichissement par lots

Chaque vague peut donner plusieurs branches courtes, typiquement de 15 à 25
exercices, à adapter au périmètre et à la revue. Prendre un petit ensemble de
compétences qui se complètent et traiter les trois niveaux quand les prérequis
sont prêts. Les nombres sont les écarts du relevé initial, pas un nombre
d'exercices à livrer d'un coup.

| Ordre | Vague | Chapitres | Écart | Dépendance ou preuve spécifique |
| --- | --- | --- | ---: | --- |
| 1 | Repères, fichiers et droits | `00`, `12` | 46 | Stabiliser chemins et guillemets ; expliciter les différences de droits et de liens. |
| 2 | Flux, processus et environnement | `11`, `13` | 48 | Partir des repères de la vague 1 ; vérifier descripteurs, jobs et portabilité avant les cas Advanced. |
| 3 | Édition de la ligne de commande | `01` | 22 | Vérifier les touches en Ghostty et, avant les diagnostics, les parcours dédiés et mixtes. |
| 4 | Navigation, recherche et sélection | `02`, `06` | 39 | S'appuyer sur les chemins et les flux ; vérifier les modules combinés, l'annulation et les chemins atypiques. |
| 5 | Lecture, édition et aperçu | `07`, `08`, `09` | 76 | Commencer par l'inspection et l'édition locales ; vérifier ensuite les compositions et les diagnostics de rendu. |
| 6 | Git, historique et interfaces | `10`, `05` | 54 | Commencer par les états Git avant les panneaux Lazygit ; couvrir Atuin local sans exiger l'intégration shell sauf scénario explicite. |
| | **Total** | | **285** | |

Une frontière partagée a un responsable pédagogique : chemins et résolution
en `00`, effets sur le système de fichiers en `12`, flux en `11`, portée et
processus en `13`. Appliquer la même règle à la navigation/recherche
(`02`/`06`), à l'inspection/rendu (`07`/`09`) et aux commandes/interfaces Git
(`10`/`05`). Une deuxième mise en situation n'est utile que si elle évalue une
décision différente. Avant le premier lot qui élargit un chapitre, préciser les
limites dans son rapport ; si une compétence exige de modifier le moteur ou la
structure historique, la remplacer ou laisser la case ouverte.

### Acceptation d'un lot

1. Relier chaque ajout à une compétence et à une section officielle précise.
   Relire la consigne pour contrôler l'autonomie, l'absence de réponse dévoilée,
   la justesse du niveau et l'utilité de l'explication.
2. Vérifier les réponses principales, les variantes équivalentes et des erreurs
   plausibles avec le correcteur existant. Si une vérification réelle s'impose,
   employer une fixture isolée ; l'application ne doit jamais exécuter les
   réponses saisies.
3. Contrôler les clés, niveaux, modules et filtres exacts avec des tests ciblés,
   puis la suite Go. Essayer les parcours guidés ou mélangés concernés, la
   sélection sans module tiers et les combinaisons de modules pertinentes.
   Contrôler la projection des acquis avec des fixtures JSON v1 si le lot touche
   à l'éligibilité. Conserver clés, format de progression et préférences.
4. Noter le décompte avant/après, les limites et les résultats de recette dans
   le rapport du lot. Distinguer l'essai sous pseudo-terminal de la recette
   humaine dans Ghostty, et signaler les plateformes et versions non testées.

La [spécification](SPEC.md) fixe les garanties communes. Ce référentiel ne
change pas le moteur, les modes de séance, le schéma des données ni
[P3 — Internationalisation FR/EN](../TODO.md#p3--internationalisation-fren).
