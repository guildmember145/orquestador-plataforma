// services/task-orchestrator-service/pkg/config/config.go
package config // <-- Declaración del paquete

import (
    "log"
    "os"
    "github.com/joho/godotenv" // Asegúrate de tener esta dependencia si la usas
)

type Config struct {
    Port                string
    AuthServiceBaseURL  string
}

var AppConfig Config

func LoadConfig() {
    // Cargar .env desde el directorio actual del ejecutable (task-orchestrator-service)
    // Esto es importante porque el .env debe estar en la raíz de este servicio
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found for task-orchestrator-service, using environment variables")
    }
    AppConfig.Port = getEnv("PORT", "8080") // Puerto interno para task-orchestrator-service
    AppConfig.AuthServiceBaseURL = getEnv("AUTH_SERVICE_BASE_URL", "http://auth_service:5000/api/baas/v1/auth")
    log.Println("Configuration loaded for task-orchestrator-service")
}

func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    // log.Printf("Using fallback for %s: %s", key, fallback) // Puedes descomentar para debug
    return fallback
}