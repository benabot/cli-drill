package catalog

import (
	"sort"
	"strings"
	"unicode"
)

type EntryType string

const (
	EntryShortcut EntryType = "shortcut"
	EntryAlias    EntryType = "alias"
	EntryFunction EntryType = "function"
	EntryTool     EntryType = "tool"
	EntryWorkflow EntryType = "workflow"
	EntryConcept  EntryType = "concept"
	EntryBinding  EntryType = "binding"
	EntryChapter  EntryType = "chapter"
)

type Source struct {
	Path string `json:"path" yaml:"path,omitempty"`
	Line int    `json:"line,omitempty" yaml:"line,omitempty"`
}

type Entry struct {
	ID      string    `json:"id" yaml:"id"`
	Name    string    `json:"name" yaml:"name"`
	Type    EntryType `json:"type" yaml:"type"`
	Summary string    `json:"summary,omitempty" yaml:"summary,omitempty"`
	Command string    `json:"command,omitempty" yaml:"command,omitempty"`
	Source  Source    `json:"source,omitempty" yaml:"source,omitempty"`
	Sources []Source  `json:"sources,omitempty" yaml:"sources,omitempty"`
	Tags    []string  `json:"tags,omitempty" yaml:"tags,omitempty"`
}

type Catalog struct {
	entries []Entry
	index   map[string]int
}

func New(entries ...Entry) Catalog {
	c := Catalog{}
	for _, entry := range entries {
		c.Add(entry)
	}
	return c
}

func (c *Catalog) Add(entry Entry) {
	if c.index == nil {
		c.index = map[string]int{}
		for i, existing := range c.entries {
			c.index[dedupeKey(existing)] = i
		}
	}
	if entry.ID == "" {
		entry.ID = NormalizeID(entry.Name)
	}
	entry = normalizeEntry(entry)
	key := dedupeKey(entry)
	if existingIndex, ok := c.index[key]; ok {
		c.entries[existingIndex] = mergeEntry(c.entries[existingIndex], entry)
		return
	}
	c.index[key] = len(c.entries)
	c.entries = append(c.entries, entry)
}

func (c Catalog) Entries() []Entry {
	entries := append([]Entry(nil), c.entries...)
	Sort(entries)
	return entries
}

func (c Catalog) FilterByType(entryType EntryType) []Entry {
	var matches []Entry
	for _, entry := range c.entries {
		if entry.Type == entryType {
			matches = append(matches, entry)
		}
	}
	Sort(matches)
	return matches
}

func (c Catalog) Search(query string) []Entry {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return c.Entries()
	}

	var matches []Entry
	for _, entry := range c.entries {
		haystack := strings.ToLower(strings.Join([]string{
			entry.ID,
			entry.Name,
			string(entry.Type),
			entry.Summary,
			entry.Command,
			strings.Join(entry.Tags, " "),
		}, " "))
		if strings.Contains(haystack, query) {
			matches = append(matches, entry)
		}
	}
	Sort(matches)
	return matches
}

func (c Catalog) Find(idOrName string) (Entry, bool) {
	needle := strings.ToLower(strings.TrimSpace(idOrName))
	for _, entry := range c.entries {
		if strings.ToLower(entry.ID) == needle || strings.ToLower(entry.Name) == needle {
			return entry, true
		}
	}
	return Entry{}, false
}

func Sort(entries []Entry) {
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Type != entries[j].Type {
			return entries[i].Type < entries[j].Type
		}
		return entries[i].Name < entries[j].Name
	})
}

func NormalizeID(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var out strings.Builder
	lastDash := false
	for _, r := range value {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			out.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			out.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(out.String(), "-")
}

func dedupeKey(entry Entry) string {
	return string(entry.Type) + "\x00" + entry.ID
}

func normalizeEntry(entry Entry) Entry {
	if entry.Source.Path != "" || entry.Source.Line != 0 {
		entry.Sources = appendSource(entry.Sources, entry.Source)
	}
	if entry.Source.Path == "" && len(entry.Sources) > 0 {
		entry.Source = entry.Sources[0]
	}
	entry.Tags = uniqueStrings(entry.Tags)
	return entry
}

func mergeEntry(existing, incoming Entry) Entry {
	incoming = normalizeEntry(incoming)
	if existing.Name == "" {
		existing.Name = incoming.Name
	}
	if existing.Summary == "" {
		existing.Summary = incoming.Summary
	} else if incoming.Summary != "" && !containsString(splitSummary(existing.Summary), incoming.Summary) {
		existing.Summary += "; " + incoming.Summary
	}
	if existing.Command == "" {
		existing.Command = incoming.Command
	}
	for _, source := range incoming.Sources {
		existing.Sources = appendSource(existing.Sources, source)
	}
	if existing.Source.Path == "" && len(existing.Sources) > 0 {
		existing.Source = existing.Sources[0]
	}
	existing.Tags = uniqueStrings(append(existing.Tags, incoming.Tags...))
	return existing
}

func appendSource(sources []Source, source Source) []Source {
	if source.Path == "" && source.Line == 0 {
		return sources
	}
	for _, existing := range sources {
		if existing == source {
			return sources
		}
	}
	return append(sources, source)
}

func uniqueStrings(values []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, value := range values {
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		out = append(out, value)
	}
	return out
}

func splitSummary(summary string) []string {
	return strings.Split(summary, "; ")
}

func containsString(values []string, needle string) bool {
	for _, value := range values {
		if value == needle {
			return true
		}
	}
	return false
}
