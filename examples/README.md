# Running the Demo

This directory contains a sample Northwind database schema to demo ERD Viewer.

## Quick Start

From the project root:

```bash
docker compose up --build
```

Then open http://localhost:3000 in your browser. That's it!

## Running Locally (without Docker)

1. Start just the Postgres database:
   ```bash
   docker compose up -d postgres
   ```

2. Run ERD Viewer locally:
   ```bash
   # If you have the binary installed
   erd-viewer --username demo --db northwind --host localhost --port 5432
   # Password: demo

   # Or build and run
   make run ARGS="--username demo --db northwind --host localhost --port 5432"
   # Password: demo
   ```

3. Your browser will open automatically with the interactive ERD!

## What You'll See

The Northwind database is a classic sample database with:
- **Customers** - who place orders
- **Orders** - linked to customers and employees
- **Order Details** - the line items for each order
- **Products** - what's being sold
- **Categories** - product categories
- **Suppliers** - who supplies the products
- **Employees** - who process orders (with self-referential reporting structure)
- **Shippers** - who ship the orders

All the foreign key relationships are visible in the diagram. Hover over a table to see its connections light up. Drag tables around to organize them however makes sense.

## Cleanup

```bash
docker-compose down
```

To completely remove the data:
```bash
docker-compose down -v
```

## Using Your Own Database

Just point ERD Viewer at any Postgres database:

```bash
erd-viewer --username youruser --db yourdb --host yourhost --schema yourschema
```
