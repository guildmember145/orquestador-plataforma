// services/task-orchestrator-service/internal/engine/executor.go
package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/guildmember145/task-orchestrator-service/internal/workflow"
)

var httpClient = &http.Client{Timeout: 30 * time.Second}

// LogEntry define la estructura de una línea de log individual que guardaremos como JSON.
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Status    string    `json:"status"` // "INFO", "ERROR", "ACTION_OUTPUT"
}

// ExecuteWorkflow ahora acepta el 'Store' para poder guardar los resultados.
func ExecuteWorkflow(wf workflow.Workflow, store workflow.Store) {
	log.Printf("ENGINE: >>> Starting execution for Workflow ID %s, Name: '%s' <<<", wf.ID, wf.Name)

	executionLogs := []LogEntry{}
	addLog := func(message string, status string) {
		entry := LogEntry{Timestamp: time.Now().UTC(), Message: message, Status: status}
		executionLogs = append(executionLogs, entry)
		log.Printf("[%s] %s", status, message)
	}

	// 1. Crear el registro de ejecución inicial en la base de datos
	initialLog, _ := json.Marshal([]LogEntry{
		{Timestamp: time.Now().UTC(), Message: fmt.Sprintf("Starting execution for Workflow '%s'", wf.Name), Status: "INFO"},
	})
	execution := &workflow.ExecutionLog{
		ID:          uuid.New(),
		WorkflowID:  wf.ID,
		UserID:      wf.UserID,
		Status:      "running",
		TriggeredAt: time.Now().UTC(),
		Logs:        initialLog,
	}
	store.CreateExecution(execution)

	// 2. Defer se asegura de que el estado final se guarde siempre
	defer func() {
		if r := recover(); r != nil {
			addLog(fmt.Sprintf("Panic recovered during execution: %v", r), "ERROR")
			execution.Status = "failed"
		}
		now := time.Now().UTC()
		execution.CompletedAt = &now
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
		var actionErr error

		switch action.Type {
		case workflow.ActionTypeLogMessage:
			msg, ok := action.Config["message"].(string)
			if !ok {
				actionErr = fmt.Errorf("action '%s': 'message' not found in config or is not a string", action.Name)
			} else {
				addLog(msg, "ACTION_OUTPUT")
			}

		case workflow.ActionTypeHTTPEndpoint:
			// --- INICIO DE LA LÓGICA COMPLETA PARA HTTP ---
			url, urlOk := action.Config["url"].(string)
			if !urlOk || url == "" {
				actionErr = fmt.Errorf("action '%s': 'url' is required for http_endpoint", action.Name)
				break
			}

			method, _ := action.Config["method"].(string)
			if method == "" {
				method = http.MethodGet
			}
			method = strings.ToUpper(method)

			var reqBody io.Reader
			if bodyData, exists := action.Config["body"]; exists && bodyData != nil {
				if bodyStr, isStr := bodyData.(string); isStr {
					reqBody = strings.NewReader(bodyStr)
				} else {
					jsonBody, err := json.Marshal(bodyData)
					if err != nil {
						actionErr = fmt.Errorf("failed to marshal 'body' to JSON for action '%s': %w", action.Name, err)
						break
					}
					reqBody = bytes.NewBuffer(jsonBody)
				}
			}

			req, err := http.NewRequestWithContext(context.Background(), method, url, reqBody)
			if err != nil {
				actionErr = fmt.Errorf("failed to create HTTP request for action '%s': %w", action.Name, err)
				break
			}

			if headersData, exists := action.Config["headers"].(map[string]interface{}); exists {
				for key, val := range headersData {
					if valStr, isStr := val.(string); isStr {
						req.Header.Set(key, valStr)
					}
				}
			}
			if reqBody != nil && req.Header.Get("Content-Type") == "" {
				if _, isMap := action.Config["body"].(map[string]interface{}); isMap {
					req.Header.Set("Content-Type", "application/json")
				}
			}

			resp, err := httpClient.Do(req)
			if err != nil {
				actionErr = fmt.Errorf("failed to execute HTTP request for action '%s': %w", action.Name, err)
				break
			}
			defer resp.Body.Close()

			respBodyBytes, _ := io.ReadAll(resp.Body)
			addLog(fmt.Sprintf("HTTP Status: %s, Response Body: %.500s", resp.Status, string(respBodyBytes)), "ACTION_OUTPUT")

			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				actionErr = fmt.Errorf("HTTP request for action '%s' returned non-2xx status: %s", action.Name, resp.Status)
			}
			// --- FIN DE LA LÓGICA COMPLETA PARA HTTP ---

		default:
			actionErr = fmt.Errorf("unknown action type '%s'", action.Type)
		}

		if actionErr != nil {
			addLog(fmt.Sprintf("Error processing Action '%s': %v", action.Name, actionErr), "ERROR")
			overallSuccess = false
		}
	}

	if overallSuccess {
		execution.Status = "completed"
	} else {
		execution.Status = "failed"
	}
}