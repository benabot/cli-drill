package app

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/benabot/cli-drill/data"
	"github.com/benabot/cli-drill/internal/tui"
)

type fdBuffer struct {
	bytes.Buffer
	fd uintptr
}

func (b fdBuffer) Fd() uintptr {
	return b.fd
}

func TestDirectoryUsesScanCatalogWhenConfigIsProvided(t *testing.T) {
	root := t.TempDir()
	writeAppTestFile(t, root, "zsh/modules/aliases.zsh", "alias cgs='git status --short'\n")
	configPath := writeAppTestConfig(t, root)

	var out bytes.Buffer
	cmd := NewRootCommand(data.Chapters())
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"--config", configPath, "directory", "--type", "alias"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v\n%s", err, out.String())
	}

	got := out.String()
	if !bytes.Contains([]byte(got), []byte("alias\tcgs\tcgs")) {
		t.Fatalf("expected scanned alias in directory output, got:\n%s", got)
	}
}

func TestScanSummaryPrintsCountsWithoutEntries(t *testing.T) {
	root := t.TempDir()
	writeAppTestFile(t, root, "zsh/modules/aliases.zsh", "alias cgs='git status --short'\n")
	configPath := writeAppTestConfig(t, root)

	out, err := executeTestCommand("--config", configPath, "scan", "--summary")
	if err != nil {
		t.Fatalf("Execute returned error: %v\n%s", err, out)
	}

	for _, want := range []string{"aliases: 1", "functions: 0", "tools: 0", "concepts: 0", "workflows: 0", "chapters: 0", "total: 1"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected summary to contain %q, got:\n%s", want, out)
		}
	}
	if strings.Contains(out, "cgs\talias") {
		t.Fatalf("summary should not print entries, got:\n%s", out)
	}
}

func TestScanTypeFiltersAndRejectsUnknownType(t *testing.T) {
	root := t.TempDir()
	writeAppTestFile(t, root, "zsh/modules/aliases.zsh", "alias cgs='git status --short'\n")
	configPath := writeAppTestConfig(t, root)

	out, err := executeTestCommand("--config", configPath, "scan", "--type", "alias")
	if err != nil {
		t.Fatalf("Execute returned error: %v\n%s", err, out)
	}
	if !strings.Contains(out, "alias\tcgs\tcgs\tgit status --short") {
		t.Fatalf("expected alias output, got:\n%s", out)
	}

	out, err = executeTestCommand("--config", configPath, "scan", "--type", "nope")
	if err == nil {
		t.Fatalf("expected invalid type error, got output:\n%s", out)
	}
	if !strings.Contains(err.Error(), "unknown entry type") {
		t.Fatalf("expected clear invalid type error, got: %v", err)
	}
}

func TestShowFindsEntryByIDAndName(t *testing.T) {
	root := t.TempDir()
	writeAppTestFile(t, root, "zsh/modules/aliases.zsh", "alias cgs='git status --short'\n")
	configPath := writeAppTestConfig(t, root)

	for _, query := range []string{"cgs", "cgs"} {
		out, err := executeTestCommand("--config", configPath, "show", query)
		if err != nil {
			t.Fatalf("Execute returned error: %v\n%s", err, out)
		}
		for _, want := range []string{"id: cgs", "name: cgs", "type: alias", "command: git status --short", "source: zsh/modules/aliases.zsh:1"} {
			if !strings.Contains(out, want) {
				t.Fatalf("expected show output to contain %q, got:\n%s", want, out)
			}
		}
	}
}

func TestShowReportsAmbiguousName(t *testing.T) {
	root := t.TempDir()
	writeAppTestFile(t, root, "zsh/modules/aliases.zsh", "alias dup='git status --short'\n")
	writeAppTestFile(t, root, "zsh/modules/functions.zsh", "dup() {\n  echo dup\n}\n")
	configPath := writeAppTestConfigWithFunctions(t, root)

	out, err := executeTestCommand("--config", configPath, "show", "dup")
	if err == nil {
		t.Fatalf("expected ambiguous show error, got output:\n%s", out)
	}
	if !strings.Contains(err.Error(), "ambiguous entry") {
		t.Fatalf("expected ambiguous entry error, got: %v", err)
	}
}

func TestTrainMatchesTerminalControlShortcutInput(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	progressPath := filepath.Join(t.TempDir(), "progress.json")
	chapters := fstest.MapFS{
		"shortcuts.yaml": &fstest.MapFile{Data: []byte(`id: shortcuts
title: Shortcuts
items:
  - id: ctrl-a
    type: shortcut
    exercise_type: free-answer
    prompt: Ctrl A?
    answer:
      primary: Ctrl+A
  - id: ctrl-e
    type: shortcut
    exercise_type: free-answer
    prompt: Ctrl E?
    answer:
      primary: Ctrl+E
  - id: ctrl-l
    type: shortcut
    exercise_type: free-answer
    prompt: Ctrl L?
    answer:
      primary: Ctrl+L
  - id: ctrl-l-bare
    type: shortcut
    exercise_type: free-answer
    prompt: Ctrl L bare?
    answer:
      primary: Ctrl+L
`)},
	}

	var out bytes.Buffer
	cmd := NewRootCommand(chapters)
	cmd.SetIn(strings.NewReader("\x01\n\x05\n\x0c\nL\n"))
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"--progress", progressPath, "train", "shortcuts"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v\n%s", err, out.String())
	}

	got := out.String()
	if strings.Count(got, "Correct.") != 3 {
		t.Fatalf("expected three control shortcut answers to be correct, got:\n%s", got)
	}
	if !strings.Contains(got, "Pas encore. Reponse attendue: Ctrl+L") {
		t.Fatalf("expected bare L to be rejected, got:\n%s", got)
	}
}

func TestTrainKeySequenceMatchesOneControlBytePerExercise(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	progressPath := filepath.Join(t.TempDir(), "progress.json")
	chapters := fstest.MapFS{
		"shortcuts.yaml": &fstest.MapFile{Data: []byte(`id: shortcuts
title: Shortcuts
items:
  - id: ctrl-a
    type: shortcut
    exercise_type: key-sequence
    prompt: Press Ctrl+A.
    answer:
      primary: Ctrl+A
  - id: ctrl-l
    type: shortcut
    exercise_type: key-sequence
    prompt: Press Ctrl+L.
    answer:
      primary: Ctrl+L
  - id: ctrl-l-bare
    type: shortcut
    exercise_type: key-sequence
    prompt: Press Ctrl+L.
    answer:
      primary: Ctrl+L
`)},
	}

	var out bytes.Buffer
	cmd := NewRootCommand(chapters)
	cmd.SetIn(strings.NewReader("\x01\n\x0c\nL"))
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"--progress", progressPath, "train", "shortcuts"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v\n%s", err, out.String())
	}

	got := out.String()
	if strings.Count(got, "Correct.") != 2 {
		t.Fatalf("expected two key sequence answers to be correct, got:\n%s", got)
	}
	if !strings.Contains(got, "Recu: L") || !strings.Contains(got, "Pas encore.") {
		t.Fatalf("expected bare L to be rejected, got:\n%s", got)
	}
	if !strings.Contains(got, "Correct: 2/3") || !strings.Contains(got, "À revoir: 1") {
		t.Fatalf("expected chapter summary with one missed exercise, got:\n%s", got)
	}
}

func TestRootRunsCLITrainingAndReturnsToTUIWhenTUIRequestsKeySequenceChapter(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	progressPath := filepath.Join(t.TempDir(), "progress.json")
	chapters := fstest.MapFS{
		"shortcuts.yaml": &fstest.MapFile{Data: []byte(`id: shortcuts
title: Shortcuts
items:
  - id: ctrl-a
    type: shortcut
    exercise_type: key-sequence
    prompt: Press Ctrl+A.
    answer:
      primary: Ctrl+A
`)},
	}
	calls := 0
	restore := stubTUIRunner(func(opts tui.Options) (tui.Result, error) {
		calls++
		if calls == 1 {
			return tui.Result{RunCLITrainChapterID: "shortcuts"}, nil
		}
		return tui.Result{}, nil
	})
	defer restore()

	var out bytes.Buffer
	cmd := NewRootCommand(chapters)
	cmd.SetIn(strings.NewReader("\x01\n"))
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"--progress", progressPath})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v\n%s", err, out.String())
	}

	got := out.String()
	if !strings.Contains(got, "Recu: Ctrl+A") || !strings.Contains(got, "Correct.") {
		t.Fatalf("expected TUI intent to run CLI key-sequence training, got:\n%s", got)
	}
	if calls != 2 {
		t.Fatalf("expected TUI to relaunch after CLI training, calls = %d", calls)
	}
}

func TestCommandBarFormat(t *testing.T) {
	got := commandBar("Enter next", "r retry", "Esc quit")
	want := "────────────────────────────────────────\nEnter next · r retry · Esc quit"
	if got != want {
		t.Fatalf("unexpected command bar:\n%s", got)
	}
}

func TestRenderKeySequenceQuestion(t *testing.T) {
	got := renderKeySequenceQuestion(keySequenceQuestionView{
		Title:  "Shortcuts",
		Index:  1,
		Total:  3,
		Prompt: "Press Ctrl+A.",
	})

	for _, want := range []string{
		"Shortcuts",
		"Progression: 1/3",
		"Press Ctrl+A.",
		"Waiting for key...",
		"h help · Esc quit",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected question render to contain %q, got:\n%s", want, got)
		}
	}
}

func TestRenderKeySequenceQuestionDefaultStyleHasNoANSI(t *testing.T) {
	got := renderKeySequenceQuestion(keySequenceQuestionView{
		Title:  "Shortcuts",
		Index:  1,
		Total:  3,
		Prompt: "Press Ctrl+A.",
	})

	if strings.Contains(got, "\x1b[") {
		t.Fatalf("default render should stay ANSI-free for stable tests, got:\n%q", got)
	}
}

func TestRenderKeySequenceQuestionStyledKeepsReadableText(t *testing.T) {
	got := renderKeySequenceQuestion(keySequenceQuestionView{
		Title:  "Shortcuts",
		Index:  1,
		Total:  3,
		Prompt: "Press Ctrl+A.",
		Style:  keySequenceStyle{Color: true},
	})

	for _, want := range []string{"Shortcuts", "Progression: 1/3", "Press Ctrl+A.", "h help · Esc quit"} {
		if !strings.Contains(got, want) {
			t.Fatalf("styled render should preserve readable text %q, got:\n%s", want, got)
		}
	}
	if !strings.Contains(got, "\x1b[") {
		t.Fatalf("styled render should contain ANSI styling, got:\n%q", got)
	}
}

func TestKeySequenceStyleUsesTerminalFDWriter(t *testing.T) {
	restore := stubTerminalDetector(func(fd uintptr) bool {
		return fd == 42
	})
	defer restore()
	t.Setenv("NO_COLOR", "")
	t.Setenv("TERM", "xterm-256color")

	style := keySequenceStyleFor(Options{Out: &fdBuffer{fd: 42}})

	if !style.Color {
		t.Fatal("expected terminal fd writer to enable key-sequence color")
	}
}

func TestKeySequenceStyleDisablesColorForNoColorAndDumbTerminal(t *testing.T) {
	restore := stubTerminalDetector(func(fd uintptr) bool {
		return fd == 42
	})
	defer restore()

	t.Run("NO_COLOR", func(t *testing.T) {
		t.Setenv("NO_COLOR", "1")
		t.Setenv("TERM", "xterm-256color")
		style := keySequenceStyleFor(Options{Out: &fdBuffer{fd: 42}})
		if style.Color {
			t.Fatal("expected NO_COLOR=1 to disable key-sequence color")
		}
	})

	t.Run("TERM dumb", func(t *testing.T) {
		t.Setenv("NO_COLOR", "")
		t.Setenv("TERM", "dumb")
		style := keySequenceStyleFor(Options{Out: &fdBuffer{fd: 42}})
		if style.Color {
			t.Fatal("expected TERM=dumb to disable key-sequence color")
		}
	})
}

func TestRenderKeySequenceQuestionWithHelp(t *testing.T) {
	got := renderKeySequenceQuestion(keySequenceQuestionView{
		Title:    "Shortcuts",
		Index:    1,
		Total:    3,
		Prompt:   "Press Ctrl+A.",
		ShowHelp: true,
	})

	for _, want := range []string{
		"Press Ctrl+A.",
		"Waiting for key...",
		"Help: press the requested shortcut.",
		"h hide help · Esc quit",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected help question render to contain %q, got:\n%s", want, got)
		}
	}
}

func TestReadKeySequenceUsesRawReaderForTerminalInput(t *testing.T) {
	restoreTerminal := stubTerminalDetector(func(fd uintptr) bool {
		return fd == os.Stdin.Fd()
	})
	defer restoreTerminal()
	called := false
	restoreRaw := stubRawKeyReader(func(file *os.File) (string, error) {
		called = true
		if file != os.Stdin {
			t.Fatalf("raw reader received %v, want os.Stdin", file)
		}
		return "Ctrl+A", nil
	})
	defer restoreRaw()

	got, err := readKeySequence(Options{In: os.Stdin}, bufio.NewReader(strings.NewReader("x")))
	if err != nil {
		t.Fatalf("readKeySequence returned error: %v", err)
	}
	if got != "Ctrl+A" {
		t.Fatalf("readKeySequence = %q, want Ctrl+A", got)
	}
	if !called {
		t.Fatal("expected terminal input to use raw key reader")
	}
}

func TestRenderKeySequenceFeedbackCorrect(t *testing.T) {
	got := renderKeySequenceFeedback(keySequenceFeedbackView{
		Received:    "Ctrl+A",
		Correct:     true,
		Explanation: "Ctrl+A moves to line start.",
	})

	for _, want := range []string{
		"Recu: Ctrl+A",
		"Correct.",
		"Ctrl+A moves to line start.",
		"Enter next · Esc quit",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected correct feedback to contain %q, got:\n%s", want, got)
		}
	}
	for _, unwanted := range []string{"r retry", "s solution"} {
		if strings.Contains(got, unwanted) {
			t.Fatalf("correct feedback should not contain %q, got:\n%s", unwanted, got)
		}
	}
}

func TestRenderKeySequenceFeedbackIncorrect(t *testing.T) {
	got := renderKeySequenceFeedback(keySequenceFeedbackView{
		Received: "Ctrl+L",
		Correct:  false,
	})

	for _, want := range []string{
		"Recu: Ctrl+L",
		"Pas encore.",
		"Enter next · r retry · s solution · Esc quit",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected incorrect feedback to contain %q, got:\n%s", want, got)
		}
	}
}

func TestRenderKeySequenceSolution(t *testing.T) {
	got := renderKeySequenceSolution(keySequenceSolutionView{
		Expected:    "Ctrl+W",
		Explanation: "Ctrl+W deletes the previous word.",
	})

	for _, want := range []string{
		"Solution: Ctrl+W",
		"Ctrl+W deletes the previous word.",
		"Enter next · r retry · Esc quit",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected solution render to contain %q, got:\n%s", want, got)
		}
	}
	if strings.Contains(got, "s solution") {
		t.Fatalf("solution render should not contain solution command again, got:\n%s", got)
	}
}

func TestRenderKeySequenceSummary(t *testing.T) {
	got := renderKeySequenceSummaryView(keySequenceSummaryView{
		Correct: 2,
		Total:   3,
		Missed:  1,
	})

	for _, want := range []string{
		"Chapitre termine.",
		"Correct: 2/3",
		"À revoir: 1",
		"Enter review missed · Esc quit",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected summary render to contain %q, got:\n%s", want, got)
		}
	}

	got = renderKeySequenceSummaryView(keySequenceSummaryView{
		Correct: 3,
		Total:   3,
		Missed:  0,
	})
	if strings.Contains(got, "Enter review missed") {
		t.Fatalf("perfect summary should not contain review footer, got:\n%s", got)
	}
}

func TestTrainKeySequenceSupportsRetryCommand(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	progressPath := filepath.Join(t.TempDir(), "progress.json")
	chapters := fstest.MapFS{
		"shortcuts.yaml": &fstest.MapFile{Data: []byte(`id: shortcuts
title: Shortcuts
items:
  - id: ctrl-l
    type: shortcut
    exercise_type: key-sequence
    prompt: Press Ctrl+L.
    answer:
      primary: Ctrl+L
`)},
	}

	var out bytes.Buffer
	cmd := NewRootCommand(chapters)
	cmd.SetIn(strings.NewReader("Lr\x0c\n"))
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"--progress", progressPath, "train", "shortcuts"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v\n%s", err, out.String())
	}

	got := out.String()
	if strings.Count(got, "Progression: 1/1\n\nPress Ctrl+L.") != 2 {
		t.Fatalf("expected retry to render the same exercise twice, got:\n%s", got)
	}
	if !strings.Contains(got, "Recu: L") || !strings.Contains(got, "Pas encore.") {
		t.Fatalf("expected first attempt to be wrong, got:\n%s", got)
	}
	if !strings.Contains(got, "Recu: Ctrl+L") || !strings.Contains(got, "Correct.") {
		t.Fatalf("expected retry attempt to be correct, got:\n%s", got)
	}
}

func TestTrainKeySequenceShowsHelpWithHBeforeAnswer(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	progressPath := filepath.Join(t.TempDir(), "progress.json")
	chapters := fstest.MapFS{
		"shortcuts.yaml": &fstest.MapFile{Data: []byte(`id: shortcuts
title: Shortcuts
items:
  - id: ctrl-a
    type: shortcut
    exercise_type: key-sequence
    prompt: Press Ctrl+A.
    answer:
      primary: Ctrl+A
`)},
	}

	var out bytes.Buffer
	cmd := NewRootCommand(chapters)
	cmd.SetIn(strings.NewReader("h\x01\n"))
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"--progress", progressPath, "train", "shortcuts"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v\n%s", err, out.String())
	}

	got := out.String()
	if !strings.Contains(got, "Help: press the requested shortcut.") {
		t.Fatalf("expected help text, got:\n%s", got)
	}
	if !strings.Contains(got, "h hide help · Esc quit") {
		t.Fatalf("expected help footer to advertise h hide help, got:\n%s", got)
	}
	if !strings.Contains(got, "Recu: Ctrl+A") || !strings.Contains(got, "Correct.") {
		t.Fatalf("expected answer after help to be accepted, got:\n%s", got)
	}
}

func TestTrainKeySequenceHelpToggleHidesHelp(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	progressPath := filepath.Join(t.TempDir(), "progress.json")
	chapters := fstest.MapFS{
		"shortcuts.yaml": &fstest.MapFile{Data: []byte(`id: shortcuts
title: Shortcuts
items:
  - id: ctrl-a
    type: shortcut
    exercise_type: key-sequence
    prompt: Press Ctrl+A.
    answer:
      primary: Ctrl+A
`)},
	}

	var out bytes.Buffer
	cmd := NewRootCommand(chapters)
	cmd.SetIn(strings.NewReader("hh\x01\n"))
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"--progress", progressPath, "train", "shortcuts"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v\n%s", err, out.String())
	}

	got := out.String()
	if strings.Count(got, "Help: press the requested shortcut.") != 1 {
		t.Fatalf("expected help to be shown once before being hidden, got:\n%s", got)
	}
	if !strings.Contains(got, "h hide help · Esc quit") || !strings.Contains(got, "h help · Esc quit") {
		t.Fatalf("expected help toggle footers, got:\n%s", got)
	}
	if !strings.Contains(got, "Recu: Ctrl+A") || !strings.Contains(got, "Correct.") {
		t.Fatalf("expected answer after help toggle to be accepted, got:\n%s", got)
	}
}

func TestTrainKeySequenceHelpEnterIsIgnored(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	progressPath := filepath.Join(t.TempDir(), "progress.json")
	chapters := fstest.MapFS{
		"shortcuts.yaml": &fstest.MapFile{Data: []byte(`id: shortcuts
title: Shortcuts
items:
  - id: ctrl-a
    type: shortcut
    exercise_type: key-sequence
    prompt: Press Ctrl+A.
    answer:
      primary: Ctrl+A
`)},
	}

	var out bytes.Buffer
	cmd := NewRootCommand(chapters)
	cmd.SetIn(strings.NewReader("h\n\x01\n"))
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"--progress", progressPath, "train", "shortcuts"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v\n%s", err, out.String())
	}

	got := out.String()
	if !strings.Contains(got, "Help: press the requested shortcut.") {
		t.Fatalf("expected help text, got:\n%s", got)
	}
	if strings.Contains(got, "Recu: Enter") || strings.Contains(got, "Pas encore.") {
		t.Fatalf("enter while capture help is visible should be ignored, got:\n%s", got)
	}
	if !strings.Contains(got, "Recu: Ctrl+A") || !strings.Contains(got, "Correct.") {
		t.Fatalf("expected answer after ignored enter to be accepted, got:\n%s", got)
	}
	if !strings.Contains(got, "Correct: 1/1") || !strings.Contains(got, "À revoir: 0") {
		t.Fatalf("help and ignored enter should not add review items, got:\n%s", got)
	}
}

func TestTrainKeySequenceHelpDoesNotRecordHAsAnswer(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	progressPath := filepath.Join(t.TempDir(), "progress.json")
	chapters := fstest.MapFS{
		"shortcuts.yaml": &fstest.MapFile{Data: []byte(`id: shortcuts
title: Shortcuts
items:
  - id: ctrl-a
    type: shortcut
    exercise_type: key-sequence
    prompt: Press Ctrl+A.
    answer:
      primary: Ctrl+A
`)},
	}

	var out bytes.Buffer
	cmd := NewRootCommand(chapters)
	cmd.SetIn(strings.NewReader("h\x01\n"))
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"--progress", progressPath, "train", "shortcuts"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v\n%s", err, out.String())
	}

	got := out.String()
	if strings.Contains(got, "Recu: h") || strings.Contains(got, "Pas encore.") {
		t.Fatalf("h help should not be recorded as a wrong answer, got:\n%s", got)
	}
	if !strings.Contains(got, "Correct: 1/1") || !strings.Contains(got, "À revoir: 0") {
		t.Fatalf("h help should not add the exercise to review, got:\n%s", got)
	}
}

func TestTrainKeySequenceHelpStateResetsBetweenQuestions(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	progressPath := filepath.Join(t.TempDir(), "progress.json")
	chapters := fstest.MapFS{
		"shortcuts.yaml": &fstest.MapFile{Data: []byte(`id: shortcuts
title: Shortcuts
items:
  - id: ctrl-a
    type: shortcut
    exercise_type: key-sequence
    prompt: Press Ctrl+A.
    answer:
      primary: Ctrl+A
  - id: ctrl-e
    type: shortcut
    exercise_type: key-sequence
    prompt: Press Ctrl+E.
    answer:
      primary: Ctrl+E
`)},
	}

	var out bytes.Buffer
	cmd := NewRootCommand(chapters)
	cmd.SetIn(strings.NewReader("h\x01\n\x05\n"))
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"--progress", progressPath, "train", "shortcuts"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v\n%s", err, out.String())
	}

	got := out.String()
	secondQuestion := strings.LastIndex(got, "Progression: 2/2")
	if secondQuestion == -1 {
		t.Fatalf("expected second question render, got:\n%s", got)
	}
	afterSecondQuestion := got[secondQuestion:]
	if strings.Contains(afterSecondQuestion, "h hide help · Esc quit") {
		t.Fatalf("help state should reset before the next question, got:\n%s", afterSecondQuestion)
	}
	if !strings.Contains(afterSecondQuestion, "h help · Esc quit") {
		t.Fatalf("expected default help footer on next question, got:\n%s", afterSecondQuestion)
	}
	if strings.Contains(afterSecondQuestion, "Help: press the requested shortcut.") {
		t.Fatalf("help text should not carry into next question, got:\n%s", afterSecondQuestion)
	}
}

func TestTrainKeySequenceDoesNotRetryAfterCorrectAnswer(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	progressPath := filepath.Join(t.TempDir(), "progress.json")
	chapters := fstest.MapFS{
		"shortcuts.yaml": &fstest.MapFile{Data: []byte(`id: shortcuts
title: Shortcuts
items:
  - id: ctrl-a
    type: shortcut
    exercise_type: key-sequence
    prompt: Press Ctrl+A.
    answer:
      primary: Ctrl+A
`)},
	}

	var out bytes.Buffer
	cmd := NewRootCommand(chapters)
	cmd.SetIn(strings.NewReader("\x01r\n"))
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"--progress", progressPath, "train", "shortcuts"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v\n%s", err, out.String())
	}

	got := out.String()
	if strings.Count(got, "Progression: 1/1\n\nPress Ctrl+A.") != 1 {
		t.Fatalf("correct retry command should not re-render the exercise, got:\n%s", got)
	}
	if strings.Contains(got, "Recu: r") || strings.Contains(got, "r retry") {
		t.Fatalf("correct retry command should only keep the next/quit footer, got:\n%s", got)
	}
}

func TestTrainKeySequenceDoesNotShowSolutionAfterCorrectAnswer(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	progressPath := filepath.Join(t.TempDir(), "progress.json")
	chapters := fstest.MapFS{
		"shortcuts.yaml": &fstest.MapFile{Data: []byte(`id: shortcuts
title: Shortcuts
items:
  - id: ctrl-a
    type: shortcut
    exercise_type: key-sequence
    prompt: Press Ctrl+A.
    answer:
      primary: Ctrl+A
`)},
	}

	var out bytes.Buffer
	cmd := NewRootCommand(chapters)
	cmd.SetIn(strings.NewReader("\x01s\n"))
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"--progress", progressPath, "train", "shortcuts"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v\n%s", err, out.String())
	}

	got := out.String()
	if strings.Contains(got, "Solution:") || strings.Contains(got, "Recu: s") || strings.Contains(got, "s solution") {
		t.Fatalf("correct solution command should only keep the next/quit footer, got:\n%s", got)
	}
}

func TestTrainKeySequenceSolutionCommandShowsExpectedAnswer(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	progressPath := filepath.Join(t.TempDir(), "progress.json")
	chapters := fstest.MapFS{
		"shortcuts.yaml": &fstest.MapFile{Data: []byte(`id: shortcuts
title: Shortcuts
items:
  - id: ctrl-l
    type: shortcut
    exercise_type: key-sequence
    prompt: Press Ctrl+L.
    answer:
      primary: Ctrl+L
    explanation: Ctrl+L clears the screen.
`)},
	}

	var out bytes.Buffer
	cmd := NewRootCommand(chapters)
	cmd.SetIn(strings.NewReader("Ls\n"))
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"--progress", progressPath, "train", "shortcuts"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v\n%s", err, out.String())
	}

	got := out.String()
	if !strings.Contains(got, "Solution: Ctrl+L") {
		t.Fatalf("expected solution command to print expected answer, got:\n%s", got)
	}
	if !strings.Contains(got, "Ctrl+L clears the screen.") {
		t.Fatalf("expected solution command to print explanation, got:\n%s", got)
	}
}

func TestTrainKeySequenceReviewsMissedExercises(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	progressPath := filepath.Join(t.TempDir(), "progress.json")
	chapters := fstest.MapFS{
		"shortcuts.yaml": &fstest.MapFile{Data: []byte(`id: shortcuts
title: Shortcuts
items:
  - id: ctrl-a
    type: shortcut
    exercise_type: key-sequence
    prompt: Press Ctrl+A.
    answer:
      primary: Ctrl+A
  - id: ctrl-e
    type: shortcut
    exercise_type: key-sequence
    prompt: Press Ctrl+E.
    answer:
      primary: Ctrl+E
`)},
	}

	var out bytes.Buffer
	cmd := NewRootCommand(chapters)
	cmd.SetIn(strings.NewReader("L\n\x05\n\n\x01\n"))
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"--progress", progressPath, "train", "shortcuts"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v\n%s", err, out.String())
	}

	got := out.String()
	if !strings.Contains(got, "Correct: 1/2") || !strings.Contains(got, "À revoir: 1") {
		t.Fatalf("expected initial summary with one missed exercise, got:\n%s", got)
	}
	if strings.Count(got, "Progression: 1/1\n\nPress Ctrl+A.") != 1 {
		t.Fatalf("expected review pass to include only missed Ctrl+A exercise, got:\n%s", got)
	}
	if !strings.Contains(got, "Correct: 1/1") || !strings.Contains(got, "À revoir: 0") {
		t.Fatalf("expected review summary to clear missed exercise, got:\n%s", got)
	}
}

func TestTrainKeySequenceExitsCleanlyOnEscapeAndCtrlC(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "escape", input: "\x1b\x01"},
		{name: "ctrl c", input: "\x03\x01"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("XDG_CONFIG_HOME", t.TempDir())
			progressPath := filepath.Join(t.TempDir(), "progress.json")
			chapters := fstest.MapFS{
				"shortcuts.yaml": &fstest.MapFile{Data: []byte(`id: shortcuts
title: Shortcuts
items:
  - id: ctrl-a
    type: shortcut
    exercise_type: key-sequence
    prompt: Press Ctrl+A.
    answer:
      primary: Ctrl+A
  - id: ctrl-e
    type: shortcut
    exercise_type: key-sequence
    prompt: Press Ctrl+E.
    answer:
      primary: Ctrl+E
`)},
			}

			var out bytes.Buffer
			cmd := NewRootCommand(chapters)
			cmd.SetIn(strings.NewReader(tt.input))
			cmd.SetOut(&out)
			cmd.SetErr(&out)
			cmd.SetArgs([]string{"--progress", progressPath, "train", "shortcuts"})

			if err := cmd.Execute(); err != nil {
				t.Fatalf("Execute returned error: %v\n%s", err, out.String())
			}

			got := out.String()
			if !strings.Contains(got, "Session interrompue.") {
				t.Fatalf("expected clean interruption message, got:\n%s", got)
			}
			if strings.Contains(got, "2/2") || strings.Contains(got, "Correct.") {
				t.Fatalf("expected session to stop before next exercise, got:\n%s", got)
			}
		})
	}
}

func writeAppTestConfigWithFunctions(t *testing.T, root string) string {
	t.Helper()
	path := filepath.Join(root, "config.toml")
	content := `dotfiles_path = "` + root + `"
shell = "zsh"

[paths]
aliases = ["zsh/modules/aliases.zsh"]
functions = ["zsh/modules/functions.zsh"]
docs = []

[security]
exclude = []
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile config returned error: %v", err)
	}
	return path
}

func TestInitPrintWritesDefaultConfigToStdoutOnly(t *testing.T) {
	out, err := executeTestCommand("init", "--print")
	if err != nil {
		t.Fatalf("Execute returned error: %v\n%s", err, out)
	}
	if !strings.Contains(out, `dotfiles_path = "~/dotfiles"`) || !strings.Contains(out, `shell = "zsh"`) {
		t.Fatalf("expected default TOML config, got:\n%s", out)
	}
	if strings.Contains(out, "Created config:") {
		t.Fatalf("init --print should not create a file, got:\n%s", out)
	}
}

func executeTestCommand(args ...string) (string, error) {
	var out bytes.Buffer
	cmd := NewRootCommand(data.Chapters())
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return out.String(), err
}

func writeAppTestConfig(t *testing.T, root string) string {
	t.Helper()
	path := filepath.Join(root, "config.toml")
	content := `dotfiles_path = "` + root + `"
shell = "zsh"

[paths]
aliases = ["zsh/modules/aliases.zsh"]
functions = []
docs = []

[security]
exclude = []
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile config returned error: %v", err)
	}
	return path
}

func writeAppTestFile(t *testing.T, root, name, content string) {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
}

func stubTerminalDetector(fn func(uintptr) bool) func() {
	previous := isTerminalFD
	isTerminalFD = fn
	return func() {
		isTerminalFD = previous
	}
}

func stubRawKeyReader(fn func(*os.File) (string, error)) func() {
	previous := readRawKeyFromTerminal
	readRawKeyFromTerminal = fn
	return func() {
		readRawKeyFromTerminal = previous
	}
}

func stubTUIRunner(fn func(tui.Options) (tui.Result, error)) func() {
	previous := runTUIWithOptions
	runTUIWithOptions = fn
	return func() {
		runTUIWithOptions = previous
	}
}
