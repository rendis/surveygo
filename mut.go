package surveygo

import (
	"encoding/json"
	"fmt"
	"github.com/rendis/surveygo/part"
)

func (s *Survey) RemoveQuestion(nameId string) error {
	ins, err := s.getInternal()
	if err != nil {
		return err
	}
	if ins == nil {
		return fmt.Errorf("json survey is empty, title: '%s', version: '%s'", *s.Title, *s.Version)
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

	return s.internalUpdate(ins)
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

	ins, err := s.getInternal()
	if err != nil {
		return err
	}
	if ins == nil {
		return fmt.Errorf("json survey is empty, title: '%s', version: '%s'", *s.Title, *s.Version)
	}

	ins.Questions = append(ins.Questions, question)
	return s.internalUpdate(ins)
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

	ins, err := s.getInternal()
	if err != nil {
		return err
	}
	if ins == nil {
		return fmt.Errorf("json survey is empty, title: '%s', version: '%s'", *s.Title, *s.Version)
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

	return s.internalUpdate(ins)
}

func (s *Survey) internalUpdate(ins *jsonSurvey) error {
	paths, err := ins.getNameIdPaths()
	if err != nil {
		return err
	}

	b, err := json.Marshal(ins)
	if err != nil {
		return err
	}

	js := string(b)
	s.FullJsonSurvey = &js
	s.NameIdPaths = paths
	return nil
}

func (s *Survey) getInternal() (*jsonSurvey, error) {
	if s.FullJsonSurvey == nil {
		return &jsonSurvey{
			Title:       s.Title,
			Version:     s.Version,
			Description: s.Description,
			Questions:   []part.Question{},
		}, nil
	}

	var ins jsonSurvey
	err := json.Unmarshal([]byte(*s.FullJsonSurvey), &ins)
	if err != nil {
		return nil, err
	}

	return &ins, nil
}
