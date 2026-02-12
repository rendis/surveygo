# Survey JSON Structure Reference

Quick reference for understanding SurveyGo JSON structure.

## Top-Level Survey

```json
{
  "nameId": "survey_identifier",
  "title": "Survey Title",
  "version": "1.0.0",
  "description": "Optional description",
  "questions": {
    /* map of question objects */
  },
  "groups": {
    /* map of group objects */
  },
  "groupsOrder": ["group1", "group2"],
  "metadata": {
    /* optional metadata */
  }
}
```

| Field         | Type     | Description                              |
| ------------- | -------- | ---------------------------------------- |
| `nameId`      | string   | Unique survey identifier (required)      |
| `title`       | string   | Survey title (required)                  |
| `version`     | string   | Version string (required)                |
| `description` | ?string  | Optional description (nullable)          |
| `questions`   | object   | Map of question objects, keyed by nameId |
| `groups`      | object   | Map of group objects, keyed by nameId    |
| `groupsOrder` | array    | Order of group nameIds for display       |
| `metadata`    | object   | Optional additional data                 |

## Groups

**What is a group?** Container that organizes related questions. Groups control display order and conditional visibility.

```json
"groups": {
  "grp_contact": {
    "nameId": "grp_contact",
    "title": "Contact Information",
    "description": "Tell us how to reach you",
    "questionsIds": ["name", "email", "phone"],
    "groupsOrder": ["nested_grp_1"],
    "hidden": false,
    "disabled": false,
    "isExternalSurvey": false,
    "allowRepeat": false,
    "dependsOn": [ /* conditional logic */ ],
    "metadata": { }
  }
}
```

| Field              | Type    | Description                                          |
| ------------------ | ------- | ---------------------------------------------------- |
| `nameId`           | string  | Unique group identifier (required)                   |
| `title`            | ?string | Group title (optional, nullable)                     |
| `description`      | ?string | Group description (optional, nullable)               |
| `questionsIds`     | array   | List of question nameIds in this group (required)    |
| `groupsOrder`      | array   | Order of nested group nameIds (optional)             |
| `hidden`           | boolean | Hide group (default: false)                          |
| `disabled`         | boolean | Disable group (default: false)                       |
| `isExternalSurvey` | boolean | Mark as external survey (default: false)             |
| `allowRepeat`      | boolean | Allow repeating group (default: false)               |
| `dependsOn`        | array   | Conditional visibility rules                         |
| `metadata`         | object  | Optional additional data                             |
| `position`         | number  | Auto-calculated display position                     |

## Questions

**What is a question?** Individual survey item that collects user input.

```json
"questions": {
  "event_rating": {
    "nameId": "event_rating",
    "visible": true,
    "type": "radio",
    "label": "How would you rate the event?",
    "required": true,
    "value": { /* type-specific config */ },
    "dependsOn": [ /* conditional logic */ ],
    "metadata": { }
  }
}
```

| Field       | Type    | Description                                 |
| ----------- | ------- | ------------------------------------------- |
| `nameId`    | string  | Unique question identifier (required)       |
| `visible`   | boolean | Question visibility (default: true)         |
| `type`      | string  | Question type (required, see types below)   |
| `label`     | string  | Display text for the question               |
| `required`  | boolean | Whether answer is required (default: false) |
| `value`     | object  | Type-specific configuration (required)      |
| `dependsOn` | array   | Conditional visibility rules                |
| `metadata`  | object  | Optional additional data                    |
| `position`  | number  | Auto-calculated display position            |
| `disabled`  | boolean | Whether question is disabled                |
| `answerExpr`| string  | Optional [expr-lang/expr](https://github.com/expr-lang/expr) expression for custom answer processing. Overrides default type-based extraction. Env: `ans` ([]any) + `options` (map[nameId→label], choice types only) |

## Question Types

### Common Value Fields (QBase)

All question type values inherit these optional fields:

| Field         | Type    | Description                                  |
| ------------- | ------- | -------------------------------------------- |
| `placeholder` | ?string | Placeholder/hint text (optional, nullable)   |
| `metadata`    | object  | Additional metadata for the value (optional) |
| `collapsible` | ?bool   | Whether the question is collapsible          |
| `collapsed`   | ?bool   | Whether the question starts collapsed        |
| `color`       | ?string | Color coding for the question                |
| `defaults`    | array   | List of default values (optional)            |

### Choice-Based — Simple (with options)

Questions with selectable options:

| Type            | Description                        | Use Case                       |
| --------------- | ---------------------------------- | ------------------------------ |
| `single_select` | Dropdown/select with single choice | Long list of options           |
| `multi_select`  | Multiple selection dropdown        | Select multiple from long list |
| `radio`         | Radio buttons, single choice       | 2-5 mutually exclusive options |
| `checkbox`      | Checkboxes, multiple selection     | 2-5 options, multiple allowed  |

**Value structure:**

```json
"value": {
  "placeholder": "Select an option",
  "options": [
    {
      "nameId": "option_id",
      "label": "Display Label",
      "value": "optional_value",
      "groupsIds": ["grp_to_show_if_selected"],
      "metadata": { }
    }
  ]
}
```

### Choice-Based — Toggle

| Type     | Description   | Use Case               |
| -------- | ------------- | ---------------------- |
| `toggle` | Toggle switch | On/off, yes/no choices |

**Value structure:**

```json
"value": {
  "onLabel": "Yes",
  "offLabel": "No",
  "default": false
}
```

| Field      | Type    | Description                              |
| ---------- | ------- | ---------------------------------------- |
| `onLabel`  | string  | Label for the "on" state (required)      |
| `offLabel` | string  | Label for the "off" state (required)     |
| `default`  | boolean | Default toggle state (optional)          |

### Choice-Based — Slider

| Type     | Description    | Use Case                |
| -------- | -------------- | ----------------------- |
| `slider` | Slider control | Numeric range selection |

**Value structure:**

```json
"value": {
  "min": 0,
  "max": 100,
  "step": 1,
  "default": 50,
  "unit": "years"
}
```

| Field     | Type   | Description                                      |
| --------- | ------ | ------------------------------------------------ |
| `min`     | number | Minimum slider value (required)                  |
| `max`     | number | Maximum slider value (required)                  |
| `step`    | number | Step increment, min 1 (required)                 |
| `default` | number | Default slider value (optional)                  |
| `unit`    | string | Unit label, e.g. "years", "months" (optional)    |

### Text-Based

Text input questions:

| Type                    | Description                  | Type-Specific Value Fields                          |
| ----------------------- | ---------------------------- | --------------------------------------------------- |
| `input_text`            | Single-line text             | `min`, `max`                                        |
| `text_area`             | Multi-line text              | `min`, `max`                                        |
| `email`                 | Email input with validation  | `allowedDomains`                                    |
| `telephone`             | Phone number input           | `allowedCountryCodes`                               |
| `identification_number` | ID number input              | _(QBase fields only)_                               |
| `date_time`             | Date/time picker             | `type`, `format`                                    |
| `information`           | Display-only text (no input) | `text`                                              |

**input_text / text_area value structure:**

```json
"value": {
  "placeholder": "Enter text here",
  "min": 10,
  "max": 500
}
```

**email value structure:**

```json
"value": {
  "placeholder": "user@example.com",
  "allowedDomains": ["example.com", "company.org"]
}
```

**telephone value structure:**

```json
"value": {
  "placeholder": "+1 555-0100",
  "allowedCountryCodes": ["+1", "+44", "+56"]
}
```

**date_time value structure:**

```json
"value": {
  "placeholder": "Select a date",
  "type": "date",
  "format": "2006-01-02"
}
```

`type` options: `date`, `time`, `datetime`

**information value structure:**

```json
"value": {
  "text": "This is informational text displayed to the user."
}
```

### Asset-Based (file uploads)

File upload questions. All share base fields (`maxFiles`, `minFiles`, `maxSize`, `allowedContentTypes`, `tags`) plus type-specific fields.

| Type       | Description          | Type-Specific Fields |
| ---------- | -------------------- | -------------------- |
| `image`    | Image file upload    | `altText`            |
| `video`    | Video file upload    | `caption`            |
| `audio`    | Audio file upload    | `caption`            |
| `document` | Document file upload | `caption`            |

**image value structure:**

```json
"value": {
  "maxFiles": 3,
  "minFiles": 1,
  "maxSize": 5000000,
  "allowedContentTypes": ["image/jpeg", "image/png"],
  "altText": "Description",
  "tags": ["tag1", "tag2"]
}
```

**video / audio / document value structure:**

```json
"value": {
  "maxFiles": 3,
  "minFiles": 1,
  "maxSize": 5000000,
  "allowedContentTypes": ["video/mp4"],
  "caption": "Optional description",
  "tags": ["tag1", "tag2"]
}
```

| Field                 | Type   | Description                                           |
| --------------------- | ------ | ----------------------------------------------------- |
| `maxFiles`            | number | Max files allowed (optional, 0 = treat as 1)         |
| `minFiles`            | number | Min files required (optional, 0 = treat as 1)        |
| `maxSize`             | number | Max file size in bytes (optional)                     |
| `allowedContentTypes` | array  | Allowed MIME types (optional)                         |
| `tags`                | array  | Keywords/tags (optional, all types)                   |
| `altText`             | string | Alt text description (optional, `image` only)         |
| `caption`             | string | Caption text, max 255 chars (optional, `video`/`audio`/`document` only) |

### External

| Type                | Description                      |
| ------------------- | -------------------------------- |
| `external_question` | Integration with external survey |

**Value structure:**

```json
"value": {
  "externalType": "nps_survey",
  "description": "Net Promoter Score survey",
  "src": "https://example.com/survey",
  "defaults": ["default_val"]
}
```

| Field          | Type    | Description                                  |
| -------------- | ------- | -------------------------------------------- |
| `externalType` | string  | Type identifier for external question (required) |
| `description`  | ?string | Description text (optional, nullable)        |
| `src`          | ?string | Source URL/reference (optional, nullable)     |
| `defaults`     | array   | Default values (optional)                    |

## NameId Format

**What is a nameId?** Unique identifier used for surveys, groups, questions, and options.

**Format:** Must match regex `^[a-zA-Z][a-zA-Z\d_-]{1,62}[a-zA-Z\d]$`

**Rules:**

- Start with letter (a-z, A-Z)
- End with letter or digit
- Middle: letters, digits, underscore, hyphen
- Length: 3-64 characters

**Valid:**

- `event_rating`
- `grp-contact-info`
- `Q1`
- `user_email_2024`

**Invalid:**

- `_private` (starts with underscore)
- `rating-` (ends with hyphen)
- `q` (too short)
- `9options` (starts with digit)

## DependsOn Conditional Logic

**Purpose:** Show/hide questions or groups based on answers to other questions.

**Structure:** Array of arrays = OR of ANDs

- Outer array: OR logic (any inner array matches = visible)
- Inner arrays: AND logic (all conditions must match)

```json
"dependsOn": [
  [
    { "questionNameId": "rating", "optionNameId": "bad" }
  ],
  [
    { "questionNameId": "rating", "optionNameId": "okay" },
    { "questionNameId": "attend_again", "optionNameId": "no" }
  ]
]
```

**Interpretation:**

- Show if: `rating == "bad"` OR (`rating == "okay"` AND `attend_again == "no"`)

**Simple example (single condition):**

```json
"dependsOn": [
  [
    { "questionNameId": "contact_preference", "optionNameId": "yes" }
  ]
]
```

**Complex example (multiple conditions):**

```json
"dependsOn": [
  [
    { "questionNameId": "event_rating", "optionNameId": "terrible" },
    { "questionNameId": "would_attend", "optionNameId": "no" }
  ],
  [
    { "questionNameId": "event_rating", "optionNameId": "terrible" },
    { "questionNameId": "improvements", "optionNameId": "security" }
  ]
]
```

Show if: (`rating == "terrible"` AND `attend == "no"`) OR (`rating == "terrible"` AND `improvements includes "security"`)

## Quick Analysis Checklist

When reviewing a survey JSON:

1. **Top level:** Check nameId, title, version are present
2. **Groups:** Verify all groups in `groupsOrder` exist in `groups` map
3. **Questions:** Each group's `questionsIds` must reference valid questions
4. **NameIds:** All nameIds follow format rules and are unique
5. **DependsOn:** Referenced questions and options exist
6. **Question types:** Each question's `type` matches its `value` structure
7. **Choice options:** Simple choice questions have valid options array with nameIds
8. **Toggle/Slider:** Verify required fields (`onLabel`/`offLabel` or `min`/`max`/`step`)
