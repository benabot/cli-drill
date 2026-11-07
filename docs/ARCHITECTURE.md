# CLI-Drill V2 — Architecture cible

**Cible décidée le 7 novembre 2026 ; non encore implémentée.**
La branche documentaire ne change aucun package Go ni aucune dépendance.
L'[état courant](../PROJECT_STATE.md) décrit le code disponible ; la
[spécification](SPEC.md) définit les comportements V2 et les garanties de données.

## Flux cible

```text
Chapitres YAML intégrés -> chargement et validation
                                |
                         entraînement CLI / TUI
                                |
                    correction -> progression JSON locale
```

Le module reste à la racine, le point d'entrée dans `cmd/cli-drill/`,
l'embarquement dans `data/defaults.go` et le corpus dans `data/chapters/`.
Les données éditables sont des sources pédagogiques maintenues dans le projet ;
aucun dossier utilisateur ne remplace automatiquement le corpus intégré.

## Packages conservés

| Package | Responsabilité cible |
| --- | --- |
| `internal/app` | Commandes Cobra conservées, choix des chapitres et orchestration TUI / entraînement dédié. |
| `internal/chapter` | Modèle pédagogique, chargement du corpus intégré et validation YAML. |
| `internal/exercise` | Types d'exercices, correction pure et normalisation adaptée au contexte. |
| `internal/progress` | Compteurs, compatibilité JSON et persistance sûre. |
| `internal/tui` | Navigation, saisie, filtres, feedback et rendu Bubble Tea. |
| `internal/xdg` | Résolution du chemin actuel de progression et respect de `XDG_DATA_HOME`. |

Conserver ces frontières. Aucun package de services, interface de stockage,
registre de moteurs ou framework de session supplémentaire n'est nécessaire.
Répartir éventuellement les gros fichiers dans leur package existant lorsque
cela aide un changement ciblé ; ne pas engager une réorganisation générale.

## Packages et responsabilités à retirer

| Élément actuel | Retrait prévu et condition |
| --- | --- |
| `internal/detect` | Scanner et filtrage de sources dotfiles ; retirer leurs consommateurs CLI/TUI. |
| `internal/shell/zsh` | Parsing statique des alias et fonctions ; aucun remplacement par un autre shell. |
| `internal/markdown` | Extraction depuis les documents utilisateur ; aucun parser Markdown requis par le corpus YAML. |
| `internal/config` | Chargement et création du TOML dotfiles ; ne pas supprimer les fichiers TOML existants. |
| `internal/catalog` | Annuaire, déduplication et recherche ; découpler auparavant les types pédagogiques qui le référencent. |
| Génération dans `chapter` | Retirer `GenerateFromCatalog`, écriture de leçons générées et conversion `ToCatalog`. |
| Chemins de configuration dans `xdg` | Retirer les helpers devenus inutilisés, sans nettoyage des dossiers utilisateur. |

`chapter.Item` dépend actuellement de `catalog.EntryType` et `catalog.Source`.
Les métadonnées pédagogiques encore utiles doivent devenir des données locales
simples de `chapter`, sans nouvelle abstraction de catalogue. Le découplage
ne doit pas imposer une conversion des YAML utilisateur conservés.

## CLI, TUI et vérifications du terminal

La CLI cible conserve `chapters`, `train [chapter-id]`, `stats`, `reset --yes`
et l'option `--progress`. Retirer `init`, `scan`, `generate`, `directory`,
`search`, `show` et `--config`, ainsi que les écrans et états associés.

Réutiliser le navigateur de chapitres pour la vue Keyboard shortcuts et les
composants de progression pour My progress. Conserver les moteurs textuel et
de capture dédié : leurs entrées et leur rendu ont des contraintes différentes.
Le retour à l'accueil après capture reste orchestré par `app`.

Les contrôles techniques du terminal restent nécessaires : entrée/sortie TTY,
dimensions, rendu couleur et sans couleur, activation et restauration du mode
brut. L'abandon de l'introspection ne concerne que les dotfiles et la détection
des outils installés. Ne pas supprimer `x/term`, `TERM` ou `NO_COLOR` au titre
de cette simplification ; conserver leurs usages techniques appropriés.

Avant d'activer la TUI ou la capture brute, vérifier les capacités nécessaires
et traiter leurs erreurs. Le premier chapitre choisi par défaut est textuel.
Si la TUI ne peut fonctionner, guider vers un entraînement CLI textuel sans
couleur ; si la capture n'est pas possible, expliquer la limite sans la valider
artificiellement. Aucun support universel d'ANSI ou de raccourcis macOS n'est promis.

## Règles partagées et données

Partager la correction et l'enregistrement des tentatives entre les interfaces :
réponse soumise, correction, mise à jour puis sauvegarde. Le feedback et la
navigation restent propres à chaque interface. Une fonction concrète ou une
fonction de rappel suffit si un point de partage est nécessaire ; ne pas
introduire un bus d'événements ni fusionner les boucles d'interface.

La progression reste en JSON version 1, avec ses champs, chemins et clés actuels.
La [politique de compatibilité](SPEC.md#7-progression-et-donnees-heritees) impose
la conservation intégrale des données héritées. Aucun code de migration ne doit
purger, déplacer ou réécrire automatiquement les JSON, YAML ou TOML utilisateur.
L'entraînement enrichit la progression sans retirer les anciens résultats ;
les indicateurs du corpus actif filtrent les identifiants dans leur projection.

Prévoir des sauvegardes atomiques et des erreurs visibles. Ne pas remplacer une
progression invalide ou non prise en charge par une valeur par défaut écrite.
L'écriture atomique ne résout pas les mises à jour de plusieurs processus :
le nouveau MVP vise un seul écrivain par fichier, sans ajout immédiat de base
de données ou de mécanisme de synchronisation.

## Risques de migration

- Les scripts utilisant les commandes retirées et les utilisateurs habitués à
  leur corpus personnel devront adapter leur parcours, même si leurs fichiers
  restent intacts. Ne pas ajouter de commandes de remplacement par défaut.
- Les anciens identifiants peuvent désigner des apprentissages personnels :
  ne pas réutiliser ces clés pour le nouveau corpus et lui attribuer de faux acquis.
- La suppression de packages peut emporter des dépendances encore utiles à la
  capture ou aux couleurs ; vérifier les usages avant de retirer une dépendance.
- Un premier chapitre exigeant le mode brut exclurait certains terminaux ;
  garder le démarrage textuel et tester les capacités manquantes.
- Les anciens défauts de persistance et de correction existent encore jusqu'aux
  PR P1 ; ne pas annoncer la fiabilité V2 après les seuls retraits P0.

## Dépendances

Conserver Cobra, Bubble Tea, Bubbles, Lip Gloss, `x/term` et YAML.
Retirer `BurntSushi/toml` seulement après suppression effective de ses usages,
dans une future PR de code. Ne modifier ni `go.mod` ni `go.sum` dans la mission
documentaire. Aucune nouvelle dépendance n'est nécessaire à cette cible.

## Séquence de migration et vérification

Le [backlog](../TODO.md#branches-atomiques-proposees) découpe les retraits en PR
vers `dev`, chacune compilable et testable. D'abord fixer le corpus intégré,
puis retirer la surface dotfiles, enfin ses packages et dépendances ; ensuite
fiabiliser l'entraînement, réviser le corpus et adapter la navigation.

Les tests de migration devront utiliser des copies synthétiques de données
version 1 et d'anciens fichiers YAML/TOML. Vérifier qu'un premier lancement
ne les modifie pas, et que les entraînements conservent les résultats historiques.
Tester le terminal disponible et ses limites, la saisie, les sorties et le
retour TUI. Les cas détaillés figurent dans la [spécification](SPEC.md).

Pour cette mission documentaire : cohérence des liens, lecture du diff,
`git diff --check` et contrôle de la liste autorisée de fichiers suffisent.
Les futures PR Go devront exécuter les tests pertinents, `go test ./...`,
`go vet ./...`, le contrôle de formatage et les commandes d'aide concernées.
