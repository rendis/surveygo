package text

import (
	"fmt"
	"github.com/rendis/surveygo/v2/question/types"
)

// Telephone represents a telephone question type.
// QuestionType: types.QTypeTelephone
type Telephone struct {
	types.QBase `json:",inline" bson:",inline"`

	// AllowedCountryCodes list of allowed country codes for the telephone field.
	// Validations:
	// - optional
	// - if defined, each country code must have a length of at least 1
	AllowedCountryCodes []string `json:"allowedCountryCodes,omitempty" bson:"allowedCountryCodes,omitempty" validate:"omitempty,dive,min=1"`
}

// CastToTelephone casts the given interface to a Telephone type.
func CastToTelephone(questionValue any) (*Telephone, error) {
	c, ok := questionValue.(*Telephone)
	if !ok {
		return nil, fmt.Errorf("invalid type, expected *text.Telephone, got %T", questionValue)
	}
	return c, nil
}
