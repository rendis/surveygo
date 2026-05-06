package reviewer

import (
	"strings"
	"testing"

	"github.com/rendis/surveygo/v2/question/types"
)

func TestReviewTextReturnsCastErrorsInsteadOfPanicking(t *testing.T) {
	tests := []struct {
		name    string
		qtype   types.QuestionType
		answers []any
	}{
		{
			name:    "input text",
			qtype:   types.QTypeInputText,
			answers: []any{"hello"},
		},
		{
			name:    "text area",
			qtype:   types.QTypeTextArea,
			answers: []any{"hello"},
		},
		{
			name:    "email",
			qtype:   types.QTypeEmail,
			answers: []any{"user@example.com"},
		},
		{
			name:    "telephone",
			qtype:   types.QTypeTelephone,
			answers: []any{"+56912345678"},
		},
		{
			name:    "date time",
			qtype:   types.QTypeDateTime,
			answers: []any{"2026-05-06"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("ReviewText panicked: %v", r)
				}
			}()

			err := ReviewText(map[string]any{}, tt.answers, tt.qtype)
			if err == nil {
				t.Fatal("expected cast error, got nil")
			}
			if !strings.Contains(err.Error(), "invalid type") {
				t.Fatalf("expected invalid type error, got %q", err.Error())
			}
		})
	}
}
