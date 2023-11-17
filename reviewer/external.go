package reviewer

import (
	"github.com/rendis/surveygo/v2/question/types"
)

// ReviewExternal validates the answers for the given external type.
func ReviewExternal(_ any, _ []any, _ types.QuestionType) error {
	return nil
}
