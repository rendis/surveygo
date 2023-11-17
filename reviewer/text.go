package reviewer

import (
	"fmt"
	"github.com/rendis/surveygo/v2/question/types"
	"github.com/rendis/surveygo/v2/question/types/text"
	"regexp"
	"strings"
)

// emailRegex is a regex to validate email.
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z]{2,})+$")

// textAnswerReviewers is a map of text type to its review function.
var textAnswerReviewers = map[types.QuestionType]func(question any, answer any) error{
	types.QTypeTextArea:    reviewFreeText,
	types.QTypeInputText:   reviewFreeText,
	types.QTypeEmail:       reviewEmail,
	types.QTypeTelephone:   reviewTelephone,
	types.QTypeInformation: dummyReview,
}

// ReviewText validates format of the answers for the given text type.
func ReviewText(questionValue any, answers []any, qt types.QuestionType) error {
	if len(answers) != 1 {
		return fmt.Errorf("text type can only have one answer. got: %v", answers)
	}

	validator, ok := textAnswerReviewers[qt]
	if !ok {
		return fmt.Errorf("invalid text type '%s'. supported types: %v", qt, types.QTypeChoiceTypes)
	}
	return validator(questionValue, answers[0])
}

// reviewFreeText validates the answers for a text type.
func reviewFreeText(questionValue any, answer any) error {
	freeText, _ := text.CastToFreeText(questionValue)

	// cast answer to string
	a, ok := answer.(string)
	if !ok {
		return fmt.Errorf("answer is not a string. got: %v", answer)
	}

	l := len(a)

	if freeText.Min != nil && l < *freeText.Min {
		return fmt.Errorf("answer length is less than min length '%d'. got: '%s' (%d)", *freeText.Min, a, l)
	}

	if freeText.Max != nil && l > *freeText.Max {
		return fmt.Errorf("answer length is greater than max length '%d'. got: '%s' (%d)", *freeText.Max, a, l)
	}

	return nil
}

// reviewEmail validates the answers for an email type.
func reviewEmail(questionValue any, answer any) error {
	email, _ := text.CastToEmail(questionValue)

	// cast answer to string
	a, ok := answer.(string)
	if !ok {
		return fmt.Errorf("answer is not a string. got: %v", answer)
	}

	a = strings.TrimSpace(a)

	// validate email format
	if !emailRegex.MatchString(a) {
		return fmt.Errorf("answer is not a valid email. got: '%s'", a)
	}

	// validate domain
	if email.AllowedDomains != nil && len(email.AllowedDomains) > 0 {
		for _, allowedDomain := range email.AllowedDomains {
			if strings.HasSuffix(a, allowedDomain) {
				return nil
			}
		}
		return fmt.Errorf("answer domain is not allowed. got '%s'", a)
	}

	return nil
}

// reviewTelephone validates the answers for a telephone type.
func reviewTelephone(questionValue any, answer any) error {
	phone, _ := text.CastToTelephone(questionValue)

	// cast answer to string
	a, ok := answer.(string)
	if !ok {
		return fmt.Errorf("answer is not a string. got: %v", answer)
	}

	a = strings.ReplaceAll(a, " ", "")
	a = strings.ReplaceAll(a, "-", "")

	// validate country code
	if phone.AllowedCountryCodes != nil && len(phone.AllowedCountryCodes) > 0 {
		for _, allowedCountryCode := range phone.AllowedCountryCodes {
			if strings.HasPrefix(a, allowedCountryCode) {
				return nil
			}
		}
		return fmt.Errorf("answer country code is not allowed. got '%s'", a)
	}

	return nil
}

func dummyReview(_ any, _ any) error {
	return nil
}
