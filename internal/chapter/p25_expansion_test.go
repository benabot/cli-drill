package chapter

import (
	"slices"
	"testing"

	"github.com/benabot/cli-drill/data"
)

type expansionCase struct {
	key            string
	difficulty     Difficulty
	tools          []string
	correct, wrong string
}

func checkExpansionCases(t *testing.T, cases []expansionCase) {
	t.Helper()
	chapters, err := LoadFS(data.Chapters(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	available := AvailableTools(chapters)
	items := make(map[string]Item)
	for _, ch := range chapters {
		for _, item := range ch.Items {
			items[ch.ID+"/"+item.ID] = item
		}
	}
	for _, tc := range cases {
		t.Run(tc.key, func(t *testing.T) {
			item, ok := items[tc.key]
			if !ok {
				t.Fatal("exercise missing")
			}
			if item.Difficulty != tc.difficulty || !slices.Equal(item.Tools, tc.tools) {
				t.Fatalf("metadata = %s/%v", item.Difficulty, item.Tools)
			}
			if !item.MatchAnswer(tc.correct) || item.MatchAnswer(tc.wrong) {
				t.Fatal("unexpected correction")
			}
			if !item.EligibleForTools(available, tc.tools) {
				t.Fatal("exercise unavailable with required modules")
			}
			for i := range tc.tools {
				without := slices.Delete(slices.Clone(tc.tools), i, i+1)
				if item.EligibleForTools(available, without) {
					t.Fatalf("available without %s", tc.tools[i])
				}
			}
			for _, choice := range item.Choices {
				if item.MatchAnswer(choice) != (choice == item.Answer.Primary) {
					t.Fatalf("unexpected choice correction for %q", choice)
				}
			}
		})
	}
}

func TestP25Milestone1Exercises(t *testing.T) {
	checkExpansionCases(t, []expansionCase{
		{"11-flux-shell/printf-diagnostic-stderr", "intermediate", []string{}, "printf '%s\\n' 'Source absente' >&2", "printf '%s\\n' 'Source absente'"},
		{"11-flux-shell/tail-follow-rotated-log", "intermediate", []string{}, "tail -F service.log", "tail -f service.log"},
		{"11-flux-shell/pipefail-head-sigpipe", "advanced", []string{}, "Le producteur a écrit après la fermeture du pipe.", "head a échoué malgré son statut zéro."},
		{"12-fichiers-permissions/mktemp-work-directory", "intermediate", []string{}, "mktemp -d /tmp/cli-drill.XXXXXX", "mktemp /tmp/cli-drill.XXXXXX"},
		{"12-fichiers-permissions/ls-compare-file-inodes", "intermediate", []string{}, "ls -i brouillon.txt note.txt", "ls -l brouillon.txt note.txt"},
		{"12-fichiers-permissions/rename-destination-write", "advanced", []string{}, "L’écriture dans le dossier de destination.", "L’écriture dans le contenu du fichier."},
		{"12-fichiers-permissions/mv-macos-no-clobber", "intermediate", []string{}, "mv -n brouillon.txt rapport.txt", "mv -i brouillon.txt rapport.txt"},
		{"13-processus-environnement/jobs-with-process-ids", "intermediate", []string{}, "jobs -l", "jobs -p"},
		{"13-processus-environnement/env-client-without-proxy", "intermediate", []string{}, "env -u HTTP_PROXY client", "HTTP_PROXY= client"},
		{"13-processus-environnement/kill-zero-observation", "advanced", []string{}, "Le PID désignait une cible à laquelle tu avais accès.", "Le service attendu fonctionnait normalement."},
		{"13-processus-environnement/nohup-background-collect", "advanced", []string{}, "nohup collecte </dev/null >collecte.log 2>&1 &", "collecte </dev/null >collecte.log 2>&1 &"},
		{"00-premiers-pas/man-keyword-checksum", "beginner", []string{}, "man -k checksum", "man checksum"},
		{"00-premiers-pas/dirname-export-path", "intermediate", []string{}, "dirname build/exports/rapport.txt", "basename build/exports/rapport.txt"},
		{"00-premiers-pas/realpath-existing-export", "intermediate", []string{}, "realpath -e courant/bilan.txt", "realpath -E courant/bilan.txt"},
	})
}
func TestP25Milestone2Exercises(t *testing.T) {
	checkExpansionCases(t, []expansionCase{
		{"10-workflows-dotfiles/git-current-branch", "beginner", []string{"git"}, "git branch --show-current", "git branch --list"},
		{"10-workflows-dotfiles/git-show-last-commit", "beginner", []string{"git"}, "git show HEAD", "git log -1"},
		{"10-workflows-dotfiles/git-log-oneline", "beginner", []string{"git"}, "git log --oneline", "git log --stat"},
		{"10-workflows-dotfiles/git-worktree-diff-path", "beginner", []string{"git"}, "git diff -- config.ini", "git diff --stat -- config.ini"},
		{"10-workflows-dotfiles/git-add-single-path", "beginner", []string{"git"}, "git add -- docs/guide.md", "git add ."},
		{"10-workflows-dotfiles/git-create-branch-in-place", "beginner", []string{"git"}, "git branch relire-doc", "git switch -c relire-doc"},
		{"10-workflows-dotfiles/git-show-historical-file", "beginner", []string{"git"}, "git show v1.2:config.ini", "git show v1.2 -- config.ini"},
		{"10-workflows-dotfiles/git-diff-index-only", "intermediate", []string{"git"}, "git diff --cached", "git diff"},
		{"10-workflows-dotfiles/git-unstage-keep-edit", "intermediate", []string{"git"}, "git restore --staged -- config.ini", "git restore --worktree -- config.ini"},
		{"10-workflows-dotfiles/git-add-patch-file", "intermediate", []string{"git"}, "git add -p -- settings.yml", "git add -- settings.yml"},
		{"10-workflows-dotfiles/git-ignore-source-rule", "intermediate", []string{"git"}, "git check-ignore -v -- build/cache.db", "git status --ignored"},
		{"10-workflows-dotfiles/git-follow-renamed-path", "intermediate", []string{"git"}, "git log --follow -- docs/guide.md", "git log -- docs/guide.md"},
		{"10-workflows-dotfiles/git-log-only-topic", "intermediate", []string{"git"}, "git log main..correction", "git log correction..main"},
		{"10-workflows-dotfiles/git-log-symmetric-divergence", "intermediate", []string{"git"}, "git log --left-right locale...revue", "git log --left-right locale..revue"},
		{"10-workflows-dotfiles/git-show-commit-file-list", "intermediate", []string{"git"}, "git show --name-only a1b2c3d", "git show --stat a1b2c3d"},
		{"10-workflows-dotfiles/git-grep-tag-content", "intermediate", []string{"git"}, "git grep TIMEOUT v2.0", "git grep TIMEOUT"},
		{"10-workflows-dotfiles/git-diff-cached-remote-ref", "intermediate", []string{"git"}, "git diff origin/main HEAD", "git diff origin/main"},
		{"10-workflows-dotfiles/git-status-two-columns", "advanced", []string{"git"}, "Une version est indexée, puis le fichier a encore changé.", "Le fichier a été modifié deux fois dans l’index."},
		{"10-workflows-dotfiles/git-diff-from-common-base", "advanced", []string{"git"}, "git diff main...correctif", "git diff main correctif"},
		{"10-workflows-dotfiles/git-reflog-before-reset", "advanced", []string{"git"}, "b420ffe", "a310ffe"},
		{"10-workflows-dotfiles/git-detached-versus-unborn", "advanced", []string{"git"}, "Dans A uniquement.", "Dans B uniquement."},
		{"10-workflows-dotfiles/git-branch-missing-upstream", "advanced", []string{"git"}, "La branche revue n’a pas de branche amont configurée.", "origin/revue a forcément été supprimée sur le serveur."},
		{"10-workflows-dotfiles/git-stash-inspect-patch", "advanced", []string{"git"}, "git stash show -p 'stash@{2}'", "git stash apply 'stash@{2}'"},
		{"10-workflows-dotfiles/git-show-conflict-stage", "advanced", []string{"git"}, "git show :2:config.ini", "git show :3:config.ini"},
		{"10-workflows-dotfiles/git-ignore-already-tracked", "advanced", []string{"git"}, "La règle ne retire pas un fichier déjà suivi de l’index.", "La règle prendra effet dès la prochaine ouverture du dépôt."},
		{"10-workflows-dotfiles/git-bisect-linear-boundary", "advanced", []string{"git"}, "D", "C"},
		{"10-workflows-dotfiles/git-log-pickaxe-regex", "advanced", []string{"git"}, "git log -G timeout", "git log -S timeout"},
		{"05-outils-quotidiens/atuin-local-stats", "beginner", []string{"atuin"}, "atuin stats", "atuin search"},
		{"05-outils-quotidiens/atuin-search-with-dates", "beginner", []string{"atuin"}, "atuin search --format '{time} {command}' cargo", "atuin search --cmd-only cargo"},
		{"05-outils-quotidiens/lazygit-open-repo-path", "beginner", []string{"lazygit"}, "lazygit -p /atelier/demo", "lazygit --filter /atelier/demo"},
		{"05-outils-quotidiens/lazygit-current-branch-panel", "beginner", []string{"lazygit"}, "main", "revue"},
		{"05-outils-quotidiens/lazygit-modified-files-panel", "beginner", []string{"lazygit"}, "Files", "Commits"},
		{"05-outils-quotidiens/lazygit-view-selected-diff", "beginner", []string{"lazygit"}, "Dans le panneau de diff à droite.", "Dans la liste des branches locales."},
		{"05-outils-quotidiens/lazygit-context-key-help", "beginner", []string{"lazygit"}, "?", "h"},
		{"05-outils-quotidiens/lazygit-commits-panel", "beginner", []string{"lazygit"}, "Commits", "Files"},
		{"05-outils-quotidiens/atuin-search-working-dir", "intermediate", []string{"atuin"}, "atuin search --cwd /atelier/api make", "atuin search /atelier/api make"},
		{"05-outils-quotidiens/atuin-search-date-window", "intermediate", []string{"atuin"}, "atuin search --after 2025-03-01 --before 2025-03-03 docker", "atuin search --before 2025-03-01 --after 2025-03-03 docker"},
		{"05-outils-quotidiens/atuin-open-interactive-query", "intermediate", []string{"atuin"}, "atuin search -i kubectl", "atuin search kubectl"},
		{"05-outils-quotidiens/atuin-exclude-one-exit-code", "intermediate", []string{"atuin"}, "atuin search --exclude-exit 23 rsync", "atuin search --exit 23 rsync"},
		{"05-outils-quotidiens/lazygit-staged-vs-unstaged", "intermediate", []string{"lazygit"}, "config.yml", "notes.md"},
		{"05-outils-quotidiens/lazygit-enter-line-staging", "intermediate", []string{"lazygit"}, "Entrée", "Espace"},
		{"05-outils-quotidiens/lazygit-review-indexed-pane", "intermediate", []string{"lazygit"}, "Staged", "Unstaged"},
		{"05-outils-quotidiens/lazygit-find-local-branch", "intermediate", []string{"lazygit"}, "Filtrer le panneau Branches et lire le résultat.", "Valider le changement de branche sur revue."},
		{"05-outils-quotidiens/lazygit-view-stash-safely", "intermediate", []string{"lazygit"}, "Sélectionner l’entrée et lire son aperçu.", "Appliquer le stash sur la branche actuelle."},
		{"05-outils-quotidiens/atuin-cwd-period-failure", "advanced", []string{"atuin"}, "atuin search --cwd /atelier/api --after 2025-03-01 --exit 2 deploy", "atuin search --cwd /atelier/api --before 2025-03-01 --exit 2 deploy"},
		{"05-outils-quotidiens/atuin-empty-directory-scope", "advanced", []string{"atuin"}, "Élargir la recherche ou la limiter à /atelier/ui.", "Changer le code de sortie en zéro."},
		{"05-outils-quotidiens/atuin-shell-hook-gap", "advanced", []string{"atuin"}, "Cette entrée peut manquer dans Atuin.", "Les deux historiques sont forcément identiques."},
		{"05-outils-quotidiens/atuin-print-command-no-execute", "advanced", []string{"atuin"}, "Aucun, la recherche affiche du texte.", "La dernière suppression est exécutée à nouveau."},
		{"05-outils-quotidiens/lazygit-two-states-one-file", "advanced", []string{"lazygit"}, "Il reste une modification locale non indexée.", "Il entre dans le commit avec le port."},
		{"05-outils-quotidiens/lazygit-stage-only-fix-lines", "advanced", []string{"lazygit"}, "Indexer seulement les lignes de la correction.", "Indexe le fichier complet depuis Files."},
		{"05-outils-quotidiens/lazygit-conflict-current-side", "advanced", []string{"lazygit"}, "Current.", "Incoming."},
		{"05-outils-quotidiens/lazygit-read-upstream-divergence", "advanced", []string{"lazygit"}, "Deux propres à revue, un propre à l’amont.", "Deux fichiers indexés et un fichier non suivi."},
		{"05-outils-quotidiens/lazygit-detached-head-recovery", "advanced", []string{"lazygit"}, "Créer une branche sur le commit courant.", "Supprimer la référence au commit courant."},
		{"05-outils-quotidiens/lazygit-stash-identify-content", "advanced", []string{"lazygit"}, "Parcourir les aperçus des trois entrées.", "Appliquer chaque entrée pour vérifier le délai."},
	})
}
func TestP25Milestone3Exercises(t *testing.T) {
	checkExpansionCases(t, []expansionCase{
		{"02-navigation-shell/eza-visible-dirs-only", "beginner", []string{"eza"}, "eza -D", "eza -d"},
		{"02-navigation-shell/eza-sort-modified", "beginner", []string{"eza"}, "eza --sort=modified exports", "eza --sort=name exports"},
		{"02-navigation-shell/zoxide-query-without-cd", "beginner", []string{"zoxide"}, "zoxide query factures", "z factures"},
		{"02-navigation-shell/yazi-open-project-dir", "beginner", []string{"yazi"}, "yazi projets/api", "yazi"},
		{"02-navigation-shell/yazi-back-to-parent", "beginner", []string{"yazi"}, "h", "H"},
		{"02-navigation-shell/eza-long-size-units", "intermediate", []string{"eza"}, "eza -l --binary", "eza -l --bytes"},
		{"02-navigation-shell/eza-see-link-target", "intermediate", []string{"eza"}, "eza -l déploiement", "eza -l --dereference déploiement"},
		{"02-navigation-shell/eza-tree-exclude-build", "intermediate", []string{"eza"}, "eza --tree --ignore-glob 'build|node_modules' projet", "eza --tree 'build|node_modules' projet"},
		{"02-navigation-shell/zoxide-add-known-directory", "intermediate", []string{"zoxide"}, "zoxide add /atelier/api", "zoxide query /atelier/api"},
		{"02-navigation-shell/zoxide-pick-homonym", "intermediate", []string{"zoxide", "fzf"}, "zi docs", "z docs"},
		{"02-navigation-shell/yazi-copy-current-dirpath", "intermediate", []string{"yazi"}, "c puis d", "c puis c"},
		{"02-navigation-shell/zoxide-rank-not-mtime", "advanced", []string{"zoxide"}, "Les visites et leur récence alimentent le score.", "La dernière modification du contenu décide toujours."},
		{"02-navigation-shell/zoxide-remove-stale-path", "advanced", []string{"zoxide"}, "zoxide remove /atelier/ancien", "zoxide add /atelier/ancien"},
		{"02-navigation-shell/zoxide-missing-shell-hook", "advanced", []string{"zoxide"}, "L’initialisation de zoxide dans ce shell.", "Un nouvel index des fichiers à l’intérieur du dossier."},
		{"02-navigation-shell/fzf-cancel-keep-cwd", "advanced", []string{"fzf"}, "Il reste le même puisque la sélection échoue.", "Il devient ton dossier personnel."},
		{"02-navigation-shell/yazi-quit-parent-shell-cwd", "advanced", []string{"yazi"}, "Dans /atelier.", "Dans /atelier/api."},
		{"02-navigation-shell/eza-tree-follow-linked-dir", "advanced", []string{"eza"}, "Le suivi des liens vers les dossiers dans l’arbre.", "Le tri des dossiers avant les fichiers."},
		{"02-navigation-shell/eza-tree-literal-bracket-path", "advanced", []string{"eza"}, "eza --tree 'versions/[beta]'", "eza --tree versions/[beta]"},
		{"02-navigation-shell/fd-fzf-dir-depth-root", "advanced", []string{"fd", "fzf"}, "fd -t d -d 2 . projets | fzf", "fd -t d . projets | fzf"},
		{"02-navigation-shell/fd-fzf-directory-nul-input", "advanced", []string{"fd", "fzf"}, "fd -t d -0 . projets | fzf --read0", "fd -t d -0 . projets | fzf"},
		{"02-navigation-shell/fd-base-directory-relative-output", "advanced", []string{"fd"}, "La base de recherche change pour fd seulement.", "fd a changé le dossier du shell en /atelier/projets."},
		{"06-recherche-fichiers-contenu/rg-in-single-log", "beginner", []string{"rg"}, "rg timeout api.log", "rg timeout"},
		{"06-recherche-fichiers-contenu/rg-scope-one-subtree", "beginner", []string{"rg"}, "rg migrated db/migrations", "rg migrated db"},
		{"06-recherche-fichiers-contenu/rg-case-insensitive-query", "beginner", []string{"rg"}, "rg -i warning daemon.log", "rg warning daemon.log"},
		{"06-recherche-fichiers-contenu/rg-count-lines-per-file", "beginner", []string{"rg"}, "rg -c retry logs", "rg --count-matches retry logs"},
		{"06-recherche-fichiers-contenu/rg-invert-health-lines", "beginner", []string{"rg"}, "rg -v /health access.log", "rg /health access.log"},
		{"06-recherche-fichiers-contenu/rg-context-line-marker", "beginner", []string{"rg"}, "La ligne 9.", "La ligne 8."},
		{"06-recherche-fichiers-contenu/fd-regex-prefix-name", "beginner", []string{"fd"}, "fd -t f '^backup'", "fd -t f backup"},
		{"06-recherche-fichiers-contenu/fzf-filter-no-interaction", "beginner", []string{"fzf"}, "fzf --filter api", "fzf --query api"},
		{"06-recherche-fichiers-contenu/rg-explicit-exclusion-last", "intermediate", []string{"rg"}, "rg -g '*.env' -g '!*.sample.env' TOKEN configs", "rg -g '!*.sample.env' -g '*.env' TOKEN configs"},
		{"06-recherche-fichiers-contenu/rg-reinclude-only-audit-log", "advanced", []string{"rg"}, "rg -g audit.log WARN exports", "rg --no-ignore WARN exports"},
		{"06-recherche-fichiers-contenu/rg-nested-repo-ignore-vcs", "advanced", []string{"rg"}, "rg --no-ignore-vcs FIXME src", "rg --no-ignore FIXME src"},
		{"06-recherche-fichiers-contenu/fd-fzf-preserve-file-nul-stream", "advanced", []string{"fd", "fzf"}, "fd -t f -0 . imports | fzf --read0 --print0", "fd -t f -0 . imports | fzf --read0"},
		{"06-recherche-fichiers-contenu/fzf-guard-empty-selection", "advanced", []string{"fzf"}, "choix=$(fzf) && less -- \"$choix\"", "less -- \"$(fzf)\""},
		{"06-recherche-fichiers-contenu/fzf-multi-space-safe", "advanced", []string{"fzf"}, "Lire chaque résultat NUL et citer la variable de chemin.", "Lire les résultats séparés par espaces."},
		{"06-recherche-fichiers-contenu/rg-binary-nul-diagnostic", "advanced", []string{"rg"}, "--text", "--hidden"},
		{"06-recherche-fichiers-contenu/rg-exit-error-versus-no-match", "advanced", []string{"rg"}, "Le statut 2 indique une erreur de recherche.", "Le statut 2 indique zéro correspondance sans erreur."},
		{"06-recherche-fichiers-contenu/rg-match-across-lines", "advanced", []string{"rg"}, "rg -U 'BEGIN\\nEND' trace.log", "rg 'BEGIN\\nEND' trace.log"},
		{"06-recherche-fichiers-contenu/rg-structured-match-stream", "advanced", []string{"rg"}, "rg --json token imports", "rg -n --with-filename token imports"},
	})
}
func TestP25Milestone4Exercises(t *testing.T) {
	checkExpansionCases(t, []expansionCase{
		{"07-lecture-preview/less-open-big-log", "beginner", []string{}, "less serveur.log", "cat serveur.log"},
		{"07-lecture-preview/less-forward-one-screen", "beginner", []string{}, "Espace", "Retour arrière"},
		{"07-lecture-preview/less-quit-to-shell", "beginner", []string{}, "q", "s"},
		{"07-lecture-preview/less-start-forward-search", "beginner", []string{}, "/", "?"},
		{"07-lecture-preview/bat-read-ordered-files", "beginner", []string{"bat"}, "bat entete.txt conclusion.txt", "bat conclusion.txt entete.txt"},
		{"07-lecture-preview/bat-read-standard-input", "beginner", []string{"bat"}, "bat -", "bat stdin"},
		{"07-lecture-preview/bat-disable-pager", "beginner", []string{"bat"}, "bat --paging=never guide.txt", "bat --paging=always guide.txt"},
		{"07-lecture-preview/bat-plain-presentation", "beginner", []string{"bat"}, "bat --style=plain guide.txt", "bat --color=never guide.txt"},
		{"07-lecture-preview/bat-show-invisible-characters", "beginner", []string{"bat"}, "bat -A import.txt", "bat -n import.txt"},
		{"07-lecture-preview/bat-force-json-language", "beginner", []string{"bat"}, "bat -l json reponse", "bat -n reponse"},
		{"07-lecture-preview/less-search-previous-match", "intermediate", []string{}, "N majuscule", "n minuscule"},
		{"07-lecture-preview/less-start-at-last-page", "intermediate", []string{}, "less +G serveur.log", "less serveur.log"},
		{"07-lecture-preview/less-enable-line-numbers", "intermediate", []string{}, "less -N config.ini", "less -n config.ini"},
		{"07-lecture-preview/less-no-visual-wrap", "intermediate", []string{}, "less -S trace.txt", "less trace.txt"},
		{"07-lecture-preview/less-next-opened-file", "intermediate", []string{}, ":n", "n"},
		{"07-lecture-preview/bat-highlight-review-range", "intermediate", []string{"bat"}, "bat -H 21:23 service.toml", "bat -r 21:23 service.toml"},
		{"07-lecture-preview/bat-name-anonymous-json-stream", "intermediate", []string{"bat"}, "bat --file-name reponse.json -", "bat reponse.json"},
		{"07-lecture-preview/bat-list-supported-languages", "intermediate", []string{"bat"}, "bat --list-languages", "bat --list-themes"},
		{"07-lecture-preview/bat-read-dash-file", "intermediate", []string{"bat"}, "bat -- -notes.txt", "bat -notes.txt"},
		{"07-lecture-preview/less-pause-live-follow", "advanced", []string{}, "Ctrl+C, puis reprendre la navigation.", "q, puis chercher dans le shell."},
		{"07-lecture-preview/less-repeat-backward-direction", "advanced", []string{}, "Vers le début du fichier.", "Vers la fin du fichier."},
		{"07-lecture-preview/less-color-sgr-pass-through", "advanced", []string{}, "less -R", "less -r"},
		{"07-lecture-preview/bat-no-color-for-consumer", "advanced", []string{"bat"}, "bat --color=never --paging=never releve.json", "bat --color=always --paging=never releve.json"},
		{"07-lecture-preview/bat-override-pager-environment", "advanced", []string{"bat"}, "Désactiver la pagination pour cet appel.", "Choisir un autre thème Bat."},
		{"07-lecture-preview/bat-inspect-binary-byte-view", "advanced", []string{"bat"}, "Afficher les caractères avec le mode visible de Bat.", "Forcer la couleur ANSI seulement."},
		{"07-lecture-preview/bat-read-unknown-stream-end", "advanced", []string{"bat"}, "bat --paging=never -r 25: -", "bat --paging=never -r :25 -"},
		{"08-micro/nano-open-file", "beginner", []string{}, "nano notes.txt", "less notes.txt"},
		{"08-micro/nano-delete-before-cursor", "beginner", []string{}, "Retour arrière", "Suppr"},
		{"08-micro/nano-start-forward-search", "beginner", []string{}, "Ctrl+F", "Ctrl+B"},
		{"08-micro/nano-save-confirm-filename", "beginner", []string{}, "Ctrl+O", "Ctrl+X"},
		{"08-micro/nano-quit-unmodified", "beginner", []string{}, "Ctrl+X", "Ctrl+K"},
		{"08-micro/nano-cut-line", "beginner", []string{}, "Ctrl+K", "Ctrl+X"},
		{"08-micro/nano-paste-cut-line", "beginner", []string{}, "Ctrl+U", "Ctrl+V"},
		{"08-micro/micro-cut-selected-text", "beginner", []string{"micro"}, "Ctrl+X", "Ctrl+C"},
		{"08-micro/micro-paste-clipboard", "beginner", []string{"micro"}, "Ctrl+V", "Ctrl+U"},
		{"08-micro/micro-redo-undo", "beginner", []string{"micro"}, "Ctrl+Y", "Ctrl+Z"},
		{"08-micro/micro-open-editor-help", "beginner", []string{"micro"}, "Ctrl+G", "Ctrl+H"},
		{"08-micro/nano-open-config-at-line", "intermediate", []string{}, "nano +42 service.ini", "nano service.ini 42"},
		{"08-micro/nano-jump-end-buffer", "intermediate", []string{}, "Alt+/", "Alt+\\"},
		{"08-micro/nano-next-forward-match", "intermediate", []string{}, "Alt+F", "Alt+B"},
		{"08-micro/nano-save-copy-other-name", "intermediate", []string{}, "Remplacer le nom proposé par revue.md, puis confirmer.", "Confirmer brouillon.md puis renommer le tampon."},
		{"08-micro/nano-move-line-by-cut-paste", "intermediate", []string{}, "Couper la ligne puis coller au début de la ligne cible.", "Copier la ligne puis coller au début de la ligne cible."},
		{"08-micro/micro-goto-diagnostic-line", "intermediate", []string{"micro"}, "goto 42", "jump 42"},
		{"08-micro/micro-find-next-from-buffer", "intermediate", []string{"micro"}, "Ctrl+N", "Ctrl+P"},
		{"08-micro/micro-replace-with-confirmation", "intermediate", []string{"micro"}, "replace old_host new_host", "replace old_host new_host -a"},
		{"08-micro/micro-save-copy-spaced-name", "intermediate", []string{"micro"}, "save \"revue finale.md\"", "save revue finale.md"},
		{"08-micro/micro-vsplit-compare-files", "intermediate", []string{"micro"}, "vsplit modele.ini", "hsplit modele.ini"},
		{"08-micro/nano-softwrap-preserves-data", "advanced", []string{}, "Le retour visuel à la ligne.", "L’insertion de sauts de ligne lors de la frappe."},
		{"08-micro/nano-save-elsewhere-permission", "advanced", []string{}, "Enregistrer le tampon sous un chemin accessible.", "Quitter et refuser l’enregistrement."},
		{"08-micro/micro-replace-literal-dots", "advanced", []string{"micro"}, "replace v1.2 v1.3 -l -a", "replace v1.2 v1.3 -a"},
		{"08-micro/micro-local-json-filetype", "advanced", []string{"micro"}, "setlocal filetype json", "set filetype json"},
		{"08-micro/micro-remapped-shortcut-diagnostic", "advanced", []string{"micro"}, "Le fichier bindings.json de Micro.", "Le format des fins de ligne du fichier."},
		{"09-markdown-glow/markdown-read-readme-source", "beginner", []string{}, "less README.md", "glow README.md"},
		{"09-markdown-glow/markdown-add-section-heading", "beginner", []string{}, "## Installation", "# Installation"},
		{"09-markdown-glow/markdown-add-unordered-entry", "beginner", []string{}, "- Installation", "1. Installation"},
		{"09-markdown-glow/markdown-link-label-target", "beginner", []string{}, "[Guide](docs/guide.md)", "[docs/guide.md](Guide)"},
		{"09-markdown-glow/markdown-quote-one-line", "beginner", []string{}, "> Relancer le service", "- Relancer le service"},
		{"09-markdown-glow/markdown-inline-command-code", "beginner", []string{}, "`git status`", "**git status**"},
		{"09-markdown-glow/markdown-fenced-example-reading", "beginner", []string{}, "Afficher l’extrait dans un bloc de code.", "Transformer les lignes en liens."},
		{"09-markdown-glow/glow-preview-stdin", "beginner", []string{"glow"}, "glow -", "glow ."},
		{"09-markdown-glow/glow-choose-light-style", "beginner", []string{"glow"}, "glow -s light guide.md", "glow -w 80 guide.md"},
		{"09-markdown-glow/markdown-source-vs-rendered", "beginner", []string{}, "Le lecteur A, qui montre la source.", "Le lecteur B, qui montre le titre rendu."},
		{"09-markdown-glow/markdown-separate-paragraphs", "intermediate", []string{}, "Insérer une ligne vide entre les deux lignes.", "Mettre un seul espace en début de seconde ligne."},
		{"09-markdown-glow/markdown-unintended-indented-code", "intermediate", []string{}, "Retirer l’indentation de la phrase.", "Ajouter quatre espaces supplémentaires."},
		{"09-markdown-glow/markdown-relative-link-from-guide", "intermediate", []string{}, "../README.md", "./README.md"},
		{"09-markdown-glow/markdown-nested-bullet-layout", "intermediate", []string{}, "Sous le P de Préparer.", "Sous le premier tiret, à gauche."},
		{"09-markdown-glow/markdown-gfm-table-delimiter", "intermediate", []string{}, "| --- | --- |", "| Nom | État |"},
		{"09-markdown-glow/markdown-gfm-table-column-mismatch", "intermediate", []string{}, "Ramener la ligne de séparation à deux cellules.", "Ajouter une troisième cellule aux lignes de données."},
		{"09-markdown-glow/markdown-escape-formatting-with-code-span", "intermediate", []string{}, "`**urgent**`", "**urgent**"},
		{"09-markdown-glow/glow-open-local-document-browser", "intermediate", []string{"glow"}, "glow -t docs", "glow -p docs"},
		{"09-markdown-glow/glow-review-after-reading-source", "intermediate", []string{"glow"}, "less notes.md && glow notes.md", "glow notes.md && less notes.md"},
		{"09-markdown-glow/markdown-missing-fence-swallows-heading", "advanced", []string{}, "Fermer le bloc avant le titre de la section.", "Ajouter un dièse au titre de la section."},
		{"09-markdown-glow/markdown-fence-contains-triple-backticks", "advanced", []string{}, "Quatre accents graves à l’ouverture et à la fermeture.", "Trois accents graves autour de l’extrait."},
		{"09-markdown-glow/markdown-gfm-cell-literal-pipe", "advanced", []string{}, "rouge\\|bleu", "rouge|bleu"},
		{"09-markdown-glow/markdown-keep-code-under-list-item", "advanced", []string{}, "Les aligner sous le D de Déploiement.", "Les aligner avec le tiret de la puce."},
		{"09-markdown-glow/markdown-hard-break-vs-wrap", "advanced", []string{}, "Un antislash avant la fin de la ligne source.", "Une largeur de terminal fixée à 80 colonnes."},
		{"09-markdown-glow/glow-preview-balanced-section", "advanced", []string{"glow"}, "Garder le titre et tout le bloc avec ses délimiteurs.", "Terminer l’extrait au milieu du bloc."},
	})
}

func TestP25GuidedCoreAndOptionalModules(t *testing.T) {
	chapters, err := LoadFS(data.Chapters(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	available := AvailableTools(chapters)
	for _, tc := range []struct {
		chapter, module, first string
		core, full             [3]int
	}{
		{"07-lecture-preview", "bat", "less-open-big-log", [3]int{4, 5, 3}, [3]int{11, 11, 8}},
		{"08-micro", "micro", "nano-open-file", [3]int{7, 5, 2}, [3]int{16, 11, 5}},
		{"09-markdown-glow", "glow", "markdown-read-readme-source", [3]int{8, 7, 5}, [3]int{11, 11, 6}},
		{"10-workflows-dotfiles", "git", "dotfiles-status-shell", [3]int{}, [3]int{10, 10, 10}},
	} {
		t.Run(tc.chapter, func(t *testing.T) {
			index := slices.IndexFunc(chapters, func(ch Chapter) bool { return ch.ID == tc.chapter })
			if index < 0 || !chapters[index].Guided() {
				t.Fatal("guided chapter missing")
			}
			ch := chapters[index]
			for i, level := range []Difficulty{DifficultyBeginner, DifficultyIntermediate, DifficultyAdvanced} {
				filter := SessionFilter{Difficulty: level, AvailableTools: available}
				core := filter.EligibleItems(ch, TextExercises)
				if len(core) != tc.core[i] {
					t.Fatalf("%s without modules: %d exercises, want %d", level, len(core), tc.core[i])
				}
				filter.EnabledTools = []string{tc.module}
				learn := filter.Prepare(ch, TextExercises, ModeLearn, func(int, func(int, int)) {
					t.Fatal("guided order shuffled")
				})
				if len(learn.Items) != tc.full[i] || i == 0 && learn.Items[0].ID != tc.first {
					t.Fatalf("%s: unexpected guided session", level)
				}
				review := filter.Prepare(ch, TextExercises, ModeReview, func(n int, swap func(int, int)) {
					for j := 0; j < n/2; j++ {
						swap(j, n-1-j)
					}
				})
				if len(review.Items) != len(learn.Items) {
					t.Fatal("review changed the size")
				}
				for j, item := range learn.Items {
					if j < len(core) && item.ID != core[j].ID {
						t.Fatalf("%s: optional exercise precedes the core", item.ID)
					}
					if item.ID != review.Items[len(review.Items)-1-j].ID {
						t.Fatal("review changed the set or skipped the shuffle")
					}
				}
			}
		})
	}
}

func TestP25IndependentToolSessions(t *testing.T) {
	chapters, err := LoadFS(data.Chapters(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	available := AvailableTools(chapters)
	for _, id := range []string{"02-navigation-shell", "05-outils-quotidiens", "06-recherche-fichiers-contenu"} {
		t.Run(id, func(t *testing.T) {
			index := slices.IndexFunc(chapters, func(ch Chapter) bool { return ch.ID == id })
			if index < 0 || chapters[index].Guided() {
				t.Fatal("independent chapter missing")
			}
			ch := chapters[index]
			for _, level := range []Difficulty{DifficultyBeginner, DifficultyIntermediate, DifficultyAdvanced} {
				filter := SessionFilter{Difficulty: level, AvailableTools: available}
				if len(filter.EligibleItems(ch, TextExercises)) != 0 {
					t.Fatalf("%s is available without its modules", level)
				}
				filter.EnabledTools = available
				eligible := filter.EligibleItems(ch, TextExercises)
				session := filter.Prepare(ch, TextExercises, ModeReview, func(n int, swap func(int, int)) {
					for i := 0; i < n/2; i++ {
						swap(i, n-1-i)
					}
				})
				if len(eligible) != 10 || len(session.Items) != len(eligible) {
					t.Fatalf("%s: unexpected eligible or review count", level)
				}
				for i, item := range session.Items {
					if item.Difficulty != level || item.ID != eligible[len(eligible)-1-i].ID {
						t.Fatalf("%s: review changed the eligible set or skipped the shuffle", item.ID)
					}
				}
			}
		})
	}
}
