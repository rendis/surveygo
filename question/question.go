package question

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rendis/surveygo/v2/question/types"
	"github.com/rendis/surveygo/v2/question/types/asset"
	"github.com/rendis/surveygo/v2/question/types/choice"
	"github.com/rendis/surveygo/v2/question/types/external"
	"github.com/rendis/surveygo/v2/question/types/text"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// BaseQuestion is a struct that contains common fields for all types of questions in a survey.
type BaseQuestion struct {
	// NameId is the identifier of the question.
	// Validations:
	// - required
	// - valid name id
	NameId string `json:"nameId" bson:"nameId" validate:"required,validNameId"`

	// Visible is a flag that indicates if the question is visible.
	Visible bool `json:"visible,omitempty" bson:"visible,omitempty"`

	// QTyp is the type of question, such as single_select, multi_select, radio, checkbox, or text_area.
	// Validations:
	// - required
	// - must be a valid question type
	QTyp types.QuestionType `json:"type,omitempty" bson:"type,omitempty" validate:"required,questionType"`

	// Label is a label for the question.
	// Validations:
	// - required
	// - min length: 1
	Label string `json:"label,omitempty" bson:"label,omitempty" validate:"omitempty,min=1"`

	// Required indicates whether the question is required. Defaults to false.
	Required bool `json:"required,omitempty" bson:"required,omitempty"`

	// Metadata is a map of key-value pairs that can be used to store additional information about the question.
	// Validations:
	// - optional
	Metadata map[string]any `json:"metadata,omitempty" bson:"metadata,omitempty" validate:"omitempty"`

	// Position is the position of the question in the survey.
	// This field is calculated automatically by the system.
	// Validations:
	// - required
	// - min: 1
	Position int `json:"position,omitempty" bson:"position,omitempty" validate:"omitempty,min=1"`

	// Disabled indicates whether the question is disabled. Defaults to false.
	Disabled bool `json:"disabled,omitempty" bson:"disabled,omitempty"`
}

// Question is a struct that represents a question in a survey.
type Question struct {
	// BaseQuestion contains common fields for all types of questions.
	BaseQuestion `bson:",inline"`

	// Value is the value of the question, which can be of different types depending on the type of question.
	// Validations:
	// - required
	Value any `json:"value,omitempty" bson:"value,omitempty" validate:"required"`
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (q *Question) UnmarshalJSON(b []byte) error {
	var bq BaseQuestion

	if err := json.Unmarshal(b, &bq); err != nil {
		return err
	}

	var realQuestion *Question
	var err error

	// unmarshal the question based on its type
	switch bq.QTyp {
	case types.QTypeSingleSelect, types.QTypeMultipleSelect, types.QTypeRadio, types.QTypeCheckbox, types.QTypeToggle:
		realQuestion, err = unmarshalJSONQuestionByType[choice.Choice](b)
	case types.QTypeSlider:
		realQuestion, err = unmarshalJSONQuestionByType[choice.Slider](b)
	case types.QTypeTextArea, types.QTypeInputText:
		realQuestion, err = unmarshalJSONQuestionByType[text.FreeText](b)
	case types.QTypeEmail:
		realQuestion, err = unmarshalJSONQuestionByType[text.Email](b)
	case types.QTypeTelephone:
		realQuestion, err = unmarshalJSONQuestionByType[text.Telephone](b)
	case types.QTypeInformation:
		realQuestion, err = unmarshalJSONQuestionByType[text.InformationText](b)
	case types.QTypeIdentificationNumber:
		realQuestion, err = unmarshalJSONQuestionByType[text.IdentificationNumber](b)
	case types.QTypeDateTime:
		realQuestion, err = unmarshalJSONQuestionByType[text.DateTime](b)
	case types.QTypeExternalQuestion:
		realQuestion, err = unmarshalJSONQuestionByType[external.ExternalQuestion](b)
	case types.QTypeImage:
		realQuestion, err = unmarshalJSONQuestionByType[asset.ImageAsset](b)
	case types.QTypeVideo:
		realQuestion, err = unmarshalJSONQuestionByType[asset.VideoAsset](b)
	case types.QTypeAudio:
		realQuestion, err = unmarshalJSONQuestionByType[asset.AudioAsset](b)
	case types.QTypeDocument:
		realQuestion, err = unmarshalJSONQuestionByType[asset.DocumentAsset](b)
	default:
		return fmt.Errorf("invalid question type: %s", bq.QTyp)
	}

	if err != nil {
		return errors.Join(fmt.Errorf("error unmarshalling question '%s'", bq.NameId), err)
	}

	*q = *realQuestion
	return nil
}

func (q *Question) UnmarshalBSONValue(typ bsontype.Type, b []byte) error {
	var bq BaseQuestion

	if err := bson.Unmarshal(b, &bq); err != nil {
		return fmt.Errorf("BSON unmarshal error, %s", err)
	}

	var value any
	var err error

	// unmarshal the question based on its type
	switch bq.QTyp {
	case types.QTypeSingleSelect, types.QTypeMultipleSelect, types.QTypeRadio, types.QTypeCheckbox, types.QTypeToggle:
		value, err = unmarshalBSONQuestionByType[choice.Choice](b)
	case types.QTypeSlider:
		value, err = unmarshalBSONQuestionByType[choice.Slider](b)
	case types.QTypeTextArea, types.QTypeInputText:
		value, err = unmarshalBSONQuestionByType[text.FreeText](b)
	case types.QTypeEmail:
		value, err = unmarshalBSONQuestionByType[text.Email](b)
	case types.QTypeTelephone:
		value, err = unmarshalBSONQuestionByType[text.Telephone](b)
	case types.QTypeInformation:
		value, err = unmarshalBSONQuestionByType[text.InformationText](b)
	case types.QTypeIdentificationNumber:
		value, err = unmarshalBSONQuestionByType[text.IdentificationNumber](b)
	case types.QTypeDateTime:
		value, err = unmarshalBSONQuestionByType[text.DateTime](b)
	case types.QTypeExternalQuestion:
		value, err = unmarshalBSONQuestionByType[external.ExternalQuestion](b)
	case types.QTypeImage:
		value, err = unmarshalBSONQuestionByType[asset.ImageAsset](b)
	case types.QTypeVideo:
		value, err = unmarshalBSONQuestionByType[asset.VideoAsset](b)
	case types.QTypeAudio:
		value, err = unmarshalBSONQuestionByType[asset.AudioAsset](b)
	case types.QTypeDocument:
		value, err = unmarshalBSONQuestionByType[asset.DocumentAsset](b)
	default:
		return fmt.Errorf("invalid question type: %s", bq.QTyp)
	}

	if err != nil {
		return errors.Join(fmt.Errorf("error unmarshalling question '%s'", bq.NameId), err)
	}

	*q = Question{
		BaseQuestion: bq,
		Value:        value,
	}
	return nil
}

// unmarshalJSONQuestionByType returns a question of a specific type.
func unmarshalJSONQuestionByType[T any](b []byte) (*Question, error) {
	// build a temporary struct with the base question and the value of the specific type
	var tq = struct {
		BaseQuestion
		Value *T `json:"value"`
	}{}

	if err := json.Unmarshal(b, &tq); err != nil {
		return nil, err
	}

	if tq.Value == nil {
		return nil, fmt.Errorf("value is not defined")
	}

	return &Question{
		BaseQuestion: tq.BaseQuestion,
		Value:        tq.Value,
	}, nil
}

func unmarshalBSONQuestionByType[T any](b []byte) (*T, error) {
	var tq = struct {
		Value *T `bson:"value"`
	}{}

	if err := bson.Unmarshal(b, &tq); err != nil {
		return nil, err
	}

	if tq.Value == nil {
		return nil, fmt.Errorf("value is not defined")
	}

	return tq.Value, nil
}
