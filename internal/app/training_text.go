package app

import (
	"fmt"
	"io"
	"strings"

	"github.com/benabot/cli-drill/internal/chapter"
	"github.com/benabot/cli-drill/internal/exercise"
	"github.com/charmbracelet/lipgloss"
)

type textTrainingView struct {
	out                     io.Writer
	title, label, secondary lipgloss.Style
	success, failure        lipgloss.Style
}

func newTextTrainingView(out io.Writer) textTrainingView {
	view := textTrainingView{out: out}
	if terminalStylesAllowed(out) {
		style := lipgloss.NewRenderer(out).NewStyle()
		view.title = style.Bold(true).Foreground(lipgloss.Color("6"))
		view.label = style.Bold(true)
		view.secondary = style.Faint(true)
		view.success = style.Bold(true).Foreground(lipgloss.Color("2"))
		view.failure = style.Bold(true).Foreground(lipgloss.Color("1"))
	}
	return view
}

func (v textTrainingView) chapterTitle(title string) {
	_, _ = fmt.Fprintln(v.out, v.title.Render("Chapitre : "+title))
}

func (v textTrainingView) separator() {
	width := 40
	if writerIsTerminal(v.out) {
		fd, _ := writerFD(v.out)
		if columns, _, err := terminalSize(fd); err == nil && columns > 0 {
			width = min(width, columns)
		}
	}
	_, _ = fmt.Fprintln(v.out, "\n"+v.secondary.Render(strings.Repeat("-", width)))
}

func (v textTrainingView) question(item chapter.Item, number, total int) {
	_, _ = fmt.Fprintf(v.out, "\n%s\n\n%s\n", v.label.Render(fmt.Sprintf("Exercice %d/%d", number, total)), item.Prompt)
	if item.ExerciseType == exercise.TypeMultipleChoice {
		for _, choice := range item.Choices {
			_, _ = fmt.Fprintf(v.out, "  - %s\n", choice)
		}
	}
	_, _ = fmt.Fprintf(v.out, "\n%s\n> ", v.secondary.Render("Réponse"))
}

func (v textTrainingView) feedback(item chapter.Item, correct bool) {
	if correct {
		_, _ = fmt.Fprintf(v.out, "\n%s\n", v.success.Render("Correct."))
	} else {
		_, _ = fmt.Fprintf(v.out, "\n%s\n%s %s\n", v.failure.Render("Réponse incorrecte."),
			v.label.Render("Réponse attendue :"), item.Answer.Primary)
	}
	if item.Explanation != "" {
		_, _ = fmt.Fprintf(v.out, "\n%s\n%s\n", v.secondary.Render("Explication"), item.Explanation)
	}
}

func (v textTrainingView) summary(title string, pass trainingPass, acquired int) {
	v.separator()
	_, _ = fmt.Fprintf(v.out, "%s\nBonnes réponses : %d/%d\nExercices réalisés : %d\nErreurs : %d\nAcquis cumulés (historique) : %d\n",
		v.title.Render(title), pass.Correct, pass.Total, pass.Total, len(pass.Missed), acquired)
	if pass.immediate.Attempts > 0 {
		_, _ = fmt.Fprintf(v.out, "Réessais immédiats : %d/%d bonnes réponses\n", pass.immediate.Correct, pass.immediate.Attempts)
	}
}
