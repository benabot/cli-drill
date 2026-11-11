package app

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/benabot/cli-drill/internal/chapter"
	"github.com/charmbracelet/x/ansi"
)

func writeChapters(out io.Writer, chapters []chapter.Chapter) error {
	if !writerIsTerminal(out) {
		return writeChapterRows(out, chapters, "exercise")
	}
	fd, _ := writerFD(out)
	width, _, err := terminalSize(fd)
	if err != nil || width < 1 {
		width = 80
	}
	if width <= 80 {
		return writeCompactChapters(out, chapters, width)
	}
	table := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(table, "Identifiant\tChapitre\tTotal"); err != nil {
		return err
	}
	if err := writeChapterRows(table, chapters, "exercice"); err != nil {
		return err
	}
	return table.Flush()
}

func writeChapterRows(out io.Writer, chapters []chapter.Chapter, unit string) error {
	for _, ch := range chapters {
		label := unit
		if len(ch.Items) != 1 {
			label += "s"
		}
		if _, err := fmt.Fprintf(out, "%s\t%s\t%d %s\n", ch.ID, ch.Title, len(ch.Items), label); err != nil {
			return err
		}
	}
	return nil
}

func writeCompactChapters(out io.Writer, chapters []chapter.Chapter, width int) error {
	rows := [][3]string{{"Identifiant", "Chapitre", "Total"}}
	var columns [3]int
	for _, ch := range chapters {
		rows = append(rows, [3]string{ch.ID, ch.Title, strconv.Itoa(len(ch.Items))})
	}
	for _, row := range rows {
		for i, value := range row {
			columns[i] = max(columns[i], ansi.StringWidth(value))
		}
	}
	if columns[0]+columns[1]+columns[2]+4 <= width {
		for _, row := range rows {
			if _, err := fmt.Fprintln(out, padChapterColumn(row[0], columns[0]+2)+padChapterColumn(row[1], columns[1]+2)+row[2]); err != nil {
				return err
			}
		}
		return nil
	}
	if _, err := fmt.Fprintln(out, ansi.Wrap("Identifiant (total) et chapitre", width, "")); err != nil {
		return err
	}
	for _, row := range rows[1:] {
		header := row[0] + " (" + row[2] + ")"
		if _, err := fmt.Fprintln(out, ansi.Wrap(header, width, "")+"\n"+ansi.Wrap(row[1], width, "")); err != nil {
			return err
		}
	}
	return nil
}

func padChapterColumn(value string, width int) string {
	return value + strings.Repeat(" ", width-ansi.StringWidth(value))
}
