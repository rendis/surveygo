package render

import (
	"testing"

	surveygo "github.com/rendis/surveygo/v2"
	"github.com/rendis/surveygo/v2/question"
	"github.com/rendis/surveygo/v2/question/types"
	"github.com/rendis/surveygo/v2/question/types/choice"
	"github.com/rendis/surveygo/v2/question/types/text"
)

// TestBuildSections_TruthTable verifies all four branches of the section-type
// truth table:
//
//	AllowRepeat=false                             → group
//	AllowRepeat=true, RepeatDescendants > 0       → repeat-list
//	AllowRepeat=true, RepeatDescendants=0, con ms → repeat-list
//	AllowRepeat=true, RepeatDescendants=0, sin ms → repeat-table (flattened)
func TestBuildSections_TruthTable(t *testing.T) {
	tests := []struct {
		name       string
		survey     *surveygo.Survey
		answers    surveygo.Answers
		wantType   string
		wantCols   []string // only for repeat-table: expected column nameIds
		wantFields []string // only for group: expected field nameIds
	}{
		{
			name:     "flat_group",
			survey:   flatGroupSurvey(),
			answers:  surveygo.Answers{"q-name": {"Juan"}, "q-email": {"j@mail.com"}},
			wantType: "group",
			wantFields: []string{"q-name", "q-email"},
		},
		{
			name:     "repeat_with_repeat_descendants",
			survey:   repeatWithRepeatDescendantsSurvey(),
			answers:  surveygo.Answers{},
			wantType: "repeat-list",
		},
		{
			name:     "repeat_with_multi_select",
			survey:   repeatWithMultiSelectSurvey(),
			answers:  surveygo.Answers{},
			wantType: "repeat-list",
		},
		{
			name:    "repeat_no_ms_flattened_table",
			survey:  repeatNoMsFlattenedSurvey(),
			answers: surveygo.Answers{
				"grp-person": {
					map[string]any{
						"q-first": []any{"Juan"},
						"q-last":  []any{"Perez"},
						"q-email": []any{"j@mail.com"},
						"q-phone": []any{"+56", "912345678"},
					},
					map[string]any{
						"q-first": []any{"Maria"},
						"q-last":  []any{"Lopez"},
						"q-email": []any{"m@mail.com"},
						"q-phone": []any{"+56", "987654321"},
					},
				},
			},
			wantType: "repeat-table",
			wantCols: []string{"q-first", "q-last", "q-email", "q-phone"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			card, err := AnswersToJSON(tt.survey, tt.answers)
			if err != nil {
				t.Fatalf("AnswersToJSON: %v", err)
			}

			if len(card.Sections) == 0 {
				t.Fatal("expected at least 1 section, got 0")
			}
			sec := card.Sections[0]

			if sec.Type != tt.wantType {
				t.Errorf("section type = %q, want %q", sec.Type, tt.wantType)
			}

			if tt.wantCols != nil {
				if len(sec.Columns) != len(tt.wantCols) {
					t.Fatalf("columns count = %d, want %d", len(sec.Columns), len(tt.wantCols))
				}
				for i, want := range tt.wantCols {
					if sec.Columns[i].NameId != want {
						t.Errorf("column[%d].NameId = %q, want %q", i, sec.Columns[i].NameId, want)
					}
				}
				if len(sec.Rows) != 2 {
					t.Errorf("rows count = %d, want 2", len(sec.Rows))
				}
			}

			if tt.wantFields != nil {
				if len(sec.Fields) != len(tt.wantFields) {
					t.Fatalf("fields count = %d, want %d", len(sec.Fields), len(tt.wantFields))
				}
				for i, want := range tt.wantFields {
					if sec.Fields[i].NameId != want {
						t.Errorf("field[%d].NameId = %q, want %q", i, sec.Fields[i].NameId, want)
					}
				}
			}
		})
	}
}

// --- Survey builders ---

// flatGroupSurvey: AllowRepeat=false → group
func flatGroupSurvey() *surveygo.Survey {
	return &surveygo.Survey{
		NameId: "s-flat", Title: "Flat", Version: "1",
		GroupsOrder: []string{"grp-main"},
		Groups: map[string]*question.Group{
			"grp-main": {NameId: "grp-main", Title: strPtr("Main"), QuestionsIds: []string{"q-name", "q-email"}},
		},
		Questions: map[string]*question.Question{
			"q-name":  textQuestion("q-name", "Nombre"),
			"q-email": emailQuestion("q-email", "Email"),
		},
	}
}

// repeatWithRepeatDescendantsSurvey: AllowRepeat=true, RepeatDescendants > 0 → repeat-list
func repeatWithRepeatDescendantsSurvey() *surveygo.Survey {
	return &surveygo.Survey{
		NameId: "s-nested-repeat", Title: "Nested Repeat", Version: "1",
		GroupsOrder: []string{"grp-outer"},
		Groups: map[string]*question.Group{
			"grp-outer": {NameId: "grp-outer", Title: strPtr("Outer"), AllowRepeat: true, GroupsOrder: []string{"grp-inner"}},
			"grp-inner": {NameId: "grp-inner", Title: strPtr("Inner"), AllowRepeat: true, QuestionsIds: []string{"q-val"}},
		},
		Questions: map[string]*question.Question{
			"q-val": textQuestion("q-val", "Value"),
		},
	}
}

// repeatWithMultiSelectSurvey: AllowRepeat=true, RepeatDescendants=0, con multi_select → repeat-list
func repeatWithMultiSelectSurvey() *surveygo.Survey {
	return &surveygo.Survey{
		NameId: "s-repeat-ms", Title: "Repeat MS", Version: "1",
		GroupsOrder: []string{"grp-main"},
		Groups: map[string]*question.Group{
			"grp-main": {NameId: "grp-main", Title: strPtr("Main"), AllowRepeat: true, GroupsOrder: []string{"grp-child"}},
			"grp-child": {NameId: "grp-child", Title: strPtr("Child"), QuestionsIds: []string{"q-name", "q-perms"}},
		},
		Questions: map[string]*question.Question{
			"q-name": textQuestion("q-name", "Nombre"),
			"q-perms": {
				BaseQuestion: question.BaseQuestion{NameId: "q-perms", QTyp: types.QTypeMultipleSelect, Label: "Permisos", Visible: true},
				Value: &choice.Choice{Options: []*choice.Option{
					{NameId: "opt-read", Label: "Leer"},
					{NameId: "opt-write", Label: "Escribir"},
				}},
			},
		},
	}
}

// repeatNoMsFlattenedSurvey: AllowRepeat=true, RepeatDescendants=0, sin multi_select → repeat-table
// Has non-repeatable children whose questions should be flattened into columns.
func repeatNoMsFlattenedSurvey() *surveygo.Survey {
	return &surveygo.Survey{
		NameId: "s-repeat-flat", Title: "Repeat Flat", Version: "1",
		GroupsOrder: []string{"grp-person"},
		Groups: map[string]*question.Group{
			"grp-person":  {NameId: "grp-person", Title: strPtr("Persona"), AllowRepeat: true, GroupsOrder: []string{"grp-name", "grp-contact"}},
			"grp-name":    {NameId: "grp-name", Title: strPtr("Nombre"), QuestionsIds: []string{"q-first", "q-last"}},
			"grp-contact": {NameId: "grp-contact", Title: strPtr("Contacto"), QuestionsIds: []string{"q-email", "q-phone"}},
		},
		Questions: map[string]*question.Question{
			"q-first": textQuestion("q-first", "Nombre"),
			"q-last":  textQuestion("q-last", "Apellido"),
			"q-email": emailQuestion("q-email", "Email"),
			"q-phone": phoneQuestion("q-phone", "Teléfono"),
		},
	}
}

// --- Question helpers ---

func textQuestion(nameId, label string) *question.Question {
	return &question.Question{
		BaseQuestion: question.BaseQuestion{NameId: nameId, QTyp: types.QTypeInputText, Label: label, Visible: true},
		Value:        &text.FreeText{},
	}
}

func emailQuestion(nameId, label string) *question.Question {
	return &question.Question{
		BaseQuestion: question.BaseQuestion{NameId: nameId, QTyp: types.QTypeEmail, Label: label, Visible: true},
		Value:        &text.FreeText{},
	}
}

func phoneQuestion(nameId, label string) *question.Question {
	return &question.Question{
		BaseQuestion: question.BaseQuestion{NameId: nameId, QTyp: types.QTypeTelephone, Label: label, Visible: true},
		Value:        &text.Telephone{},
	}
}
