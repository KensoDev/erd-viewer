<div align="center">
  <img src="logo.svg" alt="ERD Viewer" width="200"/>
</div>

# ERD Viewer

Point it at your Postgres database, get an interactive diagram in your browser. Drag tables around, hover to see relationships, search to filter.

## Quick demo
![Demo](https://assets.avi.io/Monosnap_screencast_2026-03-12_12-42-55.gif)

## What It Does

Connects to PostgreSQL, introspects the schema, renders an interactive D3.js diagram. Everything runs locally, single binary, no dependencies.

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

## License

MIT
