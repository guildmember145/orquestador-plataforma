// services/task-orchestrator-service/internal/handlers/workflow_handler.go
package handlers

import (
	"fmt"
	"net/http"
	"time"
	"log"
	"strings"

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

// CreateWorkflowHandler crea un nuevo workflow.
func (h *WorkflowHandler) CreateWorkflowHandler(c *gin.Context) {
	var req transport.CreateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// Lógica de validación en etapas (ya la tenías funcionando)
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
		UserID:      userID.String(), // Convertimos UUID a string
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

// GetWorkflowsHandler lista los workflows del usuario.
func (h *WorkflowHandler) GetWorkflowsHandler(c *gin.Context) {
    userIDClaim, _ := c.Get("userID")
    userWorkflows, err := h.Store.GetWorkflowsByUserID(userIDClaim.(string))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve workflows"})
        return
    }
    if userWorkflows == nil {
        c.JSON(http.StatusOK, []workflow.Workflow{})
        return
    }
    c.JSON(http.StatusOK, userWorkflows)
}

// GetWorkflowByIDHandler obtiene un workflow específico.
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

// UpdateWorkflowHandler actualiza un workflow existente.
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

// DeleteWorkflowHandler elimina un workflow.
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

// GetWorkflowExecutionsHandler obtiene el historial de ejecuciones.
func (h *WorkflowHandler) GetWorkflowExecutionsHandler(c *gin.Context) {
    userIDClaim, exists := c.Get("userID")
    if !exists {
        log.Printf("ERROR: userID not found in context")
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    userID := userIDClaim.(string)
    workflowIDParam := c.Param("workflow_id")
    
    log.Printf("INFO: Getting executions for workflow %s, user %s", workflowIDParam, userID)
    
    workflowID, err := uuid.Parse(workflowIDParam)
    if err != nil {
        log.Printf("ERROR: Invalid workflow ID format: %s", workflowIDParam)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID format"})
        return
    }

    executions, err := h.Store.GetExecutionsByWorkflowID(userID, workflowID)
    if err != nil {
        log.Printf("ERROR: Failed to get executions for workflow %s, user %s: %v", workflowID, userID, err)
        
        // Si el error es por acceso denegado, devolver 403
        if strings.Contains(err.Error(), "workflow not found or access denied") {
            c.JSON(http.StatusForbidden, gin.H{"error": "Access to workflow denied"})
            return
        }
        
        // Para otros errores de base de datos
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to retrieve executions",
            "details": err.Error(), // Solo para debugging, remover en producción
        })
        return
    }

    log.Printf("SUCCESS: Retrieved %d executions for workflow %s", len(executions), workflowID)

    // Si no hay ejecuciones, devolver array vacío
    if executions == nil {
        executions = []*workflow.ExecutionLog{}
    }

    c.JSON(http.StatusOK, executions)
}