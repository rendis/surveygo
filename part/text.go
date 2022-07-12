package part

import "fmt"

type TextArea struct {
	Placeholder *string `json:"placeholder"` // Optional
	Min         *int    `json:"min"`         // Optional
	Max         *int    `json:"max"`         // Optional
}

func TextAreaUnmarshallValidator(t *TextArea) error {
	if t.Min != nil && t.Max != nil && *t.Min > *t.Max {
		return fmt.Errorf("min is greater than max")
	}
	return nil
}
