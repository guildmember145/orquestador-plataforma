// services/task-orchestrator-service/internal/handlers/workflow_handler.go
package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/guildmember145/task-orchestrator-service/internal/scheduler" // <--- AÑADE ESTA IMPORTACIÓN
	"github.com/guildmember145/task-orchestrator-service/internal/workflow"
	"github.com/guildmember145/task-orchestrator-service/pkg/transport"
)

var validate = validator.New()

// CreateWorkflowHandler crea un nuevo workflow
// MODIFICADO: Acepta appScheduler como parámetro
func CreateWorkflowHandler(c *gin.Context, appScheduler *scheduler.Scheduler) {
	var req transport.CreateWorkflowRequest
	// ... (tu validación de JSON y structs que ya funciona) ...
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
        return
    }
    log.Println("Validating CreateWorkflowRequest structure (top level)...")
    if err := validate.Struct(req); err != nil {
        log.Printf("Top-level validation failed: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Top-level validation failed: " + err.Error()})
        return
    }
    log.Println("Top-level validation passed.")
    log.Println("Validating TriggerDefinition structure...")
    if err := validate.Struct(req.Trigger); err != nil {
        log.Printf("Trigger validation failed: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Trigger validation failed: " + err.Error()})
        return
    }
    log.Println("Trigger validation passed.")
    log.Println("Validating ActionDefinition structures...")
    if len(req.Actions) > 0 {
        for i, action := range req.Actions {
            if err := validate.Struct(action); err != nil {
                log.Printf("Action %d validation failed: %v", i, err)
                c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Action %d validation failed: %s", i, err.Error())})
                return
            }
        }
    }
    log.Println("Actions validation passed.")

	userID, _ := c.Get("userID")
	if userID == nil || userID.(string) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found or invalid"})
		return
	}

	newWorkflow := &workflow.Workflow{
		ID:          uuid.New(),
		UserID:      userID.(string),
		Name:        req.Name,
		Description: req.Description,
		Trigger:     req.Trigger,
		Actions:     req.Actions,
		IsEnabled:   req.IsEnabled,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := workflow.SaveWorkflow(newWorkflow); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save workflow: " + err.Error()})
		return
	}

	appScheduler.ReloadAndRescheduleWorkflows() // <--- NUEVA LÍNEA: Notificar al scheduler
	c.JSON(http.StatusCreated, newWorkflow)
}

// GetWorkflowsHandler (sin cambios necesarios para el scheduler por ahora)
func GetWorkflowsHandler(c *gin.Context) {
    // ... (código existente) ...
	userID, _ := c.Get("userID")
	if userID == nil || userID.(string) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found or invalid"})
		return
	}
	userWorkflows, err := workflow.GetWorkflowsByUserID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve workflows: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, userWorkflows)
}


// GetWorkflowByIDHandler (sin cambios necesarios para el scheduler por ahora)
func GetWorkflowByIDHandler(c *gin.Context) {
    // ... (código existente) ...
	userID, _ := c.Get("userID")
	if userID == nil || userID.(string) == "" {
		 c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found or invalid"})
		return
	}
	workflowIDParam := c.Param("workflow_id")
	workflowID, err := uuid.Parse(workflowIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID format"})
		return
	}
	wf, found := workflow.GetWorkflowByID(userID.(string), workflowID)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found or access denied"})
		return
	}
	c.JSON(http.StatusOK, wf)
}

// UpdateWorkflowHandler actualiza un workflow existente
// MODIFICADO: Acepta appScheduler como parámetro
func UpdateWorkflowHandler(c *gin.Context, appScheduler *scheduler.Scheduler) {
    // ... (tu lógica de obtener userID, workflowID, verificar existencia, bind & validate JSON req) ...
	userID, _ := c.Get("userID")
	if userID == nil || userID.(string) == "" {
		 c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found or invalid"})
		return
	}
	workflowIDParam := c.Param("workflow_id")
	workflowID, err := uuid.Parse(workflowIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID format"})
		return
	}
	existingWorkflow, found := workflow.GetWorkflowByID(userID.(string), workflowID)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found or access denied"})
		return
	}
	var req transport.CreateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}
	// Reutilizar la lógica de validación escalonada
    log.Println("Validating UpdateWorkflowRequest (top level)...")
    if err_val := validate.Struct(req); err_val != nil {
        log.Printf("Update top-level validation failed: %v", err_val)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Update top-level validation failed: " + err_val.Error()})
        return
    }
    log.Println("Update top-level validation passed.")
    log.Println("Validating Update TriggerDefinition structure...")
    if err_val := validate.Struct(req.Trigger); err_val != nil {
        log.Printf("Update Trigger validation failed: %v", err_val)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Update Trigger validation failed: " + err_val.Error()})
        return
    }
    log.Println("Update Trigger validation passed.")
    log.Println("Validating Update ActionDefinition structures...")
    if len(req.Actions) > 0 {
        for i, action := range req.Actions {
            if err_val := validate.Struct(action); err_val != nil {
                log.Printf("Update Action %d validation failed: %v", i, err_val)
                c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Update Action %d validation failed: %s", i, err_val.Error())})
                return
            }
        }
    }
    log.Println("Update Actions validation passed.")


	existingWorkflow.Name = req.Name
	existingWorkflow.Description = req.Description
	existingWorkflow.Trigger = req.Trigger
	existingWorkflow.Actions = req.Actions
	existingWorkflow.IsEnabled = req.IsEnabled
	existingWorkflow.UpdatedAt = time.Now()

	if err_save := workflow.SaveWorkflow(existingWorkflow); err_save != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update workflow: " + err_save.Error()})
		return
	}
	appScheduler.ReloadAndRescheduleWorkflows() // <--- NUEVA LÍNEA: Notificar al scheduler
	c.JSON(http.StatusOK, existingWorkflow)
}

// DeleteWorkflowHandler elimina un workflow
// MODIFICADO: Acepta appScheduler como parámetro
func DeleteWorkflowHandler(c *gin.Context, appScheduler *scheduler.Scheduler) {
    // ... (tu lógica de obtener userID, workflowID) ...
	userID, _ := c.Get("userID")
	 if userID == nil || userID.(string) == "" {
		 c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found or invalid"})
		return
	}
	workflowIDParam := c.Param("workflow_id")
	workflowID, err := uuid.Parse(workflowIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID format"})
		return
	}

	if !workflow.DeleteWorkflow(userID.(string), workflowID) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found or access denied"})
		return
	}
	appScheduler.ReloadAndRescheduleWorkflows() // <--- NUEVA LÍNEA: Notificar al scheduler
	c.Status(http.StatusNoContent)
}