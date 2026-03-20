package erd

import (
	"fmt"
	"strings"
)

// PlantUMLExporter generates PlantUML ER diagrams
type PlantUMLExporter struct{}

// NewPlantUMLExporter creates a new PlantUML exporter
func NewPlantUMLExporter() *PlantUMLExporter {
	return &PlantUMLExporter{}
}

// Export generates PlantUML syntax for selected tables
func (e *PlantUMLExporter) Export(schema *SchemaData, selectedTables []string) (string, error) {
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

	var sb strings.Builder

	// PlantUML header
	sb.WriteString("@startuml\n")
	sb.WriteString("!theme plain\n")
	sb.WriteString("skinparam linetype ortho\n\n")

	if schema.Title != "" {
		sb.WriteString(fmt.Sprintf("title %s\n\n", schema.Title))
	}

	// Define entities
	for _, table := range tables {
		sb.WriteString(e.buildEntity(table))
		sb.WriteString("\n")
	}

	// Define relationships
	if len(fks) > 0 {
		sb.WriteString("\n' Relationships\n")
		for _, fk := range fks {
			sb.WriteString(e.buildRelationship(fk))
			sb.WriteString("\n")
		}
	}

	sb.WriteString("@enduml\n")

	return sb.String(), nil
}

func (e *PlantUMLExporter) buildEntity(table Table) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("entity \"%s\" as %s {\n", table.Name, SanitizeIdentifier(table.Name)))

	// Group primary keys first
	var primaryKeys []Column
	var regularColumns []Column

	for _, col := range table.Columns {
		if col.IsPK {
			primaryKeys = append(primaryKeys, col)
		} else {
			regularColumns = append(regularColumns, col)
		}
	}

	// Write primary keys
	for _, col := range primaryKeys {
		sb.WriteString(e.buildColumn(col, true))
		sb.WriteString("\n")
	}

	if len(primaryKeys) > 0 && len(regularColumns) > 0 {
		sb.WriteString("  --\n")
	}

	// Write regular columns
	for _, col := range regularColumns {
		sb.WriteString(e.buildColumn(col, false))
		sb.WriteString("\n")
	}

	sb.WriteString("}\n")

	return sb.String()
}

func (e *PlantUMLExporter) buildColumn(col Column, isPK bool) string {
	var parts []string

	// Add PK marker
	if isPK {
		parts = append(parts, "*")
	}

	// Column name
	parts = append(parts, col.Name)

	// Column type
	parts = append(parts, fmt.Sprintf(": %s", col.Type))

	// Nullable marker
	if !col.Nullable && !isPK {
		parts = append(parts, "<<NOT NULL>>")
	}

	return fmt.Sprintf("  %s", strings.Join(parts, " "))
}

func (e *PlantUMLExporter) buildRelationship(fk ForeignKey) string {
	fromTable := SanitizeIdentifier(fk.FromTable)
	toTable := SanitizeIdentifier(fk.ToTable)

	// Using standard ERD notation: many-to-one relationship
	// The "from" table has many records pointing to one record in "to" table
	return fmt.Sprintf("%s }o--|| %s : \"%s → %s\"",
		fromTable,
		toTable,
		fk.FromCol,
		fk.ToCol,
	)
}

// SanitizeIdentifier converts table names to valid PlantUML identifiers
func SanitizeIdentifier(name string) string {
	// Replace invalid characters with underscores
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, ".", "_")
	name = strings.ReplaceAll(name, " ", "_")
	return name
}
