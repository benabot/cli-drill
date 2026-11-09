package chapter

import "slices"

// AvailableTools returns the sorted, unique module identifiers in the corpus.
func AvailableTools(chapters []Chapter) []string {
	tools := make(map[string]struct{})
	for _, ch := range chapters {
		for _, item := range ch.Items {
			for _, tool := range item.Tools {
				tools[tool] = struct{}{}
			}
		}
	}
	available := make([]string, 0, len(tools))
	for tool := range tools {
		available = append(available, tool)
	}
	slices.Sort(available)
	return available
}

// EligibleForTools requires every tool to be both available and enabled.
// Items with no required tools are always eligible.
func (item Item) EligibleForTools(available, enabled []string) bool {
	for _, tool := range item.Tools {
		if !slices.Contains(available, tool) || !slices.Contains(enabled, tool) {
			return false
		}
	}
	return true
}

// UnavailableTools reports saved identifiers absent from the current corpus.
func UnavailableTools(available, enabled []string) []string {
	var unavailable []string
	for _, tool := range enabled {
		if !slices.Contains(available, tool) {
			unavailable = append(unavailable, tool)
		}
	}
	return unavailable
}
