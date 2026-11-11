package app

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/spf13/cobra"
)

func choicesCLICommand(path string) (*bytes.Buffer, *cobra.Command) {
	cmd := NewRootCommand(fstest.MapFS{
		"choices.yaml": {Data: []byte(`id: choices
title: Synthetic choices
items:
  - id: search
    type: tool
    exercise_type: multiple-choice
    prompt: Which search tool?
    choices: [grep, ripgrep, fd]
    answer:
      primary: rg
      accepted: [ripgrep, historical-rg]
`)},
	})
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetErr(&output)
	cmd.SetArgs([]string{"--progress", path, "train", "choices"})
	return &output, cmd
}

func TestCLIChoicesAreVisibleBeforeReadingAnswer(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")
	output, cmd := choicesCLICommand(path)
	cmd.SetIn(recordingReader(func([]byte) (int, error) {
		text := output.String()
		prompt := strings.Index(text, "Which search tool?")
		input := strings.LastIndex(text, "> ")
		if prompt < 0 || input <= prompt {
			t.Fatalf("missing question or input prompt: %q", text)
		}
		previous := prompt
		for _, choice := range []string{"grep", "ripgrep", "fd"} {
			position := strings.Index(text, choice)
			if position <= previous || position >= input {
				t.Errorf("choice %q must appear in YAML order between question and input prompt: %q", choice, text)
			}
			previous = position
		}
		return 0, io.EOF
	}))
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("viewing choices must not create progress: %v", err)
	}
}

func TestCLIChoicesKeepCorrectionAndV1Recording(t *testing.T) {
	for _, tc := range []struct {
		name              string
		inputs            []string
		attempts, correct int
	}{
		{"primary outside choices", []string{"rg\n"}, 1, 1},
		{"displayed accepted label", []string{"ripgrep\n"}, 1, 1},
		{"historical variant outside choices", []string{"historical-rg\n"}, 1, 1},
		{"concept case rule", []string{"RIPGREP\n"}, 1, 1},
		{"wrong displayed label", []string{"grep\n"}, 1, 0},
		{"numeric position is not an answer", []string{"2\n"}, 1, 0},
		{"nonempty EOF submitted once", []string{"ripgrep"}, 1, 1},
		{"wrong EOF submitted once", []string{"fd"}, 1, 0},
		{"empty EOF", []string{""}, 0, 0},
		{"blank lines and empty EOF", []string{"\n \n"}, 0, 0},
		{"spaces before EOF", []string{"  "}, 0, 0},
		{"blank retries then answer", []string{"\n \nripgrep\n"}, 1, 1},
		{"escape", []string{"\x1b"}, 0, 0},
		{"ctrl c with unfinished answer", []string{"ripgrep\x03"}, 0, 0},
		{"repeated acquisition and wrong answer", []string{"ripgrep\n", "rg\n", "historical-rg\n", "grep\n"}, 4, 3},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "progress.json")
			if err := os.WriteFile(path, []byte(recordingHistory), 0o600); err != nil {
				t.Fatal(err)
			}
			correctFeedback, wrongFeedback := 0, 0
			for _, input := range tc.inputs {
				output, cmd := choicesCLICommand(path)
				cmd.SetArgs([]string{"--progress", path, "train", "choices", "--mode", "review"})
				cmd.SetIn(strings.NewReader(input))
				if err := cmd.Execute(); err != nil {
					t.Fatal(err)
				}
				correctFeedback += strings.Count(output.String(), "Correct.")
				wrongFeedback += strings.Count(output.String(), "Réponse incorrecte.\nRéponse attendue : rg")
			}
			if correctFeedback != tc.correct || wrongFeedback != tc.attempts-tc.correct {
				t.Fatalf("feedback: got %d correct/%d wrong, want %d/%d", correctFeedback, wrongFeedback, tc.correct, tc.attempts-tc.correct)
			}
			var acquired []string
			if tc.correct > 0 {
				acquired = []string{"choices/search"}
			}
			assertRecordingScore(t, path, "choices", tc.attempts, tc.correct, acquired...)
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if tc.attempts == 0 && string(data) != recordingHistory {
				t.Fatalf("empty or interrupted input rewrote historical progress: %s", data)
			}
			var fields map[string]json.RawMessage
			if err := json.Unmarshal(data, &fields); err != nil || len(fields) != 5 || string(fields["version"]) != "1" {
				t.Fatalf("v1 JSON schema changed: %s (%v)", data, err)
			}
		})
	}
}
