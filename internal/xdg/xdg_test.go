package xdg

import (
	"path/filepath"
	"testing"
)

func TestProgressFilePreservesExistingLocation(t *testing.T) {
	home := t.TempDir()
	data := t.TempDir()
	for _, tc := range []struct {
		name, home, data, want string
	}{
		{"default", home, "", filepath.Join(home, ".local", "share", "cli-drill", "progress.json")},
		{"XDG override", home, data, filepath.Join(data, "cli-drill", "progress.json")},
		{"unavailable home", "", "", filepath.Join("cli-drill", "progress.json")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("HOME", tc.home)
			t.Setenv("USERPROFILE", tc.home)
			t.Setenv("XDG_DATA_HOME", tc.data)
			if got := ProgressFile(); got != tc.want {
				t.Fatalf("ProgressFile() = %q, want existing location %q", got, tc.want)
			}
		})
	}
}
