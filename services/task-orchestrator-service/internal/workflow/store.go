// services/task-orchestrator-service/internal/workflow/store.go
package workflow

import "github.com/google/uuid"

// Store define la interfaz para las operaciones de almacenamiento de workflows.
type Store interface {
    SaveWorkflow(wf *Workflow) error
    GetWorkflowsByUserID(userID string) ([]*Workflow, error)
    GetWorkflowByID(userID string, workflowID uuid.UUID) (*Workflow, bool)
    DeleteWorkflow(userID string, workflowID uuid.UUID) bool
    GetAllEnabledScheduledWorkflows() ([]*Workflow, error)
    CreateExecution(exec *ExecutionLog) error
    UpdateExecution(exec *ExecutionLog) error
}