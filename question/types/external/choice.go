package external

import "github.com/rendis/surveygo/v2/question/types"

// ExternalQuestion represents an external question type.
// Types:
// - types.QTypeExternalQuestion
type ExternalQuestion struct {
	types.QBase

	// Default is the default value for the question.
	// Validations:
	// - optional
	Default *string `json:"default" bson:"default" validate:"omitempty"`

	// QuestionType is the type of the external question.
	QuestionType types.QuestionType `json:"questionType" bson:"questionType" validate:"required,questionType"`

	// ExternalType is the type of the external source.
	// Validations:
	// - required
	// - min length: 1
	ExternalType string `json:"externalType" bson:"externalType" validate:"required,min=1"`

	// Description is a description for the choice field.
	// Validations:
	// - optional
	Description *string `json:"description" bson:"description" validate:"omitempty"`

	// Src is the source of the external source.
	// Validations:
	// - optional
	Src *string `json:"src" bson:"src" validate:"omitempty"`
}
