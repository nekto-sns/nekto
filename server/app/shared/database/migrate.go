package database

import (
	"context"
	"embed"
	"fmt"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

func RunMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version TEXT PRIMARY KEY,
		applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	`)
	if err != nil {
		return fmt.Errorf("failed to create migration table: %w", err)
	}

	entries, err := migrationFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		filename := entry.Name()

		var exists bool
		err := pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", filename).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check migration history for %s: %w", filename, err)
		}

		if !exists {
			fmt.Printf("[Migration] Applying: %s\n", filename)

			content, err := migrationFS.ReadFile("migrations/" + filename)
			if err != nil {
				return fmt.Errorf("failed to read migration file %s: %w", filename, err)
			}

			tx, err := pool.Begin(ctx)
			if err != nil {
				return fmt.Errorf("failed to begin transaction for %s: %w", filename, err)
			}
			defer tx.Rollback(ctx)

			if _, err := tx.Exec(ctx, string(content)); err != nil {
				return fmt.Errorf("error executing SQL in %s: %w", filename, err)
			}

			if _, err := tx.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", filename); err != nil {
				return fmt.Errorf("failed to record migration history for %s: %w", filename, err)
			}

			if err := tx.Commit(ctx); err != nil {
				return fmt.Errorf("failed to commit transaction for %s: %w", filename, err)
			}
			fmt.Printf("[Migration] Success: %s\n", filename)
		}
	}

	return nil
}
