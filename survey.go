package surveygo

import (
	"github.com/rendis/surveygo/v2/question"
)

// Answers is a map with the answers provided by the user.
// The key is the question NameId (Question.NameId).
type Answers map[string][]any

// Survey is a struct representation of a survey.
type Survey struct {
	// Title is the title of the survey.
	// Validations:
	//	- required
	//	- min length: 1
	Title string `json:"title,omitempty" bson:"title,omitempty" validate:"required,min=1"`

	// Version is the version of the survey.
	// Validations:
	//	- required
	//	- min length: 1
	Version string `json:"version,omitempty" bson:"version,omitempty" validate:"required,min=1"`

	// Description is the description of the survey.
	// Validations:
	//	- optional
	//	- min length: 1
	Description *string `json:"description,omitempty" bson:"description,omitempty" validate:"omitempty"`

	// Questions is a map with all the questions in the survey.
	// The key is the question NameId (Question.NameId).
	// Validations:
	//	- required
	//	- min length: 1
	//	- each question must be valid
	Questions map[string]*question.Question `json:"questions,omitempty" bson:"questions,omitempty" validate:"required,dive"`

	// Groups is a map with all the groups in the survey.
	// The key is the group NameId (Group.NameId).
	// Validations:
	//	- required
	//	- min length: 1
	//	- each group must be valid
	Groups map[string]*question.Group `json:"groups,omitempty" bson:"groups,omitempty" validate:"required,dive"`

	// GroupsOrder is a list of group name ids that defines the order of the groups in the survey.
	// Validations:
	//	- required
	//	- min length: 1
	GroupsOrder []string `json:"groupsOrder,omitempty" bson:"groupsOrder,omitempty" validate:"required"`
}
