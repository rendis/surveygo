package surveygo

import (
	"errors"
	"fmt"
	"github.com/rendis/devtoolkit"
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

// ReviewAnswers verifies if the answers provided are valid for this survey.
// Args:
// * ans: the answers to check.
// Returns:
// * map[string]bool:
//   - key: the name id of the missing question
//   - value: if the question is required or not
//   - error: if an error occurred
func (s *Survey) ReviewAnswers(ans Answers) (*SurveyResume, error) {
	var invalidAnswers []*InvalidAnswerError

	var correctAnswersCount = make(map[string]int)
	var groupsCount = make(map[string]int)

	for nameId, values := range ans {
		// if nameId is a question
		if s.Questions[nameId] != nil {
			if invalid := s.reviewQuestion(nameId, values, correctAnswersCount); invalid != nil {
				invalidAnswers = append(invalidAnswers, invalid)
			}
			continue
		}

		// if nameId is a group
		if s.Groups[nameId] != nil {
			if invalids := s.reviewGroup(nameId, values, correctAnswersCount, groupsCount); len(invalids) > 0 {
				invalidAnswers = append(invalidAnswers, invalids...)
			}
			continue
		}

		// if nameId is not a question or a group
		invalidAnswers = append(invalidAnswers, &InvalidAnswerError{
			QuestionNameId: nameId,
			Answer:         values,
			Error:          fmt.Sprintf("question '%s' not found", nameId),
		})
	}

	if len(invalidAnswers) > 0 {
		return &SurveyResume{InvalidAnswers: invalidAnswers}, nil
	}

	return s.getSurveyResume(correctAnswersCount, groupsCount), nil
}

// reviewQuestion verifies if the answer provided is valid for the given question.
func (s *Survey) reviewQuestion(questionNameID string, answers []any, correctAnswersCount map[string]int) *InvalidAnswerError {
	q := s.Questions[questionNameID]
	reviewerFn, err := reviewer.GetQuestionReviewer(q.QTyp)
	if err != nil {
		return &InvalidAnswerError{
			QuestionNameId: questionNameID,
			Answer:         answers,
			Error:          err.Error(),
		}
	}

	if err = reviewerFn(q.Value, answers, q.QTyp); err != nil {
		return &InvalidAnswerError{
			QuestionNameId: questionNameID,
			Answer:         answers,
			Error:          err.Error(),
		}
	}

	// update correct answers count
	correctAnswersCount[questionNameID]++
	return nil
}

// reviewGroup verifies if the answers provided are valid for the given group.
func (s *Survey) reviewGroup(groupNameID string, nestedAnswers []any, correctAnswersCount map[string]int, groupsCount map[string]int) []*InvalidAnswerError {
	groupAnswers, err := reviewer.ExtractGroupNestedAnswers(nestedAnswers)
	if err != nil {
		return []*InvalidAnswerError{{
			QuestionNameId: groupNameID,
			Answer:         nestedAnswers,
			Error:          err.Error(),
		}}
	}
	var invalidAnswers []*InvalidAnswerError

	if err != nil {
		invalidAnswers = append(invalidAnswers, &InvalidAnswerError{
			QuestionNameId: groupNameID,
			Answer:         nestedAnswers,
			Error:          err.Error(),
		})
		return invalidAnswers
	}

	for _, groupedAnswers := range groupAnswers {
		for questionNameID, answers := range groupedAnswers {
			if invalid := s.reviewQuestion(questionNameID, answers, correctAnswersCount); invalid != nil {
				invalidAnswers = append(invalidAnswers, invalid)
			}
		}
	}

	// update correct answers count
	groupsCount[groupNameID] += len(groupAnswers)
	return invalidAnswers
}

// getSurveyResume returns the resume of the survey based on the answers provided.
func (s *Survey) getSurveyResume(correctAnswersCount map[string]int, groupsCount map[string]int) *SurveyResume {
	var resume = &SurveyResume{
		TotalsResume: TotalsResume{
			UnansweredQuestions: make(map[string]bool),
		},
		GroupsResume:      make(map[string]*GroupTotalsResume),
		ExternalSurveyIds: make(map[string]string),
	}

	visibleQuestionWithGroup := s.getVisibleQuestionFromActiveGroups()

	for questionId, groupId := range visibleQuestionWithGroup {
		q := s.Questions[questionId]

		// init group resume if not exists
		if _, ok := resume.GroupsResume[groupId]; !ok {
			resume.GroupsResume[groupId] = &GroupTotalsResume{
				TotalsResume: TotalsResume{
					UnansweredQuestions: map[string]bool{},
				},
				AnswerGroups: groupsCount[groupId],
			}
		}

		questionResponsesCount := correctAnswersCount[questionId]
		questionOccurrenceCount := devtoolkit.IfThenElse(groupsCount[groupId] == 0, 1, groupsCount[groupId])

		// update totals
		resume.TotalQuestions += questionOccurrenceCount
		resume.GroupsResume[groupId].TotalQuestions += questionOccurrenceCount
		if q.Required {
			resume.TotalRequiredQuestions += questionOccurrenceCount
			resume.GroupsResume[groupId].TotalRequiredQuestions += questionOccurrenceCount
		}

		// update answers
		resume.TotalQuestionsAnswered += questionResponsesCount
		resume.GroupsResume[groupId].TotalQuestionsAnswered += questionResponsesCount
		if q.Required {
			resume.TotalRequiredQuestionsAnswered += questionResponsesCount
			resume.GroupsResume[groupId].TotalRequiredQuestionsAnswered += questionResponsesCount
		}

		// update unanswered questions
		if questionResponsesCount == 0 {
			resume.UnansweredQuestions[q.NameId] = q.Required
			resume.GroupsResume[groupId].UnansweredQuestions[q.NameId] = q.Required
		}
	}

	// update external survey ids
	for _, g := range s.Groups {
		if g.IsExternalSurvey {
			resume.ExternalSurveyIds[g.NameId] = g.QuestionsIds[0]
		}
	}

	return resume
}

// getVisibleQuestionFromActiveGroups returns a maps with the visible questions within its active groups nameId.
// and active groups are the groups that are visible and enabled.
func (s *Survey) getVisibleQuestionFromActiveGroups() map[string]string {
	var questionWithGroup = map[string]string{}
	for _, group := range s.Groups {
		// skip hidden && disabled groups
		if group.Hidden || group.Disabled {
			continue
		}
		for _, questionNameId := range group.QuestionsIds {
			// get only visible question
			if s.Questions[questionNameId].Visible {
				questionWithGroup[questionNameId] = group.NameId
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
