<div align="center">
  <img src="logo.svg" alt="ERD Viewer Logo" width="200"/>
  <h1>ERD Viewer</h1>
  <p><em>Visualize your PostgreSQL schema without the headache</em></p>
</div>

---

I got tired of squinting at `\d+ table_name` in psql and alt-tabbing between terminal windows to understand how my database tables relate to each other. You know the drill - you're debugging some gnarly query and you can't remember if `users.id` references `orders.user_id` or the other way around, and now you're three Google searches deep into "postgres show foreign keys".

So I built this. Point it at your Postgres database, and it spins up an interactive diagram in your browser. Drag tables around. Hover to see relationships. Search to filter. That's it.

## What It Does

ERD Viewer connects to your PostgreSQL database, introspects the schema, and renders an interactive entity-relationship diagram using D3.js. You get:

- **Drag-and-drop interface** - rearrange tables however makes sense to you
- **Smart highlighting** - hover over a table to see all its relationships light up
- **Real-time search** - filter tables as you type
- **Auto-layout** - starts with a sensible grid layout, but you can mess with it however you want
- **Primary keys & foreign keys** - clearly marked with visual indicators
- **Zero configuration** - no config files, no setup, just run it and open the URL

## Why It Works

Look, I know there are fancy ER diagram tools out there. But they either:
1. Require uploading your schema to some cloud service (no thanks)
2. Need 47 dependencies and a PhD to install
3. Generate static images that are useless the moment your schema changes
4. Cost money

This tool runs locally, starts in seconds, and updates instantly when you point it at a different schema. The whole thing is a single binary. That's the way it should be.

## Quick Demo

Want to see it in action first? We've got a Docker Compose setup with the Northwind sample database:

```bash
docker compose up --build
```

That's it. Open http://localhost:3000 in your browser and you'll see an interactive diagram of the Northwind schema. Drag tables around, hover to see relationships, search to filter.

The demo spins up:
- PostgreSQL with the Northwind database pre-loaded
- ERD Viewer connected to it on port 3000

See [examples/README.md](examples/README.md) for more details.

## Installation

You need Go installed. If you don't have it, grab it from [golang.org](https://golang.org/dl/).

```bash
go install github.com/kensodev/erd-viewer/cmd/erd-viewer@latest
```

Or clone and build:

```bash
git clone https://github.com/kensodev/erd-viewer.git
cd erd-viewer
make install
```

## Usage

Basic usage:

```bash
erd-viewer --username your_user --db your_database
```

It'll prompt for your password (securely, no echo), connect to `localhost:5432`, introspect the `public` schema, and print the URL. Open it in your browser.

### Options

```bash
erd-viewer \
  --host localhost \
  --port 5432 \
  --username postgres \
  --db myapp_development \
  --schema public \
  --exclude "schema_migrations,ar_internal_metadata" \
  --title "My App ERD"
```

- `--host` - Database host (default: localhost)
- `--port` - Database port (default: 5432)
- `--username` - Database username (required)
- `--db` - Database name (required)
- `--schema` - Schema to introspect (default: public)
- `--exclude` - Comma-separated list of tables to skip
- `--title` - Custom title for the diagram

## Real-World Example

I use this when I'm onboarding to a new codebase. Instead of reading through 47 migration files to understand the data model, I just:

```bash
erd-viewer --username dev --db my_new_project
```

And boom - instant visual understanding of the whole schema. Drag the `users` table to the middle, see what references it, trace the relationships. Makes way more sense than staring at SQL dumps.

## How It Works

1. Connects to your Postgres database
2. Queries `information_schema` for tables and columns
3. Queries `pg_catalog` for foreign key relationships
4. Starts an HTTP server on a random free port
5. Renders an interactive D3.js visualization in your browser
6. Everything stays local - no data leaves your machine

The frontend is pure JavaScript with D3.js for the graph layout. Tables are SVG elements that you can drag around. Foreign key relationships are Bezier curves that update as you move things.

## Development

```bash
# Run tests
make test

# Build binary
make build

# Run locally
make run ARGS="--username postgres --db mydb"

# Clean up
make clean
```

## Architecture

The code is organized like a principal engineer would write it (that's the goal, anyway):

```
├── cmd/erd-viewer/          # CLI entry point
├── internal/
│   ├── db/                  # Database introspection logic
│   ├── server/              # HTTP server
│   └── browser/             # Browser launcher
└── web/
    ├── static/
    │   ├── css/             # Styles
    │   └── js/              # D3.js visualization
    └── templates/           # HTML
```

## Testing

Run the tests with:

```bash
make test
```

The test suite covers database introspection, schema parsing, and foreign key detection. No mocks - we use testcontainers to spin up real Postgres instances for integration tests, because that's how you know it actually works.

## Contributing

Found a bug? Open an issue. Want to add a feature? Open a PR. Standard GitHub workflow.

## License

MIT. Do whatever you want with it.

## Why "ERD Viewer"?

Because I'm bad at naming things and this tool views ERDs. If you have a better name, I'm all ears.

---

Built by [KensoDev](https://github.com/kensodev) because SQL `\d+` wasn't cutting it anymore.
