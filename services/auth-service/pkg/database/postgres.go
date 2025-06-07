// services/auth-service/pkg/database/postgres.go
package database

import (
    "context"
    "log"
    "os"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/guildmember145/auth-service/pkg/config"
)

func ConnectDB() *pgxpool.Pool {
    pool, err := pgxpool.New(context.Background(), config.AppConfig.DatabaseURL)
    if err != nil {
        log.Fatalf("Unable to connect to database: %v\n", err)
        os.Exit(1)
    }
    if err := pool.Ping(context.Background()); err != nil {
        log.Fatalf("Database ping failed: %v\n", err)
        os.Exit(1)
    }
    log.Println("Database connection successful for auth-service.")
    return pool
}

func RunMigrations(pool *pgxpool.Pool) {
    log.Println("Running database migrations for auth-service...")
    createTableSQL := `
    CREATE TABLE IF NOT EXISTS users (
        id UUID PRIMARY KEY,
        username VARCHAR(50) UNIQUE NOT NULL,
        email VARCHAR(255) UNIQUE NOT NULL,
        password_hash VARCHAR(255) NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );`
    _, err := pool.Exec(context.Background(), createTableSQL)
    if err != nil {
        log.Fatalf("Failed to create 'users' table: %v\n", err)
        os.Exit(1)
    }
    log.Println("Auth-service database migrations completed successfully.")
}