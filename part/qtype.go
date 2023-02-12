package part

import (
	"encoding/json"
	"fmt"
)

// QuestionType represents the different types of questions that can exist in a survey.
type QuestionType string

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
	switch v {
	case string(QTypeSingleSelect):
		return QTypeSingleSelect, nil
	case string(QTypeMultipleSelect):
		return QTypeMultipleSelect, nil
	case string(QTypeRadio):
		return QTypeRadio, nil
	case string(QTypeCheckbox):
		return QTypeCheckbox, nil
	case string(QTypeTextArea):
		return QTypeTextArea, nil
	case string(QTypeInputText):
		return QTypeInputText, nil
	case string(QTypeEmail):
		return QTypeEmail, nil
	case string(QTypeTelephone):
		return QTypeTelephone, nil
	default:
		return "", fmt.Errorf("invalid question type '%s'", v)
	}
}
