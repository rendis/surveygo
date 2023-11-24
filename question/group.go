package question

// Group is a struct that represents a group of questions.
type Group struct {
	// NameId is the name id of the group.
	// Validations:
	// - required
	// - valid name id
	NameId string `json:"nameId,omitempty" bson:"nameId,omitempty" validate:"required,validNameId"`

	// Title is the title of the group.
	// Validations:
	// - optional
	Title *string `json:"title,omitempty" bson:"title,omitempty" validate:"omitempty"`

	// Description is the description of the group.
	// Validations:
	// - optional
	Description *string `json:"description,omitempty" bson:"description,omitempty" validate:"omitempty"`

	// Hidden is a flag that indicates if the group is hidden (Default: false).
	// Validations:
	// - optional
	Hidden bool `json:"hidden,omitempty" bson:"hidden,omitempty" validate:"omitempty"`

	// Disabled is a flag that indicates if the group is disabled (Default: false).
	// Validations:
	// - optional
	Disabled bool `json:"disabled,omitempty" bson:"disabled,omitempty" validate:"omitempty"`

	// IsExternalSurvey is a flag that indicates if the group is an external survey.
	// When a group is an external survey, it means that:
	// - QuestionsIds length must be 1
	// - The question will be an external survey question id
	// Validations:
	// - validIfExternalSurvey
	IsExternalSurvey bool `json:"isExternalSurvey,omitempty" bson:"isExternalSurvey,omitempty" validate:"validIfExternalSurvey"`

	// AllowRepeat is a flag that indicates if the group can be repeated (Default: false).
	// Validations:
	// - optional
	AllowRepeat bool `json:"allowRepeat,omitempty" bson:"allowRepeat,omitempty" validate:"omitempty"`

	// QuestionsIds is a list of question ids that are associated with this group or the external survey id.
	// Validations:
	//	- required
	// 	- each id must be valid:
	//		* length must be greater than 0
	QuestionsIds []string `json:"questionsIds,omitempty" bson:"questionsIds,omitempty" validate:"required,dive,min=1"`

	// Metadata is a map with additional metadata for the group.
	// Validations:
	// - optional
	Metadata map[string]any `json:"metadata,omitempty" bson:"metadata,omitempty" validate:"omitempty"`
}

// RemoveQuestionId removes the question with the specified name ID from the group.
// Returns true if the question was removed, false otherwise.
func (g *Group) RemoveQuestionId(nameId string) bool {
	for i, id := range g.QuestionsIds {
		if id == nameId {
			g.QuestionsIds = append(g.QuestionsIds[:i], g.QuestionsIds[i+1:]...)
			return true
		}
	}
	return false
}

// AddQuestionId adds the question with the specified name ID to the group.
// Args:
// - nameId: the name ID of the question to add.
// - position: the position where to add the question. If position is -1, the question will be added at the end of the group.
func (g *Group) AddQuestionId(nameId string, position int) {
	// check if question already exists in the group
	for _, id := range g.QuestionsIds {
		if id == nameId {
			return
		}
	}

	// add question at the end of the group
	if position < 0 || position >= len(g.QuestionsIds) {
		g.QuestionsIds = append(g.QuestionsIds, nameId)
		return
	}

	// add question at the specified position
	g.QuestionsIds = append(g.QuestionsIds, "")
	copy(g.QuestionsIds[position+1:], g.QuestionsIds[position:])
	g.QuestionsIds[position] = nameId
}
