package zsh

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	"github.com/benabot/cli-drill/internal/catalog"
)

type Alias struct {
	Name   string
	Value  string
	Source string
	Line   int
}

type Function struct {
	Name   string
	Source string
	Line   int
}

var (
	aliasPattern      = regexp.MustCompile(`^\s*alias\s+([A-Za-z0-9_+\-.,:]+)=(.*)$`)
	functionPattern   = regexp.MustCompile(`^\s*([A-Za-z_][A-Za-z0-9_-]*)\s*\(\s*\)\s*\{?.*$`)
	functionKeywordRe = regexp.MustCompile(`^\s*function\s+([A-Za-z_][A-Za-z0-9_-]*)(?:\s*\(\s*\))?\s*\{?.*$`)
)

func ParseAliases(r io.Reader, source string) ([]Alias, error) {
	scanner := bufio.NewScanner(r)
	var aliases []Alias
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		matches := aliasPattern.FindStringSubmatch(line)
		if len(matches) != 3 {
			continue
		}

		aliases = append(aliases, Alias{
			Name:   matches[1],
			Value:  unquoteAliasValue(matches[2]),
			Source: source,
			Line:   lineNumber,
		})
	}

	return aliases, scanner.Err()
}

func ParseFunctions(r io.Reader, source string) ([]Function, error) {
	scanner := bufio.NewScanner(r)
	var functions []Function
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		if matches := functionKeywordRe.FindStringSubmatch(line); len(matches) == 2 {
			functions = append(functions, Function{Name: matches[1], Source: source, Line: lineNumber})
			continue
		}
		if matches := functionPattern.FindStringSubmatch(line); len(matches) == 2 {
			functions = append(functions, Function{Name: matches[1], Source: source, Line: lineNumber})
		}
	}

	return functions, scanner.Err()
}

func AliasEntries(aliases []Alias) []catalog.Entry {
	entries := make([]catalog.Entry, 0, len(aliases))
	for _, alias := range aliases {
		entries = append(entries, catalog.Entry{
			ID:      catalog.NormalizeID(alias.Name),
			Name:    alias.Name,
			Type:    catalog.EntryAlias,
			Summary: alias.Value,
			Command: alias.Value,
			Source:  catalog.Source{Path: alias.Source, Line: alias.Line},
		})
	}
	return entries
}

func FunctionEntries(functions []Function) []catalog.Entry {
	entries := make([]catalog.Entry, 0, len(functions))
	for _, function := range functions {
		entries = append(entries, catalog.Entry{
			ID:      catalog.NormalizeID(function.Name),
			Name:    function.Name,
			Type:    catalog.EntryFunction,
			Summary: "ZSH function",
			Source:  catalog.Source{Path: function.Source, Line: function.Line},
		})
	}
	return entries
}

func unquoteAliasValue(value string) string {
	value = strings.TrimSpace(stripInlineComment(value))
	if len(value) < 2 {
		return value
	}

	first := value[0]
	last := value[len(value)-1]
	if (first == '\'' && last == '\'') || (first == '"' && last == '"') {
		return value[1 : len(value)-1]
	}
	return value
}

func stripInlineComment(value string) string {
	inSingleQuote := false
	inDoubleQuote := false
	for i, r := range value {
		switch r {
		case '\'':
			if !inDoubleQuote {
				inSingleQuote = !inSingleQuote
			}
		case '"':
			if !inSingleQuote {
				inDoubleQuote = !inDoubleQuote
			}
		case '#':
			if !inSingleQuote && !inDoubleQuote && i > 0 && value[i-1] == ' ' {
				return value[:i]
			}
		}
	}
	return value
}
