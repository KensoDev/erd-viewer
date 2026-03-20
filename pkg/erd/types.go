package erd

// Column represents a single database column with its metadata
type Column struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Nullable bool   `json:"nullable"`
	IsPK     bool   `json:"isPK"`
}

// Table represents a database table with its columns
type Table struct {
	Name    string   `json:"name"`
	Columns []Column `json:"columns"`
}

// ForeignKey represents a foreign key relationship between tables
type ForeignKey struct {
	FromTable string `json:"fromTable"`
	FromCol   string `json:"fromCol"`
	ToTable   string `json:"toTable"`
	ToCol     string `json:"toCol"`
}

// SchemaData contains all the ERD information for a database schema
type SchemaData struct {
	Title  string       `json:"title"`
	Tables []Table      `json:"tables"`
	FKs    []ForeignKey `json:"fks"`
}
