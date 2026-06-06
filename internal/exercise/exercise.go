package exercise

import "strings"

type Type string

const (
	TypeFreeAnswer     Type = "free-answer"
	TypeMultipleChoice Type = "multiple-choice"
	TypeScenario       Type = "scenario"
	TypeSimpleShellSim Type = "simple-shell-sim"
	TypeKeySequence    Type = "key-sequence"
)

type Answer struct {
	Primary  string   `json:"primary" yaml:"primary"`
	Accepted []string `json:"accepted,omitempty" yaml:"accepted,omitempty"`
}

func MatchAnswer(input string, answer Answer) bool {
	normalizedInput := NormalizeAnswer(input)
	candidates := append([]string{answer.Primary}, answer.Accepted...)
	for _, candidate := range candidates {
		if normalizedInput == NormalizeAnswer(candidate) {
			return true
		}
	}
	return false
}

func NormalizeAnswer(value string) string {
	value = trimAnswerEdges(value)
	if shortcut, ok := normalizeRawControlShortcut(value); ok {
		return shortcut
	}

	value = strings.ToLower(value)
	value = strings.Join(strings.Fields(value), " ")
	value = normalizePipeSpacing(value)
	value = normalizeControlShortcut(value)
	return value
}

func trimAnswerEdges(value string) string {
	return strings.Trim(value, " \r\n")
}

func normalizePipeSpacing(value string) string {
	parts := strings.Split(value, "|")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return strings.Join(parts, "|")
}

func normalizeRawControlShortcut(value string) (string, bool) {
	if len(value) != 1 {
		return "", false
	}
	notation, ok := KeyByteNotation(value[0])
	if !ok {
		return "", false
	}
	return strings.ToLower(notation), true
}

func normalizeControlShortcut(value string) string {
	if len(value) == 2 && value[0] == '^' && isASCIILetter(value[1]) {
		return "ctrl+" + string(value[1])
	}

	parts := strings.Fields(strings.NewReplacer("+", " ", "-", " ").Replace(value))
	if len(parts) == 2 && (parts[0] == "ctrl" || parts[0] == "control") && len(parts[1]) == 1 && isASCIILetter(parts[1][0]) {
		return "ctrl+" + parts[1]
	}

	return value
}

func isASCIILetter(value byte) bool {
	return value >= 'a' && value <= 'z'
}

func KeyByteNotation(value byte) (string, bool) {
	if value == 0x1b {
		return "Esc", true
	}
	if value >= 1 && value <= 26 {
		return "Ctrl+" + string(rune('A'+value-1)), true
	}
	return "", false
}
