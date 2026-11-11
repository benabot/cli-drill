# P2.5 — Lecture, édition et Markdown

Jalon 4 de la [campagne de finalisation](CURRICULUM.md#enrichissement-par-lots).
Base publiée `d8b383a`. Branche locale `content/p25-unix-git-modern-expansion`.

**369 → 446 exercices**, soit **77 ajouts**
(31 / 28 / 18 Beginner / Intermediate / Advanced).

| Chapitre | Ajouts B / I / A |
| --- | --- |
| [`07-lecture-preview`](../data/chapters/07-lecture-preview.yaml) | 10 / 9 / 7 |
| [`08-micro`](../data/chapters/08-micro.yaml) | 11 / 10 / 5 |
| [`09-markdown-glow`](../data/chapters/09-markdown-glow.yaml) | 10 / 9 / 6 |

Les matrices de `07`, `08` et `09` ont été adaptées avant la rédaction.
Les 77 ajouts remplacent la projection historique de 76. Le volume de
`07` reste à +26. `08` passe de +24 à +26 pour permettre le parcours
Nano puis Micro. `09` passe de +26 à +25 et atteint 29 exercices. Les
options spécialisées écartées incluent les thèmes Bat installés, des
réglages de présentation Micro et des scénarios exclusivement Glow.
Les 14 exercices publiés de ces trois chapitres sont conservés.

Les trois chapitres rejoignent la famille `unix-shell` et utilisent
l’apprentissage guidé. Le filtre consulte les outils de l’exercice.
Le socle donne 12 exercices less, 14 Nano et 20 Markdown/source avec
`tools: []`. La pratique hors de CLI-Drill suppose que le lecteur ou
l’éditeur est disponible sur la machine. Les modules Bat, Micro et Glow
restent activables séparément. Le scénario historique fd + fzf + Glow
conserve ses trois prérequis.

La coupe et le collage se retravaillent avec deux éditeurs, tandis que
la lecture en `less` se prolonge dans `09` par la comparaison de la source
et du rendu. Les tableaux sont identifiés comme extension GFM.
La CLI laisse la fin de ligne dans la réponse reçue. Le correcteur peut
la conserver lorsqu’elle suit un antislash final. Le saut forcé Markdown
est donc évalué par un diagnostic à choix. Le moteur reste inchangé.

Sous pseudo-terminal sur macOS 27.2.1, essais de dernière page dans less
668, de coupe et d’écriture dans Nano 9.3 et de remplacement littéral
puis sauvegarde dans Micro 2.0.15. Les essais CLI ont vérifié le flux
et le mode visible de Bat 0.26.1 et l’entrée standard de Glow 3.0.0.
Le tableau GFM a été vérifié avec Pandoc en mode GFM. Les touches
d’éditeur restent des réponses textuelles dans CLI-Drill. La recette P2.6
doit couvrir les raccourcis et le rendu dans des terminaux représentatifs,
les modules désactivés et la durée des séances.

La CLI compilée a démarré les trois parcours Beginner avec des
préférences vierges. Les séances ont proposé respectivement 4, 7 et 8
exercices. Le premier exercice de chaque séance a été validé et les
acquis synthétiques existants sont restés présents.

## Vérifications automatisées

`go test -count=1 ./...`, `go vet ./...`, contrôle `gofmt` et
`git diff --check` réussis. Les cas ciblés couvrent réponse, erreur
plausible, niveau et modules. Le contrôle du corpus couvre les variantes
déclarées. La suite de séances vérifie le filtrage et la conservation
de la progression JSON v1 avec des fixtures.

## Références par exercice

### `07`

- `less-open-big-log`, `less-forward-one-screen`, `less-quit-to-shell`, `less-start-forward-search`, `less-search-previous-match`, `less-start-at-last-page`, `less-enable-line-numbers`, `less-no-visual-wrap`, `less-next-opened-file`, `less-pause-live-follow`, `less-repeat-backward-direction`, `less-color-sgr-pass-through`. [LESS:man](https://github.com/gwsw/less/blob/v668/less.nro.VER)
- `bat-read-ordered-files`, `bat-read-standard-input`, `bat-disable-pager`, `bat-plain-presentation`, `bat-show-invisible-characters`, `bat-force-json-language`, `bat-highlight-review-range`, `bat-name-anonymous-json-stream`, `bat-list-supported-languages`, `bat-read-dash-file`, `bat-no-color-for-consumer`, `bat-override-pager-environment`, `bat-inspect-binary-byte-view`, `bat-read-unknown-stream-end`. [BAT:man](https://github.com/sharkdp/bat/blob/v0.26.1/src/bin/bat/app.rs)

### `08`

- `nano-open-file`, `nano-open-config-at-line`, `nano-save-copy-other-name`, `nano-softwrap-preserves-data`, `nano-save-elsewhere-permission`. [NANO:manual](https://www.nano-editor.org/dist/latest/nano.html)
- `nano-delete-before-cursor`, `nano-start-forward-search`, `nano-save-confirm-filename`, `nano-quit-unmodified`, `nano-cut-line`, `nano-paste-cut-line`, `nano-jump-end-buffer`, `nano-next-forward-match`, `nano-move-line-by-cut-paste`. [NANO:keys](https://www.nano-editor.org/dist/latest/cheatsheet.html)
- `micro-cut-selected-text`, `micro-paste-clipboard`, `micro-redo-undo`, `micro-open-editor-help`, `micro-find-next-from-buffer`, `micro-remapped-shortcut-diagnostic`. [MIC:keys](https://github.com/zyedidia/micro/blob/v2.0.15/runtime/help/defaultkeys.md)
- `micro-goto-diagnostic-line`, `micro-replace-with-confirmation`, `micro-save-copy-spaced-name`, `micro-vsplit-compare-files`, `micro-replace-literal-dots`. [MIC:commands](https://github.com/zyedidia/micro/blob/v2.0.15/runtime/help/commands.md)
- `micro-local-json-filetype`. [MIC:options](https://github.com/zyedidia/micro/blob/v2.0.15/runtime/help/options.md)

### `09`

- `markdown-read-readme-source`. [LESS:man](https://github.com/gwsw/less/blob/v668/less.nro.VER)
- `markdown-add-section-heading`, `markdown-source-vs-rendered`. [MD:atx](https://spec.commonmark.org/0.31.2/#atx-headings)
- `markdown-add-unordered-entry`, `markdown-nested-bullet-layout`, `markdown-keep-code-under-list-item`. [MD:lists](https://spec.commonmark.org/0.31.2/#list-items)
- `markdown-link-label-target`, `markdown-relative-link-from-guide`. [MD:links](https://spec.commonmark.org/0.31.2/#links)
- `markdown-quote-one-line`. [MD:quotes](https://spec.commonmark.org/0.31.2/#block-quotes)
- `markdown-inline-command-code`, `markdown-escape-formatting-with-code-span`. [MD:code](https://spec.commonmark.org/0.31.2/#code-spans)
- `markdown-fenced-example-reading`, `markdown-missing-fence-swallows-heading`, `markdown-fence-contains-triple-backticks`. [MD:fences](https://spec.commonmark.org/0.31.2/#fenced-code-blocks)
- `markdown-separate-paragraphs`. [MD:paragraphs](https://spec.commonmark.org/0.31.2/#paragraphs)
- `markdown-unintended-indented-code`. [MD:indent](https://spec.commonmark.org/0.31.2/#indented-code-blocks)
- `markdown-gfm-table-delimiter`, `markdown-gfm-table-column-mismatch`, `markdown-gfm-cell-literal-pipe`. [MD:tables](https://github.github.com/gfm/#tables-extension-)
- `markdown-hard-break-vs-wrap`. [MD:breaks](https://spec.commonmark.org/0.31.2/#hard-line-breaks)
- `glow-preview-stdin`, `glow-choose-light-style`, `glow-open-local-document-browser`, `glow-review-after-reading-source`, `glow-preview-balanced-section`. [GLOW:usage](https://github.com/charmbracelet/glow/tree/v3.0.0#usage)
