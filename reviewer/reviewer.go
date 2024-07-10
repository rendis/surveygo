package reviewer

import (
	"fmt"
	"github.com/rendis/surveygo/v2/question/types"
)

// QuestionReviewer defines the function signature for a question validator.
type QuestionReviewer func(question any, answers []any, qt types.QuestionType) error

// GroupAnswers is a map with the answers provided by the user.
// Each item is a group of answers for the different questions in the group.
// The key is the question NameId (Question.NameId).
type GroupAnswers []map[string][]any

// GetQuestionReviewer returns the QuestionReviewer for the given question type.
func GetQuestionReviewer(qt types.QuestionType) (QuestionReviewer, error) {
	switch {
	case types.IsChoiceType(qt):
		return ReviewChoice, nil
	case types.IsTextType(qt):
		return ReviewText, nil
	case types.IsExternalType(qt):
		return ReviewExternal, nil
	case types.IsAssetType(qt):
		return ReviewAsset, nil
	}
	return nil, fmt.Errorf("unknown question type: %s", qt)
}

// ExtractGroupNestedAnswers extracts the nested answers
func ExtractGroupNestedAnswers(groupAnswersPack []any) (GroupAnswers, error) {
	if len(groupAnswersPack) == 0 {
		return nil, nil
	}

	var resp = make(GroupAnswers, 0)

	for _, p := range groupAnswersPack {
		answersPack, ok := p.([]map[string][]any)
		if !ok {
			return nil, fmt.Errorf(
				"invalid answersPack, expected '[]map[string][]any', got '%T' at index '%d'", p, len(resp),
			)
		}

		for _, answers := range answersPack {
			resp = append(resp, answers)
		}
	}

	return resp, nil
}
