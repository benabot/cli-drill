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
		{"first eligible text", []string{"train"}, "Chapitre : Premiers pas Unix", false},
		{"exact override", []string{"train", "02-navigation-shell", "--difficulty", "beginner"}, "1/2", false},
		{"override discarded", []string{"train", "02-navigation-shell"}, "1/1", false},
		{"reading without Bat", []string{"train", "07-lecture-preview", "--difficulty", "beginner", "--mode", "learn"}, "Exercice 1/4\n\nOuvre serveur.log dans less", false},
		{"editing without Micro", []string{"train", "08-micro", "--difficulty", "beginner", "--mode", "learn"}, "Exercice 1/7\n\nOuvre notes.txt avec Nano", false},
		{"advanced Markdown without Glow", []string{"train", "09-markdown-glow", "--difficulty", "advanced"}, "1/5", false},
		{"new intermediate exercises", []string{"train", "00-premiers-pas"}, "1/14", false},
		{"invalid difficulty", []string{"train", "--difficulty", "expert"}, "choose beginner, intermediate or advanced", true},
		{"explicit empty difficulty", []string{"train", "--difficulty="}, "invalid difficulty", true},
		{"advanced Unix without modules", []string{"train", "--difficulty", "advanced"}, "Chapitre : Premiers pas Unix", false},
		{"mode on independent chapter", []string{"train", "02-navigation-shell", "--mode", "learn"}, "only available for guided chapters", true},
		{"unknown mode", []string{"train", "00-premiers-pas", "--mode", "other"}, "choose learn or review", true},
		{"learn", []string{"train", "00-premiers-pas", "--difficulty", "beginner", "--mode", "learn"}, "1/12", false},
		{"review", []string{"train", "00-premiers-pas", "--difficulty", "beginner", "--mode", "review"}, "1/12", false},
		{"mixed level requires terminal", []string{"train", "01-raccourcis-terminal"}, "requires a terminal", true},
		{"eligible capture requires terminal", []string{"train", "01-raccourcis-terminal", "--difficulty", "beginner"}, "requires a terminal", true},
		{"guided diagnostics without terminal", []string{"train", "01-raccourcis-terminal", "--difficulty", "advanced", "--mode", "learn"}, "Exercice 1/8\n\nDans Bash avec Readline Emacs standard, la recherche", false},
		{"review diagnostics without terminal", []string{"train", "01-raccourcis-terminal", "--difficulty", "advanced", "--mode", "review"}, "Exercice 1/8", false},
		{"Unix text guided", []string{"train", "14-recherche-texte-unix", "--difficulty", "beginner", "--mode", "learn"}, "Exercice 1/14\n\nDans api.log", false},
		{"Unix text review", []string{"train", "14-recherche-texte-unix", "--mode", "review"}, "Exercice 1/22", false},
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
	for _, want := range []string{"Active acquired (intermediate, enabled CLI modules): 1/117", "Historical totals (all sessions)", "Completed exercises: 3", "Chapters with attempts: 2", "Streak: 4"} {
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
