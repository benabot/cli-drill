package chapter

import "slices"

const BatchSize = 6

// VisitItems selects the eligible items for the chosen format and mode.
// Only short learning visits exclude previously acquired exercises.
func (f SessionFilter) VisitItems(c Chapter, scope ExerciseScope, mode SessionMode, all bool, acquired []string) []Item {
	items := f.EligibleItems(c, scope)
	if !all && mode != ModeReview {
		items = slices.DeleteFunc(items, func(item Item) bool {
			return slices.Contains(acquired, c.ID+"/"+item.ID)
		})
	}
	return items
}

// PrepareVisit fixes the order for this visit.
func (f SessionFilter) PrepareVisit(c Chapter, scope ExerciseScope, mode SessionMode, all bool, acquired []string) Chapter {
	c.Items = f.VisitItems(c, scope, mode, all, acquired)
	return orderSession(c, mode, nil)
}

// Visit consumes the prepared order in batches. Reviewing mistakes does not
// change this cursor.
type Visit struct {
	remaining Chapter
	all       bool
}

func NewVisit(prepared Chapter, all bool) Visit {
	return Visit{remaining: prepared, all: all}
}

func (v Visit) HasNext() bool {
	return len(v.remaining.Items) != 0
}

func (v *Visit) Next() Chapter {
	batch := v.remaining
	if !v.all {
		batch.Items = batch.Items[:min(BatchSize, len(batch.Items))]
	}
	v.remaining.Items = v.remaining.Items[len(batch.Items):]
	return batch
}

type MissedAnswer struct {
	Item     Item
	Answered string
}

type BatchResult struct {
	Total   int
	Correct int
	Missed  []MissedAnswer
}

func (b *BatchResult) Record(item Item, answered string, correct bool) {
	b.Total++
	if correct {
		b.Correct++
	} else {
		b.Missed = append(b.Missed, MissedAnswer{Item: item, Answered: answered})
	}
}

func (b BatchResult) RetryItems() []Item {
	items := make([]Item, len(b.Missed))
	for i, missed := range b.Missed {
		items[i] = missed.Item
	}
	return items
}
