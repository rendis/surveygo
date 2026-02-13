package render

import (
	"fmt"

	surveygo "github.com/rendis/surveygo/v2"
)

// ReportColumn is a single column in a multi-record report.
type ReportColumn struct {
	Header     string // display header (label, or "label - option" for multi-select)
	QuestionID string // question nameId for answer lookup
	QType      string // surveygo question type
	OptionID   string // non-empty for multi-select/checkbox boolean columns
	GroupID    string // immediate group nameId this column belongs to
}

// ReportColumns resolves ordered column definitions from a survey definition.
// It builds the group tree and extracts questions, returning exported column types
// along with the tree (needed for ReportRows and TopLevelGroup).
func ReportColumns(survey *surveygo.Survey) ([]ReportColumn, *GroupTree, error) {
	tree, err := buildGroupTree(survey)
	if err != nil {
		return nil, nil, fmt.Errorf("building group tree: %w", err)
	}

	questions, err := extractGroupQuestions(survey)
	if err != nil {
		return nil, nil, fmt.Errorf("extracting questions: %w", err)
	}

	gqIndex := make(map[string]GroupQuestions, len(questions))
	for _, gq := range questions {
		gqIndex[gq.GroupNameId] = gq
	}

	var cols []csvColumn
	for _, root := range tree.Roots {
		buildColumns(root, survey, gqIndex, &cols)
	}

	report := make([]ReportColumn, len(cols))
	for i, c := range cols {
		report[i] = ReportColumn{
			Header:     c.header,
			QuestionID: c.questionID,
			QType:      c.qType,
			OptionID:   c.optionID,
			GroupID:    c.groupID,
		}
	}

	return report, tree, nil
}

// ReportRows extracts row data for one set of answers against resolved columns.
// Returns 1+ rows (multiple when repeatable groups expand via cartesian product).
// Each row is map[header]string matching ReportColumn.Header keys.
func ReportRows(survey *surveygo.Survey, tree *GroupTree, columns []ReportColumn, answers surveygo.Answers, cm *CheckMark) ([]map[string]string, error) {
	questions, err := extractGroupQuestions(survey)
	if err != nil {
		return nil, fmt.Errorf("extracting questions: %w", err)
	}

	gqIndex := make(map[string]GroupQuestions, len(questions))
	for _, gq := range questions {
		gqIndex[gq.GroupNameId] = gq
	}

	cols := make([]csvColumn, len(columns))
	for i, c := range columns {
		cols[i] = csvColumn{
			header:     c.Header,
			questionID: c.QuestionID,
			qType:      c.QType,
			optionID:   c.OptionID,
			groupID:    c.GroupID,
		}
	}

	rows := []map[string]string{make(map[string]string)}
	for _, root := range tree.Roots {
		rows = fillRows(root, answers, survey, gqIndex, cols, rows, cm)
	}

	return rows, nil
}
