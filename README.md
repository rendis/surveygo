# Survey Processing Library Documentation

## Table of Contents

1. [Overview](#overview)
    - [Definitions and Rules within a Survey](#definitions-and-rules-within-a-survey)
2. [Base Structures](#base-structures)
    - [Survey](#survey)
    - [Question](#question)
    - [Group](#group)
3. [Question Structures](#question-structures)
    - [Types of Questions](#types-of-questions)
    - [Choice](#choice)
    - [Text](#text)
    - [External Question](#external-question)
4. [Functions](#functions)
    - [operation.go](#operationgo)
    - [operation_de_serializers.go](#operation_de_serializersgo)
    - [operation_group.go](#operation_groupgo)
    - [operation_question.go](#operation_questiongo)

## Overview

`surveygo` facilitates the creation and management of surveys. It provides data structures and methods for creating
surveys, questions, and groups of questions, as well as handling validations.

### Definitions and Rules within a survey

- All `survey` structures will have a **unique** identifier called `nameId`.
- Every reference to a `nameId` must be unique, or in other words:
    - A question can only be associated with one group.
    - A group can only be associated with one question or the initial set of groups `groupsOrder`.

## Base Structures

### Survey

Structure representing a complete survey.

**Fields**

- `title`: Survey title. (Required)
- `version`: Survey version. (Required)
- `description`: Survey description. (Optional)
- `questions`: Map of questions. (Required)
- `groups`: Map of groups. (Required)
- `groupsOrder`: Order of the groups. (Required)

### Question

Structure representing a question within a survey.

**Fields**

- `nameId`: Question identifier. (Required)
- `visible`: Indicates if the question is visible. (Required)
- `type`: Type of the question. (Required)
- `label`: Question label. (Required)
- `required`: Indicates if the question is mandatory to answer. (Required)
- `Value`: Object representing the value of the question. Varies depending on the type of question. (Required)

### Group

Structure representing a group of questions in a survey.

**Fields**

- `nameId`: Group identifier. (Required)
- `title`: Group title. (Optional)
- `description`: Group description. (Optional)
- `visible`: Indicates if the group is visible. (Required)
- `isExternalSurvey`: Indicates if the group is an external survey. (Optional)
- `questionsIds`: Identifiers of the questions that belong to the group. (Required)
  <br>If the group is an external survey, this field will indicate the identifier of the external survey.

## Question Structures

### Types of Questions

#### Choice:

- `single_select`: Single select
- `multi_select`: Multiple select
- `radio`: Single select
- `checkbox`: Multiple select

#### Text:

- `email`: Email
- `telephone`: Telephone
- `text_area`: Free text
- `input_text`: Free text
- `information`: Information field, not editable

#### External Questions:

- `external_question`: External question

### Choice

Structure for all questions in the `Choice` group.

- `placeholder`: Placeholder text for the question. (Optional)
- `defaults`: List of default values for the question. Each value must be a valid `nameId` of an option. (Optional)
- `options`: Question options. (Required)
    - `nameId`: Option identifier. (Required)
    - `label`: Option label. (Required)
    - `groupsIds`: Identifiers of the groups to be displayed when the option is selected. (Optional)

### Text

**Email** (`email`)

- `placeholder`: Placeholder text for the question. (Optional)
- `allowedDomains`: Allowed domains for the email. (Optional)

**Telephone** (`telephone`)

- `placeholder`: Placeholder text for the question. (Optional)
- `allowedCountryCodes`: List of allowed country codes. (Optional)

**FreeText** (`input_text` and `text_area`)

- `placeholder`: Placeholder text for the question. (Optional)
- `min`: Minimum length of the text. (Optional)
- `max`: Maximum length of the text. (Optional)

**Information** (`information`)

- `text`: Text to be displayed. (Required)

### External Question

Used to create an external questions.

**External Question** (`external_question`)

- `placeholder`: Placeholder for the question. (Optional)
- `defaults`: List of default values for the question. (Optional)
- `questionType`: Type of the question. Refer to [Types of Questions](#types-of-questions). (Required)
- `externalType`: Type of the external question. (Required)
- `description`: Description of the external question. (Optional)
- `src`: Source of the external question. (Optional)

## Functions

For the complete list of available functions and methods, please refer to the files:

- [operation.go](operation.go): Basic survey operations (construction, validation, answers review, etc.).
- [operation_de_serializers.go](operation_de_serializers.go): Survey serialization and deserialization.
- [operation_group.go](operation_group.go): Operations on groups (add, remove, etc.).
- [operation_question.go](operation_question.go): Operations on questions (add, remove, etc.).