package app

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/benabot/cli-drill/data"
	"github.com/benabot/cli-drill/internal/xdg"
)

func TestCLISessionDifficultyAndModules(t *testing.T) {
	isolateChapterHome(t)
	prefs := writeAppTestFileContent(t, xdg.PreferencesFile(), `{"version":1,"enabled_tools":["fd","fzf","retired-tool"],"default_difficulty":"intermediate"}`)
	path := filepath.Join(t.TempDir(), "progress.json")

	for _, tc := range []struct {
		name     string
		args     []string
		want     string
		rejected bool
	}{
		{"first eligible text", []string{"train"}, "1/1", false},
		{"exact override", []string{"train", "02-navigation-shell", "--difficulty", "beginner"}, "1/2", false},
		{"override discarded", []string{"train", "02-navigation-shell"}, "1/1", false},
		{"advanced blocked by modules", []string{"train", "09-markdown-glow", "--difficulty", "advanced"}, "enabled CLI modules", true},
		{"level missing", []string{"train", "00-premiers-pas"}, "No exercises at intermediate difficulty", true},
		{"invalid difficulty", []string{"train", "--difficulty", "expert"}, "choose beginner, intermediate or advanced", true},
		{"explicit empty difficulty", []string{"train", "--difficulty="}, "invalid difficulty", true},
		{"advanced Unix without modules", []string{"train", "--difficulty", "advanced"}, "Chapitre : Flux du shell", false},
		{"mode on independent chapter", []string{"train", "02-navigation-shell", "--mode", "learn"}, "only available for guided chapters", true},
		{"unknown mode", []string{"train", "00-premiers-pas", "--mode", "other"}, "choose learn or review", true},
		{"learn", []string{"train", "00-premiers-pas", "--difficulty", "beginner", "--mode", "learn"}, "1/8", false},
		{"review", []string{"train", "00-premiers-pas", "--difficulty", "beginner", "--mode", "review"}, "1/8", false},
		{"capture empty before terminal check", []string{"train", "01-raccourcis-terminal"}, "No exercises at intermediate difficulty", true},
		{"eligible capture requires terminal", []string{"train", "01-raccourcis-terminal", "--difficulty", "beginner"}, "requires a terminal", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			output, err := runSessionCLI(path, "", tc.args...)
			if (err != nil) != tc.rejected || !strings.Contains(output, tc.want) {
				t.Fatalf("Execute = %v, output %q (want %q)", err, output, tc.want)
			}
			if _, err := os.Stat(path); !os.IsNotExist(err) {
				t.Fatalf("empty input or selection wrote progress: %v", err)
			}
			assertPreferencesBytes(t, xdg.PreferencesFile(), prefs)
		})
	}
}

func TestCLIStatsProjectsActiveAcquiredWithoutWritingHistory(t *testing.T) {
	isolateChapterHome(t)
	writeAppTestFileContent(t, xdg.PreferencesFile(), `{"version":1,"enabled_tools":["fd","fzf","retired-tool"],"default_difficulty":"intermediate"}`)
	path := filepath.Join(t.TempDir(), "progress.json")
	historical := writeAppTestFileContent(t, path, `{"version":1,"completed_exercises":["02-navigation-shell/fd-fzf-dir","02-navigation-shell/retired","00-premiers-pas/pwd"],"chapter_scores":{"02-navigation-shell":{"attempts":7,"correct":3},"retired":{"attempts":2,"correct":1}},"last_session":null,"streak":4}`)
	output, err := runSessionCLI(path, "", "stats")
	for _, want := range []string{"Active acquired (intermediate, enabled CLI modules): 1/12", "Historical totals (all sessions)", "Completed exercises: 3", "Chapters with attempts: 2", "Streak: 4"} {
		if err != nil || !strings.Contains(output, want) {
			t.Fatalf("stats = %q (%v), want %q", output, err, want)
		}
	}
	assertPreferencesBytes(t, path, historical)
}

func runSessionCLI(path, input string, args ...string) (string, error) {
	cmd := NewRootCommand(data.Chapters())
	var out bytes.Buffer
	cmd.SetIn(strings.NewReader(input))
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs(append([]string{"--progress", path}, args...))
	err := cmd.Execute()
	return out.String(), err
}

func writeAppTestFileContent(t *testing.T, path, content string) []byte {
	t.Helper()
	writeAppTestFile(t, filepath.Dir(path), filepath.Base(path), content)
	return []byte(content)
}
