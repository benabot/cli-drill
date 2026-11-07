package app

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/benabot/cli-drill/internal/progress"
)

func TestInvalidProgressStructureIsPreserved(t *testing.T) {
	for name, input := range map[string]string{
		"null version":        `{"version":null}`,
		"null score":          `{"version":1,"chapter_scores":{"retired":null}}`,
		"null attempts":       `{"chapter_scores":{"retired":{"attempts":null,"correct":3}}}`,
		"null correct":        `{"chapter_scores":{"retired":{"attempts":7,"correct":null}}}`,
		"null streak":         `{"streak":null}`,
		"null acquisition":    `{"completed_exercises":["retired/item",null]}`,
		"repeated version":    `{"version":null,"version":1}`,
		"case variant":        `{"VERSION":null}`,
		"string version":      `{"version":"1"}`,
		"object acquisitions": `{"completed_exercises":{}}`,
		"number acquisition":  `{"completed_exercises":[7]}`,
		"array score":         `{"chapter_scores":{"retired":[]}}`,
		"fractional attempts": `{"chapter_scores":{"retired":{"attempts":1.5}}}`,
	} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "progress.json")
			original := []byte(input)
			if err := os.WriteFile(path, original, 0o600); err != nil {
				t.Fatal(err)
			}
			_, err := progress.Load(path)
			assertRejectedProgress(t, path, original, err)
			for _, command := range []string{"stats", "train"} {
				err := executeProgressCommand(path, command)
				assertRejectedProgress(t, path, original, err)
			}
			// Protect callers that supply a default state after a failed load too.
			assertRejectedProgress(t, path, original, progress.Save(path, progress.Default()))
			entries, err := os.ReadDir(filepath.Dir(path))
			if err != nil || len(entries) != 1 {
				t.Fatalf("unexpected temporary files: %v (%v)", entries, err)
			}
		})
	}
}

func TestTrainingPreservesHistoricalProgressStructures(t *testing.T) {
	when := time.Date(2026, 6, 6, 12, 0, 0, 0, time.UTC)
	for _, tc := range []struct {
		name, input string
		want        progress.Progress
	}{
		{name: "missing fields", input: `{}`, want: progress.Default()},
		{
			name: "nullable containers", input: `{"completed_exercises":null,"chapter_scores":null,"last_session":null}`,
			want: progress.Progress{Version: 1, ChapterScores: map[string]progress.ChapterScore{}},
		},
		{
			name: "missing counters", input: `{"completed_exercises":["retired/item"],"chapter_scores":{"retired":{}},"streak":4}`,
			want: progress.Progress{Version: 1, CompletedExercises: []string{"retired/item"},
				ChapterScores: map[string]progress.ChapterScore{"retired": {}}, Streak: 4},
		},
		{
			name: "missing correct", input: `{"chapter_scores":{"retired":{"attempts":7}}}`,
			want: progress.Progress{Version: 1, CompletedExercises: []string{},
				ChapterScores: map[string]progress.ChapterScore{"retired": {Attempts: 7}}},
		},
		{
			name: "full v1", input: `{"version":1,"completed_exercises":["retired/item"],"chapter_scores":{"retired":{"attempts":7,"correct":3}},"last_session":"2026-06-06T12:00:00Z","streak":4}`,
			want: progress.Progress{Version: 1, CompletedExercises: []string{"retired/item"},
				ChapterScores: map[string]progress.ChapterScore{"retired": {Attempts: 7, Correct: 3}},
				LastSession:   &when, Streak: 4},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "progress.json")
			if err := os.WriteFile(path, []byte(tc.input), 0o600); err != nil {
				t.Fatal(err)
			}
			loaded, err := progress.Load(path)
			if err != nil || !reflect.DeepEqual(loaded, tc.want) {
				t.Fatalf("historical defaults: got %#v, want %#v (%v)", loaded, tc.want, err)
			}
			if err := executeProgressCommand(path, "stats"); err != nil {
				t.Fatal(err)
			}
			unchanged, err := os.ReadFile(path)
			if err != nil || string(unchanged) != tc.input {
				t.Fatalf("stats changed the original: %q (%v)", unchanged, err)
			}
			if err := executeProgressCommand(path, "train"); err != nil {
				t.Fatal(err)
			}
			got, err := progress.Load(path)
			if err != nil || got.LastSession == nil {
				t.Fatalf("reload after training: %#v (%v)", got, err)
			}
			want := tc.want
			want.MarkCompleted("embedded", "pwd", true)
			want.LastSession = got.LastSession
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("training lost historical results: got %#v, want %#v", got, want)
			}
		})
	}
}

func executeProgressCommand(path, command string) error {
	cmd := NewRootCommand(fstest.MapFS{
		"embedded.yaml": {Data: []byte(embeddedTestChapter("Storage fixture", "pwd"))},
	})
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetErr(&output)
	cmd.SetIn(strings.NewReader("pwd\n"))
	cmd.SetArgs([]string{"--progress", path, command})
	return cmd.Execute()
}

func assertRejectedProgress(t *testing.T, path string, original []byte, err error) {
	t.Helper()
	if err == nil || !strings.Contains(err.Error(), path) {
		t.Errorf("expected a contextual progress error: %v", err)
	}
	got, readErr := os.ReadFile(path)
	if readErr != nil || !bytes.Equal(got, original) {
		t.Fatalf("original progress changed: %q (%v)", got, readErr)
	}
}
