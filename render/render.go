package render

import (
	"fmt"

	surveygo "github.com/rendis/surveygo/v2"
)

// AnswersToCSV generates a CSV from survey answers.
func AnswersToCSV(survey *surveygo.Survey, answers surveygo.Answers) ([]byte, error) {
	tree, err := buildGroupTree(survey)
	if err != nil {
		return nil, fmt.Errorf("building group tree: %w", err)
	}
	questions, err := extractGroupQuestions(survey)
	if err != nil {
		return nil, fmt.Errorf("extracting questions: %w", err)
	}
	return generateCSV(survey, tree, questions, answers)
}

// AnswersToJSON builds a structured SurveyCard from survey answers.
func AnswersToJSON(survey *surveygo.Survey, answers surveygo.Answers) (*SurveyCard, error) {
	tree, err := buildGroupTree(survey)
	if err != nil {
		return nil, fmt.Errorf("building group tree: %w", err)
	}
	questions, err := extractGroupQuestions(survey)
	if err != nil {
		return nil, fmt.Errorf("extracting questions: %w", err)
	}
	return buildSurveyCard(survey, tree, questions, answers)
}

// AnswersToHTML renders survey answers as HTML and CSS independently.
func AnswersToHTML(survey *surveygo.Survey, answers surveygo.Answers) (*HTMLResult, error) {
	card, err := AnswersToJSON(survey, answers)
	if err != nil {
		return nil, err
	}
	html, err := generateHTML(card)
	if err != nil {
		return nil, err
	}
	return &HTMLResult{HTML: html, CSS: defaultCSS()}, nil
}

// AnswersToTipTap builds a TipTap-compatible document from survey answers.
func AnswersToTipTap(survey *surveygo.Survey, answers surveygo.Answers) (*TipTapNode, error) {
	card, err := AnswersToJSON(survey, answers)
	if err != nil {
		return nil, err
	}
	doc := buildTipTapDoc(card)
	return &doc, nil
}

// AnswersTo generates multiple output formats in a single pass.
// Only the formats enabled in opts are computed.
func AnswersTo(survey *surveygo.Survey, answers surveygo.Answers, opts OutputOptions) (*AnswersResult, error) {
	tree, err := buildGroupTree(survey)
	if err != nil {
		return nil, fmt.Errorf("building group tree: %w", err)
	}
	questions, err := extractGroupQuestions(survey)
	if err != nil {
		return nil, fmt.Errorf("extracting questions: %w", err)
	}

	result := &AnswersResult{}

	if opts.CSV {
		result.CSV, err = generateCSV(survey, tree, questions, answers)
		if err != nil {
			return nil, fmt.Errorf("generating CSV: %w", err)
		}
	}

	if opts.JSON || opts.HTML || opts.TipTap {
		card, cardErr := buildSurveyCard(survey, tree, questions, answers)
		if cardErr != nil {
			return nil, fmt.Errorf("building survey card: %w", cardErr)
		}

		if opts.JSON {
			result.JSON = card
		}
		if opts.HTML {
			htmlBytes, htmlErr := generateHTML(card)
			if htmlErr != nil {
				return nil, fmt.Errorf("generating HTML: %w", htmlErr)
			}
			result.HTML = &HTMLResult{HTML: htmlBytes, CSS: defaultCSS()}
		}
		if opts.TipTap {
			doc := buildTipTapDoc(card)
			result.TipTap = &doc
		}
	}

	return result, nil
}

// DefinitionTreeHTML renders the survey group hierarchy as interactive HTML bytes.
func DefinitionTreeHTML(survey *surveygo.Survey) ([]byte, error) {
	tree, err := buildGroupTree(survey)
	if err != nil {
		return nil, fmt.Errorf("building group tree: %w", err)
	}
	return renderTreeToBytes(tree)
}

// DefinitionTreeJSON builds the hierarchical group tree with cycle detection.
func DefinitionTreeJSON(survey *surveygo.Survey) (*GroupTree, error) {
	return buildGroupTree(survey)
}

// DefinitionTree returns both HTML and JSON representations of the group tree.
func DefinitionTree(survey *surveygo.Survey) (*TreeResult, error) {
	tree, err := buildGroupTree(survey)
	if err != nil {
		return nil, fmt.Errorf("building group tree: %w", err)
	}
	html, err := renderTreeToBytes(tree)
	if err != nil {
		return nil, fmt.Errorf("rendering tree HTML: %w", err)
	}
	return &TreeResult{HTML: html, JSON: tree}, nil
}
