# surveygo API Reference

## Contents

- [Survey Struct](#survey-struct)
- [Answer Types](#answer-types)
- [Creation and Parsing](#creation-and-parsing)
- [Serialization](#serialization)
- [Validation](#validation)
- [Question Operations](#question-operations)
- [Group Operations](#group-operations)
- [Answer Operations](#answer-operations)
- [Query Operations](#query-operations)
- [Resume Types](#resume-types)
- [Reviewer Package](#reviewer-package)

## Survey Struct

```go
// survey.go
type Survey struct {
    NameId      string                           `json:"nameId"`
    Title       string                           `json:"title,omitempty"`
    Version     string                           `json:"version,omitempty"`
    Description *string                          `json:"description,omitempty"`
    Questions   map[string]*question.Question    `json:"questions,omitempty"`
    Groups      map[string]*question.Group       `json:"groups,omitempty"`
    GroupsOrder []string                         `json:"groupsOrder,omitempty"`
    Metadata    map[string]any                   `json:"metadata,omitempty"`
}
```

## Answer Types

```go
// survey.go
type Answers map[string][]any  // key: question or group nameId, value: answer array
```

For grouped answers, value elements are `map[string][]any` (each map is one answer set).

## Creation and Parsing

```go
// operation_serde.go
func ParseFromJsonStr(jsonSurvey string) (*Survey, error)  // parse + validate + consistency check + position update
func ParseFromBytes(b []byte) (*Survey, error)              // same as above from bytes

// operation.go
func NewSurvey(title, version string, description *string) (*Survey, error)  // empty survey, no NameId set
```

## Serialization

```go
// operation_serde.go
func (s *Survey) ToMap() (map[string]any, error)
func (s *Survey) ToJson() (string, error)
```

## Validation

```go
// operation.go
func (s *Survey) ValidateSurvey() error  // runs checkConsistency()

// validator.go
var SurveyValidator *validator.Validate                    // global validator instance
func TranslateValidationErrors(err error) []error          // translates validator errors to readable messages

// Custom validators registered:
// - "questionType": valid question type string
// - "validNameId": matches ^[a-zA-Z][a-zA-Z\d_-]{1,62}[a-zA-Z\d]$
// - "validIfExternalSurvey": if IsExternalSurvey, QuestionsIds length must be 1
```

## Question Operations

All methods auto-validate and run consistency checks.

```go
// operation_question.go

// Add (error if nameId exists)
func (s *Survey) AddQuestion(q *question.Question) error
func (s *Survey) AddQuestionJson(qs string) error
func (s *Survey) AddQuestionBytes(qb []byte) error
func (s *Survey) AddQuestionMap(qm map[string]any) error

// AddOrUpdate (upsert)
func (s *Survey) AddOrUpdateQuestion(q *question.Question) error
func (s *Survey) AddOrUpdateQuestionJson(qs string) error
func (s *Survey) AddOrUpdateQuestionBytes(qb []byte) error
func (s *Survey) AddOrUpdateQuestionMap(qm map[string]any) error

// Update (error if nameId not found)
func (s *Survey) UpdateQuestion(uq *question.Question) error
func (s *Survey) UpdateQuestionJson(uq string) error
func (s *Survey) UpdateQuestionBytes(uq []byte) error
func (s *Survey) UpdateQuestionMap(uq map[string]any) error

// Remove — also cleans DependsOn references and removes from groups
func (s *Survey) RemoveQuestion(questionNameId string) error

// Group placement
func (s *Survey) AddQuestionToGroup(questionNameId, groupNameId string, position int) error  // position -1 = append
func (s *Survey) RemoveQuestionFromGroup(questionNameId, groupNameId string) error
func (s *Survey) UpdateGroupQuestions(groupNameId string, questionsIds []string) error

// Query
func (s *Survey) GetQuestionsAssignments() map[string]string  // question nameId -> group nameId ("" if unassigned)
func (s *Survey) GetAssetQuestions() []*question.Question       // all image/video/audio/document questions
```

## Group Operations

```go
// operation_group.go

// Add
func (s *Survey) AddGroup(g *question.Group) error
func (s *Survey) AddGroupJson(gs string) error
func (s *Survey) AddGroupBytes(gb []byte) error
func (s *Survey) AddGroupMap(g map[string]any) error

// Update
func (s *Survey) UpdateGroup(ug *question.Group) error
func (s *Survey) UpdateGroupJson(ug string) error
func (s *Survey) UpdateGroupBytes(ug []byte) error
func (s *Survey) UpdateGroupMap(ug map[string]any) error

// AddOrUpdate (upsert)
func (s *Survey) AddOrUpdateGroupJson(g string) error
func (s *Survey) AddOrUpdateGroupBytes(g []byte) error
func (s *Survey) AddOrUpdateGroupMap(g map[string]any) error

// Remove — also removes from option GroupsIds and GroupsOrder
func (s *Survey) RemoveGroup(groupNameId string) error

// Order and state
func (s *Survey) UpdateGroupsOrder(order []string) error
func (s *Survey) EnableGroup(groupNameId string) error
func (s *Survey) DisableGroup(groupNameId string) error
```

## Answer Operations

```go
// operation_answers.go
func (s *Survey) ReviewAnswers(ans Answers) (*SurveyResume, error)
func (s *Survey) TranslateAnswers(ans Answers, ignoreUnknown bool) (Answers, error)
func (s *Survey) GroupAnswersByType(ans Answers) map[types.QuestionType]Answers
```

**TranslateAnswers behavior:**

- Text types: value passed through unchanged
- Simple choice types: nameId -> Option.Value (if set) or nameId
- `ignoreUnknown: true` skips unknown nameIds instead of erroring

## Query Operations

```go
// operation_answers.go
func (s *Survey) GetDisabledQuestions() map[string]bool    // question disabled OR group disabled
func (s *Survey) GetEnabledQuestions() map[string]bool     // question enabled AND group enabled
func (s *Survey) GetRequiredQuestions() map[string]bool
func (s *Survey) GetOptionalQuestions() map[string]bool
func (s *Survey) GetRequiredAndOptionalQuestions() map[string]bool  // value: true=required, false=optional
```

## Resume Types

```go
// operation.go
type SurveyResume struct {
    TotalsResume                                              // inline
    ExternalSurveyIds map[string]string                       // groupNameId -> externalSurveyId
    GroupsResume      map[string]*GroupTotalsResume            // per-group stats
    InvalidAnswers    []*InvalidAnswerError                    // validation errors
}

type TotalsResume struct {
    TotalQuestions                 int
    TotalRequiredQuestions         int
    TotalQuestionsAnswered         int
    TotalRequiredQuestionsAnswered int
    UnansweredQuestions            map[string]bool             // nameId -> isRequired
}

type GroupTotalsResume struct {
    TotalsResume                                              // inline
    AnswerGroups int                                          // number of repeated answer sets
}

type InvalidAnswerError struct {
    QuestionNameId string
    Answer         any
    Error          string
}
```

## Reviewer Package

```go
// reviewer/reviewer.go
type QuestionReviewer func(question any, answers []any, qt types.QuestionType) error
type GroupAnswers []map[string][]any

func GetQuestionReviewer(qt types.QuestionType) (QuestionReviewer, error)
func ExtractGroupNestedAnswers(groupAnswersPack []any) (GroupAnswers, error)
```

`GetQuestionReviewer` returns type-specific validators: `ReviewChoice`, `ReviewText`, `ReviewAsset`, `ReviewExternal`.
