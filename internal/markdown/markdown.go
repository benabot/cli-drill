package markdown

import (
	"bufio"
	"io"
	"strings"
)

type Document struct {
	Headings   []Heading
	ListItems  []string
	CodeBlocks []CodeBlock
	Tables     []Table
}

type Heading struct {
	Level int
	Text  string
	Line  int
}

type CodeBlock struct {
	Language string
	Content  string
	Line     int
}

type Table struct {
	Rows []string
	Line int
}

func Parse(r io.Reader) (Document, error) {
	scanner := bufio.NewScanner(r)
	doc := Document{}
	lineNumber := 0
	inCode := false
	codeStart := 0
	codeLanguage := ""
	var codeLines []string
	var tableRows []string
	tableStart := 0

	flushTable := func() {
		if len(tableRows) > 0 {
			doc.Tables = append(doc.Tables, Table{Rows: append([]string(nil), tableRows...), Line: tableStart})
			tableRows = nil
			tableStart = 0
		}
	}

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "```") {
			if inCode {
				doc.CodeBlocks = append(doc.CodeBlocks, CodeBlock{
					Language: codeLanguage,
					Content:  strings.Join(codeLines, "\n"),
					Line:     codeStart,
				})
				inCode = false
				codeLines = nil
				codeLanguage = ""
				continue
			}
			flushTable()
			inCode = true
			codeStart = lineNumber
			codeLanguage = strings.TrimSpace(strings.TrimPrefix(trimmed, "```"))
			continue
		}
		if inCode {
			codeLines = append(codeLines, line)
			continue
		}

		if headingLevel := parseHeadingLevel(trimmed); headingLevel > 0 {
			flushTable()
			doc.Headings = append(doc.Headings, Heading{
				Level: headingLevel,
				Text:  strings.TrimSpace(trimmed[headingLevel:]),
				Line:  lineNumber,
			})
			continue
		}
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			flushTable()
			doc.ListItems = append(doc.ListItems, strings.TrimSpace(trimmed[2:]))
			continue
		}
		if strings.Contains(trimmed, "|") {
			if len(tableRows) == 0 {
				tableStart = lineNumber
			}
			tableRows = append(tableRows, trimmed)
			continue
		}
		flushTable()
	}
	flushTable()
	if inCode {
		doc.CodeBlocks = append(doc.CodeBlocks, CodeBlock{Language: codeLanguage, Content: strings.Join(codeLines, "\n"), Line: codeStart})
	}
	return doc, scanner.Err()
}

func parseHeadingLevel(line string) int {
	level := 0
	for _, r := range line {
		if r != '#' {
			break
		}
		level++
	}
	if level == 0 || level > 6 || len(line) <= level || line[level] != ' ' {
		return 0
	}
	return level
}
