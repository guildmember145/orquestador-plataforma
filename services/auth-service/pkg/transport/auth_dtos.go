package transport

import "github.com/google/uuid"

type RegisterRequest struct {
    Username string `json:"username" validate:"required,min=3,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"` // OWASP: Longitud y complejidad
}

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
    AccessToken  string `json:"access_token"`
    // RefreshToken string `json:"refresh_token"` // Añadir luego
}

type UserResponse struct {
    ID       uuid.UUID `json:"id"`
    Username string    `json:"username"`
    Email    string    `json:"email"`
}