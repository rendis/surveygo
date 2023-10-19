package reviewer

import (
	"fmt"
	"github.com/rendis/surveygo/v2/question/types"
)

// QuestionReviewer defines the function signature for a question validator.
type QuestionReviewer func(question any, answers []any, qt types.QuestionType) error

// GetQuestionReviewer returns the QuestionReviewer for the given question type.
func GetQuestionReviewer(qt types.QuestionType) (QuestionReviewer, error) {
	switch {
	case types.IsChoiceType(qt):
		return ReviewChoice, nil
	case types.IsTextType(qt):
		return ReviewText, nil
	case types.IsExternalType(qt):
		return ReviewExternal, nil
	}
	return nil, fmt.Errorf("unknown question type: %s", qt)
}
