package external

import "github.com/rendis/surveygo/v2/question/types"

// ExternalQuestion represents an external question type.
// Types:
// - types.QTypeExternalQuestion
type ExternalQuestion struct {
	types.QBase `bson:",inline"`

	// Defaults is the list of default values for the external question field.
	// Validations:
	// - optional
	Defaults []string `json:"defaults,omitempty" bson:"defaults,omitempty" validate:"omitempty"`

	// ExternalType is the type of the external question field
	// Validations:
	// - required
	// - min length: 1
	ExternalType string `json:"externalType,omitempty" bson:"externalType,omitempty" validate:"required,min=1"`

	// Description is a description for the choice field.
	// Validations:
	// - optional
	Description *string `json:"description,omitempty" bson:"description,omitempty" validate:"omitempty"`

	// Src is the source of the external source.
	// Validations:
	// - optional
	Src *string `json:"src,omitempty" bson:"src,omitempty" validate:"omitempty"`
}
