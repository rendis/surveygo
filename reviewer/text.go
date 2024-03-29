package reviewer

import (
	"fmt"
	"github.com/rendis/surveygo/v2/question/types"
	"github.com/rendis/surveygo/v2/question/types/text"
	"regexp"
	"strings"
	"time"
)

// emailRegex is a regex to validate email.
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z]{2,})+$")

// textAnswerReviewers is a map of text type to its review function.
var textAnswerReviewers = map[types.QuestionType]func(question any, answer []any) error{
	types.QTypeTextArea:             reviewFreeText,
	types.QTypeInputText:            reviewFreeText,
	types.QTypeEmail:                reviewEmail,
	types.QTypeTelephone:            reviewTelephone,
	types.QTypeDateTime:             reviewDateTime,
	types.QTypeInformation:          dummyReview,
	types.QTypeIdentificationNumber: dummyReview,
}

// ReviewText validates format of the answers for the given text type.
func ReviewText(questionValue any, answers []any, qt types.QuestionType) error {
	validator, ok := textAnswerReviewers[qt]
	if !ok {
		return fmt.Errorf("invalid text type '%s'. supported types: %v", qt, types.QTypeChoiceTypes)
	}
	return validator(questionValue, answers)
}

// reviewFreeText validates the answers for a text type.
func reviewFreeText(questionValue any, answers []any) error {
	if len(answers) != 1 {
		return fmt.Errorf("text type can only have one answer. got: %v", answers)
	}

	answer := answers[0]
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
func reviewEmail(questionValue any, answers []any) error {
	if len(answers) != 1 {
		return fmt.Errorf("text type can only have one answer. got: %v", answers)
	}

	answer := answers[0]
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
func reviewTelephone(questionValue any, answers []any) error {
	if len(answers) == 0 || len(answers) > 2 {
		return fmt.Errorf("text type can only have one [phone number] or two [country code, phone number] answers. got: %v", answers)
	}

	phone, _ := text.CastToTelephone(questionValue)

	if len(phone.AllowedCountryCodes) == 0 {
		return nil
	}

	countryCodeAnswer, ok := answers[0].(string)
	if !ok {
		// try to cast to int
		countryCodeAnswerInt, ok := answers[0].(int)
		if !ok {
			return fmt.Errorf("answer country code must be a string or an int. got: %v", answers[0])
		}
		countryCodeAnswer = fmt.Sprintf("%d", countryCodeAnswerInt)
	}

	// validate allowed country codes
	for _, allowedCountryCode := range phone.AllowedCountryCodes {
		if strings.HasPrefix(countryCodeAnswer, allowedCountryCode) {
			return nil
		}
	}

	return fmt.Errorf("answer country code is not allowed. got '%s'", countryCodeAnswer)
}

// reviewDateTime validates the answers for a date time type.
func reviewDateTime(questionValue any, answers []any) error {
	if len(answers) != 1 {
		return fmt.Errorf("date time type can only have one answer. got: %v", answers)
	}

	answer, ok := answers[0].(string)
	if !ok {
		return fmt.Errorf("date time answer must be a string. got: %v", answers[0])
	}

	dateTime, _ := text.CastToDateTime(questionValue)

	dateTimeFormat := dateTime.Format

	// check if answer has the correct format
	if _, err := time.Parse(dateTimeFormat, answer); err != nil {
		return fmt.Errorf("answer is not a valid date time format '%s'. got: '%s'", dateTimeFormat, answer)
	}

	return nil
}

func dummyReview(_ any, _ []any) error {
	return nil
}
