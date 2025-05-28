package config

import (
    "log"
    "os"
    "strconv"
    "time"

    "github.com/joho/godotenv"
)

type Config struct {
    Port                    string
    JWTSecretKey            string
    JWTExpiration           time.Duration
    RefreshTokenExpiration  time.Duration
}

var AppConfig Config

func LoadConfig() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }

    AppConfig.Port = getEnv("PORT", "5000")
    AppConfig.JWTSecretKey = getEnv("JWT_SECRET_KEY", "default_secret_key_please_change") // Cambia el default

    jwtExpMinutes, _ := strconv.Atoi(getEnv("JWT_EXPIRATION_MINUTES", "15"))
    AppConfig.JWTExpiration = time.Duration(jwtExpMinutes) * time.Minute

    refreshExpHours, _ := strconv.Atoi(getEnv("REFRESH_TOKEN_EXPIRATION_HOURS", "168"))
    AppConfig.RefreshTokenExpiration = time.Duration(refreshExpHours) * time.Hour
}

func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallback
}