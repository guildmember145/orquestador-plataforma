package main

import (
    "log"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/guildmember145/auth-service/internal/handlers"   
    "github.com/guildmember145/auth-service/internal/middleware" 
    "github.com/guildmember145/auth-service/pkg/config"       
)

func main() {
    config.LoadConfig()

    router := gin.Default() // Logger y Recovery middleware por defecto

    // CORS Middleware (importante para desarrollo con Vue en otro puerto)
    router.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Sé más restrictivo en producción
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    })


    // Rutas públicas
    authRoutes := router.Group("/api/baas/v1/auth")
    {
        authRoutes.POST("/register", handlers.RegisterHandler)
        authRoutes.POST("/login", handlers.LoginHandler)
        // Esta ruta es para que otros servicios validen un token,
        // podría requerir su propia forma de autenticación (ej. API key de servicio a servicio)
        // o estar en una red interna no expuesta públicamente.
        // Por ahora, la dejamos abierta en la red interna del Podman.
        authRoutes.POST("/validate_token", handlers.ValidateTokenHandler)
    }

    // Rutas protegidas (ejemplo)
    userRoutes := router.Group("/api/baas/v1/users")
    userRoutes.Use(middleware.AuthMiddleware()) // Aplicar middleware de autenticación
    {
        userRoutes.GET("/me", handlers.GetMeHandler)
    }

    addr := fmt.Sprintf(":%s", config.AppConfig.Port)
    log.Printf("Auth service starting on %s", addr)
    if err := router.Run(addr); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}