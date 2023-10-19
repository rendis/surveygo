package text

import (
	"github.com/rendis/surveygo/v2/question/types"
)

// InformationText represents an information text question type.
// Types:
// - types.QTypeInformation
type InformationText struct {
	types.Base

	// Text is the text to be displayed.
	// Validations:
	// - required
	// - min length: 1
	Text string `json:"text" bson:"text" validate:"required,min=1"`
}
