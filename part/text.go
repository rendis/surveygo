package part

import "fmt"

// TextArea represents a text area field type.
type TextArea struct {
	// Placeholder is an optional placeholder for the text area field.
	Placeholder *string `json:"placeholder"`

	// Min is an optional minimum length for the text area field.
	Min *int `json:"min"`

	// Max is an optional maximum length for the text area field.
	Max *int `json:"max"`
}

// TextAreaUnmarshallValidator validates the unmarshalled text area field.
func TextAreaUnmarshallValidator(t *TextArea) error {
	if t == nil {
		return fmt.Errorf("text area is nil")
	}
	if t.Min != nil && t.Max != nil && *t.Min > *t.Max {
		return fmt.Errorf("min length is greater than max length")
	}
	return nil
}
