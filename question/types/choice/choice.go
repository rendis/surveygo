package choice

import (
	"fmt"
	"github.com/rendis/surveygo/v2/question/types"
)

// Choice represents a choice question type.
// Types:
// - types.QTypeSingleSelect
// - types.QTypeMultipleSelect
// - types.QTypeRadio
// - types.QTypeCheckbox
type Choice struct {
	types.QBase `json:",inline" bson:",inline"`

	// Options is a list of options for the choice field.
	// Validations:
	// - required
	// - at least one option
	// - each option must be valid
	Options []*Option `json:"options,omitempty" bson:"options,omitempty" validate:"required,min=1,dive"`
}

// Option represents a single option in a choice widget.
type Option struct {
	// NameId is the identifier of the option.
	// Validations:
	// - required
	// - valid name id
	NameId string `json:"nameId" bson:"nameId" validate:"required,validNameId"`

	// Label is a label for the option.
	// Validations:
	// - required
	// - min length: 1
	Label string `json:"label,omitempty" bson:"label,omitempty" validate:"required,min=1"`

	// Value is the value of the option.
	// Validations:
	// - optional
	// - min length: 1
	Value any `json:"value,omitempty" bson:"value,omitempty" validate:"omitempty,min=1"`

	// GroupsIds is a list of group ids that are associated with this option.
	// Validations:
	// - optional
	GroupsIds []string `json:"groupsIds,omitempty" bson:"groupsIds,omitempty" validate:"omitempty"`

	// Metadata is a map of metadata for the option.
	// Validations:
	// - optional
	Metadata map[string]any `json:"metadata,omitempty" bson:"metadata,omitempty" validate:"omitempty"`
}

// GetOptionsGroups returns a map with each option and its associated groups.
// Key: Option name id
// Value: Groups ids
func (c *Choice) GetOptionsGroups() map[string][]string {
	var dependencies = map[string][]string{}
	for _, option := range c.Options {
		dependencies[option.NameId] = option.GroupsIds
	}
	return dependencies
}

// RemoveGroupId removes the group with the specified name ID from the choice.
// Returns true if the group was removed, false otherwise.
func (c *Choice) RemoveGroupId(groupId string) bool {
	for _, option := range c.Options {
		for i, id := range option.GroupsIds {
			if id == groupId {
				option.GroupsIds = append(option.GroupsIds[:i], option.GroupsIds[i+1:]...)
				return true
			}
		}
	}
	return false
}

// CastToChoice casts the given interface to a Choice type.
func CastToChoice(i any) (*Choice, error) {
	c, ok := i.(*Choice)
	if !ok {
		return nil, fmt.Errorf("invalid type, expected *choice.Choice, got %T", i)
	}
	return c, nil
}
