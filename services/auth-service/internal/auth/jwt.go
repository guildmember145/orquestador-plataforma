package auth

import (
    "fmt"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/guildmember145/auth-service/pkg/config" // Ajusta tu path
    "github.com/google/uuid"
)

type Claims struct {
    UserID   uuid.UUID `json:"user_id"`
    Username string    `json:"username"`
    jwt.RegisteredClaims
}

func GenerateAccessToken(userID uuid.UUID, username string) (string, error) {
    expirationTime := time.Now().Add(config.AppConfig.JWTExpiration)
    claims := &Claims{
        UserID:   userID,
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "auth-service", // O tu identificador de emisor
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(config.AppConfig.JWTSecretKey))
}

// Podrías añadir GenerateRefreshToken de forma similar con mayor expiración

func ValidateToken(tokenString string) (*Claims, error) {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(config.AppConfig.JWTSecretKey), nil
    })

    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }
    return claims, nil
}