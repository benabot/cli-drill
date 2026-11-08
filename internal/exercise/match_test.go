package exercise

import "testing"

func TestMatchAnswerPreservesLiteralMeaning(t *testing.T) {
	for _, tc := range []struct{ expected, input string }{
		{"bat README.md", "bat readme.md"},
		{"grep ERROR fichier.log", "grep error fichier.log"},
		{"ls -A", "ls -a"},
		{"git status", "GIT STATUS"},
		{`echo "hello world"`, `echo "hello  world"`},
		{`echo 'a | b'`, `echo 'a|b'`},
		{`echo a\| b`, `echo a\|b`},
		{`echo hello\ `, `echo hello\`},
		{"echo hello world", "echo hello\nworld"},
		{"echo hello world", "echo hello\u00a0world"},
		{"echo ok > out", "echo ok >> out"},
		{"cat < in", "cat > in"},
		{"echo 2 > out", "echo 2> out"},
		{"echo ok | cat", "echo ok || cat"},
		{"echo ok || cat", "echo ok | | cat"},
		{`echo '$HOME'`, `echo "$HOME"`},
		{"cmd -a file", "cmd file -a"},
		{"cmd -a file", "cmd -a other"},
		{"ctrl a", "^A"},
		{"h", "H"},
		{"", ""},
		{"pwd", ""},
	} {
		t.Run(tc.expected+"/"+tc.input, func(t *testing.T) {
			for _, match := range []func(string, Answer) bool{MatchAnswer, MatchShellAnswer} {
				if match(tc.input, Answer{Primary: tc.expected}) {
					t.Fatalf("incorrect response %q matched %q", tc.input, tc.expected)
				}
			}
		})
	}
}

func TestMatchKeySequenceAcceptsPrimaryAndVariants(t *testing.T) {
	answer := Answer{
		Primary:  "Ctrl+A",
		Accepted: []string{"^A", "control-a", "ctrl a"},
	}

	for _, input := range []string{"Ctrl+A", " ctrl+a ", "^a", "CONTROL-A", "ctrl   a"} {
		if !MatchKeySequence(input, answer) {
			t.Fatalf("expected %q to match %#v", input, answer)
		}
	}
}

func TestMatchShellAnswerNormalizesSimpleSpacing(t *testing.T) {
	answer := Answer{Primary: "fd -e md | fzf"}

	for _, input := range []string{"fd -e md|fzf", " fd   -e   md  |  fzf ", "fd\t-e md | fzf\r\n"} {
		if !MatchShellAnswer(input, answer) {
			t.Fatalf("expected %q to match shell answer", input)
		}
	}
}

func TestTextAnswersAcceptUnescapedEdgeTabs(t *testing.T) {
	answer := Answer{Primary: "rg TODO | fzf", Accepted: []string{`rg "TODO" | fzf`}}
	if !MatchConceptAnswer("\tfzf\t\r\n", Answer{Primary: "fzf"}) ||
		!MatchShellAnswer("\trg \"TODO\" | fzf\t\r\n", answer) ||
		!MatchAnswer("\tmdpreview\t", Answer{Primary: "mdpreview"}) {
		t.Fatal("non-significant edge tabs rejected")
	}
	for _, tc := range []struct{ primary, incorrect string }{
		{"echo \"hello\tworld\"", `echo "hello world"`},
		{"echo hello\\\t", `echo hello\`},
	} {
		if !MatchShellAnswer(tc.primary, Answer{Primary: tc.primary}) ||
			MatchShellAnswer(tc.incorrect, Answer{Primary: tc.primary}) {
			t.Fatalf("significant tab changed: %q", tc.primary)
		}
	}
}

func TestMatchKeySequenceNormalizesControlShortcutVariants(t *testing.T) {
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
			if !MatchKeySequence(tt.input, Answer{Primary: tt.primary}) {
				t.Fatalf("expected %q to match %q", tt.input, tt.primary)
			}
		})
	}
}

func TestMatchKeySequenceRejectsUnrelatedAndMalformedKeys(t *testing.T) {
	for _, input := range []string{"L", "l", "Ctrl+A", "ctrl++l", "ctrl--l", "ctrll", ""} {
		if MatchKeySequence(input, Answer{Primary: "Ctrl+L"}) {
			t.Fatalf("expected %q to be rejected for Ctrl+L", input)
		}
	}
	if MatchKeySequence("H", Answer{Primary: "h"}) {
		t.Fatal("bare keys must preserve case")
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

func TestMatchKeySequenceMatchesRawControlBytes(t *testing.T) {
	for key := byte(1); key <= 26; key++ {
		notation, _ := KeyByteNotation(key)
		if !MatchKeySequence(string(key), Answer{Primary: notation}) {
			t.Fatalf("raw %#x did not match %s", key, notation)
		}
	}
	if !MatchKeySequence(" \x0c ", Answer{Primary: "Ctrl+L"}) ||
		!MatchKeySequence(" \t ", Answer{Primary: "Ctrl+I"}) ||
		!MatchKeySequence("\x1b", Answer{Primary: "Esc"}) {
		t.Fatal("lost raw key or edge padding support")
	}
}

func TestMatchersAcceptOnlyExplicitComplexVariants(t *testing.T) {
	for _, primary := range []string{`echo "hello  world"`, `echo 'a | b'`, `echo a\ b`,
		`echo hello\ `, "echo ok > out", "echo ok >> out", "cat < in", `glow "$(fd -e md | fzf)"`} {
		for _, match := range []func(string, Answer) bool{MatchAnswer, MatchShellAnswer} {
			if !match(primary, Answer{Primary: primary}) || !match("  "+primary+"  \n", Answer{Primary: primary}) {
				t.Fatalf("lost literal or edge-padded answer %q", primary)
			}
		}
	}
	answer := Answer{Primary: "fd -t d | fzf", Accepted: []string{"fd --type d | fzf", ""}}
	for _, input := range []string{answer.Primary, answer.Accepted[0]} {
		if !MatchShellAnswer(input, answer) {
			t.Fatalf("lost explicit variant %q", input)
		}
	}
	for _, input := range []string{"fd -d | fzf", "fd -t f | fzf", "find . | fzf", "", "   "} {
		if MatchShellAnswer(input, answer) {
			t.Fatalf("invented variant %q", input)
		}
	}
}

func TestConceptNamesHaveLimitedCaseTolerance(t *testing.T) {
	for _, input := range []string{"FD", " fd \r\n"} {
		if !MatchConceptAnswer(input, Answer{Primary: "fd"}) {
			t.Fatalf("pedagogical tool name rejected: %q", input)
		}
	}
	for _, tc := range []struct{ primary, input string }{
		{"grep ERROR file", "grep error file"}, {"README.md", "readme.md"}, {"ls -A", "ls -a"},
	} {
		if MatchConceptAnswer(tc.input, Answer{Primary: tc.primary}) {
			t.Fatalf("name tolerance changed a command or path: %q", tc.input)
		}
	}
}

func TestConceptAnswerPreservesExplicitShortcutVariant(t *testing.T) {
	answer := Answer{Primary: "atuin", Accepted: []string{"Ctrl+R"}}
	for _, input := range []string{"atuin", "ATUIN", "Ctrl+R", "control-r", "^R", "\x12"} {
		if !MatchConceptAnswer(input, answer) {
			t.Fatalf("declared shortcut lost its notation variants: %q", input)
		}
	}
	for _, input := range []string{"R", "Ctrl+A"} {
		if MatchConceptAnswer(input, answer) {
			t.Fatalf("unrelated key matched conceptual answer: %q", input)
		}
	}
}

func TestMatchAnswerRejectsUnknownAnswer(t *testing.T) {
	answer := Answer{Primary: "rg"}

	if MatchAnswer("fd", answer) {
		t.Fatal("expected unrelated answer to be rejected")
	}
}
