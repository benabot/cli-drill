package tui

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/benabot/cli-drill/internal/chapter"
	"github.com/benabot/cli-drill/internal/exercise"
	"github.com/benabot/cli-drill/internal/progress"
	tea "github.com/charmbracelet/bubbletea"
)

const recordingHistory = `{"version":1,"completed_exercises":["retired/item"],"chapter_scores":{"retired":{"attempts":7,"correct":3}},"last_session":"2026-06-06T12:00:00Z","streak":4}`

func recordingChapter() chapter.Chapter {
	result := chapter.Chapter{ID: "text", Title: "Text"}
	for _, kind := range []exercise.Type{exercise.TypeFreeAnswer, exercise.TypeScenario, exercise.TypeMultipleChoice, exercise.TypeSimpleShellSim} {
		result.Items = append(result.Items, chapter.Item{ID: string(kind), ExerciseType: kind,
			Prompt: "Command?", Choices: []string{"pwd", "wrong"}, Answer: exercise.Answer{Primary: "pwd"}})
	}
	return result
}

func recordingModel(t *testing.T, path, entry string) model {
	t.Helper()
	m := newModel(Options{Chapters: []chapter.Chapter{recordingChapter()}, ProgressPath: path})
	m.home.Select(homeActionIndex(entry))
	m = press(t, m, tea.KeyEnter)
	if entry == homeActionBrowse {
		m = press(t, m, tea.KeyEnter)
		m = press(t, m, tea.KeyEnter)
	}
	assertScreen(t, m, screenTraining)
	return m
}

func recordingPath(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "progress.json")
	if err := os.WriteFile(path, []byte(recordingHistory), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func assertSavedRecording(t *testing.T, path string, attempts, correct int, acquired ...string) progress.Progress {
	t.Helper()
	state, err := progress.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if state.ChapterScores["text"] != (progress.ChapterScore{Attempts: attempts, Correct: correct}) {
		t.Fatalf("recorded score: %#v, want %d/%d", state.ChapterScores, correct, attempts)
	}
	if state.ChapterScores["retired"] != (progress.ChapterScore{Attempts: 7, Correct: 3}) || state.Streak != 4 {
		t.Fatalf("history lost: %#v", state)
	}
	want := append([]string{"retired/item"}, acquired...)
	if !slices.Equal(state.CompletedExercises, want) {
		t.Fatalf("acquisitions: got %q, want %q", state.CompletedExercises, want)
	}
	return state
}

func assertQuitCommand(t *testing.T, cmd tea.Cmd) {
	t.Helper()
	if cmd == nil {
		t.Fatal("expected a quit command")
	}
	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Fatalf("expected tea.QuitMsg, got %T", msg)
	}
}

func TestTUIRecordsEachTextSubmissionBeforeAdvancing(t *testing.T) {
	for _, entry := range []string{homeActionStartTraining, homeActionBrowse} {
		t.Run(entry, func(t *testing.T) {
			path := recordingPath(t)
			m := recordingModel(t, path, entry)
			correct := 0
			var acquired []string
			for i, step := range []struct{ answer, acquired string }{
				{"pwd", "text/free-answer"}, {"wrong", ""},
				{"pwd", "text/multiple-choice"}, {"wrong", ""},
			} {
				if m.currentItem().ExerciseType == exercise.TypeMultipleChoice {
					for range slices.Index(m.currentItem().Choices, step.answer) {
						m = press(t, m, tea.KeyDown)
					}
				} else {
					m.input.SetValue(step.answer)
				}
				m = press(t, m, tea.KeyEnter)
				if step.acquired != "" {
					correct++
					acquired = append(acquired, step.acquired)
				}
				state := assertSavedRecording(t, path, i+1, correct, acquired...)
				if !reflect.DeepEqual(m.progress, state) || m.summary.Total != i+1 || m.summary.Correct != correct {
					t.Fatalf("displayed and saved results diverge: %#v, %#v", m.summary, m.progress)
				}
				restarted := newModel(Options{Chapters: m.chapters, ProgressPath: path})
				if !reflect.DeepEqual(restarted.progress, state) {
					t.Fatal("restart lost submitted answers")
				}
				m = press(t, m, tea.KeyEnter)
				assertSavedRecording(t, path, i+1, correct, acquired...)
			}
			assertScreen(t, m, screenTrainingSummary)
			m = press(t, m, tea.KeyRunes, 'g')
			assertScreen(t, m, screenHome)
			assertSavedRecording(t, path, 4, 2, acquired...)
		})
	}
}

func TestTUIInterruptionsKeepOnlySubmittedAnswers(t *testing.T) {
	for _, tc := range []struct {
		name      string
		submitted int
		key       tea.KeyType
	}{
		{"escape before submission", 0, tea.KeyEsc},
		{"ctrl c before submission", 0, tea.KeyCtrlC},
		{"escape on feedback", 1, tea.KeyEsc},
		{"ctrl c on feedback", 1, tea.KeyCtrlC},
		{"escape after several answers", 2, tea.KeyEsc},
		{"ctrl c after several answers", 2, tea.KeyCtrlC},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := recordingPath(t)
			m := recordingModel(t, path, homeActionStartTraining)
			for i := 0; i < tc.submitted; i++ {
				if i > 0 {
					m = press(t, m, tea.KeyEnter)
				}
				m.input.SetValue("pwd")
				m = press(t, m, tea.KeyEnter)
			}
			if tc.submitted == 0 {
				m.input.SetValue("unfinished")
			}
			updated, cmd := m.Update(tea.KeyMsg{Type: tc.key})
			m = updated.(model)
			if tc.key == tea.KeyCtrlC {
				assertQuitCommand(t, cmd)
			} else {
				assertScreen(t, m, screenChapterDetail)
				m = press(t, m, tea.KeyRunes, 'g')
				assertScreen(t, m, screenHome)
			}
			acquired := []string{"text/free-answer", "text/scenario"}[:tc.submitted]
			assertSavedRecording(t, path, tc.submitted, tc.submitted, acquired...)
			if tc.submitted == 0 {
				got, err := os.ReadFile(path)
				if err != nil || string(got) != recordingHistory {
					t.Fatalf("abandon changed history: %q (%v)", got, err)
				}
			}
		})
	}
}

func TestTUIRepeatedEnterAndNewSession(t *testing.T) {
	path := recordingPath(t)
	m := recordingModel(t, path, homeActionStartTraining)
	m.input.SetValue("pwd")
	m = press(t, m, tea.KeyEnter)
	for range 8 {
		m = press(t, m, tea.KeyEnter)
	}
	assertSavedRecording(t, path, 1, 1, "text/free-answer")
	if m.itemIndex != 1 || m.summary.Total != 1 {
		t.Fatalf("repeated enter created extra submissions: index=%d, %#v", m.itemIndex, m.summary)
	}
	m = press(t, m, tea.KeyEsc)
	m = press(t, m, tea.KeyEnter)
	m.input.SetValue("pwd")
	m = press(t, m, tea.KeyEnter)
	assertSavedRecording(t, path, 2, 2, "text/free-answer")
}

func recordingSaveFailure(t *testing.T) (model, string, string) {
	t.Helper()
	path := recordingPath(t)
	target := path + ".target"
	if err := os.Rename(path, target); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, path); err != nil {
		t.Fatal(err)
	}
	m := recordingModel(t, path, homeActionStartTraining)
	m.input.SetValue("pwd")
	return m, path, target
}

func TestTUISaveFailureRetainsAnswerForRetry(t *testing.T) {
	m, path, target := recordingSaveFailure(t)
	for range 2 {
		m = press(t, m, tea.KeyEnter)
		assertContains(t, m.View(), "Progress could not be saved")
		if m.feedback != "" || m.summary.Total != 0 || m.itemIndex != 0 || m.input.Value() != "pwd" {
			t.Fatal("failed save advanced or consumed the answer")
		}
		got, err := os.ReadFile(target)
		if err != nil || string(got) != recordingHistory {
			t.Fatalf("failed save changed history: %q (%v)", got, err)
		}
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(target, path); err != nil {
		t.Fatal(err)
	}
	m = press(t, m, tea.KeyEnter)
	assertSavedRecording(t, path, 1, 1, "text/free-answer")
	assertNotContains(t, m.View(), "Progress could not be saved")
	m = press(t, m, tea.KeyEnter)
	assertSavedRecording(t, path, 1, 1, "text/free-answer")
}

func TestTUISaveFailureCanQuitWithoutPersisting(t *testing.T) {
	m, path, target := recordingSaveFailure(t)
	m = press(t, m, tea.KeyEnter)
	assertScreen(t, m, screenTraining)
	assertContains(t, m.View(), "Progress could not be saved")
	if m.feedback != "" || m.summary.Total != 0 || m.summary.Correct != 0 || len(m.summary.Missed) != 0 ||
		m.itemIndex != 0 || m.input.Value() != "pwd" {
		t.Fatal("failed save advanced or consumed the answer")
	}
	saved := assertSavedRecording(t, path, 0, 0)
	if !reflect.DeepEqual(m.progress, saved) {
		t.Fatal("unsaved answer is presented as persisted")
	}
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	assertQuitCommand(t, cmd)
	final := updated.(model)
	if !reflect.DeepEqual(final.progress, saved) || final.summary.Total != 0 ||
		final.summary.Correct != 0 || len(final.summary.Missed) != 0 {
		t.Fatal("quitting counted the unsaved answer")
	}
	if got, err := os.ReadFile(target); err != nil || string(got) != recordingHistory {
		t.Fatalf("quitting after failure changed history: %q (%v)", got, err)
	}
	restarted := newModel(Options{Chapters: m.chapters, ProgressPath: path})
	if !reflect.DeepEqual(restarted.progress, saved) {
		t.Fatal("unsaved answer appeared after restart")
	}
}

func TestTUIReadFailureIsVisibleAndNeverOverwritten(t *testing.T) {
	path := recordingPath(t)
	original := []byte(`{"version":null}`)
	if err := os.WriteFile(path, original, 0o600); err != nil {
		t.Fatal(err)
	}
	m := recordingModel(t, path, homeActionStartTraining)
	if view := m.View(); !strings.Contains(view, "Progress could not be loaded") ||
		!strings.Contains(strings.Join(strings.Fields(view), ""), path) {
		t.Fatal("initial read error is hidden while training")
	}
	m.input.SetValue("pwd")
	m = press(t, m, tea.KeyEnter)
	got, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(got, original) || m.summary.Total != 0 {
		t.Fatalf("invalid progress was overwritten or counted: %q (%v)", got, err)
	}
}

func TestTUIKeyRoutesDoNotRecordNavigation(t *testing.T) {
	for _, entry := range []string{homeActionBrowse, homeActionShortcuts} {
		t.Run(entry, func(t *testing.T) {
			path := recordingPath(t)
			m := newModel(Options{Chapters: []chapter.Chapter{keySequenceChapter()}, ProgressPath: path})
			m.home.Select(homeActionIndex(entry))
			if entry == homeActionBrowse {
				m = press(t, m, tea.KeyEnter)
				m = press(t, m, tea.KeyEnter)
			}
			updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
			if updated.(model).result.CaptureChapter.ID != keySequenceChapter().ID {
				t.Fatal("key route did not delegate to the CLI engine")
			}
			assertQuitCommand(t, cmd)
			got, err := os.ReadFile(path)
			if err != nil || string(got) != recordingHistory {
				t.Fatalf("navigation rewrote progress: %q (%v)", got, err)
			}
		})
	}
}

func TestTUIPerfectSessionSurvivesQuitAndRestart(t *testing.T) {
	path := recordingPath(t)
	m := recordingModel(t, path, homeActionStartTraining)
	for range 4 {
		if m.currentItem().ExerciseType != exercise.TypeMultipleChoice {
			m.input.SetValue("pwd")
		}
		m = press(t, m, tea.KeyEnter)
		m = press(t, m, tea.KeyEnter)
	}
	assertScreen(t, m, screenHome)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	assertQuitCommand(t, cmd)
	restarted := newModel(Options{Chapters: []chapter.Chapter{recordingChapter()}, ProgressPath: path})
	saved := assertSavedRecording(t, path, 4, 4,
		"text/free-answer", "text/scenario", "text/multiple-choice", "text/simple-shell-sim")
	if !reflect.DeepEqual(restarted.progress, saved) {
		t.Fatal("completed session lost after exit/restart")
	}
}
