package render

import (
	"encoding/csv"
	"strings"
	"testing"

	surveygo "github.com/rendis/surveygo/v2"
	"github.com/rendis/surveygo/v2/question"
	"github.com/rendis/surveygo/v2/question/types"
	"github.com/rendis/surveygo/v2/question/types/choice"
)

// newMultiSelectSurvey builds a minimal survey with one multi-select question.
func newMultiSelectSurvey() *surveygo.Survey {
	return &surveygo.Survey{
		NameId:      "test-survey",
		Title:       "Test",
		Version:     "1",
		GroupsOrder: []string{"grp-main"},
		Groups: map[string]*question.Group{
			"grp-main": {
				NameId:       "grp-main",
				Title:        strPtr("Main"),
				QuestionsIds: []string{"q-colors"},
			},
		},
		Questions: map[string]*question.Question{
			"q-colors": {
				BaseQuestion: question.BaseQuestion{
					NameId:  "q-colors",
					QTyp:    types.QTypeMultipleSelect,
					Label:   "Colors",
					Visible: true,
				},
				Value: &choice.Choice{
					Options: []*choice.Option{
						{NameId: "red", Label: "Red"},
						{NameId: "green", Label: "Green"},
						{NameId: "blue", Label: "Blue"},
					},
				},
			},
		},
	}
}

// newToggleSurvey builds a minimal survey with one toggle question.
func newToggleSurvey() *surveygo.Survey {
	return &surveygo.Survey{
		NameId:      "test-survey",
		Title:       "Test",
		Version:     "1",
		GroupsOrder: []string{"grp-main"},
		Groups: map[string]*question.Group{
			"grp-main": {
				NameId:       "grp-main",
				Title:        strPtr("Main"),
				QuestionsIds: []string{"q-agree"},
			},
		},
		Questions: map[string]*question.Question{
			"q-agree": {
				BaseQuestion: question.BaseQuestion{
					NameId:  "q-agree",
					QTyp:    types.QTypeToggle,
					Label:   "Agree",
					Visible: true,
				},
				Value: &choice.Toggle{
					OnLabel:  "Yes",
					OffLabel: "No",
				},
			},
		},
	}
}

// parseCSV parses CSV bytes into a header row and data rows.
func parseCSV(t *testing.T, data []byte) ([]string, [][]string) {
	t.Helper()
	r := csv.NewReader(strings.NewReader(string(data)))
	records, err := r.ReadAll()
	if err != nil {
		t.Fatalf("parsing CSV: %v", err)
	}
	if len(records) < 1 {
		t.Fatal("expected at least a header row")
	}
	return records[0], records[1:]
}

func TestCSV_MultiSelect_DefaultCheckMark(t *testing.T) {
	s := newMultiSelectSurvey()
	answers := surveygo.Answers{
		"q-colors": {"red", "blue"},
	}

	data, err := AnswersToCSV(s, answers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	headers, rows := parseCSV(t, data)

	// Expect 3 option columns.
	if len(headers) != 3 {
		t.Fatalf("expected 3 columns, got %d: %v", len(headers), headers)
	}

	expectedHeaders := []string{"Colors - Red", "Colors - Green", "Colors - Blue"}
	for i, h := range expectedHeaders {
		if headers[i] != h {
			t.Errorf("header[%d]: expected %q, got %q", i, h, headers[i])
		}
	}

	if len(rows) != 1 {
		t.Fatalf("expected 1 data row, got %d", len(rows))
	}

	// Default: "true"/"false"
	expected := []string{"true", "false", "true"}
	for i, v := range expected {
		if rows[0][i] != v {
			t.Errorf("row[0][%d]: expected %q, got %q", i, v, rows[0][i])
		}
	}
}

func TestCSV_MultiSelect_CustomCheckMark(t *testing.T) {
	s := newMultiSelectSurvey()
	answers := surveygo.Answers{
		"q-colors": {"green"},
	}

	cm := &CheckMark{Selected: "x", NotSelected: ""}
	data, err := AnswersToCSV(s, answers, cm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, rows := parseCSV(t, data)
	if len(rows) != 1 {
		t.Fatalf("expected 1 data row, got %d", len(rows))
	}

	expected := []string{"", "x", ""}
	for i, v := range expected {
		if rows[0][i] != v {
			t.Errorf("row[0][%d]: expected %q, got %q", i, v, rows[0][i])
		}
	}
}

func TestCSV_MultiSelect_EmojiCheckMark(t *testing.T) {
	s := newMultiSelectSurvey()
	answers := surveygo.Answers{
		"q-colors": {"red", "green", "blue"},
	}

	cm := &CheckMark{Selected: "\u2705", NotSelected: "-"}
	data, err := AnswersToCSV(s, answers, cm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, rows := parseCSV(t, data)
	expected := []string{"\u2705", "\u2705", "\u2705"}
	for i, v := range expected {
		if rows[0][i] != v {
			t.Errorf("row[0][%d]: expected %q, got %q", i, v, rows[0][i])
		}
	}
}

func TestCSV_MultiSelect_NoneSelected(t *testing.T) {
	s := newMultiSelectSurvey()
	answers := surveygo.Answers{}

	cm := &CheckMark{Selected: "Yes", NotSelected: "No"}
	data, err := AnswersToCSV(s, answers, cm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, rows := parseCSV(t, data)
	expected := []string{"No", "No", "No"}
	for i, v := range expected {
		if rows[0][i] != v {
			t.Errorf("row[0][%d]: expected %q, got %q", i, v, rows[0][i])
		}
	}
}

func TestCSV_Toggle_DefaultCheckMark(t *testing.T) {
	s := newToggleSurvey()
	answers := surveygo.Answers{
		"q-agree": {true},
	}

	data, err := AnswersToCSV(s, answers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, rows := parseCSV(t, data)
	if rows[0][0] != "true" {
		t.Errorf("expected %q, got %q", "true", rows[0][0])
	}
}

func TestCSV_Toggle_CustomCheckMark(t *testing.T) {
	s := newToggleSurvey()
	answers := surveygo.Answers{
		"q-agree": {true},
	}

	cm := &CheckMark{Selected: "Si", NotSelected: "No"}
	data, err := AnswersToCSV(s, answers, cm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, rows := parseCSV(t, data)
	if rows[0][0] != "Si" {
		t.Errorf("expected %q, got %q", "Si", rows[0][0])
	}
}

func TestCSV_Toggle_FalseCustomCheckMark(t *testing.T) {
	s := newToggleSurvey()
	answers := surveygo.Answers{
		"q-agree": {false},
	}

	cm := &CheckMark{Selected: "Si", NotSelected: "No"}
	data, err := AnswersToCSV(s, answers, cm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, rows := parseCSV(t, data)
	if rows[0][0] != "No" {
		t.Errorf("expected %q, got %q", "No", rows[0][0])
	}
}

func TestCSV_AnswersTo_WithCheckMark(t *testing.T) {
	s := newMultiSelectSurvey()
	answers := surveygo.Answers{
		"q-colors": {"red"},
	}

	result, err := AnswersTo(s, answers, OutputOptions{
		CSV:       true,
		CheckMark: &CheckMark{Selected: "1", NotSelected: "0"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, rows := parseCSV(t, result.CSV)
	expected := []string{"1", "0", "0"}
	for i, v := range expected {
		if rows[0][i] != v {
			t.Errorf("row[0][%d]: expected %q, got %q", i, v, rows[0][i])
		}
	}
}

func TestCSV_Checkbox_CustomCheckMark(t *testing.T) {
	s := &surveygo.Survey{
		NameId:      "test-survey",
		Title:       "Test",
		Version:     "1",
		GroupsOrder: []string{"grp-main"},
		Groups: map[string]*question.Group{
			"grp-main": {
				NameId:       "grp-main",
				Title:        strPtr("Main"),
				QuestionsIds: []string{"q-features"},
			},
		},
		Questions: map[string]*question.Question{
			"q-features": {
				BaseQuestion: question.BaseQuestion{
					NameId:  "q-features",
					QTyp:    types.QTypeCheckbox,
					Label:   "Features",
					Visible: true,
				},
				Value: &choice.Choice{
					Options: []*choice.Option{
						{NameId: "wifi", Label: "WiFi"},
						{NameId: "pool", Label: "Pool"},
					},
				},
			},
		},
	}

	answers := surveygo.Answers{
		"q-features": {"pool"},
	}

	cm := &CheckMark{Selected: "X", NotSelected: "-"}
	data, err := AnswersToCSV(s, answers, cm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	headers, rows := parseCSV(t, data)
	if len(headers) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(headers))
	}

	expected := []string{"-", "X"}
	for i, v := range expected {
		if rows[0][i] != v {
			t.Errorf("row[0][%d]: expected %q, got %q", i, v, rows[0][i])
		}
	}
}
