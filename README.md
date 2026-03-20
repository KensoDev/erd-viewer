<div align="center">
  <img src="logo.svg" alt="ERD Viewer" width="200"/>
</div>

# ERD Viewer

Point it at your Postgres database, get an interactive diagram in your browser. Drag tables around, hover to see relationships, search to filter.

**Now available as a library!** Use the ERD Viewer components programmatically in your own Go applications.

## Quick demo
![Demo](https://assets.avi.io/Monosnap_screencast_2026-03-12_12-42-55.gif)

## What It Does

Connects to PostgreSQL, introspects the schema, renders an interactive D3.js diagram. Everything runs locally, single binary, no dependencies.

## Features

- **Interactive Visualization**: Drag tables around, hover to see relationships, search to filter
- **Table Selection**: Click to select tables for export (green highlight indicates selection)
- **Export to Draw.io**: Generate XML compatible with Confluence Draw.io plugin
- **Export to PlantUML**: Generate PlantUML ER diagram syntax
- **Select All/Deselect All**: Quickly select or deselect all tables for export
- **Real-time Selection Count**: See how many tables are currently selected

## Installation

### As a CLI Tool

Download the latest version from the [Releases](https://github.com/KensoDev/erd-viewer/releases) page that matches your OS.

### As a Library

```bash
go get github.com/kensodev/erd-viewer
```

## Usage

### CLI Usage

Flags: `-host`, `-port`, `-username`, `-db`, `-schema`, `-exclude`, `-title`, `-listen`

```bash
# Run with PostgreSQL
erd-viewer -username postgres -db mydb

# Demo with Docker
docker compose up --build
# Open http://localhost:3000
```

### Library Usage

Use ERD Viewer components programmatically in your own applications:

```go
package main

import (
    "github.com/kensodev/erd-viewer/pkg/erd"
    "github.com/kensodev/erd-viewer/pkg/erd/postgres"
    "github.com/kensodev/erd-viewer/pkg/webview"
)

func main() {
    // Create schema data programmatically
    schema := &erd.SchemaData{
        Title: "My Database",
        Tables: []erd.Table{ /* ... */ },
        FKs: []erd.ForeignKey{ /* ... */ },
    }

    // Export to PlantUML
    exporter := erd.NewPlantUMLExporter()
    output, _ := exporter.Export(schema, []string{"users", "posts"})
    println(output)

    // Or introspect a PostgreSQL database
    // introspector := postgres.NewIntrospector(conn)
    // schema, _ := introspector.IntrospectSchema(ctx, "public", map[string]bool{})

    // Or start a custom web server
    // server, _ := webview.New(webview.Config{
    //     SchemaData: schema,
    //     ListenAddr: "127.0.0.1:8080",
    //     Assets: &webview.EmbedAssets{FS: web.Files},
    // })
    // server.Start()
}
```

See the [examples](./examples/custom-webview) directory for more detailed usage examples.

### Exporting Diagrams

1. Click on tables to select them (green highlight indicates selection)
2. Use "Select All" to select all tables or "Deselect All" to clear selection
3. Click "Export to Draw.io" to download an XML file compatible with Confluence
4. Click "Export to PlantUML" to download a PlantUML diagram file
5. Import the exported files into your documentation tools

## Library Architecture

The project is now split into reusable components:

### Core Library (`pkg/erd`)
- **Types**: Core data structures (`SchemaData`, `Table`, `Column`, `ForeignKey`)
- **Exporters**: PlantUML and Draw.io exporters
- **Interface**: `Exporter` interface for custom export formats

### Database Introspection (`pkg/erd/postgres`)
- PostgreSQL-specific introspector
- Easy to extend for other databases (MySQL, SQLite, etc.)

### Web View (`pkg/webview`)
- Reusable HTTP server for ERD visualization
- Pluggable asset provider interface
- Customizable with your own HTML/CSS/JS

### Benefits
- Use components independently or together
- Build custom ERD tools for your specific needs
- Integrate ERD generation into CI/CD pipelines
- Create custom exporters for different diagram formats
- Support additional databases with custom introspectors

## License

MIT
