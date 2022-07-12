package surveygo

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

func Parse(j string) (*Survey, error) {
	b := []byte(j)
	return ParseBytes(b)
}

func ParseBytes(b []byte) (*Survey, error) {
	uglyJson := gjson.ParseBytes(b).Get("@ugly").String()

	var ins internalSurvey
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
		Title:       *ins.Title,
		Version:     *ins.Version,
		Description: ins.Description,
		NameIdPaths: paths,
		JsonSurvey:  &uglyJson,
	}, nil
}
