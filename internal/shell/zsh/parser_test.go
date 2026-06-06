package zsh

import (
	"strings"
	"testing"
)

func TestParseAliasesDetectsSimpleAliasesAndIgnoresComments(t *testing.T) {
	input := strings.NewReader(`
# alias old='rm -rf nope'
alias cgs='git status --short'
alias lg="lazygit"
  alias grep='rg'
not_an_alias foo='bar'
`)

	aliases, err := ParseAliases(input, "aliases.zsh")
	if err != nil {
		t.Fatalf("ParseAliases returned error: %v", err)
	}

	if len(aliases) != 3 {
		t.Fatalf("expected 3 aliases, got %d: %#v", len(aliases), aliases)
	}

	assertAlias(t, aliases[0], "cgs", "git status --short", 3)
	assertAlias(t, aliases[1], "lg", "lazygit", 4)
	assertAlias(t, aliases[2], "grep", "rg", 5)
}

func TestParseFunctionsDetectsSupportedZSHFunctionForms(t *testing.T) {
	input := strings.NewReader(`
# old_func() {
y() {
  cd "$HOME"
}

function preview() {
  glow "$1"
}

function mkcd {
  mkdir -p "$1"
}
`)

	functions, err := ParseFunctions(input, "functions.zsh")
	if err != nil {
		t.Fatalf("ParseFunctions returned error: %v", err)
	}

	if len(functions) != 3 {
		t.Fatalf("expected 3 functions, got %d: %#v", len(functions), functions)
	}

	assertFunction(t, functions[0], "y", 3)
	assertFunction(t, functions[1], "preview", 7)
	assertFunction(t, functions[2], "mkcd", 11)
}

func assertAlias(t *testing.T, alias Alias, name, value string, line int) {
	t.Helper()
	if alias.Name != name || alias.Value != value || alias.Line != line {
		t.Fatalf("unexpected alias: got %#v, want name=%q value=%q line=%d", alias, name, value, line)
	}
}

func assertFunction(t *testing.T, function Function, name string, line int) {
	t.Helper()
	if function.Name != name || function.Line != line {
		t.Fatalf("unexpected function: got %#v, want name=%q line=%d", function, name, line)
	}
}
