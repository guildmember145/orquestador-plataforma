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
	
	// Crear tabla workflows
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
	log.Println("✓ 'workflows' table created successfully")

	// Crear tabla workflow_executions
	createWorkflowExecutionsTableSQL := `
    CREATE TABLE IF NOT EXISTS workflow_executions (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        workflow_id UUID NOT NULL,
        user_id VARCHAR(255) NOT NULL,
        status VARCHAR(50) NOT NULL DEFAULT 'running',
        triggered_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        completed_at TIMESTAMPTZ,
        logs JSONB DEFAULT '[]'::jsonb,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        
        -- Restricciones y referencias
        CONSTRAINT fk_workflow_executions_workflow_id 
            FOREIGN KEY (workflow_id) REFERENCES workflows(id) ON DELETE CASCADE
    );
    `
	_, err = pool.Exec(context.Background(), createWorkflowExecutionsTableSQL)
	if err != nil {
		log.Fatalf("Failed to create 'workflow_executions' table: %v\n", err)
		os.Exit(1)
	}
	log.Println("✓ 'workflow_executions' table created successfully")

	// Crear índices para mejorar el rendimiento
	createIndexesSQL := `
    CREATE INDEX IF NOT EXISTS idx_workflow_executions_workflow_id ON workflow_executions(workflow_id);
    CREATE INDEX IF NOT EXISTS idx_workflow_executions_user_id ON workflow_executions(user_id);
    CREATE INDEX IF NOT EXISTS idx_workflow_executions_status ON workflow_executions(status);
    CREATE INDEX IF NOT EXISTS idx_workflow_executions_triggered_at ON workflow_executions(triggered_at DESC);
    `
	_, err = pool.Exec(context.Background(), createIndexesSQL)
	if err != nil {
		log.Printf("Warning: Failed to create some indexes: %v", err)
		// No terminamos la aplicación por esto, solo advertimos
	} else {
		log.Println("✓ Indexes created successfully")
	}

	// Crear función y trigger para updated_at automático
	createTriggerSQL := `
    CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
    BEGIN
        NEW.updated_at = NOW();
        RETURN NEW;
    END;
    $$ language 'plpgsql';

    CREATE TRIGGER IF NOT EXISTS update_workflow_executions_updated_at 
        BEFORE UPDATE ON workflow_executions 
        FOR EACH ROW 
        EXECUTE FUNCTION update_updated_at_column();
    `
	_, err = pool.Exec(context.Background(), createTriggerSQL)
	if err != nil {
		log.Printf("Warning: Failed to create trigger: %v", err)
		// No es crítico, continuamos
	} else {
		log.Println("✓ Auto-update trigger created successfully")
	}

	log.Println("🎉 Orchestrator database migrations completed successfully.")
}