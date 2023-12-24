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

// TranslateAnswers translates the nameIDs of the answers to the values provided in each question (if any, otherwise the nameID is used).
// Translations:
// * text type: the value is the same passed in the answer
// * simple choice type: the value is the value, if any, of the choice with the same nameID as the answer
func (s *Survey) TranslateAnswers(ans Answers) (Answers, error) {
	var res = make(Answers, len(ans))

	for nameId, answers := range ans {
		q, ok := s.Questions[nameId]
		if !ok {
			return nil, fmt.Errorf("question '%s' not found", nameId)
		}

		// if text type, the value is the same passed in the answer
		if types.IsTextType(q.QTyp) {
			res[nameId] = answers
			continue
		}

		// if simple choice type, the value is the value, if any, of the choice with the same nameID as the answer
		if types.IsSimpleChoiceType(q.QTyp) {
			translatedValues := make([]any, len(answers))
			c, err := choice.CastToChoice(q)
			if err != nil {
				return nil, err
			}

			var optionsMap = make(map[string]*choice.Option)
			for _, option := range c.Options {
				optionsMap[option.NameId] = &option
			}

			for _, answer := range answers {
				answeredNameId, ok := answer.(string)
				if !ok {
					return nil, fmt.Errorf("invalid type, expected string, got '%T'", answer)
				}

				option, ok := optionsMap[answeredNameId]
				if !ok {
					return nil, fmt.Errorf("option '%s' not found", answeredNameId)
				}

				// if the option has a value, use it, otherwise use the answered name id
				if option.Value == nil {
					res[nameId] = append(res[nameId], option.Value)
				}

				res[nameId] = append(res[nameId], answeredNameId)
			}

			res[nameId] = translatedValues
			continue
		}

	}

	return res, nil
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
			if q, ok := s.Questions[questionNameId]; ok && q.Visible {
				questionWithGroup[questionNameId] = group.NameId
			}
		}
	}
	return questionWithGroup
}
