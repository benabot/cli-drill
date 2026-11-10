package tui

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/benabot/cli-drill/data"
	"github.com/benabot/cli-drill/internal/chapter"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func TestModernExpansionFits80ColumnsBeforeAndAfterSubmission(t *testing.T) {
	chapters, err := chapter.LoadFS(data.Chapters(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	ids := []string{
		"eza-tree-depth", "eza-hidden-long", "zoxide-list-matches",
		"atuin-failed-search", "rg-numbered", "rg-matching-files",
		"rg-fixed-dot", "rg-go-glob", "fd-go-files", "fd-hidden-notes",
		"rg-files-select", "bat-number-lines", "bat-line-range",
		"bat-json-stdin", "micro-undo", "micro-line-numbers",
		"glow-pager", "glow-width",
	}
	seen := 0
	for _, ch := range chapters {
		for _, item := range ch.Items {
			if !slices.Contains(ids, item.ID) {
				continue
			}
			t.Run(ch.ID+"/"+item.ID, func(t *testing.T) {
				selected := ch
				selected.Items = []chapter.Item{item}
				m := newModel(Options{
					Chapters:     []chapter.Chapter{selected},
					ProgressPath: filepath.Join(t.TempDir(), "progress.json"),
					LoadPreferences: func() ([]string, chapter.Difficulty, error) {
						return item.Tools, item.Difficulty, nil
					},
				})
				updated, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
				m = startTestTraining(t, updated.(model))
				for _, line := range strings.Split(item.Prompt, "\n") {
					assertContains(t, m.View(), line)
				}
				for _, answer := range []string{"incorrect", item.Answer.Primary} {
					m.input.SetValue(answer)
					m = press(t, m, tea.KeyEnter)
					view := m.View()
					for _, line := range strings.Split(item.Explanation, "\n") {
						assertContains(t, view, line)
					}
					if width, height := lipgloss.Size(view); width > 80 || height > 24 {
						t.Fatalf("feedback exceeds 80x24: %dx%d\n%s", width, height, view)
					}
					m = press(t, m, tea.KeyEsc)
					m = press(t, m, tea.KeyEnter)
				}
			})
			seen++
		}
	}
	if seen != len(ids) {
		t.Fatalf("found %d new exercises, want %d", seen, len(ids))
	}
}
