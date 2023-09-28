package surveygo

import (
	"errors"
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/rendis/surveygo/v2/question/types"
	"regexp"
)

// questionNameIdRegexValidator is a regular expression used to validate the format of the "nameId" field in a Question.
var questionNameIdRegexValidator = regexp.MustCompile(`^[a-zA-Z][a-zA-Z\d_-]{1,62}[a-zA-Z\d]$`)

var validatorTranslator = newValidatorTranslator()
var SurveyValidator = newSurveyValidator(validatorTranslator)

// newValidatorTranslator returns a new validator translator.
func newValidatorTranslator() ut.Translator {
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	return trans
}

// newSurveyValidator returns a new validator with custom validators registered.
func newSurveyValidator(trans ut.Translator) *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())

	// Register custom validators
	for name, fn := range customValidators {
		_ = v.RegisterValidation(name, fn)
	}

	// Register custom translations
	for _, fn := range customTranslations {
		tag, register, translation := fn()
		_ = v.RegisterTranslation(tag, trans, register, translation)
	}

	return v
}

func TranslateValidationErrors(err error) []error {
	var validationErrorsTranslations []error
	for _, err := range err.(validator.ValidationErrors) {
		msg := errors.New(err.Translate(validatorTranslator))
		validationErrorsTranslations = append(validationErrorsTranslations, fmt.Errorf("%s", msg))
	}
	return validationErrorsTranslations
}

// ------------ Custom validators ------------ //

// customValidators is a map of custom validators.
var customValidators = map[string]func(fl validator.FieldLevel) bool{
	"questionType":          questionType,
	"validNameId":           validNameId,
	"validIfExternalSurvey": validIfExternalSurvey,
}

// questionType is a custom validator. It checks if the question is a valid question types.QuestionType
func questionType(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	if _, err := types.ParseToQuestionType(s); err != nil {
		return false
	}
	return true
}

// validNameId is a custom validator. It checks if the question nameId is valid.
func validNameId(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	return questionNameIdRegexValidator.MatchString(s)
}

// validIfExternalSurvey is a custom validator. It checks group.IsExternalSurvey and group.QuestionsIds.
// Rules:
// - If Group.IsExternalSurvey is false, ignore validation.
// - If Group.IsExternalSurvey is true, Group.QuestionsIds length must be 1.
func validIfExternalSurvey(fl validator.FieldLevel) bool {
	isExternalSurvey := fl.Field().Bool()
	if !isExternalSurvey {
		return true
	}

	questionsIds := fl.Parent().FieldByName("QuestionsIds").Interface().([]string)
	return len(questionsIds) == 1
}

// ------------ Custom translations ------------ //

// customTranslations is a map of custom translations.
var customTranslations = []func() (string, validator.RegisterTranslationsFunc, validator.TranslationFunc){
	minTranslation, validNameIdTranslation, validIfExternalTranslation, questionTypeTranslation,
}

func minTranslation() (string, validator.RegisterTranslationsFunc, validator.TranslationFunc) {
	tag := "min"
	register := func(ut ut.Translator) error {
		return ut.Add(tag, "Key: '{0}', Tag: '{1}', Error: '{2}' length must be >= {3}.", true)
	}

	translation := func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.StructNamespace(), tag, fe.Field(), fe.Param())
		return t
	}

	return tag, register, translation
}

func validNameIdTranslation() (string, validator.RegisterTranslationsFunc, validator.TranslationFunc) {
	tag := "validNameId"

	register := func(ut ut.Translator) error {
		return ut.Add(tag, "Key: '{0}', Tag: '{1}', Error: '{2}' must be a valid id. It should start with a letter, followed by 1-62 alphanumeric characters, (_) or (-), and end with an alphanumeric character.", true)
	}

	translation := func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.StructNamespace(), tag, fe.Field(), fe.Tag(), fe.Param())
		return t
	}

	return tag, register, translation
}

func validIfExternalTranslation() (string, validator.RegisterTranslationsFunc, validator.TranslationFunc) {
	tag := "validIfExternalSurvey"

	register := func(ut ut.Translator) error {
		return ut.Add(tag, "Key: '{0}', Tag: '{1}', Error: If 'IsExternalSurvey' is true, 'QuestionsIds' length must be 1.", true)
	}

	translation := func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.StructNamespace(), tag)
		return t
	}

	return tag, register, translation
}

func questionTypeTranslation() (string, validator.RegisterTranslationsFunc, validator.TranslationFunc) {
	tag := "questionType"

	register := func(ut ut.Translator) error {
		return ut.Add(tag, "Key: '{0}', Tag: '{1}', Error: '{2}' must be a valid question type.", true)
	}

	translation := func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.StructNamespace(), tag, fe.Field())
		return t
	}

	return tag, register, translation
}
