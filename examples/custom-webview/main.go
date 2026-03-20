package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/kensodev/erd-viewer/pkg/erd"
	"github.com/kensodev/erd-viewer/pkg/erd/postgres"
	"github.com/kensodev/erd-viewer/pkg/webview"
)

// CustomAssets demonstrates how to provide your own web assets
type CustomAssets struct {
	basePath string
}

func (c *CustomAssets) ReadFile(name string) ([]byte, error) {
	// In a real implementation, you would load custom HTML/CSS/JS files
	// from your own filesystem or embed them using go:embed
	return os.ReadFile(c.basePath + "/" + name)
}

func main() {
	// Example 1: Using the library for programmatic export
	fmt.Println("Example 1: Programmatic Export")
	fmt.Println("================================")

	// Create sample schema data
	schema := &erd.SchemaData{
		Title: "Example Database",
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
					{Name: "content", Type: "text", IsPK: false, Nullable: true},
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
		log.Fatalf("PlantUML export failed: %v", err)
	}
	fmt.Println("PlantUML Output:")
	fmt.Println(plantUML)
	fmt.Println()

	// Export to Draw.io
	drawioExporter := erd.NewDrawioExporter()
	drawio, err := drawioExporter.Export(schema, []string{"users", "posts"})
	if err != nil {
		log.Fatalf("Draw.io export failed: %v", err)
	}
	fmt.Printf("Draw.io XML Output: %d bytes generated\n\n", len(drawio))

	// Example 2: Connect to a real database and introspect
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		fmt.Println("Example 2: Database Introspection")
		fmt.Println("==================================")

		ctx := context.Background()
		conn, err := pgx.Connect(ctx, dsn)
		if err != nil {
			log.Printf("Database connection failed (skipping): %v", err)
		} else {
			defer conn.Close(ctx)

			introspector := postgres.NewIntrospector(conn)
			realSchema, err := introspector.IntrospectSchema(ctx, "public", map[string]bool{})
			if err != nil {
				log.Printf("Introspection failed: %v", err)
			} else {
				fmt.Printf("Found %d tables and %d foreign keys\n", len(realSchema.Tables), len(realSchema.FKs))

				// Export the real schema
				output, _ := plantUMLExporter.Export(realSchema, getAllTableNames(realSchema))
				fmt.Println("\nReal Database Schema (PlantUML):")
				fmt.Println(output)
			}
		}
	}

	// Example 3: Custom web view (commented out - requires custom assets)
	// To use this, you would need to create your own HTML/CSS/JS files
	/*
	fmt.Println("Example 3: Custom Web View")
	fmt.Println("===========================")

	server, err := webview.New(webview.Config{
		SchemaData: schema,
		ListenAddr: "127.0.0.1:8080",
		Assets:     &CustomAssets{basePath: "./custom-web-assets"},
	})
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	fmt.Printf("Custom web view available at %s\n", server.URL())
	log.Fatal(server.Start())
	*/
}

func getAllTableNames(schema *erd.SchemaData) []string {
	names := make([]string, len(schema.Tables))
	for i, table := range schema.Tables {
		names[i] = table.Name
	}
	return names
}
