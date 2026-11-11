<div align="center">

# CLI-Drill

<!-- Keep this anchor without href to prevent GitHub image auto-linking. -->
<a name="cli-drill-hero">
  <img src="docs/assets/cli-drill-hero.webp" width="820"
    alt="Bannière CLI-Drill et menu Browse chapters dans le terminal.">
</a>

**Learn the terminal by doing.**

[Français](#français) · [English](#english) · [État du projet](PROJECT_STATE.md) · [Recette P2.6](docs/P26_FINAL_ACCEPTANCE.md)

</div>

## Français

CLI-Drill est un entraînement local à la pratique du terminal. La branche
`dev` intègre **446 exercices dans 13 chapitres**. Elle couvre les commandes
Unix, l’édition de ligne dans Bash et Zsh, Git et des outils CLI optionnels.
Les exercices combinent commandes simulées, questions à choix et capture
de raccourcis. Les commandes soumises sont corrigées comme du texte.

<a id="installer-cli-drill"></a>

### Installation locale

Depuis un checkout de `dev`, avec Go (le module déclare Go 1.23) :

```sh
go run ./cmd/cli-drill
```

Le binaire embarque le corpus. L’entraînement fonctionne hors ligne. La TUI
demande un terminal en entrée et en sortie. Pour un terminal sans TUI, utiliser
`go run ./cmd/cli-drill train`.

<a id="demarrage-rapide"></a>

### Démarrage rapide

Dans la TUI, **Start training** ouvre le premier chapitre textuel éligible.
**Browse chapters** permet de chercher avec `/` et d’ouvrir une fiche avec
Entrée. La fiche propose `d` pour le niveau, `m` pour le mode et `a` pour
passer au parcours intégral. Utiliser `PgUp` et `PgDn` pour lire un contenu
long, `h` pour l’aide de navigation et Échap pour revenir.

```sh
go run ./cmd/cli-drill chapters
go run ./cmd/cli-drill train
go run ./cmd/cli-drill train 00-premiers-pas --mode review
go run ./cmd/cli-drill train 14-recherche-texte-unix --difficulty intermediate --all
go run ./cmd/cli-drill stats
go run ./cmd/cli-drill prefs
```

Un lot contient jusqu’à six exercices. **Apprendre** part des exercices
non acquis. **Réviser** mélange tous les exercices du filtre. **Continuer**
poursuit la visite, et **Revoir mes erreurs** ouvre une reprise ciblée.
Au bilan, utiliser `c`, `r` ou `q` pour terminer. En CLI textuelle,
valider la lettre par Entrée. `--all` regroupe le filtre entier dans un lot.

**Keyboard shortcuts** ouvre la capture dédiée. La compatibilité dépend des
octets transmis par le clavier et l’émulateur. Pour retirer les couleurs,
lancer avec `NO_COLOR=1`. Les cibles sont macOS, Linux et Unix/BSD.
Les environnements effectivement vérifiés sont indiqués dans le
[rapport de recette](docs/P26_FINAL_ACCEPTANCE.md).

### Préférences et progression

Git essentiel est actif en l’absence de fichier de préférences. Activer
les autres modules depuis **Preferences** dans la TUI, ou avec
`go run ./cmd/cli-drill prefs --enable fd`. Les modules pédagogiques
ne nécessitent aucune installation. Les choix enregistrés priment sur les
valeurs initiales.

Chaque réponse soumise enrichit la progression locale JSON v1. Apprendre
repart des acquis après une interruption. L’ordre et le curseur d’une visite
restent en mémoire. La progression est enregistrée sous
`$XDG_DATA_HOME/cli-drill/progress.json`, ou dans
`~/.local/share/cli-drill/progress.json` par défaut. `--progress` sélectionne
un autre fichier. Utiliser un seul processus écrivain par fichier.

## English

CLI-Drill is a local terminal trainer. The `dev` branch includes **446 exercises
across 13 chapters** covering Unix commands, Bash and Zsh line editing, Git
and optional CLI tools. Exercises use simulated commands, multiple choice and
live shortcut capture. Submitted commands are checked as text.

### Installation

From a checkout of `dev`, with Go (the module declares Go 1.23):

```sh
go run ./cmd/cli-drill
```

The corpus is built in and training works offline. The TUI requires a terminal
for both input and output. For text training, run `go run ./cmd/cli-drill train`.

### Quick start

In the TUI, choose **Start training** or **Browse chapters**. Use `/` to
search, Enter to open a chapter, `d` to change difficulty, `m` to switch
between Learn and Review, and `a` for the full chapter selection. Use
`PgUp` and `PgDn` to read long content, `h` for navigation help and Escape
to go back.

```sh
go run ./cmd/cli-drill chapters
go run ./cmd/cli-drill train 00-premiers-pas --mode review
go run ./cmd/cli-drill train 14-recherche-texte-unix --difficulty intermediate --all
go run ./cmd/cli-drill stats
go run ./cmd/cli-drill prefs
```

Batches contain up to six exercises. Learn selects exercises you have not
acquired. Review draws from the entire filtered set. At the summary, use
`c` to continue, `r` to review mistakes and `q` to finish. In the text CLI,
press Enter after the letter. Use `--all` to run the whole filtered set in one
batch. Live capture depends on the bytes sent by the keyboard and terminal.

Git is enabled when there is no preferences file. Other modules are optional.
Saved preferences take precedence. Progress is stored in local JSON v1.
Use `--progress` for a separate file and one writer per file. Acquired exercises
carry over to the next Learn visit. The visit cursor is kept in memory.

The current corpus and some interface labels are in French. See the
[project state](PROJECT_STATE.md) and the [P2.6 acceptance report](docs/P26_FINAL_ACCEPTANCE.md)
for tested environments and remaining checks.

<a id="development"></a>
<a id="developper-et-contribuer"></a>

## Développement / Development

```sh
go test -count=1 ./...
go test -race -count=1 ./...
go vet ./...
```

Pour proposer une correction, décrire le problème et les validations effectuées.
To propose a fix, describe the issue and include your verification results.
[Issues](https://github.com/benabot/cli-drill/issues)

[Spécification](docs/SPEC.md) · [Architecture](docs/ARCHITECTURE.md) · [Décisions](docs/DECISIONS.md) · [Backlog](TODO.md) · [Consignes](AGENTS.md)

Créé par / Created by [benabot](https://github.com/benabot).
