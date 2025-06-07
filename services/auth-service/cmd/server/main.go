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

	userStore := user.NewPostgresUserStore(dbPool)
	authHandler := handlers.NewAuthHandler(userStore)

	router := gin.Default()

	// --- INICIO DE LA CORRECCIÓN DE CORS ---
	// Reemplazamos nuestro middleware manual con uno profesional.
	corsConfig := cors.DefaultConfig()
	// Permitimos explícitamente el origen de nuestro frontend Vue.
	corsConfig.AllowOrigins = []string{"http://localhost:3003"}
	// Permitimos credenciales (importante para tokens o cookies en el futuro)
	corsConfig.AllowCredentials = true
	// Especificamos las cabeceras que el frontend puede enviar.
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}

	router.Use(cors.New(corsConfig))
	// --- FIN DE LA CORRECCIÓN DE CORS ---

	// El resto de tus rutas no cambia...
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