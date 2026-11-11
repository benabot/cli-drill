package tui

import (
	"fmt"
	"path/filepath"
	"slices"
	"testing"

	"github.com/benabot/cli-drill/internal/chapter"
	"github.com/benabot/cli-drill/internal/progress"
	tea "github.com/charmbracelet/bubbletea"
)

func TestVisitSummaryRetriesOnlyMissesAndContinues(t *testing.T) {
	ch := visitTestChapter()
	path := filepath.Join(t.TempDir(), "progress.json")
	m := startTestTraining(t, newModel(Options{Chapters: []chapter.Chapter{ch}, ProgressPath: path}))
	for i := range 6 {
		answer := m.currentItem().Answer.Primary
		if i == 1 {
			answer = "faux"
		}
		m = submitVisitTest(t, m, answer)
	}
	assertScreen(t, m, screenTrainingSummary)
	assertContains(t, m.View(), "Bonnes réponses : 5/6")
	assertContains(t, m.View(), "c Continuer")
	assertContains(t, m.View(), "r Revoir mes erreurs")
	m = press(t, m, tea.KeyRunes, 'r')
	if len(m.session.Items) != 1 || m.currentItem().ID != "item-2" {
		t.Fatalf("review restarted the batch: %v", tuiSessionIDs(m.session.Items))
	}
	m = submitVisitTest(t, m, m.currentItem().Answer.Primary)
	assertContains(t, m.View(), "Bilan de la reprise")
	assertContains(t, m.View(), "Lot initial : 5/6 bonnes réponses")
	assertNotContains(t, m.View(), "r Revoir mes erreurs")
	m = press(t, m, tea.KeyRunes, 'c')
	if !slices.Equal(tuiSessionIDs(m.session.Items), []string{"item-7", "item-8"}) {
		t.Fatalf("next batch repeated exercises: %v", tuiSessionIDs(m.session.Items))
	}
	for range 2 {
		m = submitVisitTest(t, m, m.currentItem().Answer.Primary)
	}
	assertScreen(t, m, screenTrainingSummary)
	assertNotContains(t, m.View(), "c Continuer")
	state, err := progress.Load(path)
	if err != nil || len(state.CompletedExercises) != 8 || state.ChapterScores[ch.ID] != (progress.ChapterScore{Attempts: 9, Correct: 8}) {
		t.Fatalf("progress: %#v (%v)", state, err)
	}
	m = press(t, m, tea.KeyEsc)
	assertContains(t, m.View(), "Tous les exercices de ce filtre sont acquis")
	m = press(t, m, tea.KeyEnter)
	assertScreen(t, m, screenChapterDetail)
	m = press(t, m, tea.KeyRunes, 'm')
	m = press(t, m, tea.KeyRunes, 'a')
	m = press(t, m, tea.KeyEnter)
	assertScreen(t, m, screenTraining)
	if len(m.session.Items) != 8 || m.visit.HasNext() {
		t.Fatalf("full review did not include the entire eligible set: %v", tuiSessionIDs(m.session.Items))
	}
}

func TestVisitCanContinueDespiteAnUncorrectedError(t *testing.T) {
	ch := visitTestChapter()
	path := filepath.Join(t.TempDir(), "progress.json")
	opts := Options{Chapters: []chapter.Chapter{ch}, ProgressPath: path}
	m := startTestTraining(t, newModel(opts))
	for i := range 6 {
		answer := m.currentItem().Answer.Primary
		if i == 0 {
			answer = "faux"
		}
		m = submitVisitTest(t, m, answer)
	}
	m = press(t, m, tea.KeyRunes, 'c')
	if m.currentItem().ID != "item-7" {
		t.Fatalf("continuation blocked by an error: %s", m.currentItem().ID)
	}
	m = startTestTraining(t, newModel(opts))
	if !slices.Equal(tuiSessionIDs(m.session.Items), []string{"item-1", "item-7", "item-8"}) {
		t.Fatalf("restart lost unacquired exercises: %v", tuiSessionIDs(m.session.Items))
	}
}

func TestFullFormatIsPassedToCapture(t *testing.T) {
	ch := keySequenceChapter()
	m := newTestModel(t, Options{Chapters: []chapter.Chapter{ch}})
	m.openDetail(0, chapter.AllExercises)
	m = press(t, m, tea.KeyRunes, 'a')
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	result := updated.(model).result
	assertQuitCommand(t, cmd)
	if !result.AllExercises || result.CaptureChapter.ID != ch.ID {
		t.Fatalf("full format lost during capture delegation: %#v", result)
	}
}

func submitVisitTest(t *testing.T, m model, answer string) model {
	t.Helper()
	m.input.SetValue(answer)
	m = press(t, m, tea.KeyEnter)
	return press(t, m, tea.KeyEnter)
}

func visitTestChapter() chapter.Chapter {
	ch := textChapter()
	item := ch.Items[0]
	ch.Items = nil
	for i := 1; i <= 8; i++ {
		item.ID = fmt.Sprintf("item-%d", i)
		item.Prompt = fmt.Sprintf("Exercise %d", i)
		item.Answer.Primary = fmt.Sprintf("answer%d", i)
		ch.Items = append(ch.Items, item)
	}
	return ch
}
