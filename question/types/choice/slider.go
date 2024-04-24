package choice

import (
	"fmt"
	"github.com/rendis/surveygo/v2/question/types"
)

// Slider represents a choice question type.
// Types:
// - types.QTypeSlider
type Slider struct {
	types.QBase `json:",inline" bson:",inline"`

	// Min is the minimum value for the slider.
	// Validations:
	// - required
	Min int `json:"min,omitempty" bson:"min,omitempty" validate:"required"`

	// Max is the maximum value for the slider.
	// Validations:
	// - required
	Max int `json:"max,omitempty" bson:"max,omitempty" validate:"required"`

	// Step is the step value for the slider.
	// Validations:
	// - required
	// - min: 1
	Step int `json:"step,omitempty" bson:"step,omitempty" validate:"required,min=1"`

	// Default is the default value for the slider.
	// Validations:
	// - optional
	Default int `json:"default,omitempty" bson:"default,omitempty" validate:"omitempty"`

	// Unit is the unit of the slider (e.g. years, months, days, etc.).
	// Validations:
	// - optional
	// - min length: 1
	Unit string `json:"unit,omitempty" bson:"unit,omitempty" validate:"omitempty,min=1"`
}

// CastToSlider casts an interface to a Slider type.
func CastToSlider(i any) (*Slider, error) {
	c, ok := i.(*Slider)
	if !ok {
		return nil, fmt.Errorf("invalid type, expected *choice.Slider, got %T", i)
	}
	return c, nil
}
