package tui

import (
	"strings"
	"testing"

	"github.com/benabot/cli-drill/internal/chapter"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func TestHomeOffersOnlyTrainingActions(t *testing.T) {
	m := newTestModel(t, Options{})
	want := []string{"Start training", "Browse chapters", "Keyboard shortcuts", "My progress", "Quit"}
	items := m.home.Items()
	if len(items) != len(want) {
		t.Fatalf("home has %d actions, want %d", len(items), len(want))
	}
	for i, title := range want {
		if items[i].(listItem).Title() != title {
			t.Errorf("action %d: %q, want %q", i, items[i].(listItem).Title(), title)
		}
	}
	for _, obsolete := range []string{"Dotfiles", "Configuration", "Directory", "cli-drill init"} {
		assertNotContains(t, m.View(), obsolete)
	}
}

func TestStartTrainingSkipsEmptyAndKeyCaptureChapters(t *testing.T) {
	text := textChapter()
	text.ID = "text"
	m := newTestModel(t, Options{Chapters: []chapter.Chapter{{ID: "empty"}, keySequenceChapter(), text}})
	m = press(t, m, tea.KeyEnter)
	assertScreen(t, m, screenTraining)
	if m.currentChapter().ID != "text" || m.currentItem() == nil {
		t.Fatal("Start training did not choose the immediately usable text chapter")
	}
	assertContains(t, m.View(), "Question:")
}

func TestKeyboardShortcutsReuseBrowserAndOriginalChapterIndices(t *testing.T) {
	first, second := keySequenceChapter(), keySequenceChapter()
	first.ID, first.Title = "keys-one", "Keys one"
	second.ID, second.Title = "keys-two", "Keys two"
	m := newTestModel(t, Options{Chapters: []chapter.Chapter{textChapter(), first, textChapter(), second}})
	m.home.Select(homeActionIndex("Keyboard shortcuts"))
	m = press(t, m, tea.KeyEnter)
	assertScreen(t, m, screenChapters)
	if len(m.chapterList.Items()) != 2 {
		t.Fatalf("keyboard browser contains non-capture chapters: %d items", len(m.chapterList.Items()))
	}
	m = press(t, m, tea.KeyRunes, 'j')
	m = press(t, m, tea.KeyEnter)
	assertScreen(t, m, screenChapterDetail)
	if m.currentChapter().ID != "keys-two" {
		t.Fatalf("filtered list lost the source chapter index: %s", m.currentChapter().ID)
	}
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if updated.(model).result.RunCLITrainChapterID != "keys-two" || cmd == nil {
		t.Fatal("keyboard chapter did not delegate to the existing capture engine")
	}
	m = press(t, m, tea.KeyEsc)
	assertScreen(t, m, screenChapters)
	if len(m.chapterList.Items()) != 2 {
		t.Fatal("back lost the keyboard-only browser")
	}
	m = press(t, m, tea.KeyRunes, 'g')
	m.home.Select(homeActionIndex("Browse chapters"))
	m = press(t, m, tea.KeyEnter)
	if len(m.chapterList.Items()) != 4 {
		t.Fatal("Browse chapters did not restore the complete corpus")
	}
}

func TestChapterFilterOwnsPrintableKeysAndBackspace(t *testing.T) {
	m := newTestModel(t, Options{})
	m.home.Select(homeActionIndex("Browse chapters"))
	m = press(t, m, tea.KeyEnter)
	m = press(t, m, tea.KeyRunes, '/')
	for _, key := range []rune{'g', 'h', 'q'} {
		updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{key}})
		m = updated.(model)
		if cmd != nil {
			if _, quit := cmd().(tea.QuitMsg); quit {
				t.Fatalf("filter input %q quit the TUI", key)
			}
		}
		assertScreen(t, m, screenChapters)
	}
	m = press(t, m, tea.KeyBackspace)
	if m.chapterList.FilterValue() != "gh" || m.showHelp {
		t.Fatalf("filter did not receive input: %q", m.chapterList.FilterValue())
	}
	m = press(t, m, tea.KeyEsc)
	assertScreen(t, m, screenChapters)
	m = press(t, m, tea.KeyRunes, 'g')
	assertScreen(t, m, screenHome)
}

func TestMyProgressHasNoDirectoryCounter(t *testing.T) {
	m := newTestModel(t, Options{})
	m.home.Select(homeActionIndex("My progress"))
	m = press(t, m, tea.KeyEnter)
	assertScreen(t, m, screenStats)
	assertContains(t, m.View(), "My progress")
	assertNotContains(t, m.View(), "Directory")
	assertNotContains(t, m.View(), "enter select")
	assertNotContains(t, m.View(), "j/k move")
	m = press(t, m, tea.KeyEsc)
	assertScreen(t, m, screenHome)
}

func TestNavigationScreensFitTerminalDimensions(t *testing.T) {
	for _, action := range []string{"Home", "Browse chapters", "Keyboard shortcuts", "My progress"} {
		t.Run(action, func(t *testing.T) {
			m := newTestModel(t, Options{Chapters: []chapter.Chapter{keySequenceChapter()}})
			updated, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
			m = updated.(model)
			if action != "Home" {
				m.home.Select(homeActionIndex(action))
				m = press(t, m, tea.KeyEnter)
			}
			view := m.View()
			if height := strings.Count(view, "\n") + 1; height > 24 {
				t.Fatalf("screen hides its header/footer: %d lines for a 24-line terminal", height)
			}
			for _, line := range strings.Split(view, "\n") {
				if width := lipgloss.Width(line); width > 80 {
					t.Fatalf("screen exceeds terminal width: %d columns", width)
				}
			}
		})
	}
}
