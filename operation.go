package surveygo

import (
	"encoding/json"
	"fmt"
	"github.com/rendis/surveygo/part"
)

// RemoveQuestion removes the question with the specified name ID from the survey
func (s *Survey) RemoveQuestion(nameId string) error {
	// get the internal representation of the survey.
	ins, err := s.getInternal()
	if err != nil {
		return err
	}

	// check if the internal representation is empty.
	if ins == nil {
		return fmt.Errorf("json survey is empty, title: '%s', version: '%s'", *s.Title, *s.Version)
	}

	// iterate through the questions to find the one with the specified name ID.
	var found bool
	for i, q := range ins.Questions {
		if *q.NameId == nameId {
			// remove the question from the slice of questions.
			ins.Questions = append(ins.Questions[:i], ins.Questions[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("question '%s' not found", nameId)
	}

	// update the internal representation of the survey.
	return s.internalUpdate(ins)
}

// AddQuestionMap adds a question to the survey given its representation as a map[string]any
func (s *Survey) AddQuestionMap(question map[string]any) error {
	b, _ := json.Marshal(question)
	return s.AddQuestionBytes(b)
}

// AddQuestionJson adds a question to the survey given its representation as a JSON string
func (s *Survey) AddQuestionJson(question string) error {
	return s.AddQuestionBytes([]byte(question))
}

// AddQuestionBytes adds a question to the survey given its representation as a byte array
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

// UpdateQuestionMap updates an existing question in the survey with the data provided as a map.
func (s *Survey) UpdateQuestionMap(question map[string]any) error {
	b, _ := json.Marshal(question)
	return s.UpdateQuestionBytes(b)
}

// UpdateQuestionJson updates an existing question in the survey with the data provided as a JSON string.
func (s *Survey) UpdateQuestionJson(question string) error {
	return s.UpdateQuestionBytes([]byte(question))
}

// UpdateQuestionBytes updates a question in the survey given its representation as a byte array
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

// internalUpdate updates the survey after a modification to its internal representation
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

// getInternal returns the internal representation of the survey as a jsonSurvey struct
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
