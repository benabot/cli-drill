package app

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/benabot/cli-drill/data"
	"github.com/benabot/cli-drill/internal/chapter"
	"github.com/benabot/cli-drill/internal/progress"
	"github.com/benabot/cli-drill/internal/xdg"
)

func TestShortTrainingContinuesAndReviewsOnlyErrors(t *testing.T) {
	isolateChapterHome(t)
	path := filepath.Join(t.TempDir(), "progress.json")
	input := "pwd\nls\nls -a\ncd docs\nfaux\ncd /atelier\nr\ncd ..\nc\nmkdir essais\n"
	output, err := runSessionCLI(path, input, "train")
	for _, want := range []string{"Exercice 6/6", "Bilan du lot\nBonnes réponses : 5/6",
		"Bilan de la reprise\nBonnes réponses : 1/1", "Exercice 1/6", "Acquis cumulés (historique) : 7"} {
		if err != nil || !strings.Contains(output, want) {
			t.Fatalf("missing %q: %v\n%s", want, err, output)
		}
	}
	state, err := progress.Load(path)
	if err != nil || state.ChapterScores["00-premiers-pas"] != (progress.ChapterScore{Attempts: 8, Correct: 7}) ||
		!slices.Contains(state.CompletedExercises, "00-premiers-pas/cd-parent") || len(state.CompletedExercises) != 7 {
		t.Fatalf("unexpected progression: %#v (%v)", state, err)
	}
	resumed, err := runSessionCLI(path, "", "train")
	if err != nil || !strings.Contains(resumed, "Exercice 1/5") || strings.Contains(resumed, "sans changer de dossier") {
		t.Fatalf("unexpected restart: %v\n%s", err, resumed)
	}
}

func TestGitDefaultRespectsExistingPreferences(t *testing.T) {
	isolateChapterHome(t)
	path := xdg.PreferencesFile()
	out, err := runPreferencesCLI(t, "prefs")
	if err != nil || !strings.Contains(out, "git: enabled") {
		t.Fatalf("fresh defaults: %v\n%s", err, out)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("reading defaults wrote preferences: %v", err)
	}
	trainingPath := filepath.Join(t.TempDir(), "progress.json")
	out, err = runSessionCLI(trainingPath, "", "train", "10-workflows-dotfiles", "--all")
	if err != nil || !strings.Contains(out, "Exercice 1/10") {
		t.Fatalf("Git chapter not available with fresh defaults: %v\n%s", err, out)
	}
	if _, err := runPreferencesCLI(t, "prefs", "--disable", "git"); err != nil {
		t.Fatal(err)
	}
	if out, err := runPreferencesCLI(t, "prefs"); err != nil || !strings.Contains(out, "git: disabled") {
		t.Fatalf("explicit disable did not persist: %v\n%s", err, out)
	}
	for _, fixture := range []string{`{"version":1}`, `{"version":1,"enabled_tools":[]}`, `{"version":1,"enabled_tools":["fd"]}`} {
		before := writeAppTestFileContent(t, path, fixture)
		out, err = runPreferencesCLI(t, "prefs")
		if err != nil || !strings.Contains(out, "git: disabled") {
			t.Fatalf("saved preferences %s: %v\n%s", fixture, err, out)
		}
		assertPreferencesBytes(t, path, before)
	}
}

func TestStatsIgnoresHistoricalChaptersWithoutAttempts(t *testing.T) {
	isolateChapterHome(t)
	path := filepath.Join(t.TempDir(), "progress.json")
	before := writeAppTestFileContent(t, path, `{"version":1,"chapter_scores":{"old":{"attempts":0},"retired":{"attempts":4,"correct":2}}}`)
	out, err := runSessionCLI(path, "", "stats")
	if err != nil || !strings.Contains(out, "Chapters with attempts: 1") {
		t.Fatalf("stats: %v\n%s", err, out)
	}
	assertPreferencesBytes(t, path, before)
}

func TestAllAcquiredOffersReviewWithoutChangingProgress(t *testing.T) {
	isolateChapterHome(t)
	chapters, err := chapter.LoadFS(data.Chapters(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	ch := chapters[0]
	filter := chapter.SessionFilter{Difficulty: chapter.DifficultyBeginner}
	path := filepath.Join(t.TempDir(), "progress.json")
	state := progress.Default()
	for _, item := range filter.EligibleItems(ch, chapter.TextExercises) {
		state.MarkCompleted(ch.ID, item.ID, true)
	}
	if err := progress.Save(path, state); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ mode, want string }{
		{"learn", "tous les exercices de ce filtre sont acquis. Choisis Réviser avec --mode review"},
		{"review", "Exercice 1/6"},
	} {
		out, err := runSessionCLI(path, "", "train", ch.ID, "--mode", tc.mode)
		if err != nil || !strings.Contains(out, tc.want) {
			t.Fatalf("%s: %v\n%s", tc.mode, err, out)
		}
		assertPreferencesBytes(t, path, before)
	}
}
