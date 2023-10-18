package surveygo

import (
	"encoding/json"
)

// ToMap returns a map representation of the survey.
func (s *Survey) ToMap() (map[string]any, error) {
	r := make(map[string]any)
	// s to map
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	// unmarshal to map[string]any
	if err = json.Unmarshal(b, &r); err != nil {
		return nil, err
	}

	return r, nil
}

// ToJson returns a JSON string representation of the survey.
func (s *Survey) ToJson() (string, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
