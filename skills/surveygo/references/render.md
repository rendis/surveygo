# Render Package Reference

## Contents

- [Render Package Reference](#render-package-reference)
  - [Contents](#contents)
  - [Import](#import)
  - [Answer Output Functions](#answer-output-functions)
  - [Definition Tree Functions](#definition-tree-functions)
  - [Tabular Row Output](#tabular-row-output)
  - [Output Types](#output-types)
    - [SurveyCard](#surveycard)
    - [HTMLResult](#htmlresult)
    - [TipTapNode](#tiptapnode)
    - [GroupTree](#grouptree)
    - [Result Types](#result-types)
  - [Options and Helpers](#options-and-helpers)
  - [AnswerExpr](#answerexpr)

## Import

```go
import "github.com/rendis/surveygo/v2/render"
```

## Answer Output Functions

File: `render/render.go`

```go
// Single format outputs
func AnswersToCSV(survey *Survey, answers Answers, checkMark ...*CheckMark) ([]byte, error)
func AnswersToJSON(survey *Survey, answers Answers) (*SurveyCard, error)
func AnswersToHTML(survey *Survey, answers Answers) (*HTMLResult, error)
func AnswersToTipTap(survey *Survey, answers Answers) (*TipTapNode, error)

// Multi-format single pass — only computes formats enabled in opts
func AnswersTo(survey *Survey, answers Answers, opts OutputOptions) (*AnswersResult, error)
```

`AnswersToCSV` accepts optional `CheckMark` to customize selected/not-selected strings for multi-select, checkbox, and toggle columns.

## Definition Tree Functions

Visualize the survey group hierarchy.

```go
func DefinitionTreeJSON(survey *Survey) (*GroupTree, error)    // JSON tree structure
func DefinitionTreeHTML(survey *Survey) ([]byte, error)         // interactive HTML (go-echarts)
func DefinitionTree(survey *Survey) (*TreeResult, error)        // both JSON + HTML
```

## Tabular Row Output

Returns the same data as `AnswersToCSV` but as Go types (`[][]string`) instead of serialized CSV bytes.

File: `render/render.go`

```go
func AnswersToRows(survey *Survey, answers Answers, checkMark ...*CheckMark) ([][]string, error)
```

Returns a matrix where `matrix[0]` is the header row and `matrix[1:]` are data rows. Repeatable groups expand via cartesian product (same logic as `AnswersToCSV`). Optional `CheckMark` controls selected/not-selected strings for boolean columns (multi-select, checkbox, toggle).

## Output Types

### SurveyCard

Structured JSON output for a survey response.

```go
type SurveyCard struct {
    SurveyId string    `json:"surveyId"`
    Title    string    `json:"title"`
    Sections []Section `json:"sections"`
}

// Section.Type determined by truth table:
//   AllowRepeat=false                                          → "group" (Fields)
//   AllowRepeat=true, RepeatDescendants > 0                    → "repeat-list" (Instances)
//   AllowRepeat=true, RepeatDescendants=0, has multi_select    → "repeat-list" (Instances)
//   AllowRepeat=true, RepeatDescendants=0, no multi_select     → "repeat-table" (Columns+Rows, flattens descendant questions)
type Section struct {
    Type      string     `json:"type"`
    NameId    string     `json:"nameId"`
    Title     string     `json:"title"`
    Fields    []Field    `json:"fields,omitempty"`
    Columns   []Column   `json:"columns,omitempty"`
    Rows      []Row      `json:"rows,omitempty"`
    Instances []Instance `json:"instances,omitempty"`
    Sections  []Section  `json:"sections,omitempty"`  // nested sections
}

type Field struct {
    Type   string `json:"type"`
    NameId string `json:"nameId"`
    Label  string `json:"label"`
    Value  any    `json:"value"`
}

type Column struct {
    NameId    string      `json:"nameId"`
    Label     string      `json:"label"`
    FieldType string      `json:"fieldType"`
    Options   []OptionRef `json:"options,omitempty"`
}

type Row = map[string]any

type Instance struct {
    Sections []Section `json:"sections"`
}

type OptionRef struct {
    NameId   string `json:"nameId"`
    Label    string `json:"label"`
    Selected bool   `json:"selected"`
}
```

### HTMLResult

```go
type HTMLResult struct {
    HTML []byte `json:"html"`
    CSS  []byte `json:"css"`   // independent CSS
}

// Replace CSS href in the HTML (default: "card.css")
func (r *HTMLResult) WithCSSPath(cssPath string) *HTMLResult
```

### TipTapNode

TipTap/ProseMirror document tree.

```go
type TipTapNode struct {
    Type    string         `json:"type"`
    Attrs   map[string]any `json:"attrs,omitempty"`
    Content []TipTapNode   `json:"content,omitempty"`
    Text    string         `json:"text,omitempty"`
    Marks   []TipTapMark   `json:"marks,omitempty"`
}

type TipTapMark struct {
    Type  string         `json:"type"`
    Attrs map[string]any `json:"attrs,omitempty"`
}
```

### GroupTree

```go
type GroupTree struct {
    Roots []*GroupNode          `json:"roots"`
    Index map[string]*GroupNode `json:"-"`  // O(1) lookup, not serialized
}

type GroupNode struct {
    NameId            string       `json:"nameId"`
    AllowRepeat       bool         `json:"allowRepeat,omitempty"`
    RepeatDescendants int          `json:"repeatDescendants"`          // count of AllowRepeat descendants (recursive, excludes self)
    Children          []*GroupNode `json:"children,omitempty"`
}
```

### Result Types

```go
type AnswersResult struct {
    CSV    []byte      `json:"csv,omitempty"`
    JSON   *SurveyCard `json:"json,omitempty"`
    HTML   *HTMLResult `json:"html,omitempty"`
    TipTap *TipTapNode `json:"tiptap,omitempty"`
}

type TreeResult struct {
    HTML []byte     `json:"html,omitempty"`
    JSON *GroupTree `json:"json,omitempty"`
}
```

## Options and Helpers

```go
type OutputOptions struct {
    CSV    bool
    JSON   bool
    HTML   bool
    TipTap bool
    CheckMark *CheckMark  // nil = "true"/"false"
}

type CheckMark struct {
    Selected    string
    NotSelected string
}
```

## AnswerExpr

Optional field on `BaseQuestion`. Evaluated by [expr-lang/expr](https://github.com/expr-lang/expr) in the render package's `resolveValue()` function.

**Environment variables:**

- `ans` — `[]any`, the raw answer array
- `options` — `map[string]string`, maps option nameId to label (only for choice types, empty otherwise)

**Behavior:**

- When `AnswerExpr` is empty: default type-based extraction applies (ExtractTextValue, ExtractPhoneValue, etc.)
- When set: expression result replaces default extraction
- On error: silent fallback returns `(nil, false)` — never writes to stderr
- Used by render package to customize how answers appear in output

**Example expressions:**

```plaintext
ans[0]                                    // first answer value
options[ans[0]]                           // translate option nameId to label
len(ans) > 0 ? ans[0] : "N/A"            // with fallback
```
