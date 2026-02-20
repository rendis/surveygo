package render

import (
	"bytes"
	"encoding/csv"
	"fmt"

	surveygo "github.com/rendis/surveygo/v2"
)

// multiSelectTypes identifies question types that expand to one boolean column per option.
var multiSelectTypes = map[string]bool{
	"multi_select": true,
	"checkbox":     true,
}

// csvColumn represents a single CSV column header with metadata for value extraction.
type csvColumn struct {
	header     string // column header name
	questionID string // question nameId for answer lookup
	qType      string // question type
	optionID   string // non-empty for multi_select/checkbox boolean columns
}

func generateCSV(survey *surveygo.Survey, tree *GroupTree, questions []GroupQuestions, answers surveygo.Answers, cm *CheckMark) ([]byte, error) {
	gqIndex := make(map[string]GroupQuestions, len(questions))
	for _, gq := range questions {
		gqIndex[gq.GroupNameId] = gq
	}

	// 1. Build column headers via DFS of group tree.
	var cols []csvColumn
	for _, root := range tree.Roots {
		buildColumns(root, survey, gqIndex, &cols)
	}

	// 2. Build rows via cartesian product DFS.
	rows := []map[string]string{make(map[string]string)}
	for _, root := range tree.Roots {
		rows = fillRows(root, answers, survey, gqIndex, cols, rows, cm)
	}

	// 3. Write CSV.
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	headers := make([]string, len(cols))
	for i, c := range cols {
		headers[i] = c.header
	}
	if err := w.Write(headers); err != nil {
		return nil, fmt.Errorf("writing CSV headers: %w", err)
	}

	for _, row := range rows {
		record := make([]string, len(cols))
		for i, c := range cols {
			record[i] = row[c.header]
		}
		if err := w.Write(record); err != nil {
			return nil, fmt.Errorf("writing CSV row: %w", err)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, fmt.Errorf("flushing CSV: %w", err)
	}

	return buf.Bytes(), nil
}

func buildColumns(node *GroupNode, survey *surveygo.Survey, gqIndex map[string]GroupQuestions, cols *[]csvColumn) {
	if gq, ok := gqIndex[node.NameId]; ok {
		for _, q := range gq.Questions {
			if q.AnswerExpr != "" {
				*cols = append(*cols, csvColumn{
					header:     questionHeader(q),
					questionID: q.NameId,
					qType:      q.QuestionType,
				})
				continue
			}
			if multiSelectTypes[q.QuestionType] {
				for _, opt := range q.Options {
					*cols = append(*cols, csvColumn{
						header:     optionHeader(q, opt),
						questionID: q.NameId,
						qType:      q.QuestionType,
						optionID:   opt.NameId,
					})
				}
			} else {
				*cols = append(*cols, csvColumn{
					header:     questionHeader(q),
					questionID: q.NameId,
					qType:      q.QuestionType,
				})
			}
		}
	}

	for _, child := range node.Children {
		buildColumns(child, survey, gqIndex, cols)
	}
}

func fillRows(node *GroupNode, answers surveygo.Answers, survey *surveygo.Survey, gqIndex map[string]GroupQuestions, cols []csvColumn, rows []map[string]string, cm *CheckMark) []map[string]string {
	if node.AllowRepeat {
		rows = expandRepeatGroup(node, answers, survey, gqIndex, cols, rows, cm)
	} else {
		fillGroupValues(node, answers, gqIndex, rows, cm)
		for _, child := range node.Children {
			rows = fillRows(child, answers, survey, gqIndex, cols, rows, cm)
		}
	}
	return rows
}

func expandRepeatGroup(node *GroupNode, answers surveygo.Answers, survey *surveygo.Survey, gqIndex map[string]GroupQuestions, cols []csvColumn, rows []map[string]string, cm *CheckMark) []map[string]string {
	instances := extractGroupInstances(answers[node.NameId])
	if len(instances) == 0 {
		return rows
	}

	var expanded []map[string]string
	for _, inst := range instances {
		for _, row := range rows {
			cloned := cloneRow(row)
			fillGroupValues(node, inst, gqIndex, []map[string]string{cloned}, cm)
			expanded = append(expanded, cloned)
		}
	}

	result := make([]map[string]string, 0, len(expanded))
	for i, inst := range instances {
		batchSize := len(rows)
		start := i * batchSize
		end := start + batchSize
		batch := expanded[start:end]

		for _, child := range node.Children {
			batch = fillRows(child, inst, survey, gqIndex, cols, batch, cm)
		}
		result = append(result, batch...)
	}

	return result
}

func fillGroupValues(node *GroupNode, answers surveygo.Answers, gqIndex map[string]GroupQuestions, rows []map[string]string, cm *CheckMark) {
	gq, ok := gqIndex[node.NameId]
	if !ok {
		return
	}

	selMark, notSelMark := "true", "false"
	if cm != nil {
		selMark, notSelMark = cm.Selected, cm.NotSelected
	}

	for _, q := range gq.Questions {
		ans := answers[q.NameId]

		if q.AnswerExpr != "" {
			if val, ok := evalAnswerExprString(q.AnswerExpr, ans, q.Options); ok {
				for _, row := range rows {
					row[questionHeader(q)] = val
				}
				continue
			}
		}

		if multiSelectTypes[q.QuestionType] {
			selected := make(map[string]bool)
			for _, v := range extractMultiSelectValues(ans) {
				selected[v] = true
			}
			for _, opt := range q.Options {
				colName := optionHeader(q, opt)
				val := notSelMark
				if selected[opt.NameId] {
					val = selMark
				}
				for _, row := range rows {
					row[colName] = val
				}
			}
		} else {
			val := extractCSVValue(q.QuestionType, ans, selMark, notSelMark)
			for _, row := range rows {
				row[questionHeader(q)] = val
			}
		}
	}
}

func extractCSVValue(qType string, ans []any, selMark, notSelMark string) string {
	switch qType {
	case "input_text", "email", "date_time", "identification_number":
		return extractTextValue(ans)
	case "telephone":
		return extractPhoneValue(ans)
	case "single_select", "radio":
		return extractSelectValue(ans)
	case "external_question":
		val, label := extractExternalValue(ans)
		if label != "" {
			return label
		}
		return val
	case "toggle":
		if extractToggleValue(ans) {
			return selMark
		}
		return notSelMark
	default:
		return extractTextValue(ans)
	}
}

func questionHeader(q QuestionInfo) string {
	if q.Label != "" {
		return q.Label
	}
	return q.NameId
}

func optionHeader(q QuestionInfo, opt OptionInfo) string {
	lbl := opt.Label
	if lbl == "" {
		lbl = opt.NameId
	}
	return questionHeader(q) + " - " + lbl
}

func cloneRow(row map[string]string) map[string]string {
	c := make(map[string]string, len(row))
	for k, v := range row {
		c[k] = v
	}
	return c
}
