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
	Title string `json:"title" bson:"title" validate:"required,min=1"`

	// Version is the version of the survey.
	// Validations:
	//	- required
	//	- min length: 1
	Version string `json:"version" bson:"version" validate:"required,min=1"`

	// Description is the description of the survey.
	// Validations:
	//	- optional
	//	- min length: 1
	Description *string `json:"description" bson:"description" validate:"omitempty"`

	// Questions is a map with all the questions in the survey.
	// The key is the question NameId (Question.NameId).
	// Validations:
	//	- required
	//	- min length: 1
	//	- each question must be valid
	Questions map[string]*question.Question `json:"questions" bson:"questions" validate:"required,dive"`

	// Groups is a map with all the groups in the survey.
	// The key is the group NameId (Group.NameId).
	// Validations:
	//	- required
	//	- min length: 1
	//	- each group must be valid
	Groups map[string]*question.Group `json:"groups" bson:"groups" validate:"required,dive"`

	// GroupsOrder is a list of group name ids that defines the order of the groups in the survey.
	// Validations:
	//	- required
	//	- min length: 1
	GroupsOrder []string `json:"groupsOrder" bson:"groupsOrder" validate:"required"`
}
