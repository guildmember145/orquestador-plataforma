// services/task-orchestrator-service/pkg/database/postgres.go
package database

import (
	"context"
	"log"
	"os"

	"github.com/guildmember145/task-orchestrator-service/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB() *pgxpool.Pool {
	log.Println("Connecting to database...")
	// --- CORRECCIÓN AQUÍ: de DatabaseUrl a DatabaseURL ---
	pool, err := pgxpool.New(context.Background(), config.AppConfig.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Database ping failed: %v\n", err)
		os.Exit(1)
	}

	log.Println("Database connection successful.")
	return pool
}

func RunMigrations(pool *pgxpool.Pool) {
	log.Println("Running database migrations for task-orchestrator...")
	createWorkflowsTableSQL := `
    CREATE TABLE IF NOT EXISTS workflows (
        id UUID PRIMARY KEY,
        user_id UUID NOT NULL,
        name VARCHAR(255) NOT NULL,
        description TEXT,
        trigger JSONB NOT NULL,
        actions JSONB NOT NULL,
        is_enabled BOOLEAN NOT NULL DEFAULT TRUE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );
    `
	_, err := pool.Exec(context.Background(), createWorkflowsTableSQL)
	if err != nil {
		log.Fatalf("Failed to create 'workflows' table: %v\n", err)
		os.Exit(1)
	}
	log.Println("Orchestrator database migrations completed successfully.")
}