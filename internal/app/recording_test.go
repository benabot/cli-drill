package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/benabot/cli-drill/internal/progress"
	"github.com/benabot/cli-drill/internal/tui"
)

const recordingHistory = `{"version":1,"completed_exercises":["retired/item"],"chapter_scores":{"retired":{"attempts":7,"correct":3}},"last_session":"2026-06-06T12:00:00Z","streak":4}`

func recordingFS() fstest.MapFS {
	var text strings.Builder
	text.WriteString("id: text\ntitle: Text\nitems:\n")
	for _, kind := range []string{"free-answer", "scenario", "multiple-choice", "simple-shell-sim"} {
		text.WriteString("  - id: " + kind + "\n    type: concept\n    exercise_type: " + kind + "\n    prompt: Command?\n    choices: [pwd, wrong]\n    answer:\n      primary: pwd\n")
	}
	return fstest.MapFS{
		"01-text.yaml": {Data: []byte(text.String())},
		"02-keys.yaml": {Data: []byte("id: keys\ntitle: Keys\nitems:\n  - id: ctrl-a\n    type: shortcut\n    exercise_type: key-sequence\n    prompt: Press Ctrl+A\n    answer:\n      primary: Ctrl+A\n")},
	}
}

func runRecordingCLI(path string, input io.Reader, args ...string) (string, error) {
	cmd := NewRootCommand(recordingFS())
	var output bytes.Buffer
	cmd.SetIn(input)
	cmd.SetOut(&output)
	cmd.SetErr(&output)
	cmd.SetArgs(append([]string{"--progress", path}, args...))
	err := cmd.Execute()
	return output.String(), err
}

func assertRecordingScore(t *testing.T, path, chapterID string, attempts, correct int) progress.Progress {
	t.Helper()
	state, err := progress.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := state.ChapterScores[chapterID]; got != (progress.ChapterScore{Attempts: attempts, Correct: correct}) {
		t.Fatalf("%s: got %#v, want %d attempts/%d correct", chapterID, got, attempts, correct)
	}
	if state.ChapterScores["retired"] != (progress.ChapterScore{Attempts: 7, Correct: 3}) || state.Streak != 4 ||
		len(state.CompletedExercises) == 0 || state.CompletedExercises[0] != "retired/item" {
		t.Fatalf("historical results changed: %#v", state)
	}
	return state
}

func TestCLIRecordsOnlySubmittedText(t *testing.T) {
	for _, tc := range []struct {
		name, input       string
		attempts, correct int
	}{
		{"empty EOF", "", 0, 0},
		{"blank lines", "\n \n", 0, 0},
		{"nonempty EOF", "pwd", 1, 1},
		{"incorrect EOF", "wrong", 1, 0},
		{"partial session", "pwd\n", 1, 1},
		{"all text types", "pwd\nwrong\npwd\nwrong\n", 4, 2},
		{"repeated enter", "\n\npwd\n\n\n", 1, 1},
		{"escape before answer", "\x1b", 0, 0},
		{"ctrl c before answer", "\x03", 0, 0},
		{"escape after answer", "pwd\nunfinished\x1b\n", 1, 1},
		{"ctrl c after answer", "pwd\nunfinished\x03", 1, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "progress.json")
			if err := os.WriteFile(path, []byte(recordingHistory), 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := runRecordingCLI(path, strings.NewReader(tc.input), "train"); err != nil {
				t.Fatal(err)
			}
			state := assertRecordingScore(t, path, "text", tc.attempts, tc.correct)
			if len(state.CompletedExercises) != 1+tc.correct {
				t.Fatalf("unexpected acquisitions: %#v", state.CompletedExercises)
			}
			if tc.attempts == 0 {
				got, err := os.ReadFile(path)
				if err != nil || string(got) != recordingHistory {
					t.Fatalf("consultation/abandon rewrote history: %q (%v)", got, err)
				}
			}
		})
	}
}

type recordingReader func([]byte) (int, error)

func (read recordingReader) Read(data []byte) (int, error) { return read(data) }

func TestCLIRecordsBeforeReadingAgain(t *testing.T) {
	for _, chapterID := range []string{"text", "keys"} {
		t.Run(chapterID, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "progress.json")
			if err := os.WriteFile(path, []byte(recordingHistory), 0o600); err != nil {
				t.Fatal(err)
			}
			interrupted := errors.New("input interrupted")
			reads := 0
			input := recordingReader(func(data []byte) (int, error) {
				reads++
				if reads == 1 {
					if chapterID == "keys" {
						return copy(data, "\x01"), nil
					}
					return copy(data, "pwd\n"), nil
				}
				assertRecordingScore(t, path, chapterID, 1, 1)
				return 0, interrupted
			})
			if _, err := runRecordingCLI(path, input, "train", chapterID); !errors.Is(err, interrupted) {
				t.Fatalf("lost input error: %v", err)
			}
			assertRecordingScore(t, path, chapterID, 1, 1)
		})
	}
}

func TestCLIKeyRecordingDoesNotCountControls(t *testing.T) {
	for _, tc := range []struct {
		name, input       string
		attempts, correct int
	}{
		{"help and enter", "h\n\x1b", 0, 0},
		{"validated then escape", "\x01\x1b", 1, 1},
		{"retry and solution", "Xsr\x01\x1b", 2, 1},
		{"review missed", "X\n\n\x01\x1b", 2, 1},
		{"repeated feedback enter", "\x01\n\n\n", 1, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "progress.json")
			if err := os.WriteFile(path, []byte(recordingHistory), 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := runRecordingCLI(path, strings.NewReader(tc.input), "train", "keys"); err != nil {
				t.Fatal(err)
			}
			state := assertRecordingScore(t, path, "keys", tc.attempts, tc.correct)
			if len(state.CompletedExercises) != 1+tc.correct {
				t.Fatalf("duplicate acquisition: %#v", state.CompletedExercises)
			}
		})
	}
}

func TestRecordingKeepsV1StatsAndReset(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")
	if err := os.WriteFile(path, []byte(recordingHistory), 0o600); err != nil {
		t.Fatal(err)
	}
	for range 2 {
		if _, err := runRecordingCLI(path, strings.NewReader("pwd\npwd\npwd\npwd\n"), "train", "text"); err != nil {
			t.Fatal(err)
		}
	}
	state := assertRecordingScore(t, path, "text", 8, 8)
	if len(state.CompletedExercises) != 5 {
		t.Fatalf("repeated session duplicated acquisitions: %#v", state.CompletedExercises)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil || len(fields) != 5 || string(fields["version"]) != "1" {
		t.Fatalf("v1 schema changed: %s (%v)", data, err)
	}
	stats, err := runRecordingCLI(path, strings.NewReader(""), "stats")
	if err != nil || stats != "Completed exercises: 5\nChapters with attempts: 2\nStreak: 4\n" {
		t.Fatalf("stats: %q (%v)", stats, err)
	}
	if _, err := runRecordingCLI(path, strings.NewReader(""), "reset"); err == nil {
		t.Fatal("reset must require confirmation")
	}
	if got, err := os.ReadFile(path); err != nil || !bytes.Equal(got, data) {
		t.Fatalf("unconfirmed reset changed progress: %v", err)
	}
	if _, err := runRecordingCLI(path, strings.NewReader(""), "reset", "--yes"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("confirmed reset: %v", err)
	}
}

func TestCLIPropagatesSaveFailure(t *testing.T) {
	for _, tc := range []struct{ chapter, input string }{{"text", "pwd\n"}, {"keys", "\x01"}} {
		t.Run(tc.chapter, func(t *testing.T) {
			dir := t.TempDir()
			path, target := filepath.Join(dir, "progress.json"), filepath.Join(dir, "target.json")
			if err := os.WriteFile(target, []byte(recordingHistory), 0o600); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(target, path); err != nil {
				t.Fatal(err)
			}
			output, err := runRecordingCLI(path, strings.NewReader(tc.input), "train", tc.chapter)
			if err == nil || !strings.Contains(err.Error(), "save progress") {
				t.Fatalf("lost save error: %v", err)
			}
			if strings.Contains(output, "Correct.") {
				t.Fatal("success feedback appeared before a failed save")
			}
			if got, err := os.ReadFile(target); err != nil || string(got) != recordingHistory {
				t.Fatalf("save failure changed target: %q (%v)", got, err)
			}
		})
	}
}

func TestCLIAbandonDoesNotCreateProgress(t *testing.T) {
	for _, tc := range []struct{ chapter, input string }{
		{"text", ""}, {"text", "\n \n"}, {"text", "\x1b"}, {"text", "\x03"},
		{"keys", ""}, {"keys", "h\n\x1b"},
	} {
		t.Run(tc.chapter+"/"+tc.input, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "new", "progress.json")
			if _, err := runRecordingCLI(path, strings.NewReader(tc.input), "train", tc.chapter); err != nil {
				t.Fatal(err)
			}
			if _, err := os.Stat(filepath.Dir(path)); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("abandon created a progress directory: %v", err)
			}
		})
	}
}

func TestTUIKeyIntentPersistsBeforeReturningToTUI(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")
	if err := os.WriteFile(path, []byte(recordingHistory), 0o600); err != nil {
		t.Fatal(err)
	}
	calls := 0
	restore := stubTUIRunner(func(opts tui.Options) (tui.Result, error) {
		calls++
		if opts.ProgressPath != path {
			t.Fatalf("TUI progress path changed: %s", opts.ProgressPath)
		}
		if calls == 1 {
			return tui.Result{RunCLITrainChapterID: "keys"}, nil
		}
		assertRecordingScore(t, path, "keys", 1, 1)
		return tui.Result{}, nil
	})
	defer restore()
	if _, err := runRecordingCLI(path, strings.NewReader("\x01\x1b")); err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Fatalf("expected a return to TUI, got %d launches", calls)
	}
}
