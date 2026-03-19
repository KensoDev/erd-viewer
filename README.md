<div align="center">
  <img src="logo.svg" alt="ERD Viewer" width="200"/>
</div>

# ERD Viewer

Point it at your Postgres database, get an interactive diagram in your browser. Drag tables around, hover to see relationships, search to filter.

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

## Usage

Download the latest veriosn from the [Releases](https://github.com/KensoDev/erd-viewer/releases) page that matches your os.

Flags: `-host`, `-port`, `-username`, `-db`, `-schema`, `-exclude`, `-title`, `-listen`

```bash
# Run
erd-viewer -username postgres -db mydb

# Demo with Docker
docker compose up --build
# Open http://localhost:3000
```

### Exporting Diagrams

1. Click on tables to select them (green highlight indicates selection)
2. Use "Select All" to select all tables or "Deselect All" to clear selection
3. Click "Export to Draw.io" to download an XML file compatible with Confluence
4. Click "Export to PlantUML" to download a PlantUML diagram file
5. Import the exported files into your documentation tools

## License

MIT
