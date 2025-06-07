// services/task-orchestrator-service/internal/handlers/workflow_handler.go
package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/guildmember145/task-orchestrator-service/internal/scheduler"
	"github.com/guildmember145/task-orchestrator-service/internal/workflow"
	"github.com/guildmember145/task-orchestrator-service/pkg/transport"
)

var validate = validator.New()

// WorkflowHandler contiene las dependencias para los manejadores de workflows.
type WorkflowHandler struct {
	Store     workflow.Store
	Scheduler *scheduler.Scheduler
}

// NewWorkflowHandler crea una nueva instancia de WorkflowHandler.
func NewWorkflowHandler(store workflow.Store, scheduler *scheduler.Scheduler) *WorkflowHandler {
	return &WorkflowHandler{Store: store, Scheduler: scheduler}
}

// CreateWorkflowHandler ahora es un método de WorkflowHandler.
func (h *WorkflowHandler) CreateWorkflowHandler(c *gin.Context) {
	var req transport.CreateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// Lógica de validación en etapas
	if err := validate.Struct(req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Top-level validation failed: " + err.Error()}); return }
	if err := validate.Struct(req.Trigger); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Trigger validation failed: " + err.Error()}); return }
	if len(req.Actions) > 0 {
		for i, action := range req.Actions {
			if err := validate.Struct(action); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Action %d validation failed: %s", i, err.Error())}); return }
		}
	}

	userIDClaim, _ := c.Get("userID")
    userID, err := uuid.Parse(userIDClaim.(string))
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid User ID format in token"}); return
    }

	newWorkflow := &workflow.Workflow{
		ID:          uuid.New(),
		UserID:      userID.String(),
		Name:        req.Name,
		Description: req.Description,
		Trigger:     req.Trigger,
		Actions:     req.Actions,
		IsEnabled:   req.IsEnabled,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if err := h.Store.SaveWorkflow(newWorkflow); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save workflow"})
		return
	}

	h.Scheduler.ReloadAndRescheduleWorkflows()
	c.JSON(http.StatusCreated, newWorkflow)
}

// GetWorkflowsHandler ahora es un método de WorkflowHandler.
func (h *WorkflowHandler) GetWorkflowsHandler(c *gin.Context) {
    userIDClaim, _ := c.Get("userID")
    if userIDClaim == nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"}); return
    }

    userWorkflows, err := h.Store.GetWorkflowsByUserID(userIDClaim.(string))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve workflows"})
        return
    }

    // --- INICIO DE LA CORRECCIÓN ---
    // Si el slice de workflows es nulo (porque el usuario no tiene ninguno),
    // nos aseguramos de devolver un array JSON vacío [] en lugar de 'null'.
    if userWorkflows == nil {
        c.JSON(http.StatusOK, []workflow.Workflow{})
        return
    }
    // --- FIN DE LA CORRECCIÓN ---

    c.JSON(http.StatusOK, userWorkflows)
}


// --- INICIO DE LOS NUEVOS MÉTODOS CORREGIDOS ---

// GetWorkflowByIDHandler ahora es un método de WorkflowHandler.
func (h *WorkflowHandler) GetWorkflowByIDHandler(c *gin.Context) {
	userIDClaim, _ := c.Get("userID")
	workflowIDParam := c.Param("workflow_id")
	workflowID, err := uuid.Parse(workflowIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID format"})
		return
	}

	wf, found := h.Store.GetWorkflowByID(userIDClaim.(string), workflowID)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found or access denied"})
		return
	}
	c.JSON(http.StatusOK, wf)
}

// UpdateWorkflowHandler ahora es un método de WorkflowHandler.
func (h *WorkflowHandler) UpdateWorkflowHandler(c *gin.Context) {
	userIDClaim, _ := c.Get("userID")
	workflowIDParam := c.Param("workflow_id")
	workflowID, err := uuid.Parse(workflowIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID format"})
		return
	}

	existingWorkflow, found := h.Store.GetWorkflowByID(userIDClaim.(string), workflowID)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found or access denied"})
		return
	}

	var req transport.CreateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
    // Puedes añadir la validación en etapas aquí también si lo deseas...

	existingWorkflow.Name = req.Name
	existingWorkflow.Description = req.Description
	existingWorkflow.Trigger = req.Trigger
	existingWorkflow.Actions = req.Actions
	existingWorkflow.IsEnabled = req.IsEnabled
	existingWorkflow.UpdatedAt = time.Now().UTC()

	if err := h.Store.SaveWorkflow(existingWorkflow); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update workflow"})
		return
	}

	h.Scheduler.ReloadAndRescheduleWorkflows()
	c.JSON(http.StatusOK, existingWorkflow)
}

// DeleteWorkflowHandler ahora es un método de WorkflowHandler.
func (h *WorkflowHandler) DeleteWorkflowHandler(c *gin.Context) {
	userIDClaim, _ := c.Get("userID")
	workflowIDParam := c.Param("workflow_id")
	workflowID, err := uuid.Parse(workflowIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID format"})
		return
	}

	if !h.Store.DeleteWorkflow(userIDClaim.(string), workflowID) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found or access denied"})
		return
	}

	h.Scheduler.ReloadAndRescheduleWorkflows()
	c.Status(http.StatusNoContent)
}
// --- FIN DE LOS NUEVOS MÉTODOS CORREGIDOS ---