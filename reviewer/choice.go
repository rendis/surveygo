package reviewer

import (
	"errors"
	"fmt"
	"github.com/rendis/surveygo/v2/question/types"
	"github.com/rendis/surveygo/v2/question/types/choice"
)

// choiceAnswerReviewers is a map of choice type to its review function.
var choiceAnswerReviewers = map[types.QuestionType]func(question *choice.Choice, answers []any) error{
	types.QTypeCheckbox:       reviewCheckbox,
	types.QTypeSingleSelect:   reviewSingleSelect,
	types.QTypeMultipleSelect: reviewMultipleSelect,
	types.QTypeRadio:          reviewRadio,
}

// ReviewChoice validates the answers for the given choice type.
// obj: the JSON representation of the survey.
// answers: the list of answers to validate.
// qt: the type of the choice field.
func ReviewChoice(questionValue any, answers []any, qt types.QuestionType) error {
	if len(answers) == 0 {
		return nil
	}
	validator, ok := choiceAnswerReviewers[qt]
	if !ok {
		return fmt.Errorf("invalid choice type '%s'. supported types: %v", qt, types.QTypeChoiceTypes)
	}
	c, _ := choice.CastToChoice(questionValue)
	return validator(c, answers)
}

// reviewCheckbox validates the answers for a checkbox type.
func reviewCheckbox(questionValue *choice.Choice, answers []any) error {
	return choiceContainsAllAnswers(questionValue, answers)
}

// reviewSingleSelect validates the answers for a single select type.
func reviewSingleSelect(questionValue *choice.Choice, answers []any) error {
	if len(answers) > 1 {
		return fmt.Errorf("single select can only have one answer. got: %v", answers)
	}
	return choiceContainsAllAnswers(questionValue, answers)
}

// reviewMultipleSelect validates the answers for a multiple select type.
func reviewMultipleSelect(questionValue *choice.Choice, answers []any) error {
	return choiceContainsAllAnswers(questionValue, answers)
}

// reviewRadio validates the answers for a radio type.
func reviewRadio(questionValue *choice.Choice, answers []any) error {
	if len(answers) > 1 {
		return fmt.Errorf("radio can only have one answer. got: %v", answers)
	}
	return choiceContainsAllAnswers(questionValue, answers)
}

// choiceContainsAllAnswers validates that the given answers are contained in the choice options.
func choiceContainsAllAnswers(questionValue *choice.Choice, answers []any) error {

	var errs []error

	options := questionValue.GetOptionsGroups()
	for _, answer := range answers {
		if _, ok := options[answer.(string)]; !ok {
			errs = append(errs, fmt.Errorf("answer '%s' not found in options", answer))
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
