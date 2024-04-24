package choice

import (
	"github.com/rendis/surveygo/v2/question/types"
)

// Toggle represents a toggle question type.
// Types:
// - types.QTypeToggle
type Toggle struct {
	types.QBase `json:",inline" bson:",inline"`

	// Default is the default value for the toggle.
	// Validations:
	// - optional
	Default bool `json:"default,omitempty" bson:"default,omitempty" validate:"omitempty"`

	// OnLabel is the label for the "on" state.
	// Validations:
	// - required
	// - min length: 1
	OnLabel string `json:"onLabel,omitempty" bson:"onLabel,omitempty" validate:"required,min=1"`

	// OffLabel is the label for the "off" state.
	// Validations:
	// - required
	// - min length: 1
	OffLabel string `json:"offLabel,omitempty" bson:"offLabel,omitempty" validate:"required,min=1"`
}
