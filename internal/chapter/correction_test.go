package chapter

import (
	"testing"

	"github.com/benabot/cli-drill/data"
	"github.com/benabot/cli-drill/internal/exercise"
)

func TestEmbeddedCorrectionAcceptsEveryDeclaredAnswer(t *testing.T) {
	chapters, err := LoadFS(data.Chapters(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	for _, ch := range chapters {
		for _, item := range ch.Items {
			t.Run(ch.ID+"/"+item.ID, func(t *testing.T) {
				for _, answer := range append([]string{item.Answer.Primary}, item.Answer.Accepted...) {
					if answer == "" || !item.MatchAnswer(answer) {
						t.Fatalf("declared answer rejected: %q", answer)
					}
				}
				if item.MatchAnswer("") || item.MatchAnswer("definitely-wrong-answer") {
					t.Fatal("empty or unrelated answer accepted")
				}
				if item.ExerciseType == exercise.TypeMultipleChoice {
					for _, choice := range item.Choices {
						// All current QCMs have exactly one correct displayed choice.
						if got, want := item.MatchAnswer(choice), choice == item.Answer.Primary; got != want {
							t.Fatalf("choice %q: got %v, want %v", choice, got, want)
						}
					}
				}
			})
		}
	}
}

func TestItemCorrectionUsesPedagogicalContext(t *testing.T) {
	for _, kind := range []exercise.Type{exercise.TypeFreeAnswer, exercise.TypeScenario, exercise.TypeMultipleChoice, exercise.TypeSimpleShellSim} {
		t.Run(string(kind), func(t *testing.T) {
			item := Item{Type: "workflow", ExerciseType: kind, Answer: exercise.Answer{Primary: "grep ERROR file"}}
			if !item.MatchAnswer(" grep ERROR file \n") || item.MatchAnswer("grep error file") {
				t.Fatal("command case or edge-padding rule changed with exercise type")
			}
			item.Type = "shortcut"
			item.Answer.Primary = "Ctrl+L"
			if kind != exercise.TypeSimpleShellSim && (!item.MatchAnswer("^l") || item.MatchAnswer("l")) {
				t.Fatal("textual shortcut lost its dedicated comparison")
			}
		})
	}
	for _, topic := range []string{"alias", "function", "workflow"} {
		item := Item{Type: topic, ExerciseType: exercise.TypeFreeAnswer, Answer: exercise.Answer{Primary: "mc"}}
		if item.MatchAnswer("MC") || !item.MatchAnswer("mc") {
			t.Fatalf("identifier case lost for %s", topic)
		}
	}
	for _, topic := range []string{"tool", "concept"} {
		item := Item{Type: topic, ExerciseType: exercise.TypeScenario, Answer: exercise.Answer{Primary: "fd"}}
		if !item.MatchAnswer("FD") {
			t.Fatalf("bare pedagogical name lost its case tolerance for %s", topic)
		}
		item.ExerciseType = exercise.TypeSimpleShellSim
		if item.MatchAnswer("FD") {
			t.Fatal("shell command inherited a name's case tolerance")
		}
	}
	item := Item{ExerciseType: exercise.TypeKeySequence, Answer: exercise.Answer{Primary: "Ctrl+A"}}
	if !item.MatchAnswer("\x01") || item.MatchAnswer("A") {
		t.Fatal("key capture comparison changed")
	}
}
