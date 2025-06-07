// services/task-orchestrator-service/pkg/config/config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	AuthServiceBaseURL string
	DatabaseURL        string // Nomenclatura correcta (URL en mayúsculas)
}

var AppConfig Config

func LoadConfig() {
	// Cargar .env desde el directorio actual del ejecutable
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found for task-orchestrator-service, using environment variables")
	}

	// 1. Usamos el puerto 9090 como default para evitar conflictos
	AppConfig.Port = getEnv("PORT", "9090")
	
    // 2. URL para comunicación interna entre contenedores
	AppConfig.AuthServiceBaseURL = getEnv("AUTH_SERVICE_BASE_URL", "http://auth_service:5000/api/baas/v1/auth")
	
    // 3. La URL de la BD es requerida. La app fallará si no se provee.
	AppConfig.DatabaseURL = getEnv("DATABASE_URL", "")

	log.Println("Configuration loaded for task-orchestrator-service")
}

// getEnv obtiene una variable de entorno o devuelve un valor por defecto.
// Si el valor por defecto es una cadena vacía "", la variable es considerada requerida
// y el programa terminará si no se encuentra.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	if fallback == "" {
		log.Fatalf("FATAL: Required environment variable %s is not set.", key)
	}

	return fallback
}