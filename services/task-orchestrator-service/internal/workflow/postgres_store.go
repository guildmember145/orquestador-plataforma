// services/task-orchestrator-service/internal/workflow/postgres_store.go
package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("workflow not found")

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