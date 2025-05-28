package user

import (
    "time"
    "github.com/google/uuid"
)

type User struct {
    ID           uuid.UUID `json:"id"`
    Username     string    `json:"username" validate:"required,min=3,max=50"`
    Email        string    `json:"email" validate:"required,email"`
    PasswordHash string    `json:"-"` // No exponer en JSON
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

// Almacenamiento en memoria para empezar (simplificado)
// En un sistema real, esto sería una base de datos.
var userStore = make(map[string]*User) // email como clave
var userStoreByID = make(map[uuid.UUID]*User)

// Funciones de ejemplo para interactuar con el store (CRUD simplificado)
func FindByEmail(email string) (*User, bool) {
    user, found := userStore[email]
    return user, found
}

func FindByID(id uuid.UUID) (*User, bool) {
    user, found := userStoreByID[id]
    return user, found
}

func Save(user *User) error {
    // Aquí iría la lógica de validación antes de guardar
    userStore[user.Email] = user
    userStoreByID[user.ID] = user
    return nil
}