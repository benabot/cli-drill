package tui

import (
	"fmt"

	"github.com/benabot/cli-drill/internal/catalog"
	"github.com/benabot/cli-drill/internal/chapter"
	"github.com/benabot/cli-drill/internal/exercise"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type screen int

const (
	screenMenu screen = iota
	screenChapters
	screenTraining
	screenDirectory
	screenStats
)

type listItem struct {
	title       string
	description string
	index       int
}

func (i listItem) Title() string       { return i.title }
func (i listItem) Description() string { return i.description }
func (i listItem) FilterValue() string { return i.title + " " + i.description }

type model struct {
	screen       screen
	menu         list.Model
	chapterList  list.Model
	directory    list.Model
	input        textinput.Model
	chapters     []chapter.Chapter
	entries      []catalog.Entry
	chapterIndex int
	itemIndex    int
	feedback     string
}

func Run(chapters []chapter.Chapter, entries []catalog.Entry) error {
	delegate := list.NewDefaultDelegate()

	menu := list.New([]list.Item{
		listItem{title: "Continuer", description: "Reprendre le premier chapitre disponible"},
		listItem{title: "Choisir un chapitre", description: fmt.Sprintf("%d chapitres disponibles", len(chapters))},
		listItem{title: "Annuaire", description: fmt.Sprintf("%d entrees disponibles", len(entries))},
		listItem{title: "Stats", description: "Voir la progression de cette session TUI"},
		listItem{title: "Quitter", description: "Fermer cli-drill"},
	}, delegate, 80, 18)
	menu.Title = "cli-drill"
	menu.SetShowStatusBar(false)

	chapterItems := make([]list.Item, 0, len(chapters))
	for i, chapter := range chapters {
		chapterItems = append(chapterItems, listItem{
			title:       chapter.ID,
			description: fmt.Sprintf("%s - %d exercices", chapter.Title, len(chapter.Items)),
			index:       i,
		})
	}
	chapterList := list.New(chapterItems, delegate, 80, 18)
	chapterList.Title = "Chapitres"
	chapterList.SetShowStatusBar(false)

	directoryItems := make([]list.Item, 0, len(entries))
	for i, entry := range entries {
		directoryItems = append(directoryItems, listItem{
			title:       entry.Name,
			description: fmt.Sprintf("%s - %s", entry.Type, entry.Summary),
			index:       i,
		})
	}
	directory := list.New(directoryItems, delegate, 80, 18)
	directory.Title = "Annuaire"
	directory.SetShowStatusBar(false)

	input := textinput.New()
	input.Placeholder = "Votre reponse"
	input.Focus()

	_, err := tea.NewProgram(model{
		screen:      screenMenu,
		menu:        menu,
		chapterList: chapterList,
		directory:   directory,
		input:       input,
		chapters:    chapters,
		entries:     entries,
	}).Run()
	return err
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc", "backspace":
			if m.screen != screenMenu {
				m.screen = screenMenu
				m.feedback = ""
				m.input.SetValue("")
				return m, nil
			}
		case "enter":
			return m.handleEnter()
		}
	case tea.WindowSizeMsg:
		m.menu.SetSize(msg.Width, msg.Height-2)
		m.chapterList.SetSize(msg.Width, msg.Height-2)
		m.directory.SetSize(msg.Width, msg.Height-2)
	}

	var cmd tea.Cmd
	switch m.screen {
	case screenMenu:
		m.menu, cmd = m.menu.Update(msg)
	case screenChapters:
		m.chapterList, cmd = m.chapterList.Update(msg)
	case screenDirectory:
		m.directory, cmd = m.directory.Update(msg)
	case screenTraining:
		m.input, cmd = m.input.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	style := lipgloss.NewStyle().Padding(1, 2)
	var view string
	switch m.screen {
	case screenMenu:
		view = m.menu.View()
	case screenChapters:
		view = m.chapterList.View()
	case screenDirectory:
		view = m.directory.View()
	case screenStats:
		view = m.statsView()
	case screenTraining:
		view = m.trainingView()
	}
	view += "\n\nq: quitter  esc: menu  enter: selectionner"
	return style.Render(view)
}

func (m model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.screen {
	case screenMenu:
		selected, ok := m.menu.SelectedItem().(listItem)
		if !ok {
			return m, nil
		}
		switch selected.title {
		case "Continuer":
			return m.startTraining(0), nil
		case "Choisir un chapitre":
			m.screen = screenChapters
			return m, nil
		case "Annuaire":
			m.screen = screenDirectory
			return m, nil
		case "Stats":
			m.screen = screenStats
			return m, nil
		case "Quitter":
			return m, tea.Quit
		}
	case screenChapters:
		selected, ok := m.chapterList.SelectedItem().(listItem)
		if !ok {
			return m, nil
		}
		return m.startTraining(selected.index), nil
	case screenTraining:
		return m.answerOrAdvance(), nil
	}
	return m, nil
}

func (m model) startTraining(index int) model {
	if len(m.chapters) == 0 {
		m.screen = screenStats
		return m
	}
	if index < 0 || index >= len(m.chapters) {
		index = 0
	}
	m.screen = screenTraining
	m.chapterIndex = index
	m.itemIndex = 0
	m.feedback = ""
	m.input.SetValue("")
	m.input.Focus()
	return m
}

func (m model) answerOrAdvance() model {
	current := m.currentItem()
	if current == nil {
		m.screen = screenMenu
		return m
	}
	if m.feedback != "" {
		m.itemIndex++
		m.feedback = ""
		m.input.SetValue("")
		if m.currentItem() == nil {
			m.screen = screenStats
		}
		return m
	}

	if exercise.MatchAnswer(m.input.Value(), current.Answer) {
		m.feedback = "Correct. " + current.Explanation
	} else {
		m.feedback = "Pas encore. Reponse attendue: " + current.Answer.Primary
	}
	return m
}

func (m model) currentItem() *chapter.Item {
	if m.chapterIndex < 0 || m.chapterIndex >= len(m.chapters) {
		return nil
	}
	chapter := m.chapters[m.chapterIndex]
	if m.itemIndex < 0 || m.itemIndex >= len(chapter.Items) {
		return nil
	}
	return &chapter.Items[m.itemIndex]
}

func (m model) trainingView() string {
	current := m.currentItem()
	if current == nil {
		return "Aucun exercice disponible."
	}
	selected := m.chapters[m.chapterIndex]
	title := lipgloss.NewStyle().Bold(true).Render(selected.Title)
	return fmt.Sprintf(
		"%s\nProgression: %d/%d\n\nQuestion:\n%s\n\nReponse:\n%s\n\n%s",
		title,
		m.itemIndex+1,
		len(selected.Items),
		current.Prompt,
		m.input.View(),
		m.feedback,
	)
}

func (m model) statsView() string {
	total := 0
	for _, chapter := range m.chapters {
		total += len(chapter.Items)
	}
	return fmt.Sprintf("Chapitres: %d\nExercices disponibles: %d\nEntrees annuaire: %d", len(m.chapters), total, len(m.entries))
}
