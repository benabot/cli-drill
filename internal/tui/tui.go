package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/benabot/cli-drill/internal/catalog"
	"github.com/benabot/cli-drill/internal/chapter"
	"github.com/benabot/cli-drill/internal/config"
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
	screenConfig
	screenScan
	screenChapters
	screenChapterDetail
	screenTraining
	screenDirectory
	screenStats
)

const (
	homeActionStartTraining = "Start training"
	homeActionBrowse        = "Browse chapters"
	homeActionScan          = "Scan dotfiles"
	homeActionDirectory     = "Directory"
	homeActionStats         = "Stats"
	homeActionConfig        = "Configuration"
	homeActionQuit          = "Quit"
)

const (
	footerNavigation        = "enter select · j/k move · esc back · g home · h help · q quit"
	footerChapterDetail     = "enter start training · esc chapters · g home · h help · q quit"
	footerKeySequenceDetail = "enter start key training · esc chapters · g home · h help · q quit"
	footerTraining          = "enter submit/next · esc chapter · g home · h help · q quit"
	footerStatic            = "esc back · g home · h help · q quit"
)

var homeActions = []listItem{
	{title: homeActionStartTraining, description: "Begin the first available chapter"},
	{title: homeActionBrowse, description: "Pick a chapter and inspect its exercises"},
	{title: homeActionScan, description: "Read-only guidance for scanning configured dotfiles"},
	{title: homeActionDirectory, description: "Browse embedded shortcuts, aliases, tools, and workflows"},
	{title: homeActionStats, description: "Review local progress"},
	{title: homeActionConfig, description: "Check dotfiles configuration status"},
	{title: homeActionQuit, description: "Close cli-drill"},
}

type Options struct {
	Chapters     []chapter.Chapter
	Entries      []catalog.Entry
	ConfigPath   string
	ProgressPath string
}

type Result struct {
	RunCLITrainChapterID string
}

type listItem struct {
	title       string
	description string
	index       int
}

func (i listItem) Title() string       { return i.title }
func (i listItem) Description() string { return i.description }
func (i listItem) FilterValue() string { return i.title + " " + i.description }

type configStatus struct {
	Path              string
	Exists            bool
	LoadError         string
	DotfilesPath      string
	DotfilesPathCheck string
	DotfilesExists    bool
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
	screen        screen
	home          list.Model
	chapterList   list.Model
	directory     list.Model
	input         textinput.Model
	chapters      []chapter.Chapter
	entries       []catalog.Entry
	config        configStatus
	progress      progress.Progress
	progressError string
	chapterIndex  int
	itemIndex     int
	feedback      string
	showHelp      bool
	result        Result
}

func Run(chapters []chapter.Chapter, entries []catalog.Entry) error {
	_, err := RunWithOptions(Options{Chapters: chapters, Entries: entries})
	return err
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
	home.Styles.Title = styles.ScreenTitle

	chapterItems := make([]list.Item, 0, len(opts.Chapters))
	for i, chapter := range opts.Chapters {
		chapterItems = append(chapterItems, listItem{
			title:       chapter.Title,
			description: fmt.Sprintf("%s - %d exercises", chapter.ID, len(chapter.Items)),
			index:       i,
		})
	}
	chapterList := list.New(chapterItems, delegate, 80, 20)
	chapterList.Title = "Chapters"
	chapterList.SetShowStatusBar(false)
	chapterList.Styles.Title = styles.ScreenTitle

	directoryItems := make([]list.Item, 0, len(opts.Entries))
	for i, entry := range opts.Entries {
		directoryItems = append(directoryItems, listItem{
			title:       entry.Name,
			description: fmt.Sprintf("%s - %s", entry.Type, entry.Summary),
			index:       i,
		})
	}
	directory := list.New(directoryItems, delegate, 80, 20)
	directory.Title = "Directory"
	directory.SetShowStatusBar(false)
	directory.Styles.Title = styles.ScreenTitle

	input := textinput.New()
	input.Placeholder = "Votre reponse"
	input.Focus()

	localProgress, progressErr := loadProgress(opts.ProgressPath)
	return model{
		screen:        screenHome,
		home:          home,
		chapterList:   chapterList,
		directory:     directory,
		input:         input,
		chapters:      opts.Chapters,
		entries:       opts.Entries,
		config:        loadConfigStatus(opts.ConfigPath),
		progress:      localProgress,
		progressError: progressErr,
	}
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
		m.home.SetSize(msg.Width, msg.Height-2)
		m.chapterList.SetSize(msg.Width, msg.Height-2)
		m.directory.SetSize(msg.Width, msg.Height-2)
	}

	var cmd tea.Cmd
	switch m.screen {
	case screenHome:
		m.home, cmd = m.home.Update(msg)
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
	case screenHome:
		view = m.homeView()
	case screenConfig:
		view = m.configView()
	case screenScan:
		view = m.scanView()
	case screenChapters:
		view = m.chaptersView()
	case screenChapterDetail:
		view = m.chapterDetailView()
	case screenDirectory:
		view = m.directoryView()
	case screenStats:
		view = m.statsView()
	case screenTraining:
		view = m.trainingView()
	}
	if m.showHelp {
		view += "\n\n" + m.helpView()
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
			return m.startTraining(0), nil
		case homeActionBrowse:
			m.screen = screenChapters
			m.feedback = ""
			return m, nil
		case homeActionScan:
			m.screen = screenScan
			return m, nil
		case homeActionDirectory:
			m.screen = screenDirectory
			return m, nil
		case homeActionStats:
			m.screen = screenStats
			return m, nil
		case homeActionConfig:
			m.screen = screenConfig
			return m, nil
		case homeActionQuit:
			return m, tea.Quit
		}
	case screenChapters:
		selected, ok := m.chapterList.SelectedItem().(listItem)
		if !ok {
			return m, nil
		}
		m.chapterIndex = selected.index
		m.screen = screenChapterDetail
		return m, nil
	case screenChapterDetail:
		if current := m.currentChapter(); current != nil && chapterUsesKeySequence(*current) {
			m.result.RunCLITrainChapterID = current.ID
			return m, tea.Quit
		}
		return m.startTraining(m.chapterIndex), nil
	case screenTraining:
		return m.answerOrAdvance(), nil
	}
	return m, nil
}

func (m *model) back() {
	switch m.screen {
	case screenHome:
		return
	case screenChapterDetail:
		m.screen = screenChapters
	case screenTraining:
		m.screen = screenChapterDetail
		m.feedback = ""
		m.input.SetValue("")
	default:
		m.screen = screenHome
	}
}

func (m *model) goHome() {
	m.screen = screenHome
	m.feedback = ""
	m.input.SetValue("")
}

func (m model) startTraining(index int) model {
	if len(m.chapters) == 0 {
		m.screen = screenStats
		return m
	}
	if index < 0 || index >= len(m.chapters) {
		index = 0
	}
	m.chapterIndex = index
	if chapterUsesKeySequence(m.chapters[index]) {
		m.screen = screenChapterDetail
		m.feedback = ""
		m.input.SetValue("")
		return m
	}
	m.screen = screenTraining
	m.itemIndex = 0
	m.feedback = ""
	m.input.SetValue("")
	m.input.Focus()
	return m
}

func (m model) answerOrAdvance() model {
	current := m.currentItem()
	if current == nil {
		m.screen = screenHome
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

func (m model) currentChapter() *chapter.Chapter {
	if m.chapterIndex < 0 || m.chapterIndex >= len(m.chapters) {
		return nil
	}
	return &m.chapters[m.chapterIndex]
}

func (m model) homeView() string {
	styles := newTUIStyles()
	return strings.Join([]string{
		styles.AppTitle.Render("cli-drill"),
		styles.Subtitle.Render("Train terminal shortcuts, aliases, shell tools, and dotfiles workflows."),
		"",
		styles.Info.Render(m.configSummaryLine()),
		styles.Separator.Render("────────────────────────────────────────"),
		"",
		m.home.View(),
	}, "\n")
}

func (m model) configView() string {
	styles := newTUIStyles()
	lines := []string{
		styles.ScreenTitle.Render("Configuration"),
		styles.Subtitle.Render("Read-only status. No config file is created from the TUI."),
		styles.Separator.Render("────────────────────────────────────────"),
		"",
		"Config file: " + valueOrFallback(m.config.Path, "unknown"),
		styles.Info.Render("Status: " + m.configStatusLabel()),
	}
	if m.config.LoadError != "" {
		lines = append(lines, styles.Warning.Render("Error: "+m.config.LoadError))
	}
	if m.config.DotfilesPath != "" {
		lines = append(lines, "Dotfiles path: "+m.config.DotfilesPath)
		if m.config.DotfilesPathCheck != "" {
			lines = append(lines, "Checked path: "+m.config.DotfilesPathCheck)
		}
		lines = append(lines, "Path exists: "+yesNo(m.config.DotfilesExists))
	}
	lines = append(lines,
		"",
		"To preview a starter config, run:",
		"cli-drill init --print",
	)
	return strings.Join(lines, "\n")
}

func (m model) scanView() string {
	styles := newTUIStyles()
	return strings.Join([]string{
		styles.ScreenTitle.Render("Scan dotfiles"),
		styles.Subtitle.Render("This screen is read-only in this pass."),
		styles.Separator.Render("────────────────────────────────────────"),
		"",
		"No dotfiles scan is launched here.",
		"Use the CLI when you want to inspect configured dotfiles:",
		"cli-drill scan --summary",
	}, "\n")
}

func (m model) chaptersView() string {
	styles := newTUIStyles()
	return strings.Join([]string{
		styles.ScreenTitle.Render("Chapters"),
		styles.Subtitle.Render("Choose a chapter to inspect before training."),
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
	lines := []string{
		styles.ScreenTitle.Render(current.Title),
		styles.Subtitle.Render(current.Description),
		styles.Separator.Render("────────────────────────────────────────"),
		"",
		"Chapter: " + current.ID,
		fmt.Sprintf("Exercises: %d", len(current.Items)),
		fmt.Sprintf("Progress: %d/%d completed", m.completedCount(current.ID), len(current.Items)),
		fmt.Sprintf("Score: %d/%d correct", score.Correct, score.Attempts),
		"",
	}
	if chapterUsesKeySequence(*current) {
		lines = append(lines,
			styles.Info.Render("Ce chapitre entraîne les vrais raccourcis clavier."),
			"cli-drill ouvre un mode dédié pour capturer les touches.",
		)
	} else {
		lines = append(lines, "Press enter to start this chapter.")
	}
	return strings.Join(lines, "\n")
}

func (m model) directoryView() string {
	styles := newTUIStyles()
	return strings.Join([]string{
		styles.ScreenTitle.Render("Directory"),
		styles.Subtitle.Render("Embedded entries from chapters. Dotfiles scan is not launched from here."),
		styles.Separator.Render("────────────────────────────────────────"),
		"",
		m.directory.View(),
	}, "\n")
}

func (m model) trainingView() string {
	styles := newTUIStyles()
	current := m.currentItem()
	if current == nil {
		return "Aucun exercice disponible."
	}
	selected := m.chapters[m.chapterIndex]
	title := styles.ScreenTitle.Render(selected.Title)
	return fmt.Sprintf(
		"%s\n%s\n%s\n\n%s\n%s\n\n%s\n%s\n\n%s",
		title,
		styles.Muted.Render(fmt.Sprintf("Progression: %d/%d", m.itemIndex+1, len(selected.Items))),
		styles.Separator.Render("────────────────────────────────────────"),
		styles.Muted.Render("Question:"),
		current.Prompt,
		styles.Muted.Render("Reponse:"),
		m.input.View(),
		feedbackText(m.feedback, styles),
	)
}

func (m model) statsView() string {
	styles := newTUIStyles()
	total := 0
	completed := 0
	for _, chapter := range m.chapters {
		total += len(chapter.Items)
		completed += m.completedCount(chapter.ID)
	}
	lines := []string{
		styles.ScreenTitle.Render("Stats"),
		styles.Separator.Render("────────────────────────────────────────"),
		fmt.Sprintf("Chapters: %d", len(m.chapters)),
		fmt.Sprintf("Exercises available: %d", total),
		fmt.Sprintf("Completed locally: %d/%d", completed, total),
		fmt.Sprintf("Directory entries: %d", len(m.entries)),
	}
	if m.progressError != "" {
		lines = append(lines, "", "Progress could not be loaded: "+m.progressError)
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
	var footer string
	switch m.screen {
	case screenHome:
		footer = footerNavigation
	case screenChapters, screenDirectory:
		footer = footerNavigation
	case screenChapterDetail:
		if current := m.currentChapter(); current != nil && chapterUsesKeySequence(*current) {
			footer = footerKeySequenceDetail
		} else {
			footer = footerChapterDetail
		}
	case screenTraining:
		footer = footerTraining
	default:
		footer = footerStatic
	}
	return styles.Footer.Render(footer)
}

func (m model) configSummaryLine() string {
	if m.config.Path == "" {
		return "Dotfiles: config path unknown"
	}
	if !m.config.Exists {
		return "Dotfiles: not configured"
	}
	if m.config.LoadError != "" {
		return "Dotfiles: config error"
	}
	if m.config.DotfilesPath == "" {
		return "Dotfiles: configured without dotfiles_path"
	}
	if m.config.DotfilesExists {
		return "Dotfiles: " + m.config.DotfilesPath
	}
	return "Dotfiles: " + m.config.DotfilesPath + " (missing)"
}

func (m model) configStatusLabel() string {
	if !m.config.Exists {
		return "not configured"
	}
	if m.config.LoadError != "" {
		return "error"
	}
	return "configured"
}

func (m model) completedCount(chapterID string) int {
	prefix := chapterID + "/"
	count := 0
	for _, completed := range m.progress.CompletedExercises {
		if strings.HasPrefix(completed, prefix) {
			count++
		}
	}
	return count
}

func chapterUsesKeySequence(chapter chapter.Chapter) bool {
	for _, item := range chapter.Items {
		if item.ExerciseType == exercise.TypeKeySequence {
			return true
		}
	}
	return false
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

func loadConfigStatus(path string) configStatus {
	status := configStatus{Path: path}
	if path == "" {
		return status
	}
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return status
		}
		status.LoadError = err.Error()
		return status
	}
	status.Exists = true

	cfg, err := config.Load(path)
	if err != nil {
		status.LoadError = err.Error()
		return status
	}
	if err := cfg.Validate(); err != nil {
		status.LoadError = err.Error()
	}
	status.DotfilesPath = cfg.DotfilesPath
	status.DotfilesPathCheck = expandHome(cfg.DotfilesPath)
	if status.DotfilesPathCheck == "" {
		return status
	}
	if _, err := os.Stat(status.DotfilesPathCheck); err == nil {
		status.DotfilesExists = true
	}
	return status
}

func loadProgress(path string) (progress.Progress, string) {
	if path == "" {
		return progress.Default(), ""
	}
	localProgress, err := progress.Load(path)
	if err != nil {
		return progress.Default(), err.Error()
	}
	return localProgress, ""
}

func homeActionIndex(title string) int {
	for i, action := range homeActions {
		if action.title == title {
			return i
		}
	}
	return 0
}

func expandHome(path string) string {
	if path == "" {
		return ""
	}
	if path == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return home
	}
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, strings.TrimPrefix(path, "~/"))
	}
	return path
}

func valueOrFallback(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

func titleStyle() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true)
}
