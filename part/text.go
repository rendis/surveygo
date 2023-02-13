package part

import "fmt"

// baseText is a struct that contains common fields for all types of text fields.
type baseText struct {
	// Placeholder is an optional placeholder text for the question.
	Placeholder *string `json:"placeholder"`
}

// TextArea represents a text area field type.
// It is used to represent the value of a question of type QTypeTextArea and QTypeInputText.
type TextArea struct {
	baseText

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

// Email represents an email field type.
type Email struct {
	baseText

	// AllowedDomains is an optional list of allowed domains for the email field.
	AllowedDomains []string `json:"allowedDomains"`
}

// EmailUnmarshallValidator validates the unmarshalled email field.
func EmailUnmarshallValidator(e *Email) error {
	if e == nil {
		return fmt.Errorf("email is nil")
	}
	return nil
}

// Telephone represents a telephone field type.
type Telephone struct {
	baseText

	// AllowedCountryCodes is an optional list of allowed country codes for the telephone field.
	AllowedCountryCodes []string `json:"allowedCountryCodes"`
}

// TelephoneUnmarshallValidator validates the unmarshalled telephone field.
func TelephoneUnmarshallValidator(t *Telephone) error {
	if t == nil {
		return fmt.Errorf("telephone is nil")
	}
	return nil
}
