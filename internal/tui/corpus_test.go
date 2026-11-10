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
	if err != nil || state.ChapterScores["00-premiers-pas"] != (progress.ChapterScore{Attempts: 10, Correct: 10}) {
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
