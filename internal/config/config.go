package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	DotfilesPath string         `toml:"dotfiles_path"`
	Shell        string         `toml:"shell"`
	Paths        PathConfig     `toml:"paths"`
	Security     SecurityConfig `toml:"security"`
}

type PathConfig struct {
	Aliases   []string `toml:"aliases"`
	Functions []string `toml:"functions"`
	Docs      []string `toml:"docs"`
}

type SecurityConfig struct {
	Exclude []string `toml:"exclude"`
}

func Default() Config {
	return Config{
		DotfilesPath: "~/dotfiles",
		Shell:        "zsh",
		Paths: PathConfig{
			Aliases:   []string{"zsh/modules/aliases.zsh"},
			Functions: []string{"zsh/modules/functions.zsh"},
			Docs:      []string{"README.md", "docs/tools-inventory.md", "docs/cli-tools-usage.md"},
		},
		Security: SecurityConfig{Exclude: DefaultExcludes()},
	}
}

func DefaultExcludes() []string {
	return []string{
		"~/.config/secrets",
		"~/.ssh",
		"~/.gnupg",
		"~/.config/gh/hosts.yml",
		"~/.config/zed/settings.json",
	}
}

func Decode(r io.Reader) (Config, error) {
	cfg := Default()
	if _, err := toml.NewDecoder(r).Decode(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func Load(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()
	return Decode(file)
}

func Save(path string, cfg Config, force bool) error {
	if path == "" {
		return errors.New("config path is required")
	}
	if !force {
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("config already exists: %s", path)
		} else if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(cfg); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0o644)
}

func (c Config) Validate() error {
	if c.DotfilesPath == "" {
		return errors.New("dotfiles_path is required")
	}
	if c.Shell != "zsh" {
		return fmt.Errorf("unsupported shell %q: MVP supports zsh only", c.Shell)
	}
	return nil
}

func (c Config) Encode() ([]byte, error) {
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(c); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
