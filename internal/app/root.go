package app

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/benabot/cli-drill/internal/chapter"
	"github.com/benabot/cli-drill/internal/exercise"
	"github.com/benabot/cli-drill/internal/progress"
	"github.com/benabot/cli-drill/internal/tui"
	"github.com/benabot/cli-drill/internal/xdg"
	"github.com/charmbracelet/x/term"
	"github.com/spf13/cobra"
)

type Options struct {
	ProgressPath string
	In           io.Reader
	Out          io.Writer
	DefaultFS    fs.FS
	AllExercises bool
}

var (
	isTerminalFD           = term.IsTerminal
	terminalSize           = term.GetSize
	checkTUITerminal       = tuiTerminalError
	readRawKeyFromTerminal = readRawKey
	runTUIWithOptions      = tui.RunWithOptions
	canCaptureKeys         = terminalCanCaptureKeys
	errCaptureUnavailable  = errors.New("Keyboard capture requires a terminal for input and output. Choose a text chapter to continue.")
)

func NewRootCommand(defaultFS fs.FS) *cobra.Command {
	opts := Options{
		In:        os.Stdin,
		Out:       os.Stdout,
		DefaultFS: defaultFS,
	}

	root := &cobra.Command{
		Use:   "cli-drill",
		Short: "Learn the terminal by doing",
		Long:  "Learn the terminal by doing. Practise commands and keyboard shortcuts with built-in chapters, locally and offline.",
		RunE: func(cmd *cobra.Command, args []string) error {
			runOpts := commandOptions(cmd, opts)
			if err := checkTUITerminal(runOpts); err != nil {
				cmd.SilenceUsage, cmd.SilenceErrors = true, true
				return err
			}
			if err := checkProgressOverride(runOpts); err != nil {
				return err
			}
			chapters, err := loadChapters(runOpts)
			if err != nil {
				return err
			}
			for {
				ui := tuiOptions(chapters, progressPath(runOpts), xdg.PreferencesFile())
				ui.In, ui.Out = runOpts.In, runOpts.Out
				if !canCaptureKeys(runOpts) {
					ui.CaptureUnavailable = errCaptureUnavailable.Error()
				}
				result, err := runTUIWithOptions(ui)
				if err != nil {
					cmd.SilenceUsage, cmd.SilenceErrors = true, true
					return fmt.Errorf("could not start the interactive interface. Run cli-drill train for text training: %w", err)
				}
				if result.CaptureChapter.ID == "" {
					return nil
				}
				runOpts.AllExercises = result.AllExercises
				if err := runTraining(runOpts, result.CaptureChapter); err != nil {
					return err
				}
			}
		},
	}
	root.PersistentFlags().StringVar(&opts.ProgressPath, "progress", "", "path to progress JSON")

	root.AddCommand(newChaptersCommand(&opts))
	root.AddCommand(newTrainCommand(&opts))
	root.AddCommand(newStatsCommand(&opts))
	root.AddCommand(newResetCommand(&opts))
	root.AddCommand(newPreferencesCommand(&opts))

	return root
}

func newChaptersCommand(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "chapters",
		Short: "List available chapters",
		RunE: func(cmd *cobra.Command, args []string) error {
			runOpts := commandOptions(cmd, *opts)
			chapters, err := loadChapters(runOpts)
			if err != nil {
				return err
			}
			return writeChapters(runOpts.Out, chapters)
		},
	}
}

func newTrainCommand(opts *Options) *cobra.Command {
	var level, mode string
	var all bool
	cmd := &cobra.Command{
		Use:   "train [chapter-id]",
		Short: "Train from a chapter in the terminal",
		Long: "Train in batches of up to six. Learn focuses on exercises you have not acquired. " +
			"Review mixes all eligible exercises. Use --all for the full filtered chapter in one batch. " +
			"Without an ID, start the first eligible text chapter. " +
			"At the summary, use c to continue, r to review mistakes or q to finish.",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			runOpts := commandOptions(cmd, *opts)
			if err := checkProgressOverride(runOpts); err != nil {
				return err
			}
			chapters, err := loadChapters(runOpts)
			if err != nil {
				return err
			}
			difficulty := chapter.Difficulty(level)
			if cmd.Flags().Changed("difficulty") && !validDefaultDifficulty(difficulty) {
				return fmt.Errorf("invalid difficulty %q; choose beginner, intermediate or advanced", level)
			}
			sessionMode := chapter.SessionMode(mode)
			if cmd.Flags().Changed("mode") && sessionMode != chapter.ModeLearn && sessionMode != chapter.ModeReview {
				return fmt.Errorf("invalid mode %q; choose learn or review", mode)
			}
			filter := cliSessionFilter(cmd, chapters)
			if difficulty != "" {
				filter.Difficulty = difficulty
			}

			var selected chapter.Chapter
			scope := chapter.AllExercises
			if len(args) == 1 {
				var ok bool
				selected, ok = findChapter(chapters, args[0])
				if !ok {
					return fmt.Errorf("chapter not found: %s", args[0])
				}
			} else {
				index := filter.FirstTextChapter(chapters)
				if index < 0 {
					return fmt.Errorf("No eligible text exercises at %s difficulty with the enabled CLI modules. Try --difficulty or check cli-drill prefs.", filter.Difficulty)
				}
				selected = chapters[index]
				scope = chapter.TextExercises
			}
			if len(filter.EligibleItems(selected, scope)) == 0 {
				cmd.SilenceUsage, cmd.SilenceErrors = true, true
				if groups := filter.MissingModuleGroups(selected, scope); len(groups) > 0 {
					example, _, _ := strings.Cut(groups[0], " + ")
					return fmt.Errorf("%s: %s Missing modules (activate one group): %s. Enable each module separately, for example: cli-drill prefs --enable %s.",
						selected.ID, filter.EmptyReason(selected, scope), strings.Join(groups, " or "), example)
				}
				return fmt.Errorf("%s: %s Try --difficulty or check cli-drill prefs.", selected.ID, filter.EmptyReason(selected, scope))
			}
			state, err := progress.Load(progressPath(runOpts))
			if err != nil {
				return err
			}
			prepared := filter.PrepareVisit(selected, scope, sessionMode, all, state.CompletedExercises)
			if len(prepared.Items) == 0 {
				_, _ = fmt.Fprintf(runOpts.Out, "%s : tous les exercices de ce filtre sont acquis. Choisis Réviser avec --mode review, ou utilise --all.\n", selected.Title)
				return nil
			}
			runOpts.AllExercises = all
			return runTrainingFrom(runOpts, prepared, state)
		},
	}
	cmd.Flags().StringVar(&level, "difficulty", "", "session difficulty (beginner, intermediate, advanced), defaults to preferences")
	cmd.Flags().StringVar(&mode, "mode", "", "learn (default) or review")
	cmd.Flags().BoolVar(&all, "all", false, "train on all eligible exercises in one batch")
	return cmd
}

func cliSessionFilter(cmd *cobra.Command, chapters []chapter.Chapter) chapter.SessionFilter {
	prefs, err := loadPreferences(xdg.PreferencesFile())
	if err != nil {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Preferences unavailable. Using defaults. %v\n", err)
	}
	return chapter.SessionFilter{
		AvailableTools: chapter.AvailableTools(chapters),
		EnabledTools:   prefs.EnabledTools,
		Difficulty:     prefs.DefaultDifficulty,
	}
}

func findChapter(chapters []chapter.Chapter, id string) (chapter.Chapter, bool) {
	for _, candidate := range chapters {
		if candidate.ID == id {
			return candidate, true
		}
	}
	return chapter.Chapter{}, false
}

func newStatsCommand(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "stats",
		Short: "Show local progress",
		RunE: func(cmd *cobra.Command, args []string) error {
			runOpts := commandOptions(cmd, *opts)
			if err := checkProgressOverride(runOpts); err != nil {
				return err
			}
			state, err := progress.Load(progressPath(runOpts))
			if err != nil {
				return err
			}
			chapters, err := loadChapters(runOpts)
			if err != nil {
				return err
			}
			filter := cliSessionFilter(cmd, chapters)
			acquired, total := filter.ActiveCounts(chapters, state.CompletedExercises)
			_, _ = fmt.Fprintf(runOpts.Out, "Active acquired (%s, enabled CLI modules): %d/%d\n", filter.Difficulty, acquired, total)
			_, _ = fmt.Fprintln(runOpts.Out, "Historical totals (all sessions)")
			_, _ = fmt.Fprintf(runOpts.Out, "Completed exercises: %d\n", len(state.CompletedExercises))
			_, _ = fmt.Fprintf(runOpts.Out, "Chapters with attempts: %d\n", state.ChaptersWithAttempts())
			_, _ = fmt.Fprintf(runOpts.Out, "Streak: %d\n", state.Streak)
			return nil
		},
	}
}

func newResetCommand(opts *Options) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset local progress",
		RunE: func(cmd *cobra.Command, args []string) error {
			runOpts := commandOptions(cmd, *opts)
			if !yes {
				return errors.New("reset requires --yes")
			}
			if err := checkProgressOverride(runOpts); err != nil {
				return err
			}
			if err := progress.Reset(progressPath(runOpts)); err != nil {
				return err
			}
			_, _ = fmt.Fprintln(runOpts.Out, "Progress reset.")
			return nil
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "confirm progress reset")
	return cmd
}

type keySequenceAction int

const (
	keySequenceNext keySequenceAction = iota
	keySequenceRetry
	keySequenceQuit
)

type keySequenceResult struct {
	Action      keySequenceAction
	Answered    bool
	Correct     bool
	AnsweredKey string
}

func runKeySequenceExercise(opts Options, reader *bufio.Reader, selected chapter.Chapter, item chapter.Item, index, total int, state *progress.Progress) (keySequenceResult, error) {
	showHelp := false
	for {
		writeKeySequenceScreen(opts, renderKeySequenceQuestion(keySequenceQuestionView{
			Title:    selected.Title,
			Index:    index + 1,
			Total:    total,
			Prompt:   item.Prompt,
			ShowHelp: showHelp,
			Style:    keySequenceStyleFor(opts),
		}))
		input, aborted, err := readTrainingInput(opts, reader, item.ExerciseType)
		if err != nil {
			return keySequenceResult{Action: keySequenceQuit}, err
		}
		if aborted {
			_, _ = fmt.Fprintln(opts.Out, "Session interrompue.")
			return keySequenceResult{Action: keySequenceQuit}, nil
		}
		if input == "h" {
			showHelp = !showHelp
			continue
		}
		if input == "" || input == "Enter" {
			continue
		}

		correct := item.MatchAnswer(input)
		recorded, err := progress.RecordAnswer(progressPath(opts), selected.ID, item.ID, correct)
		if err != nil {
			return keySequenceResult{Action: keySequenceQuit}, err
		}
		*state = recorded
		explanation := ""
		if correct {
			explanation = item.Explanation
		}
		feedback := keySequenceFeedbackView{
			Received:    input,
			Correct:     correct,
			Explanation: explanation,
			Style:       keySequenceStyleFor(opts),
		}
		writeKeySequenceScreen(opts, renderKeySequenceFeedback(feedback))

		action, err := readKeySequenceFeedbackAction(opts, reader, item, feedback)
		return keySequenceResult{Action: action, Answered: true, Correct: correct, AnsweredKey: capturedKeyLabel(input)}, err
	}
}

type keySequenceQuestionView struct {
	Title    string
	Index    int
	Total    int
	Prompt   string
	ShowHelp bool
	Style    keySequenceStyle
}

type keySequenceFeedbackView struct {
	Received    string
	Correct     bool
	Explanation string
	Style       keySequenceStyle
}

type keySequenceSolutionView struct {
	Expected    string
	Explanation string
	Style       keySequenceStyle
}

func renderKeySequenceQuestion(view keySequenceQuestionView) string {
	style := view.Style
	footer := style.commandBar("h help", "Esc quit")
	lines := []string{
		style.title(view.Title),
		style.progress(fmt.Sprintf("Progression: %d/%d", view.Index, view.Total)),
		"",
		style.prompt(view.Prompt),
		"",
		style.muted("Waiting for key..."),
	}
	if view.ShowHelp {
		lines = append(lines, "", style.help("Help: press the requested shortcut. Normal keys count as answers."))
		footer = style.commandBar("h hide help", "Esc quit")
	}
	lines = append(lines, "", footer)
	return strings.Join(lines, "\n") + "\n"
}

func renderKeySequenceFeedback(view keySequenceFeedbackView) string {
	style := view.Style
	var builder strings.Builder
	_, _ = fmt.Fprintf(&builder, "\n%s %s\n", style.muted("Recu:"), capturedKeyLabel(view.Received))
	if view.Correct {
		_, _ = fmt.Fprintln(&builder, style.success("Correct."))
		if view.Explanation != "" {
			_, _ = fmt.Fprintln(&builder, view.Explanation)
		}
		_, _ = fmt.Fprintf(&builder, "\n%s\n", style.commandBar("Enter next", "Esc quit"))
		return builder.String()
	}
	_, _ = fmt.Fprintln(&builder, style.error("Pas encore."))
	_, _ = fmt.Fprintf(&builder, "\n%s\n", style.commandBar("Enter next", "r retry", "s solution", "Esc quit"))
	return builder.String()
}

type keySequenceStyle struct {
	Color bool
}

func keySequenceStyleFor(opts Options) keySequenceStyle {
	return keySequenceStyle{Color: terminalStylesAllowed(opts.Out)}
}

func terminalStylesAllowed(out io.Writer) bool {
	return writerIsTerminal(out) && os.Getenv("NO_COLOR") == "" && os.Getenv("TERM") != "dumb"
}

func (s keySequenceStyle) commandBar(commands ...string) string {
	return s.separator() + "\n" + s.footer(strings.Join(commands, " · "))
}

func (s keySequenceStyle) separator() string {
	return s.muted("────────────────────────────────────────")
}

func (s keySequenceStyle) title(value string) string    { return s.paint("1;36", value) }
func (s keySequenceStyle) progress(value string) string { return s.paint("2", value) }
func (s keySequenceStyle) prompt(value string) string   { return s.paint("1", value) }
func (s keySequenceStyle) success(value string) string  { return s.paint("1;32", value) }
func (s keySequenceStyle) error(value string) string    { return s.paint("1;31", value) }
func (s keySequenceStyle) solution(value string) string { return s.paint("1;33", value) }
func (s keySequenceStyle) help(value string) string     { return s.paint("36", value) }
func (s keySequenceStyle) footer(value string) string   { return s.paint("2", value) }
func (s keySequenceStyle) muted(value string) string    { return s.paint("2", value) }

func (s keySequenceStyle) paint(code, value string) string {
	if !s.Color || value == "" {
		return value
	}
	return "\x1b[" + code + "m" + value + "\x1b[0m"
}

func writeKeySequenceScreen(opts Options, content string) {
	if writerIsTerminal(opts.Out) && os.Getenv("TERM") != "dumb" {
		_, _ = fmt.Fprint(opts.Out, "\x1b[2J\x1b[H")
	}
	_, _ = fmt.Fprint(opts.Out, content)
}

func readKeySequenceFeedbackAction(opts Options, reader *bufio.Reader, item chapter.Item, feedback keySequenceFeedbackView) (keySequenceAction, error) {
	for {
		input, err := readKeySequence(opts, reader)
		if err != nil {
			return keySequenceNext, err
		}
		switch input {
		case "":
			continue
		case "Enter":
			return keySequenceNext, nil
		case "r", "R":
			if feedback.Correct {
				writeKeySequenceScreen(opts, renderKeySequenceFeedback(feedback))
				continue
			}
			return keySequenceRetry, nil
		case "s", "S":
			if feedback.Correct {
				writeKeySequenceScreen(opts, renderKeySequenceFeedback(feedback))
				continue
			}
			writeKeySequenceScreen(opts, renderKeySequenceSolution(keySequenceSolutionView{
				Expected:    item.Answer.Primary,
				Explanation: item.Explanation,
				Style:       feedback.Style,
			}))
		case "Esc", "Ctrl+C":
			_, _ = fmt.Fprintln(opts.Out, "\nSession interrompue.")
			return keySequenceQuit, nil
		default:
			writeKeySequenceScreen(opts, renderKeySequenceFeedback(feedback))
		}
	}
}

func renderKeySequenceSolution(view keySequenceSolutionView) string {
	style := view.Style
	var builder strings.Builder
	_, _ = fmt.Fprintf(&builder, "\n%s %s\n", style.solution("Solution:"), view.Expected)
	if view.Explanation != "" {
		_, _ = fmt.Fprintln(&builder, view.Explanation)
	}
	_, _ = fmt.Fprintf(&builder, "\n%s\n", style.commandBar("Enter next", "r retry", "Esc quit"))
	return builder.String()
}

func readTrainingInput(opts Options, reader *bufio.Reader, exerciseType exercise.Type) (string, bool, error) {
	if exerciseType != exercise.TypeKeySequence {
		input, err := reader.ReadString('\n')
		if strings.ContainsRune(input, '\x03') {
			return input, true, err
		}
		line := strings.TrimRight(input, "\r\n")
		if escape := strings.IndexByte(line, '\x1b'); escape >= 0 {
			if escape == len(line)-1 {
				// A standalone Escape can follow an unfinished answer.
				return input, true, err
			}
			// Do not evaluate a line whose terminal controls may have
			// changed the intended answer or menu action.
			return "", false, err
		}
		return input, false, err
	}

	input, err := readKeySequence(opts, reader)
	if err != nil {
		return "", false, err
	}
	if input == "Esc" || input == "Ctrl+C" {
		return input, true, nil
	}
	return input, false, nil
}

func readKeySequence(opts Options, reader *bufio.Reader) (string, error) {
	if file, ok := opts.In.(*os.File); ok && isTerminalFD(file.Fd()) {
		return readRawKeyFromTerminal(file)
	}

	return readCaptureKey(func(time.Duration) (byte, error) {
		return reader.ReadByte()
	})
}

type fdProvider interface {
	Fd() uintptr
}

func writerIsTerminal(w io.Writer) bool {
	fd, ok := writerFD(w)
	return ok && isTerminalFD(fd)
}

func tuiTerminalError(opts Options) error {
	input, ok := opts.In.(fdProvider)
	if !ok || !isTerminalFD(input.Fd()) || !writerIsTerminal(opts.Out) {
		return errors.New("the interactive interface requires a terminal for input and output. Run cli-drill train for text training")
	}
	if os.Getenv("TERM") == "dumb" {
		return errors.New("the interactive interface is unavailable with TERM=dumb. Run cli-drill train for text training")
	}
	return nil
}

func terminalCanCaptureKeys(opts Options) bool {
	input, ok := opts.In.(*os.File)
	return rawCaptureSupported && ok && isTerminalFD(input.Fd()) && writerIsTerminal(opts.Out)
}

func writerFD(w io.Writer) (uintptr, bool) {
	provider, ok := w.(fdProvider)
	if !ok {
		return 0, false
	}
	return provider.Fd(), true
}

func loadChapters(opts Options) ([]chapter.Chapter, error) {
	return chapter.LoadFS(opts.DefaultFS, "*.yaml")
}

func commandOptions(cmd *cobra.Command, opts Options) Options {
	opts.In = cmd.InOrStdin()
	opts.Out = cmd.OutOrStdout()
	return opts
}

func progressPath(opts Options) string {
	if opts.ProgressPath != "" {
		return opts.ProgressPath
	}
	return xdg.ProgressFile()
}

// The schemas both use version 1. Do not let an explicit progress path make
// preferences look like an empty progress file that training or reset can replace.
func checkProgressOverride(opts Options) error {
	if opts.ProgressPath == "" {
		return nil
	}
	override := opts.ProgressPath
	prefs := xdg.PreferencesFile()
	a, aErr := filepath.Abs(override)
	b, bErr := filepath.Abs(prefs)
	conflict := aErr == nil && bErr == nil && a == b
	if !conflict {
		overrideInfo, overrideErr := os.Stat(override)
		prefsInfo, prefsErr := os.Stat(prefs)
		conflict = overrideErr == nil && prefsErr == nil && os.SameFile(overrideInfo, prefsInfo)
	}
	if !conflict && filepath.Base(override) == filepath.Base(prefs) {
		// Also protect the reserved filename through a parent alias before
		// preferences.json itself exists.
		overrideDir, overrideErr := os.Stat(filepath.Dir(override))
		prefsDir, prefsErr := os.Stat(filepath.Dir(prefs))
		conflict = overrideErr == nil && prefsErr == nil && os.SameFile(overrideDir, prefsDir)
	}
	if conflict {
		return errors.New("--progress cannot point to the preferences file")
	}
	return nil
}

func tuiOptions(chapters []chapter.Chapter, progressPath, preferencesPath string) tui.Options {
	var loaded preferences
	var loadErr error
	return tui.Options{
		Chapters:     chapters,
		ProgressPath: progressPath,
		LoadPreferences: func() ([]string, chapter.Difficulty, error) {
			loaded, loadErr = loadPreferences(preferencesPath)
			return loaded.EnabledTools, loaded.DefaultDifficulty, loadErr
		},
		SavePreferences: func(tools []string, level chapter.Difficulty) error {
			if loadErr != nil {
				return loadErr
			}
			updated := loaded
			updated.EnabledTools, updated.DefaultDifficulty = tools, level
			if err := savePreferences(preferencesPath, loaded, updated); err != nil {
				return err
			}
			loaded = updated
			return nil
		},
	}
}
