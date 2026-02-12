package render

import "bytes"

// GroupNode represents a node in the group hierarchy tree.
type GroupNode struct {
	NameId      string       `json:"nameId"`
	AllowRepeat bool         `json:"allowRepeat,omitempty"`
	Children    []*GroupNode `json:"children,omitempty"`
}

// GroupTree holds the hierarchical tree and a flat index for O(1) lookup.
type GroupTree struct {
	Roots []*GroupNode          `json:"roots"`
	Index map[string]*GroupNode `json:"-"`
}

// QuestionInfo is the processed output for a question.
type QuestionInfo struct {
	NameId       string       `json:"nameId"`
	Label        string       `json:"label,omitempty"`
	QuestionType string       `json:"questionType"`
	Format       string       `json:"format,omitempty"`
	ExternalType string       `json:"externalType,omitempty"`
	Options      []OptionInfo `json:"options,omitempty"`
	AnswerExpr   string       `json:"answerExpr,omitempty"`
}

// OptionInfo is the processed output for a select/choice option.
type OptionInfo struct {
	NameId string `json:"nameId"`
	Label  string `json:"label"`
	Value  any    `json:"value,omitempty"`
}

// GroupQuestions maps a group to its direct questions.
type GroupQuestions struct {
	GroupNameId string         `json:"groupNameId"`
	Questions   []QuestionInfo `json:"questions"`
}

// SurveyCard is the JSON-renderable card output for a survey response.
type SurveyCard struct {
	SurveyId string    `json:"surveyId"`
	Title    string    `json:"title"`
	Sections []Section `json:"sections"`
}

// Section represents a group rendered as a card section.
// Type determines which fields are populated: "group" uses Fields,
// "repeat-table" uses Columns+Rows, "repeat-list" uses Instances.
type Section struct {
	Type      string     `json:"type"`
	NameId    string     `json:"nameId"`
	Title     string     `json:"title"`
	Fields    []Field    `json:"fields,omitempty"`
	Columns   []Column   `json:"columns,omitempty"`
	Rows      []Row      `json:"rows,omitempty"`
	Instances []Instance `json:"instances,omitempty"`
	Sections  []Section  `json:"sections,omitempty"`
}

// Field represents a single question rendered inside a "group" section.
type Field struct {
	Type   string `json:"type"`
	NameId string `json:"nameId"`
	Label  string `json:"label"`
	Value  any    `json:"value"`
}

// Column describes a column in a "repeat-table" section.
type Column struct {
	NameId    string      `json:"nameId"`
	Label     string      `json:"label"`
	FieldType string      `json:"fieldType"`
	Options   []OptionRef `json:"options,omitempty"`
}

// Row is a single row in a "repeat-table" section, keyed by column nameId.
type Row = map[string]any

// Instance is a single entry in a "repeat-list" section.
type Instance struct {
	Sections []Section `json:"sections"`
}

// OptionRef represents an option in a multi-select column or field.
type OptionRef struct {
	NameId   string `json:"nameId"`
	Label    string `json:"label"`
	Selected bool   `json:"selected"`
}

// TipTapNode represents a node in a TipTap/ProseMirror document tree.
type TipTapNode struct {
	Type    string         `json:"type"`
	Attrs   map[string]any `json:"attrs,omitempty"`
	Content []TipTapNode   `json:"content,omitempty"`
	Text    string         `json:"text,omitempty"`
	Marks   []TipTapMark   `json:"marks,omitempty"`
}

// TipTapMark represents an inline mark (bold, italic, etc.) on a text node.
type TipTapMark struct {
	Type  string         `json:"type"`
	Attrs map[string]any `json:"attrs,omitempty"`
}

// OutputOptions specifies which output formats to generate.
type OutputOptions struct {
	CSV    bool
	JSON   bool
	HTML   bool
	TipTap bool
}

// HTMLResult contains HTML body and CSS as separate byte slices.
type HTMLResult struct {
	HTML []byte `json:"html"`
	CSS  []byte `json:"css"`
}

// WithCSSPath returns a copy of HTMLResult with the CSS link href replaced in the HTML.
func (r *HTMLResult) WithCSSPath(cssPath string) *HTMLResult {
	replaced := bytes.Replace(r.HTML, []byte(`href="card.css"`), []byte(`href="`+cssPath+`"`), 1)
	return &HTMLResult{HTML: replaced, CSS: r.CSS}
}

// AnswersResult contains the requested output formats.
// Fields are nil when not requested via OutputOptions.
type AnswersResult struct {
	CSV    []byte      `json:"csv,omitempty"`
	JSON   *SurveyCard `json:"json,omitempty"`
	HTML   *HTMLResult `json:"html,omitempty"`
	TipTap *TipTapNode `json:"tiptap,omitempty"`
}

// TreeResult contains both representations of the group tree.
type TreeResult struct {
	HTML []byte     `json:"html,omitempty"`
	JSON *GroupTree `json:"json,omitempty"`
}
