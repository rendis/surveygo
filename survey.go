package surveygo

import (
	"encoding/json"
	"fmt"
	"github.com/rendis/surveygo/check"
	"github.com/rendis/surveygo/part"
	"github.com/tidwall/gjson"
	"strings"
)

// Answers is a map of question nameId to its answer value(s).
type Answers map[string][]any

// Survey is a struct representation of a survey.
type Survey struct {
	// Title is the title of the survey.
	Title *string `json:"title"`

	// Version is the version of the survey.
	Version *string `json:"version"`

	// Description is the description of the survey.
	Description *string `json:"description"`

	// NameIdPaths is a map of question nameId to its path in the survey json.
	NameIdPaths map[string]string `json:"idPaths"`

	// FullJsonSurvey is the full json representation of the survey.
	FullJsonSurvey *string `json:"fullJsonSurvey"`
}

// Check verifies if the answers provided are valid for this survey.
func (s *Survey) Check(aws Answers) error {
	// parse the full json survey into gjson object
	gres := gjson.Parse(*s.FullJsonSurvey)

	// iterate through each answer
	for nameId, values := range aws {
		// find the path of the question in the json survey
		path, ok := s.NameIdPaths[nameId]
		if !ok {
			return fmt.Errorf("nameId not found: %s", nameId)
		}

		// get the question object from the json survey
		obj := gres.Get(path)

		// get the type of the question
		typ := obj.Get("type")

		// parse the string type to question type
		qt, err := part.ParseToQuestionType(typ.String())
		if err != nil {
			return err
		}

		// validate the answers based on question type
		switch qt {
		case part.QTypeSingleSelect, part.QTypeMultipleSelect, part.QTypeRadio, part.QTypeCheckbox:
			err = check.ValidateChoice(obj, values, qt)
			if err != nil {
				return fmt.Errorf("check error for nameId: '%s', path: '%s', error: %s", nameId, path, err)
			}
		}
	}
	return nil
}

// ToMap returns the survey as a map representation.
func (s *Survey) ToMap() (map[string]any, error) {
	r := make(map[string]any)
	if s.FullJsonSurvey != nil {
		err := json.Unmarshal([]byte(*s.FullJsonSurvey), &r)
		if err != nil {
			return nil, err
		}
	}
	return r, nil
}

// jsonSurvey represents the internal structure of the survey.
type jsonSurvey struct {
	// Title is the title of the survey.
	Title *string `json:"title"`

	// Version is the version of the survey.
	Version *string `json:"version"`

	// Description is the description of the survey.
	Description *string `json:"description"`

	// Questions is the list of questions in the survey.
	Questions []part.Question `json:"questions"`
}

// getNameIdPaths returns a map of nameIds to their respective paths in the survey's json.
func (s *jsonSurvey) getNameIdPaths() (map[string]string, error) {
	base := []string{"questions"}
	var paths = make(map[string]string)
	for i, question := range s.Questions {
		pathsForQuestion := question.GetNameIdPaths(append(base, fmt.Sprintf("%d", i)))
		for _, ip := range pathsForQuestion {
			if _, ok := paths[ip.NameId]; ok {
				return nil, fmt.Errorf("duplicate nameId: %s", ip.NameId)
			}
			paths[ip.NameId] = strings.Join(ip.Path, ".")
		}
	}
	return paths, nil
}

// validate checks that the required fields are present in the survey.
func (s *jsonSurvey) validate() error {
	if s.Title == nil || *s.Title == "" {
		return fmt.Errorf("survey title is required")
	}

	if s.Version == nil || *s.Version == "" {
		return fmt.Errorf("survey version is required")
	}

	return nil
}

// marshal returns a string representation of the jsonSurvey.
func (s *jsonSurvey) marshal() (string, error) {
	bytes, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
