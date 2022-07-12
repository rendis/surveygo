package surveygo

import (
	"encoding/json"
	"github.com/rendis/surveygo/part"
)

//func AddSubQuestion(survey *Survey, nameId, questionJson string) error {
//
//}

func AddQuestion(survey *Survey, questionJson string) error {
	return AddQuestionBytes(survey, []byte(questionJson))
}

func AddQuestionBytes(survey *Survey, questionByte []byte) error {
	var question part.Question
	err := json.Unmarshal(questionByte, &question)
	if err != nil {
		return err
	}

	var ins internalSurvey
	err = json.Unmarshal([]byte(survey.JsonSurvey), &ins)
	if err != nil {
		return err
	}
	ins.Questions = append(ins.Questions, question)

	paths, err := ins.getNameIdPaths()
	if err != nil {
		return err
	}

	b, err := json.Marshal(ins)
	if err != nil {
		return err
	}
	survey.JsonSurvey = string(b)
	survey.NameIdPaths = paths
	return nil
}
