// services/auth-service/internal/user/postgres_store.go
package user

import (
    "context"
    "fmt"
    "log"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
)


type Store interface {
    Save(user *User) error
    FindByEmail(email string) (*User, error)
    FindByID(id uuid.UUID) (*User, error)
}

type PostgresUserStore struct {
    DB *pgxpool.Pool
}

// NewPostgresUserStore crea una nueva instancia de PostgresUserStore.
func NewPostgresUserStore(db *pgxpool.Pool) *PostgresUserStore {
    return &PostgresUserStore{DB: db}
}

// Save inserta un nuevo usuario en la tabla 'users'.
func (s *PostgresUserStore) Save(user *User) error {
    query := `
        INSERT INTO users (id, username, email, password_hash, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
    _, err := s.DB.Exec(context.Background(), query,
        user.ID, user.Username, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt)

    if err != nil {
        log.Printf("Error saving user to database: %v", err)
        // Aquí podrías manejar errores específicos de la BD, como 'violación de clave única'.
        return fmt.Errorf("could not save user: %w", err)
    }
    return nil
}

// FindByEmail busca un usuario por su dirección de email.
func (s *PostgresUserStore) FindByEmail(email string) (*User, error) {
    query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE email = $1`

    user := &User{}
    err := s.DB.QueryRow(context.Background(), query, email).Scan(
        &user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)

    if err != nil {
        // pgx.ErrNoRows es un error común que significa que no se encontró el usuario.
        // Lo manejamos para no devolver un error genérico.
        log.Printf("Error finding user by email '%s': %v", email, err)
        return nil, fmt.Errorf("user not found")
    }
    return user, nil
}

// FindByID busca un usuario por su UUID.
func (s *PostgresUserStore) FindByID(id uuid.UUID) (*User, error) {
    query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE id = $1`

    user := &User{}
    err := s.DB.QueryRow(context.Background(), query, id).Scan(
        &user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)

    if err != nil {
        log.Printf("Error finding user by ID '%s': %v", id, err)
        return nil, fmt.Errorf("user not found")
    }
    return user, nil
}