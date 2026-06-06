package exercise

import "testing"

func TestMatchAnswerAcceptsPrimaryAndVariants(t *testing.T) {
	answer := Answer{
		Primary:  "Ctrl+A",
		Accepted: []string{"^A", "control-a", "ctrl a"},
	}

	for _, input := range []string{"Ctrl+A", " ctrl+a ", "^a", "CONTROL-A", "ctrl   a"} {
		if !MatchAnswer(input, answer) {
			t.Fatalf("expected %q to match %#v", input, answer)
		}
	}
}

func TestMatchAnswerNormalizesSimpleShellSpacing(t *testing.T) {
	answer := Answer{Primary: "fd -e md | fzf"}

	for _, input := range []string{"fd -e md|fzf", " fd   -e   md  |  fzf "} {
		if !MatchAnswer(input, answer) {
			t.Fatalf("expected %q to match shell answer", input)
		}
	}
}

func TestMatchAnswerNormalizesControlShortcutVariants(t *testing.T) {
	tests := []struct {
		name    string
		primary string
		input   string
	}{
		{name: "ctrl l accepts caret notation", primary: "Ctrl+L", input: "^L"},
		{name: "ctrl a accepts caret notation", primary: "Ctrl+A", input: "^A"},
		{name: "ctrl w accepts spaced ctrl", primary: "Ctrl+W", input: "ctrl w"},
		{name: "ctrl r accepts control dash", primary: "Ctrl+R", input: "control-r"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !MatchAnswer(tt.input, Answer{Primary: tt.primary}) {
				t.Fatalf("expected %q to match %q", tt.input, tt.primary)
			}
		})
	}
}

func TestMatchAnswerRejectsBareLetterForControlShortcut(t *testing.T) {
	if MatchAnswer("L", Answer{Primary: "Ctrl+L"}) {
		t.Fatal("expected bare letter to be rejected for Ctrl+L")
	}
}

func TestMatchAnswerRejectsUnknownAnswer(t *testing.T) {
	answer := Answer{Primary: "rg"}

	if MatchAnswer("fd", answer) {
		t.Fatal("expected unrelated answer to be rejected")
	}
}
