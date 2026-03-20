package erd

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// DrawioExporter generates Draw.io compatible XML
type DrawioExporter struct{}

// mxGraphModel represents the root element for Draw.io diagrams
type mxGraphModel struct {
	XMLName xml.Name `xml:"mxGraphModel"`
	Dx      int      `xml:"dx,attr"`
	Dy      int      `xml:"dy,attr"`
	Grid    int      `xml:"grid,attr"`
	Root    mxRoot   `xml:"root"`
}

type mxRoot struct {
	XMLName xml.Name `xml:"root"`
	Cells   []mxCell `xml:"mxCell"`
}

type mxCell struct {
	XMLName  xml.Name    `xml:"mxCell"`
	ID       string      `xml:"id,attr"`
	Value    string      `xml:"value,attr,omitempty"`
	Style    string      `xml:"style,attr,omitempty"`
	Vertex   *int        `xml:"vertex,attr,omitempty"`
	Edge     *int        `xml:"edge,attr,omitempty"`
	Parent   string      `xml:"parent,attr,omitempty"`
	Source   string      `xml:"source,attr,omitempty"`
	Target   string      `xml:"target,attr,omitempty"`
	Geometry *mxGeometry `xml:"mxGeometry,omitempty"`
}

type mxGeometry struct {
	XMLName  xml.Name `xml:"mxGeometry"`
	X        float64  `xml:"x,attr,omitempty"`
	Y        float64  `xml:"y,attr,omitempty"`
	Width    float64  `xml:"width,attr,omitempty"`
	Height   float64  `xml:"height,attr,omitempty"`
	Relative *int     `xml:"relative,attr,omitempty"`
	As       string   `xml:"as,attr"`
}

// NewDrawioExporter creates a new Draw.io exporter
func NewDrawioExporter() *DrawioExporter {
	return &DrawioExporter{}
}

// Export generates Draw.io XML for selected tables
func (e *DrawioExporter) Export(schema *SchemaData, selectedTables []string) (string, error) {
	if len(selectedTables) == 0 {
		return "", fmt.Errorf("no tables selected")
	}

	// Create a map for quick lookup
	selectedMap := make(map[string]bool)
	for _, t := range selectedTables {
		selectedMap[t] = true
	}

	// Filter tables and foreign keys
	var tables []Table
	for _, table := range schema.Tables {
		if selectedMap[table.Name] {
			tables = append(tables, table)
		}
	}

	var fks []ForeignKey
	for _, fk := range schema.FKs {
		if selectedMap[fk.FromTable] && selectedMap[fk.ToTable] {
			fks = append(fks, fk)
		}
	}

	// Create Draw.io structure
	model := mxGraphModel{
		Dx:   1426,
		Dy:   827,
		Grid: 1,
		Root: mxRoot{
			Cells: e.generateCells(tables, fks),
		},
	}

	// Marshal to XML
	output, err := xml.MarshalIndent(model, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal XML: %w", err)
	}

	return xml.Header + string(output), nil
}

func (e *DrawioExporter) generateCells(tables []Table, fks []ForeignKey) []mxCell {
	cells := []mxCell{}
	vertex := 1
	edge := 1
	idCounter := 2

	// Add root cells (required by Draw.io)
	cells = append(cells, mxCell{ID: "0"})
	cells = append(cells, mxCell{ID: "1", Parent: "0"})

	// Table ID mapping
	tableIDs := make(map[string]string)

	// Add tables as vertices with proper table structure
	x, y := 40.0, 40.0
	for i, table := range tables {
		tableID := fmt.Sprintf("%d", idCounter)
		tableIDs[table.Name] = tableID
		idCounter++

		// Calculate dimensions based on content
		width := 200.0
		rowHeight := 26.0
		headerHeight := 30.0
		totalHeight := headerHeight + float64(len(table.Columns))*rowHeight

		// Add table container (swimlane)
		cells = append(cells, mxCell{
			ID:     tableID,
			Value:  table.Name,
			Style:  "swimlane;fontStyle=1;childLayout=stackLayout;horizontal=1;startSize=30;horizontalStack=0;resizeParent=1;resizeParentMax=0;resizeLast=0;collapsible=1;marginBottom=0;whiteSpace=wrap;html=1;fillColor=#dae8fc;strokeColor=#6c8ebf;",
			Vertex: &vertex,
			Parent: "1",
			Geometry: &mxGeometry{
				X:      x,
				Y:      y,
				Width:  width,
				Height: totalHeight,
				As:     "geometry",
			},
		})

		// Add each column as a child cell
		for j, col := range table.Columns {
			colID := fmt.Sprintf("%d", idCounter)
			idCounter++

			pkIndicator := ""
			if col.IsPK {
				pkIndicator = "🔑 "
			}

			nullable := ""
			if !col.Nullable {
				nullable = " NOT NULL"
			}

			colValue := fmt.Sprintf("%s%s: %s%s", pkIndicator, col.Name, col.Type, nullable)

			fillColor := "#ffffff"
			if col.IsPK {
				fillColor = "#fff2cc"
			}

			cells = append(cells, mxCell{
				ID:     colID,
				Value:  colValue,
				Style:  fmt.Sprintf("text;strokeColor=none;fillColor=%s;align=left;verticalAlign=middle;spacingLeft=4;spacingRight=4;overflow=hidden;points=[[0,0.5],[1,0.5]];portConstraint=eastwest;rotatable=0;whiteSpace=wrap;html=1;", fillColor),
				Vertex: &vertex,
				Parent: tableID,
				Geometry: &mxGeometry{
					Y:      float64(j) * rowHeight,
					Width:  width,
					Height: rowHeight,
					As:     "geometry",
				},
			})
		}

		// Position tables in a grid
		x += 250
		if (i+1)%3 == 0 {
			x = 40
			y += totalHeight + 50
		}
	}

	// Add foreign keys as edges
	for _, fk := range fks {
		sourceID, sourceExists := tableIDs[fk.FromTable]
		targetID, targetExists := tableIDs[fk.ToTable]

		if !sourceExists || !targetExists {
			continue
		}

		edgeID := fmt.Sprintf("%d", idCounter)
		idCounter++

		relative := 1
		cells = append(cells, mxCell{
			ID:     edgeID,
			Value:  fmt.Sprintf("%s → %s", fk.FromCol, fk.ToCol),
			Style:  "edgeStyle=entityRelationEdgeStyle;rounded=0;orthogonalLoop=1;jettySize=auto;html=1;endArrow=ERmany;startArrow=ERone;endFill=0;startFill=0;",
			Edge:   &edge,
			Parent: "1",
			Source: sourceID,
			Target: targetID,
			Geometry: &mxGeometry{
				Relative: &relative,
				As:       "geometry",
			},
		})
	}

	return cells
}

// EscapeHTML escapes HTML special characters
func EscapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}
