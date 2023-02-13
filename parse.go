package surveygo

import (
	"encoding/json"
	"fmt"
	"github.com/rendis/surveygo/part"
	"github.com/tidwall/gjson"
)

// NewSurvey creates a new Survey instance with the given title, version, and description.
// Title and version are required and an error is returned if they are not provided.
func NewSurvey(title, version, description *string) (*Survey, error) {
	if title == nil || version == nil || *title == "" || *version == "" {
		return nil, fmt.Errorf("title and version are required")
	}

	j := &jsonSurvey{
		Title:       title,
		Version:     version,
		Description: description,
		Questions:   []part.Question{},
	}
	js, _ := j.marshal()

	return &Survey{
		Title:          title,
		Version:        version,
		Description:    description,
		FullJsonSurvey: &js,
	}, nil
}

// Parse converts the given json string into a Survey instance.
func Parse(j string) (*Survey, error) {
	return ParseBytes([]byte(j))
}

// ParseBytes converts the given json byte slice into a Survey instance.
func ParseBytes(b []byte) (*Survey, error) {
	uglyJson := gjson.ParseBytes(b).Get("@ugly").String()

	var jSvy jsonSurvey
	err := json.Unmarshal([]byte(uglyJson), &jSvy)
	if err != nil {
		return nil, err
	}

	err = jSvy.validate()
	if err != nil {
		return nil, err
	}

	paths, required, err := jSvy.getNameIdPaths()
	if err != nil {
		return nil, err
	}

	return &Survey{
		Title:           jSvy.Title,
		Version:         jSvy.Version,
		Description:     jSvy.Description,
		NameIdPaths:     paths,
		RequiredNameIds: required,
		FullJsonSurvey:  &uglyJson,
	}, nil
}
