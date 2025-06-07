// services/auth-service/internal/handlers/auth_handler.go
package handlers

import (
    "net/http"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
    "github.com/google/uuid"
    "github.com/guildmember145/auth-service/internal/auth"
    "github.com/guildmember145/auth-service/internal/user"
    "github.com/guildmember145/auth-service/pkg/transport"
)

var validate = validator.New()

type AuthHandler struct {
    UserStore user.Store
}

// Definimos una interfaz para el store para facilitar las pruebas en el futuro
type UserStoreInterface interface {
    Save(user *user.User) error
    FindByEmail(email string) (*user.User, error)
    FindByID(id uuid.UUID) (*user.User, error)
}

// NewAuthHandler ahora acepta la interfaz, no la implementación concreta
func NewAuthHandler(userStore user.Store) *AuthHandler {
    return &AuthHandler{UserStore: userStore}
}

func (h *AuthHandler) RegisterHandler(c *gin.Context) {
    var req transport.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()}); return }
    if err := validate.Struct(req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed: " + err.Error()}); return }

    _, err := h.UserStore.FindByEmail(req.Email)
    if err == nil {
        c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"}); return
    }

    hashedPassword, err := auth.HashPassword(req.Password)
    if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"}); return }

    newUser := &user.User{
        ID:           uuid.New(),
        Username:     req.Username,
        Email:        req.Email,
        PasswordHash: hashedPassword,
        CreatedAt:    time.Now().UTC(),
        UpdatedAt:    time.Now().UTC(),
    }

    if err := h.UserStore.Save(newUser); err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"}); return }
    c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user_id": newUser.ID})
}

func (h *AuthHandler) LoginHandler(c *gin.Context) {
    var req transport.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"}); return }
    if err := validate.Struct(req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed: " + err.Error()}); return }

    foundUser, err := h.UserStore.FindByEmail(req.Email)
    if err != nil || !auth.CheckPasswordHash(req.Password, foundUser.PasswordHash) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"}); return
    }

    accessToken, err := auth.GenerateAccessToken(foundUser.ID, foundUser.Username)
    if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"}); return }
    c.JSON(http.StatusOK, transport.LoginResponse{AccessToken: accessToken})
}

// --- INICIO DE LA LÓGICA CORREGIDA ---
func (h *AuthHandler) GetMeHandler(c *gin.Context) {
    userIDClaim, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token claims"}); return
    }

    userID, err := uuid.Parse(userIDClaim.(string))
    if err != nil {
         c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format in token"}); return
    }

    currentUser, err := h.UserStore.FindByID(userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"}); return
    }

    c.JSON(http.StatusOK, transport.UserResponse{
        ID:       currentUser.ID,
        Username: currentUser.Username,
        Email:    currentUser.Email,
    })
}

func (h *AuthHandler) ValidateTokenHandler(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" { c.JSON(http.StatusUnauthorized, gin.H{"valid": false, "error": "Authorization header required"}); return }

    tokenString := strings.TrimPrefix(authHeader, "Bearer ")
    if tokenString == authHeader {
        c.JSON(http.StatusUnauthorized, gin.H{"valid": false, "error": "Bearer token required"}); return
    }

    claims, err := auth.ValidateToken(tokenString)
    if err != nil { c.JSON(http.StatusUnauthorized, gin.H{"valid": false, "error": err.Error()}); return }

    // Verificamos que el usuario aún exista en la BD
    _, err = h.UserStore.FindByID(claims.UserID)
    if err != nil {
         c.JSON(http.StatusUnauthorized, gin.H{"valid": false, "error": "User in token no longer exists"}); return
    }

    c.JSON(http.StatusOK, gin.H{"valid": true, "user_id": claims.UserID, "username": claims.Username})
}
// --- FIN DE LA LÓGICA CORREGIDA ---