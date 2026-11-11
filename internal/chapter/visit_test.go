package chapter

import (
	"fmt"
	"slices"
	"testing"
)

func TestVisitBatchesKeepFilterOrderAndAcquisitionPolicy(t *testing.T) {
	ch := Chapter{ID: "practice", LearningOrder: LearningOrderGuided}
	for i := 1; i <= 13; i++ {
		ch.Items = append(ch.Items, Item{ID: fmt.Sprint(i), Difficulty: DifficultyBeginner})
	}
	ch.Items = append(ch.Items,
		Item{ID: "higher-level", Difficulty: DifficultyAdvanced},
		Item{ID: "requires-tools", Difficulty: DifficultyBeginner, Tools: []string{"fd", "fzf"}},
	)
	filter := SessionFilter{Difficulty: DifficultyBeginner, AvailableTools: []string{"fd", "fzf"}, EnabledTools: []string{"fd"}}
	acquired := []string{"practice/2", "retired/3"}
	allEligible := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13"}
	learning := slices.Delete(slices.Clone(allEligible), 1, 2)
	for _, tc := range []struct {
		name  string
		mode  SessionMode
		order LearningOrder
		all   bool
		want  []string
		sizes []int
	}{
		{"guided", ModeLearn, LearningOrderGuided, false, learning, []int{6, 6}},
		{"independent", ModeLearn, LearningOrderIndependent, false, learning, []int{6, 6}},
		{"review", ModeReview, LearningOrderGuided, false, allEligible, []int{6, 6, 1}},
		{"full", ModeLearn, LearningOrderGuided, true, allEligible, []int{13}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ch.LearningOrder = tc.order
			visit := NewVisit(filter.PrepareVisit(ch, TextExercises, tc.mode, tc.all, acquired), tc.all)
			var got []string
			for _, size := range tc.sizes {
				batch := visit.Next()
				if len(batch.Items) != size {
					t.Fatalf("batch size %d, want %d", len(batch.Items), size)
				}
				got = append(got, sessionIDs(batch.Items)...)
			}
			if visit.HasNext() {
				t.Fatal("visit still has exercises")
			}
			if tc.mode == ModeLearn && tc.order == LearningOrderGuided {
				if !slices.Equal(got, tc.want) {
					t.Fatalf("wrong editorial order: %v", got)
				}
			} else {
				want := slices.Clone(tc.want)
				slices.Sort(got)
				slices.Sort(want)
				if !slices.Equal(got, want) {
					t.Fatalf("wrong coverage or duplicates: %v", got)
				}
			}
		})
	}
}
