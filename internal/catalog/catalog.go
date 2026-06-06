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
	Tags    []string  `json:"tags,omitempty" yaml:"tags,omitempty"`
}

type Catalog struct {
	entries []Entry
}

func New(entries ...Entry) Catalog {
	c := Catalog{}
	for _, entry := range entries {
		c.Add(entry)
	}
	return c
}

func (c *Catalog) Add(entry Entry) {
	if entry.ID == "" {
		entry.ID = NormalizeID(entry.Name)
	}
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
