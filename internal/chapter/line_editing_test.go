package chapter

import (
	"slices"
	"testing"

	"github.com/benabot/cli-drill/data"
	"github.com/benabot/cli-drill/internal/exercise"
)

func TestWave3LineEditing(t *testing.T) {
	chapters, err := LoadFS(data.Chapters(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	index := slices.IndexFunc(chapters, func(ch Chapter) bool { return ch.ID == "01-raccourcis-terminal" })
	if index < 0 {
		t.Fatal("missing line editing chapter")
	}
	ch := chapters[index]
	if ch.Family != FamilyKeyboardShortcuts || ch.LearningOrder != LearningOrderGuided {
		t.Fatalf("unexpected family/order: %s / %s", ch.Family, ch.LearningOrder)
	}
	items := make(map[string]Item)
	for _, item := range ch.Items {
		items[item.ID] = item
		if len(item.Tools) != 0 {
			t.Errorf("%s: unexpected modules %v", item.ID, item.Tools)
		}
	}
	for _, tc := range []struct {
		id             string
		level          Difficulty
		format         exercise.Type
		primary, wrong string
		raw            byte
	}{
		{"ctrl-b-character-left", DifficultyBeginner, exercise.TypeKeySequence, "Ctrl+B", "Ctrl+A", 2},
		{"ctrl-f-character-right", DifficultyBeginner, exercise.TypeKeySequence, "Ctrl+F", "Ctrl+E", 6},
		{"ctrl-h-delete-before-cursor", DifficultyBeginner, exercise.TypeKeySequence, "Backspace", "Ctrl+D", 127},
		{"tab-complete-unique-name", DifficultyBeginner, exercise.TypeKeySequence, "Tab", "Ctrl+D", 9},
		{"ctrl-d-delete-under-cursor", DifficultyIntermediate, exercise.TypeKeySequence, "Ctrl+D", "Ctrl+H", 4},
		{"ctrl-t-transpose-characters", DifficultyIntermediate, exercise.TypeKeySequence, "Ctrl+T", "Ctrl+W", 20},
		{"ctrl-p-previous-history", DifficultyIntermediate, exercise.TypeKeySequence, "Ctrl+P", "Ctrl+R", 16},
		{"ctrl-n-next-history", DifficultyIntermediate, exercise.TypeKeySequence, "Ctrl+N", "Ctrl+P", 14},
		{"ctrl-underscore-undo-edit", DifficultyIntermediate, exercise.TypeKeySequence, "Ctrl+_", "Ctrl+Y", 31},
		{"ctrl-g-restore-history-draft", DifficultyIntermediate, exercise.TypeKeySequence, "Ctrl+G", "Ctrl+R", 7},
		{"readline-backward-word-edit", DifficultyIntermediate, exercise.TypeMultipleChoice, "git switch |sauvegarde", "git |switch sauvegarde", 0},
		{"readline-forward-word-edit", DifficultyIntermediate, exercise.TypeMultipleChoice, "git status| --short", "git status |--short", 0},
		{"readline-cut-word-forward", DifficultyIntermediate, exercise.TypeMultipleChoice, "grep | serveur.log", "grep |ERREUR", 0},
		{"bash-abandon-incomplete-input", DifficultyIntermediate, exercise.TypeMultipleChoice, "Interrompre la saisie avec Ctrl+C.", "Appuyer sur Entrée pour valider la ligne.", 0},
		{"readline-refine-ambiguous-completion", DifficultyIntermediate, exercise.TypeMultipleChoice, "Ajouter n au préfixe puis redemander la complétion.", "Relancer la même complétion sans préciser le nom.", 0},
		{"readline-edit-found-history", DifficultyAdvanced, exercise.TypeMultipleChoice, "Reprendre l’édition avec un déplacement du curseur.", "Valider le résultat avec Entrée puis le corriger.", 0},
		{"ctrl-d-eof-versus-delete", DifficultyAdvanced, exercise.TypeMultipleChoice, "A efface b. B signale une fin d’entrée.", "A et B signalent une fin d’entrée.", 0},
		{"readline-unexpected-vi-mode", DifficultyAdvanced, exercise.TypeMultipleChoice, "Readline était en mode commande Vi.", "L’émulateur avalait tous les caractères tapés.", 0},
		{"terminal-consumed-editing-key", DifficultyAdvanced, exercise.TypeMultipleChoice, "L’émulateur a consommé la touche.", "Readline a reçu la touche et ignoré sa liaison.", 0},
		{"tty-flow-control-resume-output", DifficultyAdvanced, exercise.TypeMultipleChoice, "Le contrôle de flux suspend la sortie.", "Bash a quitté le mode d’édition.", 0},
		{"ctrl-u-zle-readline-scope", DifficultyAdvanced, exercise.TypeMultipleChoice, "Readline laisse cd. ZLE vide la ligne.", "Readline vide la ligne. ZLE laisse cd.", 0},
		{"zle-wordchars-path-boundary", DifficultyAdvanced, exercise.TypeMultipleChoice, "Readline revient avant api. ZLE revient avant src.", "Readline revient avant src. ZLE revient avant api.", 0},
		{"zle-multiline-editing-boundary", DifficultyAdvanced, exercise.TypeMultipleChoice, "Au début de la deuxième ligne logique.", "Au début de toute la commande.", 0},
	} {
		t.Run(tc.id, func(t *testing.T) {
			item, ok := items[tc.id]
			if !ok {
				t.Fatal("exercise is missing from the embedded corpus")
			}
			if item.Difficulty != tc.level || item.ExerciseType != tc.format || item.Answer.Primary != tc.primary {
				t.Fatalf("unexpected exercise: %#v", item)
			}
			if !item.MatchAnswer(tc.primary) || item.MatchAnswer(tc.wrong) {
				t.Fatalf("want %q accepted and %q rejected", tc.primary, tc.wrong)
			}
			for _, choice := range item.Choices {
				if item.MatchAnswer(choice) != (choice == tc.primary) {
					t.Errorf("unexpected verdict for choice %q", choice)
				}
			}
			if tc.raw != 0 {
				input := string([]byte{tc.raw})
				if notation, ok := exercise.KeyByteNotation(tc.raw); ok {
					input = notation
				}
				if !item.MatchAnswer(input) {
					t.Errorf("decoded terminal byte %#02x rejected", tc.raw)
				}
				if item.MatchAnswer("Esc") || item.MatchAnswer("Ctrl+C") {
					t.Error("reserved control accepted")
				}
			}
		})
	}

	for _, tc := range []struct {
		level       Difficulty
		total, keys int
		first, last string
	}{
		{DifficultyBeginner, 12, 12, "ctrl-b-character-left", "ctrl-w"},
		{DifficultyIntermediate, 11, 6, "ctrl-d-delete-under-cursor", "readline-refine-ambiguous-completion"},
		{DifficultyAdvanced, 8, 0, "readline-edit-found-history", "zle-multiline-editing-boundary"},
	} {
		filter := SessionFilter{Difficulty: tc.level}
		all := filter.EligibleItems(ch, AllExercises)
		if len(all) != tc.total || len(filter.EligibleItems(ch, KeyExercises)) != tc.keys || len(filter.EligibleItems(ch, TextExercises)) != tc.total-tc.keys {
			t.Fatalf("wrong capture/text counts for %s", tc.level)
		}
		if all[0].ID != tc.first || all[len(all)-1].ID != tc.last {
			t.Fatalf("unexpected editorial order for %s", tc.level)
		}
		learn := filter.Prepare(ch, AllExercises, ModeLearn, func(int, func(int, int)) {
			t.Fatal("Learn shuffled a guided chapter")
		})
		if !slices.EqualFunc(learn.Items, all, func(a, b Item) bool { return a.ID == b.ID }) {
			t.Fatalf("Learn lost editorial order for %s", tc.level)
		}
		reversed := slices.Clone(all)
		slices.Reverse(reversed)
		review := filter.Prepare(ch, AllExercises, ModeReview, func(n int, swap func(int, int)) {
			for i := 0; i < n/2; i++ {
				swap(i, n-1-i)
			}
		})
		idsEqual := func(a, b Item) bool { return a.ID == b.ID }
		if !slices.EqualFunc(review.Items, reversed, idsEqual) || !slices.EqualFunc(filter.EligibleItems(ch, AllExercises), all, idsEqual) {
			t.Fatalf("Review did not shuffle a copy of the eligible items for %s", tc.level)
		}
	}
}
