package render

import (
	"testing"

	surveygo "github.com/rendis/surveygo/v2"
)

func TestEvalAnswerExpr_SimpleIndex(t *testing.T) {
	result, ok := evalAnswerExpr(`ans[0]`, []any{"hello"}, nil)
	if !ok {
		t.Fatal("expected ok")
	}
	if result != "hello" {
		t.Fatalf("expected 'hello', got %v", result)
	}
}

func TestEvalAnswerExpr_PhoneConcat(t *testing.T) {
	result, ok := evalAnswerExpr(`ans[0] + " " + ans[1]`, []any{"+56", "912345678"}, nil)
	if !ok {
		t.Fatal("expected ok")
	}
	if result != "+56 912345678" {
		t.Fatalf("expected '+56 912345678', got %v", result)
	}
}

func TestEvalAnswerExpr_PhoneNumberOnly(t *testing.T) {
	result, ok := evalAnswerExpr(`ans[1]`, []any{"+56", "912345678"}, nil)
	if !ok {
		t.Fatal("expected ok")
	}
	if result != "912345678" {
		t.Fatalf("expected '912345678', got %v", result)
	}
}

func TestEvalAnswerExpr_ToggleSiNo(t *testing.T) {
	result, ok := evalAnswerExpr(`ans[0] ? "Si" : "No"`, []any{true}, nil)
	if !ok {
		t.Fatal("expected ok")
	}
	if result != "Si" {
		t.Fatalf("expected 'Si', got %v", result)
	}

	result, ok = evalAnswerExpr(`ans[0] ? "Si" : "No"`, []any{false}, nil)
	if !ok {
		t.Fatal("expected ok")
	}
	if result != "No" {
		t.Fatalf("expected 'No', got %v", result)
	}
}

func TestEvalAnswerExpr_WithOptions(t *testing.T) {
	options := []OptionInfo{
		{NameId: "run", Label: "RUN (Rol Unico Nacional)"},
		{NameId: "passport", Label: "Pasaporte"},
	}
	result, ok := evalAnswerExpr(`options[ans[0]]`, []any{"run"}, options)
	if !ok {
		t.Fatal("expected ok")
	}
	if result != "RUN (Rol Unico Nacional)" {
		t.Fatalf("expected label, got %v", result)
	}
}

func TestEvalAnswerExpr_SyntaxError(t *testing.T) {
	_, ok := evalAnswerExpr(`ans[`, []any{"hello"}, nil)
	if ok {
		t.Fatal("expected failure on syntax error")
	}
}

func TestEvalAnswerExpr_RuntimeError(t *testing.T) {
	_, ok := evalAnswerExpr(`ans[5]`, []any{"hello"}, nil)
	if ok {
		t.Fatal("expected failure on index out of bounds")
	}
}

func TestEvalAnswerExpr_NilAns(t *testing.T) {
	result, ok := evalAnswerExpr(`ans`, nil, nil)
	if !ok {
		t.Fatal("expected ok")
	}
	arr, isSlice := result.([]any)
	if !isSlice || len(arr) != 0 {
		t.Fatalf("expected empty slice, got %v (%T)", result, result)
	}
}

func TestEvalAnswerExprString_Coercion(t *testing.T) {
	result, ok := evalAnswerExprString(`ans[0]`, []any{true}, nil)
	if !ok {
		t.Fatal("expected ok")
	}
	if result != "true" {
		t.Fatalf("expected 'true', got %v", result)
	}
}

func TestEvalAnswerExprString_NilAns(t *testing.T) {
	result, ok := evalAnswerExprString(`ans`, nil, nil)
	if !ok {
		t.Fatal("expected ok")
	}
	if result != "[]" {
		t.Fatalf("expected '[]', got %v", result)
	}
}

func TestResolveValue_WithAnswerExpr(t *testing.T) {
	qi := QuestionInfo{
		NameId:       "q1",
		QuestionType: "telephone",
		AnswerExpr:   `ans[1]`,
	}
	answers := surveygo.Answers{"q1": {"+56", "912345678"}}
	result := resolveValue(qi, answers)
	if result != "912345678" {
		t.Fatalf("expected '912345678', got %v", result)
	}
}

func TestResolveValue_AnswerExprFallback(t *testing.T) {
	qi := QuestionInfo{
		NameId:       "q1",
		QuestionType: "telephone",
		AnswerExpr:   `invalid!!!`,
	}
	answers := surveygo.Answers{"q1": {"+56", "912345678"}}
	result := resolveValue(qi, answers)
	// Should fall back to ExtractPhoneValue
	if result != "+56 912345678" {
		t.Fatalf("expected fallback '+56 912345678', got %v", result)
	}
}

func TestResolveValue_WithoutAnswerExpr(t *testing.T) {
	qi := QuestionInfo{
		NameId:       "q1",
		QuestionType: "telephone",
	}
	answers := surveygo.Answers{"q1": {"+56", "912345678"}}
	result := resolveValue(qi, answers)
	if result != "+56 912345678" {
		t.Fatalf("expected '+56 912345678', got %v", result)
	}
}

func TestResolveValue_SelectWithOptionsExpr(t *testing.T) {
	qi := QuestionInfo{
		NameId:       "q1",
		QuestionType: "single_select",
		AnswerExpr:   `options[ans[0]]`,
		Options: []OptionInfo{
			{NameId: "run", Label: "RUN"},
			{NameId: "passport", Label: "Pasaporte"},
		},
	}
	answers := surveygo.Answers{"q1": {"run"}}
	result := resolveValue(qi, answers)
	if result != "RUN" {
		t.Fatalf("expected 'RUN', got %v", result)
	}
}
