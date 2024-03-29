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

	// Metadata is a map of key-value pairs that can be used to store additional information about the question.
	// Validations:
	// - optional
	Metadata map[string]any `json:"metadata,omitempty" bson:"metadata,omitempty" validate:"omitempty"`

	// Collapsible is a flag that indicates if the question is collapsible.
	// Validations:
	// - optional
	Collapsible *bool `json:"collapsible,omitempty" bson:"collapsible,omitempty"`

	// Collapsed is a flag that indicates if the question is collapsed.
	// Validations:
	// - optional
	Collapsed *bool `json:"collapsed,omitempty" bson:"collapsed,omitempty"`

	// Color is the color of the question.
	// Validations:
	// - optional
	// - min length: 1
	Color *string `json:"color,omitempty" bson:"color,omitempty" validate:"omitempty,min=1"`

	// Defaults is the list of default values for the question.
	// Validations:
	// - optional
	Defaults []string `json:"defaults,omitempty" bson:"defaults,omitempty" validate:"omitempty"`
}

// QuestionType represents the different types of questions that can exist in a survey.
type QuestionType string

// Any new type, depending on the type of question, should be added to the following maps:
// - Maps:
//   - QTypeChoiceTypes if type is a choice type
//   - QTypeTextTypes if type is a text type
//   - QTypeAssetTypes if type is an asset type
//
// - Add to serde functions:
//   - question.UnmarshalJSON
//   - question.UnmarshalBSONValue
//
// - Add reviewer function to reviewer.GetQuestionReviewer
//
// if the new type cant be added to either of the above, then:
// - add new types below
// - create a new map of QuestionType and add the new types to it (e.g. QTypeChoiceTypes or QTypeTextTypes)
// - create new struct for the new type (e.g. choice.Choice)
// - add new map to ParseToQuestionType function
// - add a new function to check if the new type is of the new slice (e.g. IsChoiceType or IsTextType)
// - add reviewer function for the new type in reviewer.GetQuestionReviewer function
// - add news struct to the switch case in question.UnmarshalJSON and question.UnmarshalBSONValue functions
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

	// QTypeToggle represents a toggle field type
	QTypeToggle = "toggle"

	// QTypeSlider represents a slider field type
	QTypeSlider = "slider"

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

	// QTypeIdentificationNumber represents an identification number field type
	QTypeIdentificationNumber = "identification_number"

	// QTypeDateTime represents a date time field type
	QTypeDateTime = "date_time"

	//------ Asset types ------//

	// QTypeImage represents an image field type
	QTypeImage = "image"

	// QTypeVideo represents a video field type
	QTypeVideo = "video"

	// QTypeAudio represents an audio field type
	QTypeAudio = "audio"

	// QTypeDocument represents a document field type
	QTypeDocument = "document"

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
var QTypeChoiceTypes = joinMaps(QTypeSimpleChoiceTypes, QTypeComplexChoiceTypes)

var QTypeSimpleChoiceTypes = map[QuestionType]bool{
	QTypeSingleSelect:   true,
	QTypeMultipleSelect: true,
	QTypeRadio:          true,
	QTypeCheckbox:       true,
}

var QTypeComplexChoiceTypes = map[QuestionType]bool{
	QTypeToggle: true,
	QTypeSlider: true,
}

// QTypeTextTypes groups all text types.
var QTypeTextTypes = map[QuestionType]bool{
	QTypeTextArea:             true,
	QTypeInputText:            true,
	QTypeEmail:                true,
	QTypeTelephone:            true,
	QTypeInformation:          true,
	QTypeIdentificationNumber: true,
	QTypeDateTime:             true,
}

// QTypeExternalQuestions groups all external types.
var QTypeExternalQuestions = map[QuestionType]bool{
	QTypeExternalQuestion: true,
}

// QTypeAssetTypes groups all asset types.
var QTypeAssetTypes = map[QuestionType]bool{
	QTypeImage:    true,
	QTypeVideo:    true,
	QTypeAudio:    true,
	QTypeDocument: true,
}

// IsChoiceType returns true if the question type is a choice type, false otherwise.
func IsChoiceType(qt QuestionType) bool {
	return QTypeChoiceTypes[qt]
}

// IsSimpleChoiceType returns true if the question type is a simple choice type, false otherwise.
func IsSimpleChoiceType(qt QuestionType) bool {
	return QTypeSimpleChoiceTypes[qt]
}

// IsComplexChoiceType returns true if the question type is a complex choice type, false otherwise.
func IsComplexChoiceType(qt QuestionType) bool {
	return QTypeComplexChoiceTypes[qt]
}

// IsTextType returns true if the question type is a text type, false otherwise.
func IsTextType(qt QuestionType) bool {
	return QTypeTextTypes[qt]
}

// IsExternalType returns true if the question type is an external type, false otherwise.
func IsExternalType(qt QuestionType) bool {
	return QTypeExternalQuestions[qt]
}

// IsAssetType returns true if the question type is an asset type, false otherwise.
func IsAssetType(qt QuestionType) bool {
	return QTypeAssetTypes[qt]
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

	if _, ok := QTypeAssetTypes[tmpQT]; ok {
		return tmpQT, nil
	}

	return "", fmt.Errorf("invalid question type '%s'", v)
}

func joinMaps[T any](maps ...map[QuestionType]T) map[QuestionType]T {
	res := make(map[QuestionType]T)
	for _, m := range maps {
		for k, v := range m {
			res[k] = v
		}
	}
	return res
}
