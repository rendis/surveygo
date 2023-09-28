package surveygo

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rendis/surveygo/v2/question"
	"github.com/rendis/surveygo/v2/question/types"
	"github.com/rendis/surveygo/v2/question/types/choice"
)

// NewSurvey creates a new Survey instance with the given title, version, and description.
// Args:
//   - nameId: the name id of the survey (required)
//   - title: the title of the survey (required)
//   - version: the version of the survey (required)
//   - description: the description of the survey (optional)
//
// Returns:
//   - *Survey: the new survey instance
//   - error: if an error occurred
func NewSurvey(nameId, title, version string, description *string) (*Survey, error) {
	return &Survey{
		NameId:      nameId,
		Title:       title,
		Version:     version,
		Description: description,
		Questions:   map[string]question.Question{},
		Groups:      map[string]question.Group{},
		GroupsOrder: []string{},
	}, nil
}

// Parse converts the given json string into a Survey instance.
func Parse(jsonSurvey string) (*Survey, error) {
	return ParseBytes([]byte(jsonSurvey))
}

// ParseBytes converts the given json byte slice into a Survey instance.
func ParseBytes(b []byte) (*Survey, error) {
	var survey = &Survey{}

	// unmarshal the json survey into the survey struct
	if err := json.Unmarshal(b, survey); err != nil {
		return nil, errors.Join(fmt.Errorf("error unmarshalling json survey"), err)
	}

	// validate the survey struct
	if err := SurveyValidator.Struct(survey); err != nil {
		errs := TranslateValidationErrors(err)
		errs = append([]error{fmt.Errorf("error validating survey")}, errs...)
		return nil, errors.Join(errs...)
	}

	// check survey consistency
	if err := survey.checkConsistency(); err != nil {
		return nil, err
	}

	return survey, nil
}

// Answers is a map with the answers provided by the user.
// The key is the question NameId (Question.NameId).
type Answers map[string][]any

// Survey is a struct representation of a survey.
type Survey struct {
	// NameId is the name id of the survey.
	// Validations:
	//	- required
	//	- valid name id
	NameId string `json:"nameId" bson:"nameId" validate:"required,validNameId"`

	// Title is the title of the survey.
	// Validations:
	//	- required
	//	- min length: 1
	Title string `json:"title" bson:"title" validate:"required,min=1"`

	// Version is the version of the survey.
	// Validations:
	//	- required
	//	- min length: 1
	Version string `json:"version" bson:"version" validate:"required,min=1"`

	// Description is the description of the survey.
	// Validations:
	//	- optional
	//	- min length: 1
	Description *string `json:"description" bson:"description" validate:"omitempty"`

	// Questions is a map with all the questions in the survey.
	// The key is the question NameId (Question.NameId).
	// Validations:
	//	- required
	//	- min length: 1
	//	- each question must be valid
	Questions map[string]question.Question `json:"questions" bson:"questions" validate:"required,dive"`

	// Groups is a map with all the groups in the survey.
	// The key is the group NameId (Group.NameId).
	// Validations:
	//	- required
	//	- min length: 1
	//	- each group must be valid
	Groups map[string]question.Group `json:"groups" bson:"groups" validate:"required,dive"`

	// GroupsOrder is a list of group name ids that defines the order of the groups in the survey.
	// Validations:
	//	- required
	//	- min length: 1
	GroupsOrder []string `json:"groupsOrder" bson:"groupsOrder" validate:"required"`
}

// checkConsistency checks the consistency of the survey.
func (s *Survey) checkConsistency() error {
	var errs []error

	// check questions
	optionsProcessed := map[string]bool{} // key: option name id, value: true if the option was processed
	groupsProcessed := map[string]bool{}  // key: group name id, value: true if the group was processed
	for k, q := range s.Questions {
		if k != q.NameId {
			errs = append(errs, fmt.Errorf("question key '%s' does not match question name id '%s'", k, q.NameId))
			continue
		}

		if !types.IsChoiceType(q.QTyp) {
			continue
		}

		c, _ := choice.CastToChoice(q.Value)
		var options = c.GetOptionsGroups() // key: option name id, value: list of group name ids
		for optionNameId, groupsIds := range options {

			// check if the option name id was processed
			if optionsProcessed[optionNameId] {
				errs = append(errs, fmt.Errorf("duplicate option id '%s'", optionNameId))
				continue
			}
			optionsProcessed[optionNameId] = true

			// check the group name ids
			for _, groupNameId := range groupsIds {
				// check if the group name id exists
				if _, ok := s.Groups[groupNameId]; !ok {
					errs = append(errs, fmt.Errorf("group id '%s' not found for option id '%s'", groupNameId, optionNameId))
				}

				// check if the group id was processed
				if groupsProcessed[groupNameId] {
					errs = append(errs, fmt.Errorf("group id '%s' is duplicated for options", groupNameId))
				}
				groupsProcessed[groupNameId] = true
			}
		}
	}

	// check groups
	questionsProcessed := map[string]bool{} // key: question name id, value: true if the question was processed
	for k, g := range s.Groups {
		if k != g.NameId {
			errs = append(errs, fmt.Errorf("group key '%s' does not match group name id '%s'", k, g.NameId))
			continue
		}

		// skip external groups
		if g.IsExternalSurvey {
			continue
		}

		for _, questionNameId := range g.QuestionsIds {
			// check if the question name id exists
			if _, ok := s.Questions[questionNameId]; !ok {
				errs = append(errs, fmt.Errorf("question id '%s' not found for group id '%s'", questionNameId, g.NameId))
			}

			// check if the question name id was processed
			if questionsProcessed[questionNameId] {
				errs = append(errs, fmt.Errorf("question id '%s' in multiple groups", questionNameId))
				continue
			}
			questionsProcessed[questionNameId] = true
		}
	}

	// check groups order
	for _, groupNameId := range s.GroupsOrder {
		// check if the group name id exists
		if _, ok := s.Groups[groupNameId]; !ok {
			errs = append(errs, fmt.Errorf("group id '%s' in groups order not found", groupNameId))
		}

		// check if the group name id was processed
		if groupsProcessed[groupNameId] {
			errs = append(errs, fmt.Errorf("group id '%s' found multiple times", groupNameId))
		}
	}

	// build the error
	if len(errs) > 0 {
		var consErr = fmt.Errorf("error checking survey consistency")
		errs = append([]error{consErr}, errs...)
		return errors.Join(errs...)
	}

	return nil
}
