package tui

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/benabot/cli-drill/data"
	"github.com/benabot/cli-drill/internal/chapter"
	"github.com/benabot/cli-drill/internal/exercise"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func TestUnixPipelineFeedbackCanBeReadAtEverySize(t *testing.T) {
	chapters, err := chapter.LoadFS(data.Chapters(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	ch := chapters[slices.IndexFunc(chapters, func(c chapter.Chapter) bool { return c.ID == "14-recherche-texte-unix" })]
	item := ch.Items[slices.IndexFunc(ch.Items, func(i chapter.Item) bool { return i.ID == "compose-sed-awk-block-total" })]
	ch.Items = []chapter.Item{item}
	path := filepath.Join(t.TempDir(), "progress.json")
	m := startTestTraining(t, newModel(Options{
		Chapters: []chapter.Chapter{ch}, ProgressPath: path,
		LoadPreferences: func() ([]string, chapter.Difficulty, error) { return nil, chapter.DifficultyAdvanced, nil },
	}))
	m = resizeTest(t, m, 80, 24)
	m.input.SetValue("faux")
	m = press(t, m, tea.KeyEnter)
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, size := range [][2]int{{80, 24}, {60, 20}, {40, 12}, {80, 24}} {
		m = resizeTest(t, m, size[0], size[1])
		m = press(t, m, tea.KeyPgUp)
		assertContains(t, m.View(), "Dans suivi.txt")
		view := allTestPages(t, m, size[0], size[1])
		assertContains(t, view, "END { print n + 0 }'")
		assertContains(t, view, "une fois la lecture terminée.")
		assertContains(t, view, "enter next")
		if m.currentItem().ID != item.ID || m.summary.Total != 1 {
			t.Fatalf("navigation altered the exercise or score: %s/%d", m.currentItem().ID, m.summary.Total)
		}
	}
	after, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatalf("navigation altered progress: %v", err)
	}
}

func TestNarrowQuestionKeepsInputAndNavigationVisible(t *testing.T) {
	ch := textChapter()
	ch.Items[0].Prompt = strings.Repeat("Une consigne détaillée à relire.\n", 17) + "Dernière ligne."
	m := startTestTraining(t, newTestModel(t, Options{Chapters: []chapter.Chapter{ch}}))
	m = resizeTest(t, m, 40, 12)
	m = press(t, m, tea.KeyRunes, 'a')
	assertContains(t, m.View(), "> a")
	assertContains(t, m.View(), "PgUp/PgDn")
	assertContains(t, allTestPages(t, m, 40, 12), "Dernière ligne.")
	short := resizeTest(t, m, 40, 4)
	assertContains(t, short.View(), "> a")
	assertContains(t, allTestPages(t, short, 40, 4), "Dernière ligne.")
	m = press(t, m, tea.KeyPgDown)
	assertContains(t, m.View(), "> a")
	m = press(t, m, tea.KeyEsc)
	assertScreen(t, m, screenChapterDetail)
}

func TestNarrowQCMFollowsSelectionAndShowsQuestionInFeedback(t *testing.T) {
	ch := textChapter()
	ch.Items[0].ExerciseType = exercise.TypeMultipleChoice
	ch.Items[0].Prompt = "Consigne de choix.\n" + strings.Repeat("Une précision à lire.\n", 12)
	ch.Items[0].Choices = []string{"Premier choix", "Deuxième choix", "Troisième choix", "Quatrième choix"}
	ch.Items[0].Answer.Primary = "Premier choix"
	ch.Items[0].Explanation = strings.Repeat("Une explication complète.\n", 12) + "Fin du feedback."
	m := startTestTraining(t, newTestModel(t, Options{Chapters: []chapter.Chapter{ch}}))
	m = resizeTest(t, m, 40, 12)
	for range 3 {
		m = press(t, m, tea.KeyDown)
	}
	assertContains(t, m.View(), "> Quatrième choix")
	m = press(t, m, tea.KeyEnter)
	assertContains(t, m.View(), "Consigne de choix.")
	assertContains(t, m.View(), incorrectAnswer)
	view := allTestPages(t, m, 40, 12)
	assertContains(t, view, "Premier choix")
	assertContains(t, view, "Fin du feedback.")
}

func TestNarrowBrowserAndMissingModulesHelp(t *testing.T) {
	ch := textChapter()
	ch.Items[0].Tools = []string{"fd", "fzf"}
	m := resizeTest(t, newTestModel(t, Options{Chapters: []chapter.Chapter{ch}}), 40, 12)
	m.home.Select(homeActionIndex(homeActionBrowse))
	m = press(t, m, tea.KeyEnter)
	assertContains(t, m.View(), "/ search")
	m = press(t, m, tea.KeyEnter)
	view := allTestPages(t, m, 40, 12)
	assertContains(t, view, "fd + fzf")
	assertContains(t, view, "Preferences")
	m = press(t, m, tea.KeyRunes, 'h')
	assertContains(t, m.View(), "Help")
	m = press(t, m, tea.KeyEsc)
	assertScreen(t, m, screenChapters)
}

func TestNarrowProgressErrorRemainsReadable(t *testing.T) {
	path := filepath.Join(t.TempDir(), strings.Repeat("long-path-", 12)+".json")
	if err := os.WriteFile(path, []byte(`{"version":99}`), 0o600); err != nil {
		t.Fatal(err)
	}
	m := startTestTraining(t, newModel(Options{Chapters: []chapter.Chapter{textChapter()}, ProgressPath: path}))
	m = resizeTest(t, m, 40, 12)
	view := allTestPages(t, m, 40, 12)
	assertContains(t, view, "Progress could not be loaded")
	assertContains(t, strings.Join(strings.Fields(view), " "), "unsupported version 99")
	assertContains(t, view, "enter submit")
}

func resizeTest(t *testing.T, m model, width, height int) model {
	t.Helper()
	next, _ := m.Update(tea.WindowSizeMsg{Width: width, Height: height})
	return next.(model)
}

func allTestPages(t *testing.T, m model, width, height int) string {
	t.Helper()
	var pages []string
	for range 40 {
		view := m.View()
		if w, h := lipgloss.Size(view); w > width || h > height {
			t.Errorf("view at %dx%d exceeds bounds: %dx%d\n%s", width, height, w, h, view)
		}
		if len(pages) > 0 && pages[len(pages)-1] == view {
			break
		}
		pages = append(pages, view)
		m = press(t, m, tea.KeyPgDown)
	}
	return strings.Join(pages, "\n")
}
