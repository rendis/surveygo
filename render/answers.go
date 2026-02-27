package render

import (
	"encoding/json"
	"strings"

	surveygo "github.com/rendis/surveygo/v2"
)

func parseAnswers(data []byte) (surveygo.Answers, error) {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	answers := make(surveygo.Answers, len(raw))
	for k, v := range raw {
		switch val := v.(type) {
		case []any:
			answers[k] = val
		case nil:
			// skip null values
		default:
			// wrap scalar or object in a single-element slice
			answers[k] = []any{val}
		}
	}
	return answers, nil
}

func extractTextValue(ans []any) string {
	if len(ans) == 0 {
		return ""
	}
	s, _ := ans[0].(string)
	return s
}

func extractPhoneValue(ans []any) string {
	if len(ans) == 0 {
		return ""
	}
	var result string
	for i, v := range ans {
		s, _ := v.(string)
		if i > 0 && s != "" {
			result += " "
		}
		result += s
	}
	return result
}

func extractSelectValue(ans []any) string {
	if len(ans) == 0 {
		return ""
	}
	s, _ := ans[0].(string)
	return s
}

func extractMultiSelectValues(ans []any) []string {
	if len(ans) == 0 {
		return nil
	}
	var result []string
	for _, v := range ans {
		switch val := v.(type) {
		case string:
			result = append(result, val)
		case []any:
			if len(val) > 0 {
				if s, ok := val[0].(string); ok {
					result = append(result, s)
				}
			}
		}
	}
	return result
}

func extractExternalValue(ans []any) (string, string) {
	var value, label string
	if len(ans) > 0 {
		value, _ = ans[0].(string)
	}
	if len(ans) > 1 {
		label, _ = ans[1].(string)
	}
	return value, label
}

func extractToggleValue(ans []any) bool {
	if len(ans) == 0 {
		return false
	}
	switch v := ans[0].(type) {
	case bool:
		return v
	case string:
		return strings.EqualFold(v, "true") || v == "1"
	default:
		return false
	}
}

func extractGroupInstances(ans []any) []surveygo.Answers {
	if len(ans) == 0 {
		return nil
	}
	var result []surveygo.Answers
	for _, v := range ans {
		if m, ok := v.(map[string]any); ok {
			inst := make(surveygo.Answers)
			for k, val := range m {
				if arr, ok := val.([]any); ok {
					inst[k] = arr
				}
			}
			result = append(result, inst)
		}
	}
	return result
}
