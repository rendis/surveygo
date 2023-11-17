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
var textAnswerReviewers = map[types.QuestionType]func(question any, answer string) error{
	types.QTypeTextArea:    reviewFreeText,
	types.QTypeInputText:   reviewFreeText,
	types.QTypeEmail:       reviewEmail,
	types.QTypeTelephone:   reviewTelephone,
	types.QTypeInformation: dummyReview,
}

// ReviewText validates format of the answers for the given text type.
func ReviewText(questionValue any, answers []string, qt types.QuestionType) error {
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
func reviewFreeText(questionValue any, answer string) error {
	freeText, _ := text.CastToFreeText(questionValue)
	l := len(answer)

	if freeText.Min != nil && l < *freeText.Min {
		return fmt.Errorf("answer length is less than min length '%d'. got: '%s' (%d)", *freeText.Min, answer, l)
	}

	if freeText.Max != nil && l > *freeText.Max {
		return fmt.Errorf("answer length is greater than max length '%d'. got: '%s' (%d)", *freeText.Max, answer, l)
	}

	return nil
}

// reviewEmail validates the answers for an email type.
func reviewEmail(questionValue any, answer string) error {
	email, _ := text.CastToEmail(questionValue)
	answer = strings.TrimSpace(answer)

	// validate email format
	if !emailRegex.MatchString(answer) {
		return fmt.Errorf("answer is not a valid email. got: '%s'", answer)
	}

	// validate domain
	if email.AllowedDomains != nil && len(email.AllowedDomains) > 0 {
		for _, allowedDomain := range email.AllowedDomains {
			if strings.HasSuffix(answer, allowedDomain) {
				return nil
			}
		}
		return fmt.Errorf("answer domain is not allowed. got '%s'", answer)
	}

	return nil
}

// reviewTelephone validates the answers for a telephone type.
func reviewTelephone(questionValue any, answer string) error {
	phone, _ := text.CastToTelephone(questionValue)
	answer = strings.ReplaceAll(answer, " ", "")
	answer = strings.ReplaceAll(answer, "-", "")

	// validate country code
	if phone.AllowedCountryCodes != nil && len(phone.AllowedCountryCodes) > 0 {
		for _, allowedCountryCode := range phone.AllowedCountryCodes {
			if strings.HasPrefix(answer, allowedCountryCode) {
				return nil
			}
		}
		return fmt.Errorf("answer country code is not allowed. got '%s'", answer)
	}

	return nil
}

func dummyReview(_ any, _ string) error {
	return nil
}
