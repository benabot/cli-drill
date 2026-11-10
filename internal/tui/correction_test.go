package tui

import (
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/benabot/cli-drill/internal/chapter"
	"github.com/benabot/cli-drill/internal/exercise"
	"github.com/benabot/cli-drill/internal/progress"
	tea "github.com/charmbracelet/bubbletea"
)

func TestTUICorrectionWithHistoricalExercises(t *testing.T) {
	chapters, err := chapter.LoadFS(os.DirFS("../../testdata/legacy-chapters"), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	// Same expected outcomes as the CLI test, including the textual shortcut.
	// Actual key-sequence exercises remain delegated to the CLI capture engine.
	for _, tc := range []struct {
		chapter, item, input string
		correct              bool
	}{
		{"02-navigation-shell", "fzf-pick", " FZF ", true},
		{"02-navigation-shell", "fzf-pick", "\tfzf\t", true},
		{"02-navigation-shell", "zoxide-jump", "z", true},
		{"02-navigation-shell", "fd-find-dir", "fd", true},
		{"02-navigation-shell", "fd-find-dir", "bat", false},
		{"02-navigation-shell", "fd-fzf-dir", "fd --type d | fzf", true},
		{"02-navigation-shell", "fd-fzf-dir", "fd   -t d|fzf", true},
		{"02-navigation-shell", "fd-fzf-dir", "fd -t D | fzf", false},
		{"07-lecture-preview", "bat-readme", "bat readme.md", false},
		{"07-lecture-preview", "bat-readme", " bat README.md ", true},
		{"08-micro", "micro-save", "^s", true},
		{"08-micro", "micro-save", "s", false},
		{"05-outils-quotidiens", "atuin", "control-r", true},
		{"06-recherche-fichiers-contenu", "rg-fzf-todo", "\trg \"TODO\" | fzf\t", true},
	} {
		t.Run(tc.chapter+"/"+tc.item+"/"+tc.input, func(t *testing.T) {
			var selected chapter.Chapter
			for _, ch := range chapters {
				for _, item := range ch.Items {
					if ch.ID == tc.chapter && item.ID == tc.item {
						selected = ch
						selected.Items = []chapter.Item{item}
					}
				}
			}
			if selected.ID == "" {
				t.Fatal("embedded exercise not found")
			}
			path := recordingPath(t)
			m := newModel(Options{Chapters: []chapter.Chapter{selected}, ProgressPath: path})
			m.home.Select(homeActionIndex(homeActionBrowse))
			for range 3 {
				m = press(t, m, tea.KeyEnter)
			}
			assertScreen(t, m, screenTraining)
			if selected.Items[0].ExerciseType == exercise.TypeMultipleChoice {
				index := slices.Index(selected.Items[0].Choices, tc.input)
				if index < 0 {
					t.Fatalf("input is not a declared choice: %q", tc.input)
				}
				for range index {
					m = press(t, m, tea.KeyDown)
				}
			} else {
				m = press(t, m, tea.KeyRunes, []rune(tc.input)...)
			}
			m = press(t, m, tea.KeyEnter)
			correct := 0
			acquired := []string{"retired/item"}
			if tc.correct {
				correct = 1
				acquired = append(acquired, tc.chapter+"/"+tc.item)
			}
			if strings.HasPrefix(m.feedback, "Correct.") != tc.correct || m.summary.Total != 1 || m.summary.Correct != correct {
				t.Fatalf("wrong correction or summary: %s / %#v", m.feedback, m.summary)
			}
			for range 2 {
				state, err := progress.Load(path)
				if err != nil {
					t.Fatal(err)
				}
				if state.ChapterScores[tc.chapter] != (progress.ChapterScore{Attempts: 1, Correct: correct}) ||
					state.ChapterScores["retired"] != (progress.ChapterScore{Attempts: 7, Correct: 3}) ||
					!slices.Equal(state.CompletedExercises, acquired) {
					t.Fatalf("wrong persisted score or acquisitions: %#v", state)
				}
				restarted := newModel(Options{Chapters: []chapter.Chapter{selected}, ProgressPath: path})
				if !reflect.DeepEqual(restarted.progress, state) {
					t.Fatal("restart lost correction results")
				}
				m = press(t, m, tea.KeyEnter) // Feedback validation must not record again.
			}
		})
	}
}

func TestTUIRejectsMultilinePasteBeforeSanitizing(t *testing.T) {
	for _, tc := range []struct {
		name, text string
		clipboard  bool
	}{
		{"terminal LF", "fd -t d\n| fzf", false},
		{"terminal CR", "fd -t d\r| fzf", false},
		{"terminal CRLF", "fd -t d\r\n| fzf", false},
		{"clipboard LF", "fd -t d\n| fzf", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := recordingPath(t)
			ch := recordingChapter()
			ch.Items = ch.Items[3:]
			ch.Items[0].Answer.Primary = "fd -t d | fzf"
			m := startTestTraining(t, newModel(Options{Chapters: []chapter.Chapter{ch}, ProgressPath: path}))
			m.input.SetValue("fd -t d | fzf")
			var message tea.Msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tc.text), Paste: true}
			if tc.clipboard {
				message = trainingPasteMsg{text: tc.text}
			}
			updated, _ := m.Update(message)
			m = updated.(model)
			if m.input.Value() != "fd -t d | fzf" || m.input.Err == nil || !strings.Contains(m.View(), "single line") {
				t.Fatalf("paste changed input or warning is absent: %q / %v", m.input.Value(), m.input.Err)
			}
			m = press(t, m, tea.KeyEnter)
			if m.summary.Total != 0 || m.feedback != "" {
				t.Fatal("rejected paste produced a submission")
			}
			if got, err := os.ReadFile(path); err != nil || string(got) != recordingHistory {
				t.Fatalf("rejected paste changed original progress: %s (%v)", got, err)
			}
			message = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" "), Paste: true}
			if tc.clipboard {
				message = trainingPasteMsg{text: " "}
			}
			updated, _ = m.Update(message)
			m = updated.(model)
			m = press(t, m, tea.KeyEnter)
			assertSavedRecording(t, path, 1, 1, "text/simple-shell-sim")
		})
	}
}
