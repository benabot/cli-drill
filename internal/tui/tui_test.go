package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/benabot/cli-drill/internal/catalog"
	"github.com/benabot/cli-drill/internal/chapter"
	"github.com/benabot/cli-drill/internal/exercise"
	tea "github.com/charmbracelet/bubbletea"
)

func TestHomeViewExplainsAppAndShowsConfigStatus(t *testing.T) {
	m := newTestModel(t, Options{ConfigPath: filepath.Join(t.TempDir(), "config.toml")})

	view := m.View()

	assertContains(t, view, "cli-drill")
	assertContains(t, view, "Train terminal shortcuts")
	assertContains(t, view, "Dotfiles: not configured")
	assertContains(t, view, "Start training")
	assertContains(t, view, "Browse chapters")
	assertContains(t, view, "Configuration")
	assertContains(t, view, "enter select")
	assertContains(t, view, "esc back")
	assertContains(t, view, "g home")
}

func TestConfigStatusMissingIsReadOnlyGuidance(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.toml")
	m := newTestModel(t, Options{ConfigPath: configPath})
	m.home.Select(homeActionIndex("Configuration"))
	m = press(t, m, tea.KeyEnter)

	view := m.View()

	assertScreen(t, m, screenConfig)
	assertContains(t, view, "Configuration")
	assertContains(t, view, configPath)
	assertContains(t, view, "Status: not configured")
	assertContains(t, view, "cli-drill init --print")
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		t.Fatalf("expected config screen to avoid creating %s, got err=%v", configPath, err)
	}
}

func TestConfigStatusPresentShowsDotfilesPath(t *testing.T) {
	dir := t.TempDir()
	dotfilesPath := filepath.Join(dir, "dotfiles")
	if err := os.Mkdir(dotfilesPath, 0o755); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(dir, "config.toml")
	data := "dotfiles_path = " + quoteTOML(dotfilesPath) + "\nshell = \"zsh\"\n"
	if err := os.WriteFile(configPath, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}

	m := newTestModel(t, Options{ConfigPath: configPath})
	m.home.Select(homeActionIndex("Configuration"))
	m = press(t, m, tea.KeyEnter)

	view := m.View()

	assertContains(t, view, "Status: configured")
	assertContains(t, view, "Dotfiles path: "+dotfilesPath)
	assertContains(t, view, "Path exists: yes")
}

func TestNavigationHomeChaptersDetailBackHome(t *testing.T) {
	m := newTestModel(t, Options{})

	m.home.Select(homeActionIndex("Browse chapters"))
	m = press(t, m, tea.KeyEnter)
	assertScreen(t, m, screenChapters)

	m = press(t, m, tea.KeyEnter)
	assertScreen(t, m, screenChapterDetail)
	assertContains(t, m.View(), "Raccourcis terminal")
	assertContains(t, m.View(), "Exercises: 1")

	m = press(t, m, tea.KeyEsc)
	assertScreen(t, m, screenChapters)

	m = press(t, m, tea.KeyRunes, 'g')
	assertScreen(t, m, screenHome)
}

func TestChapterDetailStartsTraining(t *testing.T) {
	m := newTestModel(t, Options{})
	m.home.Select(homeActionIndex("Browse chapters"))
	m = press(t, m, tea.KeyEnter)
	m = press(t, m, tea.KeyEnter)

	m = press(t, m, tea.KeyEnter)

	assertScreen(t, m, screenTraining)
	assertContains(t, m.View(), "Question:")
	assertContains(t, m.View(), "Aller au debut de ligne")
}

func TestChapterDetailShowsCLIOnlyGuidanceForKeySequence(t *testing.T) {
	m := newTestModel(t, Options{Chapters: []chapter.Chapter{keySequenceChapter()}})
	m.home.Select(homeActionIndex("Browse chapters"))
	m = press(t, m, tea.KeyEnter)
	m = press(t, m, tea.KeyEnter)

	view := m.View()

	assertScreen(t, m, screenChapterDetail)
	assertContains(t, view, "Ce chapitre entraîne les vrais raccourcis clavier.")
	assertContains(t, view, "cli-drill ouvre un mode dédié pour capturer les touches.")
	assertNotContains(t, strings.ToLower(view), "run cli-drill again")
}

func TestChapterDetailKeySequenceFooterOffersRunInCLI(t *testing.T) {
	m := newTestModel(t, Options{Chapters: []chapter.Chapter{keySequenceChapter()}})
	m.home.Select(homeActionIndex("Browse chapters"))
	m = press(t, m, tea.KeyEnter)
	m = press(t, m, tea.KeyEnter)

	assertContains(t, m.View(), "enter start key training · esc chapters · g home · h help · q quit")
}

func TestChapterDetailDoesNotStartBubbleTrainingForKeySequence(t *testing.T) {
	m := newTestModel(t, Options{Chapters: []chapter.Chapter{keySequenceChapter()}})
	m.home.Select(homeActionIndex("Browse chapters"))
	m = press(t, m, tea.KeyEnter)
	m = press(t, m, tea.KeyEnter)

	m = press(t, m, tea.KeyEnter)

	assertScreen(t, m, screenChapterDetail)
	assertContains(t, m.View(), "mode dédié")
}

func TestChapterDetailKeySequenceEnterReturnsCLITrainIntent(t *testing.T) {
	m := newTestModel(t, Options{Chapters: []chapter.Chapter{keySequenceChapter()}})
	m.home.Select(homeActionIndex("Browse chapters"))
	m = press(t, m, tea.KeyEnter)
	m = press(t, m, tea.KeyEnter)

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	next, ok := updated.(model)
	if !ok {
		t.Fatalf("expected model, got %T", updated)
	}

	assertScreen(t, next, screenChapterDetail)
	if cmd == nil {
		t.Fatal("expected Enter on key-sequence detail to quit TUI")
	}
	if next.result.RunCLITrainChapterID != "01-raccourcis-terminal" {
		t.Fatalf("RunCLITrainChapterID = %q, want 01-raccourcis-terminal", next.result.RunCLITrainChapterID)
	}
}

func TestChapterDetailStillStartsBubbleTrainingForTextChapter(t *testing.T) {
	m := newTestModel(t, Options{Chapters: []chapter.Chapter{textChapter()}})
	m.home.Select(homeActionIndex("Browse chapters"))
	m = press(t, m, tea.KeyEnter)
	m = press(t, m, tea.KeyEnter)

	m = press(t, m, tea.KeyEnter)

	assertScreen(t, m, screenTraining)
	assertContains(t, m.View(), "Question:")
	assertContains(t, m.View(), "Aller au debut de ligne")
}

func TestTrainingTextFooterIsConsistent(t *testing.T) {
	m := newTestModel(t, Options{Chapters: []chapter.Chapter{textChapter()}})
	m.home.Select(homeActionIndex("Browse chapters"))
	m = press(t, m, tea.KeyEnter)
	m = press(t, m, tea.KeyEnter)
	m = press(t, m, tea.KeyEnter)

	assertContains(t, m.View(), "enter submit/next · esc chapter · g home · h help · q quit")
}

func TestScanScreenIsGuidanceOnly(t *testing.T) {
	m := newTestModel(t, Options{})
	m.home.Select(homeActionIndex("Scan dotfiles"))
	m = press(t, m, tea.KeyEnter)

	view := m.View()

	assertScreen(t, m, screenScan)
	assertContains(t, view, "Scan dotfiles")
	assertContains(t, view, "read-only")
	assertContains(t, view, "cli-drill scan --summary")
}

func TestFooterHintsAreContextual(t *testing.T) {
	m := newTestModel(t, Options{})

	assertContains(t, m.View(), "enter select")
	assertContains(t, m.View(), "h help")
	assertContains(t, m.View(), "q quit")

	m.home.Select(homeActionIndex("Browse chapters"))
	m = press(t, m, tea.KeyEnter)
	m = press(t, m, tea.KeyEnter)

	assertContains(t, m.View(), "enter start training")
	assertContains(t, m.View(), "esc chapters")
	assertContains(t, m.View(), "g home")
}

func TestConfigScanStatsFootersAreConsistent(t *testing.T) {
	tests := []struct {
		name   string
		action string
		screen screen
	}{
		{name: "config", action: "Configuration", screen: screenConfig},
		{name: "scan", action: "Scan dotfiles", screen: screenScan},
		{name: "stats", action: "Stats", screen: screenStats},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newTestModel(t, Options{})
			m.home.Select(homeActionIndex(tt.action))
			m = press(t, m, tea.KeyEnter)

			assertScreen(t, m, tt.screen)
			assertContains(t, m.View(), "esc back · g home · h help · q quit")
		})
	}
}

func TestBackAndQuitKeys(t *testing.T) {
	m := newTestModel(t, Options{})
	m.home.Select(homeActionIndex("Stats"))
	m = press(t, m, tea.KeyEnter)
	assertScreen(t, m, screenStats)

	m = press(t, m, tea.KeyEsc)
	assertScreen(t, m, screenHome)

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Fatal("expected q to quit")
	}
}

func newTestModel(t *testing.T, opts Options) model {
	t.Helper()
	if opts.Chapters == nil {
		opts.Chapters = []chapter.Chapter{textChapter()}
	}
	if opts.Entries == nil {
		opts.Entries = []catalog.Entry{{
			ID:      "ctrl-a",
			Name:    "Ctrl+A",
			Type:    catalog.EntryShortcut,
			Summary: "Aller au debut de ligne",
		}}
	}
	return newModel(opts)
}

func textChapter() chapter.Chapter {
	return chapter.Chapter{
		ID:          "01-raccourcis-terminal",
		Title:       "Raccourcis terminal",
		Description: "Memoriser les raccourcis de navigation.",
		Items: []chapter.Item{{
			ID:           "line-start",
			Type:         catalog.EntryShortcut,
			ExerciseType: exercise.TypeFreeAnswer,
			Prompt:       "Aller au debut de ligne",
			Answer:       exercise.Answer{Primary: "Ctrl+A"},
			Explanation:  "Ctrl+A place le curseur au debut.",
		}},
	}
}

func keySequenceChapter() chapter.Chapter {
	chapter := textChapter()
	chapter.Items[0].ExerciseType = exercise.TypeKeySequence
	return chapter
}

func press(t *testing.T, m model, key tea.KeyType, runes ...rune) model {
	t.Helper()
	updated, _ := m.Update(tea.KeyMsg{Type: key, Runes: runes})
	next, ok := updated.(model)
	if !ok {
		t.Fatalf("expected model, got %T", updated)
	}
	return next
}

func assertScreen(t *testing.T, m model, want screen) {
	t.Helper()
	if m.screen != want {
		t.Fatalf("screen = %v, want %v\nview:\n%s", m.screen, want, m.View())
	}
}

func assertContains(t *testing.T, value, needle string) {
	t.Helper()
	if !strings.Contains(value, needle) {
		t.Fatalf("expected view to contain %q\nview:\n%s", needle, value)
	}
}

func assertNotContains(t *testing.T, value, needle string) {
	t.Helper()
	if strings.Contains(value, needle) {
		t.Fatalf("expected view not to contain %q\nview:\n%s", needle, value)
	}
}

func quoteTOML(value string) string {
	return `"` + strings.ReplaceAll(value, `\`, `\\`) + `"`
}
