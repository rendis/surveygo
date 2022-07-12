package surveygo

import (
	"encoding/json"
	"fmt"
	"github.com/rendis/surveygo/part"
	"github.com/tidwall/gjson"
)

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

func Parse(j string) (*Survey, error) {
	b := []byte(j)
	return ParseBytes(b)
}

func ParseBytes(b []byte) (*Survey, error) {
	uglyJson := gjson.ParseBytes(b).Get("@ugly").String()

	var ins jsonSurvey
	err := json.Unmarshal([]byte(uglyJson), &ins)
	if err != nil {
		return nil, err
	}

	err = ins.validate()
	if err != nil {
		return nil, err
	}

	paths, err := ins.getNameIdPaths()
	if err != nil {
		return nil, err
	}

	return &Survey{
		Title:          ins.Title,
		Version:        ins.Version,
		Description:    ins.Description,
		NameIdPaths:    paths,
		FullJsonSurvey: &uglyJson,
	}, nil
}
