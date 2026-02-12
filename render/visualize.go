package render

import (
	"bytes"
	"fmt"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

const legendHTML = `
<div style="position:absolute;top:12px;right:20px;background:rgba(255,255,255,0.95);border:1px solid #ccc;border-radius:6px;padding:10px 16px;font-family:sans-serif;font-size:13px;box-shadow:0 2px 6px rgba(0,0,0,0.1)">
  <div style="font-weight:600;margin-bottom:8px;font-size:14px">Leyenda</div>
  <div style="display:flex;align-items:center;margin-bottom:5px">
    <span style="display:inline-block;width:16px;height:16px;border-radius:3px;background:#4CAF50;border:2px solid #388E3C;margin-right:8px"></span>
    Grupo
  </div>
  <div style="display:flex;align-items:center">
    <span style="display:inline-block;width:16px;height:16px;border-radius:3px;background:#2196F3;border:2px solid #1565C0;margin-right:8px"></span>
    Grupo con repetición (↻)
  </div>
</div>
`

func renderTreeToBytes(tree *GroupTree) ([]byte, error) {
	var roots []opts.TreeData
	for _, root := range tree.Roots {
		roots = append(roots, *toEchartsTree(root))
	}

	tc := charts.NewTree()
	tc.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Survey Group Hierarchy"}),
		charts.WithInitializationOpts(opts.Initialization{Width: "1400px", Height: "900px"}),
		charts.WithTooltipOpts(opts.Tooltip{Show: opts.Bool(true)}),
		charts.WithLegendOpts(opts.Legend{Show: opts.Bool(false)}),
	)

	tc.AddSeries("groups", roots,
		charts.WithTreeOpts(opts.TreeChart{
			Layout:            "orthogonal",
			Orient:            "LR",
			InitialTreeDepth:  -1,
			ExpandAndCollapse: opts.Bool(true),
			Roam:              opts.Bool(true),
			Left:              "80px",
			Right:             "200px",
			Top:               "50px",
			Bottom:            "50px",
		}),
		charts.WithLabelOpts(opts.Label{
			Show:     opts.Bool(true),
			Position: "right",
			FontSize: 14,
		}),
	)

	var buf bytes.Buffer
	if err := tc.Render(&buf); err != nil {
		return nil, fmt.Errorf("rendering tree chart: %w", err)
	}

	buf.WriteString(legendHTML)
	return buf.Bytes(), nil
}

func toEchartsTree(node *GroupNode) *opts.TreeData {
	name := node.NameId

	td := &opts.TreeData{
		Name:       name,
		Symbol:     "roundRect",
		SymbolSize: []int{20, 20},
		ItemStyle:  &opts.ItemStyle{Color: "#4CAF50", BorderColor: "#388E3C"},
	}

	if node.AllowRepeat {
		td.Name = name + " \u21bb"
		td.SymbolSize = []int{24, 24}
		td.ItemStyle = &opts.ItemStyle{Color: "#2196F3", BorderColor: "#1565C0"}
	}

	for _, child := range node.Children {
		td.Children = append(td.Children, toEchartsTree(child))
	}
	return td
}
