package chapter

import (
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/benabot/cli-drill/internal/exercise"
)

func TestSessionFilters(t *testing.T) {
	ch := Chapter{Items: []Item{
		{ID: "unix"},
		{ID: "fd", Difficulty: DifficultyBeginner, Tools: []string{"fd"}},
		{ID: "fzf", Difficulty: DifficultyBeginner, Tools: []string{"fzf"}},
		{ID: "both", Difficulty: DifficultyIntermediate, Tools: []string{"fd", "fzf"}},
		{ID: "intermediate-key", Difficulty: DifficultyIntermediate, ExerciseType: exercise.TypeKeySequence},
		{ID: "advanced", Difficulty: DifficultyAdvanced},
	}}
	for _, tc := range []struct {
		name    string
		enabled []string
		level   Difficulty
		scope   ExerciseScope
		want    []string
	}{
		{"defaults", nil, DifficultyBeginner, AllExercises, []string{"unix"}},
		{"retired choice", []string{"retired"}, DifficultyBeginner, AllExercises, []string{"unix"}},
		{"one tool", []string{"fd"}, DifficultyBeginner, AllExercises, []string{"unix", "fd"}},
		{"both tools beginner", []string{"fd", "fzf"}, DifficultyBeginner, AllExercises, []string{"unix", "fd", "fzf"}},
		{"missing second tool", []string{"fd"}, DifficultyIntermediate, AllExercises, []string{"intermediate-key"}},
		{"both tools intermediate", []string{"fd", "fzf"}, DifficultyIntermediate, AllExercises, []string{"both", "intermediate-key"}},
		{"exact advanced", []string{"fd", "fzf"}, DifficultyAdvanced, AllExercises, []string{"advanced"}},
		{"only capture", []string{"fd", "fzf"}, DifficultyIntermediate, KeyExercises, []string{"intermediate-key"}},
		{"only text", []string{"fd", "fzf"}, DifficultyIntermediate, TextExercises, []string{"both"}},
		{"empty", nil, DifficultyAdvanced, KeyExercises, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			filter := SessionFilter{AvailableTools: AvailableTools([]Chapter{ch}), EnabledTools: tc.enabled, Difficulty: tc.level}
			if got := sessionIDs(filter.EligibleItems(ch, tc.scope)); !slices.Equal(got, tc.want) {
				t.Fatalf("eligible = %v, want %v", got, tc.want)
			}
		})
	}
	filter := SessionFilter{
		AvailableTools: []string{"fd"},
		EnabledTools:   []string{"fd", "fzf"}, // fzf is saved but absent from the active corpus.
		Difficulty:     DifficultyIntermediate,
	}
	if got := sessionIDs(filter.EligibleItems(ch, AllExercises)); !slices.Equal(got, []string{"intermediate-key"}) {
		t.Fatalf("a saved tool outside the corpus enabled an item: %v", got)
	}
}

func TestSessionOrder(t *testing.T) {
	filter := SessionFilter{Difficulty: DifficultyBeginner}
	ch := Chapter{ID: "lesson", LearningOrder: LearningOrderGuided, Items: []Item{
		{ID: "a"}, {ID: "hidden", Difficulty: DifficultyAdvanced}, {ID: "b"}, {ID: "c"},
	}}
	want := []string{"a", "b", "c"}
	reverse := func(n int, swap func(int, int)) {
		for i := 0; i < n/2; i++ {
			swap(i, n-1-i)
		}
	}
	learn := filter.Prepare(ch, AllExercises, ModeLearn, func(int, func(int, int)) {
		t.Fatal("guided learn shuffled")
	})
	if learn.ID != ch.ID || !slices.Equal(sessionIDs(learn.Items), want) {
		t.Fatalf("learn = %#v", learn)
	}
	for range 2 { // Repeated permutations are valid.
		review := filter.Prepare(ch, AllExercises, ModeReview, reverse)
		if !slices.Equal(sessionIDs(review.Items), []string{"c", "b", "a"}) {
			t.Fatalf("review = %v", sessionIDs(review.Items))
		}
	}
	ch.LearningOrder = LearningOrderIndependent
	for seed := uint64(0); seed < 20; seed++ {
		shuffle := rand.New(rand.NewPCG(seed, 0)).Shuffle
		ids := sessionIDs(filter.Prepare(ch, AllExercises, "", shuffle).Items)
		slices.Sort(ids)
		if !slices.Equal(ids, want) {
			t.Fatalf("seed %d produced %v, want a permutation of %v", seed, ids, want)
		}
	}
	if !slices.Equal(sessionIDs(ch.Items), []string{"a", "hidden", "b", "c"}) {
		t.Fatal("session reordered source")
	}
	for size := range 2 {
		ch.Items = ch.Items[:size]
		got := filter.Prepare(ch, AllExercises, "", func(int, func(int, int)) {
			t.Fatal("shuffled a list with fewer than two items")
		})
		if !slices.Equal(sessionIDs(got.Items), sessionIDs(ch.Items)) {
			t.Fatalf("size %d: got %v", size, sessionIDs(got.Items))
		}
	}
}

func sessionIDs(items []Item) []string {
	var ids []string
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	return ids
}
