package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

type readablePage struct {
	lines  []string
	bottom string
	rows   int
	hint   bool
}

func (p readablePage) lastOffset() int {
	return max(0, len(p.lines)-p.rows)
}

func (m model) contentWidth() int {
	if m.width == 0 {
		return 76
	}
	return max(1, m.width-2*m.horizontalPadding())
}

func (m model) horizontalPadding() int {
	if m.width > 0 && m.width < 12 {
		return 0
	}
	return 2
}

func (m model) wrap(value string) string {
	return ansi.Wrap(value, m.contentWidth(), "")
}

func (m model) wrapActions(text string) string {
	width := m.contentWidth()
	var lines []string
	line := ""
	for _, action := range strings.Split(text, " · ") {
		if line == "" {
			line = action
		} else if ansi.StringWidth(line+" · "+action) <= width {
			line += " · " + action
		} else {
			lines = append(lines, line)
			line = action
		}
	}
	return m.wrap(strings.Join(append(lines, line), "\n"))
}

func (m model) page() readablePage {
	body := m.viewBody()
	styles := newTUIStyles()
	actions := styles.Footer.Render(m.wrapActions(m.footer()))
	bottom := actions
	var input string
	if m.screen == screenTraining && m.currentItem() != nil && !m.isMultipleChoice() && m.feedback == "" {
		input = m.wrap(m.input.View())
		bottom = styles.Muted.Render("Réponse") + "\n" + input + "\n" + bottom
	}
	if m.height > 0 && lipgloss.Height(bottom) >= m.height {
		// Keep a row for the content even when the normal footer fills the window.
		body += "\n" + actions
		bottom = input
		if m.height == 1 {
			if input != "" {
				body += "\n" + input
			}
			bottom = ""
		}
	}
	lines := strings.Split(m.wrap(body), "\n")
	rows := len(lines)
	hint := false
	if m.height > 0 {
		rows = m.height
		if bottom != "" {
			rows -= lipgloss.Height(bottom)
		}
		rows = max(1, rows)
		if len(lines) > rows && rows > 1 && m.contentWidth() >= 26 {
			hint = true
			rows--
		}
	}
	return readablePage{lines: lines, bottom: bottom, rows: rows, hint: hint}
}

func (m model) View() string {
	p := m.page()
	start := min(m.scrollOffset, p.lastOffset())
	end := min(len(p.lines), start+p.rows)
	body := strings.Join(p.lines[start:end], "\n")
	bottom := p.bottom
	if p.hint {
		counter := fmt.Sprintf("PgUp/PgDn · %d–%d/%d", start+1, end, len(p.lines))
		if bottom != "" {
			bottom = "\n" + bottom
		}
		bottom = newTUIStyles().Footer.Render(counter) + bottom
	}
	view := body
	if bottom != "" {
		view += "\n" + bottom
	}
	padding := m.horizontalPadding()
	style := lipgloss.NewStyle().Padding(0, padding)
	if m.height == 0 || (m.height >= 20 && lipgloss.Height(view)+3 <= m.height) {
		style = style.Padding(1, padding)
		view = body + "\n\n" + bottom
	}
	return style.Render(view)
}

func (m *model) scroll(direction int) bool {
	p := m.page()
	last := p.lastOffset()
	if last == 0 {
		return false
	}
	// Overlap a line when possible so the next page keeps its context.
	m.scrollOffset = max(0, min(last, m.scrollOffset+direction*max(1, p.rows-1)))
	return true
}

func (m *model) focusChoice() {
	p := m.page()
	if p.lastOffset() == 0 {
		return
	}
	// Locate the choice after the question, allowing for wrapped options.
	for row := len(p.lines) - 1; row >= 0; row-- {
		if !strings.HasPrefix(ansi.Strip(p.lines[row]), "> ") {
			continue
		}
		start := min(m.scrollOffset, p.lastOffset())
		if row < start || row >= start+p.rows {
			m.scrollOffset = min(row, p.lastOffset())
		}
		return
	}
}

func (m *model) resize(width, height int) {
	m.width, m.height, m.scrollOffset = width, height, 0
	width = m.contentWidth()
	homeFooterRows := lipgloss.Height(m.wrapActions(footerText(footerHome)))
	homeHeaderRows := lipgloss.Height(m.wrap(m.homeHeader()))
	m.home.SetSize(width, max(3, height-homeHeaderRows-homeFooterRows-2))
	// Reserve the title, level, filter field, and the tallest list footer.
	chapterFooterRows := max(lipgloss.Height(m.wrapActions(footerText(footerNavigation))),
		lipgloss.Height(m.wrapActions(footerText(footerFilteredList))),
		lipgloss.Height(m.wrapActions(footerText(footerFilterInput))))
	m.chapterList.SetSize(width, max(3, height-lipgloss.Height(m.wrap(m.chaptersHeader()))-chapterFooterRows-2))
	m.input.Width = max(1, width-lipgloss.Width(m.input.Prompt)-1)
	m.preferences.resize(width, height)
}

func (m model) answerLine(label, answer string) string {
	if ansi.StringWidth(label+answer) <= m.contentWidth() {
		return label + answer
	}
	return strings.TrimSpace(label) + "\n" + ansi.Hardwrap(answer, m.contentWidth(), true)
}
