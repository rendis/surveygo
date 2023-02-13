package part

import (
	"encoding/json"
	"fmt"
	"regexp"
)

// nameIdRegex is a regular expression used to validate the format of the "nameId" field in a Question.
var nameIdRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z\d_-]{1,62}[a-zA-Z\d]$`)

// NameIdPath represents a path to a question in a survey, including its NameId.
type NameIdPath struct {
	// NameId is the identifier of the question.
	NameId string

	// Path is the location of the question within the survey.
	Path []string

	// Required indicates whether the question is required.
	Required bool
}

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

// Question is a struct that represents a question in a survey.
type Question struct {
	// baseQuestion contains common fields for all types of questions.
	baseQuestion

	// Value is the value of the question, which can be of different types depending on the type of question.
	Value any `json:"value"`
}

// QuestionPath represents a Question with its associated path.
type QuestionPath struct {
	Question
	path []string
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

// GetNameIdPaths returns a slice of NameIdPath structs containing the NameId, Required and Path of a Question and its sub-questions.
// startPath represents the starting path of the Question.
func (q *Question) GetNameIdPaths(startPath []string) []*NameIdPath {
	var paths []*NameIdPath

	// Create a slice of QuestionPath structs containing the Question and its starting path.
	queue := []QuestionPath{{*q, startPath}}

	// Loop through the queue until it's empty.
	for len(queue) > 0 {
		// Get the first QuestionPath from the queue.
		currQ := queue[0]
		queue = queue[1:]

		// Add the current Question's NameId, Required and Path to the paths slice.
		paths = append(paths, &NameIdPath{
			NameId:   *currQ.NameId,
			Required: currQ.Required,
			Path:     currQ.path,
		})

		// If the current Question is not of type Choice, continue to the next iteration.
		if !IsChoiceType(*currQ.QTyp) {
			continue
		}

		// Create the current path for the options.
		currentPath := append(startPath, "value", "options")

		// Get a queue of the sub-questions for each option.
		optQueue := getOptionsQueue(currQ.Value.(*Choice).Options, currentPath)

		// Append the options queue to the main queue.
		queue = append(queue, optQueue...)
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

// getOptionsQueue returns a slice of QuestionPath structs containing the sub-questions of the options and their associated path.
func getOptionsQueue(options []Option, path []string) []QuestionPath {
	var queue []QuestionPath

	// Loop through each option.
	for optIndex, opt := range options {
		// If the option has sub-questions, get a queue of the sub-questions.
		if opt.SubQuestions != nil && len(opt.SubQuestions) > 0 {
			newPath := make([]string, len(path))
			copy(newPath, path)
			newPath = append(newPath, fmt.Sprintf("%d", optIndex), "subQuestions")
			subQueue := getSubQuestionQueue(opt.SubQuestions, newPath)
			queue = append(queue, subQueue...)
		}
	}

	return queue
}

// getSubQuestionQueue returns a slice of QuestionPath structs containing the sub-questions and their associated path.
func getSubQuestionQueue(subQuestions []Question, path []string) []QuestionPath {
	var queue []QuestionPath

	// Loop through each sub-question.
	for subIndex, subQ := range subQuestions {
		// Create a new path for the sub-question.
		newPath := make([]string, len(path))
		copy(newPath, path)
		newPath = append(newPath, fmt.Sprintf("%d", subIndex))

		// Add the sub-question and its path to the queue.
		queue = append(queue, QuestionPath{
			Question: subQ,
			path:     newPath,
		})
	}

	return queue
}
