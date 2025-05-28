package handlers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
    "github.com/google/uuid"
    "github.com/guildmember145/auth-service/internal/auth"    
    "github.com/guildmember145/auth-service/internal/user"   
    "github.com/guildmember145/auth-service/pkg/transport" 
)

var validate = validator.New()

func RegisterHandler(c *gin.Context) {
    var req transport.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
        return
    }

    if err := validate.Struct(req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed: " + err.Error()})
        return
    }

    if _, found := user.FindByEmail(req.Email); found {
        c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
        return
    }
    // Podrías añadir validación de username único también

    hashedPassword, err := auth.HashPassword(req.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }

    newUser := &user.User{
        ID:           uuid.New(),
        Username:     req.Username,
        Email:        req.Email,
        PasswordHash: hashedPassword,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }

    if err := user.Save(newUser); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user_id": newUser.ID})
}

func LoginHandler(c *gin.Context) {
    var req transport.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
        return
    }
     if err := validate.Struct(req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed: " + err.Error()})
        return
    }


    foundUser, found := user.FindByEmail(req.Email)
    if !found || !auth.CheckPasswordHash(req.Password, foundUser.PasswordHash) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
        return
    }

    accessToken, err := auth.GenerateAccessToken(foundUser.ID, foundUser.Username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
        return
    }

    c.JSON(http.StatusOK, transport.LoginResponse{AccessToken: accessToken})
}

// Handler para /users/me (ejemplo de ruta protegida)
func GetMeHandler(c *gin.Context) {
    userIDClaim, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token claims"})
        return
    }

    userID, err := uuid.Parse(userIDClaim.(string))
    if err != nil {
         c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format in token"})
        return
    }


    currentUser, found := user.FindByID(userID)
    if !found {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    c.JSON(http.StatusOK, transport.UserResponse{
        ID:       currentUser.ID,
        Username: currentUser.Username,
        Email:    currentUser.Email,
    })
}

// Handler para validar token para otros servicios (importante para microservicios)
func ValidateTokenHandler(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"valid": false, "error": "Authorization header required"})
        return
    }

    tokenString := ""
    if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
        tokenString = authHeader[7:]
    } else {
        c.JSON(http.StatusUnauthorized, gin.H{"valid": false, "error": "Bearer token required"})
        return
    }

    claims, err := auth.ValidateToken(tokenString)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"valid": false, "error": err.Error()})
        return
    }
    // Devolver claims o solo un estado de validez
    c.JSON(http.StatusOK, gin.H{"valid": true, "user_id": claims.UserID, "username": claims.Username})
}