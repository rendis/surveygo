package text

import (
	"encoding/json"
	"fmt"
	"github.com/rendis/surveygo/v2/question/types"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// DateTime represents a date time question type.
// QuestionType: types.QTypeDateTime
type DateTime struct {
	types.QBase `bson:",inline"`

	// Format is the format of the date time field.
	// Validations:
	// - required
	Format string `json:"format" bson:"format" validate:"required"`

	// Type is the type of the date time field.
	// Validations:
	// - required
	Type DateTypeFormat `json:"type" bson:"type" validate:"required"`
}

// CastToDateTime casts the given interface to a DateTime type.
func CastToDateTime(questionValue any) (*DateTime, error) {
	c, ok := questionValue.(*DateTime)
	if !ok {
		return nil, fmt.Errorf("invalid type, expected *text.DateTime, got %T", questionValue)
	}
	return c, nil
}

type DateTypeFormat string

const (
	DateTypeFormatDate     DateTypeFormat = "date"
	DateTypeFormatTime     DateTypeFormat = "time"
	DateTypeFormatDateTime DateTypeFormat = "datetime"
)

var dateTypeMap = map[string]DateTypeFormat{
	"date":     DateTypeFormatDate,
	"time":     DateTypeFormatTime,
	"datetime": DateTypeFormatDateTime,
}

func (s *DateTypeFormat) UnmarshalJSON(b []byte) error {
	var st string
	if err := json.Unmarshal(b, &st); err != nil {
		return fmt.Errorf("unmarshal error, %s", err)
	}

	t, ok := dateTypeMap[st]
	if !ok {
		return fmt.Errorf("invalid date type format '%s'", st)
	}

	*s = t

	return nil
}

func (s *DateTypeFormat) UnmarshalBSONValue(typ bsontype.Type, raw []byte) error {
	if typ != bsontype.String {
		return fmt.Errorf("invalid bson value type '%s'", typ.String())
	}

	c, _, ok := bsoncore.ReadString(raw)
	if !ok {
		return fmt.Errorf("invalid bson value '%s'", string(raw))
	}

	t, ok := dateTypeMap[c]
	if !ok {
		return fmt.Errorf("invalid date type format '%s'", c)
	}

	*s = t
	return nil
}
