package asset

import (
	"github.com/rendis/surveygo/v2/question/types"
)

// DocumentAsset represents a document asset question type.
// Types:
// - types.QTypeDocument
type DocumentAsset struct {
	types.QBase `json:",inline" bson:",inline"`

	// Caption provides a description or additional information about the document.
	// Validations:
	// - optional
	// - max length: 255
	Caption *string `json:"caption,omitempty" bson:"caption,omitempty" validate:"omitempty,max=255"`

	// MaxSize is the maximum allowed file size for the document in bytes.
	// Validations:
	// - optional
	// - if defined, must be a positive integer
	MaxSize *int64 `json:"maxSize,omitempty" bson:"maxSize,omitempty" validate:"omitempty,gt=0"`

	// Tags are keywords or terms associated with the image.
	// Validations:
	// - optional
	Tags []string `json:"tags,omitempty" bson:"tags,omitempty" validate:"omitempty"`

	// Metadata is a map of key/value pairs that can be used to store additional information about the image.
	// Validations:
	// - optional
	Metadata map[string]any `json:"metadata,omitempty" bson:"metadata,omitempty" validate:"omitempty"`

	// AllowedContentTypes list of allowed document content types (e.g. application/pdf, application/msword, etc.).
	// Validations:
	// - optional
	AllowedContentTypes []string `json:"allowedContentTypes,omitempty" bson:"allowedContentTypes,omitempty" validate:"omitempty"`
}
