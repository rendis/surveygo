package part

import (
	"encoding/json"
	"fmt"
)

type QuestionType string

const (
	QTypeSingleSelect   QuestionType = "single_select"
	QTypeMultipleSelect QuestionType = "multi_select"
	QTypeRadio          QuestionType = "radio"
	QTypeCheckbox       QuestionType = "checkbox"
	QTypeTextArea       QuestionType = "text_area"
)

func (s *QuestionType) UnmarshalJSON(b []byte) error {
	var st string
	if err := json.Unmarshal(b, &st); err != nil {
		return err
	}

	t, err := ParseToQuestionType(st)
	if err != nil {
		return fmt.Errorf("unmarshal error, %s", err)
	}
	*s = t
	return nil
}

func (s QuestionType) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(s))
}

func ParseToQuestionType(v string) (QuestionType, error) {
	switch v {
	case string(QTypeSingleSelect):
		return QTypeSingleSelect, nil
	case string(QTypeMultipleSelect):
		return QTypeMultipleSelect, nil
	case string(QTypeRadio):
		return QTypeRadio, nil
	case string(QTypeCheckbox):
		return QTypeCheckbox, nil
	case string(QTypeTextArea):
		return QTypeTextArea, nil
	default:
		return "", fmt.Errorf("invalid question type: %s", v)
	}
}
