# Question Types Reference

## Contents

- [BaseQuestion](#basequestion)
- [Group](#group)
- [DependsOn](#dependson)
- [QBase (Value Common Fields)](#qbase)
- [Choice Types](#choice-types)
- [Slider](#slider)
- [Text Types](#text-types)
- [Asset Types](#asset-types)
- [External Type](#external-type)
- [Type Constants](#type-constants)
- [Type Check Functions](#type-check-functions)

## BaseQuestion

All questions embed `BaseQuestion`. File: `question/question.go`

```go
type BaseQuestion struct {
    NameId     string              `json:"nameId"`      // required, validNameId
    Visible    bool                `json:"visible"`
    QTyp       types.QuestionType  `json:"type"`        // required, questionType
    Label      string              `json:"label"`       // omitempty, min=1
    Required   bool                `json:"required"`
    Metadata   map[string]any      `json:"metadata"`
    Position   int                 `json:"position"`    // auto-calculated
    Disabled   bool                `json:"disabled"`
    DependsOn  [][]DependsOn       `json:"dependsOn"`   // OR of ANDs
    AnswerExpr string              `json:"answerExpr"`   // optional expr-lang expression
}

type Question struct {
    BaseQuestion                    // inline
    Value any                       `json:"value"`  // type-specific struct
}
```

## Group

File: `question/group.go`

```go
type Group struct {
    NameId           string         `json:"nameId"`           // required, validNameId
    Title            *string        `json:"title"`            // NOTE: pointer
    Description      *string        `json:"description"`
    Hidden           bool           `json:"hidden"`
    Disabled         bool           `json:"disabled"`
    IsExternalSurvey bool           `json:"isExternalSurvey"` // if true, QuestionsIds length must be 1
    AllowRepeat      bool           `json:"allowRepeat"`      // enables grouped/repeatable answers
    QuestionsIds     []string       `json:"questionsIds"`
    GroupsOrder      []string       `json:"groupsOrder"`      // nested groups
    Metadata         map[string]any `json:"metadata"`
    Position         int            `json:"position"`         // auto-calculated
    DependsOn        [][]DependsOn  `json:"dependsOn"`
}

func (g *Group) RemoveQuestionId(nameId string) bool
func (g *Group) AddQuestionId(nameId string, position int)  // position -1 = append
```

## DependsOn

File: `question/depends_on.go`

```go
type DependsOn struct {
    QuestionNameId string `json:"questionNameId"`  // must reference a choice-type question
    OptionNameId   string `json:"optionNameId"`    // must exist on that question
}
```

## QBase

Common fields embedded in all type-specific value structs. File: `question/types/types.go`

```go
type QBase struct {
    Placeholder *string        `json:"placeholder"`
    Metadata    map[string]any `json:"metadata"`
    Collapsible *bool          `json:"collapsible"`
    Collapsed   *bool          `json:"collapsed"`
    Color       *string        `json:"color"`
    Defaults    []string       `json:"defaults"`
}
```

## Choice Types

File: `question/types/choice/choice.go`

Types: `single_select`, `multi_select`, `radio`, `checkbox`, `toggle`

```go
type Choice struct {
    types.QBase
    Options []*Option  // required, min=1
}

type Option struct {
    NameId    string         `json:"nameId"`     // required, validNameId
    Label     string         `json:"label"`      // required, min=1
    Value     any            `json:"value"`      // optional — used in TranslateAnswers
    GroupsIds []string       `json:"groupsIds"`  // groups shown when this option is selected
    Metadata  map[string]any `json:"metadata"`
}

func CastToChoice(i any) (*Choice, error)
func (c *Choice) GetOptionsGroups() map[string][]string  // optionNameId -> groupIds
func (c *Choice) RemoveGroupId(groupId string) bool
```

**Answer format**: option nameId as string in the answers array.

## Slider

File: `question/types/choice/slider.go`

Type: `slider` (complex choice — no Options, uses Min/Max/Step)

```go
type Slider struct {
    types.QBase
    Min     int    `json:"min"`      // required
    Max     int    `json:"max"`      // required
    Step    int    `json:"step"`     // required, min=1
    Default int    `json:"default"`
    Unit    string `json:"unit"`
}

func CastToSlider(i any) (*Slider, error)
```

## Text Types

### FreeText

Types: `text_area`, `input_text`. File: `question/types/text/freetext.go`

```go
type FreeText struct {
    types.QBase
    Min *int  // optional, min=0
    Max *int  // optional, min=0
}

func CastToFreeText(questionValue any) (*FreeText, error)
```

### Email

Type: `email`. File: `question/types/text/email.go`

```go
type Email struct {
    types.QBase
    AllowedDomains []string  // optional, restrict to specific domains
}

func CastToEmail(questionValue any) (*Email, error)
```

### Telephone

Type: `telephone`. File: `question/types/text/telephone.go`

```go
type Telephone struct {
    types.QBase
    AllowedCountryCodes []string  // optional
}

func CastToTelephone(questionValue any) (*Telephone, error)
```

### DateTime

Type: `date_time`. File: `question/types/text/datetime.go`

```go
type DateTime struct {
    types.QBase
    Format string         // required (e.g., "2006-01-02")
    Type   DateTypeFormat // required: "date", "time", or "datetime"
}

type DateTypeFormat string  // "date" | "time" | "datetime"

func CastToDateTime(questionValue any) (*DateTime, error)
```

### InformationText

Type: `information`. File: `question/types/text/information.go`

```go
type InformationText struct {
    types.QBase
    Text string  // required, min=1
}
```

### IdentificationNumber

Type: `identification_number`. File: `question/types/text/identification_number.go`

```go
type IdentificationNumber struct {
    types.QBase
}
```

## Asset Types

All asset types share a common structure. Files: `question/types/asset/*.go`

### ImageAsset (type: `image`)

```go
type ImageAsset struct {
    types.QBase
    AltText             *string
    Tags                []string
    Metadata            map[string]any
    MaxSize             *int64           // bytes, must be > 0
    AllowedContentTypes []string
    MaxFiles            int              // 0 = treat as 1
    MinFiles            int              // 0 = treat as 1
}
```

### VideoAsset (type: `video`)

```go
type VideoAsset struct {
    types.QBase
    Caption             *string
    MaxSize             *int64
    Tags                []string
    Metadata            map[string]any
    AllowedContentTypes []string
    MaxFiles            int
    MinFiles            int
}
```

### AudioAsset (type: `audio`)

```go
type AudioAsset struct {
    types.QBase
    Caption             *string
    MaxSize             *int64
    Tags                []string
    Metadata            map[string]any
    AllowedContentTypes []string
    MaxFiles            int
    MinFiles            int
}
```

### DocumentAsset (type: `document`)

```go
type DocumentAsset struct {
    types.QBase
    Caption             *string
    MaxSize             *int64
    Tags                []string
    Metadata            map[string]any
    AllowedContentTypes []string
    MaxFiles            int
    MinFiles            int
}
```

**Asset gotcha**: `MaxFiles` and `MinFiles` default to `0` (omitempty). Consuming code must treat `0` as `1`. The library does NOT validate file count constraints.

## External Type

Type: `external_question`. File: `question/types/external/external.go`

```go
type ExternalQuestion struct {
    types.QBase
    Defaults     []string
    ExternalType string   // required, min=1
    Description  *string
    Src          *string
}
```

**No CastToExternal exists.** Use direct type assertion:

```go
ext := q.Value.(*external.ExternalQuestion)
```

## Type Constants

File: `question/types/types.go`

| Category | Type String             | Go Constant                 |
| -------- | ----------------------- | --------------------------- |
| Choice   | `single_select`         | `QTypeSingleSelect`         |
| Choice   | `multi_select`          | `QTypeMultipleSelect`       |
| Choice   | `radio`                 | `QTypeRadio`                |
| Choice   | `checkbox`              | `QTypeCheckbox`             |
| Choice   | `toggle`                | `QTypeToggle`               |
| Choice   | `slider`                | `QTypeSlider`               |
| Text     | `text_area`             | `QTypeTextArea`             |
| Text     | `input_text`            | `QTypeInputText`            |
| Text     | `email`                 | `QTypeEmail`                |
| Text     | `telephone`             | `QTypeTelephone`            |
| Text     | `information`           | `QTypeInformation`          |
| Text     | `identification_number` | `QTypeIdentificationNumber` |
| Text     | `date_time`             | `QTypeDateTime`             |
| Asset    | `image`                 | `QTypeImage`                |
| Asset    | `video`                 | `QTypeVideo`                |
| Asset    | `audio`                 | `QTypeAudio`                |
| Asset    | `document`              | `QTypeDocument`             |
| External | `external_question`     | `QTypeExternalQuestion`     |

## Type Check Functions

```go
func IsChoiceType(qt QuestionType) bool        // all choice types including slider
func IsSimpleChoiceType(qt QuestionType) bool   // single_select, multi_select, radio, checkbox
func IsComplexChoiceType(qt QuestionType) bool  // toggle, slider
func IsTextType(qt QuestionType) bool
func IsAssetType(qt QuestionType) bool
func IsExternalType(qt QuestionType) bool
func ParseToQuestionType(v string) (QuestionType, error)
```
