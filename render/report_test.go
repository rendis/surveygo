package render

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	surveygo "github.com/rendis/surveygo/v2"
)

func loadSurvey(t *testing.T, name string) *surveygo.Survey {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(testdataDir, name))
	if err != nil {
		t.Fatalf("reading survey %s: %v", name, err)
	}
	s, err := surveygo.ParseFromBytes(data)
	if err != nil {
		t.Fatalf("parsing survey %s: %v", name, err)
	}
	return s
}

func loadAnswersFile(t *testing.T, name string) surveygo.Answers {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(testdataDir, name))
	if err != nil {
		t.Fatalf("reading answers %s: %v", name, err)
	}
	var answers surveygo.Answers
	if err := json.Unmarshal(data, &answers); err != nil {
		t.Fatalf("unmarshaling answers %s: %v", name, err)
	}
	return answers
}

func TestReportColumns_MultiSelectExplosion(t *testing.T) {
	survey := loadSurvey(t, "sample_nested.json")
	cols, _, err := ReportColumns(survey)
	if err != nil {
		t.Fatalf("ReportColumns: %v", err)
	}

	// emp_contract_type is a multi_select with 3 options: fulltime, parttime, remote
	var multiSelectCols []ReportColumn
	for _, c := range cols {
		if c.QuestionID == "emp_contract_type" {
			multiSelectCols = append(multiSelectCols, c)
		}
	}

	if len(multiSelectCols) != 3 {
		t.Fatalf("expected 3 multi-select columns for emp_contract_type, got %d", len(multiSelectCols))
	}

	expectedOptions := []string{"fulltime", "parttime", "remote"}
	for i, c := range multiSelectCols {
		if c.OptionID != expectedOptions[i] {
			t.Errorf("column[%d] optionID = %q, want %q", i, c.OptionID, expectedOptions[i])
		}
		if c.QType != "multi_select" {
			t.Errorf("column[%d] qType = %q, want %q", i, c.QType, "multi_select")
		}
	}
}

func TestReportColumns_GroupIDAssignment(t *testing.T) {
	survey := loadSurvey(t, "sample_nested.json")
	cols, _, err := ReportColumns(survey)
	if err != nil {
		t.Fatalf("ReportColumns: %v", err)
	}

	// Verify specific columns belong to expected groups
	expectations := map[string]string{
		"company_name": "grp-company_info",
		"company_rut":  "grp-company_info",
		"addr_region":  "grp-addr_location",
		"addr_street":  "grp-addr_detail",
		"dept_name":    "grp-dept_info",
		"emp_name":     "grp-emp_basic",
		"emp_email":    "grp-emp_basic",
	}

	colMap := make(map[string]ReportColumn)
	for _, c := range cols {
		if c.OptionID == "" { // skip multi-select expanded columns
			colMap[c.QuestionID] = c
		}
	}

	for qID, wantGroup := range expectations {
		c, ok := colMap[qID]
		if !ok {
			t.Errorf("column %q not found", qID)
			continue
		}
		if c.GroupID != wantGroup {
			t.Errorf("column %q: groupID = %q, want %q", qID, c.GroupID, wantGroup)
		}
	}
}

func TestReportRows_RepeatExpansion(t *testing.T) {
	survey := loadSurvey(t, "sample.json")
	answers := loadAnswersFile(t, "sample_answers.json")

	cols, tree, err := ReportColumns(survey)
	if err != nil {
		t.Fatalf("ReportColumns: %v", err)
	}

	rows, err := ReportRows(survey, tree, cols, answers, nil)
	if err != nil {
		t.Fatalf("ReportRows: %v", err)
	}

	// sample_answers.json has 2 adults → 2 rows
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}

	// Check first row has Juan's data
	found := false
	for _, c := range cols {
		if c.QuestionID == "adult_first_name" {
			if rows[0][c.Header] == "Juan" {
				found = true
			}
			break
		}
	}
	if !found {
		t.Error("first row should contain adult_first_name=Juan")
	}

	// Check second row has Maria's data
	found = false
	for _, c := range cols {
		if c.QuestionID == "adult_first_name" {
			if rows[1][c.Header] == "Maria" {
				found = true
			}
			break
		}
	}
	if !found {
		t.Error("second row should contain adult_first_name=Maria")
	}
}

func TestReportRows_NestedCartesian(t *testing.T) {
	survey := loadSurvey(t, "sample_nested.json")
	answers := loadAnswersFile(t, "sample_nested_answers.json")

	cols, tree, err := ReportColumns(survey)
	if err != nil {
		t.Fatalf("ReportColumns: %v", err)
	}

	rows, err := ReportRows(survey, tree, cols, answers, nil)
	if err != nil {
		t.Fatalf("ReportRows: %v", err)
	}

	// 2 departments: Ingeniería (3 employees) + RRHH (2 employees) = 5 rows
	if len(rows) != 5 {
		t.Fatalf("expected 5 rows (3+2 from nested repeat), got %d", len(rows))
	}

	// All 5 rows should have company_name duplicated
	var companyHeader string
	for _, c := range cols {
		if c.QuestionID == "company_name" {
			companyHeader = c.Header
			break
		}
	}

	for i, row := range rows {
		if row[companyHeader] != "Tether SpA" {
			t.Errorf("row[%d] company_name = %q, want %q", i, row[companyHeader], "Tether SpA")
		}
	}
}

func TestReportRows_ToggleWithStringAnswer(t *testing.T) {
	// Build a minimal survey with a toggle question to test string "true"/"false" handling
	survey := loadSurvey(t, "sample_nested.json")

	cols, tree, err := ReportColumns(survey)
	if err != nil {
		t.Fatalf("ReportColumns: %v", err)
	}

	// Build answers with dept_active toggle as string "true" (simulating gRPC string booleans)
	answers := surveygo.Answers{
		"company_name": {"Test Co"},
		"grp-departments": {
			map[string]any{
				"dept_name":   []any{"Engineering"},
				"dept_active": []any{"true"}, // string, not bool
				"grp-dept_employees": []any{
					map[string]any{
						"emp_name": []any{"Alice"},
					},
				},
			},
		},
	}

	cm := &CheckMark{Selected: "Sí", NotSelected: "No"}
	rows, err := ReportRows(survey, tree, cols, answers, cm)
	if err != nil {
		t.Fatalf("ReportRows: %v", err)
	}

	if len(rows) == 0 {
		t.Fatal("expected at least 1 row")
	}

	// Find the dept_active toggle column
	var toggleHeader string
	for _, c := range cols {
		if c.QuestionID == "dept_active" {
			toggleHeader = c.Header
			break
		}
	}
	if toggleHeader == "" {
		t.Fatal("dept_active toggle column not found")
	}

	// String "true" should be handled by extractToggleValue and render as "Sí"
	if rows[0][toggleHeader] != "Sí" {
		t.Errorf("row[0] toggle = %q, want %q", rows[0][toggleHeader], "Sí")
	}
}

func TestReportRows_AnswerExprTransform(t *testing.T) {
	// sample_with_expr has answerExpr-based questions
	survey := loadSurvey(t, "sample_with_expr.json")
	answers := loadAnswersFile(t, "sample_with_expr_answers.json")

	cols, tree, err := ReportColumns(survey)
	if err != nil {
		t.Fatalf("ReportColumns: %v", err)
	}

	rows, err := ReportRows(survey, tree, cols, answers, nil)
	if err != nil {
		t.Fatalf("ReportRows: %v", err)
	}

	if len(rows) < 2 {
		t.Fatalf("expected at least 2 rows, got %d", len(rows))
	}

	// emergency_contact has answerExpr: options[ans[0]] → maps "true"→"Si", "false"→"No"
	var ecHeader string
	for _, c := range cols {
		if c.QuestionID == "emergency_contact" {
			ecHeader = c.Header
			break
		}
	}
	if ecHeader == "" {
		t.Fatal("emergency_contact column not found")
	}

	if rows[0][ecHeader] != "Si" {
		t.Errorf("row[0] emergency_contact = %q, want %q", rows[0][ecHeader], "Si")
	}
	if rows[1][ecHeader] != "No" {
		t.Errorf("row[1] emergency_contact = %q, want %q", rows[1][ecHeader], "No")
	}
}

func TestReportRows_EmptyAnswers(t *testing.T) {
	survey := loadSurvey(t, "sample.json")
	answers := loadAnswersFile(t, "sample_answers_empty.json")

	cols, tree, err := ReportColumns(survey)
	if err != nil {
		t.Fatalf("ReportColumns: %v", err)
	}

	rows, err := ReportRows(survey, tree, cols, answers, nil)
	if err != nil {
		t.Fatalf("ReportRows: %v", err)
	}

	// Empty answers should produce 1 row with empty values (not 0 rows)
	if len(rows) != 1 {
		t.Fatalf("expected 1 row for empty answers, got %d", len(rows))
	}

	// All values should be empty strings
	for _, c := range cols {
		if rows[0][c.Header] != "" {
			t.Errorf("column %q: expected empty value, got %q", c.Header, rows[0][c.Header])
		}
	}
}

func TestReportRows_CheckMark(t *testing.T) {
	survey := loadSurvey(t, "sample.json")
	answers := loadAnswersFile(t, "sample_answers.json")

	cols, tree, err := ReportColumns(survey)
	if err != nil {
		t.Fatalf("ReportColumns: %v", err)
	}

	cm := &CheckMark{Selected: "Sí", NotSelected: "No"}
	rows, err := ReportRows(survey, tree, cols, answers, cm)
	if err != nil {
		t.Fatalf("ReportRows: %v", err)
	}

	// adult_permissions is multi_select with 2 options: pickup, access_to_digital_content
	// First adult (Juan) has pickup selected, second (Maria) has both selected
	var pickupHeader string
	for _, c := range cols {
		if c.QuestionID == "adult_permissions" && c.OptionID == "pickup" {
			pickupHeader = c.Header
			break
		}
	}
	if pickupHeader == "" {
		t.Fatal("pickup option column not found")
	}

	if rows[0][pickupHeader] != "Sí" {
		t.Errorf("row[0] pickup = %q, want %q", rows[0][pickupHeader], "Sí")
	}
	if rows[1][pickupHeader] != "Sí" {
		t.Errorf("row[1] pickup = %q, want %q", rows[1][pickupHeader], "Sí")
	}
}

func TestTopLevelGroup(t *testing.T) {
	survey := loadSurvey(t, "sample_nested.json")
	_, tree, err := ReportColumns(survey)
	if err != nil {
		t.Fatalf("ReportColumns: %v", err)
	}

	tests := []struct {
		groupID string
		want    string
	}{
		{"grp-company", "grp-company"},           // root → itself
		{"grp-company_info", "grp-company"},       // direct child → root
		{"grp-addr_location", "grp-company"},      // nested 2 levels → root
		{"grp-departments", "grp-company"},         // repeat child → root
		{"grp-dept_info", "grp-company"},           // inside repeat → root
		{"grp-emp_basic", "grp-company"},           // nested repeat → root
		{"nonexistent", "nonexistent"},             // unknown → itself
	}

	for _, tt := range tests {
		got := tree.TopLevelGroup(tt.groupID)
		if got != tt.want {
			t.Errorf("TopLevelGroup(%q) = %q, want %q", tt.groupID, got, tt.want)
		}
	}
}
