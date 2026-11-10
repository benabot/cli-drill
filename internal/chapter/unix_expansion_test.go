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
		{"11-flux-shell", LearningOrderGuided, [3][]string{
			{"redirect-output", "append-output"},
			{"redirect-input", "count-matching-lines", "redirect-errors"},
			{"merge-output-errors"},
		}},
		{"12-fichiers-permissions", LearningOrderIndependent, [3][]string{
			{"copy-spaced-path", "rename-file", "remove-file", "remove-empty-directory", "list-permissions"},
			{"owner-execute"},
			nil,
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
				if len(item.Tools) != 0 || item.ExerciseType != exercise.TypeSimpleShellSim {
					t.Fatalf("%s: want a simulated command without modules", item.ID)
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
