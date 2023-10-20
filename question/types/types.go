package types

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// QBase is a struct that contains common fields for all types of questions.
type QBase struct {
	// Placeholder is a placeholder text.
	// Validations:
	// - optional
	// - min length: 1
	Placeholder *string `json:"placeholder,omitempty" bson:"placeholder,omitempty" validate:"omitempty,min=1"`
}

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
	//------ Choice types ------//

	// QTypeSingleSelect represents a single select field type
	QTypeSingleSelect QuestionType = "single_select"

	// QTypeMultipleSelect represents a multiple select field type
	QTypeMultipleSelect = "multi_select"

	// QTypeRadio represents a radio field type
	QTypeRadio = "radio"

	// QTypeCheckbox represents a checkbox field type
	QTypeCheckbox = "checkbox"

	//------ Text types ------//

	// QTypeTextArea represents a text area field type
	QTypeTextArea = "text_area"

	// QTypeInputText represents a text input field type
	QTypeInputText = "input_text"

	// QTypeEmail represents an email input field type
	QTypeEmail = "email"

	// QTypeTelephone represents a telephone input field type
	QTypeTelephone = "telephone"

	// QTypeInformation represents an information field type
	QTypeInformation = "information"

	//------ External types ------//

	// QTypeExternalQuestion represents a external question field type
	QTypeExternalQuestion = "external_question"
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

func (s *QuestionType) UnmarshalBSONValue(typ bsontype.Type, raw []byte) error {
	if typ != bsontype.String {
		return fmt.Errorf("invalid bson value type '%s'", typ.String())
	}

	c, _, ok := bsoncore.ReadString(raw)
	if !ok {
		return fmt.Errorf("invalid bson value '%s'", string(raw))
	}

	t, err := ParseToQuestionType(c)
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

// QTypeChoiceTypes groups all choice types.
var QTypeChoiceTypes = map[QuestionType]bool{
	QTypeSingleSelect:   true,
	QTypeMultipleSelect: true,
	QTypeRadio:          true,
	QTypeCheckbox:       true,
}

// QTypeTextTypes groups all text types.
var QTypeTextTypes = map[QuestionType]bool{
	QTypeTextArea:    true,
	QTypeInputText:   true,
	QTypeEmail:       true,
	QTypeTelephone:   true,
	QTypeInformation: true,
}

// QTypeExternalQuestions groups all external types.
var QTypeExternalQuestions = map[QuestionType]bool{
	QTypeExternalQuestion: true,
}

// IsChoiceType returns true if the question type is a choice type, false otherwise.
func IsChoiceType(qt QuestionType) bool {
	return QTypeChoiceTypes[qt]
}

// IsTextType returns true if the question type is a text type, false otherwise.
func IsTextType(qt QuestionType) bool {
	return QTypeTextTypes[qt]
}

// IsExternalType returns true if the question type is an external type, false otherwise.
func IsExternalType(qt QuestionType) bool {
	return QTypeExternalQuestions[qt]
}

// ParseToQuestionType takes a string and returns the corresponding QuestionType, or an error if the string is invalid.
func ParseToQuestionType(v string) (QuestionType, error) {
	tmpQT := QuestionType(v)
	if _, ok := QTypeChoiceTypes[tmpQT]; ok {
		return tmpQT, nil
	}

	if _, ok := QTypeTextTypes[tmpQT]; ok {
		return tmpQT, nil
	}

	if _, ok := QTypeExternalQuestions[tmpQT]; ok {
		return tmpQT, nil
	}

	return "", fmt.Errorf("invalid question type '%s'", v)
}
