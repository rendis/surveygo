package surveygo

import (
	"encoding/json"
	"fmt"
	"github.com/rendis/surveygo/part"
)

func (s *Survey) RemoveQuestion(nameId string) error {
	var ins internalSurvey
	err := json.Unmarshal([]byte(s.JsonSurvey), &ins)
	if err != nil {
		return err
	}

	var found bool
	for i, q := range ins.Questions {
		if *q.NameId == nameId {
			ins.Questions = append(ins.Questions[:i], ins.Questions[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("question '%s' not found", nameId)
	}

	return s.internalUpdate(&ins)
}

func (s *Survey) AddQuestion(questionJson string) error {
	return s.AddQuestionBytes([]byte(questionJson))
}

func (s *Survey) AddQuestionBytes(questionByte []byte) error {
	var question part.Question
	err := json.Unmarshal(questionByte, &question)
	if err != nil {
		return err
	}

	if _, ok := s.NameIdPaths[*question.NameId]; ok {
		return fmt.Errorf("question '%s' already exists", *question.NameId)
	}

	var ins internalSurvey
	err = json.Unmarshal([]byte(s.JsonSurvey), &ins)
	if err != nil {
		return err
	}

	ins.Questions = append(ins.Questions, question)
	return s.internalUpdate(&ins)
}

func (s *Survey) UpdateQuestion(questionJson string) error {
	return s.UpdateQuestionBytes([]byte(questionJson))
}

func (s *Survey) UpdateQuestionBytes(questionByte []byte) error {
	var question part.Question
	err := json.Unmarshal(questionByte, &question)
	if err != nil {
		return err
	}

	var ins internalSurvey
	err = json.Unmarshal([]byte(s.JsonSurvey), &ins)
	if err != nil {
		return err
	}

	var found bool
	for i, q := range ins.Questions {
		if *q.NameId == *question.NameId {
			ins.Questions[i] = question
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("question '%s' not found", *question.NameId)
	}

	return s.internalUpdate(&ins)
}

func (s *Survey) internalUpdate(ins *internalSurvey) error {
	paths, err := ins.getNameIdPaths()
	if err != nil {
		return err
	}

	b, err := json.Marshal(ins)
	if err != nil {
		return err
	}
	s.JsonSurvey = string(b)
	s.NameIdPaths = paths
	return nil
}
