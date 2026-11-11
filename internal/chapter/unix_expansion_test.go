package chapter

import (
	"slices"
	"testing"

	"github.com/benabot/cli-drill/data"
	"github.com/benabot/cli-drill/internal/exercise"
)

func TestUnixExpansionSessions(t *testing.T) {
	chapters, err := LoadFS(data.Chapters(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}

	for _, lesson := range []struct {
		id    string
		order LearningOrder
		ids   [3][]string // Beginner, Intermediate, Advanced.
	}{
		{"00-premiers-pas", LearningOrderGuided, [3][]string{
			{"pwd", "ls", "ls-hidden", "cd-docs", "cd-parent", "cd-absolute", "mkdir-essais", "cat-spaces", "cd-home", "man-ls", "ls-long-by-modification", "man-keyword-checksum"},
			{"cd-previous", "ls-other-directory", "mkdir-parents", "cat-hidden-path", "ls-directory-entry", "glob-print-names", "glob-bracket-list", "quote-literal-question", "literal-tilde-directory", "command-origin", "man-external-printf-section", "pwd-physical-directory", "dirname-export-path", "realpath-existing-export"},
			{"cdpath-priority", "cd-symlink-parent", "glob-quoted-prefix", "glob-no-match-shells", "glob-leading-dot", "builtin-shadowed", "subshell-directory", "relative-from-peer", "dash-leading-name", "case-sensitive-volume"},
		}},
		{"11-flux-shell", LearningOrderGuided, [3][]string{
			{"redirect-output", "append-output", "stdin-report", "preview-first-lines", "preview-last-lines", "count-file-words", "pipe-produced-words", "sort-full-lines", "dedupe-adjacent-lines", "tee-new-copy", "tail-follow-journal"},
			{"redirect-input", "count-matching-lines", "redirect-errors", "sort-dedupe-scattered", "count-sorted-runs", "sort-numeric-field", "cut-delimited-field", "tee-append-copy", "split-output-diagnostics", "heredoc-literal-start", "redirect-missing-parent", "noclobber-single-export", "show-last-exit-status", "printf-diagnostic-stderr", "tail-follow-rotated-log"},
			{"merge-output-errors", "pipefail-rightmost-failure", "bash-pipeline-status-vector", "nul-delimited-sort", "sort-before-replacement", "exec-persistent-stdout", "redirect-brace-group", "heredoc-expanded-body", "sort-locale-precedence", "newline-vs-record-count", "redirect-stderr-before-stdout", "and-or-alert-after-publish", "pipefail-head-sigpipe"},
		}},
		{"12-fichiers-permissions", LearningOrderIndependent, [3][]string{
			{"copy-spaced-path", "rename-file", "remove-file", "remove-empty-directory", "list-permissions", "touch-new-file", "copy-backup-name", "read-owner-column", "chmod-group-no-write", "recognize-link-entry", "mkdir-private-mode"},
			{"owner-execute", "copy-directory-tree", "copy-confirm-overwrite", "move-many-sources", "copy-preserve-times", "link-relative-target", "readlink-stored-target", "chmod-two-classes", "interpret-octal-mode", "missing-directory-search", "remove-tree-with-confirmation", "copy-tree-keep-symlinks", "chmod-exact-config-mode", "du-directory-kib", "df-target-filesystem-kib", "mktemp-work-directory", "ls-compare-file-inodes", "mv-macos-no-clobber"},
			{"remove-symlink-only", "umask-file-and-directory", "unlink-vs-write", "sticky-shared-directory", "setgid-inherit-group", "move-relative-symlink", "hardlink-after-replace", "chmod-conditional-search", "copy-exact-destination-gnu", "stat-mode-macos", "move-cross-filesystem-interrupted", "rename-destination-write"},
		}},
		{"13-processus-environnement", LearningOrderIndependent, [3][]string{
			{"show-environment-variable", "set-shell-variable", "ps-terminal-processes", "read-process-id", "list-exported-environment", "unset-shell-variable", "list-shell-jobs", "launch-background-job", "foreground-known-job", "empty-exported-variable", "run-executable-relative-path"},
			{"inspect-pid", "request-termination", "export-variable", "one-command-environment", "ps-pid-parent-columns", "last-background-pid", "wait-specific-child", "resume-suspended-in-background", "wait-child-exit-status", "launch-with-clean-environment", "path-prepend-for-one-call", "find-process-by-exact-name", "read-job-control-stop-state", "diagnose-command-launch-status", "jobs-with-process-ids", "env-client-without-proxy"},
			{"identify-job-process-group", "signal-service-cleanup", "background-terminal-read-stop", "path-change-clears-cache", "export-after-child-launch", "subshell-variable-isolation", "load-script-current-shell", "cwd-with-empty-environment", "read-zombie-process-state", "trap-scope-subshell", "force-stop-verified-pid", "kill-zero-observation", "nohup-background-collect"},
		}},
	} {
		t.Run(lesson.id, func(t *testing.T) {
			index := slices.IndexFunc(chapters, func(ch Chapter) bool { return ch.ID == lesson.id })
			if index < 0 {
				t.Fatal("chapter is missing from the embedded corpus")
			}
			ch := chapters[index]
			if ch.Family != FamilyUnixShell || ch.LearningOrder != lesson.order {
				t.Fatalf("unexpected family/order: %s / %s", ch.Family, ch.LearningOrder)
			}
			for _, item := range ch.Items {
				if len(item.Tools) != 0 {
					t.Fatalf("%s: unexpected tool modules %v", item.ID, item.Tools)
				}
				switch item.ExerciseType {
				case exercise.TypeSimpleShellSim, exercise.TypeFreeAnswer, exercise.TypeMultipleChoice:
				default:
					t.Fatalf("%s: unexpected exercise format %s", item.ID, item.ExerciseType)
				}
			}

			for i, level := range []Difficulty{DifficultyBeginner, DifficultyIntermediate, DifficultyAdvanced} {
				t.Run(string(level), func(t *testing.T) {
					want := lesson.ids[i]
					filter := SessionFilter{Difficulty: level} // No third-party modules.
					for _, mode := range []SessionMode{ModeLearn, ModeReview} {
						shuffles := 0
						session := filter.Prepare(ch, TextExercises, mode, func(n int, swap func(int, int)) {
							shuffles++
							for i := 0; i < n/2; i++ {
								swap(i, n-1-i)
							}
						})
						expected := slices.Clone(want)
						wantShuffles := 0
						if len(want) > 1 && (lesson.order == LearningOrderIndependent || mode == ModeReview) {
							slices.Reverse(expected)
							wantShuffles = 1
						}
						var got []string
						for _, item := range session.Items {
							got = append(got, item.ID)
						}
						if !slices.Equal(got, expected) || shuffles != wantShuffles {
							t.Fatalf("%s: got %v with %d shuffles, want %v with %d", mode, got, shuffles, expected, wantShuffles)
						}
					}
				})
			}
		})
	}
}

func TestUnixExpansionRejectsCommandsWithDifferentEffects(t *testing.T) {
	chapters, err := LoadFS(data.Chapters(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	items := make(map[string]Item)
	for _, ch := range chapters {
		for _, item := range ch.Items {
			items[ch.ID+"/"+item.ID] = item
		}
	}
	for _, tc := range []struct{ key, valid, wrong string }{
		{"11-flux-shell/redirect-output", "ls docs >inventaire.txt", "ls docs >> inventaire.txt"},
		{"11-flux-shell/append-output", "cat note.txt >> journal.txt", "cat note.txt > journal.txt"},
		{"11-flux-shell/redirect-input", "wc -l < messages.txt", "wc -l messages.txt"},
		{"11-flux-shell/count-matching-lines", "grep ERROR serveur.log|wc -l", "grep -o ERROR serveur.log | wc -l"},
		{"11-flux-shell/redirect-errors", "ls rapport.txt 2> erreurs.log", "ls rapport.txt > erreurs.log"},
		{"11-flux-shell/merge-output-errors", "ls docs absent 1> bilan.log 2>&1", "ls docs absent 2>&1 > bilan.log"},
		{"12-fichiers-permissions/copy-spaced-path", `cp 'rapport final.txt' archives/`, "cp rapport final.txt archives/"},
		{"12-fichiers-permissions/remove-empty-directory", "rmdir essais", "rm -r essais"},
		{"12-fichiers-permissions/owner-execute", "chmod u+x lancer.sh", "chmod a+x lancer.sh"},
		{"13-processus-environnement/inspect-pid", "ps -p 4242", "ps -p 4243"},
		{"13-processus-environnement/request-termination", "kill -TERM 4242", "kill -KILL 4242"},
		{"13-processus-environnement/set-shell-variable", "ATELIER_MODE=test", "export ATELIER_MODE=test"},
		{"13-processus-environnement/export-variable", "export ATELIER_MODE=test", "ATELIER_MODE=test"},
		{"13-processus-environnement/one-command-environment", "env LC_ALL=C sort noms.txt", "export LC_ALL=C; sort noms.txt"},
		{"12-fichiers-permissions/remove-tree-with-confirmation", "rm -ri ./cache-sortie", "rm -rf ./cache-sortie"},
		{"12-fichiers-permissions/copy-tree-keep-symlinks", "cp -RP config config-copie", "cp -RL config config-copie"},
		{"12-fichiers-permissions/chmod-exact-config-mode", "chmod 640 service.conf", "chmod u+rw,g+r service.conf"},
		{"12-fichiers-permissions/mkdir-private-mode", "mkdir -m 700 prive", "mkdir prive"},
		{"12-fichiers-permissions/du-directory-kib", "du -sk archives", "du -k archives"},
		{"12-fichiers-permissions/df-target-filesystem-kib", "df -Pk /mnt/backup/exports", "du -sk /mnt/backup/exports"},
		{"12-fichiers-permissions/move-cross-filesystem-interrupted", "Une copie et une suppression sont deux étapes distinctes.", "Le renommage est atomique entre les deux volumes."},
		{"13-processus-environnement/run-executable-relative-path", "./releve", "releve"},
		{"13-processus-environnement/path-prepend-for-one-call", "PATH=\"/opt/test/bin:$PATH\" analyse", "PATH=\"/opt/test/bin\" analyse"},
		{"13-processus-environnement/find-process-by-exact-name", "pgrep -x agentweb", "pgrep agentweb"},
		{"13-processus-environnement/read-job-control-stop-state", "4812", "4813"},
		{"13-processus-environnement/diagnose-command-launch-status", "127 126", "126 127"},
		{"13-processus-environnement/force-stop-verified-pid", "kill -s KILL 4821", "kill -s TERM 4821"},
		{"00-premiers-pas/ls-long-by-modification", "ls -lt livraisons", "ls -ltr livraisons"},
		{"00-premiers-pas/man-external-printf-section", "man 1 printf", "man 3 printf"},
		{"00-premiers-pas/pwd-physical-directory", "pwd -P", "pwd -L"},
		{"11-flux-shell/tail-follow-journal", "tail -f service.log", "tail -n 10 service.log"},
		{"11-flux-shell/redirect-missing-parent", "Le shell échoue avant de lancer mesurer.", "mesurer démarre puis échoue en écrivant sa sortie."},
		{"11-flux-shell/noclobber-single-export", "(set -C; exporter > rapport.txt)", "(set -C; exporter >| rapport.txt)"},
		{"11-flux-shell/redirect-stderr-before-stdout", "stderr garde la destination initiale de stdout.", "stderr suit stdout vers mesures.log."},
		{"11-flux-shell/show-last-exit-status", "printf '%s\\n' \"$?\"", "printf '%s\\n' \"$$\""},
		{"11-flux-shell/and-or-alert-after-publish", "alerter s’exécute après l’échec de publier.", "alerter ne s’exécute que si preparer échoue."},
	} {
		t.Run(tc.key, func(t *testing.T) {
			item, ok := items[tc.key]
			if !ok {
				t.Fatal("exercise is missing from the embedded corpus")
			}
			if !item.MatchAnswer(tc.valid) || item.MatchAnswer(tc.wrong) {
				t.Fatalf("want to accept %q and reject %q", tc.valid, tc.wrong)
			}
			for _, choice := range item.Choices {
				if item.MatchAnswer(choice) != (choice == item.Answer.Primary) {
					t.Errorf("unexpected verdict for choice %q", choice)
				}
			}
		})
	}
}

// Check a canonical response and a plausible misconception for every new item.
// TestV2Corpus separately checks all declared variants through the same matcher.
func TestWave1Correction(t *testing.T) {
	chapters, err := LoadFS(data.Chapters(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	items := make(map[string]Item)
	for _, ch := range chapters {
		for _, item := range ch.Items {
			items[item.ID] = item
		}
	}
	for _, tc := range []struct{ id, valid, wrong string }{
		{"cd-home", "cd ~", "cd '~'"},
		{"man-ls", "man ls", "ls --help"},
		{"cd-previous", "cd -", "cd .."},
		{"ls-other-directory", "ls /atelier/livraisons", "cd /atelier/livraisons && ls"},
		{"mkdir-parents", "mkdir -p livraisons/juin/brut", "mkdir livraisons/juin/brut"},
		{"cat-hidden-path", "cat /atelier/.etat", "cat /atelier/etat"},
		{"ls-directory-entry", "ls -d sorties", "ls sorties"},
		{"glob-print-names", "printf '%s\\n' scan-*.png", "printf '%s\\n' 'scan-*.png'"},
		{"glob-bracket-list", "plan-[135].svg", "plan-[1-5].svg"},
		{"quote-literal-question", "cat 'lot?.csv'", "cat lot?.csv"},
		{"literal-tilde-directory", "cat '~/journal.txt'", "cat ~/journal.txt"},
		{"command-origin", "command -v archiver", "which archiver"},
		{"cdpath-priority", "/partage/images", "/atelier/images"},
		{"cd-symlink-parent", "/atelier /srv", "/srv /atelier"},
		{"glob-quoted-prefix", "printf '%s\\n' \"$lot\"-*.csv", "printf '%s\\n' \"$lot-*.csv\""},
		{"glob-no-match-shells", "Bash transmet le motif ; Zsh refuse la ligne.", "Bash et Zsh transmettent le motif."},
		{"glob-leading-dot", "menu.json", ".menu.json"},
		{"builtin-shadowed", "builtin printf '%s\\n' prêt", "printf '%s\\n' prêt"},
		{"subshell-directory", "(cd exports && pwd)", "{ cd exports && pwd; }"},
		{"relative-from-peer", "../../plans/est/coupe.svg", "../plans/est/coupe.svg"},
		{"dash-leading-name", "cat -- -notes.txt", "cat '-notes.txt'"},
		{"case-sensitive-volume", "Seulement sur le volume insensible à la casse.", "Sur les deux volumes."},
		{"touch-new-file", "touch accueil.txt", "touch -c accueil.txt"},
		{"copy-backup-name", "cp devis.txt devis.bak", "mv devis.txt devis.bak"},
		{"read-owner-column", "nora", "studio"},
		{"chmod-group-no-write", "chmod g-w maquette.txt", "chmod a-w maquette.txt"},
		{"recognize-link-entry", "edition", "edition.txt"},
		{"copy-directory-tree", "cp -R modeles livraison", "cp modeles livraison"},
		{"copy-confirm-overwrite", "cp -i index.txt depot/index.txt", "cp -n index.txt depot/index.txt"},
		{"move-many-sources", "mv rouge.txt bleu.txt collecte/", "mv *.txt collecte/"},
		{"copy-preserve-times", "cp -p source.txt copie.txt", "cp source.txt copie.txt"},
		{"link-relative-target", "ln -s ../source/fiche.txt vues/fiche", "ln -s source/fiche.txt vues/fiche"},
		{"readlink-stored-target", "readlink raccourci", "readlink -f raccourci"},
		{"chmod-two-classes", "chmod g+r,o-w bilan.txt", "chmod go+r-w bilan.txt"},
		{"interpret-octal-mode", "rw-r-x--x", "rw-r-xr--"},
		{"missing-directory-search", "Exécution", "Lecture"},
		{"remove-symlink-only", "Seul acces disparaît.", "Les deux noms disparaissent."},
		{"umask-file-and-directory", "640 750", "637 750"},
		{"unlink-vs-write", "Supprimer le nom seulement.", "Modifier le contenu seulement."},
		{"sticky-shared-directory", "Mon fichier seulement.", "Les deux fichiers."},
		{"setgid-inherit-group", "atelier", "bleu"},
		{"move-relative-symlink", "/atelier/archives/source/plan.txt", "/atelier/source/plan.txt"},
		{"hardlink-after-replace", "temoin seulement.", "temoin et vue."},
		{"chmod-conditional-search", "chmod -R g+X depot", "chmod -R g+x depot"},
		{"copy-exact-destination-gnu", "cp -T plan.txt cible", "cp -i plan.txt cible"},
		{"stat-mode-macos", "stat -f %OLp fiche.txt", "stat -c %a fiche.txt"},
	} {
		t.Run(tc.id, func(t *testing.T) {
			item, ok := items[tc.id]
			if !ok {
				t.Fatal("exercise is missing from the embedded corpus")
			}
			if item.Answer.Primary != tc.valid || !item.MatchAnswer(tc.valid) || item.MatchAnswer(tc.wrong) {
				t.Fatalf("want primary %q accepted and %q rejected", tc.valid, tc.wrong)
			}
			for _, choice := range item.Choices {
				if item.MatchAnswer(choice) != (choice == tc.valid) {
					t.Errorf("unexpected verdict for choice %q", choice)
				}
			}
		})
	}
}
