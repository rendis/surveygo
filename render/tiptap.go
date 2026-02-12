package render

import (
	"fmt"
	"strings"
)

func buildTipTapDoc(card *SurveyCard) TipTapNode {
	doc := TipTapNode{Type: "doc"}

	// Survey title as h1.
	doc.Content = append(doc.Content, heading(1, card.Title))

	// Convert each top-level section at depth 0 (-> h2).
	doc.Content = append(doc.Content, sectionsToNodes(card.Sections, 0)...)

	return doc
}

// sectionsToNodes converts sections into TipTap nodes.
// depth 0 -> h2, depth 1 -> h3, depth 2+ -> bold paragraph.
func sectionsToNodes(sections []Section, depth int) []TipTapNode {
	var nodes []TipTapNode
	for _, sec := range sections {
		switch sec.Type {
		case "group":
			nodes = append(nodes, groupToNodes(sec, depth)...)
		case "repeat-table":
			nodes = append(nodes, repeatTableToNodes(sec, depth)...)
		case "repeat-list":
			nodes = append(nodes, repeatListToNodes(sec, depth)...)
		}
	}
	return nodes
}

func groupToNodes(sec Section, depth int) []TipTapNode {
	nodes := []TipTapNode{sectionTitle(depth, sec.Title)}

	if len(sec.Fields) > 0 {
		var items []TipTapNode
		for _, f := range sec.Fields {
			val := fieldValueToText(f)
			items = append(items, listItem(
				ttParagraph(boldText(f.Label+": "), textNode(val)),
			))
		}
		nodes = append(nodes, bulletList(items))
	}

	if len(sec.Sections) > 0 {
		nodes = append(nodes, sectionsToNodes(sec.Sections, depth+1)...)
	}

	return nodes
}

func repeatTableToNodes(sec Section, depth int) []TipTapNode {
	nodes := []TipTapNode{sectionTitle(depth, sec.Title)}

	if len(sec.Columns) == 0 {
		return nodes
	}

	var tableRows []TipTapNode

	// Header row.
	var headerCells []TipTapNode
	for _, col := range sec.Columns {
		headerCells = append(headerCells, tableHeader(col.Label))
	}
	tableRows = append(tableRows, tableRow(headerCells))

	// Data rows.
	for _, row := range sec.Rows {
		var cells []TipTapNode
		for _, col := range sec.Columns {
			cells = append(cells, tableCell(cellValueToText(row[col.NameId], col)))
		}
		tableRows = append(tableRows, tableRow(cells))
	}

	nodes = append(nodes, TipTapNode{
		Type:    "table",
		Content: tableRows,
	})

	return nodes
}

func repeatListToNodes(sec Section, depth int) []TipTapNode {
	nodes := []TipTapNode{sectionTitle(depth, sec.Title)}

	single := len(sec.Instances) == 1
	for i, inst := range sec.Instances {
		if i > 0 {
			nodes = append(nodes, TipTapNode{Type: "horizontalRule"})
		}
		if single {
			nodes = append(nodes, sectionsToNodes(inst.Sections, depth+1)...)
		} else {
			nodes = append(nodes, sectionTitle(depth+1, fmt.Sprintf("%s #%d", sec.Title, i+1)))
			nodes = append(nodes, sectionsToNodes(inst.Sections, depth+2)...)
		}
	}

	return nodes
}

func sectionTitle(depth int, text string) TipTapNode {
	if depth <= 1 {
		return heading(depth+2, text) // depth 0 -> h2, depth 1 -> h3
	}
	return ttParagraph(boldText(text))
}

func fieldValueToText(f Field) string {
	if f.Value == nil {
		return "\u2014"
	}

	switch f.Type {
	case "toggle":
		if b, ok := f.Value.(bool); ok {
			if b {
				return "Si"
			}
			return "No"
		}
		return "\u2014"

	case "select":
		if m, ok := f.Value.(OptionRef); ok {
			return m.Label
		}
		if m, ok := f.Value.(map[string]any); ok {
			if label, ok := m["label"].(string); ok {
				return label
			}
		}
		return "\u2014"

	case "multi-select":
		if refs, ok := f.Value.([]OptionRef); ok {
			var selected []string
			for _, r := range refs {
				if r.Selected {
					selected = append(selected, r.Label)
				}
			}
			if len(selected) == 0 {
				return "\u2014"
			}
			return strings.Join(selected, ", ")
		}
		if arr, ok := f.Value.([]any); ok {
			var selected []string
			for _, item := range arr {
				if m, ok := item.(map[string]any); ok {
					if sel, _ := m["selected"].(bool); sel {
						if label, ok := m["label"].(string); ok {
							selected = append(selected, label)
						}
					}
				}
			}
			if len(selected) == 0 {
				return "\u2014"
			}
			return strings.Join(selected, ", ")
		}
		return "\u2014"

	default:
		if s, ok := f.Value.(string); ok {
			if s == "" {
				return "\u2014"
			}
			return s
		}
		return fmt.Sprintf("%v", f.Value)
	}
}

func cellValueToText(val any, col Column) string {
	if val == nil {
		return "\u2014"
	}

	switch col.FieldType {
	case "toggle":
		if b, ok := val.(bool); ok {
			if b {
				return "Si"
			}
			return "No"
		}
		return "\u2014"

	case "select":
		if m, ok := val.(OptionRef); ok {
			return m.Label
		}
		if m, ok := val.(map[string]any); ok {
			if label, ok := m["label"].(string); ok {
				return label
			}
		}
		return "\u2014"

	case "multi-select":
		if refs, ok := val.([]OptionRef); ok {
			var selected []string
			for _, r := range refs {
				if r.Selected {
					selected = append(selected, r.Label)
				}
			}
			if len(selected) == 0 {
				return "\u2014"
			}
			return strings.Join(selected, ", ")
		}
		if arr, ok := val.([]any); ok {
			var selected []string
			for _, item := range arr {
				if m, ok := item.(map[string]any); ok {
					if sel, _ := m["selected"].(bool); sel {
						if label, ok := m["label"].(string); ok {
							selected = append(selected, label)
						}
					}
				}
			}
			if len(selected) == 0 {
				return "\u2014"
			}
			return strings.Join(selected, ", ")
		}
		return "\u2014"

	default:
		if s, ok := val.(string); ok {
			if s == "" {
				return "\u2014"
			}
			return s
		}
		return fmt.Sprintf("%v", val)
	}
}

// --- TipTap node helpers ---

func heading(level int, text string) TipTapNode {
	return TipTapNode{
		Type:    "heading",
		Attrs:   map[string]any{"level": level},
		Content: []TipTapNode{textNode(text)},
	}
}

// ttParagraph creates a TipTap paragraph node (named to avoid conflict with html/template).
func ttParagraph(children ...TipTapNode) TipTapNode {
	return TipTapNode{
		Type:    "paragraph",
		Content: children,
	}
}

func textNode(s string) TipTapNode {
	return TipTapNode{
		Type: "text",
		Text: s,
	}
}

func boldText(s string) TipTapNode {
	return TipTapNode{
		Type:  "text",
		Text:  s,
		Marks: []TipTapMark{{Type: "bold"}},
	}
}

func bulletList(items []TipTapNode) TipTapNode {
	return TipTapNode{
		Type:    "bulletList",
		Content: items,
	}
}

func listItem(children ...TipTapNode) TipTapNode {
	return TipTapNode{
		Type:    "listItem",
		Content: children,
	}
}

func tableRow(cells []TipTapNode) TipTapNode {
	return TipTapNode{
		Type:    "tableRow",
		Content: cells,
	}
}

func tableHeader(text string) TipTapNode {
	return TipTapNode{
		Type: "tableHeader",
		Content: []TipTapNode{
			ttParagraph(boldText(text)),
		},
	}
}

func tableCell(text string) TipTapNode {
	return TipTapNode{
		Type: "tableCell",
		Content: []TipTapNode{
			ttParagraph(textNode(text)),
		},
	}
}
