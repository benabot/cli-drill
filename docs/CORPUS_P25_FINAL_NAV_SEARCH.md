# P2.5 — Navigation et recherche modernes

Jalon 3 de la [campagne de finalisation](CURRICULUM.md#enrichissement-par-lots).
Base publiée `d8b383a`. Branche locale `content/p25-unix-git-modern-expansion`.

**330 → 369 exercices**, soit **39 ajouts**
(13 / 7 / 19 Beginner / Intermediate / Advanced).

| Chapitre | Ajouts B / I / A |
| --- | --- |
| [`02-navigation-shell`](../data/chapters/02-navigation-shell.yaml) | 5 / 6 / 10 |
| [`06-recherche-fichiers-contenu`](../data/chapters/06-recherche-fichiers-contenu.yaml) | 8 / 1 / 9 |

Les 39 compétences reprennent la matrice. `02` utilise les chemins pour
s’orienter ou changer de dossier. `06` traite les fichiers, le contenu et
les actions conditionnées par une sélection. Les chemins NUL sont
retravaillés à trois endroits avec des décisions distinctes : présenter
les dossiers, transmettre un fichier sélectionné et traiter plusieurs
résultats sans découpage. Les deux chapitres gardent l’ordre mélangé. La priorité entre motifs
contraires est exercée avec les globs ordonnés de ripgrep. Dans fd,
une exclusion explicite prévaut sur les règles d’ignore.

Essais macOS 27.2.1 avec eza 0.23.5, zoxide 0.10.0, Yazi 26.9.1,
ripgrep 15.2.0, fd 10.5.2 et fzf 0.75.0. Les fixtures ont vérifié tri,
arbre, exclusions, suivi d’un lien, index zoxide, base relative fd,
glob ripgrep, recherche multiligne, JSON et préservation d’un retour à
la ligne dans un chemin NUL. Sous pseudo-terminal, Échap dans fzf a
renvoyé 130 et bloqué l’action aval. Yazi a ouvert le dossier demandé
puis quitté. Les réglages personnalisés de terminal et le presse-papiers
sont des points de recette.

## Vérifications automatisées

`go test -count=1 ./...`, `go vet ./...`, contrôle `gofmt` et
`git diff --check` réussis. Les cas ciblés couvrent réponse, erreur
plausible, niveau et modules. Le contrôle du corpus couvre les variantes
déclarées. La suite de séances vérifie le filtrage et la conservation
de la progression JSON v1 avec des fixtures.

## Références par exercice

### `02`

- `eza-visible-dirs-only`, `eza-sort-modified`, `eza-long-size-units`, `eza-see-link-target`, `eza-tree-exclude-build`, `eza-tree-follow-linked-dir`, `eza-tree-literal-bracket-path`. [EZA:man](https://github.com/eza-community/eza/blob/v0.23.5/man/eza.1.md)
- `zoxide-query-without-cd`. [ZOX:query](https://github.com/ajeetdsouza/zoxide/blob/v0.10.0/man/man1/zoxide-query.1)
- `yazi-open-project-dir`, `yazi-quit-parent-shell-cwd`. [YAZI:quick](https://yazi-rs.github.io/docs/quick-start/)
- `yazi-back-to-parent`, `yazi-copy-current-dirpath`. [YAZI:keys](https://github.com/sxyazi/yazi/blob/v26.9.1/yazi-config/preset/keymap-default.toml)
- `zoxide-add-known-directory`, `zoxide-pick-homonym`, `zoxide-rank-not-mtime`, `zoxide-remove-stale-path`, `zoxide-missing-shell-hook`. [ZOX:readme](https://github.com/ajeetdsouza/zoxide/tree/v0.10.0#readme)
- `fzf-cancel-keep-cwd`, `fd-fzf-directory-nul-input`. [FZF:man](https://github.com/junegunn/fzf/blob/v0.75.0/man/man1/fzf.1)
- `fd-fzf-dir-depth-root`, `fd-base-directory-relative-output`. [FD:man](https://github.com/sharkdp/fd/blob/v10.5.2/doc/fd.1)

### `06`

- `rg-in-single-log`, `rg-scope-one-subtree`. [RG:guide](https://github.com/BurntSushi/ripgrep/blob/15.2.0/GUIDE.md)
- `rg-case-insensitive-query`, `rg-count-lines-per-file`, `rg-invert-health-lines`, `rg-context-line-marker`, `rg-explicit-exclusion-last`, `rg-reinclude-only-audit-log`, `rg-nested-repo-ignore-vcs`, `rg-binary-nul-diagnostic`, `rg-exit-error-versus-no-match`, `rg-match-across-lines`, `rg-structured-match-stream`. [RG:man](https://github.com/BurntSushi/ripgrep/blob/15.2.0/crates/core/flags/defs.rs)
- `fd-regex-prefix-name`. [FD:man](https://github.com/sharkdp/fd/blob/v10.5.2/doc/fd.1)
- `fzf-filter-no-interaction`, `fd-fzf-preserve-file-nul-stream`, `fzf-guard-empty-selection`, `fzf-multi-space-safe`. [FZF:man](https://github.com/junegunn/fzf/blob/v0.75.0/man/man1/fzf.1)
