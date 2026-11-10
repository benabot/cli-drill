package tui

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/benabot/cli-drill/data"
	"github.com/benabot/cli-drill/internal/chapter"
	"github.com/benabot/cli-drill/internal/exercise"
	"github.com/benabot/cli-drill/internal/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func embeddedChoiceChapters(t *testing.T) []chapter.Chapter {
	t.Helper()
	chapters, err := chapter.LoadFS(data.Chapters(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	var result []chapter.Chapter
	for _, ch := range chapters {
		for _, item := range ch.Items {
			if item.ExerciseType == exercise.TypeMultipleChoice {
				selected := ch
				selected.Items = []chapter.Item{item}
				result = append(result, selected)
			}
		}
	}
	if len(result) != 4 {
		t.Fatalf("expected all four embedded QCM, got %d", len(result))
	}
	return result
}

// Missing choices, ignored movement or a wrong label passed to correction must fail here.
func TestMultipleChoiceEveryEmbeddedOption(t *testing.T) {
	for _, ch := range embeddedChoiceChapters(t) {
		item := ch.Items[0]
		for index, choice := range item.Choices {
			t.Run(ch.ID+"/"+item.ID+"/"+choice, func(t *testing.T) {
				path := filepath.Join(t.TempDir(), "progress.json")
				m := startTestTraining(t, newModel(choiceTestOptions(ch, path)))
				resized, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
				m = resized.(model)
				for _, option := range item.Choices {
					assertContains(t, m.View(), option)
				}
				assertContains(t, m.View(), "> "+item.Choices[0])
				assertNotContains(t, m.View(), "Votre reponse")
				for range index {
					m = press(t, m, tea.KeyDown)
				}
				assertContains(t, m.View(), "> "+choice)
				if _, err := os.Stat(path); !os.IsNotExist(err) {
					t.Fatalf("navigation wrote progress: %v", err)
				}
				m = press(t, m, tea.KeyEnter)
				if width, height := lipgloss.Size(m.View()); width > 80 || height > 24 {
					t.Fatalf("choice feedback exceeds the terminal: %dx%d", width, height)
				}
				correct := 0
				// The primary answer is declared among choices in every existing QCM.
				if choice == item.Answer.Primary {
					correct = 1
				}
				if strings.HasPrefix(m.feedback, "Correct.") != (correct == 1) {
					t.Fatalf("wrong feedback for %q: %s", choice, m.feedback)
				}
				for _, line := range strings.Split(item.Explanation, "\n") {
					assertContains(t, m.View(), line)
				}
				if correct == 0 {
					assertContains(t, m.View(), item.Answer.Primary)
					if len(m.summary.Missed) != 1 || m.summary.Missed[0].Answered != choice {
						t.Fatalf("wrong missed answer: %#v", m.summary)
					}
				}
				for _, key := range []tea.KeyType{tea.KeyDown, tea.KeyUp} {
					m = press(t, m, key)
				}
				m = press(t, m, tea.KeyRunes, 'j') // Feedback must not change or resubmit the choice.
				m = press(t, m, tea.KeyEnter)
				state, err := progress.Load(path)
				if err != nil {
					t.Fatal(err)
				}
				if state.ChapterScores[ch.ID] != (progress.ChapterScore{Attempts: 1, Correct: correct}) {
					t.Fatalf("wrong persisted result: %#v", state)
				}
				want := []string(nil)
				if correct == 1 {
					want = []string{ch.ID + "/" + item.ID}
				}
				if !slices.Equal(state.CompletedExercises, want) {
					t.Fatalf("wrong acquired IDs: %q", state.CompletedExercises)
				}
				restarted := newModel(choiceTestOptions(ch, path))
				if restarted.progress.ChapterScores[ch.ID] != state.ChapterScores[ch.ID] {
					t.Fatal("restart lost result")
				}
			})
		}
	}
}

func TestMultipleChoiceNavigationAndInterruptions(t *testing.T) {
	ch := embeddedChoiceChapters(t)[0]
	for _, exit := range []tea.KeyType{tea.KeyEsc, tea.KeyCtrlC} {
		t.Run(exit.String(), func(t *testing.T) {
			path := recordingPath(t)
			m := startTestTraining(t, newModel(choiceTestOptions(ch, path)))
			m = press(t, m, tea.KeyUp) // Clamp at the first option.
			assertContains(t, m.View(), "> fd")
			m = press(t, m, tea.KeyRunes, 'j')
			assertContains(t, m.View(), "> bat")
			m = press(t, m, tea.KeyRunes, 'k')
			assertContains(t, m.View(), "> fd")
			for range 10 {
				m = press(t, m, tea.KeyDown)
			}
			assertContains(t, m.View(), "> vale")
			m = press(t, m, tea.KeyUp)
			assertContains(t, m.View(), "> glow")
			for _, r := range "ghq" {
				m = press(t, m, tea.KeyRunes, r)
			}
			updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}, Paste: true})
			m = updated.(model)
			assertContains(t, m.View(), "> glow")
			if m.showHelp || m.input.Value() != "" {
				t.Fatal("choice selection leaked into text/global shortcuts")
			}
			updated, cmd := m.Update(tea.KeyMsg{Type: exit})
			m = updated.(model)
			if exit == tea.KeyCtrlC {
				assertQuitCommand(t, cmd)
			} else {
				assertScreen(t, m, screenChapterDetail)
				m = press(t, m, tea.KeyEnter)
				assertContains(t, m.View(), "> fd")
			}
			got, err := os.ReadFile(path)
			if err != nil || string(got) != recordingHistory {
				t.Fatalf("interruption wrote progress: %s (%v)", got, err)
			}
		})
	}
}

func TestMultipleChoiceTransitionsResetSelectionAndPreserveText(t *testing.T) {
	ch := recordingChapter()
	ch.Items = []chapter.Item{ch.Items[2], ch.Items[2], ch.Items[3], ch.Items[0], ch.Items[1]}
	ch.Items[1].ID = "second-choice"
	m := startTestTraining(t, newTestModel(t, Options{Chapters: []chapter.Chapter{ch}}))
	m = press(t, m, tea.KeyDown)
	m = press(t, m, tea.KeyEnter)
	m = press(t, m, tea.KeyEnter)
	assertContains(t, m.View(), "> pwd")
	m = press(t, m, tea.KeyEnter)
	m = press(t, m, tea.KeyEnter)
	for range 3 {
		assertContains(t, m.View(), "Réponse")
		assertNotContains(t, m.View(), "Choix")
		m = press(t, m, tea.KeyRunes, []rune("pwd")...)
		m = press(t, m, tea.KeyEnter)
		if !strings.HasPrefix(m.feedback, "Correct.") {
			t.Fatalf("text exercise regressed: %s", m.feedback)
		}
		m = press(t, m, tea.KeyEnter)
	}
	assertScreen(t, m, screenTrainingSummary)
	if m.summary.Total != 5 || m.summary.Correct != 4 {
		t.Fatalf("wrong summary: %#v", m.summary)
	}
}

func TestKeyboardShortcutsDirectCaptureWithSingleChapter(t *testing.T) {
	path := recordingPath(t)
	m := newModel(Options{Chapters: []chapter.Chapter{textChapter(), keySequenceChapter(), textChapter()}, ProgressPath: path})
	m.home.Select(homeActionIndex(homeActionShortcuts))
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)
	if m.result.CaptureChapter.ID != keySequenceChapter().ID || m.chapterIndex != 1 {
		t.Fatalf("did not launch the single capture chapter directly: %#v", m.result)
	}
	assertQuitCommand(t, cmd)
	got, err := os.ReadFile(path)
	if err != nil || string(got) != recordingHistory {
		t.Fatalf("direct launch wrote progress: %s (%v)", got, err)
	}
}

func TestKeyboardShortcutsNoCaptureChapterDoesNotLaunch(t *testing.T) {
	m := newTestModel(t, Options{Chapters: []chapter.Chapter{textChapter()}})
	m.home.Select(homeActionIndex(homeActionShortcuts))
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)
	assertScreen(t, m, screenChapters)
	if m.result.CaptureChapter.ID != "" || cmd != nil || len(m.chapterList.Items()) != 0 {
		t.Fatal("empty shortcut browser launched training")
	}
}

func TestMultipleChoiceFeedbackInterruptionDoesNotRecordAgain(t *testing.T) {
	ch := embeddedChoiceChapters(t)[0]
	for _, exit := range []tea.KeyType{tea.KeyEsc, tea.KeyCtrlC} {
		for _, wrong := range []bool{false, true} {
			t.Run(exit.String()+"/wrong="+fmt.Sprint(wrong), func(t *testing.T) {
				path := filepath.Join(t.TempDir(), "progress.json")
				m := startTestTraining(t, newModel(choiceTestOptions(ch, path)))
				if wrong {
					m = press(t, m, tea.KeyDown)
				}
				m = press(t, m, tea.KeyEnter)
				before, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				updated, cmd := m.Update(tea.KeyMsg{Type: exit})
				if exit == tea.KeyCtrlC {
					assertQuitCommand(t, cmd)
				} else {
					assertScreen(t, updated.(model), screenChapterDetail)
				}
				after, err := os.ReadFile(path)
				if err != nil || !bytes.Equal(before, after) {
					t.Fatalf("feedback exit changed progress: %s (%v)", after, err)
				}
			})
		}
	}
}

func TestMultipleChoiceSaveFailureRetainsSelectionForRetry(t *testing.T) {
	ch := recordingChapter()
	ch.Items = ch.Items[2:3]
	path := recordingPath(t)
	target := path + ".target"
	if err := os.Rename(path, target); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, path); err != nil {
		t.Fatal(err)
	}
	m := startTestTraining(t, newModel(choiceTestOptions(ch, path)))
	m = press(t, m, tea.KeyDown)
	for range 2 {
		m = press(t, m, tea.KeyEnter)
		assertContains(t, m.View(), "Progress could not be saved")
		assertContains(t, m.View(), "> wrong")
		if m.feedback != "" || m.summary.Total != 0 || m.itemIndex != 0 {
			t.Fatal("failed save consumed the choice")
		}
		assertSavedRecording(t, target, 0, 0)
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(target, path); err != nil {
		t.Fatal(err)
	}
	m = press(t, m, tea.KeyEnter)
	assertSavedRecording(t, path, 1, 0)
	assertNotContains(t, m.View(), "Progress could not be saved")
	m = press(t, m, tea.KeyEnter)
	assertScreen(t, m, screenTrainingSummary)
	assertSavedRecording(t, path, 1, 0)
	m = press(t, m, tea.KeyEnter)
	assertScreen(t, m, screenTraining)
	assertContains(t, m.View(), "> pwd")
	assertSavedRecording(t, path, 1, 0)
}

func choiceTestOptions(ch chapter.Chapter, path string) Options {
	chapters := []chapter.Chapter{ch}
	return Options{
		Chapters:     chapters,
		ProgressPath: path,
		LoadPreferences: func() ([]string, chapter.Difficulty, error) {
			return chapter.AvailableTools(chapters), chapter.DifficultyBeginner, nil
		},
	}
}
