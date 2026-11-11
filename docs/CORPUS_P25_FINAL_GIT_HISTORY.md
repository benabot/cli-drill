# P2.5 — Git et historique

Jalon 2 de la [campagne de finalisation](CURRICULUM.md#enrichissement-par-lots).
Base publiée `d8b383a`. Branche locale `content/p25-unix-git-modern-expansion`.

**276 → 330 exercices**, soit **54 ajouts**
(15 / 19 / 20 Beginner / Intermediate / Advanced).

| Chapitre | Ajouts B / I / A |
| --- | --- |
| [`10-workflows-dotfiles`](../data/chapters/10-workflows-dotfiles.yaml) | 7 / 10 / 10 |
| [`05-outils-quotidiens`](../data/chapters/05-outils-quotidiens.yaml) | 8 / 9 / 10 |

Git suit un ordre guidé dans `10`. Les exercices passent du patch de
travail à l’index, puis à l’historique, aux branches, aux conflits et à
la récupération. `05` garde son ordre mélangé. Lazygit réinvestit
les mêmes états dans ses panneaux, avec des décisions sur le fichier
sélectionné, le fragment et l’entrée de stash. Atuin se limite à une
base locale, avec les filtres de date, de dossier et de statut.

Essais macOS 27.2.1 avec Apple Git 2.54.0, Atuin 18.24.6 et Lazygit 0.66.3.
Une copie temporaire du dépôt a permis de vérifier les références, les
diffs de l’index et de l’arbre de travail et la désindexation. Une base
Atuin vide a reçu deux entrées synthétiques pour vérifier dossier, code,
date, format et statistiques. Lazygit a ouvert puis quitté un dépôt
isolé sous pseudo-terminal. Les décisions sur les conflits, les stashes
et la sélection par lignes s’appuient aussi sur les manuels et les
scénarios simulés. Les touches Lazygit sont des réponses textuelles.

## Vérifications automatisées

`go test -count=1 ./...`, `go vet ./...`, contrôle `gofmt` et
`git diff --check` réussis. Les cas ciblés couvrent réponse, erreur
plausible, niveau et modules. Le contrôle du corpus couvre les variantes
déclarées. La suite de séances vérifie le filtrage et la conservation
de la progression JSON v1 avec des fixtures.

## Références par exercice

### `10`

- `git-current-branch`, `git-create-branch-in-place`, `git-branch-missing-upstream`. [GIT:branch](https://git-scm.com/docs/git-branch)
- `git-show-last-commit`, `git-show-historical-file`, `git-show-commit-file-list`. [GIT:show](https://git-scm.com/docs/git-show)
- `git-log-oneline`, `git-follow-renamed-path`, `git-log-only-topic`, `git-log-symmetric-divergence`. [GIT:log](https://git-scm.com/docs/git-log)
- `git-worktree-diff-path`, `git-diff-index-only`, `git-diff-cached-remote-ref`, `git-diff-from-common-base`, `git-log-pickaxe-regex`. [GIT:diff](https://git-scm.com/docs/git-diff)
- `git-add-single-path`, `git-add-patch-file`. [GIT:add](https://git-scm.com/docs/git-add)
- `git-unstage-keep-edit`. [GIT:restore](https://git-scm.com/docs/git-restore)
- `git-ignore-source-rule`, `git-ignore-already-tracked`. [GIT:check-ignore](https://git-scm.com/docs/git-check-ignore)
- `git-grep-tag-content`. [GIT:grep](https://git-scm.com/docs/git-grep)
- `git-status-two-columns`, `git-detached-versus-unborn`. [GIT:status](https://git-scm.com/docs/git-status)
- `git-reflog-before-reset`. [GIT:reflog](https://git-scm.com/docs/git-reflog)
- `git-stash-inspect-patch`. [GIT:stash](https://git-scm.com/docs/git-stash)
- `git-show-conflict-stage`. [GIT:gitrevisions](https://git-scm.com/docs/gitrevisions)
- `git-bisect-linear-boundary`. [GIT:bisect](https://git-scm.com/docs/git-bisect)

### `05`

- `atuin-local-stats`. [ATU:stats](https://docs.atuin.sh/18.24/reference/stats/)
- `atuin-search-with-dates`, `atuin-search-working-dir`, `atuin-search-date-window`, `atuin-open-interactive-query`, `atuin-exclude-one-exit-code`, `atuin-cwd-period-failure`, `atuin-empty-directory-scope`, `atuin-print-command-no-execute`. [ATU:search](https://docs.atuin.sh/18.24/reference/search/)
- `lazygit-open-repo-path`, `lazygit-current-branch-panel`, `lazygit-modified-files-panel`, `lazygit-view-selected-diff`, `lazygit-commits-panel`, `lazygit-staged-vs-unstaged`, `lazygit-review-indexed-pane`, `lazygit-view-stash-safely`, `lazygit-two-states-one-file`, `lazygit-stage-only-fix-lines`, `lazygit-conflict-current-side`, `lazygit-read-upstream-divergence`, `lazygit-detached-head-recovery`, `lazygit-stash-identify-content`. [LG:readme](https://github.com/jesseduffield/lazygit/tree/v0.66.3#features)
- `lazygit-context-key-help`, `lazygit-enter-line-staging`, `lazygit-find-local-branch`. [LG:keys](https://github.com/jesseduffield/lazygit/blob/v0.66.3/docs/keybindings/Keybindings_en.md)
- `atuin-shell-hook-gap`. [ATU:integration](https://docs.atuin.sh/cli/guide/installation/)
