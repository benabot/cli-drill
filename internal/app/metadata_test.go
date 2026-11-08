package app

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/benabot/cli-drill/internal/progress"
)

func TestPedagogicalMetadataDoesNotChangeCLITrainingOrProgress(t *testing.T) {
	isolateChapterHome(t)
	legacy := embeddedTestChapter("Metadata lesson", "pwd")
	annotated := strings.Replace(legacy, "title: Metadata lesson\n", "title: Metadata lesson\nfamily: modern-cli-macos\nlearning_order: independent\n", 1)
	annotated = strings.Replace(annotated, "    type: workflow\n", "    type: workflow\n    difficulty: advanced\n    tools: [not-installed-example]\n", 1)
	const initial = `{"version":1,"completed_exercises":["retired/old"],"chapter_scores":{"retired":{"attempts":3,"correct":1}},"last_session":null,"streak":4}`
	var outputs []string
	var states []progress.Progress
	for _, lesson := range []string{legacy, annotated} {
		path := filepath.Join(t.TempDir(), "progress.json")
		if err := os.WriteFile(path, []byte(initial), 0o600); err != nil {
			t.Fatal(err)
		}
		run := func(input string, args ...string) string {
			t.Helper()
			cmd := NewRootCommand(fstest.MapFS{"lesson.yaml": {Data: []byte(lesson)}})
			var output bytes.Buffer
			cmd.SetIn(strings.NewReader(input))
			cmd.SetOut(&output)
			cmd.SetErr(&output)
			cmd.SetArgs(append([]string{"--progress", path}, args...))
			if err := cmd.Execute(); err != nil {
				t.Fatalf("Execute: %v\n%s", err, output.String())
			}
			return output.String()
		}
		if got := run("", "chapters"); got != "embedded\tMetadata lesson\t1 exercises\n" {
			t.Fatalf("metadata changed chapter listing: %q", got)
		}
		untouched, err := os.ReadFile(path)
		if err != nil || string(untouched) != initial {
			t.Fatalf("loading metadata wrote progress: %q (%v)", untouched, err)
		}
		outputs = append(outputs, run("PWD\n", "train", "embedded"), run("pwd\n", "train", "embedded"))
		saved, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		var fields map[string]json.RawMessage
		if err := json.Unmarshal(saved, &fields); err != nil || len(fields) != 5 {
			t.Fatalf("v1 schema changed: %s (%v)", saved, err)
		}
		for _, name := range []string{"version", "completed_exercises", "chapter_scores", "last_session", "streak"} {
			if _, ok := fields[name]; !ok {
				t.Fatalf("v1 field lost: %s", name)
			}
		}
		state, err := progress.Load(path)
		if err != nil {
			t.Fatal(err)
		}
		if state.Version != 1 || state.Streak != 4 || state.LastSession == nil ||
			!reflect.DeepEqual(state.CompletedExercises, []string{"retired/old", "embedded/pwd"}) ||
			!reflect.DeepEqual(state.ChapterScores, map[string]progress.ChapterScore{
				"retired": {Attempts: 3, Correct: 1}, "embedded": {Attempts: 2, Correct: 1},
			}) {
			t.Fatalf("progress changed: %#v", state)
		}
		state.LastSession = nil // Session timestamps differ; the v1 counters and keys must not.
		states = append(states, state)
	}
	if !reflect.DeepEqual(outputs[:2], outputs[2:]) || !reflect.DeepEqual(states[0], states[1]) {
		t.Fatalf("metadata changed training or progress: %q / %#v", outputs, states)
	}
	if strings.Contains(outputs[0], "> Correct.") || !strings.Contains(outputs[1], "> Correct.") {
		t.Fatalf("command case correction changed: %q", outputs)
	}
	if entries, err := os.ReadDir(legacyConfigRoot(t)); err == nil && len(entries) != 0 {
		t.Fatalf("configuration was created: %v", entries)
	} else if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
}
