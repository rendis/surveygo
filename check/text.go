package check

import (
	"fmt"
	"github.com/rendis/surveygo/part"
	"github.com/tidwall/gjson"
	"regexp"
	"strings"
)

// emailRegex is a regex to validate email.
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z]{2,})+$")

// phoneRegex is a regex to validate phone number.
var phoneRegex = regexp.MustCompile("^(\\+[1-9]\\d{1,2})\\d{8,15}$")

// textAnswerValidator is a map of text type to its validator function.
var textAnswerValidator = map[part.QuestionType]func(obj gjson.Result, answer string) error{
	part.QTypeTextArea:  validateText,
	part.QTypeInputText: validateText,
	part.QTypeEmail:     validateEmail,
	part.QTypeTelephone: validateTelephone,
}

// ValidateText validates format of the answers for the given text type.
func ValidateText(obj gjson.Result, answers []any, qt part.QuestionType) error {
	if len(answers) != 1 {
		return fmt.Errorf("text type can only have one answer. got: %v", answers)
	}

	validator, ok := textAnswerValidator[qt]
	if !ok {
		return fmt.Errorf("invalid text type '%s'. supported types: %v", qt, part.QTypeChoiceTypes)
	}
	return validator(obj, answers[0].(string))
}

// validateText validates the answers for a text type.
func validateText(obj gjson.Result, answer string) error {
	return validateTextLength(obj, answer)
}

// validateEmail validates the answers for an email type.
func validateEmail(obj gjson.Result, answer string) error {
	if err := validateTextLength(obj, answer); err != nil {
		return err
	}

	// trim spaces
	answer = strings.TrimSpace(answer)

	if !emailRegex.MatchString(answer) {
		return fmt.Errorf("answer is not a valid email. got: '%s'", answer)
	}

	if obj.Get("value.allowedDomains").Exists() {
		allowedDomains := obj.Get("value.allowedDomains").Array()
		for _, allowedDomain := range allowedDomains {
			if strings.HasSuffix(answer, allowedDomain.String()) {
				return nil
			}
		}

		return fmt.Errorf("answer domain is not allowed. got '%s', allowed domains: %v", answer, allowedDomains)
	}

	return nil
}

// validateTelephone validates the answers for a telephone type.
func validateTelephone(obj gjson.Result, answer string) error {
	if err := validateTextLength(obj, answer); err != nil {
		return err
	}

	// remove spaces and dashes
	answer = strings.ReplaceAll(answer, " ", "")
	answer = strings.ReplaceAll(answer, "-", "")

	if !phoneRegex.MatchString(answer) {
		return fmt.Errorf("answer is not a valid telephone. got: '%s'", answer)
	}

	if obj.Get("value.allowedCountryCodes").Exists() {
		allowedCountryCodes := obj.Get("value.allowedCountryCodes").Array()
		for _, allowedCountryCode := range allowedCountryCodes {
			if strings.HasPrefix(answer, allowedCountryCode.String()) {
				return nil
			}
		}

		return fmt.Errorf("answer country code is not allowed. got '%s', allowed country codes: %v", answer, allowedCountryCodes)
	}

	return nil
}

// validateTextLength validates the length of the answers for the given text type.
func validateTextLength(obj gjson.Result, answer string) error {
	// Validate min length.
	if obj.Get("value.min").Exists() {
		min := obj.Get("value.min").Int()
		if int64(len(answer)) < min {
			return fmt.Errorf("answer length is less than min length '%d'. got: '%s' (%d)", min, answer, len(answer))
		}
	}

	// Validate max length.
	if obj.Get("value.max").Exists() {
		max := obj.Get("value.max").Int()
		if int64(len(answer)) > max {
			return fmt.Errorf("answer length is greater than max length '%d'. got: '%s' (%d)", max, answer, len(answer))
		}
	}

	return nil
}
