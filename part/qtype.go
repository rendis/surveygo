package part

import (
	"encoding/json"
	"fmt"
)

// QuestionType represents the different types of questions that can exist in a survey.
type QuestionType string

// Any new type, depending on the type of question, should be added to the following:
// - QTypeChoiceTypes if type is a choice type
// - QTypeTextTypes if type is a text type
// if the new type cant be added to either of the above, then:
// - create a new slice of QuestionType and add the new type to it
// - append the new slice to AllQuestionTypes
// - add a new function to check if the new type is of the new slice (e.g. IsChoiceType or IsTextType)
// - update this comment to include the new type :)
const (
	// QTypeSingleSelect represents a single select field type
	QTypeSingleSelect QuestionType = "single_select"

	// QTypeMultipleSelect represents a multiple select field type
	QTypeMultipleSelect = "multi_select"

	// QTypeRadio represents a radio field type
	QTypeRadio = "radio"

	// QTypeCheckbox represents a checkbox field type
	QTypeCheckbox = "checkbox"

	// QTypeTextArea represents a text area field type
	QTypeTextArea = "text_area"

	// QTypeInputText represents a text input field type
	QTypeInputText = "input_text"

	// QTypeEmail represents an email input field type
	QTypeEmail = "email"

	// QTypeTelephone represents a telephone input field type
	QTypeTelephone = "telephone"
)

var QTypeChoiceTypes = []QuestionType{
	QTypeSingleSelect, QTypeMultipleSelect, QTypeRadio, QTypeCheckbox,
}

var QTypeTextTypes = []QuestionType{
	QTypeTextArea, QTypeInputText, QTypeEmail, QTypeTelephone,
}

// AllQuestionTypes is a slice of all question types. Used for validation.
var AllQuestionTypes = append(QTypeChoiceTypes, QTypeTextTypes...)

// IsChoiceType returns true if the question type is a choice type, false otherwise.
func IsChoiceType(qt QuestionType) bool {
	for _, t := range QTypeChoiceTypes {
		if t == qt {
			return true
		}
	}
	return false
}

// IsTextType returns true if the question type is a text type, false otherwise.
func IsTextType(qt QuestionType) bool {
	for _, t := range QTypeTextTypes {
		if t == qt {
			return true
		}
	}
	return false
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *QuestionType) UnmarshalJSON(b []byte) error {
	var st string
	if err := json.Unmarshal(b, &st); err != nil {
		return fmt.Errorf("unmarshal error, %s", err)
	}

	t, err := ParseToQuestionType(st)
	if err != nil {
		return fmt.Errorf("parse error, %s", err)
	}
	*s = t
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (s *QuestionType) MarshalJSON() ([]byte, error) {
	if s == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(string(*s))
}

// ParseToQuestionType takes a string and returns the corresponding QuestionType, or an error if the string is invalid.
func ParseToQuestionType(v string) (QuestionType, error) {
	tmpQT := QuestionType(v)
	for _, qt := range AllQuestionTypes {
		if qt == tmpQT {
			return tmpQT, nil
		}
	}
	return "", fmt.Errorf("invalid question type '%s'", v)
}
