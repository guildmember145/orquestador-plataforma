// services/auth-service/pkg/config/config.go
package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                   string
	DatabaseURL            string
	JWTSecretKey           string
	JWTExpiration          time.Duration
	RefreshTokenExpiration time.Duration
}

var AppConfig Config

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found for auth-service, using environment variables")
	}

	AppConfig.Port = getEnv("PORT", "5000")
	AppConfig.DatabaseURL = getEnv("DATABASE_URL", "")
	AppConfig.JWTSecretKey = getEnv("JWT_SECRET_KEY", "default_secret")

	// --- INICIO DE LA CORRECCIÓN ---
	jwtExpStr := getEnv("JWT_EXPIRATION_MINUTES", "15")
	jwtExpMinutes, err := strconv.Atoi(jwtExpStr)
	if err != nil {
		log.Printf("Warning: Invalid JWT_EXPIRATION_MINUTES value. Defaulting to 15. Error: %v", err)
		jwtExpMinutes = 15 // Usamos un valor por defecto seguro si la conversión falla
	}
	AppConfig.JWTExpiration = time.Duration(jwtExpMinutes) * time.Minute

	refreshExpStr := getEnv("REFRESH_TOKEN_EXPIRATION_HOURS", "168")
	refreshExpHours, err := strconv.Atoi(refreshExpStr)
	if err != nil {
		log.Printf("Warning: Invalid REFRESH_TOKEN_EXPIRATION_HOURS value. Defaulting to 168. Error: %v", err)
		refreshExpHours = 168
	}
	// --- FIN DE LA CORRECCIÓN ---
	AppConfig.RefreshTokenExpiration = time.Duration(refreshExpHours) * time.Hour
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	if fallback == "" {
		log.Fatalf("FATAL: Required environment variable %s is not set.", key)
	}
	return fallback
}