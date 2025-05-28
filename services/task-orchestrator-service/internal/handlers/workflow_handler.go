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
	"github.com/guildmember145/task-orchestrator-service/internal/workflow"
	"github.com/guildmember145/task-orchestrator-service/pkg/transport"
)

var validate = validator.New()

func CreateWorkflowHandler(c *gin.Context) {
	var req transport.CreateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// Validación Paso 1: Validar la estructura CreateWorkflowRequest
	// (Asegúrate de que CreateWorkflowRequest en workflow_dtos.go NO tenga 'dive' en Trigger o Actions para esta prueba)
	log.Println("Validating CreateWorkflowRequest structure (top level)...")
	if err := validate.Struct(req); err != nil {
		log.Printf("Top-level validation failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Top-level validation failed: " + err.Error()})
		return
	}
	log.Println("Top-level validation passed.")

	// Validación Paso 2: Validar explícitamente req.Trigger (que es workflow.TriggerDefinition)
	log.Println("Validating TriggerDefinition structure...")
	if err := validate.Struct(req.Trigger); err != nil {
		log.Printf("Trigger validation failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Trigger validation failed: " + err.Error()})
		return
	}
	log.Println("Trigger validation passed.")

	// Validación Paso 3: Validar explícitamente cada elemento en req.Actions
	log.Println("Validating ActionDefinition structures...")
	if len(req.Actions) > 0 { // Solo validar si hay acciones
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

	c.JSON(http.StatusCreated, newWorkflow)
}

// GetWorkflowsHandler lista los workflows del usuario autenticado
func GetWorkflowsHandler(c *gin.Context) {
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

// GetWorkflowByIDHandler obtiene un workflow específico
func GetWorkflowByIDHandler(c *gin.Context) {
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
func UpdateWorkflowHandler(c *gin.Context) {
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

	// Validación similar a CreateWorkflowHandler para la petición de actualización
	log.Println("Validating UpdateWorkflowRequest (top level)...")
	if err := validate.Struct(req); err != nil { // Asume que CreateWorkflowRequest es adecuado para la actualización
		log.Printf("Update top-level validation failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Update top-level validation failed: " + err.Error()})
		return
	}
	log.Println("Update top-level validation passed.")

	log.Println("Validating Update TriggerDefinition structure...")
	if err := validate.Struct(req.Trigger); err != nil {
		log.Printf("Update Trigger validation failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Update Trigger validation failed: " + err.Error()})
		return
	}
	log.Println("Update Trigger validation passed.")

	log.Println("Validating Update ActionDefinition structures...")
	if len(req.Actions) > 0 {
		for i, action := range req.Actions {
			if err := validate.Struct(action); err != nil {
				log.Printf("Update Action %d validation failed: %v", i, err)
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Update Action %d validation failed: %s", i, err.Error())})
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

	if err := workflow.SaveWorkflow(existingWorkflow); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update workflow: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, existingWorkflow)
}

// DeleteWorkflowHandler elimina un workflow
func DeleteWorkflowHandler(c *gin.Context) {
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
	c.Status(http.StatusNoContent)
}