package surveygo

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ------------ Deserializers ------------- //

// ParseFromJsonStr converts the given json string into a Survey instance.
func ParseFromJsonStr(jsonSurvey string) (*Survey, error) {
	return ParseFromBytes([]byte(jsonSurvey))
}

// ParseFromBytes converts the given json byte slice into a Survey instance.
func ParseFromBytes(b []byte) (*Survey, error) {
	var survey = &Survey{}

	// unmarshal the json survey into the survey struct
	if err := json.Unmarshal(b, survey); err != nil {
		return nil, errors.Join(fmt.Errorf("error unmarshalling json survey"), err)
	}

	// validate the survey struct
	if err := SurveyValidator.Struct(survey); err != nil {
		errs := TranslateValidationErrors(err)
		errs = append([]error{fmt.Errorf("error validating survey")}, errs...)
		return nil, errors.Join(errs...)
	}

	// check survey consistency
	if err := survey.checkConsistency(); err != nil {
		return nil, err
	}

	// run position updater
	survey.positionUpdater()

	return survey, nil
}

// -------------- Serializers -------------- //

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
