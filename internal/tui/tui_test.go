package tui

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/benabot/cli-drill/data"
	"github.com/benabot/cli-drill/internal/chapter"
	"github.com/benabot/cli-drill/internal/exercise"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func TestHomeBannerAndTrainingTagline(t *testing.T) {
	m := newTestModel(t, Options{})

	view := m.View()

	assertContains(t, view, "cli-drill")
	assertContains(t, view, "____ _     ___")
	assertContains(t, view, homeTagline)
	assertContains(t, view, "Start training")
	assertContains(t, view, "Browse chapters")
	assertContains(t, view, "Choose a chapter and inspect its exercises.")
	assertContains(t, view, "enter select")
	assertContains(t, view, "esc back")
	assertContains(t, view, "g home")
}

func TestHomeUsesOnlyCentralizedFooter(t *testing.T) {
	m := newTestModel(t, Options{})

	view := m.View()
	wantFooter := footerText(footerHome)

	if count := strings.Count(view, wantFooter); count != 1 {
		t.Fatalf("centralized footer count = %d, want 1\nview:\n%s", count, view)
	}
	assertNotContains(t, view, "↑/k up")
	assertNotContains(t, view, "? more")
}

func TestHomeBannerFallsBackWhenTerminalIsNarrow(t *testing.T) {
	styles := newTUIStyles()
	banner := renderHomeBanner(20, styles)

	assertContains(t, banner, "cli-drill")
	assertNotContains(t, banner, "____ _     ___")
}

func TestChaptersKeepDescriptionInDetailAndShowCompactMetadata(t *testing.T) {
	first := textChapter()
	second := textChapter()
	second.ID = "02-navigation-shell"
	second.Title = "Navigation shell"
	second.Description = "Move between directories and inspect files from the terminal. The full description belongs on the chapter detail screen."
	m := newTestModel(t, Options{Chapters: []chapter.Chapter{first, second}})
	resized, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = resized.(model)
	m.home.Select(homeActionIndex("Browse chapters"))
	m = press(t, m, tea.KeyEnter)

	view := m.View()
	wantFooter := footerText(footerNavigation)

	assertContains(t, view, "│ Raccourcis terminal")
	assertContains(t, view, "Exercices éligibles : 1/1 · 01-raccourcis-terminal")
	assertNotContains(t, view, "Memoriser les raccourcis")
	assertContains(t, view, "Exercices éligibles : 1/1 · 02-navigation-shell")
	assertNotContains(t, view, "Move between directories")
	assertContains(t, m.chapterList.Items()[1].FilterValue(), second.Description)
	if width, height := lipgloss.Size(view); width > 80 || height > 24 {
		t.Fatalf("chapter list exceeds the terminal: %dx%d", width, height)
	}
	if count := strings.Count(view, wantFooter); count != 1 {
		t.Fatalf("centralized footer count = %d, want 1\nview:\n%s", count, view)
	}
	assertNotContains(t, view, "↑/k up")

	m = press(t, m, tea.KeyRunes, 'j')
	assertContains(t, m.View(), "│ Navigation shell")
	assertNotContains(t, m.View(), "│ Raccourcis terminal")
	m = press(t, m, tea.KeyEnter)
	assertScreen(t, m, screenChapterDetail)
	assertContains(t, m.View(), "Navigation shell")
	assertContains(t, strings.Join(strings.Fields(m.View()), " "), second.Description)
}

func TestCompactChapterListPagesThroughEmbeddedCorpusAt80Columns(t *testing.T) {
	chapters, err := chapter.LoadFS(data.Chapters(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	m := newTestModel(t, Options{Chapters: chapters})
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = updated.(model)
	m.home.Select(homeActionIndex(homeActionBrowse))
	m = press(t, m, tea.KeyEnter)
	firstPage := m.chapterList.Paginator.Page
	for i, ch := range chapters {
		view := m.View()
		assertContains(t, view, "│ "+ch.Title)
		assertContains(t, view, " · "+ch.ID)
		if width, height := lipgloss.Size(view); width > 80 || height > 24 {
			t.Fatalf("%s exceeds the terminal: %dx%d", ch.ID, width, height)
		}
		if i+1 < len(chapters) {
			m = press(t, m, tea.KeyRunes, 'j')
		}
	}
	if m.chapterList.Paginator.Page == firstPage {
		t.Fatal("chapter navigation did not reach a second page")
	}
	m = press(t, m, tea.KeyRunes, 'k')
	assertContains(t, m.View(), "│ "+chapters[len(chapters)-2].Title)
}

func TestNavigationHomeChaptersDetailBackHome(t *testing.T) {
	m := newTestModel(t, Options{})

	m.home.Select(homeActionIndex("Browse chapters"))
	m = press(t, m, tea.KeyEnter)
	assertScreen(t, m, screenChapters)

	m = press(t, m, tea.KeyEnter)
	assertScreen(t, m, screenChapterDetail)
	assertContains(t, m.View(), "Raccourcis terminal")
	assertContains(t, m.View(), "Exercices au total : 1")

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
	assertContains(t, m.View(), "Question")
	assertContains(t, m.View(), "Aller au debut de ligne")
}

func TestChapterDetailShowsCLIOnlyGuidanceForKeySequence(t *testing.T) {
	m := newTestModel(t, Options{Chapters: []chapter.Chapter{keySequenceChapter()}})
	m.home.Select(homeActionIndex("Browse chapters"))
	m = press(t, m, tea.KeyEnter)
	m = press(t, m, tea.KeyEnter)

	view := m.View()

	assertScreen(t, m, screenChapterDetail)
	assertContains(t, view, "Capture réelle des raccourcis dans un mode dédié.")
	assertNotContains(t, strings.ToLower(view), "run cli-drill again")
}

func TestChapterDetailKeySequenceFooterOffersRunInCLI(t *testing.T) {
	m := newTestModel(t, Options{Chapters: []chapter.Chapter{keySequenceChapter()}})
	m.home.Select(homeActionIndex("Browse chapters"))
	m = press(t, m, tea.KeyEnter)
	m = press(t, m, tea.KeyEnter)

	assertContains(t, m.View(), "Entrée Capturer · Échap Chapitres · g Accueil · h Aide · q Quitter")
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
	if next.result.CaptureChapter.ID != "01-raccourcis-terminal" {
		t.Fatalf("CaptureChapter.ID = %q, want 01-raccourcis-terminal", next.result.CaptureChapter.ID)
	}
}

func TestChapterDetailStillStartsBubbleTrainingForTextChapter(t *testing.T) {
	m := newTestModel(t, Options{Chapters: []chapter.Chapter{textChapter()}})
	m.home.Select(homeActionIndex("Browse chapters"))
	m = press(t, m, tea.KeyEnter)
	m = press(t, m, tea.KeyEnter)

	m = press(t, m, tea.KeyEnter)

	assertScreen(t, m, screenTraining)
	assertContains(t, m.View(), "Question")
	assertContains(t, m.View(), "Aller au debut de ligne")
}

func TestTrainingTextInputKeepsPrintableKeysInTheAnswer(t *testing.T) {
	m := newTestModel(t, Options{Chapters: []chapter.Chapter{textChapter()}})
	m.home.Select(homeActionIndex("Browse chapters"))
	m = press(t, m, tea.KeyEnter)
	m = press(t, m, tea.KeyEnter)
	m = press(t, m, tea.KeyEnter)

	for _, key := range []rune{'c', 'g', 's'} {
		m = press(t, m, tea.KeyRunes, key)
	}

	if m.input.Value() != "cgs" {
		t.Fatalf("input = %q, want cgs", m.input.Value())
	}
	assertScreen(t, m, screenTraining)
	if m.showHelp {
		t.Fatal("h should be typed into the answer instead of opening help")
	}
}

func TestTrainingTextInputDisablesGlobalSingleKeyShortcuts(t *testing.T) {
	m := newTrainingModel(t)

	for _, key := range []rune{'g', 'h', 'q', 'c'} {
		updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{key}})
		next, ok := updated.(model)
		if !ok {
			t.Fatalf("expected model, got %T", updated)
		}
		if cmd != nil {
			if _, ok := cmd().(tea.QuitMsg); ok {
				t.Fatalf("%q should not quit while typing", key)
			}
		}
		m = next
	}

	assertScreen(t, m, screenTraining)
	if m.showHelp {
		t.Fatal("h should not open help while typing")
	}
	if m.input.Value() != "ghqc" {
		t.Fatalf("input = %q, want ghqc", m.input.Value())
	}
}

func TestTrainingTextInputFooterIsContextual(t *testing.T) {
	m := newTrainingModel(t)

	assertContains(t, m.View(), "enter submit · esc chapter · ctrl+c quit")
	assertNotContains(t, m.View(), "g home")
	assertNotContains(t, m.View(), "h help")
	assertNotContains(t, m.View(), "q quit")

	m.input.SetValue("Ctrl+A")
	m = press(t, m, tea.KeyEnter)

	assertContains(t, m.View(), "enter next · esc chapter · ctrl+c quit")
}

func TestTrainingTextInputEscReturnsToChapterDetail(t *testing.T) {
	m := newTrainingModel(t)

	m = press(t, m, tea.KeyEsc)

	assertScreen(t, m, screenChapterDetail)
}

func TestTrainingTextInputCtrlCQuits(t *testing.T) {
	m := newTrainingModel(t)

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Fatal("expected ctrl+c to quit while typing")
	}
}

func TestTrainingSummaryShowsMissedAnswersAndOffersRetry(t *testing.T) {
	m := newTrainingModel(t)
	assertNotContains(t, m.View(), "Ctrl+A")
	assertNotContains(t, m.View(), "Explication")
	m.input.SetValue("Ctrl+L")
	m = press(t, m, tea.KeyEnter)
	assertContains(t, m.View(), "Réponse incorrecte.")
	assertContains(t, m.View(), "Réponse attendue : Ctrl+A")
	assertContains(t, m.View(), "Explication")
	assertContains(t, m.View(), "Ctrl+A place le curseur au debut.")
	m = press(t, m, tea.KeyEnter)

	assertScreen(t, m, screenTrainingSummary)
	view := m.View()
	assertContains(t, view, "Bonnes réponses : 0/1")
	assertContains(t, view, "À revoir")
	assertContains(t, view, "Aller au debut de ligne")
	assertContains(t, view, "Réponse attendue : Ctrl+A")
	assertContains(t, view, "Ta réponse : Ctrl+L")
	assertContains(t, view, "r Revoir mes erreurs")

	m = press(t, m, tea.KeyEnter)
	assertScreen(t, m, screenTraining)
}

func TestTrainingShowsSummaryWithoutMissedAnswers(t *testing.T) {
	m := newTrainingModel(t)
	m.input.SetValue("Ctrl+A")
	m = press(t, m, tea.KeyEnter)
	m = press(t, m, tea.KeyEnter)

	assertScreen(t, m, screenTrainingSummary)
	view := m.View()
	assertNotContains(t, view, "À revoir")
	assertNotContains(t, view, "r Revoir mes erreurs")
}

func TestTrainingSummaryEscReturnsToChapterDetail(t *testing.T) {
	m := newMissedTrainingSummary(t)

	m = press(t, m, tea.KeyEsc)

	assertScreen(t, m, screenChapterDetail)
}

func TestTrainingSummaryQQuits(t *testing.T) {
	m := newMissedTrainingSummary(t)

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Fatal("expected q to quit from the training summary")
	}
}

func TestTrainingSummaryGReturnsHome(t *testing.T) {
	m := newMissedTrainingSummary(t)

	m = press(t, m, tea.KeyRunes, 'g')

	assertScreen(t, m, screenHome)
}

func TestFooterHintsAreContextual(t *testing.T) {
	m := newTestModel(t, Options{})

	assertContains(t, m.View(), "enter select")
	assertContains(t, m.View(), "h help")
	assertContains(t, m.View(), "q quit")

	m.home.Select(homeActionIndex("Browse chapters"))
	m = press(t, m, tea.KeyEnter)
	m = press(t, m, tea.KeyEnter)

	assertContains(t, m.View(), "Entrée Démarrer")
	assertContains(t, m.View(), "Échap Chapitres")
	assertContains(t, m.View(), "g Accueil")
}

func TestBackAndQuitKeys(t *testing.T) {
	m := newTestModel(t, Options{})
	m.home.Select(homeActionIndex("My progress"))
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
	if opts.ProgressPath == "" {
		opts.ProgressPath = filepath.Join(t.TempDir(), "progress.json")
	}
	return newModel(opts)
}

func newTrainingModel(t *testing.T) model {
	t.Helper()
	m := newTestModel(t, Options{Chapters: []chapter.Chapter{textChapter()}})
	m.home.Select(homeActionIndex("Browse chapters"))
	m = press(t, m, tea.KeyEnter)
	m = press(t, m, tea.KeyEnter)
	m = press(t, m, tea.KeyEnter)
	assertScreen(t, m, screenTraining)
	return m
}

func newMissedTrainingSummary(t *testing.T) model {
	t.Helper()
	m := newTrainingModel(t)
	m.input.SetValue("Ctrl+L")
	m = press(t, m, tea.KeyEnter)
	m = press(t, m, tea.KeyEnter)
	assertScreen(t, m, screenTrainingSummary)
	return m
}

func textChapter() chapter.Chapter {
	return chapter.Chapter{
		ID:          "01-raccourcis-terminal",
		Title:       "Raccourcis terminal",
		Description: "Memoriser les raccourcis de navigation.",
		Items: []chapter.Item{{
			ID:           "line-start",
			Type:         "shortcut",
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

func homeActionIndex(title string) int {
	for i, action := range homeActions {
		if action.title == title {
			return i
		}
	}
	return 0
}

func startTestTraining(t *testing.T, m model) model {
	t.Helper()
	m.home.Select(homeActionIndex(homeActionStartTraining))
	m = press(t, m, tea.KeyEnter)
	assertScreen(t, m, screenTraining)
	return m
}
