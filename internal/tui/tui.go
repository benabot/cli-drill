package tui

import (
	"errors"
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/benabot/cli-drill/internal/chapter"
	"github.com/benabot/cli-drill/internal/exercise"
	"github.com/benabot/cli-drill/internal/progress"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type screen int

const (
	screenHome screen = iota
	screenChapters
	screenChapterDetail
	screenTraining
	screenTrainingSummary
	screenStats
	screenPreferences
)

const (
	homeActionStartTraining = "Start training"
	homeActionBrowse        = "Browse chapters"
	homeActionShortcuts     = "Keyboard shortcuts"
	homeActionStats         = "My progress"
	homeActionPreferences   = "Preferences"
	homeActionQuit          = "Quit"
)

const homeTagline = "Learn the terminal by doing."

var homeASCIIBanner = []string{
	"   ____ _     ___               _      _ _ _",
	"  / ___| |   |_ _|           __| |_ __(_) | |",
	" | |   | |    | |   _____   / _` | '__| | | |",
	" | |___| |___ | |  |_____| | (_| | |  | | | |",
	"  \\____|_____|___|          \\__,_|_|  |_|_|_|",
}

type footerContext int

const (
	footerNavigation footerContext = iota
	footerReadOnly
	footerFilterInput
	footerFilteredList
	footerChapterDetail
	footerKeySequenceDetail
	footerTextInput
	footerChoiceInput
	footerTextFeedback
	footerTrainingSummary
	footerTrainingSummaryRetry
	footerPreferences
	footerHome
)

var homeActions = []listItem{
	{title: homeActionStartTraining, description: "Begin a built-in chapter. No setup required."},
	{title: homeActionBrowse, description: "Choose a chapter and inspect its exercises."},
	{title: homeActionShortcuts, description: "Practise keyboard shortcuts with real key capture."},
	{title: homeActionStats, description: "Review local training progress."},
	{title: homeActionPreferences, description: "Choose CLI modules and a default difficulty."},
	{title: homeActionQuit, description: "Close cli-drill."},
}

type Options struct {
	Chapters           []chapter.Chapter
	ProgressPath       string
	LoadPreferences    func() ([]string, chapter.Difficulty, error)
	SavePreferences    func([]string, chapter.Difficulty) error
	CaptureUnavailable string
}

type Result struct {
	CaptureChapter chapter.Chapter
}

// Keep the clipboard's raw text until it passes the same checks as terminal paste.
type trainingPasteMsg struct {
	text string
	err  error
}

func pasteTrainingAnswer() tea.Msg {
	text, err := clipboard.ReadAll()
	return trainingPasteMsg{text: text, err: err}
}

type listItem struct {
	title       string
	description string
	index       int
}

func (i listItem) Title() string       { return i.title }
func (i listItem) Description() string { return i.description }
func (i listItem) FilterValue() string { return i.title + " " + i.description }

type trainingMiss struct {
	Prompt   string
	Expected string
	Answered string
}

type trainingSummary struct {
	Total   int
	Correct int
	Missed  []trainingMiss
}

func (s trainingSummary) retryAvailable() bool {
	return len(s.Missed) > 0
}

type tuiStyleSet struct {
	AppTitle    lipgloss.Style
	ScreenTitle lipgloss.Style
	Subtitle    lipgloss.Style
	Muted       lipgloss.Style
	Info        lipgloss.Style
	Warning     lipgloss.Style
	Footer      lipgloss.Style
	Separator   lipgloss.Style
}

func newTUIStyles() tuiStyleSet {
	return tuiStyleSet{
		AppTitle:    lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("81")),
		ScreenTitle: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86")),
		Subtitle:    lipgloss.NewStyle().Foreground(lipgloss.Color("252")),
		Muted:       lipgloss.NewStyle().Foreground(lipgloss.Color("244")),
		Info:        lipgloss.NewStyle().Foreground(lipgloss.Color("111")),
		Warning:     lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("214")),
		Footer:      lipgloss.NewStyle().Foreground(lipgloss.Color("246")),
		Separator:   lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
	}
}

type model struct {
	screen             screen
	home               list.Model
	chapterList        list.Model
	preferences        preferencesScreen
	input              textinput.Model
	chapters           []chapter.Chapter
	progress           progress.Progress
	progressPath       string
	progressError      string
	chapterIndex       int
	itemIndex          int
	choiceIndex        int
	feedback           string
	summary            trainingSummary
	showHelp           bool
	result             Result
	filter             chapter.SessionFilter
	preferencesError   error
	detailDifficulty   chapter.Difficulty
	detailMode         chapter.SessionMode
	detailScope        chapter.ExerciseScope
	browserNotice      string
	captureUnavailable string
	session            chapter.Chapter
}

func RunWithOptions(opts Options) (Result, error) {
	finalModel, err := tea.NewProgram(newModel(opts)).Run()
	if err != nil {
		return Result{}, err
	}
	if final, ok := finalModel.(model); ok {
		return final.result, nil
	}
	return Result{}, nil
}

func newModel(opts Options) model {
	styles := newTUIStyles()
	delegate := list.NewDefaultDelegate()
	delegate.SetSpacing(0)
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.Foreground(lipgloss.Color("252"))
	delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.Foreground(lipgloss.Color("244"))
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("213")).
		BorderForeground(lipgloss.Color("177"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color("177")).
		BorderForeground(lipgloss.Color("177"))

	homeItems := make([]list.Item, 0, len(homeActions))
	for i, action := range homeActions {
		action.index = i
		homeItems = append(homeItems, action)
	}
	home := list.New(homeItems, delegate, 80, 24)
	home.Title = "Actions"
	home.SetShowStatusBar(false)
	home.SetFilteringEnabled(false)
	home.SetShowHelp(false)
	home.Styles.Title = styles.ScreenTitle

	chapterList := list.New(nil, delegate, 80, 20)
	chapterList.Title = "Chapters"
	chapterList.SetShowStatusBar(false)
	chapterList.SetShowHelp(false)
	chapterList.Styles.Title = styles.ScreenTitle

	input := textinput.New()
	input.Placeholder = "Votre reponse"
	input.Focus()

	localProgress, progressErr := loadProgress(opts.ProgressPath)
	preferences := newPreferencesScreen(opts, delegate)
	return model{
		screen:             screenHome,
		home:               home,
		chapterList:        chapterList,
		preferences:        preferences,
		input:              input,
		chapters:           opts.Chapters,
		progress:           localProgress,
		progressPath:       opts.ProgressPath,
		progressError:      progressErr,
		filter:             chapter.SessionFilter{AvailableTools: preferences.available, Difficulty: chapter.DifficultyBeginner},
		captureUnavailable: opts.CaptureUnavailable,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case trainingPasteMsg:
		if m.screen != screenTraining || m.feedback != "" || m.isMultipleChoice() {
			return m, nil
		}
		if msg.err != nil {
			m.input.Err = msg.err
			return m, nil
		}
		return m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(msg.text), Paste: true})
	case tea.KeyMsg:
		if m.screen == screenPreferences {
			return m.updatePreferences(msg)
		}
		if m.screen == screenChapters && msg.String() != "ctrl+c" {
			if m.chapterList.FilterState() == list.Filtering {
				var cmd tea.Cmd
				m.chapterList, cmd = m.chapterList.Update(msg)
				return m, cmd
			}
			if m.chapterList.FilterState() == list.FilterApplied && msg.String() == "esc" {
				m.chapterList.ResetFilter()
				return m, nil
			}
		}
		if m.screen == screenTrainingSummary {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "g":
				m.goHome()
				return m, nil
			case "esc":
				m.back()
				return m, nil
			case "enter":
				return m.handleEnter()
			}
			return m, nil
		}
		if m.screen == screenTraining {
			if m.feedback == "" && !m.isMultipleChoice() && msg.Type == tea.KeyRunes && strings.ContainsAny(string(msg.Runes), "\r\n") {
				m.input.Err = errors.New("Paste rejected: enter the answer on a single line.")
				return m, nil
			}
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "esc":
				m.back()
				return m, nil
			case "enter":
				return m.handleEnter()
			case "ctrl+v":
				if m.feedback == "" && !m.isMultipleChoice() {
					return m, pasteTrainingAnswer
				}
			}
			if m.feedback == "" {
				if m.isMultipleChoice() {
					if !msg.Paste {
						switch msg.String() {
						case "up", "k":
							m.choiceIndex = max(0, m.choiceIndex-1)
						case "down", "j":
							m.choiceIndex = max(0, min(len(m.currentItem().Choices)-1, m.choiceIndex+1))
						}
					}
					return m, nil
				}
				var cmd tea.Cmd
				m.input, cmd = m.input.Update(msg)
				return m, cmd
			}
			return m, nil
		}
		if m.screen == screenHome && !msg.Paste && msg.String() == "p" {
			m.openPreferences()
			return m, nil
		}
		if m.screen == screenChapterDetail && !msg.Paste && m.changeSessionChoice(msg.String()) {
			return m, nil
		}
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "h":
			m.showHelp = !m.showHelp
			return m, nil
		case "g":
			m.goHome()
			return m, nil
		case "esc", "backspace":
			m.back()
			return m, nil
		case "enter":
			return m.handleEnter()
		}
	case tea.WindowSizeMsg:
		// Reserve outer padding, the chapter heading and the contextual footer.
		width := max(1, msg.Width-4)
		m.home.SetSize(width, homeListHeight(msg.Height))
		m.chapterList.SetSize(width, max(1, msg.Height-9))
		m.preferences.resize(width, msg.Height)
	}

	var cmd tea.Cmd
	switch m.screen {
	case screenHome:
		m.home, cmd = m.home.Update(msg)
	case screenChapters:
		m.chapterList, cmd = m.chapterList.Update(msg)
	case screenTraining:
		m.input, cmd = m.input.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	style := lipgloss.NewStyle().Padding(1, 2)
	if m.screen == screenPreferences {
		return style.Render(m.preferencesView())
	}
	var view string
	switch m.screen {
	case screenHome:
		view = m.homeView()
	case screenChapters:
		view = m.chaptersView()
	case screenChapterDetail:
		view = m.chapterDetailView()
	case screenStats:
		view = m.statsView()
	case screenTraining:
		view = m.trainingView()
	case screenTrainingSummary:
		view = m.trainingSummaryView()
	}
	if m.showHelp {
		view += "\n\n" + m.helpView()
	}
	if m.progressError != "" {
		view += "\n\n" + newTUIStyles().Warning.Render(m.progressError)
	}
	if m.preferencesError != nil {
		view += "\n\n" + newTUIStyles().Warning.Render(fmt.Sprintf("Preferences unavailable. Using defaults. %v", m.preferencesError))
	}
	view += "\n\n" + m.footer()
	return style.Render(view)
}

func (m model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.screen {
	case screenHome:
		selected, ok := m.home.SelectedItem().(listItem)
		if !ok {
			return m, nil
		}
		switch selected.title {
		case homeActionStartTraining:
			m.openChapters(false)
			i := m.filter.FirstTextChapter(m.chapters)
			if i < 0 {
				m.browserNotice = fmt.Sprintf("No eligible text exercises at %s difficulty.\nCheck enabled CLI modules or open a chapter to choose another level.", m.filter.Difficulty)
				return m, nil
			}
			m.chapterList.Select(i)
			m.openDetail(i, chapter.TextExercises)
			return m.launchSession()
		case homeActionBrowse:
			m.openChapters(false)
			return m, nil
		case homeActionShortcuts:
			m.openChapters(true)
			if items := m.chapterList.Items(); len(items) == 1 {
				m.openDetail(items[0].(listItem).index, chapter.KeyExercises)
				return m.launchSession()
			}
			return m, nil
		case homeActionStats:
			m.refreshSessionFilter()
			m.screen = screenStats
			return m, nil
		case homeActionPreferences:
			m.openPreferences()
			return m, nil
		case homeActionQuit:
			return m, tea.Quit
		}
	case screenChapters:
		selected, ok := m.chapterList.SelectedItem().(listItem)
		if !ok {
			return m, nil
		}
		scope := chapter.AllExercises
		if m.chapterList.Title == homeActionShortcuts {
			scope = chapter.KeyExercises
		}
		m.openDetail(selected.index, scope)
		return m, nil
	case screenChapterDetail:
		m.refreshSessionFilter()
		return m.launchSession()
	case screenTraining:
		return m.answerOrAdvance(), nil
	case screenTrainingSummary:
		if m.summary.retryAvailable() {
			m.resetTextTraining()
			return m, nil
		}
	}
	return m, nil
}

func (m *model) openChapters(shortcutsOnly bool) {
	m.refreshSessionFilter()
	m.browserNotice = ""
	scope := chapter.AllExercises
	if shortcutsOnly {
		scope = chapter.KeyExercises
	}
	items := make([]list.Item, 0, len(m.chapters))
	for i, current := range m.chapters {
		count := len(m.filter.EligibleItems(current, scope))
		if shortcutsOnly && count == 0 {
			continue
		}
		items = append(items, listItem{
			title:       current.Title,
			description: chapterListDescription(current, count),
			index:       i,
		})
	}
	m.chapterList.ResetFilter()
	m.chapterList.SetItems(items)
	m.chapterList.Select(0)
	m.chapterList.Title = "Chapters"
	if shortcutsOnly {
		m.chapterList.Title = homeActionShortcuts
		if len(items) == 0 {
			m.browserNotice = fmt.Sprintf("No eligible keyboard shortcuts at %s difficulty.\nCheck enabled CLI modules or choose another level via Browse chapters.", m.filter.Difficulty)
		}
	}
	m.screen = screenChapters
	m.feedback = ""
	m.showHelp = false
}

func (m *model) back() {
	switch m.screen {
	case screenHome:
		return
	case screenChapterDetail:
		m.clearSession()
		m.screen = screenChapters
	case screenTraining, screenTrainingSummary:
		m.clearSession()
		m.refreshSessionFilter()
		m.screen = screenChapterDetail
	default:
		m.screen = screenHome
	}
}

func (m *model) goHome() {
	m.clearSession()
	m.screen = screenHome
}

func (m model) answerOrAdvance() model {
	current := m.currentItem()
	if current == nil {
		m.screen = screenHome
		return m
	}
	if m.feedback != "" {
		m.itemIndex++
		m.choiceIndex = 0
		m.feedback = ""
		m.input.SetValue("")
		if m.currentItem() == nil {
			if m.summary.retryAvailable() {
				m.screen = screenTrainingSummary
			} else {
				m.goHome()
			}
		}
		return m
	}

	answered := m.input.Value()
	if m.isMultipleChoice() {
		if m.choiceIndex < 0 || m.choiceIndex >= len(current.Choices) {
			return m
		}
		answered = current.Choices[m.choiceIndex]
	} else if m.input.Err != nil || strings.Trim(answered, " \r\n") == "" {
		return m
	}
	correct := current.MatchAnswer(answered)
	recorded, err := progress.RecordAnswer(m.progressPath, m.session.ID, current.ID, correct)
	if err != nil {
		m.progressError = "Progress could not be saved: " + err.Error()
		return m
	}
	m.progress = recorded
	m.progressError = ""
	m.summary.Total++
	if correct {
		m.summary.Correct++
		m.feedback = "Correct. " + current.Explanation
	} else {
		m.summary.Missed = append(m.summary.Missed, trainingMiss{
			Prompt:   current.Prompt,
			Expected: current.Answer.Primary,
			Answered: answered,
		})
		m.feedback = "Pas encore. Reponse attendue: " + current.Answer.Primary
	}
	return m
}

func (m model) currentItem() *chapter.Item {
	if m.itemIndex < 0 || m.itemIndex >= len(m.session.Items) {
		return nil
	}
	return &m.session.Items[m.itemIndex]
}

func (m model) isMultipleChoice() bool {
	current := m.currentItem()
	return current != nil && current.ExerciseType == exercise.TypeMultipleChoice
}

func (m model) currentChapter() *chapter.Chapter {
	if m.chapterIndex < 0 || m.chapterIndex >= len(m.chapters) {
		return nil
	}
	return &m.chapters[m.chapterIndex]
}

func (m model) homeView() string {
	styles := newTUIStyles()
	lines := []string{
		renderHomeBanner(m.home.Width(), styles),
		styles.Subtitle.Render(homeTagline),
		"",
	}
	lines = append(lines,
		styles.Separator.Render("────────────────────────────────────────"),
		"",
		m.home.View(),
	)
	return strings.Join(lines, "\n")
}

func renderHomeBanner(width int, styles tuiStyleSet) string {
	if width > 0 && width < homeBannerWidth() {
		return styles.AppTitle.Render("cli-drill")
	}
	return styles.AppTitle.Render(strings.Join(homeASCIIBanner, "\n"))
}

func homeBannerWidth() int {
	width := 0
	for _, line := range homeASCIIBanner {
		width = max(width, lipgloss.Width(line))
	}
	return width
}

func homeListHeight(windowHeight int) int {
	const homeChromeHeight = 14
	return max(7, windowHeight-homeChromeHeight)
}

func (m model) chaptersView() string {
	styles := newTUIStyles()
	if m.browserNotice != "" {
		return strings.Join([]string{
			styles.ScreenTitle.Render(m.chapterList.Title),
			styles.Warning.Render(m.browserNotice),
			m.chapterList.View(),
		}, "\n\n")
	}
	return strings.Join([]string{
		styles.ScreenTitle.Render(m.chapterList.Title),
		styles.Subtitle.Render(fmt.Sprintf("Choose a chapter to inspect · Default difficulty: %s", m.filter.Difficulty)),
		styles.Separator.Render("────────────────────────────────────────"),
		"",
		m.chapterList.View(),
	}, "\n")
}

func (m model) chapterDetailView() string {
	styles := newTUIStyles()
	current := m.currentChapter()
	if current == nil {
		return "No chapter selected."
	}
	score := m.progress.ChapterScores[current.ID]
	filter := m.detailFilter()
	eligible := filter.EligibleItems(*current, m.detailScope)
	acquired := chapter.AcquiredCount(current.ID, eligible, m.progress.CompletedExercises)
	lines := []string{
		styles.ScreenTitle.Render(current.Title),
		styles.Subtitle.Render(current.Description),
		styles.Separator.Render("────────────────────────────────────────"),
		"",
		"Chapter: " + current.ID,
		fmt.Sprintf("Exercises: %d in corpus", len(current.Items)),
		fmt.Sprintf("Session difficulty: %s · d change", filter.Difficulty),
		sessionLevelCounts(*current, m.detailScope, filter),
		fmt.Sprintf("Active acquired: %d/%d eligible", acquired, len(eligible)),
		fmt.Sprintf("Historical chapter score: %d/%d correct", score.Correct, score.Attempts),
		"",
	}
	if current.Guided() {
		mode := "Learn (editorial order)"
		if m.detailMode == chapter.ModeReview {
			mode = "Review (shuffled)"
		}
		lines = append(lines, "Mode: "+mode+" · m change")
	} else {
		lines = append(lines, "Order: shuffled")
	}
	prepared := *current
	prepared.Items = eligible
	switch {
	case len(eligible) == 0:
		lines = append(lines, styles.Warning.Render(filter.EmptyReason(*current, m.detailScope)))
	case prepared.HasKeyExercises():
		lines = append(lines,
			styles.Info.Render("Ce chapitre entraîne les vrais raccourcis clavier."),
			"cli-drill ouvre un mode dédié pour capturer les touches.",
		)
		if m.captureUnavailable != "" {
			lines = append(lines, styles.Warning.Render(m.captureUnavailable))
		}
	default:
		lines = append(lines, "Press enter to start this chapter.")
	}
	return strings.Join(lines, "\n")
}

func (m model) trainingView() string {
	styles := newTUIStyles()
	current := m.currentItem()
	if current == nil {
		return "Aucun exercice disponible."
	}
	selected := m.session
	title := styles.ScreenTitle.Render(selected.Title)
	inputView := m.input.View()
	answerLabel := "Reponse:"
	feedback := m.feedback
	if m.isMultipleChoice() {
		answerLabel = "Choix:"
		choices := make([]string, 0, len(current.Choices))
		for i, choice := range current.Choices {
			if i == m.choiceIndex {
				choices = append(choices, styles.Info.Render("> "+choice))
			} else {
				choices = append(choices, "  "+choice)
			}
		}
		inputView = strings.Join(choices, "\n")
		if strings.HasPrefix(feedback, "Pas encore.") && current.Explanation != "" {
			feedback += "\n" + current.Explanation
		}
	} else if m.input.Err != nil {
		inputView += "\n" + styles.Warning.Render(m.input.Err.Error())
	}
	return fmt.Sprintf(
		"%s\n%s\n%s\n\n%s\n%s\n\n%s\n%s\n\n%s",
		title,
		styles.Muted.Render(fmt.Sprintf("Progression: %d/%d", m.itemIndex+1, len(selected.Items))),
		styles.Separator.Render("────────────────────────────────────────"),
		styles.Muted.Render("Question:"),
		current.Prompt,
		styles.Muted.Render(answerLabel),
		inputView,
		feedbackText(feedback, styles),
	)
}

func (m model) trainingSummaryView() string {
	styles := newTUIStyles()
	title := "Training"
	if current := m.currentChapter(); current != nil {
		title = current.Title
	}
	lines := []string{
		styles.ScreenTitle.Render("Chapter complete"),
		styles.Subtitle.Render(title),
		styles.Separator.Render("────────────────────────────────────────"),
		"",
		fmt.Sprintf("Score: %d/%d correct", m.summary.Correct, m.summary.Total),
	}
	if len(m.summary.Missed) == 0 {
		return strings.Join(lines, "\n")
	}

	lines = append(lines, "", styles.Warning.Render("Missed:"))
	for i, missed := range m.summary.Missed {
		if i == 3 {
			lines = append(lines, fmt.Sprintf("And %d more.", len(m.summary.Missed)-i))
			break
		}
		lines = append(lines,
			missed.Prompt,
			styles.Muted.Render("Expected: "+missed.Expected),
			styles.Muted.Render("Answered: "+valueOrFallback(missed.Answered, "(empty)")),
		)
	}
	return strings.Join(lines, "\n")
}

func (m model) statsView() string {
	styles := newTUIStyles()
	acquired, total := m.filter.ActiveCounts(m.chapters, m.progress.CompletedExercises)
	lines := []string{
		styles.ScreenTitle.Render("My progress"),
		styles.Separator.Render("────────────────────────────────────────"),
		fmt.Sprintf("Active acquired (%s, enabled CLI modules): %d/%d", m.filter.Difficulty, acquired, total),
		fmt.Sprintf("Historical acquired: %d", len(m.progress.CompletedExercises)),
		fmt.Sprintf("Chapters with historical scores: %d", len(m.progress.ChapterScores)),
	}
	return strings.Join(lines, "\n")
}

func (m model) helpView() string {
	styles := newTUIStyles()
	return strings.Join([]string{
		styles.ScreenTitle.Render("Help"),
		"enter select/start",
		"j/k or arrows move in lists",
		"esc back",
		"g home",
		"q quit",
	}, "\n")
}

func (m model) footer() string {
	styles := newTUIStyles()
	if m.screen == screenPreferences {
		return styles.Footer.Render(m.preferences.footerText())
	}
	return styles.Footer.Render(footerText(m.footerContext()))
}

func (m model) footerContext() footerContext {
	if m.screen == screenChapters {
		switch m.chapterList.FilterState() {
		case list.Filtering:
			return footerFilterInput
		case list.FilterApplied:
			return footerFilteredList
		}
	}
	switch m.screen {
	case screenHome:
		return footerHome
	case screenChapters:
		return footerNavigation
	case screenStats:
		return footerReadOnly
	case screenPreferences:
		return footerPreferences
	case screenChapterDetail:
		if current := m.currentChapter(); current != nil {
			selected := *current
			selected.Items = m.detailFilter().EligibleItems(selected, m.detailScope)
			if len(selected.Items) == 0 || (selected.HasKeyExercises() && m.captureUnavailable != "") {
				return footerReadOnly
			}
			if selected.HasKeyExercises() {
				return footerKeySequenceDetail
			}
		}
		return footerChapterDetail
	case screenTraining:
		if m.feedback != "" {
			return footerTextFeedback
		}
		if m.isMultipleChoice() {
			return footerChoiceInput
		}
		return footerTextInput
	case screenTrainingSummary:
		if m.summary.retryAvailable() {
			return footerTrainingSummaryRetry
		}
		return footerTrainingSummary
	default:
		return footerNavigation
	}
}

func footerText(context footerContext) string {
	switch context {
	case footerHome:
		return "enter select · j/k · p preferences · esc back · g home · h help · q quit"
	case footerPreferences:
		return "↑/↓ j/k · enter change · s save · esc back"
	case footerReadOnly:
		return "esc back · g home · h help · q quit"
	case footerFilterInput:
		return "type to filter · enter apply · esc cancel · ctrl+c quit"
	case footerFilteredList:
		return "enter select · j/k move · esc clear filter · g home · h help · q quit"
	case footerChapterDetail:
		return "enter start training · esc chapters · g home · h help · q quit"
	case footerKeySequenceDetail:
		return "enter start key training · esc chapters · g home · h help · q quit"
	case footerTextInput:
		return "enter submit · esc chapter · ctrl+c quit"
	case footerChoiceInput:
		return "↑/↓ j/k move · enter submit · esc chapter · ctrl+c quit"
	case footerTextFeedback:
		return "enter next · esc chapter · ctrl+c quit"
	case footerTrainingSummaryRetry:
		return "enter retry · esc chapter · g home · q quit"
	case footerTrainingSummary:
		return "esc chapter · g home · q quit"
	default:
		return "enter select · j/k move · esc back · g home · h help · q quit"
	}
}

func feedbackText(value string, styles tuiStyleSet) string {
	if value == "" {
		return ""
	}
	if strings.HasPrefix(value, "Correct.") {
		return styles.Info.Render(value)
	}
	if strings.HasPrefix(value, "Pas encore.") {
		return styles.Warning.Render(value)
	}
	return value
}

func loadProgress(path string) (progress.Progress, string) {
	if path == "" {
		return progress.Default(), ""
	}
	localProgress, err := progress.Load(path)
	if err != nil {
		return progress.Default(), "Progress could not be loaded: " + err.Error()
	}
	return localProgress, ""
}

func valueOrFallback(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func chapterListDescription(chapter chapter.Chapter, eligible int) string {
	count := exerciseCountLabel(len(chapter.Items))
	if eligible != len(chapter.Items) || eligible == 0 {
		count = fmt.Sprintf("%d/%d eligible", eligible, len(chapter.Items))
	}
	metadata := chapter.ID + " · " + count
	if chapter.Description == "" {
		return metadata
	}
	return metadata + " · " + chapter.Description
}

func exerciseCountLabel(count int) string {
	if count == 1 {
		return "1 exercise"
	}
	return fmt.Sprintf("%d exercises", count)
}
