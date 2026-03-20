package erd

import (
	"strings"
	"testing"
)

func getTestSchema() *SchemaData {
	return &SchemaData{
		Title: "Test Schema",
		Tables: []Table{
			{
				Name: "users",
				Columns: []Column{
					{Name: "id", Type: "integer", Nullable: false, IsPK: true},
					{Name: "email", Type: "varchar", Nullable: false, IsPK: false},
					{Name: "name", Type: "varchar", Nullable: true, IsPK: false},
				},
			},
			{
				Name: "posts",
				Columns: []Column{
					{Name: "id", Type: "integer", Nullable: false, IsPK: true},
					{Name: "user_id", Type: "integer", Nullable: false, IsPK: false},
					{Name: "title", Type: "varchar", Nullable: false, IsPK: false},
					{Name: "content", Type: "text", Nullable: true, IsPK: false},
				},
			},
		},
		FKs: []ForeignKey{
			{
				FromTable: "posts",
				FromCol:   "user_id",
				ToTable:   "users",
				ToCol:     "id",
			},
		},
	}
}

func TestDrawioExporter_Export(t *testing.T) {
	schema := getTestSchema()
	exporter := NewDrawioExporter()

	tests := []struct {
		name           string
		selectedTables []string
		wantErr        bool
		checkContains  []string
	}{
		{
			name:           "export all tables",
			selectedTables: []string{"users", "posts"},
			wantErr:        false,
			checkContains: []string{
				"<?xml",
				"mxGraphModel",
				"value=\"users\"",
				"value=\"posts\"",
				"🔑 id: integer NOT NULL",
				"email: varchar NOT NULL",
				"user_id → id",
			},
		},
		{
			name:           "export single table",
			selectedTables: []string{"users"},
			wantErr:        false,
			checkContains: []string{
				"<?xml",
				"value=\"users\"",
				"email: varchar NOT NULL",
			},
		},
		{
			name:           "no tables selected",
			selectedTables: []string{},
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := exporter.Export(schema, tt.selectedTables)

			if (err != nil) != tt.wantErr {
				t.Errorf("Export() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				for _, check := range tt.checkContains {
					if !strings.Contains(output, check) {
						t.Errorf("Export() output does not contain %q", check)
					}
				}
			}
		})
	}
}

func TestPlantUMLExporter_Export(t *testing.T) {
	schema := getTestSchema()
	exporter := NewPlantUMLExporter()

	tests := []struct {
		name           string
		selectedTables []string
		wantErr        bool
		checkContains  []string
	}{
		{
			name:           "export all tables",
			selectedTables: []string{"users", "posts"},
			wantErr:        false,
			checkContains: []string{
				"@startuml",
				"@enduml",
				"entity \"users\"",
				"entity \"posts\"",
				"* id : integer",
				"email : varchar",
				"user_id : integer",
				"posts }o--|| users",
			},
		},
		{
			name:           "export single table",
			selectedTables: []string{"users"},
			wantErr:        false,
			checkContains: []string{
				"@startuml",
				"entity \"users\"",
				"* id : integer",
			},
		},
		{
			name:           "no tables selected",
			selectedTables: []string{},
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := exporter.Export(schema, tt.selectedTables)

			if (err != nil) != tt.wantErr {
				t.Errorf("Export() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				for _, check := range tt.checkContains {
					if !strings.Contains(output, check) {
						t.Errorf("Export() output does not contain %q", check)
					}
				}
			}
		})
	}
}

func TestPlantUMLExporter_SanitizeIdentifier(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple_table", "simple_table"},
		{"table-with-dash", "table_with_dash"},
		{"table.with.dot", "table_with_dot"},
		{"table with space", "table_with_space"},
		{"mixed-table.name", "mixed_table_name"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := SanitizeIdentifier(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeIdentifier(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDrawioExporter_EscapeHTML(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple", "simple"},
		{"<tag>", "&lt;tag&gt;"},
		{"a & b", "a &amp; b"},
		{"quote\"test", "quote&quot;test"},
		{"apostrophe'test", "apostrophe&#39;test"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := EscapeHTML(tt.input)
			if result != tt.expected {
				t.Errorf("EscapeHTML(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
