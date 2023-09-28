package surveygo

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rendis/surveygo/v2/question"
	"github.com/rendis/surveygo/v2/question/types"
	"github.com/rendis/surveygo/v2/question/types/choice"
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
	// TotalAnsweredQuestions number of answered questions in the survey
	TotalAnsweredQuestions int `json:"totalAnsweredQuestions" bson:"totalAnsweredQuestions"`
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

// ValidateSurvey validates the survey.
func (s *Survey) ValidateSurvey() error {
	return s.checkConsistency()
}

// ToMap returns a map representation of the survey.
func (s *Survey) ToMap() (map[string]any, error) {
	r := make(map[string]any)
	// s to map
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	// unmarshal to map[string]any
	if err = json.Unmarshal(b, &r); err != nil {
		return nil, err
	}

	return r, nil
}

// ToJson returns a JSON string representation of the survey.
func (s *Survey) ToJson() (string, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// RemoveQuestion removes a question from the survey given its nameId.
func (s *Survey) RemoveQuestion(nameId string) error {
	_, ok := s.Questions[nameId]
	if !ok {
		return fmt.Errorf("question '%s' not found", nameId)
	}

	// remove question from group
	questionWithGroup := s.getQuestionFromGroups()
	if groupNameId, ok := questionWithGroup[nameId]; ok {
		g := s.Groups[groupNameId]
		g.RemoveQuestionId(nameId)
	}

	// remove question from survey
	delete(s.Questions, nameId)

	return nil
}

// AddQuestionMap adds a question to the survey given its representation as a map[string]any
func (s *Survey) AddQuestionMap(question map[string]any) error {
	b, _ := json.Marshal(question)
	return s.AddQuestionBytes(b)
}

// AddQuestionJson adds a question to the survey given its representation as a JSON string
func (s *Survey) AddQuestionJson(question string) error {
	return s.AddQuestionBytes([]byte(question))
}

// AddQuestionBytes adds a question to the survey given its representation as a byte array
func (s *Survey) AddQuestionBytes(q []byte) error {
	// unmarshal question
	var pq *question.Question = &question.Question{}
	err := json.Unmarshal(q, pq)
	if err != nil {
		return err
	}

	// add question to survey
	return s.addQuestion(pq)
}

// UpdateQuestionMap updates an existing question in the survey with the data provided as a map.
func (s *Survey) UpdateQuestionMap(question map[string]any) error {
	b, _ := json.Marshal(question)
	return s.UpdateQuestionBytes(b)
}

// UpdateQuestionJson updates an existing question in the survey with the data provided as a JSON string.
func (s *Survey) UpdateQuestionJson(question string) error {
	return s.UpdateQuestionBytes([]byte(question))
}

// UpdateQuestionBytes updates a question in the survey given its representation as a byte array
func (s *Survey) UpdateQuestionBytes(q []byte) error {
	// unmarshal question
	var pq *question.Question
	err := json.Unmarshal(q, pq)
	if err != nil {
		return err
	}

	// remove question from survey
	if err = s.RemoveQuestion(pq.NameId); err != nil {
		return err
	}

	// add question to survey
	return s.addQuestion(pq)
}

// addQuestion adds a question to the survey.
func (s *Survey) addQuestion(pq *question.Question) error {
	// validate question
	if err := SurveyValidator.Struct(pq); err != nil {
		errs := TranslateValidationErrors(err)
		errs = append([]error{fmt.Errorf("error validating question")}, errs...)
		return errors.Join(errs...)
	}

	// check if question already exists
	if _, ok := s.Questions[pq.NameId]; ok {
		return fmt.Errorf("question nameId '%s' already exists", pq.NameId)
	}

	// if question is choice type, check if options groups
	if types.IsChoiceType(pq.QTyp) {
		c, _ := choice.CastToChoice(pq.Value)
		optionsGroups := c.GetOptionsGroups()
		for _, ogs := range optionsGroups {
			for _, og := range ogs {
				// check if options group exists
				if _, ok := s.Groups[og]; !ok {
					return fmt.Errorf("options group nameId '%s' not found", og)
				}
			}
		}
	}

	// add question to survey
	s.Questions[pq.NameId] = *pq

	return nil
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
			resume.TotalAnsweredQuestions++
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
