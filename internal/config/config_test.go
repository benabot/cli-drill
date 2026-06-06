package config

import (
	"strings"
	"testing"
)

func TestDecodeConfigFromTOML(t *testing.T) {
	input := strings.NewReader(`
dotfiles_path = "~/dotfiles"
shell = "zsh"

[paths]
aliases = ["zsh/modules/aliases.zsh"]
functions = ["zsh/modules/functions.zsh"]
docs = ["README.md"]

[security]
exclude = ["~/.ssh"]
`)

	cfg, err := Decode(input)
	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	if cfg.DotfilesPath != "~/dotfiles" || cfg.Shell != "zsh" {
		t.Fatalf("unexpected config root fields: %#v", cfg)
	}
	if len(cfg.Paths.Aliases) != 1 || cfg.Paths.Aliases[0] != "zsh/modules/aliases.zsh" {
		t.Fatalf("unexpected aliases: %#v", cfg.Paths.Aliases)
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestValidateRejectsUnsupportedShell(t *testing.T) {
	cfg := Default()
	cfg.Shell = "bash"

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected unsupported shell error")
	}
}
