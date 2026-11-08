package tui

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
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

func assertSavedRecording(t *testing.T, path string, attempts, correct int) progress.Progress {
	t.Helper()
	state, err := progress.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if state.ChapterScores["text"] != (progress.ChapterScore{Attempts: attempts, Correct: correct}) {
		t.Fatalf("recorded score: %#v, want %d/%d", state.ChapterScores, correct, attempts)
	}
	if state.ChapterScores["retired"] != (progress.ChapterScore{Attempts: 7, Correct: 3}) || state.Streak != 4 ||
		len(state.CompletedExercises) == 0 || state.CompletedExercises[0] != "retired/item" {
		t.Fatalf("history lost: %#v", state)
	}
	return state
}

func TestTUIRecordsEachTextSubmissionBeforeAdvancing(t *testing.T) {
	for _, entry := range []string{homeActionStartTraining, homeActionBrowse} {
		t.Run(entry, func(t *testing.T) {
			path := recordingPath(t)
			m := recordingModel(t, path, entry)
			correct := 0
			for i, answer := range []string{"pwd", "wrong", "pwd", "wrong"} {
				m.input.SetValue(answer)
				m = press(t, m, tea.KeyEnter)
				if answer == "pwd" {
					correct++
				}
				state := assertSavedRecording(t, path, i+1, correct)
				if !reflect.DeepEqual(m.progress, state) || m.summary.Total != i+1 || m.summary.Correct != correct {
					t.Fatalf("displayed and saved results diverge: %#v, %#v", m.summary, m.progress)
				}
				restarted := newModel(Options{Chapters: m.chapters, ProgressPath: path})
				if !reflect.DeepEqual(restarted.progress, state) {
					t.Fatal("restart lost submitted answers")
				}
				m = press(t, m, tea.KeyEnter)
				assertSavedRecording(t, path, i+1, correct)
			}
			assertScreen(t, m, screenTrainingSummary)
			m = press(t, m, tea.KeyRunes, 'g')
			assertScreen(t, m, screenHome)
			state := assertSavedRecording(t, path, 4, 2)
			if len(state.CompletedExercises) != 3 {
				t.Fatalf("unexpected acquisitions: %#v", state.CompletedExercises)
			}
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
				if cmd == nil {
					t.Fatal("Ctrl+C did not quit")
				}
			} else {
				assertScreen(t, m, screenChapterDetail)
				m = press(t, m, tea.KeyRunes, 'g')
				assertScreen(t, m, screenHome)
			}
			assertSavedRecording(t, path, tc.submitted, tc.submitted)
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
	assertSavedRecording(t, path, 1, 1)
	if m.itemIndex != 1 || m.summary.Total != 1 {
		t.Fatalf("repeated enter created extra submissions: index=%d, %#v", m.itemIndex, m.summary)
	}
	m = press(t, m, tea.KeyEsc)
	m = press(t, m, tea.KeyEnter)
	m.input.SetValue("pwd")
	m = press(t, m, tea.KeyEnter)
	state := assertSavedRecording(t, path, 2, 2)
	if len(state.CompletedExercises) != 2 {
		t.Fatalf("legitimate retry duplicated acquisition: %#v", state.CompletedExercises)
	}
}

func TestTUISaveFailureRetainsAnswerForRetry(t *testing.T) {
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
	assertSavedRecording(t, path, 1, 1)
	assertNotContains(t, m.View(), "Progress could not be saved")
	m = press(t, m, tea.KeyEnter)
	assertSavedRecording(t, path, 1, 1)
}

func TestTUIReadFailureIsVisibleAndNeverOverwritten(t *testing.T) {
	path := recordingPath(t)
	original := []byte(`{"version":null}`)
	if err := os.WriteFile(path, original, 0o600); err != nil {
		t.Fatal(err)
	}
	m := recordingModel(t, path, homeActionStartTraining)
	if !strings.Contains(m.View(), path) {
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
			m = press(t, m, tea.KeyEnter)
			m = press(t, m, tea.KeyEnter)
			updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
			if updated.(model).result.RunCLITrainChapterID != keySequenceChapter().ID || cmd == nil {
				t.Fatal("key route did not delegate to the CLI engine")
			}
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
		m.input.SetValue("pwd")
		m = press(t, m, tea.KeyEnter)
		m = press(t, m, tea.KeyEnter)
	}
	assertScreen(t, m, screenHome)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Fatal("normal exit did not quit")
	}
	restarted := newModel(Options{Chapters: []chapter.Chapter{recordingChapter()}, ProgressPath: path})
	saved := assertSavedRecording(t, path, 4, 4)
	if !reflect.DeepEqual(restarted.progress, saved) || len(saved.CompletedExercises) != 5 {
		t.Fatal("completed session lost after exit/restart")
	}
}
