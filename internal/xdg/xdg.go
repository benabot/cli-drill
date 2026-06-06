package xdg

import (
	"os"
	"path/filepath"
)

func ConfigFile() string {
	return filepath.Join(configDir(), "cli-drill", "config.toml")
}

func ChaptersDir() string {
	return filepath.Join(configDir(), "cli-drill", "chapters")
}

func ProgressFile() string {
	return filepath.Join(dataDir(), "cli-drill", "progress.json")
}

func configDir() string {
	if dir, err := os.UserConfigDir(); err == nil && dir != "" {
		return dir
	}
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".config")
	}
	return "."
}

func dataDir() string {
	if dir := os.Getenv("XDG_DATA_HOME"); dir != "" {
		return dir
	}
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".local", "share")
	}
	return "."
}
