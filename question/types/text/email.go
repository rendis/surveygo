package text

import (
	"fmt"
	"github.com/rendis/surveygo/v2/question/types"
)

// Email represents an email question type.
// QuestionType: types.QTypeEmail
type Email struct {
	types.QBase `bson:",inline"`

	// AllowedDomains list of allowed domains for the email field.
	// Validations:
	// - optional
	// - if defined, each domain must have a length of at least 1
	AllowedDomains []string `json:"allowedDomains,omitempty" bson:"allowedDomains,omitempty" validate:"omitempty,dive,min=1"`
}

// CastToEmail casts the given interface to an Email type.
func CastToEmail(questionValue any) (*Email, error) {
	c, ok := questionValue.(*Email)
	if !ok {
		return nil, fmt.Errorf("invalid type, expected *text.Email, got %T", questionValue)
	}
	return c, nil
}
