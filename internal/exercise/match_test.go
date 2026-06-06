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

func TestKeyByteNotationMapsControlBytes(t *testing.T) {
	tests := []struct {
		name string
		b    byte
		want string
	}{
		{name: "ctrl a", b: 0x01, want: "Ctrl+A"},
		{name: "ctrl e", b: 0x05, want: "Ctrl+E"},
		{name: "ctrl u", b: 0x15, want: "Ctrl+U"},
		{name: "ctrl k", b: 0x0b, want: "Ctrl+K"},
		{name: "ctrl w", b: 0x17, want: "Ctrl+W"},
		{name: "ctrl y", b: 0x19, want: "Ctrl+Y"},
		{name: "ctrl r", b: 0x12, want: "Ctrl+R"},
		{name: "ctrl l", b: 0x0c, want: "Ctrl+L"},
		{name: "ctrl c", b: 0x03, want: "Ctrl+C"},
		{name: "ctrl z", b: 0x1a, want: "Ctrl+Z"},
		{name: "esc", b: 0x1b, want: "Esc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := KeyByteNotation(tt.b)
			if !ok || got != tt.want {
				t.Fatalf("KeyByteNotation(%#x) = %q, %v; want %q, true", tt.b, got, ok, tt.want)
			}
		})
	}
}

func TestMatchAnswerMatchesRawControlBytes(t *testing.T) {
	if !MatchAnswer("\x01", Answer{Primary: "Ctrl+A"}) {
		t.Fatal("expected raw Ctrl+A byte to match Ctrl+A")
	}
	if !MatchAnswer("\x0c", Answer{Primary: "Ctrl+L"}) {
		t.Fatal("expected raw Ctrl+L byte to match Ctrl+L")
	}
	if MatchAnswer("L", Answer{Primary: "Ctrl+L"}) {
		t.Fatal("expected bare L to be rejected for Ctrl+L")
	}
}

func TestMatchAnswerRejectsUnknownAnswer(t *testing.T) {
	answer := Answer{Primary: "rg"}

	if MatchAnswer("fd", answer) {
		t.Fatal("expected unrelated answer to be rejected")
	}
}
