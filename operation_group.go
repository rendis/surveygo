package surveygo

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rendis/surveygo/v2/question"
	"github.com/rendis/surveygo/v2/question/types"
	"github.com/rendis/surveygo/v2/question/types/choice"
)

// UpdateGroupsOrder updates the groups order in the survey.
// It also validates the group and checks if the group is consistent with the survey.
func (s *Survey) UpdateGroupsOrder(order []string) error {
	var errs []error

	// check if all groups exist and are unique
	groups := map[string]bool{}
	for _, id := range order {
		// check if group exists
		if _, ok := s.Groups[id]; !ok {
			errs = append(errs, fmt.Errorf("group '%s' not found", id))
			continue
		}

		// check if group is duplicated
		if _, ok := groups[id]; ok {
			errs = append(errs, fmt.Errorf("group '%s' is duplicated", id))
			continue
		}
		groups[id] = true
	}

	// update groups order
	s.GroupsOrder = order

	// check consistency
	return s.checkConsistency()
}

// AddGroup adds a group to the survey.
// It also validates the group and checks if the group is consistent with the survey.
func (s *Survey) AddGroup(g *question.Group) error {
	return s.addGroup(g)
}

// AddGroupMap adds a group to the survey given its representation as a map[string]any
// It also validates the group and checks if the group is consistent with the survey.
func (s *Survey) AddGroupMap(g map[string]any) error {
	b, _ := json.Marshal(g)
	return s.AddGroupBytes(b)
}

// AddGroupJson adds a group to the survey given its representation as a JSON string
// It also validates the group and checks if the group is consistent with the survey.
func (s *Survey) AddGroupJson(g string) error {
	return s.AddGroupBytes([]byte(g))
}

// AddGroupBytes adds a group to the survey given its representation as a byte array
// It also validates the group and checks if the group is consistent with the survey.
func (s *Survey) AddGroupBytes(g []byte) error {
	// unmarshal group
	var pg *question.Group
	err := json.Unmarshal(g, pg)
	if err != nil {
		return err
	}

	// add group to survey
	return s.addGroup(pg)
}

// UpdateGroupMap updates an existing group in the survey with the data provided as a map.
// It also validates the group and checks if the group is consistent with the survey.
func (s *Survey) UpdateGroupMap(ug map[string]any) error {
	b, _ := json.Marshal(ug)
	return s.UpdateGroupBytes(b)
}

// UpdateGroupJson updates an existing group in the survey with the data provided as a JSON string.
// It also validates the group and checks if the group is consistent with the survey.
func (s *Survey) UpdateGroupJson(ug string) error {
	return s.UpdateGroupBytes([]byte(ug))
}

// UpdateGroupBytes updates a group in the survey given its representation as a byte array
// It also validates the group and checks if the group is consistent with the survey.
func (s *Survey) UpdateGroupBytes(ug []byte) error {
	// unmarshal group
	var pg *question.Group
	err := json.Unmarshal(ug, pg)
	if err != nil {
		return err
	}

	group, ok := s.Groups[pg.NameId]
	if !ok {
		return fmt.Errorf("group '%s' not found", pg.NameId)
	}

	// update group
	group.Title = pg.Title
	group.Description = pg.Description
	group.Visible = pg.Visible
	group.IsExternalSurvey = pg.IsExternalSurvey
	group.QuestionsIds = pg.QuestionsIds

	// check consistency
	return s.checkConsistency()
}

// AddOrUpdateGroupMap adds or updates a group in the survey given its representation as a map.
// It also validates the group and checks if the group is consistent with the survey.
func (s *Survey) AddOrUpdateGroupMap(g map[string]any) error {
	b, _ := json.Marshal(g)
	return s.AddOrUpdateGroupBytes(b)
}

// AddOrUpdateGroupJson adds or updates a group in the survey given its representation as a JSON string.
// It also validates the group and checks if the group is consistent with the survey.
func (s *Survey) AddOrUpdateGroupJson(g string) error {
	return s.AddOrUpdateGroupBytes([]byte(g))
}

// AddOrUpdateGroupBytes adds or updates a group in the survey given its representation as a byte array.
// It also validates the group and checks if the group is consistent with the survey.
func (s *Survey) AddOrUpdateGroupBytes(g []byte) error {
	// unmarshal group
	var pg *question.Group
	err := json.Unmarshal(g, pg)
	if err != nil {
		return err
	}

	// check if group already exists
	if _, ok := s.Groups[pg.NameId]; ok {
		return s.UpdateGroupBytes(g)
	}

	// add group to survey
	return s.addGroup(pg)
}

// RemoveGroup removes a group from the survey given its nameId.
// It also validates the group and checks if the group is consistent with the survey.
func (s *Survey) RemoveGroup(groupNameId string) error {
	// check if group exists
	if _, ok := s.Groups[groupNameId]; !ok {
		return fmt.Errorf("group '%s' not found", groupNameId)
	}

	// remove group from options groups
	for _, q := range s.Questions {
		if types.IsChoiceType(q.QTyp) {
			c, _ := choice.CastToChoice(q.Value)
			if removed := c.RemoveGroupId(groupNameId); removed {
				break
			}
		}
	}

	// remove group from groups order
	for i, id := range s.GroupsOrder {
		if id == groupNameId {
			s.GroupsOrder = append(s.GroupsOrder[:i], s.GroupsOrder[i+1:]...)
			break
		}
	}

	// remove group from survey
	delete(s.Groups, groupNameId)

	// check consistency
	return s.checkConsistency()
}

// AddQuestionToGroup adds a question to a group in the survey.
// Args:
// * questionNameId: the nameId of the question to add.
// * groupNameId: the nameId of the group to add the question to.
// * position: the position of the question in the group. If position is -1, the question will be added at the end of the group.
// It also validates the group and checks if the group is consistent with the survey.
func (s *Survey) AddQuestionToGroup(questionNameId, groupNameId string, position int) error {
	// check if question exists
	if _, ok := s.Questions[questionNameId]; !ok {
		return fmt.Errorf("question '%s' not found", questionNameId)
	}

	// check if group exists
	if _, ok := s.Groups[groupNameId]; !ok {
		return fmt.Errorf("group '%s' not found", groupNameId)
	}

	// add question to group
	s.Groups[groupNameId].AddQuestionId(questionNameId, position)

	// check consistency
	return s.checkConsistency()
}

// RemoveQuestionFromGroup removes a question from a group in the survey.
// Args:
// * questionNameId: the nameId of the question to remove.
// * groupNameId: the nameId of the group to remove the question from.
// It also validates the group and checks if the group is consistent with the survey.
func (s *Survey) RemoveQuestionFromGroup(questionNameId, groupNameId string) error {
	// check if question exists
	if _, ok := s.Questions[questionNameId]; !ok {
		return fmt.Errorf("question '%s' not found", questionNameId)
	}

	// check if group exists
	if _, ok := s.Groups[groupNameId]; !ok {
		return fmt.Errorf("group '%s' not found", groupNameId)
	}

	// remove question from group
	s.Groups[groupNameId].RemoveQuestionId(questionNameId)

	// check consistency
	return s.checkConsistency()
}

// UpdateGroupQuestions updates the questions of a group in the survey.
// Args:
// * groupNameId: the nameId of the group to update.
// * questionsIds: the list of questions ids to update the group with.
// It also validates the group and checks if the group is consistent with the survey.
func (s *Survey) UpdateGroupQuestions(groupNameId string, questionsIds []string) error {
	// check if group exists
	if _, ok := s.Groups[groupNameId]; !ok {
		return fmt.Errorf("group '%s' not found", groupNameId)
	}

	// check if questions exist
	for _, questionNameId := range questionsIds {
		if _, ok := s.Questions[questionNameId]; !ok {
			return fmt.Errorf("question '%s' not found", questionNameId)
		}
	}

	// update group questions
	s.Groups[groupNameId].QuestionsIds = questionsIds

	// check consistency
	return s.checkConsistency()
}

// addGroup adds a group to the survey.
// It also validates the group and checks if the group is consistent with the survey.
func (s *Survey) addGroup(pg *question.Group) error {
	// validate group
	if err := SurveyValidator.Struct(pg); err != nil {
		errs := TranslateValidationErrors(err)
		errs = append([]error{fmt.Errorf("error validating group")}, errs...)
		return errors.Join(errs...)
	}

	// check if group already exists
	if _, ok := s.Groups[pg.NameId]; ok {
		return fmt.Errorf("group nameId '%s' already exists", pg.NameId)
	}

	// add group to survey
	s.Groups[pg.NameId] = pg

	// check consistency
	return s.checkConsistency()
}
