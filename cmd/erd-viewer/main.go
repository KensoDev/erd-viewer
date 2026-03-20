package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/jackc/pgx/v5"
	"golang.org/x/term"

	"github.com/kensodev/erd-viewer/pkg/erd/postgres"
	"github.com/kensodev/erd-viewer/pkg/webview"
	"github.com/kensodev/erd-viewer/web"
)

func main() {
	host := flag.String("host", "localhost", "Database host")
	port := flag.Int("port", 5432, "Database port")
	username := flag.String("username", "", "Database username (required)")
	dbName := flag.String("db", "", "Database name (required)")
	schema := flag.String("schema", "public", "Schema to introspect")
	exclude := flag.String("exclude", "", "Comma-separated tables to exclude")
	title := flag.String("title", "", "Diagram title")
	listen := flag.String("listen", "127.0.0.1:0", "Address to listen on (use 0.0.0.0:3000 for Docker)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "ERD Viewer - Visualize your PostgreSQL database schema interactively\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *username == "" || *dbName == "" {
		fmt.Fprintln(os.Stderr, "Error: --username and --db are required")
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*host, *port, *username, *dbName, *schema, *exclude, *title, *listen); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(host string, port int, username, dbName, schema, excludeStr, title, listenAddr string) error {
	// Read password from env var or prompt
	var pw []byte
	if envPass := os.Getenv("PGPASSWORD"); envPass != "" {
		pw = []byte(envPass)
	} else {
		fmt.Fprintf(os.Stderr, "Password for %s@%s/%s: ", username, host, dbName)
		var err error
		pw, err = term.ReadPassword(int(syscall.Stdin))
		fmt.Fprintln(os.Stderr)
		if err != nil {
			return fmt.Errorf("reading password: %w", err)
		}
	}

	// Connect to database
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", username, string(pw), host, port, dbName)
	ctx := context.Background()

	fmt.Fprintln(os.Stderr, "Connecting...")
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer conn.Close(ctx)

	// Build exclude set
	exclude := toSet(strings.Split(excludeStr, ","))

	// Introspect schema
	fmt.Fprintf(os.Stderr, "Introspecting schema %q...\n", schema)
	introspector := postgres.NewIntrospector(conn)
	schemaData, err := introspector.IntrospectSchema(ctx, schema, exclude)
	if err != nil {
		return fmt.Errorf("introspection failed: %w", err)
	}

	if len(schemaData.Tables) == 0 {
		return fmt.Errorf("no tables found in schema %q", schema)
	}

	fmt.Fprintf(os.Stderr, "Found %d tables, %d foreign key relationships.\n", len(schemaData.Tables), len(schemaData.FKs))

	// Set title
	if title == "" {
		title = fmt.Sprintf("ERD — %s / %s", dbName, schema)
	}
	schemaData.Title = title

	// Start HTTP server with the default web view
	srv, err := webview.New(webview.Config{
		SchemaData: schemaData,
		ListenAddr: listenAddr,
		Assets:     &webview.EmbedAssets{FS: web.Files},
	})
	if err != nil {
		return err
	}

	url := srv.URL()
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "🚀 ERD Viewer is running at %s\n", url)
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Open your browser and navigate to the URL above.\n")
	fmt.Fprintf(os.Stderr, "Press Ctrl+C to stop.\n")
	fmt.Fprintf(os.Stderr, "\n")

	return srv.Start()
}

func toSet(ss []string) map[string]bool {
	m := make(map[string]bool, len(ss))
	for _, s := range ss {
		if s != "" {
			m[s] = true
		}
	}
	return m
}
