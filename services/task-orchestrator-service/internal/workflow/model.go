package workflow

import (
	"fmt"
    "time"
    "github.com/google/uuid"
)

// TriggerType define los tipos de disparadores que podemos tener
type TriggerType string

const (
    TriggerTypeSchedule TriggerType = "schedule"
    TriggerTypeWebhook  TriggerType = "webhook"
    // Podríamos añadir más tipos en el futuro (ej. TriggerTypeEvent)
)

// ActionType define los tipos de acciones que podemos realizar
type ActionType string

const (
    ActionTypeLogMessage      ActionType = "log_message"       // Simplemente escribe un mensaje en el log
    ActionTypeHTTPEndpoint    ActionType = "http_endpoint"     // Llama a una URL externa
    // Podríamos añadir más (ej. ActionTypeSendEmail)
)

// TriggerDefinition contiene la configuración para un disparador
type TriggerDefinition struct {
    Type   TriggerType            `json:"type" validate:"required,oneof=schedule webhook"`
    Config map[string]interface{} `json:"config"` // Ej. para schedule: {"cron": "0 * * * *"}, para webhook: {} (podría generar URL)
}

// ActionDefinition contiene la configuración para una acción
type ActionDefinition struct {
    Type        ActionType             `json:"type" validate:"required,oneof=log_message http_endpoint"`
    Name        string                 `json:"name" validate:"required"` // Un nombre descriptivo para el paso de acción
    Config      map[string]interface{} `json:"config"`                   // Ej. para log_message: {"message": "Tarea ejecutada"}, para http_endpoint: {"url": "...", "method": "GET/POST", "body": "..."}
    DependsOn   []string               `json:"depends_on,omitempty"`     // Nombres de otras acciones de las que depende (para flujos secuenciales/paralelos básicos)
}

// Workflow es la estructura principal para nuestra definición de tarea
type Workflow struct {
    ID              uuid.UUID          `json:"id"`
    UserID          string             `json:"user_id"` // Vendrá del token JWT
    Name            string             `json:"name" validate:"required,min=3,max=100"`
    Description     string             `json:"description,omitempty"`
    Trigger         TriggerDefinition  `json:"trigger" validate:"required"`
    Actions         []ActionDefinition `json:"actions" validate:"required,min=1,dive"` // 'dive' valida cada elemento del slice
    IsEnabled       bool               `json:"is_enabled"`
    CreatedAt       time.Time          `json:"created_at"`
    UpdatedAt       time.Time          `json:"updated_at"`
    LastRunAt       *time.Time         `json:"last_run_at,omitempty"` // Puntero para que pueda ser nulo
    NextRunAt       *time.Time         `json:"next_run_at,omitempty"` // Para triggers de schedule
}

// --- Almacenamiento en Memoria (Temporal) ---
// Usaremos un mapa para guardar los workflows, la clave será el UserID
// y el valor será otro mapa de WorkflowID a Workflow.
// Esto es para desarrollo; lo reemplazaremos con una BD.
var workflowStore = make(map[string]map[uuid.UUID]*Workflow)

// Funciones de ayuda para el store en memoria (ejemplos)
func GetWorkflowsByUserID(userID string) ([]*Workflow, error) {
    userWorkflows, ok := workflowStore[userID]
    if !ok {
        return []*Workflow{}, nil // No hay workflows para este usuario aún
    }

    list := make([]*Workflow, 0, len(userWorkflows))
    for _, wf := range userWorkflows {
        list = append(list, wf)
    }
    return list, nil
}

func GetWorkflowByID(userID string, workflowID uuid.UUID) (*Workflow, bool) {
    userWorkflows, ok := workflowStore[userID]
    if !ok {
        return nil, false
    }
    wf, found := userWorkflows[workflowID]
    return wf, found
}

func SaveWorkflow(wf *Workflow) error {
    if wf.ID == uuid.Nil {
        wf.ID = uuid.New() // Asignar nuevo ID si es nuevo
    }
    if _, ok := workflowStore[wf.UserID]; !ok {
        workflowStore[wf.UserID] = make(map[uuid.UUID]*Workflow)
    }
    // Simulación de validación simple
    if wf.Name == "" {
        return fmt.Errorf("workflow name cannot be empty")
    }
    workflowStore[wf.UserID][wf.ID] = wf
    return nil
}

 func DeleteWorkflow(userID string, workflowID uuid.UUID) bool {
    userWorkflows, ok := workflowStore[userID]
    if !ok {
        return false // No hay workflows para este usuario
    }
    if _, found := userWorkflows[workflowID]; found {
        delete(userWorkflows, workflowID)
        return true
    }
    return false
}


type InMemoryWorkflowStore struct{}

func (s *InMemoryWorkflowStore) GetAllEnabledScheduledWorkflows() ([]*Workflow, error) {
    allScheduled := []*Workflow{}
    // Iteramos sobre todos los usuarios y sus workflows
    // Esto es simple para un store en memoria; una BD lo haría más eficientemente.
    for _, userWorkflows := range workflowStore { // workflowStore es nuestro mapa global en memoria
        for _, wf := range userWorkflows {
            if wf.IsEnabled && wf.Trigger.Type == TriggerTypeSchedule {
                allScheduled = append(allScheduled, wf)
            }
        }
    }
    return allScheduled, nil
}