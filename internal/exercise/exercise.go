package exercise

import "strings"

type Type string

const (
	TypeFreeAnswer     Type = "free-answer"
	TypeMultipleChoice Type = "multiple-choice"
	TypeScenario       Type = "scenario"
	TypeSimpleShellSim Type = "simple-shell-sim"
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
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.Join(strings.Fields(value), " ")
	value = normalizePipeSpacing(value)
	return value
}

func normalizePipeSpacing(value string) string {
	parts := strings.Split(value, "|")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return strings.Join(parts, "|")
}
