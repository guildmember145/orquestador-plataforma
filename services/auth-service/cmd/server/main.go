// services/auth-service/cmd/server/main.go
package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/guildmember145/auth-service/internal/handlers"
	"github.com/guildmember145/auth-service/internal/middleware"
	"github.com/guildmember145/auth-service/internal/user"
	"github.com/guildmember145/auth-service/pkg/config"
	"github.com/guildmember145/auth-service/pkg/database"
)

func main() {
	config.LoadConfig()
	dbPool := database.ConnectDB()
	defer dbPool.Close()
	database.RunMigrations(dbPool)

	// Usamos la implementación de PostgreSQL
	userStore := user.NewPostgresUserStore(dbPool)
	authHandler := handlers.NewAuthHandler(userStore)

	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3003"}
    // ... (resto de tu config de CORS)
	router.Use(cors.New(corsConfig))

	authRoutes := router.Group("/api/baas/v1/auth")
	{
		authRoutes.POST("/register", authHandler.RegisterHandler)
		authRoutes.POST("/login", authHandler.LoginHandler)
		authRoutes.POST("/validate_token", authHandler.ValidateTokenHandler)
	}

	userRoutes := router.Group("/api/baas/v1/users")
	userRoutes.Use(middleware.AuthMiddleware(userStore))
	{
		userRoutes.GET("/me", authHandler.GetMeHandler)
	}

	addr := fmt.Sprintf(":%s", config.AppConfig.Port)
	log.Printf("Auth service starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}