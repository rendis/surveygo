package reviewer

import (
	"errors"
	"fmt"
	"github.com/rendis/surveygo/v2/question/types"
	"github.com/rendis/surveygo/v2/question/types/choice"
)

// choiceAnswerReviewers is a map of choice type to its review function.
var choiceAnswerReviewers = map[types.QuestionType]func(question any, answers []any) error{
	types.QTypeCheckbox:       reviewCheckbox,
	types.QTypeSingleSelect:   reviewSingleSelect,
	types.QTypeMultipleSelect: reviewMultipleSelect,
	types.QTypeRadio:          reviewRadio,
	types.QTypeToggle:         reviewToggle,
	types.QTypeSlider:         reviewSlider,
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
	return validator(questionValue, answers)
}

// reviewCheckbox validates the answers for a checkbox type.
func reviewCheckbox(questionValue any, answers []any) error {
	q, err := choice.CastToChoice(questionValue)
	if err != nil {
		return err
	}
	return choiceContainsAllAnswers(q, answers)
}

// reviewSingleSelect validates the answers for a single select type.
func reviewSingleSelect(questionValue any, answers []any) error {
	if len(answers) > 1 {
		return fmt.Errorf("single select can only have one answer. got: %v", answers)
	}

	q, err := choice.CastToChoice(questionValue)
	if err != nil {
		return err
	}

	return choiceContainsAllAnswers(q, answers)
}

// reviewMultipleSelect validates the answers for a multiple select type.
func reviewMultipleSelect(questionValue any, answers []any) error {
	q, err := choice.CastToChoice(questionValue)
	if err != nil {
		return err
	}
	return choiceContainsAllAnswers(q, answers)
}

// reviewRadio validates the answers for a radio type.
func reviewRadio(questionValue any, answers []any) error {
	// only one answer is allowed for a radio type
	if len(answers) > 1 {
		return fmt.Errorf("radio can only have one answer. got: %v", answers)
	}

	// answer len must be greater than 0
	if len(answers) == 0 {
		return fmt.Errorf("radio must have one answer. got: %v", answers)
	}

	q, err := choice.CastToChoice(questionValue)
	if err != nil {
		return err
	}

	// cast answer to string
	answer, ok := answers[0].(string)
	if !ok {
		return fmt.Errorf("invalid answer type. expected string, got %T", answers[0])
	}

	// answer must be in the options
	if _, ok = q.GetOptionsGroups()[answer]; !ok {
		return fmt.Errorf("answer '%s' not found in options", answers[0])
	}

	return choiceContainsAllAnswers(q, answers)
}

// reviewToggle validates the answers for a toggle type.
func reviewToggle(_ any, answers []any) error {
	// only one answer is allowed for a toggle type
	if len(answers) > 1 {
		return fmt.Errorf("toggle can only have one answer. got: %v", answers)
	}

	// answer len must be greater than 0
	if len(answers) == 0 {
		return fmt.Errorf("toggle must have one answer. got: %v", answers)
	}

	// cast answer to bool
	_, ok := answers[0].(bool)
	if !ok {
		return fmt.Errorf("invalid answer type for toggle. expected bool, got %T", answers[0])
	}

	return nil
}

// reviewSlider validates the answers for a slider type.
func reviewSlider(questionValue any, answers []any) error {
	// only one answer is allowed for a slider type
	if len(answers) > 1 {
		return fmt.Errorf("slider can only have one answer. got: %v", answers)
	}

	// answer len must be greater than 0
	if len(answers) == 0 {
		return fmt.Errorf("slider must have one answer. got: %v", answers)
	}

	q, err := choice.CastToSlider(questionValue)
	if err != nil {
		return err
	}

	// cast answer to int
	answer, ok := answers[0].(int)
	if !ok {
		return fmt.Errorf("invalid answer type. expected int, got %T", answers[0])
	}

	// answer must be in the range
	if answer < q.Min || answer > q.Max {
		return fmt.Errorf("answer '%d' is not in the range [%d, %d]", answer, q.Min, q.Max)
	}

	return nil
}

// choiceContainsAllAnswers validates that the given answers are contained in the choice options.
func choiceContainsAllAnswers(questionValue *choice.Choice, answers []any) error {
	var errs []error

	var optionsIDs = make(map[string]bool)
	for _, option := range questionValue.Options {
		optionsIDs[option.NameId] = true
	}

	for _, answer := range answers {
		// cast answer to string
		a, ok := answer.(string)
		if !ok {
			errs = append(errs, fmt.Errorf("invalid answer type. expected string, got %T", answer))
			continue
		}

		if !optionsIDs[a] {
			errs = append(errs, fmt.Errorf("answer '%s' not found in options", answer))
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
