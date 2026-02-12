package render

import (
	surveygo "github.com/rendis/surveygo/v2"
)

// textTypes maps question types that render as simple text fields.
var textTypes = map[string]bool{
	"input_text":            true,
	"email":                 true,
	"telephone":             true,
	"date_time":             true,
	"identification_number": true,
	"slider":                true,
}

// selectTypes maps question types that render as a single select.
var selectTypes = map[string]bool{
	"single_select": true,
	"radio":         true,
}

func buildSurveyCard(survey *surveygo.Survey, tree *GroupTree, questions []GroupQuestions, answers surveygo.Answers) (*SurveyCard, error) {
	card := &SurveyCard{
		SurveyId: survey.NameId,
		Title:    survey.Title,
	}

	gqIndex := make(map[string]GroupQuestions, len(questions))
	for _, gq := range questions {
		gqIndex[gq.GroupNameId] = gq
	}

	for _, root := range tree.Roots {
		sections := buildSections(root, survey, gqIndex, answers)
		card.Sections = append(card.Sections, sections...)
	}

	return card, nil
}

func buildSections(node *GroupNode, survey *surveygo.Survey, gqIndex map[string]GroupQuestions, answers surveygo.Answers) []Section {
	grp := survey.Groups[node.NameId]
	if grp == nil {
		return nil
	}

	if !node.AllowRepeat {
		return buildGroupSection(node, survey, gqIndex, answers)
	}

	if node.RepeatDescendants > 0 || subtreeHasMultiSelect(node, gqIndex) {
		return []Section{buildRepeatListSection(node, survey, gqIndex, answers)}
	}

	return []Section{buildRepeatTableSection(node, survey, gqIndex, answers)}
}

func buildGroupSection(node *GroupNode, survey *surveygo.Survey, gqIndex map[string]GroupQuestions, answers surveygo.Answers) []Section {
	grp := survey.Groups[node.NameId]
	sec := Section{
		Type:   "group",
		NameId: grp.NameId,
		Title:  derefStr(grp.Title),
	}

	gq, ok := gqIndex[grp.NameId]
	if ok {
		for _, qi := range gq.Questions {
			field := buildField(qi, answers)
			sec.Fields = append(sec.Fields, field)
		}
	}

	for _, child := range node.Children {
		sec.Sections = append(sec.Sections, buildSections(child, survey, gqIndex, answers)...)
	}

	return []Section{sec}
}

func buildRepeatTableSection(node *GroupNode, survey *surveygo.Survey, gqIndex map[string]GroupQuestions, answers surveygo.Answers) Section {
	grp := survey.Groups[node.NameId]
	sec := Section{
		Type:   "repeat-table",
		NameId: grp.NameId,
		Title:  derefStr(grp.Title),
	}

	allQuestions := collectSubtreeQuestions(node, gqIndex)
	if len(allQuestions) == 0 {
		return sec
	}

	for _, qi := range allQuestions {
		col := Column{
			NameId:    qi.NameId,
			Label:     qi.Label,
			FieldType: fieldTypeFromQuestion(qi),
		}
		if multiSelectTypes[qi.QuestionType] {
			for _, opt := range qi.Options {
				col.Options = append(col.Options, OptionRef{
					NameId: opt.NameId,
					Label:  opt.Label,
				})
			}
		}
		sec.Columns = append(sec.Columns, col)
	}

	instances := extractGroupInstances(answers[grp.NameId])
	for _, inst := range instances {
		row := make(Row)
		for _, qi := range allQuestions {
			row[qi.NameId] = resolveValue(qi, inst)
		}
		sec.Rows = append(sec.Rows, row)
	}

	return sec
}

func buildRepeatListSection(node *GroupNode, survey *surveygo.Survey, gqIndex map[string]GroupQuestions, answers surveygo.Answers) Section {
	grp := survey.Groups[node.NameId]
	sec := Section{
		Type:   "repeat-list",
		NameId: grp.NameId,
		Title:  derefStr(grp.Title),
	}

	instances := extractGroupInstances(answers[grp.NameId])
	for _, inst := range instances {
		var instSections []Section

		gq, ok := gqIndex[grp.NameId]
		if ok && len(gq.Questions) > 0 {
			fieldSec := Section{
				Type:   "group",
				NameId: grp.NameId,
				Title:  derefStr(grp.Title),
			}
			for _, qi := range gq.Questions {
				fieldSec.Fields = append(fieldSec.Fields, buildField(qi, inst))
			}
			instSections = append(instSections, fieldSec)
		}

		for _, child := range node.Children {
			instSections = append(instSections, buildSections(child, survey, gqIndex, inst)...)
		}

		sec.Instances = append(sec.Instances, Instance{Sections: instSections})
	}

	return sec
}

func buildField(qi QuestionInfo, answers surveygo.Answers) Field {
	return Field{
		Type:   fieldTypeFromQuestion(qi),
		NameId: qi.NameId,
		Label:  qi.Label,
		Value:  resolveValue(qi, answers),
	}
}

func resolveValue(qi QuestionInfo, answers surveygo.Answers) any {
	ans := answers[qi.NameId]

	if qi.AnswerExpr != "" {
		if result, ok := evalAnswerExpr(qi.AnswerExpr, ans, qi.Options); ok {
			return result
		}
	}

	switch {
	case qi.QuestionType == "toggle":
		return extractToggleValue(ans)

	case qi.QuestionType == "telephone":
		return extractPhoneValue(ans)

	case selectTypes[qi.QuestionType]:
		selectedId := extractSelectValue(ans)
		if selectedId == "" {
			return nil
		}
		ref := OptionRef{NameId: selectedId}
		for _, opt := range qi.Options {
			if opt.NameId == selectedId {
				ref.Label = opt.Label
				break
			}
		}
		return ref

	case multiSelectTypes[qi.QuestionType]:
		selectedIds := extractMultiSelectValues(ans)
		selectedSet := make(map[string]bool, len(selectedIds))
		for _, id := range selectedIds {
			selectedSet[id] = true
		}
		var refs []OptionRef
		for _, opt := range qi.Options {
			refs = append(refs, OptionRef{
				NameId:   opt.NameId,
				Label:    opt.Label,
				Selected: selectedSet[opt.NameId],
			})
		}
		return refs

	case qi.QuestionType == "external_question":
		value, label := extractExternalValue(ans)
		if label != "" {
			if value != "" && value != label {
				return value + " - " + label
			}
			return label
		}
		return value

	default:
		return extractTextValue(ans)
	}
}

// subtreeHasMultiSelect returns true if the node or any descendant has
// multi_select or checkbox questions.
func subtreeHasMultiSelect(node *GroupNode, gqIndex map[string]GroupQuestions) bool {
	if gq, ok := gqIndex[node.NameId]; ok {
		for _, qi := range gq.Questions {
			if multiSelectTypes[qi.QuestionType] {
				return true
			}
		}
	}
	for _, child := range node.Children {
		if subtreeHasMultiSelect(child, gqIndex) {
			return true
		}
	}
	return false
}

// collectSubtreeQuestions collects all questions from the node and its
// descendants in DFS order, flattening sub-group questions into a single slice.
func collectSubtreeQuestions(node *GroupNode, gqIndex map[string]GroupQuestions) []QuestionInfo {
	var all []QuestionInfo
	if gq, ok := gqIndex[node.NameId]; ok {
		all = append(all, gq.Questions...)
	}
	for _, child := range node.Children {
		all = append(all, collectSubtreeQuestions(child, gqIndex)...)
	}
	return all
}

func fieldTypeFromQuestion(qi QuestionInfo) string {
	switch {
	case textTypes[qi.QuestionType]:
		return "text"
	case selectTypes[qi.QuestionType]:
		return "select"
	case multiSelectTypes[qi.QuestionType]:
		return "multi-select"
	case qi.QuestionType == "external_question":
		return "external"
	case qi.QuestionType == "toggle":
		return "toggle"
	default:
		return "text"
	}
}
