package app

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/benabot/cli-drill/data"
	"github.com/benabot/cli-drill/internal/tui"
	"github.com/benabot/cli-drill/internal/xdg"
	"github.com/charmbracelet/x/ansi"
)

func TestRootDiagnosesUnavailableTUI(t *testing.T) {
	for _, tc := range []struct {
		name, term, reason string
		tty                bool
	}{
		{"pipes", "xterm-256color", "terminal for input and output", false},
		{"dumb", "dumb", "TERM=dumb", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("TERM", tc.term)
			t.Cleanup(stubTerminalDetector(func(uintptr) bool { return tc.tty }))
			previous := runTUIWithOptions
			runTUIWithOptions = func(tui.Options) (tui.Result, error) {
				t.Error("unavailable terminal reached the interactive renderer")
				return tui.Result{}, nil
			}
			t.Cleanup(func() { runTUIWithOptions = previous })
			var out fdBuffer
			cmd := NewRootCommand(data.Chapters())
			cmd.SetArgs([]string{})
			cmd.SetIn(&fdBuffer{})
			cmd.SetOut(&out)
			cmd.SetErr(&out)
			err := cmd.Execute()
			if err == nil || !strings.Contains(err.Error(), tc.reason) ||
				!strings.Contains(err.Error(), "cli-drill train") ||
				strings.Contains(out.String(), "\x1b") {
				t.Fatalf("missing usable text fallback: %v\n%s", err, out.String())
			}
		})
	}
}

func TestDumbCapturePrintsWithoutANSI(t *testing.T) {
	t.Setenv("TERM", "dumb")
	t.Setenv("CLICOLOR_FORCE", "1")
	t.Cleanup(stubTerminalDetector(func(uintptr) bool { return true }))
	var out fdBuffer
	opts := Options{Out: &out}
	writeKeySequenceScreen(opts, keySequenceStyleFor(opts).title("Test"))
	if out.String() != "Test" {
		t.Fatalf("TERM=dumb capture sent controls: %q", out.String())
	}
}

func TestMissingModulesExplainAlternativesAndActivation(t *testing.T) {
	isolateChapterHome(t)
	var out bytes.Buffer
	cmd := NewRootCommand(data.Chapters())
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"train", "06-recherche-fichiers-contenu", "--difficulty", "intermediate"})
	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "fd + fzf") ||
		!strings.Contains(err.Error(), "rg") || !strings.Contains(err.Error(), "cli-drill prefs --enable fd") {
		t.Fatalf("missing module diagnostic: %v\n%s", err, &out)
	}
	if _, err := os.Stat(xdg.PreferencesFile()); !os.IsNotExist(err) {
		t.Fatalf("the diagnostic must not create preferences: %v", err)
	}
}

func TestChaptersFitsMeasuredTerminalWidth(t *testing.T) {
	t.Cleanup(stubTerminalDetector(func(fd uintptr) bool { return fd == 42 }))
	previousSize := terminalSize
	t.Cleanup(func() { terminalSize = previousSize })
	for _, width := range []int{80, 60, 40} {
		t.Run(fmt.Sprint(width), func(t *testing.T) {
			terminalSize = func(uintptr) (int, int, error) { return width, 24, nil }
			out := &fdBuffer{fd: 42}
			cmd := NewRootCommand(data.Chapters())
			cmd.SetOut(out)
			cmd.SetArgs([]string{"chapters"})
			if err := cmd.Execute(); err != nil {
				t.Fatal(err)
			}
			listing := out.String()
			for _, line := range strings.Split(listing, "\n") {
				if ansi.StringWidth(line) > width || strings.ContainsAny(line, "\t\x1b") {
					t.Fatalf("unreadable row at %d columns: %q", width, line)
				}
			}
			for _, label := range []string{"14-recherche-texte-unix", "Recherche et traitement de texte Unix", "48"} {
				if !strings.Contains(listing, label) {
					t.Errorf("missing %q at %d columns", label, width)
				}
			}
		})
	}
}
