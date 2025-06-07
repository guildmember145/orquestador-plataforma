// services/auth-service/pkg/database/postgres.go
package database

import (
    "context"
    "log"
    "os"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/guildmember145/auth-service/pkg/config"
)

// ConnectDB establece una conexión con la base de datos PostgreSQL y devuelve un pool de conexiones.
func ConnectDB() *pgxpool.Pool {
    log.Println("Connecting to database...")
    pool, err := pgxpool.New(context.Background(), config.AppConfig.DatabaseURL)
    if err != nil {
        log.Fatalf("Unable to connect to database: %v\n", err)
        os.Exit(1)
    }

    // Hacer un ping para verificar que la conexión es exitosa
    if err := pool.Ping(context.Background()); err != nil {
        log.Fatalf("Database ping failed: %v\n", err)
        os.Exit(1)
    }

    log.Println("Database connection successful.")
    return pool
}

// RunMigrations ejecuta las sentencias SQL para crear las tablas necesarias.
func RunMigrations(pool *pgxpool.Pool) {
    log.Println("Running database migrations...")

    // Sentencia SQL para crear la tabla 'users' si no existe
    createTableSQL := `
    CREATE TABLE IF NOT EXISTS users (
        id UUID PRIMARY KEY,
        username VARCHAR(50) UNIQUE NOT NULL,
        email VARCHAR(255) UNIQUE NOT NULL,
        password_hash VARCHAR(255) NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );
    `

    _, err := pool.Exec(context.Background(), createTableSQL)
    if err != nil {
        log.Fatalf("Failed to create 'users' table: %v\n", err)
        os.Exit(1)
    }

    log.Println("Database migrations completed successfully.")
}