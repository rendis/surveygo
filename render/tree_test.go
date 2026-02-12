package render

import (
	"strings"
	"testing"

	surveygo "github.com/rendis/surveygo/v2"
	"github.com/rendis/surveygo/v2/question"
	"github.com/rendis/surveygo/v2/question/types"
	"github.com/rendis/surveygo/v2/question/types/choice"
	"github.com/rendis/surveygo/v2/question/types/text"
)

func strPtr(s string) *string { return &s }

func TestBuildGroupTree_DirectCycle(t *testing.T) {
	// grp-a -> grp-b -> grp-a (cycle via groupsOrder)
	survey := &surveygo.Survey{
		NameId:      "test",
		Title:       "Test",
		Version:     "1",
		GroupsOrder: []string{"grp-a"},
		Groups: map[string]*question.Group{
			"grp-a": {NameId: "grp-a", GroupsOrder: []string{"grp-b"}},
			"grp-b": {NameId: "grp-b", GroupsOrder: []string{"grp-a"}},
		},
		Questions: map[string]*question.Question{},
	}

	_, err := buildGroupTree(survey)
	if err == nil {
		t.Fatal("expected cycle error, got nil")
	}
	if !strings.Contains(err.Error(), "cycle detected") {
		t.Fatalf("expected 'cycle detected' in error, got: %s", err)
	}
}

func TestBuildGroupTree_SelfCycle(t *testing.T) {
	survey := &surveygo.Survey{
		NameId:      "test",
		Title:       "Test",
		Version:     "1",
		GroupsOrder: []string{"grp-a"},
		Groups: map[string]*question.Group{
			"grp-a": {NameId: "grp-a", GroupsOrder: []string{"grp-a"}},
		},
		Questions: map[string]*question.Question{},
	}

	_, err := buildGroupTree(survey)
	if err == nil {
		t.Fatal("expected cycle error, got nil")
	}
	if !strings.Contains(err.Error(), "cycle detected") {
		t.Fatalf("expected 'cycle detected' in error, got: %s", err)
	}
}

func TestBuildGroupTree_CycleViaChoiceGroupsIds(t *testing.T) {
	survey := &surveygo.Survey{
		NameId:      "test",
		Title:       "Test",
		Version:     "1",
		GroupsOrder: []string{"grp-root"},
		Groups: map[string]*question.Group{
			"grp-root": {NameId: "grp-root", GroupsOrder: []string{"grp-leaf"}},
			"grp-leaf": {NameId: "grp-leaf", QuestionsIds: []string{"q1"}},
		},
		Questions: map[string]*question.Question{
			"q1": {
				BaseQuestion: question.BaseQuestion{
					NameId:  "q1",
					QTyp:    types.QTypeSingleSelect,
					Label:   "Q1",
					Visible: true,
				},
				Value: &choice.Choice{
					Options: []*choice.Option{
						{NameId: "opt1", Label: "Opt 1", GroupsIds: []string{"grp-root"}},
					},
				},
			},
		},
	}

	_, err := buildGroupTree(survey)
	if err == nil {
		t.Fatal("expected cycle error, got nil")
	}
	if !strings.Contains(err.Error(), "cycle detected") {
		t.Fatalf("expected 'cycle detected' in error, got: %s", err)
	}
}

func TestBuildGroupTree_DeepCycle(t *testing.T) {
	survey := &surveygo.Survey{
		NameId:      "test",
		Title:       "Test",
		Version:     "1",
		GroupsOrder: []string{"grp-a"},
		Groups: map[string]*question.Group{
			"grp-a": {NameId: "grp-a", GroupsOrder: []string{"grp-b"}},
			"grp-b": {NameId: "grp-b", GroupsOrder: []string{"grp-c"}},
			"grp-c": {NameId: "grp-c", GroupsOrder: []string{"grp-d"}},
			"grp-d": {NameId: "grp-d", GroupsOrder: []string{"grp-b"}},
		},
		Questions: map[string]*question.Question{},
	}

	_, err := buildGroupTree(survey)
	if err == nil {
		t.Fatal("expected cycle error, got nil")
	}
	if !strings.Contains(err.Error(), "cycle detected") {
		t.Fatalf("expected 'cycle detected' in error, got: %s", err)
	}
}

func TestBuildGroupTree_NoCycle_SharedAcrossBranches(t *testing.T) {
	survey := &surveygo.Survey{
		NameId:      "test",
		Title:       "Test",
		Version:     "1",
		GroupsOrder: []string{"grp-root"},
		Groups: map[string]*question.Group{
			"grp-root":   {NameId: "grp-root", GroupsOrder: []string{"grp-a", "grp-b"}},
			"grp-a":      {NameId: "grp-a", QuestionsIds: []string{"qa"}},
			"grp-b":      {NameId: "grp-b", QuestionsIds: []string{"qb"}},
			"grp-shared": {NameId: "grp-shared", QuestionsIds: []string{"qs"}},
		},
		Questions: map[string]*question.Question{
			"qa": {
				BaseQuestion: question.BaseQuestion{
					NameId: "qa", QTyp: types.QTypeSingleSelect, Label: "QA", Visible: true,
				},
				Value: &choice.Choice{
					Options: []*choice.Option{{NameId: "oa", Label: "OA", GroupsIds: []string{"grp-shared"}}},
				},
			},
			"qb": {
				BaseQuestion: question.BaseQuestion{
					NameId: "qb", QTyp: types.QTypeRadio, Label: "QB", Visible: true,
				},
				Value: &choice.Choice{
					Options: []*choice.Option{{NameId: "ob", Label: "OB", GroupsIds: []string{"grp-shared"}}},
				},
			},
			"qs": {
				BaseQuestion: question.BaseQuestion{
					NameId: "qs", QTyp: types.QTypeInputText, Label: "QS", Visible: true,
				},
				Value: &text.FreeText{},
			},
		},
	}

	tree, err := buildGroupTree(survey)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	root := tree.Roots[0]
	if len(root.Children) != 2 {
		t.Fatalf("expected 2 children of root, got %d", len(root.Children))
	}
	for _, child := range root.Children {
		if len(child.Children) != 1 {
			t.Fatalf("expected 1 child of %s, got %d", child.NameId, len(child.Children))
		}
		if child.Children[0].NameId != "grp-shared" {
			t.Fatalf("expected grp-shared as child of %s, got %s", child.NameId, child.Children[0].NameId)
		}
	}
}

func TestBuildGroupTree_MissingGroup(t *testing.T) {
	survey := &surveygo.Survey{
		NameId:      "test",
		Title:       "Test",
		Version:     "1",
		GroupsOrder: []string{"grp-a"},
		Groups: map[string]*question.Group{
			"grp-a": {NameId: "grp-a", GroupsOrder: []string{"grp-missing"}},
		},
		Questions: map[string]*question.Question{},
	}

	_, err := buildGroupTree(survey)
	if err == nil {
		t.Fatal("expected error for missing group, got nil")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected 'not found' in error, got: %s", err)
	}
}

func TestBuildGroupTree_RepeatDescendants(t *testing.T) {
	// Tree:
	//   grp-root (not repeat)
	//     grp-a (repeat)
	//       grp-a1 (repeat)
	//       grp-a2 (not repeat)
	//     grp-b (not repeat)
	//       grp-b1 (repeat)
	//
	// Expected RepeatDescendants:
	//   grp-root = 3 (grp-a + grp-a1 + grp-b1)
	//   grp-a    = 1 (grp-a1)
	//   grp-a1   = 0
	//   grp-a2   = 0
	//   grp-b    = 1 (grp-b1)
	//   grp-b1   = 0
	survey := &surveygo.Survey{
		NameId:      "test",
		Title:       "Test",
		Version:     "1",
		GroupsOrder: []string{"grp-root"},
		Groups: map[string]*question.Group{
			"grp-root": {NameId: "grp-root", GroupsOrder: []string{"grp-a", "grp-b"}},
			"grp-a":    {NameId: "grp-a", AllowRepeat: true, GroupsOrder: []string{"grp-a1", "grp-a2"}},
			"grp-a1":   {NameId: "grp-a1", AllowRepeat: true},
			"grp-a2":   {NameId: "grp-a2"},
			"grp-b":    {NameId: "grp-b", GroupsOrder: []string{"grp-b1"}},
			"grp-b1":   {NameId: "grp-b1", AllowRepeat: true},
		},
		Questions: map[string]*question.Question{},
	}

	tree, err := buildGroupTree(survey)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	want := map[string]int{
		"grp-root": 3,
		"grp-a":    1,
		"grp-a1":   0,
		"grp-a2":   0,
		"grp-b":    1,
		"grp-b1":   0,
	}

	for nameId, expected := range want {
		node, ok := tree.Index[nameId]
		if !ok {
			t.Fatalf("node %q not found in index", nameId)
		}
		if node.RepeatDescendants != expected {
			t.Errorf("%s: RepeatDescendants = %d, want %d", nameId, node.RepeatDescendants, expected)
		}
	}
}

func TestBuildGroupTree_MissingQuestion(t *testing.T) {
	survey := &surveygo.Survey{
		NameId:      "test",
		Title:       "Test",
		Version:     "1",
		GroupsOrder: []string{"grp-a"},
		Groups: map[string]*question.Group{
			"grp-a": {NameId: "grp-a", QuestionsIds: []string{"q-missing"}},
		},
		Questions: map[string]*question.Question{},
	}

	_, err := buildGroupTree(survey)
	if err == nil {
		t.Fatal("expected error for missing question, got nil")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected 'not found' in error, got: %s", err)
	}
}
