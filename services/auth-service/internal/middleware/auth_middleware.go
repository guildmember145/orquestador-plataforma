package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/guildmember145/auth-service/internal/auth" // Ajusta path
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
        claims, err := auth.ValidateToken(tokenString)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
            return
        }

        // Poner información del usuario en el contexto de Gin para uso posterior
        c.Set("userID", claims.UserID.String()) // Guardamos como string para evitar problemas de tipo
        c.Set("username", claims.Username)
        c.Next()
    }
}