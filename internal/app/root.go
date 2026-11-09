package app

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

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
}

var (
	isTerminalFD           = term.IsTerminal
	readRawKeyFromTerminal = readRawKey
	runTUIWithOptions      = tui.RunWithOptions
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
			chapters, err := loadChapters(runOpts)
			if err != nil {
				return err
			}
			for {
				result, err := runTUIWithOptions(tui.Options{
					Chapters:     chapters,
					ProgressPath: progressPath(runOpts),
				})
				if err != nil {
					return err
				}
				if result.RunCLITrainChapterID == "" {
					return nil
				}
				selected, ok := findChapter(chapters, result.RunCLITrainChapterID)
				if !ok {
					return fmt.Errorf("chapter not found: %s", result.RunCLITrainChapterID)
				}
				if err := runTraining(runOpts, selected); err != nil {
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
			for _, chapter := range chapters {
				_, _ = fmt.Fprintf(runOpts.Out, "%s\t%s\t%d exercises\n", chapter.ID, chapter.Title, len(chapter.Items))
			}
			return nil
		},
	}
}

func newTrainCommand(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "train [chapter-id]",
		Short: "Train from a chapter in the terminal",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			runOpts := commandOptions(cmd, *opts)
			chapters, err := loadChapters(runOpts)
			if err != nil {
				return err
			}
			if len(chapters) == 0 {
				return errors.New("no chapters available")
			}
			selected := chapters[0]
			if len(args) == 1 {
				var ok bool
				selected, ok = findChapter(chapters, args[0])
				if !ok {
					return fmt.Errorf("chapter not found: %s", args[0])
				}
			}
			return runTraining(runOpts, selected)
		},
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
			state, err := progress.Load(progressPath(runOpts))
			if err != nil {
				return err
			}
			_, _ = fmt.Fprintf(runOpts.Out, "Completed exercises: %d\n", len(state.CompletedExercises))
			_, _ = fmt.Fprintf(runOpts.Out, "Chapters with attempts: %d\n", len(state.ChapterScores))
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

func runTraining(opts Options, selected chapter.Chapter) error {
	if len(selected.Items) == 0 {
		return fmt.Errorf("chapter %s has no exercises", selected.ID)
	}

	if _, err := progress.Load(progressPath(opts)); err != nil {
		return err
	}
	reader := bufio.NewReader(opts.In)

	textHeaderPrinted := false
	keyStats := keySequenceStats{}
	keyQuit := false
	for i := 0; i < len(selected.Items); {
		item := selected.Items[i]
		if item.ExerciseType == exercise.TypeKeySequence {
			result, err := runKeySequenceExercise(opts, reader, selected, item, i, len(selected.Items))
			if err != nil && !errors.Is(err, io.EOF) {
				return err
			}
			if result.Action == keySequenceQuit {
				keyQuit = true
				break
			}
			if result.Action == keySequenceNext {
				keyStats.Total++
				if result.Correct {
					keyStats.Correct++
				} else {
					keyStats.Missed = append(keyStats.Missed, keySequenceReviewItem{Item: item})
				}
				i++
			}
			if errors.Is(err, io.EOF) {
				break
			}
			continue
		}

		if !textHeaderPrinted {
			_, _ = fmt.Fprintf(opts.Out, "Chapitre: %s\n", selected.Title)
			textHeaderPrinted = true
		}
		_, _ = fmt.Fprintf(opts.Out, "\n%d/%d %s\n", i+1, len(selected.Items), item.Prompt)
		if item.ExerciseType == exercise.TypeMultipleChoice {
			for _, choice := range item.Choices {
				_, _ = fmt.Fprintf(opts.Out, "  - %s\n", choice)
			}
		}
		_, _ = fmt.Fprint(opts.Out, "> ")
		input, aborted, err := readTrainingInput(opts, reader, item.ExerciseType)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		if aborted {
			_, _ = fmt.Fprintln(opts.Out, "Session interrompue.")
			break
		}
		if strings.Trim(input, " \r\n") == "" {
			if errors.Is(err, io.EOF) {
				break
			}
			continue
		}
		correct := item.MatchAnswer(input)
		if _, err := progress.RecordAnswer(progressPath(opts), selected.ID, item.ID, correct); err != nil {
			return err
		}
		if correct {
			_, _ = fmt.Fprintln(opts.Out, "Correct.")
		} else {
			_, _ = fmt.Fprintf(opts.Out, "Pas encore. Reponse attendue: %s\n", item.Answer.Primary)
		}
		if item.Explanation != "" {
			_, _ = fmt.Fprintln(opts.Out, item.Explanation)
		}
		if errors.Is(err, io.EOF) {
			break
		}
		i++
	}
	if keyStats.Total > 0 && !keyQuit {
		if err := runKeySequenceReview(opts, reader, selected, keyStats); err != nil && !errors.Is(err, io.EOF) {
			return err
		}
	}
	return nil
}

type keySequenceStats struct {
	Total   int
	Correct int
	Missed  []keySequenceReviewItem
}

type keySequenceReviewItem struct {
	Item chapter.Item
}

type keySequenceAction int

const (
	keySequenceNext keySequenceAction = iota
	keySequenceRetry
	keySequenceQuit
)

type keySequenceResult struct {
	Action  keySequenceAction
	Correct bool
}

func runKeySequenceExercise(opts Options, reader *bufio.Reader, selected chapter.Chapter, item chapter.Item, index, total int) (keySequenceResult, error) {
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
		if input == "Enter" {
			continue
		}

		correct := item.MatchAnswer(input)
		if _, err := progress.RecordAnswer(progressPath(opts), selected.ID, item.ID, correct); err != nil {
			return keySequenceResult{Action: keySequenceQuit}, err
		}
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
		if err != nil {
			return keySequenceResult{Action: keySequenceNext, Correct: correct}, err
		}
		return keySequenceResult{Action: action, Correct: correct}, nil
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

type keySequenceSummaryView struct {
	Correct int
	Total   int
	Missed  int
	Style   keySequenceStyle
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
	_, _ = fmt.Fprintf(&builder, "\n%s %s\n", style.muted("Recu:"), view.Received)
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
	if !writerIsTerminal(opts.Out) {
		return keySequenceStyle{}
	}
	if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		return keySequenceStyle{}
	}
	return keySequenceStyle{Color: true}
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
	if writerIsTerminal(opts.Out) {
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

func runKeySequenceReview(opts Options, reader *bufio.Reader, selected chapter.Chapter, stats keySequenceStats) error {
	missed := stats.Missed
	for {
		renderKeySequenceSummary(opts.Out, stats.Correct, stats.Total, len(missed), keySequenceStyleFor(opts))
		if len(missed) == 0 {
			return nil
		}
		action, err := readKeySequenceReviewAction(opts, reader)
		if err != nil {
			return err
		}
		if action == keySequenceQuit {
			return nil
		}

		reviewStats := keySequenceStats{}
		for i := 0; i < len(missed); {
			review := missed[i]
			result, err := runKeySequenceExercise(opts, reader, selected, review.Item, i, len(missed))
			if err != nil {
				return err
			}
			if result.Action == keySequenceQuit {
				return nil
			}
			if result.Action == keySequenceRetry {
				continue
			}
			reviewStats.Total++
			if result.Correct {
				reviewStats.Correct++
			} else {
				reviewStats.Missed = append(reviewStats.Missed, review)
			}
			i++
		}
		stats = reviewStats
		missed = reviewStats.Missed
	}
}

func renderKeySequenceSummary(w io.Writer, correct, total, missed int, style keySequenceStyle) {
	_, _ = fmt.Fprint(w, renderKeySequenceSummaryView(keySequenceSummaryView{
		Correct: correct,
		Total:   total,
		Missed:  missed,
		Style:   style,
	}))
}

func renderKeySequenceSummaryView(view keySequenceSummaryView) string {
	style := view.Style
	var builder strings.Builder
	_, _ = fmt.Fprintf(&builder, "\n%s\nCorrect: %d/%d\nÀ revoir: %d\n", style.title("Chapitre termine."), view.Correct, view.Total, view.Missed)
	if view.Missed > 0 {
		_, _ = fmt.Fprintf(&builder, "\n%s\n", style.commandBar("Enter review missed", "Esc quit"))
	}
	return builder.String()
}

func readKeySequenceReviewAction(opts Options, reader *bufio.Reader) (keySequenceAction, error) {
	for {
		input, err := readKeySequence(opts, reader)
		if err != nil {
			return keySequenceQuit, err
		}
		switch input {
		case "Enter":
			return keySequenceNext, nil
		case "Esc", "Ctrl+C":
			_, _ = fmt.Fprintln(opts.Out, "\nSession interrompue.")
			return keySequenceQuit, nil
		default:
			_, _ = fmt.Fprintf(opts.Out, "\n%s\n", keySequenceStyleFor(opts).commandBar("Enter review missed", "Esc quit"))
		}
	}
}

func readTrainingInput(opts Options, reader *bufio.Reader, exerciseType exercise.Type) (string, bool, error) {
	if exerciseType != exercise.TypeKeySequence {
		input, err := reader.ReadString('\n')
		return input, strings.ContainsAny(input, "\x1b\x03"), err
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

	b, err := reader.ReadByte()
	if err != nil {
		return "", err
	}
	if b == '\r' || b == '\n' {
		return "Enter", nil
	}
	if notation, ok := exercise.KeyByteNotation(b); ok {
		return notation, nil
	}
	return string([]byte{b}), nil
}

func readRawKey(file *os.File) (string, error) {
	fd := file.Fd()
	state, err := term.MakeRaw(fd)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = term.Restore(fd, state)
	}()

	var b [1]byte
	if _, err := file.Read(b[:]); err != nil {
		return "", err
	}
	if b[0] == '\r' || b[0] == '\n' {
		return "Enter", nil
	}
	if notation, ok := exercise.KeyByteNotation(b[0]); ok {
		return notation, nil
	}
	return string(b[:]), nil
}

type fdProvider interface {
	Fd() uintptr
}

func writerIsTerminal(w io.Writer) bool {
	fd, ok := writerFD(w)
	return ok && isTerminalFD(fd)
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
