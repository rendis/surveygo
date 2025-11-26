package asset

import (
	"github.com/rendis/surveygo/v2/question/types"
)

// ImageAsset represents an image asset question type.
// Types:
// - types.QTypeImage
type ImageAsset struct {
	types.QBase `json:",inline" bson:",inline"`

	// AltText provides alternative text for the image to improve accessibility.
	// Validations:
	// - optional
	// - max length: 255
	AltText *string `json:"altText,omitempty" bson:"altText,omitempty" validate:"omitempty,max=255"`

	// Tags are keywords or terms associated with the image.
	// Validations:
	// - optional
	Tags []string `json:"tags,omitempty" bson:"tags,omitempty" validate:"omitempty"`

	// Metadata is a map of key/value pairs that can be used to store additional information about the image.
	// Validations:
	// - optional
	Metadata map[string]any `json:"metadata,omitempty" bson:"metadata,omitempty" validate:"omitempty"`

	// MaxSize is the maximum allowed file size for the image in bytes.
	// Validations:
	// - optional
	// - if defined, must be a positive integer
	MaxSize *int64 `json:"maxSize,omitempty" bson:"maxSize,omitempty" validate:"omitempty,gt=0"`

	// AllowedContentTypes list of allowed image types (e.g., "image/png", "image/jpeg").
	// Validations:
	// - optional
	AllowedContentTypes []string `json:"allowedContentTypes,omitempty" bson:"allowedContentTypes,omitempty" validate:"omitempty"`

	// MaxFiles is the maximum number of files that can be uploaded.
	// Validations:
	// - optional
	// - if defined, must be >= 1
	// Note: consuming code should treat 0 as default value of 1
	MaxFiles int `json:"maxFiles,omitempty" bson:"maxFiles,omitempty" validate:"omitempty,min=1"`

	// MinFiles is the minimum number of files that must be uploaded.
	// Validations:
	// - optional
	// - if defined, must be >= 0
	// Note: consuming code should treat 0 as default value of 1
	MinFiles int `json:"minFiles,omitempty" bson:"minFiles,omitempty" validate:"omitempty,min=0"`
}
