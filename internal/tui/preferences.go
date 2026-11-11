package tui

import (
	"fmt"
	"slices"
	"strings"

	"github.com/benabot/cli-drill/internal/chapter"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type preferencesScreen struct {
	list          list.Model
	load          func() ([]string, chapter.Difficulty, error)
	save          func([]string, chapter.Difficulty) error
	available     []string
	enabled       []string
	savedTools    []string
	difficulty    chapter.Difficulty
	savedLevel    chapter.Difficulty
	loadError     error
	saveError     error
	saved         bool
	width, height int
}

func newPreferencesScreen(opts Options, delegate list.DefaultDelegate) preferencesScreen {
	delegate.ShowDescription = false
	menu := list.New(nil, delegate, 76, 16)
	menu.Title = "Preferences"
	menu.Styles.Title = newTUIStyles().ScreenTitle
	menu.SetShowStatusBar(false)
	menu.SetFilteringEnabled(false)
	menu.SetShowHelp(false)
	return preferencesScreen{
		list: menu, load: opts.LoadPreferences, save: opts.SavePreferences,
		available: chapter.AvailableTools(opts.Chapters), width: 76, height: 24,
	}
}

func (p *preferencesScreen) open() {
	p.enabled, p.difficulty = nil, chapter.DifficultyBeginner
	p.loadError, p.saveError, p.saved = nil, nil, false
	if p.load != nil {
		p.enabled, p.difficulty, p.loadError = p.load()
	}
	p.savedTools, p.savedLevel = slices.Clone(p.enabled), p.difficulty
	p.refreshItems()
	p.resize(p.width, p.height)
	p.list.ResetSelected()
}

func (p *preferencesScreen) refreshItems() {
	items := []list.Item{listItem{title: "Default difficulty: " + string(p.difficulty)}}
	for _, tool := range p.available {
		status := "disabled"
		if slices.Contains(p.enabled, tool) {
			status = "enabled"
		}
		items = append(items, listItem{title: tool + ": " + status})
	}
	for _, tool := range chapter.UnavailableTools(p.available, p.enabled) {
		items = append(items, listItem{title: fmt.Sprintf("%q: removed from corpus (choice kept)", tool)})
	}
	p.list.SetItems(items)
}

func (p preferencesScreen) editable() bool {
	return p.loadError == nil && p.save != nil
}

func (p preferencesScreen) dirty() bool {
	if p.difficulty != p.savedLevel || len(p.enabled) != len(p.savedTools) {
		return true
	}
	for _, tool := range p.enabled {
		if !slices.Contains(p.savedTools, tool) {
			return true
		}
	}
	return false
}

func (p *preferencesScreen) changeSelected() {
	if !p.editable() {
		return
	}
	index := p.list.Index()
	switch {
	case index == 0:
		p.difficulty = nextDifficulty(p.difficulty)
	case index <= len(p.available):
		tool := p.available[index-1]
		if i := slices.Index(p.enabled, tool); i >= 0 {
			p.enabled = slices.Delete(slices.Clone(p.enabled), i, i+1)
		} else {
			p.enabled = append(slices.Clone(p.enabled), tool)
		}
	default:
		return // A removed module is displayed for information only.
	}
	p.saveError, p.saved = nil, false
	p.refreshItems()
	p.resize(p.width, p.height)
}

func (p *preferencesScreen) saveDraft() {
	if !p.editable() || !p.dirty() {
		return
	}
	p.saveError = p.save(p.enabled, p.difficulty)
	if p.saveError == nil {
		p.savedTools, p.savedLevel = slices.Clone(p.enabled), p.difficulty
		p.saved = true
	}
	p.resize(p.width, p.height)
}

func (p preferencesScreen) notes() string {
	styles := newTUIStyles()
	status := "No unsaved changes."
	switch {
	case p.loadError != nil:
		status = styles.Warning.Render(fmt.Sprintf("Preferences unavailable; showing defaults. %v", p.loadError))
	case p.saveError != nil:
		status = styles.Warning.Render(fmt.Sprintf("Unsaved changes. Could not save: %v", p.saveError))
	case p.saved:
		status = styles.Info.Render("Preferences saved.")
	case p.dirty():
		status = styles.Info.Render("Unsaved changes. Esc discards.")
	}
	return p.wrap("For future sessions. Modules need no installation.\n" + status)
}

func (p preferencesScreen) footerText() string {
	if !p.editable() {
		return "↑/↓ j/k · esc back"
	}
	return footerText(footerPreferences)
}

func (p preferencesScreen) wrap(text string) string {
	return lipgloss.NewStyle().Width(p.width).Render(text)
}

func (p *preferencesScreen) resize(width, height int) {
	p.width, p.height = width, height
	footerHeight := lipgloss.Height(p.wrap(p.footerText()))
	// Leave room for status, both separators and outer padding.
	p.list.SetSize(width, max(3, height-lipgloss.Height(p.notes())-footerHeight-6))
}

func (m *model) openPreferences() {
	m.preferences.open()
	m.screen = screenPreferences
	m.scrollOffset, m.showHelp = 0, false
}

func (m model) updatePreferences(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.Paste {
		return m, nil
	}
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "esc", "backspace", "g":
		m.goHome()
	case "enter":
		m.scrollOffset = 0
		m.preferences.changeSelected()
	case "s":
		m.preferences.saveDraft()
		if m.preferences.saveError != nil {
			m.scrollOffset = m.page().lastOffset()
		}
	default:
		m.scrollOffset = 0
		var cmd tea.Cmd
		m.preferences.list, cmd = m.preferences.list.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) preferencesView() string {
	return strings.Join([]string{
		m.preferences.list.View(),
		m.preferences.notes(),
	}, "\n")
}
