// services/task-orchestrator-service/internal/workflow/postgres_store.go
package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time" // <-- Asegúrate de que time esté importado

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("workflow not found")

// ExecutionLog represents a workflow execution log entry.
// CORREGIDO: Se han cambiado los tipos de TriggeredAt y CompletedAt.
type ExecutionLog struct {
	ID          uuid.UUID       `json:"id"`
	WorkflowID  uuid.UUID       `json:"workflow_id"`
	UserID      string          `json:"user_id"`
	Status      string          `json:"status"`
	TriggeredAt time.Time       `json:"triggered_at"`         // <-- CORREGIDO de string a time.Time
	CompletedAt *time.Time      `json:"completed_at,omitempty"` // <-- CORREGIDO de string a *time.Time
	Logs        json.RawMessage `json:"logs"`
}

type PostgresWorkflowStore struct {
	DB *pgxpool.Pool
}

func NewPostgresWorkflowStore(db *pgxpool.Pool) *PostgresWorkflowStore {
	return &PostgresWorkflowStore{DB: db}
}

// SaveWorkflow inserta o actualiza un workflow en la base de datos.
func (s *PostgresWorkflowStore) SaveWorkflow(wf *Workflow) error {
	// Convertimos los campos de struct/slice a JSON para guardarlos en columnas JSONB.
	triggerJSON, err := json.Marshal(wf.Trigger)
	if err != nil {
		return fmt.Errorf("failed to marshal trigger: %w", err)
	}

	actionsJSON, err := json.Marshal(wf.Actions)
	if err != nil {
		return fmt.Errorf("failed to marshal actions: %w", err)
	}

	query := `
        INSERT INTO workflows (id, user_id, name, description, trigger, actions, is_enabled, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        ON CONFLICT (id) DO UPDATE SET
            name = EXCLUDED.name,
            description = EXCLUDED.description,
            trigger = EXCLUDED.trigger,
            actions = EXCLUDED.actions,
            is_enabled = EXCLUDED.is_enabled,
            updated_at = EXCLUDED.updated_at;
    `
	_, err = s.DB.Exec(context.Background(), query,
		wf.ID, wf.UserID, wf.Name, wf.Description, triggerJSON, actionsJSON, wf.IsEnabled, wf.CreatedAt, wf.UpdatedAt)

	if err != nil {
		log.Printf("Error saving workflow to database: %v", err)
		return fmt.Errorf("could not save workflow: %w", err)
	}
	return nil
}

// scanWorkflow es una función de ayuda para escanear una fila de la BD a un struct Workflow.
func scanWorkflow(row pgx.Row) (*Workflow, error) {
	var wf Workflow
	var triggerJSON, actionsJSON []byte

	// Asumiendo que las columnas last_run_at y next_run_at no están en esta consulta.
	// Si estuvieran, necesitarías añadirlas aquí.
	err := row.Scan(
		&wf.ID, &wf.UserID, &wf.Name, &wf.Description, &triggerJSON, &actionsJSON,
		&wf.IsEnabled, &wf.CreatedAt, &wf.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(triggerJSON, &wf.Trigger); err != nil {
		return nil, fmt.Errorf("failed to unmarshal trigger: %w", err)
	}
	if err := json.Unmarshal(actionsJSON, &wf.Actions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal actions: %w", err)
	}
	return &wf, nil
}

func (s *PostgresWorkflowStore) GetWorkflowsByUserID(userID string) ([]*Workflow, error) {
	query := `SELECT id, user_id, name, description, trigger, actions, is_enabled, created_at, updated_at 
              FROM workflows WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := s.DB.Query(context.Background(), query, userID)
	if err != nil {
		return nil, fmt.Errorf("database query failed: %w", err)
	}
	defer rows.Close()

	var workflows []*Workflow
	for rows.Next() {
		wf, err := scanWorkflow(rows)
		if err != nil {
			return nil, err
		}
		workflows = append(workflows, wf)
	}
	return workflows, nil
}

func (s *PostgresWorkflowStore) GetWorkflowByID(userID string, workflowID uuid.UUID) (*Workflow, bool) {
	query := `SELECT id, user_id, name, description, trigger, actions, is_enabled, created_at, updated_at 
              FROM workflows WHERE id = $1 AND user_id = $2`

	row := s.DB.QueryRow(context.Background(), query, workflowID, userID)
	wf, err := scanWorkflow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, false
		}
		log.Printf("Error scanning workflow by ID: %v", err)
		return nil, false
	}
	return wf, true
}

func (s *PostgresWorkflowStore) DeleteWorkflow(userID string, workflowID uuid.UUID) bool {
	query := `DELETE FROM workflows WHERE id = $1 AND user_id = $2`
	cmdTag, err := s.DB.Exec(context.Background(), query, workflowID, userID)
	if err != nil {
		log.Printf("Error deleting workflow: %v", err)
		return false
	}
	return cmdTag.RowsAffected() > 0
}

func (s *PostgresWorkflowStore) GetAllEnabledScheduledWorkflows() ([]*Workflow, error) {
	// Esta consulta busca en el campo JSONB del trigger.
	query := `SELECT id, user_id, name, description, trigger, actions, is_enabled, created_at, updated_at 
              FROM workflows WHERE is_enabled = TRUE AND trigger->>'type' = 'schedule'`

	rows, err := s.DB.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("database query for scheduler failed: %w", err)
	}
	defer rows.Close()

	var workflows []*Workflow
	for rows.Next() {
		wf, err := scanWorkflow(rows)
		if err != nil {
			return nil, err
		}
		workflows = append(workflows, wf)
	}
	return workflows, nil
}

// CreateExecution ahora recibe un `ExecutionLog` con los tipos de fecha correctos.
func (s *PostgresWorkflowStore) CreateExecution(exec *ExecutionLog) error {
	query := `
        INSERT INTO workflow_executions (id, workflow_id, user_id, status, triggered_at, logs)
        VALUES ($1, $2, $3, $4, $5, $6)`

	// El campo `exec.Logs` ya viene como json.RawMessage, por lo que no necesita Marshal aquí
	// si se inicializa como json.RawMessage("[]"). Si lo inicializas como un slice de LogEntry,
	// entonces sí necesitarías json.Marshal. El código en `engine` ya hace esto,
	// pero en la creación inicial `exec.Logs` es json.RawMessage("[]"), así que se puede pasar directo.
	// Sin embargo, para consistencia, es más seguro marshallear el `exec.Logs` que es un slice de structs.
	// Pero el `exec` que llega aquí ya tiene los logs como `json.RawMessage("[]")`
	// El motor de ejecución es el que debe hacer el marshal final.
	// Para la inserción inicial, el valor de logs es simple.
	_, err := s.DB.Exec(context.Background(), query, exec.ID, exec.WorkflowID, exec.UserID, exec.Status, exec.TriggeredAt, exec.Logs)
	if err != nil {
		return fmt.Errorf("failed to insert execution record: %w", err)
	}
	return nil
}

// UpdateExecution actualiza un registro de ejecución al finalizar un workflow.
func (s *PostgresWorkflowStore) UpdateExecution(exec *ExecutionLog) error {
	query := `
        UPDATE workflow_executions 
        SET status = $1, completed_at = $2, logs = $3
        WHERE id = $4`
	
	// En la actualización, `exec.Logs` sí contiene los logs completos que han sido
	// convertidos a json.RawMessage por el motor de ejecución.
	_, err := s.DB.Exec(context.Background(), query, exec.Status, exec.CompletedAt, exec.Logs, exec.ID)
	if err != nil {
		return fmt.Errorf("failed to update execution record: %w", err)
	}
	return nil
}