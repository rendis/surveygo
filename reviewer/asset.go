package reviewer

import (
	"github.com/rendis/surveygo/v2/question/types"
)

// ReviewAsset validates the answers for an asset type.
func ReviewAsset(_ any, _ []any, _ types.QuestionType) error {
	return nil
}
