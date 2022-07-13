package check

import (
	"fmt"
	"github.com/rendis/surveygo/part"
	"github.com/tidwall/gjson"
)

func ValidateChoice(obj gjson.Result, answers []any, qt part.QuestionType) error {
	switch qt {
	case part.QTypeCheckbox:
		return validateCheckbox(obj, answers)
	case part.QTypeSingleSelect:
		return validateSingleSelect(obj, answers)
	case part.QTypeMultipleSelect:
		return validateMultipleSelect(obj, answers)
	case part.QTypeRadio:
		return validateRadio(obj, answers)
	default:
		return fmt.Errorf("invalid choice type: %s", qt)
	}
}

func validateCheckbox(obj gjson.Result, answers []any) error {
	if len(answers) == 0 {
		return nil
	}
	return validateChoiceContains(obj, answers)
}

func validateSingleSelect(obj gjson.Result, answers []any) error {
	if len(answers) == 0 {
		return nil
	}

	if len(answers) > 1 {
		return fmt.Errorf("single select can only have one answer. got: %v", answers)
	}
	return validateChoiceContains(obj, answers)
}

func validateMultipleSelect(obj gjson.Result, answers []any) error {
	if len(answers) == 0 {
		return nil
	}
	return validateChoiceContains(obj, answers)
}

func validateRadio(obj gjson.Result, answers []any) error {
	if len(answers) == 0 {
		return nil
	}
	if len(answers) > 1 {
		return fmt.Errorf("radio can only have one answer. got: %v", answers)
	}
	return validateChoiceContains(obj, answers)
}

func validateChoiceContains(obj gjson.Result, answers []any) error {
	var options []any
	obj.Get("value.options").ForEach(func(key, value gjson.Result) bool {
		options = append(options, value.Get("id").String())
		return true
	})
	for _, o := range answers {
		if !contains(options, o) {
			return fmt.Errorf("option id not found: '%s'", o)
		}
	}
	return nil
}

func contains(s []any, e any) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
