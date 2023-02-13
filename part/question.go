package part

import (
	"encoding/json"
	"fmt"
	"regexp"
)

// nameIdRegex is a regular expression used to validate the format of the "nameId" field in a Question.
var nameIdRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z\d_-]{1,62}[a-zA-Z\d]$`)

// baseQuestion is a struct that contains common fields for all types of questions in a survey.
type baseQuestion struct {
	// Order is an optional order number for the question.
	Order *int `json:"order"`

	// NameId is a required identifier for the question.
	NameId *string `json:"nameId"`

	// QTyp is the type of question, such as single_select, multi_select, radio, checkbox, or text_area.
	QTyp *QuestionType `json:"type"`

	// Label is a required label for the question.
	Label *string `json:"label"`

	// Required is an optional boolean that indicates whether the question is required. Defaults to false.
	Required bool `json:"required"`
}

// NameIdPath represents a path to a question in a survey, including its NameId.
type NameIdPath struct {
	// NameId is the identifier of the question.
	NameId string

	// Path is the location of the question within the survey.
	Path []string
}

// Question is a struct that represents a question in a survey.
type Question struct {
	// baseQuestion contains common fields for all types of questions.
	baseQuestion

	// Value is the value of the question, which can be of different types depending on the type of question.
	Value any `json:"value"`
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (q *Question) UnmarshalJSON(b []byte) error {
	var bq baseQuestion
	if err := json.Unmarshal(b, &bq); err != nil {
		return err
	}

	nameId := *bq.NameId

	if !nameIdRegex.MatchString(nameId) {
		return fmt.Errorf("invalid nameId '%s', must match %s", nameId, nameIdRegex.String())
	}

	var nq *Question
	var err error

	// Get the correct type of question based on the type of question
	switch *bq.QTyp {
	case QTypeSingleSelect, QTypeMultipleSelect, QTypeRadio, QTypeCheckbox:
		nq, err = getQuestionByValueTyp[Choice](b, ChoiceUnmarshallValidator)
	case QTypeTextArea, QTypeInputText:
		nq, err = getQuestionByValueTyp[TextArea](b, TextAreaUnmarshallValidator)
	case QTypeEmail:
		nq, err = getQuestionByValueTyp[Email](b, EmailUnmarshallValidator)
	case QTypeTelephone:
		nq, err = getQuestionByValueTyp[Telephone](b, TelephoneUnmarshallValidator)
	default:
		return fmt.Errorf("invalid question type: %s", *bq.QTyp)
	}

	if err != nil {
		return fmt.Errorf("\n - error unmarshalling question '%s'. %s", *bq.NameId, err)
	}
	*q = *nq
	return nil
}

// GetNameIdPaths returns a list of NameIdPaths for a question and its sub-questions, if any.
func (q *Question) GetNameIdPaths(from []string) []NameIdPath {
	var paths = []NameIdPath{
		{
			NameId: *q.NameId,
			Path:   from,
		},
	}

	// If the question is a choice type, get the NameIdPaths for its sub-questions, if any.
	if IsChoiceType(*q.QTyp) {
		return append(paths, q.getChoiceNameIdPaths(from)...)
	}

	return paths
}

// getChoiceNameIdPaths is a helper function that returns a list of NameIdPaths for a choice type question.
func (q *Question) getChoiceNameIdPaths(from []string) []NameIdPath {
	paths := []NameIdPath{}
	currentPath := append(from, "value", "options")
	for oi, o := range q.Value.(*Choice).Options {
		if o.SubQuestions != nil && len(o.SubQuestions) > 0 {
			for si, nq := range o.SubQuestions {
				p := nq.GetNameIdPaths(append(currentPath, fmt.Sprintf("%d", oi), "subQuestions", fmt.Sprintf("%d", si)))
				paths = append(paths, p...)
			}
		}
	}
	return paths
}

// getQuestionByValueTyp is a helper function that returns a question of a specific type.
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
