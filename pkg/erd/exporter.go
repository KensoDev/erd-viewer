package erd

// Exporter defines the interface for exporting schema data to various formats
type Exporter interface {
	// Export generates output for the selected tables in the format-specific syntax
	Export(schema *SchemaData, selectedTables []string) (string, error)
}
