<div align="center">
  <img src="logo.svg" alt="ERD Viewer" width="200"/>
</div>

# ERD Viewer

Point it at your Postgres database, get an interactive diagram in your browser. Drag tables around, hover to see relationships, search to filter.

## What It Does

Connects to PostgreSQL, introspects the schema, renders an interactive D3.js diagram. Everything runs locally, single binary, no dependencies.

## Usage

```bash
# Install
go install github.com/kensodev/erd-viewer/cmd/erd-viewer@latest

# Run
erd-viewer --username postgres --db mydb

# Demo with Docker
docker compose up --build
# Open http://localhost:3000
```

Flags: `--host`, `--port`, `--username`, `--db`, `--schema`, `--exclude`, `--title`, `--listen`

## License

MIT
