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

func TestMatchAnswerRejectsUnknownAnswer(t *testing.T) {
	answer := Answer{Primary: "rg"}

	if MatchAnswer("fd", answer) {
		t.Fatal("expected unrelated answer to be rejected")
	}
}
