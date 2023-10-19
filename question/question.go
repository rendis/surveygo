package question

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rendis/surveygo/v2/question/types"
	"github.com/rendis/surveygo/v2/question/types/choice"
	"github.com/rendis/surveygo/v2/question/types/text"
)

// BaseQuestion is a struct that contains common fields for all types of questions in a survey.
type BaseQuestion struct {
	// NameId is the identifier of the question.
	// Validations:
	// - required
	// - valid name id
	NameId string `json:"nameId" bson:"nameId" validate:"required,validNameId"`

	// Visible is a flag that indicates if the question is visible.
	Visible bool `json:"visible" bson:"visible"`

	// QTyp is the type of question, such as single_select, multi_select, radio, checkbox, or text_area.
	// Validations:
	// - required
	// - must be a valid question type
	QTyp types.QuestionType `json:"type" bson:"type" validate:"required,questionType"`

	// Label is a label for the question.
	// Validations:
	// - required
	// - min length: 1
	Label string `json:"label" bson:"label" validate:"required,min=1"`

	// Required indicates whether the question is required. Defaults to false.
	Required bool `json:"required" bson:"required"`
}

// Question is a struct that represents a question in a survey.
type Question struct {
	// BaseQuestion contains common fields for all types of questions.
	BaseQuestion

	// Value is the value of the question, which can be of different types depending on the type of question.
	// Validations:
	// - required
	Value any `json:"value" bson:"value" validate:"required"`
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
	case types.QTypeSingleSelect, types.QTypeMultipleSelect, types.QTypeRadio, types.QTypeCheckbox:
		realQuestion, err = getQuestionByType[choice.Choice](b)
	case types.QTypeTextArea, types.QTypeInputText:
		realQuestion, err = getQuestionByType[text.FreeText](b)
	case types.QTypeEmail:
		realQuestion, err = getQuestionByType[text.Email](b)
	case types.QTypeTelephone:
		realQuestion, err = getQuestionByType[text.Telephone](b)
	case types.QTypeInformation:
		realQuestion, err = getQuestionByType[text.InformationText](b)
	default:
		return fmt.Errorf("invalid question type: %s", bq.QTyp)
	}

	if err != nil {
		return errors.Join(fmt.Errorf("error unmarshalling question '%s'", bq.NameId), err)
	}

	*q = *realQuestion
	return nil
}

// getQuestionByType returns a question of a specific type.
func getQuestionByType[T any](b []byte) (*Question, error) {
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
