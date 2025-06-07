// services/auth-service/internal/middleware/auth_middleware.go
package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/guildmember145/auth-service/internal/auth"
    "github.com/guildmember145/auth-service/internal/user" // <-- AÑADIDO
)

// AuthMiddleware ahora acepta el user.Store para verificar la existencia del usuario
func AuthMiddleware(userStore user.Store) gin.HandlerFunc {
    return func(c *gin.Context) {
        // ... (la lógica para extraer el token es la misma) ...
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" { c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"}); return }
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader { c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"}); return }

        claims, err := auth.ValidateToken(tokenString)
        if err != nil { c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()}); return }

        // Verificación Adicional: ¿El usuario del token todavía existe en la BD?
        _, err = userStore.FindByID(claims.UserID)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User from token no longer exists"})
            return
        }

        c.Set("userID", claims.UserID.String())
        c.Set("username", claims.Username)
        c.Next()
    }
}