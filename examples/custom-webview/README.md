# Custom Web View Example

This example demonstrates how to use the ERD Viewer library programmatically without the CLI.

## What This Example Shows

1. **Programmatic Export**: Create schema data and export it to PlantUML or Draw.io format
2. **Database Introspection**: Connect to a PostgreSQL database and introspect its schema
3. **Custom Web View**: How to provide your own HTML/CSS/JS assets for a custom UI

## Running the Example

### Example 1 & 2: Programmatic usage

```bash
# Run with sample data only
go run main.go

# Run with actual database introspection
DATABASE_URL="postgres://user:pass@localhost:5432/dbname" go run main.go
```

### Example 3: Custom Web View

To use a custom web view:

1. Create your custom web assets directory with the following structure:
   ```
   custom-web-assets/
   ├── templates/
   │   └── index.html
   └── static/
       ├── css/
       │   └── styles.css
       └── js/
           └── app.js
   ```

2. Uncomment the "Example 3" section in main.go

3. Run the example:
   ```bash
   go run main.go
   ```

## Using the Library in Your Own Project

### Installation

```bash
go get github.com/kensodev/erd-viewer
```

### Basic Usage

```go
package main

import (
    "github.com/kensodev/erd-viewer/pkg/erd"
)

func main() {
    // Create schema data
    schema := &erd.SchemaData{
        Title: "My Database",
        Tables: []erd.Table{
            // ... your tables
        },
        FKs: []erd.ForeignKey{
            // ... your foreign keys
        },
    }

    // Export to PlantUML
    exporter := erd.NewPlantUMLExporter()
    output, err := exporter.Export(schema, []string{"table1", "table2"})
    if err != nil {
        panic(err)
    }
    println(output)
}
```

### PostgreSQL Introspection

```go
import (
    "context"
    "github.com/jackc/pgx/v5"
    "github.com/kensodev/erd-viewer/pkg/erd/postgres"
)

func main() {
    ctx := context.Background()
    conn, _ := pgx.Connect(ctx, "postgres://...")
    defer conn.Close(ctx)

    introspector := postgres.NewIntrospector(conn)
    schema, _ := introspector.IntrospectSchema(ctx, "public", map[string]bool{})

    // Use schema data...
}
```

### Custom Web Server

```go
import (
    "github.com/kensodev/erd-viewer/pkg/webview"
    "github.com/kensodev/erd-viewer/web" // Default assets
)

func main() {
    server, _ := webview.New(webview.Config{
        SchemaData: schema,
        ListenAddr: "127.0.0.1:8080",
        Assets:     &webview.EmbedAssets{FS: web.Files}, // or your own
    })

    server.Start()
}
```

## Implementing Your Own Database Introspector

To support other databases (MySQL, SQLite, etc.), implement your own introspector:

```go
type MySQLIntrospector struct {
    conn *sql.DB
}

func (i *MySQLIntrospector) IntrospectSchema(ctx context.Context, schema string, exclude map[string]bool) (*erd.SchemaData, error) {
    // Query MySQL information_schema
    // Build erd.SchemaData
    return &erd.SchemaData{...}, nil
}
```

## Implementing Custom Exporters

To add support for other diagram formats:

```go
type MermaidExporter struct{}

func (e *MermaidExporter) Export(schema *erd.SchemaData, selectedTables []string) (string, error) {
    // Generate Mermaid syntax
    return "erDiagram\n...", nil
}
```
