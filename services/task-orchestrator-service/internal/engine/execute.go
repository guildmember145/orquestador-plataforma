// services/task-orchestrator-service/internal/engine/executor.go
package engine

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"


	"github.com/google/uuid"
	"github.com/guildmember145/task-orchestrator-service/internal/workflow"
)

var httpClient = &http.Client{Timeout: 30 * time.Second}

// LogEntry define la estructura de una línea de log individual.
type LogEntry struct {
    Timestamp time.Time `json:"timestamp"`
    Message   string    `json:"message"`
    Status    string    `json:"status"` // "INFO", "ERROR", "ACTION_OUTPUT"
}

// ExecuteWorkflow ahora acepta el 'Store' para poder guardar los resultados.
func ExecuteWorkflow(wf workflow.Workflow, store workflow.Store) {
	log.Printf("ENGINE: >>> Starting execution for Workflow ID %s, Name: '%s' <<<", wf.ID, wf.Name)

	// 1. Crear el registro de ejecución en la base de datos
	execution := &workflow.ExecutionLog{
		ID:          uuid.New(),
		WorkflowID:  wf.ID,
		UserID:      wf.UserID,
		Status:      "running",
		TriggeredAt: time.Now().UTC(),
		Logs:        json.RawMessage("[]"), // Inicializar como un array JSON vacío
	}

    // Creamos un slice de Go para ir guardando los logs
    executionLogs := []LogEntry{}
    addLog := func(message string, status string) {
        entry := LogEntry{Timestamp: time.Now().UTC(), Message: message, Status: status}
        executionLogs = append(executionLogs, entry)
        log.Printf("[%s] %s", status, message) // Seguimos imprimiendo en la consola del contenedor por ahora
    }

    addLog(fmt.Sprintf("Starting execution for Workflow '%s'", wf.Name), "INFO")
    store.CreateExecution(execution) // Guardar el registro inicial

    // Defer se asegura de que el estado final se guarde incluso si hay un pánico.
    defer func() {
        if r := recover(); r != nil {
            addLog(fmt.Sprintf("Panic recovered during execution: %v", r), "ERROR")
            execution.Status = "failed"
        }

        now := time.Now().UTC()
		execution.CompletedAt = &now

        // Convertimos los logs a JSON y los guardamos
        logsJSON, err := json.Marshal(executionLogs)
        if err != nil {
            log.Printf("CRITICAL: Failed to marshal execution logs for workflow %s: %v", wf.ID, err)
        } else {
            execution.Logs = logsJSON
        }

        store.UpdateExecution(execution)
		log.Printf("ENGINE: >>> Finished execution for Workflow ID %s, Status: %s <<<", wf.ID, execution.Status)
    }()


	if len(wf.Actions) == 0 {
		addLog("Workflow has no actions to execute.", "WARNING")
        execution.Status = "completed"
		return
	}

    overallSuccess := true
	for i, action := range wf.Actions {
		addLog(fmt.Sprintf("--- Executing Action %d/%d: Name: '%s', Type: '%s' ---", i+1, len(wf.Actions), action.Name, action.Type), "INFO")

        // Lógica para cada tipo de acción...
        // (La lógica interna de log_message y http_endpoint que ya tenías va aquí,
        // pero en lugar de log.Printf, llamamos a nuestra función 'addLog')

        // Ejemplo para log_message:
        if action.Type == workflow.ActionTypeLogMessage {
            msg, ok := action.Config["message"].(string)
            if !ok {
                addLog(fmt.Sprintf("Action '%s' failed: 'message' not found in config", action.Name), "ERROR")
                overallSuccess = false
            } else {
                addLog(fmt.Sprintf("%s", msg), "ACTION_OUTPUT")
            }
        }

        // Ejemplo para http_endpoint (versión simplificada):
        if action.Type == workflow.ActionTypeHTTPEndpoint {
             // Aquí iría tu lógica completa para la acción HTTP...
             // y al final, en lugar de log.Printf, harías:
             // addLog(fmt.Sprintf("HTTP Status: %s, Body: %.100s", resp.Status, body), "ACTION_OUTPUT")
        }
	}

    if overallSuccess {
        execution.Status = "completed"
    } else {
        execution.Status = "failed"
    }
}