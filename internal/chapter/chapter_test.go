package chapter

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/benabot/cli-drill/internal/exercise"
)

func TestLoadHistoricalCorpus(t *testing.T) {
	chapters, err := LoadFS(os.DirFS("../../testdata/legacy-chapters"), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	want := []struct {
		id    string
		items int
	}{
		{"01-raccourcis-terminal", 8},
		{"02-navigation-shell", 6},
		{"03-alias-zsh", 7},
		{"04-fonctions-zsh", 7},
		{"05-outils-quotidiens", 7},
		{"06-recherche-fichiers-contenu", 7},
		{"07-lecture-preview", 6},
		{"08-micro", 6},
		{"09-markdown-glow", 7},
		{"10-workflows-dotfiles", 7},
	}
	if len(chapters) != len(want) {
		t.Fatalf("got %d chapters, want %d", len(chapters), len(want))
	}
	for i, expected := range want {
		if chapters[i].ID != expected.id || len(chapters[i].Items) != expected.items {
			t.Fatalf("chapter %d: got %s (%d items), want %s (%d items)",
				i, chapters[i].ID, len(chapters[i].Items), expected.id, expected.items)
		}
	}
	if chapters[0].Items[0].ExerciseType != exercise.TypeKeySequence {
		t.Fatal("keyboard chapter lost real key capture")
	}
	navigation := chapters[1].Items
	if navigation[0].ID != "zoxide-jump" || navigation[0].Type != "tool" ||
		navigation[0].Answer.Primary != "zoxide" ||
		!reflect.DeepEqual(navigation[0].Answer.Accepted, []string{"z"}) ||
		navigation[0].Explanation == "" {
		t.Fatalf("pedagogical metadata lost: %#v", navigation[0])
	}
	if navigation[2].ExerciseType != exercise.TypeMultipleChoice ||
		!reflect.DeepEqual(navigation[2].Choices, []string{"fd", "bat", "glow", "vale"}) {
		t.Fatalf("multiple-choice content lost: %#v", navigation[2])
	}
	if navigation[5].ExerciseType != exercise.TypeSimpleShellSim ||
		navigation[5].Answer.Primary != "fd -t d | fzf" {
		t.Fatalf("non-executing simulation content lost: %#v", navigation[5])
	}
}

func TestLoadChapterFromYAML(t *testing.T) {
	input := strings.NewReader(`
id: 01-raccourcis-terminal
title: Raccourcis terminal
description: Memoriser les raccourcis.
items:
  - id: ctrl-a
    type: shortcut
    exercise_type: free-answer
    prompt: Aller au debut de la ligne
    answer:
      primary: Ctrl+A
      accepted:
        - ^A
    explanation: Ctrl+A place le curseur au debut.
`)

	chapter, err := Load(input)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if chapter.ID != "01-raccourcis-terminal" || chapter.Title != "Raccourcis terminal" {
		t.Fatalf("unexpected chapter metadata: %#v", chapter)
	}
	if len(chapter.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(chapter.Items))
	}
	if chapter.Items[0].Answer.Primary != "Ctrl+A" {
		t.Fatalf("unexpected answer: %#v", chapter.Items[0].Answer)
	}
}

func TestLoadKeySequenceChapterItemFromYAML(t *testing.T) {
	input := strings.NewReader(`
id: shortcuts
title: Shortcuts
items:
  - id: ctrl-l
    type: shortcut
    exercise_type: key-sequence
    prompt: Appuie sur le raccourci pour nettoyer l'ecran.
    answer:
      primary: Ctrl+L
`)

	chapter, err := Load(input)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if chapter.Items[0].ExerciseType != exercise.TypeKeySequence {
		t.Fatalf("expected key-sequence exercise type, got %q", chapter.Items[0].ExerciseType)
	}
}

func TestValidateRejectsMissingItemAnswer(t *testing.T) {
	chapter := Chapter{
		ID:          "bad",
		Title:       "Bad",
		Description: "Bad chapter",
		Items: []Item{{
			ID:           "missing-answer",
			Type:         "alias",
			ExerciseType: "free-answer",
			Prompt:       "What alias?",
		}},
	}

	if err := chapter.Validate(); err == nil {
		t.Fatal("expected validation error")
	}
}
