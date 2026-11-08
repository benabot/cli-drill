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
	return matchAnswer(input, answer, trimAnswerEdges)
}

// MatchConceptAnswer tolerates case only for a bare pedagogical name.
// An explicitly accepted control shortcut retains its notation variants.
func MatchConceptAnswer(input string, answer Answer) bool {
	for _, candidate := range answer.Accepted {
		if normalizeControlShortcut(strings.ToLower(trimAnswerEdges(candidate))) != "" &&
			MatchKeySequence(input, Answer{Primary: candidate}) {
			return true
		}
	}
	return matchAnswer(input, answer, func(value string) string {
		value = trimAnswerEdges(value)
		for _, char := range value {
			if !(char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' ||
				char >= '0' && char <= '9' || char == '-' || char == '_') {
				return value
			}
		}
		return strings.ToLower(value)
	})
}

// MatchShellAnswer only relaxes spacing in plain words and single pipelines.
// Quoting, escaping, expansions and other operators require explicit variants.
func MatchShellAnswer(input string, answer Answer) bool {
	return matchAnswer(input, answer, normalizeShellSpacing)
}

func MatchKeySequence(input string, answer Answer) bool {
	return matchAnswer(input, answer, normalizeKeySequence)
}

func matchAnswer(input string, answer Answer, normalize func(string) string) bool {
	normalizedInput := normalize(input)
	if normalizedInput == "" {
		return false
	}
	candidates := append([]string{answer.Primary}, answer.Accepted...)
	for _, candidate := range candidates {
		if normalizedInput == normalize(candidate) {
			return true
		}
	}
	return false
}

func normalizeKeySequence(value string) string {
	if shortcut, ok := normalizeRawControlShortcut(value); ok {
		return shortcut
	}
	value = trimEdges(value, " \r\n")
	if shortcut, ok := normalizeRawControlShortcut(value); ok {
		return shortcut
	}
	if shortcut := normalizeControlShortcut(strings.ToLower(value)); shortcut != "" {
		return shortcut
	}
	if strings.EqualFold(value, "Esc") {
		return "esc"
	}
	return value
}

func trimAnswerEdges(value string) string {
	return trimEdges(value, " \t\r\n")
}

func trimEdges(value, padding string) string {
	value = strings.TrimLeft(value, padding)
	for len(value) > 0 && strings.ContainsRune(padding, rune(value[len(value)-1])) {
		// A trailing escaped space is an argument character, not padding.
		slashes := 0
		for i := len(value) - 2; i >= 0 && value[i] == '\\'; i-- {
			slashes++
		}
		if slashes%2 != 0 {
			break
		}
		value = value[:len(value)-1]
	}
	return value
}

func normalizeShellSpacing(value string) string {
	value = trimAnswerEdges(value)
	for _, char := range value {
		if !(char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' ||
			char >= '0' && char <= '9' || strings.ContainsRune("_-./:= \t|", char)) {
			return value
		}
	}
	if strings.Contains(value, "||") {
		return value
	}
	parts := strings.Split(value, "|")
	for i := range parts {
		parts[i] = strings.Join(strings.Fields(parts[i]), " ")
		if parts[i] == "" {
			return value
		}
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

	for _, prefix := range []string{"ctrl", "control"} {
		if !strings.HasPrefix(value, prefix) {
			continue
		}
		rest := strings.TrimPrefix(value, prefix)
		if rest == "" || !strings.ContainsRune(" +-\t", rune(rest[0])) {
			continue
		}
		rest = strings.Trim(rest, " \t")
		if strings.HasPrefix(rest, "+") || strings.HasPrefix(rest, "-") {
			rest = strings.Trim(rest[1:], " \t")
		}
		if len(rest) == 1 && isASCIILetter(rest[0]) {
			return "ctrl+" + rest
		}
	}
	return ""
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
