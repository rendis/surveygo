package surveygo

import (
	"errors"
	"fmt"
	"github.com/rendis/surveygo/v2/question"
	"github.com/rendis/surveygo/v2/question/types"
	"github.com/rendis/surveygo/v2/question/types/choice"
)

type InvalidAnswerError struct {
	QuestionNameId string `json:"questionNameId,omitempty" bson:"questionNameId,omitempty"`
	Answer         any    `json:"answer,omitempty" bson:"answer,omitempty"`
	Error          string `json:"error,omitempty" bson:"error,omitempty"`
}

type TotalsResume struct {
	//----- Questions Totals -----//
	// TotalQuestions number of questions in the group
	TotalQuestions int `json:"totalQuestions,omitempty" bson:"totalQuestions,omitempty"`
	// TotalRequiredQuestions number of required questions in the group
	TotalRequiredQuestions int `json:"totalRequiredQuestions,omitempty" bson:"totalRequiredQuestions,omitempty"`

	//----- Answers Totals  -----//
	// TotalQuestionsAnswered number of answered questions in the group
	TotalQuestionsAnswered int `json:"totalQuestionsAnswered,omitempty" bson:"totalQuestionsAnswered,omitempty"`
	// TotalRequiredQuestionsAnswered number of required questions answered in the group
	TotalRequiredQuestionsAnswered int `json:"totalRequiredQuestionsAnswered,omitempty" bson:"totalRequiredQuestionsAnswered,omitempty"`
	// UnansweredQuestions map of unanswered questions, key is the nameId of the question, value is true if the question is required
	UnansweredQuestions map[string]bool `json:"unansweredQuestions,omitempty" bson:"unansweredQuestions,omitempty"`
}

// SurveyResume contains the resume of a survey based on the answers provided.
// All values are calculated based on the answers provided and over the visible components of the survey.
type SurveyResume struct {
	TotalsResume `json:",inline" bson:",inline"`

	//----- Others Totals -----//
	// ExternalSurveyIds map of external survey ids. Key: GroupNameId, Value: ExternalSurveyId
	ExternalSurveyIds map[string]string `json:"externalSurveyIds,omitempty" bson:"externalSurveyIds,omitempty"`

	//----- Groups -----//
	// GroupsResume map of groups resume. Key: GroupNameId, Value: GroupResume
	GroupsResume map[string]*GroupTotalsResume `json:"groupsResume,omitempty" bson:"groupsResume,omitempty"`

	//----- Errors -----//
	// InvalidAnswers list of invalid answers
	InvalidAnswers []*InvalidAnswerError `json:"invalidAnswers,omitempty" bson:"invalidAnswers,omitempty"`
}

// GroupTotalsResume contains the resume of a group based on the answers provided.
type GroupTotalsResume struct {
	TotalsResume `json:",inline" bson:",inline"`
	AnswerGroups int `json:"answersGroups,omitempty" bson:"answersGroups,omitempty"`
}

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
func NewSurvey(title, version string, description *string) (*Survey, error) {
	return &Survey{
		Title:       title,
		Version:     version,
		Description: description,
		Questions:   map[string]*question.Question{},
		Groups:      map[string]*question.Group{},
		GroupsOrder: []string{},
	}, nil
}

// ValidateSurvey validates the survey.
func (s *Survey) ValidateSurvey() error {
	return s.checkConsistency()
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

		// check questions
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

		// check groups in GroupsOrder
		for _, groupNameId := range g.GroupsOrder {
			// check if the group name id exists
			if _, ok := s.Groups[groupNameId]; !ok {
				errs = append(errs, fmt.Errorf("group id '%s' in groups order of group id '%s' not found", groupNameId, g.NameId))
			}
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

// positionUpdater runs the position assignation for the survey.
// It assigns a position to each question and group.
func (s *Survey) positionUpdater() {
	qPos := 1
	gPos := 1
	for _, g := range s.GroupsOrder {
		s.Groups[g].Position = gPos
		gPos++
		for _, q := range s.Groups[g].QuestionsIds {
			s.Questions[q].Position = qPos
			qPos++
		}
	}
}
