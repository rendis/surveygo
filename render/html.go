package render

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

const defaultCSSStr = `/* Survey Card — Default Styles */

.card-body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  font-size: 14px;
  color: #1a1a1a;
  line-height: 1.5;
  max-width: 900px;
  margin: 0 auto;
  padding: 24px;
}

.card-title {
  font-size: 20px;
  font-weight: 700;
  margin-bottom: 24px;
  padding-bottom: 12px;
  border-bottom: 2px solid #e5e7eb;
}

/* --- Sections --- */

.card-group,
.card-repeat-table,
.card-repeat-list {
  margin-bottom: 20px;
  padding: 16px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  background: #fff;
}

.card-section-title {
  font-size: 15px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #f3f4f6;
}

/* --- Fields (group sections) --- */

.card-fields {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.card-field {
  display: flex;
  align-items: baseline;
  gap: 8px;
  padding: 4px 0;
}

.card-field-label {
  font-weight: 500;
  color: #6b7280;
  min-width: 180px;
  flex-shrink: 0;
}

.card-field-label:not(:empty)::after {
  content: ":";
}

.card-field-value {
  color: #1a1a1a;
}

/* --- Multi-select options --- */

.card-options {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.card-option {
  display: inline-block;
  padding: 2px 10px;
  border-radius: 12px;
  font-size: 13px;
}

.card-option--selected {
  background: #dcfce7;
  color: #166534;
  border: 1px solid #bbf7d0;
}

.card-option--not-selected {
  background: #f3f4f6;
  color: #9ca3af;
  border: 1px solid #e5e7eb;
}

/* --- Repeat-table --- */

.card-table {
  width: 100%;
  border-collapse: collapse;
}

.card-table-header {
  display: flex;
  background: #f9fafb;
  border-bottom: 2px solid #e5e7eb;
}

.card-table-header .card-col {
  flex: 1;
  padding: 8px 12px;
  font-weight: 600;
  font-size: 13px;
  color: #374151;
}

.card-row {
  display: flex;
  border-bottom: 1px solid #f3f4f6;
}

.card-row:last-child {
  border-bottom: none;
}

.card-cell {
  flex: 1;
  padding: 8px 12px;
  font-size: 14px;
}

/* --- Repeat-list instances --- */

.card-instance {
  margin: 12px 0;
  padding: 12px;
  border: 1px dashed #d1d5db;
  border-radius: 6px;
  background: #fafafa;
}

.card-instance .card-group,
.card-instance .card-repeat-table,
.card-instance .card-repeat-list {
  border: none;
  padding: 8px 0;
  margin-bottom: 8px;
  background: transparent;
}
`

func defaultCSS() []byte {
	return []byte(defaultCSSStr)
}

const defaultTemplatesStr = `
{{- define "card" -}}
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>{{.Title}}</title>
  <link rel="stylesheet" href="card.css">
</head>
<body class="card-body">
  <div class="card-title">{{.Title}}</div>
  {{- range .Sections}}
  {{template "section" .}}
  {{- end}}
</body>
</html>
{{- end -}}

{{- define "section" -}}
{{- if eq .Type "group"}}{{template "section-group" .}}
{{- else if eq .Type "repeat-table"}}{{template "section-repeat-table" .}}
{{- else if eq .Type "repeat-list"}}{{template "section-repeat-list" .}}
{{- end -}}
{{- end -}}

{{- define "section-group" -}}
<div class="card-group card-group--{{.NameId}}">
  <div class="card-section-title">{{.Title}}</div>
  {{- if .Fields}}
  <div class="card-fields">
    {{- range .Fields}}
    {{template "field" .}}
    {{- end}}
  </div>
  {{- end}}
  {{- range .Sections}}
  {{template "section" .}}
  {{- end}}
</div>
{{- end -}}

{{- define "section-repeat-table" -}}
<div class="card-repeat-table card-group--{{.NameId}}">
  <div class="card-section-title">{{.Title}}</div>
  <div class="card-table">
    <div class="card-table-header">
      {{- range .Columns}}
      <span class="card-col card-col--{{.NameId}}">{{.Label}}</span>
      {{- end}}
    </div>
    {{- range .Rows}}
    {{- $row := .}}
    <div class="card-row">
      {{- range $.Columns}}
      <span class="card-cell card-cell--{{.NameId}}">{{renderCell $row .}}</span>
      {{- end}}
    </div>
    {{- end}}
  </div>
</div>
{{- end -}}

{{- define "section-repeat-list" -}}
<div class="card-repeat-list card-group--{{.NameId}}">
  <div class="card-section-title">{{.Title}}</div>
  {{- range .Instances}}
  <div class="card-instance">
    {{- range .Sections}}
    {{template "section" .}}
    {{- end}}
  </div>
  {{- end}}
</div>
{{- end -}}

{{- define "field" -}}
<div class="card-field card-field--{{.Type}}">
  <span class="card-field-label">{{.Label}}</span>
  {{- if eq .Type "select"}}
  <span class="card-field-value">{{selectLabel .Value}}</span>
  {{- else if eq .Type "multi-select"}}
  <div class="card-field-value">
    <div class="card-options">
      {{- range optionRefs .Value}}
      <span class="card-option {{optionClass .Selected}}">{{.Label}}</span>
      {{- end}}
    </div>
  </div>
  {{- else if eq .Type "toggle"}}
  <span class="card-field-value">{{if isToggleOn .Value}}Yes{{else}}No{{end}}</span>
  {{- else}}
  <span class="card-field-value">{{textValue .Value}}</span>
  {{- end}}
</div>
{{- end -}}
`

var defaultTmpl *template.Template

func init() {
	defaultTmpl = template.Must(
		template.New("").Funcs(templateFuncMap()).Parse(defaultTemplatesStr),
	)
}

func templateFuncMap() template.FuncMap {
	return template.FuncMap{
		"selectLabel": selectLabel,
		"optionRefs":  optionRefsFn,
		"isToggleOn":  isToggleOn,
		"textValue":   textValue,
		"optionClass": optionClass,
		"renderCell":  renderCell,
	}
}

func generateHTML(card *SurveyCard) ([]byte, error) {
	var buf bytes.Buffer
	if err := defaultTmpl.ExecuteTemplate(&buf, "card", card); err != nil {
		return nil, fmt.Errorf("executing card template: %w", err)
	}
	return buf.Bytes(), nil
}

func renderSectionHTML(sec Section) string {
	var buf bytes.Buffer
	_ = defaultTmpl.ExecuteTemplate(&buf, "section", sec)
	return buf.String()
}

func renderFieldHTML(f Field) string {
	var buf bytes.Buffer
	_ = defaultTmpl.ExecuteTemplate(&buf, "field", f)
	return buf.String()
}

// --- Template FuncMap helpers ---

func selectLabel(v any) string {
	if ref, ok := v.(OptionRef); ok {
		return ref.Label
	}
	return ""
}

func optionRefsFn(v any) []OptionRef {
	if refs, ok := v.([]OptionRef); ok {
		return refs
	}
	return nil
}

func isToggleOn(v any) bool {
	b, _ := v.(bool)
	return b
}

func textValue(v any) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

func optionClass(selected bool) string {
	if selected {
		return "card-option--selected"
	}
	return "card-option--not-selected"
}

func renderCell(row Row, col Column) template.HTML {
	val := row[col.NameId]

	switch col.FieldType {
	case "select":
		if ref, ok := val.(OptionRef); ok {
			return template.HTML(template.HTMLEscapeString(ref.Label))
		}
		return ""

	case "multi-select":
		refs, ok := val.([]OptionRef)
		if !ok {
			return ""
		}
		var selected []string
		for _, ref := range refs {
			if ref.Selected {
				selected = append(selected, ref.Label)
			}
		}
		if len(selected) <= 3 {
			var escaped []string
			for _, s := range selected {
				escaped = append(escaped, template.HTMLEscapeString(s))
			}
			return template.HTML(strings.Join(escaped, ", "))
		}
		var buf bytes.Buffer
		buf.WriteString(`<div class="card-options">`)
		for _, ref := range refs {
			if ref.Selected {
				fmt.Fprintf(&buf, `<span class="card-option card-option--selected">%s</span>`, template.HTMLEscapeString(ref.Label))
			}
		}
		buf.WriteString(`</div>`)
		return template.HTML(buf.String())

	case "toggle":
		if b, ok := val.(bool); ok && b {
			return "Yes"
		}
		return "No"

	default:
		if val == nil {
			return ""
		}
		return template.HTML(template.HTMLEscapeString(fmt.Sprintf("%v", val)))
	}
}
