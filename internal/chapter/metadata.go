package chapter

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Difficulty describes an exercise, independently of its format or chapter.
type Difficulty string

const (
	DifficultyBeginner     Difficulty = "beginner"
	DifficultyIntermediate Difficulty = "intermediate"
	DifficultyAdvanced     Difficulty = "advanced"
)

// Family groups chapters without changing their historical identifiers.
type Family string

const (
	FamilyUnixShell         Family = "unix-shell"
	FamilyKeyboardShortcuts Family = "keyboard-shortcuts"
	FamilyModernCLIMacOS    Family = "modern-cli-macos"
)

// LearningOrder declares editorial intent; it does not reorder exercises.
type LearningOrder string

const (
	LearningOrderGuided      LearningOrder = "guided"
	LearningOrderIndependent LearningOrder = "independent"
)

// UnmarshalYAML applies the historical default only when learning_order is omitted.
func (c *Chapter) UnmarshalYAML(node *yaml.Node) error {
	fields, err := metadataFields(node, "family", "learning_order")
	if err != nil {
		return fmt.Errorf("chapter %v: %w", fields["id"], err)
	}
	type plain Chapter
	decoded := plain{LearningOrder: LearningOrderGuided}
	if err := node.Decode(&decoded); err != nil {
		return fmt.Errorf("chapter %v: %w", fields["id"], err)
	}
	*c = Chapter(decoded)
	return nil
}

// UnmarshalYAML defaults legacy exercises to Beginner without changing their YAML.
func (item *Item) UnmarshalYAML(node *yaml.Node) error {
	fields, err := metadataFields(node, "difficulty")
	if err != nil {
		return fmt.Errorf("item %v: %w", fields["id"], err)
	}
	if value, present := fields["tools"]; present {
		tools, ok := value.([]any)
		if !ok {
			return fmt.Errorf("item %v: tools must be a list of strings", fields["id"])
		}
		for _, tool := range tools {
			if _, ok := tool.(string); !ok {
				return fmt.Errorf("item %v: tools must contain only string identifiers", fields["id"])
			}
		}
	}
	type plain Item
	decoded := plain{Difficulty: DifficultyBeginner}
	if err := node.Decode(&decoded); err != nil {
		return fmt.Errorf("item %v: %w", fields["id"], err)
	}
	*item = Item(decoded)
	return nil
}

// Decode a mapping to distinguish omission from null and YAML scalar coercion.
// Let yaml.v3 resolve aliases and merge keys before checking metadata types.
// Unknown historical fields remain accepted by the normal struct decoder.
func metadataFields(node *yaml.Node, stringFields ...string) (map[string]any, error) {
	var fields map[string]any
	if err := node.Decode(&fields); err != nil {
		return fields, err
	}
	for _, name := range stringFields {
		if value, present := fields[name]; present {
			text, ok := value.(string)
			if !ok || text == "" {
				return fields, fmt.Errorf("%s must be a non-empty string", name)
			}
		}
	}
	return fields, nil
}
