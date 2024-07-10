package surveygo

import (
	"fmt"
	"github.com/rendis/devtoolkit"
	"github.com/rendis/surveygo/v2/question/types"
	"github.com/rendis/surveygo/v2/question/types/choice"
	"github.com/rendis/surveygo/v2/reviewer"
)

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
		if s.isQuestion(nameId) {
			if invalid := s.reviewQuestion(nameId, values, correctAnswersCount); invalid != nil {
				invalidAnswers = append(invalidAnswers, invalid)
			}
			continue
		}

		// if nameId is a group
		if s.isGroup(nameId) {
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

// TranslateAnswers translates the nameIDs of the answers to the values provided in each question (if any, otherwise the nameID is used).
// Translations:
// * text type: the value is the same passed in the answer
// * simple choice type: the value is the value, if any, of the choice with the same nameID as the answer
func (s *Survey) TranslateAnswers(ans Answers, ignoreUnknownAnswers bool) (Answers, error) {
	var res = make(Answers, len(ans))

	for nameId, answers := range ans {
		// if nameId is a question
		if s.isQuestion(nameId) {
			translations, err := s.translateAnswers(nameId, answers, ignoreUnknownAnswers)
			if err != nil {
				return nil, err
			}
			res[nameId] = translations
		}

		// if nameId is a group
		if s.isGroup(nameId) {
			groupAnswers, err := reviewer.ExtractGroupNestedAnswers(answers)
			if err != nil {
				return nil, fmt.Errorf("invalid group answers for group '%s'. %s", nameId, err)
			}

			const groupQuestionTemplate = "group.%d.%s"
			for i, groupAnswersPack := range groupAnswers {
				for questionNameId, answersPack := range groupAnswersPack {
					translations, err := s.translateAnswers(questionNameId, answersPack, ignoreUnknownAnswers)
					if err != nil {
						return nil, err
					}
					key := fmt.Sprintf(groupQuestionTemplate, i, questionNameId)
					res[key] = translations
				}
			}
		}
	}

	return res, nil
}

// GetDisabledQuestions returns a map with the name id of the disabled questions.
// Questions are disabled:
// * if the question is disabled
// * if the group of the question is disabled
func (s *Survey) GetDisabledQuestions() map[string]bool {
	var disabledQuestions = make(map[string]bool)

	for _, group := range s.Groups {
		groupDisabled := group.Disabled
		for _, questionNameId := range group.QuestionsIds {
			if q, ok := s.Questions[questionNameId]; ok && (groupDisabled || q.Disabled) {
				disabledQuestions[questionNameId] = true
			}
		}
	}

	return disabledQuestions
}

// GetEnabledQuestions returns a map with the name id of the enabled questions.
// Questions are enabled:
// * if the question is enabled and the group of the question is enabled
func (s *Survey) GetEnabledQuestions() map[string]bool {
	var enabledQuestions = make(map[string]bool)

	for _, group := range s.Groups {
		if group.Disabled {
			continue
		}

		for _, questionNameId := range group.QuestionsIds {
			if q, ok := s.Questions[questionNameId]; ok && !q.Disabled {
				enabledQuestions[questionNameId] = true
			}
		}
	}

	return enabledQuestions

}

// translateAnswers translates the nameIDs of the answers to the values provided in each question (if any, otherwise the nameID is used).
func (s *Survey) translateAnswers(nameId string, answers []any, ignoreUnknownAnswers bool) ([]any, error) {
	if len(answers) == 0 {
		return answers, nil
	}

	q, ok := s.Questions[nameId]
	if !ok {
		if ignoreUnknownAnswers {
			return answers, nil
		}
		return nil, fmt.Errorf("question '%s' not found", nameId)
	}

	// if text type, the value is the same passed in the answer
	if types.IsTextType(q.QTyp) {
		return answers, nil
	}

	// if simple choice type the translation is the value
	if types.IsSimpleChoiceType(q.QTyp) {
		var res []any
		c, err := choice.CastToChoice(q.Value)
		if err != nil {
			return nil, err
		}

		var optionsMap = make(map[string]*choice.Option)
		var optionsValuesMap = make(map[any]*choice.Option)

		for _, option := range c.Options {
			optionsMap[option.NameId] = option
			optionsValuesMap[option.Value] = option
		}

		for _, answer := range answers {
			answeredNameId, ok := answer.(string)
			if !ok {
				return nil, fmt.Errorf("invalid type, expected string, got '%T'", answer)
			}

			var option *choice.Option
			if option, ok = optionsMap[answeredNameId]; !ok {
				// if the option is not found, try to find it in the options by value
				option, ok = optionsValuesMap[answeredNameId]
			}

			if !ok {
				return nil, fmt.Errorf("option not found for question '%s' (searched by nameId and value '%s')", nameId, answeredNameId)
			}

			// if the option has a value, use it, otherwise use the answered name id
			if option.Value != nil {
				res = append(res, option.Value)
				continue
			}

			res = append(res, answeredNameId)
		}

		return res, nil
	}

	return answers, nil
}

func (s *Survey) translateGroupAnswers(groupNameID string) {}

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
			if q, ok := s.Questions[questionNameId]; ok && q.Visible {
				questionWithGroup[questionNameId] = group.NameId
			}
		}
	}
	return questionWithGroup
}

func (s *Survey) isQuestion(nameId string) bool {
	_, ok := s.Questions[nameId]
	return ok
}

func (s *Survey) isGroup(nameId string) bool {
	_, ok := s.Groups[nameId]
	return ok
}
