package render

import (
	"fmt"
	"strings"

	surveygo "github.com/rendis/surveygo/v2"
	"github.com/rendis/surveygo/v2/question/types"
	"github.com/rendis/surveygo/v2/question/types/choice"
)

func buildGroupTree(survey *surveygo.Survey) (*GroupTree, error) {
	tree := &GroupTree{
		Index:  make(map[string]*GroupNode),
		parent: make(map[string]string),
	}

	visiting := make(map[string]bool) // cycle detection: ancestors in current DFS path

	for _, rootID := range survey.GroupsOrder {
		node, err := buildNode(rootID, visiting, survey, tree)
		if err != nil {
			return nil, fmt.Errorf("building root group %q: %w", rootID, err)
		}
		tree.Roots = append(tree.Roots, node)
	}

	return tree, nil
}

func buildNode(nameID string, visiting map[string]bool, survey *surveygo.Survey, tree *GroupTree) (*GroupNode, error) {
	if visiting[nameID] {
		return nil, fmt.Errorf("cycle detected: %s (path: %s)", nameID, visitingPath(visiting))
	}

	grp, ok := survey.Groups[nameID]
	if !ok {
		return nil, fmt.Errorf("group %q not found in survey.Groups", nameID)
	}

	visiting[nameID] = true
	defer delete(visiting, nameID)

	node := &GroupNode{
		NameId:      grp.NameId,
		AllowRepeat: grp.AllowRepeat,
	}

	// Container group: has sub-groups
	for _, childID := range grp.GroupsOrder {
		child, err := buildNode(childID, visiting, survey, tree)
		if err != nil {
			return nil, fmt.Errorf("in group %q: %w", nameID, err)
		}
		node.Children = append(node.Children, child)
	}

	// Leaf group: check choice questions for GroupsIds
	for _, qID := range grp.QuestionsIds {
		q, ok := survey.Questions[qID]
		if !ok {
			return nil, fmt.Errorf("question %q referenced by group %q not found", qID, nameID)
		}
		if !types.IsSimpleChoiceType(q.QTyp) {
			continue
		}
		c, err := choice.CastToChoice(q.Value)
		if err != nil {
			continue
		}
		for _, opt := range c.Options {
			for _, gID := range opt.GroupsIds {
				child, err := buildNode(gID, visiting, survey, tree)
				if err != nil {
					return nil, fmt.Errorf("in choice option %q of question %q: %w", opt.NameId, qID, err)
				}
				node.Children = append(node.Children, child)
			}
		}
	}

	for _, c := range node.Children {
		tree.parent[c.NameId] = nameID
		node.RepeatDescendants += c.RepeatDescendants
		if c.AllowRepeat {
			node.RepeatDescendants++
		}
	}

	tree.Index[nameID] = node
	return node, nil
}

func visitingPath(visiting map[string]bool) string {
	path := make([]string, 0, len(visiting))
	for k := range visiting {
		path = append(path, k)
	}
	return strings.Join(path, " -> ")
}
