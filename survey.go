package surveygo

import (
	"fmt"
	"github.com/rendis/surveygo/check"
	"github.com/rendis/surveygo/part"
	"github.com/tidwall/gjson"
	"strings"
)

type Answers map[string][]any

type Survey struct {
	Title       *string           `json:"title"`
	Version     *string           `json:"version"`
	Description *string           `json:"description"`
	NameIdPaths map[string]string `json:"idPaths"`
	JsonSurvey  *string           `json:"jsonSurvey"`
}

func (s *Survey) Check(aws Answers) error {
	gres := gjson.Parse(*s.JsonSurvey)
	for nameId, values := range aws {
		path, ok := s.NameIdPaths[nameId]
		if !ok {
			return fmt.Errorf("nameId not found: %s", nameId)
		}
		obj := gres.Get(path)
		typ := obj.Get("type")

		qt, err := part.ParseToQuestionType(typ.String())
		if err != nil {
			return err
		}
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

type internalSurvey struct {
	Title       *string         `json:"title"`
	Version     *string         `json:"version"`
	Description *string         `json:"description"`
	Questions   []part.Question `json:"questions"`
}

func (s *internalSurvey) getNameIdPaths() (map[string]string, error) {
	base := []string{"questions"}
	var paths = make(map[string]string)
	for i, q := range s.Questions {
		p := q.GetNameIdPaths(append(base, fmt.Sprintf("%d", i)))
		for _, ip := range p {
			if _, ok := paths[ip.NameId]; ok {
				return nil, fmt.Errorf("duplicate nameId: %s", ip.NameId)
			}
			paths[ip.NameId] = strings.Join(ip.Path, ".")
		}
	}
	return paths, nil
}

func (s *internalSurvey) validate() error {
	if s.Title == nil || *s.Title == "" {
		return fmt.Errorf("survey title is required")
	}

	if s.Version == nil || *s.Version == "" {
		return fmt.Errorf("survey version is required")
	}

	return nil
}
