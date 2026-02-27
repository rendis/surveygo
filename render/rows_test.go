package render

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
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

// colIndex returns the index of the first header containing substr, or -1.
func colIndex(headers []string, substr string) int {
	for i, h := range headers {
		if strings.Contains(h, substr) {
			return i
		}
	}
	return -1
}

func TestAnswersToRows_MultiSelectHeaders(t *testing.T) {
	survey := loadSurvey(t, "sample_nested.json")
	matrix, err := AnswersToRows(survey, surveygo.Answers{})
	if err != nil {
		t.Fatalf("AnswersToRows: %v", err)
	}

	if len(matrix) == 0 {
		t.Fatal("expected at least 1 row (headers)")
	}

	headers := matrix[0]

	// emp_contract_type is multi_select with 3 options; headers use "Tipo de Contrato - <option>"
	var count int
	for _, h := range headers {
		if strings.Contains(h, "Tipo de Contrato") {
			count++
		}
	}

	if count != 3 {
		t.Fatalf("expected 3 multi-select columns for emp_contract_type, got %d", count)
	}
}

func TestAnswersToRows_RepeatExpansion(t *testing.T) {
	survey := loadSurvey(t, "sample.json")
	answers := loadAnswersFile(t, "sample_answers.json")

	matrix, err := AnswersToRows(survey, answers)
	if err != nil {
		t.Fatalf("AnswersToRows: %v", err)
	}

	// sample_answers.json has 2 adults → 1 header + 2 data rows = 3
	if len(matrix) != 3 {
		t.Fatalf("expected 3 rows (1 header + 2 data), got %d", len(matrix))
	}

	headers := matrix[0]
	nameIdx := colIndex(headers, "Primer Nombre")
	if nameIdx < 0 {
		t.Fatal("'Primer Nombre' column not found in headers")
	}

	if matrix[1][nameIdx] != "Juan" {
		t.Errorf("row[1] first name = %q, want %q", matrix[1][nameIdx], "Juan")
	}
	if matrix[2][nameIdx] != "Maria" {
		t.Errorf("row[2] first name = %q, want %q", matrix[2][nameIdx], "Maria")
	}
}

func TestAnswersToRows_NestedCartesian(t *testing.T) {
	survey := loadSurvey(t, "sample_nested.json")
	answers := loadAnswersFile(t, "sample_nested_answers.json")

	matrix, err := AnswersToRows(survey, answers)
	if err != nil {
		t.Fatalf("AnswersToRows: %v", err)
	}

	// 2 departments: Ingeniería (3 employees) + RRHH (2 employees) = 5 data rows + 1 header = 6
	if len(matrix) != 6 {
		t.Fatalf("expected 6 rows (1 header + 5 data), got %d", len(matrix))
	}

	headers := matrix[0]
	// "Razón Social" is the label for company_name
	companyIdx := colIndex(headers, "Social")
	if companyIdx < 0 {
		t.Fatal("company_name column not found in headers")
	}

	// All 5 data rows should have company_name duplicated
	for i := 1; i < len(matrix); i++ {
		if matrix[i][companyIdx] != "Tether SpA" {
			t.Errorf("row[%d] company_name = %q, want %q", i, matrix[i][companyIdx], "Tether SpA")
		}
	}
}

func TestAnswersToRows_ToggleStringTrue(t *testing.T) {
	survey := loadSurvey(t, "sample_nested.json")

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
	matrix, err := AnswersToRows(survey, answers, cm)
	if err != nil {
		t.Fatalf("AnswersToRows: %v", err)
	}

	if len(matrix) < 2 {
		t.Fatal("expected at least 2 rows (header + data)")
	}

	headers := matrix[0]
	// "¿Departamento activo?" is the label for dept_active
	toggleIdx := colIndex(headers, "activo")
	if toggleIdx < 0 {
		t.Fatal("dept_active toggle column not found")
	}

	// String "true" should be handled by extractToggleValue and render as "Sí"
	if matrix[1][toggleIdx] != "Sí" {
		t.Errorf("row[1] toggle = %q, want %q", matrix[1][toggleIdx], "Sí")
	}
}

func TestAnswersToRows_AnswerExpr(t *testing.T) {
	survey := loadSurvey(t, "sample_with_expr.json")
	answers := loadAnswersFile(t, "sample_with_expr_answers.json")

	matrix, err := AnswersToRows(survey, answers)
	if err != nil {
		t.Fatalf("AnswersToRows: %v", err)
	}

	if len(matrix) < 3 {
		t.Fatalf("expected at least 3 rows (1 header + 2 data), got %d", len(matrix))
	}

	headers := matrix[0]
	// "Contacto de Emergencia" is the label for emergency_contact
	ecIdx := colIndex(headers, "Emergencia")
	if ecIdx < 0 {
		t.Fatal("emergency_contact column not found")
	}

	// emergency_contact has answerExpr: options[ans[0]] → maps "true"→"Si", "false"→"No"
	if matrix[1][ecIdx] != "Si" {
		t.Errorf("row[1] emergency_contact = %q, want %q", matrix[1][ecIdx], "Si")
	}
	if matrix[2][ecIdx] != "No" {
		t.Errorf("row[2] emergency_contact = %q, want %q", matrix[2][ecIdx], "No")
	}
}

func TestAnswersToRows_EmptyAnswers(t *testing.T) {
	survey := loadSurvey(t, "sample.json")
	answers := loadAnswersFile(t, "sample_answers_empty.json")

	matrix, err := AnswersToRows(survey, answers)
	if err != nil {
		t.Fatalf("AnswersToRows: %v", err)
	}

	// Empty answers should produce 1 header + 1 empty data row = 2
	if len(matrix) != 2 {
		t.Fatalf("expected 2 rows (1 header + 1 data), got %d", len(matrix))
	}

	// All values in data row should be empty strings
	for i, val := range matrix[1] {
		if val != "" {
			t.Errorf("column %q: expected empty value, got %q", matrix[0][i], val)
		}
	}
}

func TestAnswersToRows_CheckMark(t *testing.T) {
	survey := loadSurvey(t, "sample.json")
	answers := loadAnswersFile(t, "sample_answers.json")

	cm := &CheckMark{Selected: "Sí", NotSelected: "No"}
	matrix, err := AnswersToRows(survey, answers, cm)
	if err != nil {
		t.Fatalf("AnswersToRows: %v", err)
	}

	headers := matrix[0]

	// adult_permissions is multi_select; find "Retirar" column (pickup option label)
	pickupIdx := colIndex(headers, "Retirar")
	if pickupIdx < 0 {
		t.Fatal("pickup option column not found")
	}

	// Both adults have pickup selected
	if matrix[1][pickupIdx] != "Sí" {
		t.Errorf("row[1] pickup = %q, want %q", matrix[1][pickupIdx], "Sí")
	}
	if matrix[2][pickupIdx] != "Sí" {
		t.Errorf("row[2] pickup = %q, want %q", matrix[2][pickupIdx], "Sí")
	}
}
