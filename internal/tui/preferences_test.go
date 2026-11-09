package tui

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/benabot/cli-drill/internal/chapter"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func TestPreferencesDraftCancelAndExplicitSave(t *testing.T) {
	enabled := []string{"retired-cli"}
	difficulty := chapter.DifficultyBeginner
	saves := 0
	opts := preferenceTestOptions()
	opts.LoadPreferences = func() ([]string, chapter.Difficulty, error) { return enabled, difficulty, nil }
	opts.SavePreferences = func(tools []string, level chapter.Difficulty) error {
		saves++
		enabled, difficulty = slices.Clone(tools), level
		return nil
	}
	m := newTestModel(t, opts)
	assertContains(t, m.View(), "p preferences")
	m = press(t, m, tea.KeyRunes, 'p')
	assertScreen(t, m, screenPreferences)
	assertContains(t, m.View(), "Default difficulty: beginner")
	assertContains(t, m.View(), "fd: disabled")
	assertContains(t, m.View(), "retired-cli")
	assertNotContains(t, m.View(), "retired-cli: enabled")

	m = press(t, m, tea.KeyEnter)
	m = press(t, m, tea.KeyDown)
	m = press(t, m, tea.KeyEnter)
	assertContains(t, m.View(), "Unsaved changes")
	assertContains(t, m.View(), "Default difficulty: intermediate")
	assertContains(t, m.View(), "fd: enabled")
	m = press(t, m, tea.KeyEsc)
	assertScreen(t, m, screenHome)
	if saves != 0 || difficulty != chapter.DifficultyBeginner || !slices.Equal(enabled, []string{"retired-cli"}) {
		t.Fatal("cancel changed preferences")
	}

	m = openTestPreferences(t, m)
	assertContains(t, m.View(), "fd: disabled")
	m = press(t, m, tea.KeyEnter)
	m = press(t, m, tea.KeyRunes, 'j')
	m = press(t, m, tea.KeyEnter)
	m = press(t, m, tea.KeyRunes, 's')
	assertContains(t, m.View(), "Preferences saved.")
	assertNotContains(t, m.View(), "Unsaved changes")
	if saves != 1 || difficulty != chapter.DifficultyIntermediate || !slices.Equal(enabled, []string{"retired-cli", "fd"}) {
		t.Fatalf("saved = %v, %s (%d saves)", enabled, difficulty, saves)
	}

	m = openTestPreferences(t, newTestModel(t, opts))
	assertContains(t, m.View(), "Default difficulty: intermediate")
	assertContains(t, m.View(), "fd: enabled")
	m = press(t, m, tea.KeyDown)
	m = press(t, m, tea.KeyEnter)
	m = press(t, m, tea.KeyRunes, 's')
	assertContains(t, m.View(), "fd: disabled")
	if saves != 2 || !slices.Equal(enabled, []string{"retired-cli"}) {
		t.Fatalf("disabling fd changed retired choices: %v (%d saves)", enabled, saves)
	}
}

func TestPreferencesFailuresRemainVisible(t *testing.T) {
	opts := preferenceTestOptions()
	opts.LoadPreferences = func() ([]string, chapter.Difficulty, error) {
		return nil, chapter.DifficultyBeginner, nil
	}
	attempts := 0
	opts.SavePreferences = func([]string, chapter.Difficulty) error {
		attempts++
		return errors.New("permission denied")
	}
	m := openTestPreferences(t, newTestModel(t, opts))
	m = press(t, m, tea.KeyEnter)
	m = press(t, m, tea.KeyRunes, 's')
	assertContains(t, m.View(), "Unsaved changes")
	assertContains(t, m.View(), "permission denied")
	assertNotContains(t, m.View(), "Preferences saved.")
	m = press(t, m, tea.KeyEsc)
	m = openTestPreferences(t, m)
	assertContains(t, m.View(), "Default difficulty: beginner")

	opts.LoadPreferences = func() ([]string, chapter.Difficulty, error) {
		return nil, chapter.DifficultyBeginner, errors.New("unsupported version")
	}
	m = openTestPreferences(t, newTestModel(t, opts))
	assertContains(t, m.View(), "showing defaults")
	assertContains(t, m.View(), "unsupported version")
	m = press(t, m, tea.KeyEnter)
	m = press(t, m, tea.KeyRunes, 's')
	assertNotContains(t, m.View(), "Preferences saved.")
	if attempts != 1 {
		t.Fatalf("invalid loaded state reached the writer: %d attempts", attempts)
	}
	m = press(t, m, tea.KeyEsc)
	m.home.Select(homeActionIndex("Start training"))
	m = press(t, m, tea.KeyEnter)
	assertScreen(t, m, screenTraining)
}

func TestPreferencesFitTerminalAndPaginate(t *testing.T) {
	opts := preferenceTestOptions()
	opts.Chapters[0].Items[0].Tools = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"}
	opts.LoadPreferences = func() ([]string, chapter.Difficulty, error) {
		return []string{"retired-cli"}, chapter.DifficultyBeginner, nil
	}
	m := newTestModel(t, opts)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 48, Height: 18})
	m = openTestPreferences(t, updated.(model))
	for range 12 {
		m = press(t, m, tea.KeyDown)
	}
	if m.preferences.list.Index() != 12 {
		t.Fatalf("arrow navigation reached index %d, want 12", m.preferences.list.Index())
	}
	view := m.View()
	assertContains(t, view, "l: disabled")
	assertContains(t, view, "esc back")
	lines := strings.Split(view, "\n")
	if len(lines) > 18 {
		t.Fatalf("preferences require %d lines:\n%s", len(lines), view)
	}
	for _, line := range lines {
		if width := lipgloss.Width(line); width > 48 {
			t.Fatalf("preferences line exceeds width: %d columns", width)
		}
	}
}

func openTestPreferences(t *testing.T, m model) model {
	t.Helper()
	m.home.Select(homeActionIndex("Preferences"))
	m = press(t, m, tea.KeyEnter)
	assertScreen(t, m, screenPreferences)
	return m
}

func preferenceTestOptions() Options {
	ch := textChapter()
	ch.Items = append(ch.Items, chapter.Item{Tools: []string{"fd", "fzf"}})
	return Options{Chapters: []chapter.Chapter{ch}}
}
