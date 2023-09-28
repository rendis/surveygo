package text

import (
	"fmt"
	"github.com/rendis/surveygo/v2/question/types"
)

// FreeText represents a free text question type.
// Types:
// - types.QTypeInputText
// - types.QTypeTextArea
type FreeText struct {
	types.Base

	// Min is an optional minimum length for the text area field.
	// Validations:
	// - optional
	// - if defined:
	//   * must be greater than or equal to 0
	//   * if max is defined, must be less than to max
	Min *int `json:"min" bson:"min" validate:"omitempty,min=0,ltfield=Max"`

	// Max is an optional maximum length for the text area field.
	// Validations:
	// - optional
	// - if defined:
	//   * must be greater than or equal to 0
	//   * if min is defined, must be greater than to min
	Max *int `json:"max" bson:"max" validate:"omitempty,min=0,gtfield=Min"`
}

// CastToFreeText casts the given interface to a FreeText type.
func CastToFreeText(questionValue any) (*FreeText, error) {
	c, ok := questionValue.(*FreeText)
	if !ok {
		return nil, fmt.Errorf("invalid type, expected *text.FreeText, got %T", questionValue)
	}
	return c, nil
}
