package surveygo

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rendis/surveygo/v2/question"
	"github.com/rendis/surveygo/v2/question/types"
	"github.com/rendis/surveygo/v2/question/types/choice"
)

// AddQuestion adds a question to the survey.
// It also validates the question and checks if the question is consistent with the survey.
func (s *Survey) AddQuestion(q *question.Question) error {
	if err := s.addQuestion(q); err != nil {
		return err
	}

	s.positionUpdater()
	return nil
}

// AddQuestionMap adds a question to the survey given its representation as a map[string]any
// It also validates the question and checks if the question is consistent with the survey.
func (s *Survey) AddQuestionMap(qm map[string]any) error {
	b, _ := json.Marshal(qm)
	return s.AddQuestionBytes(b)
}

// AddQuestionJson adds a question to the survey given its representation as a JSON string
// It also validates the question and checks if the question is consistent with the survey.
func (s *Survey) AddQuestionJson(qs string) error {
	return s.AddQuestionBytes([]byte(qs))
}

// AddQuestionBytes adds a question to the survey given its representation as a byte array
// It also validates the question and checks if the question is consistent with the survey.
func (s *Survey) AddQuestionBytes(qb []byte) error {
	// unmarshal question
	var q = &question.Question{}
	err := json.Unmarshal(qb, q)
	if err != nil {
		return err
	}

	// add question to survey
	return s.AddQuestion(q)
}

// AddOrUpdateQuestion adds a question to the survey if it does not exist, or updates it if it does.
// It also validates the question and checks if the question is consistent with the survey.
func (s *Survey) AddOrUpdateQuestion(q *question.Question) error {
	if q == nil {
		return errors.New("question is nil")
	}

	// check if question already exists
	if _, ok := s.Questions[q.NameId]; ok {
		return s.UpdateQuestion(q)
	}

	return s.AddQuestion(q)
}

// AddOrUpdateQuestionMap adds a question to the survey if it does not exist, or updates it if it does.
// It also validates the question and checks if the question is consistent with the survey.
func (s *Survey) AddOrUpdateQuestionMap(qm map[string]any) error {
	b, _ := json.Marshal(qm)
	return s.AddOrUpdateQuestionBytes(b)
}

// AddOrUpdateQuestionJson adds a question to the survey if it does not exist, or updates it if it does.
// It also validates the question and checks if the question is consistent with the survey.
func (s *Survey) AddOrUpdateQuestionJson(qs string) error {
	return s.AddOrUpdateQuestionBytes([]byte(qs))
}

// AddOrUpdateQuestionBytes adds a question to the survey if it does not exist, or updates it if it does.
// It also validates the question and checks if the question is consistent with the survey.
func (s *Survey) AddOrUpdateQuestionBytes(qb []byte) error {
	// unmarshal question
	var q = &question.Question{}
	err := json.Unmarshal(qb, q)
	if err != nil {
		return err
	}

	// add question to survey
	return s.AddOrUpdateQuestion(q)
}

// UpdateQuestion updates an existing question in the survey.
// It also validates the question and checks if the question is consistent with the survey.
func (s *Survey) UpdateQuestion(uq *question.Question) error {
	if uq == nil {
		return errors.New("question is nil")
	}

	q, ok := s.Questions[uq.NameId]
	if !ok {
		return fmt.Errorf("question '%s' not found", uq.NameId)
	}

	// update question
	q.Visible = uq.Visible
	q.QTyp = uq.QTyp
	q.Label = uq.Label
	q.Required = uq.Required
	q.Value = uq.Value
	q.Metadata = uq.Metadata
	q.Disabled = uq.Disabled
	q.DependsOn = uq.DependsOn

	// check consistency
	if err := s.checkConsistency(); err != nil {
		return err
	}

	// run position updater
	s.positionUpdater()
	return nil
}

// UpdateQuestionMap updates an existing question in the survey with the data provided as a map.
// It also validates the question and checks if the question is consistent with the survey.
func (s *Survey) UpdateQuestionMap(uq map[string]any) error {
	b, _ := json.Marshal(uq)
	return s.UpdateQuestionBytes(b)
}

// UpdateQuestionJson updates an existing question in the survey with the data provided as a JSON string.
// It also validates the question and checks if the question is consistent with the survey.
func (s *Survey) UpdateQuestionJson(uq string) error {
	return s.UpdateQuestionBytes([]byte(uq))
}

// UpdateQuestionBytes updates a question in the survey given its representation as a byte array
// It also validates the question and checks if the question is consistent with the survey.
func (s *Survey) UpdateQuestionBytes(uq []byte) error {
	// unmarshal question
	var q *question.Question
	err := json.Unmarshal(uq, q)
	if err != nil {
		return err
	}

	return s.UpdateQuestion(q)
}

// RemoveQuestion removes a question from the survey given its nameId.
// It also validates the group and checks if the group is consistent with the survey.
func (s *Survey) RemoveQuestion(questionNameId string) error {
	_, ok := s.Questions[questionNameId]
	if !ok {
		return fmt.Errorf("question '%s' not found", questionNameId)
	}

	// remove question from group
	for _, g := range s.Groups {
		if removed := g.RemoveQuestionId(questionNameId); removed {
			break
		}
	}

	// remove DependsOn references to this question from other questions
	for _, q := range s.Questions {
		q.DependsOn = removeDependsOnByQuestion(q.DependsOn, questionNameId)
	}

	// remove DependsOn references to this question from groups
	for _, g := range s.Groups {
		g.DependsOn = removeDependsOnByQuestion(g.DependsOn, questionNameId)
	}

	// remove question from survey
	delete(s.Questions, questionNameId)

	s.positionUpdater()

	return nil
}

// GetQuestionsAssignments returns a map with the question nameId as key and the group nameId as value.
// If a question is not assigned to any group, value will be empty.
func (s *Survey) GetQuestionsAssignments() map[string]string {
	var questionsAssignation = make(map[string]string)

	// get questions from groups
	for _, g := range s.Groups {
		for _, q := range g.QuestionsIds {
			questionsAssignation[q] = g.NameId
		}
	}

	// get questions not in groups
	for q := range s.Questions {
		if _, ok := questionsAssignation[q]; !ok {
			questionsAssignation[q] = ""
		}
	}

	return questionsAssignation
}

// GetAssetQuestions returns all the asset questions in the survey.
// Asset questions are questions that have a type of image, video, audio or document.
func (s *Survey) GetAssetQuestions() []*question.Question {
	var assetQuestions []*question.Question
	for _, q := range s.Questions {
		if types.IsAssetType(q.QTyp) {
			assetQuestions = append(assetQuestions, q)
		}
	}

	return assetQuestions
}

// addQuestion adds a question to the survey.
// It also validates the question and checks if the question is consistent with the survey.
func (s *Survey) addQuestion(q *question.Question) error {
	if q == nil {
		return errors.New("question is nil")
	}

	// validate question
	if err := SurveyValidator.Struct(q); err != nil {
		errs := TranslateValidationErrors(err)
		errs = append([]error{fmt.Errorf("error validating question")}, errs...)
		return errors.Join(errs...)
	}

	// check if question already exists
	if _, ok := s.Questions[q.NameId]; ok {
		return fmt.Errorf("question nameId '%s' already exists", q.NameId)
	}

	// if question is choice type, check if options groups exist
	if types.IsChoiceType(q.QTyp) {
		c, _ := choice.CastToChoice(q.Value)
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
	s.Questions[q.NameId] = q

	// check consistency
	return s.checkConsistency()
}

// removeDependsOnByQuestion removes all DependsOn conditions that reference the given questionNameId.
// Only the specific conditions referencing the question are removed; other conditions in the AND group are kept.
// If an AND group becomes empty (all conditions removed), the entire AND group is removed.
// Returns nil if all AND groups are removed.
//
// Example 1 - Removing specific conditions, keeping others:
//
//	dependsOn = [
//	  [ {questionNameId: "q1", optionNameId: "o1"}, {questionNameId: "q2", optionNameId: "o2"} ], // dependsOn: N0 (q1.o1 AND q2.o2)
//	  [ {questionNameId: "q3", optionNameId: "o3"} ], // dependsOn: N1 (single condition, no q1)
//	]
//
//	removeDependsOnByQuestion(dependsOn, "q1") returns:
//	[
//	  [ {questionNameId: "q2", optionNameId: "o2"} ], // dependsOn: N0 kept with only q2.o2 (q1.o1 removed)
//	  [ {questionNameId: "q3", optionNameId: "o3"} ], // dependsOn: N1 kept unchanged
//	]
//
// Example 2 - AND group removed when all conditions reference deleted question:
//
//	dependsOn = [
//	  [ {questionNameId: "q1", optionNameId: "o1"} ], // dependsOn: N0 (single condition)
//	  [ {questionNameId: "q2", optionNameId: "o2"} ], // dependsOn: N1 (no reference to q1)
//	  [ {questionNameId: "q1", optionNameId: "o4"} ], // dependsOn: N2 (references q1)
//	]
//
//	removeDependsOnByQuestion(dependsOn, "q1") returns:
//	[
//	  [ {questionNameId: "q2", optionNameId: "o2"} ], // dependsOn: N1 kept
//	]
//	// dependsOn: N0 and N2 removed because they only had conditions referencing q1
func removeDependsOnByQuestion(dependsOn [][]question.DependsOn, questionNameId string) [][]question.DependsOn {
	if len(dependsOn) == 0 {
		return dependsOn
	}

	var result [][]question.DependsOn
	for _, andGroup := range dependsOn {
		// filter out conditions that reference the deleted question
		var filteredAndGroup []question.DependsOn
		for _, dep := range andGroup {
			if dep.QuestionNameId != questionNameId {
				filteredAndGroup = append(filteredAndGroup, dep)
			}
		}
		// only keep AND groups that still have conditions
		if len(filteredAndGroup) > 0 {
			result = append(result, filteredAndGroup)
		}
	}
	return result
}
