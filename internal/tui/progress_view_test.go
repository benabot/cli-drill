package tui

import (
	"strings"
	"testing"

	"github.com/benabot/cli-drill/data"
	"github.com/benabot/cli-drill/internal/chapter"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func TestProgressViewWithNoHistory(t *testing.T) {
	m := newTestModel(t, Options{})
	m.home.Select(homeActionIndex(homeActionStats))
	m = press(t, m, tea.KeyEnter)
	view := m.View()
	for _, label := range []string{
		"Progression actuelle", "Niveau : Débutant",
		"Acquis : 0 sur 1 exercice éligible", "[------------]",
		"Selon le niveau et les modules activés",
		"Historique global", "Exercices acquis au total : 0", "Chapitres pratiqués : 0",
		"Les acquis restent dans l’historique après désactivation d’un module.",
		"Échap Retour",
	} {
		assertContains(t, view, label)
	}
	m = press(t, m, tea.KeyEsc)
	assertScreen(t, m, screenHome)
}

func TestProgressAndChapterViewsFitAfterResize(t *testing.T) {
	chapters, err := chapter.LoadFS(data.Chapters(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	m := newTestModel(t, Options{Chapters: chapters})
	for _, width := range []int{80, 64, 80} {
		updated, _ := m.Update(tea.WindowSizeMsg{Width: width, Height: 24})
		m = updated.(model)
		m.home.Select(homeActionIndex(homeActionStats))
		m = press(t, m, tea.KeyEnter)
		views := map[string]string{"progress": m.View()}
		m = press(t, m, tea.KeyEsc)
		for i, ch := range chapters {
			m.openDetail(i, chapter.AllExercises)
			view := m.View()
			assertContains(t, strings.Join(strings.Fields(view), " "), ch.Description)
			for _, heading := range []string{"Chapitre", "Ma séance", "Ma progression"} {
				assertContains(t, view, heading)
			}
			if !ch.Guided() {
				assertContains(t, view, "Ordre aléatoire")
				assertNotContains(t, view, "m Changer")
				m = press(t, m, tea.KeyRunes, 'm')
				if m.detailMode != "" {
					t.Fatal("m changed the order of an independent chapter")
				}
			}
			views[ch.ID] = view
		}
		for name, view := range views {
			if w, h := lipgloss.Size(view); w > width || h > 24 {
				t.Fatalf("%s exceeds %dx24: %dx%d\n%s", name, width, w, h, view)
			}
		}
		m = press(t, m, tea.KeyRunes, 'g')
		assertScreen(t, m, screenHome)
	}
}
