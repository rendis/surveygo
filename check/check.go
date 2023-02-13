package check

import (
	"fmt"
	"github.com/rendis/surveygo/part"
	"github.com/tidwall/gjson"
)

// QuestionChecker is an alias for a function that validates the answers for a question.
type QuestionChecker func(obj gjson.Result, answers []any, qt part.QuestionType) error

// GetQuestionChecker returns the QuestionChecker for the given question type.
func GetQuestionChecker(qt part.QuestionType) (QuestionChecker, error) {
	switch {
	case part.IsChoiceType(qt):
		return ValidateChoice, nil
	case part.IsTextType(qt):
		return ValidateText, nil
	}
	return nil, fmt.Errorf("unknown question type: %s", qt)
}
