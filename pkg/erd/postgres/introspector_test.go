package postgres

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

func setupTestDB(t *testing.T) (*pgx.Conn, func()) {
	t.Helper()

	// Use TEST_DATABASE_URL env var or skip if not set
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping integration tests")
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	// Create test schema
	_, err = conn.Exec(ctx, `
		DROP SCHEMA IF EXISTS test_erd CASCADE;
		CREATE SCHEMA test_erd;

		CREATE TABLE test_erd.users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) NOT NULL,
			email VARCHAR(100) NOT NULL
		);

		CREATE TABLE test_erd.posts (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			title VARCHAR(200) NOT NULL,
			content TEXT,
			FOREIGN KEY (user_id) REFERENCES test_erd.users(id)
		);

		CREATE TABLE test_erd.comments (
			id SERIAL PRIMARY KEY,
			post_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			body TEXT NOT NULL,
			FOREIGN KEY (post_id) REFERENCES test_erd.posts(id),
			FOREIGN KEY (user_id) REFERENCES test_erd.users(id)
		);

		CREATE TABLE test_erd.ignored_table (
			id SERIAL PRIMARY KEY
		);
	`)
	if err != nil {
		conn.Close(ctx)
		t.Fatalf("failed to create test schema: %v", err)
	}

	cleanup := func() {
		conn.Exec(ctx, "DROP SCHEMA IF EXISTS test_erd CASCADE")
		conn.Close(ctx)
	}

	return conn, cleanup
}

func TestIntrospector_IntrospectSchema(t *testing.T) {
	conn, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	introspector := NewIntrospector(conn)

	// Test without exclusions
	t.Run("no exclusions", func(t *testing.T) {
		data, err := introspector.IntrospectSchema(ctx, "test_erd", map[string]bool{})
		if err != nil {
			t.Fatalf("IntrospectSchema failed: %v", err)
		}

		if len(data.Tables) != 4 {
			t.Errorf("expected 4 tables, got %d", len(data.Tables))
		}

		// Check tables are sorted
		expectedOrder := []string{"comments", "ignored_table", "posts", "users"}
		for i, table := range data.Tables {
			if table.Name != expectedOrder[i] {
				t.Errorf("table %d: expected %s, got %s", i, expectedOrder[i], table.Name)
			}
		}

		// Check foreign keys
		if len(data.FKs) != 3 {
			t.Errorf("expected 3 foreign keys, got %d", len(data.FKs))
		}
	})

	// Test with exclusions
	t.Run("with exclusions", func(t *testing.T) {
		exclude := map[string]bool{"ignored_table": true}
		data, err := introspector.IntrospectSchema(ctx, "test_erd", exclude)
		if err != nil {
			t.Fatalf("IntrospectSchema failed: %v", err)
		}

		if len(data.Tables) != 3 {
			t.Errorf("expected 3 tables (excluded 1), got %d", len(data.Tables))
		}

		for _, table := range data.Tables {
			if table.Name == "ignored_table" {
				t.Error("ignored_table should have been excluded")
			}
		}
	})
}

func TestIntrospector_FetchColumns(t *testing.T) {
	conn, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	introspector := NewIntrospector(conn)

	columns, err := introspector.fetchColumns(ctx, "test_erd", map[string]bool{})
	if err != nil {
		t.Fatalf("fetchColumns failed: %v", err)
	}

	// Check users table
	userCols, ok := columns["users"]
	if !ok {
		t.Fatal("users table not found")
	}

	if len(userCols) != 3 {
		t.Errorf("expected 3 columns in users table, got %d", len(userCols))
	}

	// Check primary key detection
	var foundPK bool
	for _, col := range userCols {
		if col.Name == "id" {
			if !col.IsPK {
				t.Error("id column should be marked as primary key")
			}
			foundPK = true
		}
	}
	if !foundPK {
		t.Error("primary key not found in users table")
	}

	// Check nullable detection
	for _, col := range userCols {
		if col.Name == "id" && col.Nullable {
			t.Error("id column should not be nullable")
		}
	}
}

func TestIntrospector_FetchForeignKeys(t *testing.T) {
	conn, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	introspector := NewIntrospector(conn)

	fks, err := introspector.fetchForeignKeys(ctx, "test_erd", map[string]bool{})
	if err != nil {
		t.Fatalf("fetchForeignKeys failed: %v", err)
	}

	if len(fks) != 3 {
		t.Errorf("expected 3 foreign keys, got %d", len(fks))
	}

	// Verify specific foreign keys
	var foundPostsFK bool
	for _, fk := range fks {
		if fk.FromTable == "posts" && fk.FromCol == "user_id" &&
			fk.ToTable == "users" && fk.ToCol == "id" {
			foundPostsFK = true
		}
	}
	if !foundPostsFK {
		t.Error("expected foreign key from posts.user_id to users.id not found")
	}

	// Test exclusions
	excludeFKs, err := introspector.fetchForeignKeys(ctx, "test_erd", map[string]bool{"posts": true})
	if err != nil {
		t.Fatalf("fetchForeignKeys with exclusions failed: %v", err)
	}

	// Should exclude FKs from posts table and to posts table
	for _, fk := range excludeFKs {
		if fk.FromTable == "posts" || fk.ToTable == "posts" {
			t.Error("foreign keys involving excluded table should be filtered")
		}
	}
}
