package surveygo

import (
	"errors"
	"fmt"
	"github.com/rendis/surveygo/v2/question"
	"github.com/rendis/surveygo/v2/question/types"
	"github.com/rendis/surveygo/v2/question/types/choice"
	"github.com/rendis/surveygo/v2/reviewer"
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
	GroupsResume map[string]*TotalsResume `json:"groupsResume,omitempty" bson:"groupsResume,omitempty"`

	//----- Errors -----//
	// InvalidAnswers list of invalid answers
	InvalidAnswers []InvalidAnswerError `json:"invalidAnswers,omitempty" bson:"invalidAnswers,omitempty"`
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

// ReviewAnswers verifies if the answers provided are valid for this survey.
// Args:
// * ans: the answers to check.
// Returns:
// * map[string]bool:
//   - key: the name id of the missing question
//   - value: if the question is required or not
//   - error: if an error occurred
func (s *Survey) ReviewAnswers(ans Answers) (*SurveyResume, error) {
	var invalidAnswers []InvalidAnswerError

	for nameId, values := range ans {
		q, ok := s.Questions[nameId]

		if !ok {
			continue
		}

		checker, err := reviewer.GetQuestionReviewer(q.QTyp)
		if err != nil {
			return nil, err
		}

		if err = checker(q.Value, values, q.QTyp); err != nil {
			invalidAnswers = append(invalidAnswers, InvalidAnswerError{
				QuestionNameId: nameId,
				Answer:         values,
				Error:          err.Error(),
			})
		}
	}

	resume := s.getSurveyResume(ans)
	resume.InvalidAnswers = invalidAnswers

	return resume, nil
}

// getSurveyResume returns the resume of the survey based on the answers provided.
func (s *Survey) getSurveyResume(ans Answers) *SurveyResume {
	var resume = &SurveyResume{
		TotalsResume: TotalsResume{
			UnansweredQuestions: make(map[string]bool),
		},
		GroupsResume:      make(map[string]*TotalsResume),
		ExternalSurveyIds: make(map[string]string),
	}

	questionWithGroup := s.getVisibleQuestionFromVisibleGroups()

	for qId, gId := range questionWithGroup {
		q := s.Questions[qId]

		// init group resume if not exists
		if _, ok := resume.GroupsResume[gId]; !ok {
			resume.GroupsResume[gId] = &TotalsResume{
				UnansweredQuestions: map[string]bool{},
			}
		}

		// update totals
		resume.TotalQuestions++
		resume.GroupsResume[gId].TotalQuestions++
		if q.Required {
			resume.TotalRequiredQuestions++
			resume.GroupsResume[gId].TotalRequiredQuestions++
		}

		// update answers
		if _, ok := ans[q.NameId]; ok {
			resume.TotalQuestionsAnswered++
			resume.GroupsResume[gId].TotalQuestionsAnswered++
			if q.Required {
				resume.TotalRequiredQuestionsAnswered++
				resume.GroupsResume[gId].TotalRequiredQuestionsAnswered++
			}
			continue
		}
		resume.UnansweredQuestions[q.NameId] = q.Required
		resume.GroupsResume[gId].UnansweredQuestions[q.NameId] = q.Required
	}

	// update external survey ids
	for _, g := range s.Groups {
		if g.IsExternalSurvey {
			resume.ExternalSurveyIds[g.NameId] = g.QuestionsIds[0]
		}
	}

	return resume
}

// getVisibleQuestionFromVisibleGroups returns a maps with the visible questions within its visible groups nameId.
func (s *Survey) getVisibleQuestionFromVisibleGroups() map[string]string {
	var questionWithGroup = map[string]string{}
	for _, g := range s.Groups {
		// skip invisible group
		if !g.Visible {
			continue
		}
		for _, q := range g.QuestionsIds {
			// get only visible question
			if s.Questions[q].Visible {
				questionWithGroup[q] = g.NameId
			}
		}
	}
	return questionWithGroup
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
