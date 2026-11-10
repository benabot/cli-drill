package tui

import (
	"fmt"
	"strings"

	"github.com/benabot/cli-drill/internal/chapter"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) refreshSessionFilter() {
	m.filter.EnabledTools = nil
	m.filter.Difficulty = chapter.DifficultyBeginner
	m.preferencesError = nil
	if m.preferences.load != nil {
		tools, level, err := m.preferences.load()
		m.preferencesError = err
		if err == nil {
			m.filter.EnabledTools = tools
			if level != "" {
				m.filter.Difficulty = level
			}
		}
	}
}

func (m model) detailFilter() chapter.SessionFilter {
	filter := m.filter
	if m.detailDifficulty != "" {
		filter.Difficulty = m.detailDifficulty
	}
	return filter
}

func (m *model) openDetail(index int, scope chapter.ExerciseScope) {
	m.clearSession()
	m.chapterIndex = index
	m.detailScope = scope
	m.screen = screenChapterDetail
}

func (m *model) changeSessionChoice(key string) bool {
	current := m.currentChapter()
	if current == nil {
		return false
	}
	switch key {
	case "d":
		m.detailDifficulty = nextDifficulty(m.detailFilter().Difficulty)
		return true
	case "m":
		if !current.Guided() {
			return false
		}
		if m.detailMode == chapter.ModeReview {
			m.detailMode = chapter.ModeLearn
		} else {
			m.detailMode = chapter.ModeReview
		}
		return true
	default:
		return false
	}
}

func (m model) launchSession() (model, tea.Cmd) {
	current := m.currentChapter()
	if current == nil {
		return m, nil
	}
	filter := m.detailFilter()
	prepared := filter.Prepare(*current, m.detailScope, m.detailMode, nil)
	if len(prepared.Items) == 0 {
		return m, nil
	}
	if prepared.HasKeyExercises() {
		if m.captureUnavailable != "" {
			return m, nil
		}
		m.result.CaptureChapter = prepared
		return m, tea.Quit
	}
	m.session = prepared
	m.resetTextTraining()
	return m, nil
}

func (m *model) resetTextTraining() {
	m.screen = screenTraining
	m.itemIndex = 0
	m.choiceIndex = 0
	m.feedback = ""
	m.summary = trainingSummary{}
	m.showHelp = false
	m.input.SetValue("")
	m.input.Err = nil
	m.input.Focus()
}

func (m *model) clearSession() {
	m.detailDifficulty = ""
	m.detailMode = ""
	m.session = chapter.Chapter{}
	m.feedback = ""
	m.input.SetValue("")
}

func sessionLevelCounts(c chapter.Chapter, scope chapter.ExerciseScope, filter chapter.SessionFilter) string {
	selected := filter.Difficulty
	var counts []string
	for _, level := range []chapter.Difficulty{chapter.DifficultyBeginner, chapter.DifficultyIntermediate, chapter.DifficultyAdvanced} {
		if level == selected {
			continue
		}
		filter.Difficulty = level
		counts = append(counts, fmt.Sprintf("%s %d", sessionLevelLabel(level), len(filter.EligibleItems(c, scope))))
	}
	return "Autres niveaux · " + strings.Join(counts, " · ")
}

func sessionLevelLabel(level chapter.Difficulty) string {
	switch level {
	case chapter.DifficultyBeginner:
		return "Débutant"
	case chapter.DifficultyIntermediate:
		return "Intermédiaire"
	case chapter.DifficultyAdvanced:
		return "Avancé"
	default:
		return string(level)
	}
}

func nextDifficulty(level chapter.Difficulty) chapter.Difficulty {
	switch level {
	case chapter.DifficultyBeginner:
		return chapter.DifficultyIntermediate
	case chapter.DifficultyIntermediate:
		return chapter.DifficultyAdvanced
	default:
		return chapter.DifficultyBeginner
	}
}
