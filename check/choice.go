package check

import (
	"fmt"
	"github.com/rendis/surveygo/part"
	"github.com/tidwall/gjson"
)

// ValidateChoice validates the answers for the given choice type.
// obj: the JSON representation of the survey.
// answers: the list of answers to validate.
// qt: the type of the choice field.
func ValidateChoice(obj gjson.Result, answers []any, qt part.QuestionType) error {
	if len(answers) == 0 {
		return nil
	}

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

// validateCheckbox validates the answers for a checkbox type.
func validateCheckbox(obj gjson.Result, answers []any) error {
	return validateChoiceContains(obj, answers)
}

// validateSingleSelect validates the answers for a single select type.
func validateSingleSelect(obj gjson.Result, answers []any) error {
	if len(answers) > 1 {
		return fmt.Errorf("single select can only have one answer. got: %v", answers)
	}
	return validateChoiceContains(obj, answers)
}

// validateMultipleSelect validates the answers for a multiple select type.
func validateMultipleSelect(obj gjson.Result, answers []any) error {
	return validateChoiceContains(obj, answers)
}

// validateRadio validates the answers for a radio type.
func validateRadio(obj gjson.Result, answers []any) error {
	if len(answers) > 1 {
		return fmt.Errorf("radio can only have one answer. got: %v", answers)
	}
	return validateChoiceContains(obj, answers)
}

// validateChoiceContains checks if the options in the choice field contain all the answers.
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

// contains returns true if the given element 'e' is present in the slice 's'.
func contains(s []any, e any) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
