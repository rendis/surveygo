package surveygo

import (
	"github.com/rendis/surveygo/v2/reviewer"
)

type InvalidAnswerError struct {
	QuestionNameId string `json:"questionNameId" bson:"questionNameId"`
	Answer         any    `json:"answer" bson:"answer"`
	Error          string `json:"error" bson:"error"`
}

// SurveyResume contains the resume of a survey based on the answers provided.
// All values are calculated based on the answers provided and over the visible components of the survey.
type SurveyResume struct {
	//----- Questions Totals -----//
	// TotalQuestions number of questions in the survey
	TotalQuestions int `json:"totalQuestions" bson:"totalQuestions"`
	// TotalRequiredQuestions number of required questions in the survey
	TotalRequiredQuestions int `json:"totalRequiredQuestions" bson:"totalRequiredQuestions"`

	//----- Answers Totals  -----//
	// TotalQuestionsAnswered number of answered questions in the survey
	TotalQuestionsAnswered int `json:"totalQuestionsAnswered" bson:"totalQuestionsAnswered"`
	// TotalRequiredQuestionsAnswered number of required questions answered in the survey
	TotalRequiredQuestionsAnswered int `json:"totalRequiredQuestionsAnswered" bson:"totalRequiredQuestionsAnswered"`
	// UnansweredQuestions map of unanswered questions, key is the nameId of the question, value is true if the question is required
	UnansweredQuestions map[string]bool `json:"unansweredQuestions" bson:"unansweredQuestions"`

	//----- Others Totals -----//
	// ExternalSurveyIds map of external survey ids. Key: GroupNameId, Value: ExternalSurveyId
	ExternalSurveyIds map[string]string `json:"externalSurveyIds" bson:"externalSurveyIds"`

	//----- Errors -----//
	// InvalidAnswers list of invalid answers
	InvalidAnswers []InvalidAnswerError `json:"invalidAnswers" bson:"invalidAnswers"`
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
		UnansweredQuestions: map[string]bool{},
		ExternalSurveyIds:   map[string]string{},
	}

	questionWithGroup := s.getQuestionFromGroups()

	for _, q := range s.Questions {
		// skip invisible questions
		if !q.Visible {
			continue
		}

		// skip questions in invisible groups
		if groupNameId, ok := questionWithGroup[q.NameId]; ok {
			g := s.Groups[groupNameId]
			if !g.Visible {
				continue
			}
		}

		// update totals
		resume.TotalQuestions++
		if q.Required {
			resume.TotalRequiredQuestions++
		}

		// update answers
		if _, ok := ans[q.NameId]; ok {
			resume.TotalQuestionsAnswered++
			if q.Required {
				resume.TotalRequiredQuestionsAnswered++
			}
			continue
		}
		resume.UnansweredQuestions[q.NameId] = q.Required
	}

	// update external survey ids
	for _, g := range s.Groups {
		if g.IsExternalSurvey {
			resume.ExternalSurveyIds[g.NameId] = g.QuestionsIds[0]
		}
	}

	return resume
}

// getQuestionFromGroups returns a maps with the questions and groups information.
func (s *Survey) getQuestionFromGroups() map[string]string {
	var questionWithGroup = map[string]string{}
	for _, g := range s.Groups {
		for _, q := range g.QuestionsIds {
			questionWithGroup[q] = g.NameId
		}
	}
	return questionWithGroup
}
