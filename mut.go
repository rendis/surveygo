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

func (s *Survey) AddQuestionMap(question map[string]any) error {
	b, _ := json.Marshal(question)
	return s.AddQuestionBytes(b)
}

func (s *Survey) AddQuestionJson(question string) error {
	return s.AddQuestionBytes([]byte(question))
}

func (s *Survey) AddQuestionBytes(question []byte) error {
	var pq part.Question
	err := json.Unmarshal(question, &pq)
	if err != nil {
		return err
	}

	if _, ok := s.NameIdPaths[*pq.NameId]; ok {
		return fmt.Errorf("pq '%s' already exists", *pq.NameId)
	}

	ins, err := s.getInternal()
	if err != nil {
		return err
	}
	if ins == nil {
		return fmt.Errorf("json survey is empty, title: '%s', version: '%s'", *s.Title, *s.Version)
	}

	ins.Questions = append(ins.Questions, pq)
	return s.internalUpdate(ins)
}

func (s *Survey) UpdateQuestionMap(question map[string]any) error {
	b, _ := json.Marshal(question)
	return s.UpdateQuestionBytes(b)
}

func (s *Survey) UpdateQuestionJson(question string) error {
	return s.UpdateQuestionBytes([]byte(question))
}

func (s *Survey) UpdateQuestionBytes(question []byte) error {
	var pq part.Question
	err := json.Unmarshal(question, &pq)
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
		if *q.NameId == *pq.NameId {
			ins.Questions[i] = pq
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("pq '%s' not found", *pq.NameId)
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
