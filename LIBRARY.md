# ERD Viewer Library Documentation

ERD Viewer is now available as a reusable Go library, allowing you to integrate ERD generation and visualization into your own applications.

## Installation

```bash
go get github.com/kensodev/erd-viewer
```

## Package Structure

### `pkg/erd` - Core Library

Core data types and exporters for working with ERD data.

**Types:**
- `SchemaData` - Contains the complete ERD (tables, foreign keys, title)
- `Table` - Represents a database table with columns
- `Column` - Represents a column with type, nullability, and PK info
- `ForeignKey` - Represents a relationship between tables
- `Exporter` - Interface for implementing custom export formats

**Exporters:**
- `DrawioExporter` - Exports to Draw.io XML format
- `PlantUMLExporter` - Exports to PlantUML syntax

### `pkg/erd/postgres` - PostgreSQL Introspection

PostgreSQL-specific database introspection.

**Types:**
- `Introspector` - Introspects PostgreSQL schemas

### `pkg/webview` - Web Server

Reusable web server for ERD visualization.

**Types:**
- `Server` - HTTP server for serving ERD visualization
- `Config` - Configuration for creating a server
- `AssetProvider` - Interface for providing web assets
- `EmbedAssets` - Asset provider for embed.FS

## Usage Examples

### 1. Programmatic Export

Create ERD data and export it to various formats:

```go
package main

import (
    "fmt"
    "github.com/kensodev/erd-viewer/pkg/erd"
)

func main() {
    // Create schema data
    schema := &erd.SchemaData{
        Title: "My Application",
        Tables: []erd.Table{
            {
                Name: "users",
                Columns: []erd.Column{
                    {Name: "id", Type: "integer", IsPK: true, Nullable: false},
                    {Name: "email", Type: "varchar(255)", IsPK: false, Nullable: false},
                    {Name: "name", Type: "varchar(100)", IsPK: false, Nullable: true},
                },
            },
            {
                Name: "posts",
                Columns: []erd.Column{
                    {Name: "id", Type: "integer", IsPK: true, Nullable: false},
                    {Name: "user_id", Type: "integer", IsPK: false, Nullable: false},
                    {Name: "title", Type: "varchar(200)", IsPK: false, Nullable: false},
                },
            },
        },
        FKs: []erd.ForeignKey{
            {
                FromTable: "posts",
                FromCol:   "user_id",
                ToTable:   "users",
                ToCol:     "id",
            },
        },
    }

    // Export to PlantUML
    plantUMLExporter := erd.NewPlantUMLExporter()
    plantUML, err := plantUMLExporter.Export(schema, []string{"users", "posts"})
    if err != nil {
        panic(err)
    }
    fmt.Println(plantUML)

    // Export to Draw.io
    drawioExporter := erd.NewDrawioExporter()
    drawioXML, err := drawioExporter.Export(schema, []string{"users", "posts"})
    if err != nil {
        panic(err)
    }
    fmt.Println(drawioXML)
}
```

### 2. PostgreSQL Introspection

Connect to a PostgreSQL database and introspect its schema:

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/jackc/pgx/v5"
    "github.com/kensodev/erd-viewer/pkg/erd"
    "github.com/kensodev/erd-viewer/pkg/erd/postgres"
)

func main() {
    ctx := context.Background()

    // Connect to database
    dsn := "postgres://user:pass@localhost:5432/mydb"
    conn, err := pgx.Connect(ctx, dsn)
    if err != nil {
        log.Fatalf("Connection failed: %v", err)
    }
    defer conn.Close(ctx)

    // Introspect schema
    introspector := postgres.NewIntrospector(conn)
    schema, err := introspector.IntrospectSchema(ctx, "public", map[string]bool{
        "migrations": true, // Exclude migrations table
    })
    if err != nil {
        log.Fatalf("Introspection failed: %v", err)
    }

    fmt.Printf("Found %d tables and %d foreign keys\n",
        len(schema.Tables), len(schema.FKs))

    // Export the introspected schema
    exporter := erd.NewPlantUMLExporter()
    output, _ := exporter.Export(schema, getAllTableNames(schema))
    fmt.Println(output)
}

func getAllTableNames(schema *erd.SchemaData) []string {
    names := make([]string, len(schema.Tables))
    for i, table := range schema.Tables {
        names[i] = table.Name
    }
    return names
}
```

### 3. Web Server with Default UI

Start a web server with the default ERD Viewer UI:

```go
package main

import (
    "log"

    "github.com/kensodev/erd-viewer/pkg/erd"
    "github.com/kensodev/erd-viewer/pkg/webview"
    "github.com/kensodev/erd-viewer/web"
)

func main() {
    // Create or load schema data
    schema := &erd.SchemaData{
        Title: "My Database",
        // ... tables and FKs
    }

    // Create server with default assets
    server, err := webview.New(webview.Config{
        SchemaData: schema,
        ListenAddr: "127.0.0.1:8080",
        Assets:     &webview.EmbedAssets{FS: web.Files},
    })
    if err != nil {
        log.Fatalf("Failed to create server: %v", err)
    }

    log.Printf("Server running at %s", server.URL())
    log.Fatal(server.Start())
}
```

### 4. Web Server with Custom Assets

Provide your own custom HTML/CSS/JS for the web UI:

```go
package main

import (
    "embed"
    "log"

    "github.com/kensodev/erd-viewer/pkg/erd"
    "github.com/kensodev/erd-viewer/pkg/webview"
)

//go:embed custom-ui/*
var customAssets embed.FS

func main() {
    schema := &erd.SchemaData{
        // ... your schema data
    }

    server, err := webview.New(webview.Config{
        SchemaData: schema,
        ListenAddr: "127.0.0.1:8080",
        Assets:     &webview.EmbedAssets{FS: customAssets},
    })
    if err != nil {
        log.Fatalf("Failed to create server: %v", err)
    }

    log.Fatal(server.Start())
}
```

Your custom assets should have this structure:
```
custom-ui/
├── templates/
│   └── index.html
└── static/
    ├── css/
    │   └── styles.css
    └── js/
        └── app.js
```

## Implementing Custom Exporters

You can create custom exporters for other diagram formats by implementing the `Exporter` interface:

```go
package main

import (
    "fmt"
    "strings"

    "github.com/kensodev/erd-viewer/pkg/erd"
)

// MermaidExporter exports to Mermaid ER diagram syntax
type MermaidExporter struct{}

func (e *MermaidExporter) Export(schema *erd.SchemaData, selectedTables []string) (string, error) {
    if len(selectedTables) == 0 {
        return "", fmt.Errorf("no tables selected")
    }

    selectedMap := make(map[string]bool)
    for _, t := range selectedTables {
        selectedMap[t] = true
    }

    var sb strings.Builder
    sb.WriteString("erDiagram\n")

    // Add entities
    for _, table := range schema.Tables {
        if !selectedMap[table.Name] {
            continue
        }
        sb.WriteString(fmt.Sprintf("    %s {\n", table.Name))
        for _, col := range table.Columns {
            sb.WriteString(fmt.Sprintf("        %s %s\n", col.Type, col.Name))
        }
        sb.WriteString("    }\n")
    }

    // Add relationships
    for _, fk := range schema.FKs {
        if selectedMap[fk.FromTable] && selectedMap[fk.ToTable] {
            sb.WriteString(fmt.Sprintf("    %s ||--o{ %s : \"%s\"\n",
                fk.ToTable, fk.FromTable, fk.FromCol))
        }
    }

    return sb.String(), nil
}

func main() {
    schema := &erd.SchemaData{
        // ... your schema
    }

    exporter := &MermaidExporter{}
    output, err := exporter.Export(schema, []string{"users", "posts"})
    if err != nil {
        panic(err)
    }
    fmt.Println(output)
}
```

## Implementing Custom Database Introspectors

To support other databases like MySQL or SQLite, implement your own introspector:

```go
package main

import (
    "context"
    "database/sql"

    "github.com/kensodev/erd-viewer/pkg/erd"
)

type MySQLIntrospector struct {
    db *sql.DB
}

func NewMySQLIntrospector(db *sql.DB) *MySQLIntrospector {
    return &MySQLIntrospector{db: db}
}

func (i *MySQLIntrospector) IntrospectSchema(ctx context.Context, schema string, exclude map[string]bool) (*erd.SchemaData, error) {
    // Query MySQL information_schema to get:
    // 1. Tables and columns
    // 2. Primary keys
    // 3. Foreign keys

    // Build and return erd.SchemaData
    return &erd.SchemaData{
        Tables: []erd.Table{
            // ... tables from MySQL
        },
        FKs: []erd.ForeignKey{
            // ... foreign keys from MySQL
        },
    }, nil
}
```

## Custom Asset Providers

Implement your own asset provider for dynamic asset loading:

```go
package main

import (
    "fmt"
    "io"
    "net/http"
)

// HTTPAssetProvider loads assets from an HTTP server
type HTTPAssetProvider struct {
    baseURL string
    client  *http.Client
}

func (p *HTTPAssetProvider) ReadFile(name string) ([]byte, error) {
    resp, err := p.client.Get(p.baseURL + "/" + name)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
    }

    return io.ReadAll(resp.Body)
}
```

## API Reference

See the [GoDoc](https://pkg.go.dev/github.com/kensodev/erd-viewer) for complete API documentation.

## Examples

Complete working examples are available in the [examples](./examples/custom-webview) directory.
