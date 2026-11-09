package xdg

import (
	"os"
	"path/filepath"
)

func ProgressFile() string {
	return filepath.Join(dataDir(), "cli-drill", "progress.json")
}

func PreferencesFile() string {
	return filepath.Join(dataDir(), "cli-drill", "preferences.json")
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
