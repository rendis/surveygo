package render

import (
	"fmt"

	surveygo "github.com/rendis/surveygo/v2"
	"github.com/rendis/surveygo/v2/question"
	"github.com/rendis/surveygo/v2/question/types"
	"github.com/rendis/surveygo/v2/question/types/choice"
	"github.com/rendis/surveygo/v2/question/types/external"
	"github.com/rendis/surveygo/v2/question/types/text"
)

func extractGroupQuestions(survey *surveygo.Survey) ([]GroupQuestions, error) {
	var result []GroupQuestions

	visited := make(map[string]bool)
	if err := extractOrdered(survey.GroupsOrder, survey, &result, visited); err != nil {
		return nil, err
	}

	return result, nil
}

func extractOrdered(groupIDs []string, survey *surveygo.Survey, result *[]GroupQuestions, visited map[string]bool) error {
	for _, gID := range groupIDs {
		if visited[gID] {
			continue
		}
		visited[gID] = true

		grp, ok := survey.Groups[gID]
		if !ok {
			return fmt.Errorf("group %q not found", gID)
		}

		// Recurse into sub-groups first
		if len(grp.GroupsOrder) > 0 {
			if err := extractOrdered(grp.GroupsOrder, survey, result, visited); err != nil {
				return err
			}
		}

		// Process direct questions
		if len(grp.QuestionsIds) == 0 {
			continue
		}

		gq := GroupQuestions{GroupNameId: gID}
		for _, qID := range grp.QuestionsIds {
			q, ok := survey.Questions[qID]
			if !ok {
				return fmt.Errorf("question %q referenced by group %q not found", qID, gID)
			}

			info := mapQuestion(q)
			gq.Questions = append(gq.Questions, info)

			// Follow GroupsIds from choice options
			if types.IsSimpleChoiceType(q.QTyp) {
				c, err := choice.CastToChoice(q.Value)
				if err == nil {
					for _, opt := range c.Options {
						if len(opt.GroupsIds) > 0 {
							if err := extractOrdered(opt.GroupsIds, survey, result, visited); err != nil {
								return err
							}
						}
					}
				}
			}
		}

		*result = append(*result, gq)
	}
	return nil
}

func mapQuestion(q *question.Question) QuestionInfo {
	label := q.Label
	if label == "" {
		label = getPlaceholder(q)
	}

	info := QuestionInfo{
		NameId:       q.NameId,
		Label:        label,
		QuestionType: string(q.QTyp),
		AnswerExpr:   q.AnswerExpr,
	}

	// Options (simple choice types: single_select, multi_select, radio, checkbox)
	if types.IsSimpleChoiceType(q.QTyp) {
		c, err := choice.CastToChoice(q.Value)
		if err == nil {
			for _, opt := range c.Options {
				info.Options = append(info.Options, OptionInfo{
					NameId: opt.NameId,
					Label:  opt.Label,
					Value:  opt.Value,
				})
			}
		}
	}

	// Format (date_time)
	if q.QTyp == types.QTypeDateTime {
		dt, err := text.CastToDateTime(q.Value)
		if err == nil {
			info.Format = dt.Format
		}
	}

	// ExternalType (external_question)
	if q.QTyp == types.QTypeExternalQuestion {
		if ext, ok := q.Value.(*external.ExternalQuestion); ok {
			info.ExternalType = ext.ExternalType
		}
	}

	return info
}

// getPlaceholder extracts the Placeholder field from a question's value type.
// All value types embed types.QBase which has Placeholder *string.
func getPlaceholder(q *question.Question) string {
	if q.Value == nil {
		return ""
	}
	switch v := q.Value.(type) {
	case *choice.Choice:
		return derefStr(v.Placeholder)
	case *choice.Slider:
		return derefStr(v.Placeholder)
	case *choice.Toggle:
		return derefStr(v.Placeholder)
	case *text.FreeText:
		return derefStr(v.Placeholder)
	case *text.Email:
		return derefStr(v.Placeholder)
	case *text.Telephone:
		return derefStr(v.Placeholder)
	case *text.DateTime:
		return derefStr(v.Placeholder)
	case *text.InformationText:
		return derefStr(v.Placeholder)
	case *text.IdentificationNumber:
		return derefStr(v.Placeholder)
	case *external.ExternalQuestion:
		return derefStr(v.Placeholder)
	default:
		return ""
	}
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
