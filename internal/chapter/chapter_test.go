package chapter

import (
	"strings"
	"testing"
)

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
