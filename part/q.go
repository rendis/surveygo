package part

import (
	"encoding/json"
	"fmt"
	"regexp"
)

var nameIdRegexp = regexp.MustCompile(`^[a-zA-Z][a-zA-Z\d_-]{1,62}[a-zA-Z\d]$`)

type baseQuestion struct {
	Order  *int          `json:"order"`
	NameId *string       `json:"nameId"`
	QTyp   *QuestionType `json:"type"`
	Label  *string       `json:"label"`
}

type NameIdPath struct {
	NameId string
	Path   []string
}

type Question struct {
	baseQuestion
	Value any `json:"value"`
}

func (q *Question) UnmarshalJSON(b []byte) error {
	var bq baseQuestion
	if err := json.Unmarshal(b, &bq); err != nil {
		return err
	}

	nameId := *bq.NameId

	if !nameIdRegexp.MatchString(nameId) {
		return fmt.Errorf("invalid nameId '%s', must match %s", nameId, nameIdRegexp.String())
	}

	var nq *Question
	var err error

	switch *bq.QTyp {
	case QTypeSingleSelect, QTypeMultipleSelect, QTypeRadio, QTypeCheckbox:
		nq, err = getQuestionByValueTyp[Choice](b, ChoiceUnmarshallValidator)
	case QTypeTextArea:
		nq, err = getQuestionByValueTyp[TextArea](b, TextAreaUnmarshallValidator)
	default:
		return fmt.Errorf("invalid question type: %s", *bq.QTyp)
	}

	if err != nil {
		return fmt.Errorf("\n - error unmarshalling question '%s'. %s", *bq.NameId, err)
	}
	*q = *nq
	return nil
}

func (q *Question) GetNameIdPaths(from []string) []NameIdPath {
	var paths = []NameIdPath{
		{
			NameId: *q.NameId,
			Path:   from,
		},
	}
	switch *q.QTyp {
	case QTypeSingleSelect, QTypeMultipleSelect, QTypeRadio, QTypeCheckbox:
		currentPath := append(from, "value", "options")
		for oi, o := range q.Value.(*Choice).Options {
			if o.SubQuestions != nil && len(o.SubQuestions) > 0 {
				for si, nq := range o.SubQuestions {
					p := nq.GetNameIdPaths(append(currentPath, fmt.Sprintf("%d", oi), "subQuestions", fmt.Sprintf("%d", si)))
					paths = append(paths, p...)
				}
			}
		}
	}

	return paths
}

func getQuestionByValueTyp[T any](b []byte, unmarshallValidator func(*T) error) (*Question, error) {
	var tq = struct {
		baseQuestion
		Value *T `json:"value"`
	}{}

	if err := json.Unmarshal(b, &tq); err != nil {
		return nil, err
	}

	if tq.Value == nil {
		return nil, fmt.Errorf("value is not defined")
	}

	if err := unmarshallValidator(tq.Value); err != nil {
		return nil, err
	}

	return &Question{
		baseQuestion: tq.baseQuestion,
		Value:        tq.Value,
	}, nil
}
