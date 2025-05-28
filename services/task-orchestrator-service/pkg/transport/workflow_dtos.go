package transport

import (
    "github.com/guildmember145/task-orchestrator-service/internal/workflow"
)

type CreateWorkflowRequest struct {
    Name        string                           `json:"name" validate:"required,min=3,max=100"`
    Description string                           `json:"description,omitempty"`
    Trigger     workflow.TriggerDefinition       `json:"trigger" validate:"required,dive"`
    Actions     []workflow.ActionDefinition    `json:"actions" validate:"required,min=1,dive"`
    IsEnabled   bool                             `json:"is_enabled"`
}

// Podrías añadir UpdateWorkflowRequest si los campos son diferentes