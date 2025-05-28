// services/task-orchestrator-service/internal/services/auth_client.go
package services // <-- Declaración del paquete

import (
    "encoding/json"
    "fmt"
    "io"
    "log" // Añadido para debug
    "net/http"
    "time"

    // Asegúrate de que esta ruta sea correcta para tu estructura y el nombre de tu módulo
    "github.com/guildmember145/task-orchestrator-service/pkg/config"
)

type ValidateTokenResponse struct {
    Valid    bool   `json:"valid"`
    UserID   string `json:"user_id,omitempty"`
    Username string `json:"username,omitempty"`
    Error    string `json:"error,omitempty"`
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

func ValidateTokenWithAuthService(tokenString string) (*ValidateTokenResponse, error) {
    validateURL := fmt.Sprintf("%s/validate_token", config.AppConfig.AuthServiceBaseURL)
    log.Printf("Attempting to validate token with auth service at URL: %s", validateURL) // Log para debug

    req, err := http.NewRequest("POST", validateURL, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request to auth service: %w", err)
    }
    req.Header.Set("Authorization", "Bearer "+tokenString)
    req.Header.Set("Content-Type", "application/json")

    resp, err := httpClient.Do(req)
    if err != nil {
        log.Printf("Error calling auth service: %v", err) // Log para debug
        return nil, fmt.Errorf("failed to call auth service: %w", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read auth service response body: %w", err)
    }

    log.Printf("Auth service response status: %d, body: %s", resp.StatusCode, string(body)) // Log para debug

    if resp.StatusCode != http.StatusOK {
        var errorResp ValidateTokenResponse
        if json.Unmarshal(body, &errorResp) == nil && errorResp.Error != "" {
             return nil, fmt.Errorf("auth service validation failed (status %d): %s", resp.StatusCode, errorResp.Error)
        }
        return nil, fmt.Errorf("auth service validation failed with status %d: %s", resp.StatusCode, string(body))
    }

    var validationResponse ValidateTokenResponse
    if err := json.Unmarshal(body, &validationResponse); err != nil {
        return nil, fmt.Errorf("failed to decode auth service response: %w. Body: %s", err, string(body))
    }

    if !validationResponse.Valid {
         // No es necesariamente un error de conexión, sino que el token no es válido según el auth-service
        return &validationResponse, fmt.Errorf("token is not valid according to auth service: %s", validationResponse.Error)
    }

    return &validationResponse, nil
}