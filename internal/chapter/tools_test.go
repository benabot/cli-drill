package chapter

import (
	"slices"
	"testing"
)

func TestAvailableToolsAndEligibility(t *testing.T) {
	unix := Item{}
	one := Item{Tools: []string{"fd"}}
	multi := Item{Tools: []string{"fd", "fzf"}}
	corpus := []Chapter{{Items: []Item{{Tools: []string{"fzf", "fd"}}, unix}}, {Items: []Item{one}}}
	available := AvailableTools(corpus)
	if !slices.Equal(available, []string{"fd", "fzf"}) {
		t.Fatalf("available tools = %v", available)
	}
	for _, tc := range []struct {
		name    string
		item    Item
		enabled []string
		want    bool
	}{
		{"Unix without preferences", unix, nil, true},
		{"disabled", one, nil, false},
		{"enabled", one, []string{"fd"}, true},
		{"only one of two", multi, []string{"fd"}, false},
		{"both", multi, []string{"fd", "fzf"}, true},
		{"retired", Item{Tools: []string{"retired"}}, []string{"retired"}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.item.EligibleForTools(available, tc.enabled); got != tc.want {
				t.Errorf("eligible = %t, want %t", got, tc.want)
			}
		})
	}
}
