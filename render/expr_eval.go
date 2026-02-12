package render

import (
	"fmt"

	"github.com/expr-lang/expr"
)

// evalAnswerExpr evaluates an expr-lang expression against a raw answer slice.
// Returns (result, true) on success, or (nil, false) on error (silent fallback).
// Environment:
//   - ans: the raw []any answer data
//   - options: map[string]string (nameId -> label), only if options are provided
func evalAnswerExpr(expression string, ans []any, options []OptionInfo) (any, bool) {
	env := buildExprEnv(ans, options)

	result, err := expr.Eval(expression, env)
	if err != nil {
		return nil, false
	}
	return result, true
}

// evalAnswerExprString is a convenience wrapper that coerces the result to string.
func evalAnswerExprString(expression string, ans []any, options []OptionInfo) (string, bool) {
	result, ok := evalAnswerExpr(expression, ans, options)
	if !ok {
		return "", false
	}
	if result == nil {
		return "", true
	}
	return fmt.Sprintf("%v", result), true
}

func buildExprEnv(ans []any, options []OptionInfo) map[string]any {
	if ans == nil {
		ans = []any{}
	}
	env := map[string]any{
		"ans": ans,
	}
	if len(options) > 0 {
		optMap := make(map[string]string, len(options))
		for _, opt := range options {
			optMap[opt.NameId] = opt.Label
		}
		env["options"] = optMap
	}
	return env
}
