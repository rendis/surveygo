package asset

import (
	"github.com/rendis/surveygo/v2/question/types"
)

// AudioAsset represents an audio asset question type.
// Types:
// - types.QTypeAudio
type AudioAsset struct {
	types.QBase `bson:",inline"`

	// Caption provides a description or additional information about the audio.
	// Validations:
	// - optional
	// - max length: 255
	Caption *string `json:"caption,omitempty" bson:"caption,omitempty" validate:"omitempty,max=255"`

	// MaxSize is the maximum allowed file size for the audio in bytes.
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

	// AllowedContentTypes list of allowed audio content types (e.g. audio/mpeg, audio/wav, etc.).
	// Validations:
	// - optional
	AllowedContentTypes []string `json:"allowedContentTypes,omitempty" bson:"allowedContentTypes,omitempty" validate:"omitempty"`
}
