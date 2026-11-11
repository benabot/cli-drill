package app

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/benabot/cli-drill/internal/chapter"
	"github.com/benabot/cli-drill/internal/exercise"
	"github.com/benabot/cli-drill/internal/progress"
)

func TestCaptureIgnoresUnsupportedSequences(t *testing.T) {
	allowSimulatedCapture(t)
	ch := chapter.Chapter{ID: "keys", Items: []chapter.Item{
		{ID: "first", ExerciseType: exercise.TypeKeySequence, Answer: exercise.Answer{Primary: "Ctrl+B"}},
	}}
	input := "\x1b[D\x1bOP\x1b[[A\x1bx\x00é\x1b[200~h\ntext\x1b[201~\x02\n"
	path := filepath.Join(t.TempDir(), "progress.json")
	var out bytes.Buffer
	err := runTraining(Options{In: strings.NewReader(input), Out: &out, ProgressPath: path}, ch)
	state, loadErr := progress.Load(path)
	if err != nil || loadErr != nil || state.ChapterScores[ch.ID] != (progress.ChapterScore{Attempts: 1, Correct: 1}) {
		t.Fatalf("unsupported input affected the training: %v, %v, %#v\n%s", err, loadErr, state, &out)
	}
}

func TestCaptureLabelsPreserveCorrection(t *testing.T) {
	allowSimulatedCapture(t)
	ch := chapter.Chapter{ID: "keys", Items: []chapter.Item{
		{ID: "backspace", ExerciseType: exercise.TypeKeySequence, Answer: exercise.Answer{Primary: "Backspace", Accepted: []string{"\x7f"}}},
		{ID: "undo", ExerciseType: exercise.TypeKeySequence, Answer: exercise.Answer{Primary: "Ctrl+_", Accepted: []string{"\x1f"}}},
		{ID: "tab", ExerciseType: exercise.TypeKeySequence, Answer: exercise.Answer{Primary: "Tab", Accepted: []string{"Ctrl+I"}}},
	}}
	var out bytes.Buffer
	err := runTraining(Options{
		In: strings.NewReader("\x7f\n\x1f\n\t\n"), Out: &out,
		ProgressPath: filepath.Join(t.TempDir(), "progress.json"),
	}, ch)
	if err != nil || strings.Count(out.String(), "Correct.") != 3 {
		t.Fatalf("wrong correction: %v\n%s", err, &out)
	}
	for _, label := range []string{"DEL (Retour arrière)", "Ctrl+_", "Tab (Ctrl+I)"} {
		if !strings.Contains(out.String(), "Recu: "+label+"\n") {
			t.Errorf("missing capture label %q in %q", label, &out)
		}
	}
}

func TestMixedTrainingStopsOnTextEscapeOrEOF(t *testing.T) {
	allowSimulatedCapture(t)
	ch := chapter.Chapter{ID: "mixed", Items: []chapter.Item{
		{ID: "first", ExerciseType: exercise.TypeKeySequence, Answer: exercise.Answer{Primary: "Ctrl+B"}},
		{ID: "second", ExerciseType: exercise.TypeSimpleShellSim, Answer: exercise.Answer{Primary: "pwd"}},
	}}
	for name, lastInput := range map[string]string{"escape": "\x1b\nx\n", "eof": ""} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "progress.json")
			var out bytes.Buffer
			err := runTraining(Options{In: strings.NewReader("x\n" + lastInput), Out: &out, ProgressPath: path}, ch)
			if err != nil || strings.Contains(out.String(), "Enter review missed") ||
				strings.Contains(out.String(), "Chapitre termine.") ||
				!strings.Contains(out.String(), "Bonnes réponses : 0/1") {
				t.Fatalf("unexpected interrupted session: %v\n%s", err, &out)
			}
			state, loadErr := progress.Load(path)
			if loadErr != nil || state.ChapterScores[ch.ID] != (progress.ChapterScore{Attempts: 1}) {
				t.Fatalf("progress: %v, %#v", loadErr, state)
			}
		})
	}
}

func TestCaptureExitDoesNotCreateProgress(t *testing.T) {
	allowSimulatedCapture(t)
	ch := chapter.Chapter{ID: "keys", Items: []chapter.Item{
		{ID: "first", ExerciseType: exercise.TypeKeySequence, Answer: exercise.Answer{Primary: "Ctrl+B"}},
	}}
	for _, input := range []string{"\x1b", "\x03", "\x1b[1;2"} {
		path := filepath.Join(t.TempDir(), "progress.json")
		var out bytes.Buffer
		if err := runTraining(Options{In: strings.NewReader(input), Out: &out, ProgressPath: path}, ch); err != nil {
			t.Fatalf("input %q: %v", input, err)
		}
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("input %q created progress: %v", input, err)
		}
	}
}
