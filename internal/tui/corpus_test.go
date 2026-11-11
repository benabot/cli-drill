package tui

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/benabot/cli-drill/data"
	"github.com/benabot/cli-drill/internal/chapter"
	"github.com/benabot/cli-drill/internal/progress"
	tea "github.com/charmbracelet/bubbletea"
)

func TestStartTrainingUsesV2FoundationsAndKeepsDirectCapture(t *testing.T) {
	chapters, err := chapter.LoadFS(data.Chapters(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "progress.json")
	m := newModel(Options{Chapters: chapters, ProgressPath: path})
	m.home.Select(homeActionIndex(homeActionStartTraining))
	m = press(t, m, tea.KeyEnter)
	assertScreen(t, m, screenTraining)
	if m.chapters[m.chapterIndex].ID != "00-premiers-pas" {
		t.Fatal("Start training did not select Unix foundations")
	}
	filter := chapter.SessionFilter{Difficulty: chapter.DifficultyBeginner}
	for _, item := range filter.EligibleItems(chapters[0], chapter.TextExercises) {
		for _, line := range strings.Split(item.Prompt, "\n") {
			assertContains(t, m.View(), line)
		}
		m = press(t, m, tea.KeyRunes, []rune(item.Answer.Primary)...)
		m = press(t, m, tea.KeyEnter)
		assertContains(t, m.View(), "Correct.")
		m = press(t, m, tea.KeyEnter)
	}
	assertScreen(t, m, screenHome)
	state, err := progress.Load(path)
	if err != nil || state.ChapterScores["00-premiers-pas"] != (progress.ChapterScore{Attempts: 12, Correct: 12}) {
		t.Fatalf("TUI progress: %#v (%v)", state, err)
	}
	m.home.Select(homeActionIndex(homeActionShortcuts))
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)
	if m.result.CaptureChapter.ID != "01-raccourcis-terminal" {
		t.Fatalf("capture delegation lost: %#v", m.result)
	}
	assertQuitCommand(t, cmd)
}

func TestWave3ChapterBrowserAndShortcutScopes(t *testing.T) {
	chapters, err := chapter.LoadFS(data.Chapters(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	const id = "01-raccourcis-terminal"
	for _, tc := range []struct {
		name, action string
		level        chapter.Difficulty
		wantKeys     int
		wantText     int
	}{
		{"mixed through chapter browser", homeActionBrowse, chapter.DifficultyIntermediate, 6, 5},
		{"capture only through shortcuts", homeActionShortcuts, chapter.DifficultyIntermediate, 6, 0},
		{"advanced uses text screen", homeActionBrowse, chapter.DifficultyAdvanced, 0, 8},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "progress.json")
			m := newModel(Options{Chapters: chapters, ProgressPath: path,
				LoadPreferences: func() ([]string, chapter.Difficulty, error) { return nil, tc.level, nil }})
			m.home.Select(homeActionIndex(tc.action))
			if tc.action == homeActionBrowse {
				m = press(t, m, tea.KeyEnter)
				m = press(t, m, tea.KeyDown)
				m = press(t, m, tea.KeyEnter)
				assertScreen(t, m, screenChapterDetail)
			}
			next, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
			m = next.(model)
			if tc.wantKeys == 0 {
				assertScreen(t, m, screenTraining)
				if m.session.ID != id || len(m.session.Items) != tc.wantText || m.session.HasKeyExercises() || m.result.CaptureChapter.ID != "" {
					t.Fatalf("unexpected text session: %#v", m.session)
				}
				return
			}
			assertQuitCommand(t, cmd)
			capture := m.result.CaptureChapter
			f := chapter.SessionFilter{Difficulty: tc.level}
			if capture.ID != id || len(capture.Items) != tc.wantKeys+tc.wantText ||
				len(f.EligibleItems(capture, chapter.KeyExercises)) != tc.wantKeys ||
				len(f.EligibleItems(capture, chapter.TextExercises)) != tc.wantText {
				t.Fatalf("unexpected capture selection: %#v", capture)
			}
		})
	}
}
