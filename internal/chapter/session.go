package chapter

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"strings"

	"github.com/benabot/cli-drill/internal/exercise"
)

type SessionMode string

const (
	ModeLearn  SessionMode = "learn"
	ModeReview SessionMode = "review"
)

type ExerciseScope int

const (
	AllExercises ExerciseScope = iota
	TextExercises
	KeyExercises
)

type SessionFilter struct {
	AvailableTools []string
	EnabledTools   []string
	Difficulty     Difficulty
}

func (c Chapter) Guided() bool {
	return c.LearningOrder != LearningOrderIndependent
}

func (c Chapter) HasKeyExercises() bool {
	return slices.ContainsFunc(c.Items, KeyExercises.includes)
}

func (scope ExerciseScope) includes(item Item) bool {
	isKey := item.ExerciseType == exercise.TypeKeySequence
	return scope == AllExercises || (scope == KeyExercises && isKey) || (scope == TextExercises && !isKey)
}

func (item Item) level() Difficulty {
	if item.Difficulty == "" {
		return DifficultyBeginner
	}
	return item.Difficulty
}

// EligibleItems keeps editorial order and owns the returned slice.
func (f SessionFilter) EligibleItems(c Chapter, scope ExerciseScope) []Item {
	var items []Item
	for _, item := range c.Items {
		if scope.includes(item) && item.EligibleForTools(f.AvailableTools, f.EnabledTools) && item.level() == f.Difficulty {
			items = append(items, item)
		}
	}
	return items
}

func (f SessionFilter) FirstTextChapter(chapters []Chapter) int {
	return slices.IndexFunc(chapters, func(c Chapter) bool {
		return len(f.EligibleItems(c, TextExercises)) > 0
	})
}

// Prepare snapshots the eligible exercises and fixes their order for the session.
func (f SessionFilter) Prepare(c Chapter, scope ExerciseScope, mode SessionMode, shuffle func(int, func(int, int))) Chapter {
	c.Items = f.EligibleItems(c, scope)
	return orderSession(c, mode, shuffle)
}

func orderSession(c Chapter, mode SessionMode, shuffle func(int, func(int, int))) Chapter {
	if len(c.Items) > 1 && (!c.Guided() || mode == ModeReview) {
		if shuffle == nil {
			shuffle = rand.Shuffle
		}
		shuffle(len(c.Items), func(i, j int) {
			c.Items[i], c.Items[j] = c.Items[j], c.Items[i]
		})
	}
	return c
}

// EmptyReason explains a selection already known to contain no eligible items.
func (f SessionFilter) EmptyReason(c Chapter, scope ExerciseScope) string {
	for _, item := range c.Items {
		if scope.includes(item) && item.level() == f.Difficulty {
			return fmt.Sprintf("No eligible exercises at %s difficulty with the enabled CLI modules.", f.Difficulty)
		}
	}
	return fmt.Sprintf("No exercises at %s difficulty in this chapter.", f.Difficulty)
}

// MissingModuleGroups lists the alternative sets of modules for exercises
// at the selected level. Every module in a group must be enabled together.
func (f SessionFilter) MissingModuleGroups(c Chapter, scope ExerciseScope) []string {
	var groups []string
	for _, item := range c.Items {
		if !scope.includes(item) || item.level() != f.Difficulty {
			continue
		}
		var missing []string
		for _, tool := range item.Tools {
			if !slices.Contains(f.EnabledTools, tool) {
				missing = append(missing, tool)
			}
		}
		if len(missing) != 0 {
			slices.Sort(missing)
			groups = append(groups, strings.Join(missing, " + "))
		}
	}
	slices.Sort(groups)
	return slices.Compact(groups)
}

func AcquiredCount(chapterID string, items []Item, completed []string) int {
	count := 0
	for _, item := range items {
		if slices.Contains(completed, chapterID+"/"+item.ID) {
			count++
		}
	}
	return count
}

func (f SessionFilter) ActiveCounts(chapters []Chapter, completed []string) (acquired, total int) {
	for _, c := range chapters {
		items := f.EligibleItems(c, AllExercises)
		acquired += AcquiredCount(c.ID, items, completed)
		total += len(items)
	}
	return acquired, total
}
