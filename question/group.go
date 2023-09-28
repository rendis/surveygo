package question

// Group is a struct that represents a group of questions.
type Group struct {
	// NameId is the name id of the group.
	// Validations:
	// - required
	// - valid name id
	NameId string `json:"nameId" bson:"nameId" validate:"required,validNameId"`

	// Title is the title of the group.
	// Validations:
	// - optional
	Title *string `json:"title" bson:"title" validate:"omitempty"`

	// Description is the description of the group.
	// Validations:
	// - optional
	Description *string `json:"description" bson:"description" validate:"omitempty"`

	// Visible is a flag that indicates if the group is visible.
	Visible bool `json:"visible" bson:"visible"`

	// IsExternalSurvey is a flag that indicates if the group is an external survey.
	// When a group is an external survey, it means that:
	// - QuestionsIds length must be 1
	// - The question will be an external survey question id
	// Validations:
	// - validIfExternalSurvey
	IsExternalSurvey bool `json:"isExternalSurvey" bson:"isExternalSurvey" validate:"validIfExternalSurvey"`

	// QuestionsIds is a list of question ids that are associated with this group.
	// Validations:
	//	- required
	// 	- each question id must be valid:
	//		* length must be greater than 0
	QuestionsIds []string `json:"questionsIds" bson:"questionsIds" validate:"required,dive,min=1"`
}

// RemoveQuestionId removes the question with the specified name ID from the group.
func (g *Group) RemoveQuestionId(nameId string) {
	for i, id := range g.QuestionsIds {
		if id == nameId {
			g.QuestionsIds = append(g.QuestionsIds[:i], g.QuestionsIds[i+1:]...)
			break
		}
	}
}

// AddQuestionId adds the question with the specified name ID to the group.
func (g *Group) AddQuestionId(nameId string) {
	// check if question already exists
	for _, id := range g.QuestionsIds {
		if id == nameId {
			return
		}
	}
	g.QuestionsIds = append(g.QuestionsIds, nameId)
}
