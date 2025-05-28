// services/task-orchestrator-service/internal/middleware/auth_middleware.go
package middleware // <-- Declaración del paquete

import (
    "net/http"
    "strings"
    "log" // Para depuración si es necesario

    "github.com/gin-gonic/gin"
    // Asegúrate que esta ruta de importación sea correcta para tu estructura y el nombre de tu módulo
    "github.com/guildmember145/task-orchestrator-service/internal/services"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
            return
        }
        tokenString := parts[1]

        validationResp, err := services.ValidateTokenWithAuthService(tokenString)
        if err != nil {
            log.Printf("Error validating token with auth service: %v", err) // Loguea el error real
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token validation failed or auth service unreachable"})
            return
        }

        if !validationResp.Valid {
            log.Printf("Token deemed invalid by auth service. Details: %s", validationResp.Error)
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": validationResp.Error})
            return
        }

        c.Set("userID", validationResp.UserID)
        c.Set("username", validationResp.Username)
        c.Next()
    }
}