package text

import (
	"github.com/rendis/surveygo/v2/question/types"
)

// IdentificationNumber represents an identification number question type.
// QuestionType: types.QTypeIdentificationNumber
type IdentificationNumber struct {
	types.QBase `bson:",inline"`
}
