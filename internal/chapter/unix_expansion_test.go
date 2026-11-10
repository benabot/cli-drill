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
			{"pwd", "ls", "ls-hidden", "cd-docs", "cd-parent", "cd-absolute", "mkdir-essais", "cat-spaces", "cd-home", "man-ls"},
			{"cd-previous", "ls-other-directory", "mkdir-parents", "cat-hidden-path", "ls-directory-entry", "glob-print-names", "glob-bracket-list", "quote-literal-question", "literal-tilde-directory", "command-origin"},
			{"cdpath-priority", "cd-symlink-parent", "glob-quoted-prefix", "glob-no-match-shells", "glob-leading-dot", "builtin-shadowed", "subshell-directory", "relative-from-peer", "dash-leading-name", "case-sensitive-volume"},
		}},
		{"11-flux-shell", LearningOrderGuided, [3][]string{
			{"redirect-output", "append-output"},
			{"redirect-input", "count-matching-lines", "redirect-errors"},
			{"merge-output-errors"},
		}},
		{"12-fichiers-permissions", LearningOrderIndependent, [3][]string{
			{"copy-spaced-path", "rename-file", "remove-file", "remove-empty-directory", "list-permissions", "touch-new-file", "copy-backup-name", "read-owner-column", "chmod-group-no-write", "recognize-link-entry"},
			{"owner-execute", "copy-directory-tree", "copy-confirm-overwrite", "move-many-sources", "copy-preserve-times", "link-relative-target", "readlink-stored-target", "chmod-two-classes", "interpret-octal-mode", "missing-directory-search"},
			{"remove-symlink-only", "umask-file-and-directory", "unlink-vs-write", "sticky-shared-directory", "setgid-inherit-group", "move-relative-symlink", "hardlink-after-replace", "chmod-conditional-search", "copy-exact-destination-gnu", "stat-mode-macos"},
		}},
		{"13-processus-environnement", LearningOrderIndependent, [3][]string{
			{"show-environment-variable", "set-shell-variable"},
			{"inspect-pid", "request-termination", "export-variable", "one-command-environment"},
			nil,
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
	} {
		t.Run(tc.key, func(t *testing.T) {
			item, ok := items[tc.key]
			if !ok {
				t.Fatal("exercise is missing from the embedded corpus")
			}
			if !item.MatchAnswer(tc.valid) || item.MatchAnswer(tc.wrong) {
				t.Fatalf("want to accept %q and reject %q", tc.valid, tc.wrong)
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
