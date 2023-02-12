package part

import (
	"encoding/json"
	"fmt"
)

// Choice represents a choice field in a survey.
type Choice struct {
	// Options is a list of options for the choice field.
	Options []Option `json:"options"`
}

// Option represents a single option in a choice field.
type Option struct {
	// ID is a required identifier for the option.
	ID *string `json:"id"`

	// Label is a required label for the option.
	Label *string `json:"label"`

	// Order is an optional order number for the option.
	Order *int `json:"order"`

	// SubQuestions is an optional list of sub-questions for the option.
	SubQuestions []Question `json:"subQuestions"`
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (c *Choice) UnmarshalJSON(b []byte) error {
	// Create an alias type for Choice to prevent infinite recursion during unmarshalling.
	type alias Choice
	var a alias
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	*c = Choice(a)
	return nil
}

// ChoiceUnmarshallValidator checks if a Choice is valid.
func ChoiceUnmarshallValidator(c *Choice) error {
	if c == nil {
		return fmt.Errorf("choice is nil")
	}
	if c.Options == nil || len(c.Options) == 0 {
		return fmt.Errorf("invalid choice format, options is not defined")
	}
	return nil
}
