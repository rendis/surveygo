# AGENTS.md

This file provides guidance to AI coding agents when working with code in this repository.

For user-facing documentation (API, structures, usage), see [README.md](README.md). For JSON structure reference, see [docs/SURVEY_STRUCTURE.md](docs/SURVEY_STRUCTURE.md).

## Project Overview

SurveyGo is a Go library for creating and managing surveys with comprehensive validation, serialization, and answer processing capabilities. The library provides a structured approach to building dynamic surveys with questions, groups, and conditional logic.

## Common Development Commands

### Building and Testing

- `go build` - Build the project
- `go test ./...` - Run all tests
- `go mod tidy` - Clean up module dependencies
- `go run example/main.go` - Run the example application
- `go vet ./...` - Run Go vet for static analysis
- `go fmt ./...` - Format Go code

### Example Usage

The `example/` directory contains a complete working example showing how to:

- Parse surveys from JSON
- Add questions dynamically
- Review and validate answers
- Handle grouped answers
- Translate answers to human-readable format

## Architecture Overview

### Core Components

**Survey Structure (`survey.go`)**

- `Survey`: Main structure containing questions, groups, and metadata
- `Answers`: Map of question nameIds to answer arrays
- Hierarchical organization: Survey → Groups → Questions

**Operations (`operation*.go`)**

- `operation.go`: Core survey operations, validation, and answer review
- `operation_question.go`: Question management (add, remove, modify)
- `operation_group.go`: Group management and organization
- `operation_answers.go`: Answer processing, validation, and translation
- `operation_serde.go`: JSON serialization/deserialization

**Question System (`question/`)**

- `question.go`: Base question structure and common fields
- `group.go`: Group structure for organizing questions
- `depends_on.go`: DependsOn struct for conditional visibility logic
- `types/`: Question type definitions and implementations:
  - `choice/`: Single/multi-select, radio, checkbox questions
  - `text/`: Email, telephone, freetext, information fields
  - `external/`: External question integration
  - `asset/`: File upload questions (image, video, audio, document)

**Answer Review System (`reviewer/`)**

- Type-specific answer validators for each question category
- Handles validation logic for different question types
- Used by `ReviewAnswers()` to validate user responses

**Render Package (`render/`)**

Survey output generation. Single entry point in `render.go` — all other functions are unexported.

Public API (`render.go`):

- `AnswersToCSV(survey, answers, checkMark...*CheckMark)` → CSV bytes
- `AnswersToJSON(survey, answers)` → `*SurveyCard`
- `AnswersToHTML(survey, answers)` → `*HTMLResult` (HTML + CSS independent). `HTMLResult.WithCSSPath(path)` replaces CSS href
- `AnswersToTipTap(survey, answers)` → `*TipTapNode`
- `AnswersTo(survey, answers, opts)` → `*AnswersResult` (multiple formats, single pass)
- `DefinitionTreeJSON(survey)` → `*GroupTree`
- `DefinitionTreeHTML(survey)` → HTML bytes (go-echarts interactive tree)
- `DefinitionTree(survey)` → `*TreeResult` (HTML + JSON)
- `ReportColumns(survey)` → `([]ReportColumn, *GroupTree, error)` — ordered column definitions for multi-record reports
- `ReportRows(survey, tree, columns, answers, checkMark)` → `([]map[string]string, error)` — row data with cartesian product expansion

Internal files (all unexported):

- `types.go`: Types (GroupTree, SurveyCard, OutputOptions, CheckMark, AnswersResult, TreeResult, ReportColumn, etc.); `TopLevelGroup()` method on `GroupTree` (walks `parent` map to find root-level group)
- `answers.go`: Answer extraction helpers (extractTextValue, extractPhoneValue, etc.); `extractToggleValue` handles both `bool` and string `"true"`/`"1"`
- `tree.go`: `buildGroupTree` — DFS group hierarchy with cycle detection; computes `RepeatDescendants` (count of `AllowRepeat` descendants per node); populates `parent` map (child→parent) during DFS
- `questions.go`: `extractGroupQuestions` — adapts `*question.Question` → `QuestionInfo`
- `expr_eval.go`: expr-lang/expr evaluation with silent fallback
- `card.go`: `buildSurveyCard` — structured card with resolved answers. Section type truth table:
  - `AllowRepeat=false` → `group` (flat)
  - `AllowRepeat=true, RepeatDescendants > 0` → `repeat-list`
  - `AllowRepeat=true, RepeatDescendants=0, has multi_select/checkbox` → `repeat-list`
  - `AllowRepeat=true, RepeatDescendants=0, no multi_select/checkbox` → `repeat-table` (flattens descendant questions into columns)
- `report.go`: `ReportColumns`, `ReportRows` — public API for multi-record report column/row extraction; reuses `buildColumns`/`fillRows` from csv.go
- `csv.go`: `generateCSV` — cartesian product expansion for repeat groups
- `html.go`: `generateHTML`, `defaultCSS` — HTML rendering
- `tiptap.go`: `buildTipTapDoc` — TipTap document
- `visualize.go`: `renderTreeToBytes` — go-echarts tree visualization

Key type adaptations in `questions.go`:

- `choice.CastToChoice(q.Value)` for options extraction
- `text.CastToDateTime(q.Value)` for date format
- Direct type assertion `q.Value.(*external.ExternalQuestion)` for external type (no CastToExternal exists)
- `derefStr(*string)` helper — Group.Title is `*string` in surveygo

**Validation (`validator.go`)**

- Custom validation rules for nameIds, question types
- Struct validation using `github.com/go-playground/validator/v10`
- Internationalization support for validation messages

### Key Design Patterns

**NameId System**: All entities use unique `nameId` strings for identification. NameIds must match regex: `^[a-zA-Z][a-zA-Z\d_-]{1,62}[a-zA-Z\d]$`

**Type Safety**: Strong typing for question types, with type-specific casting and validation

**Conditional Logic**: Two mechanisms for dynamic survey flows:

- **Option-triggered groups**: Choice options can have `groupsIds` that trigger display of specific groups when selected
- **DependsOn**: Questions and groups can have `dependsOn` field for conditional visibility based on other question/option selections

**Consistency Validation**: `checkConsistency()` ensures referential integrity between questions, groups, and options

**Position Management**: Automatic position calculation for questions and groups based on `GroupsOrder`

## Important Implementation Details

### Answer Processing

- Answers are stored as `map[string][]any` where keys are question nameIds
- Multiple answers per question are supported (arrays)
- Grouped answers allow multiple response sets per group
- Translation converts raw answers to human-readable labels

### AnswerExpr (Custom Answer Processing)

Optional `AnswerExpr` field on `BaseQuestion` (`question/question.go`). Evaluated by [expr-lang/expr](https://github.com/expr-lang/expr).

- Environment: `ans` ([]any) + `options` (map[nameId]label, choice types only)
- Silent fallback: returns `(nil, false)` on error — lib never writes to stderr
- Used by render package in `resolveValue()` (`card.go`) to override default type-based extraction
- When empty, default per-type logic applies (ExtractTextValue, ExtractPhoneValue, etc.)

### CheckMark (CSV Boolean Column Customization)

`CheckMark` (`render/types.go`) customizes the strings used for selected/not-selected marks in CSV boolean columns (multi_select, checkbox, toggle).

```go
type CheckMark struct {
    Selected    string // e.g. "x", "Yes", "✅"
    NotSelected string // e.g. "", "No", "-"
}
```

- Pass via `AnswersToCSV(survey, answers, &CheckMark{...})` or `OutputOptions.CheckMark`
- When nil, defaults to `"true"`/`"false"` (backward-compatible)
- Affects: multi_select, checkbox (one column per option), toggle (single column)
- Does NOT affect single_select/radio (these write the option label, not a boolean)

### Survey Structure Rules

- Each question belongs to exactly one group
- Groups are ordered via `GroupsOrder` slice
- External surveys are supported as special group types
- Questions and groups must have unique nameIds across the survey

### DependsOn Implementation

See [README.md](README.md#dependson-conditional-logic) for more detailed documentation.

**Go struct** (`question/depends_on.go`):

```go
type DependsOn struct {
    QuestionNameId string `json:"questionNameId" bson:"questionNameId" validate:"required,validNameId"`
    OptionNameId   string `json:"optionNameId" bson:"optionNameId" validate:"required,validNameId"`
}
```

**Validation** (`operation.go` - `checkConsistency()`):

- Validates referenced `questionNameId` exists in the survey
- Validates referenced question is a choice type (has options)
- Validates referenced `optionNameId` exists on that question

**Cleanup** (`operation_question.go` - `RemoveQuestion()`):

- When a question is removed, `removeDependsOnByQuestion()` automatically cleans up all `dependsOn` references to that question from other questions and groups
- Only removes the specific condition referencing the deleted question
- If an AND group becomes empty, the entire AND group is removed

**Visibility Engine** (`operation_answers.go`):
The visibility engine evaluates `dependsOn` conditions against provided answers during `ReviewAnswers()`.

Functions:

- `evaluateDependsOn(dependsOn, ans)` - Main entry point, evaluates OR logic (any AND group matches = visible)
- `evaluateAndGroup(andGroup, ans)` - Evaluates AND logic (all conditions must be true)
- `evaluateCondition(dep, ans)` - Checks if `questionNameId` exists in answers with `optionNameId` selected

Usage in `getVisibleQuestionFromActiveGroups(ans)`:

```go
// Groups are visible if: !Hidden && !Disabled && dependsOn satisfied
if !s.evaluateDependsOn(group.DependsOn, ans) {
    continue  // skip group
}

// Questions are visible if: Visible && dependsOn satisfied
if !s.evaluateDependsOn(q.DependsOn, ans) {
    continue  // skip question
}
```

**Behavior**: Questions/groups with unsatisfied `dependsOn` are excluded from `getSurveyResume()`, meaning they don't count toward totals and required questions with unsatisfied conditions are not expected to be answered.

### Asset File Constraints

All asset types (image, video, audio, document) in `question/types/asset/` have `MaxFiles` and `MinFiles` fields.

**Default value handling**:

- Both fields default to 0 when not specified in JSON (due to `omitempty`)
- Consuming code should treat 0 as default value of 1
- This is NOT handled by the library - consuming applications must implement this logic

**Validation tags**:

```go
MaxFiles int `json:"maxFiles,omitempty" bson:"maxFiles,omitempty" validate:"omitempty,min=1"`
MinFiles int `json:"minFiles,omitempty" bson:"minFiles,omitempty" validate:"omitempty,min=0"`
```

**Note**: The `reviewer/asset.go` currently does not validate file count constraints - this should be implemented by consuming applications or added to the reviewer in the future.

### Validation Workflow

1. Structural validation using validator tags
2. Consistency checks for cross-references
3. Answer validation using type-specific reviewers
4. Resume generation with totals and error reporting

### Dependencies

- `github.com/go-playground/validator/v10` for struct validation
- `github.com/rendis/devtoolkit` for utilities
- `go.mongodb.org/mongo-driver` for BSON support
- Standard library for JSON handling
- `github.com/expr-lang/expr` for AnswerExpr expression evaluation
- `github.com/go-echarts/go-echarts/v2` for tree visualization in render package

## Agent Skill

An [Agent Skill](https://agentskills.io/specification) is bundled at `skills/surveygo/`. It provides AI coding agents with structured guidance for consuming this library.

**Structure:**

- `skills/surveygo/SKILL.md` — Core: quick start, workflows, gotchas
- `skills/surveygo/references/api.md` — Full exported API (structs, methods, signatures)
- `skills/surveygo/references/question-types.md` — All question types with fields and cast functions
- `skills/surveygo/references/render.md` — Render package API and output types

**Installation (Claude Code):**

```bash
ln -s /path/to/surveygo/skills/surveygo ~/.claude/skills/surveygo
```

A distributable `skills/surveygo.skill` package is also available.
