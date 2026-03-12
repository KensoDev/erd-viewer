package db

import (
	"context"
	"sort"

	"github.com/jackc/pgx/v5"
)

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

// Introspector handles database schema introspection
type Introspector struct {
	conn *pgx.Conn
}

// NewIntrospector creates a new database introspector
func NewIntrospector(conn *pgx.Conn) *Introspector {
	return &Introspector{conn: conn}
}

// IntrospectSchema fetches the complete schema including tables and foreign keys
func (i *Introspector) IntrospectSchema(ctx context.Context, schema string, exclude map[string]bool) (*SchemaData, error) {
	colMap, err := i.fetchColumns(ctx, schema, exclude)
	if err != nil {
		return nil, err
	}

	fks, err := i.fetchForeignKeys(ctx, schema, exclude)
	if err != nil {
		return nil, err
	}

	// Build sorted table list
	names := make([]string, 0, len(colMap))
	for n := range colMap {
		names = append(names, n)
	}
	sort.Strings(names)

	tables := make([]Table, 0, len(names))
	for _, n := range names {
		tables = append(tables, Table{Name: n, Columns: colMap[n]})
	}

	return &SchemaData{
		Tables: tables,
		FKs:    fks,
	}, nil
}

func (i *Introspector) fetchColumns(ctx context.Context, schema string, exclude map[string]bool) (map[string][]Column, error) {
	rows, err := i.conn.Query(ctx, `
		SELECT c.table_name, c.column_name, c.data_type, c.is_nullable,
		       COALESCE((
		           SELECT 'PK'
		           FROM information_schema.table_constraints tc
		           JOIN information_schema.key_column_usage kcu
		             ON tc.constraint_name = kcu.constraint_name
		            AND tc.table_schema    = kcu.table_schema
		           WHERE tc.constraint_type = 'PRIMARY KEY'
		             AND tc.table_name      = c.table_name
		             AND tc.table_schema    = c.table_schema
		             AND kcu.column_name    = c.column_name
		           LIMIT 1
		       ), '') AS key_type
		FROM information_schema.columns c
		WHERE c.table_schema = $1
		ORDER BY c.table_name, c.ordinal_position
	`, schema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tables := make(map[string][]Column)
	for rows.Next() {
		var tbl, col, typ, nullable, key string
		if err := rows.Scan(&tbl, &col, &typ, &nullable, &key); err != nil {
			return nil, err
		}
		if exclude[tbl] {
			continue
		}
		tables[tbl] = append(tables[tbl], Column{
			Name:     col,
			Type:     typ,
			Nullable: nullable == "YES",
			IsPK:     key == "PK",
		})
	}
	return tables, rows.Err()
}

func (i *Introspector) fetchForeignKeys(ctx context.Context, schema string, exclude map[string]bool) ([]ForeignKey, error) {
	rows, err := i.conn.Query(ctx, `
		SELECT src_tbl.relname, src_col.attname, tgt_tbl.relname, tgt_col.attname
		FROM pg_constraint con
		JOIN pg_class     src_tbl ON src_tbl.oid = con.conrelid
		JOIN pg_class     tgt_tbl ON tgt_tbl.oid = con.confrelid
		JOIN pg_namespace ns      ON ns.oid       = src_tbl.relnamespace
		JOIN pg_attribute src_col ON src_col.attrelid = con.conrelid
		                          AND src_col.attnum   = ANY(con.conkey)
		JOIN pg_attribute tgt_col ON tgt_col.attrelid = con.confrelid
		                          AND tgt_col.attnum   = ANY(con.confkey)
		WHERE con.contype = 'f' AND ns.nspname = $1
		ORDER BY src_tbl.relname, src_col.attname
	`, schema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fks []ForeignKey
	for rows.Next() {
		var fk ForeignKey
		if err := rows.Scan(&fk.FromTable, &fk.FromCol, &fk.ToTable, &fk.ToCol); err != nil {
			return nil, err
		}
		if exclude[fk.FromTable] || exclude[fk.ToTable] {
			continue
		}
		fks = append(fks, fk)
	}
	return fks, rows.Err()
}
