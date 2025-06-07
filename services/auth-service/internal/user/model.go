// services/auth-service/internal/user/model.go
package user

import (
    "time"
    "github.com/google/uuid"
)

type User struct {
    ID           uuid.UUID `json:"id"`
    Username     string    `json:"username" validate:"required,min=3,max=50"`
    Email        string    `json:"email" validate:"required,email"`
    PasswordHash string    `json:"-"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}