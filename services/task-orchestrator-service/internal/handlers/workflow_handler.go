package handlers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10" // Para validación
    "github.com/google/uuid"
    "github.com/guildmember145/task-orchestrator-service/internal/workflow" // Ajusta tu path
    "github.com/guildmember145/task-orchestrator-service/pkg/transport"    // Ajusta tu path
)

var validate = validator.New() // Instancia del validador

// CreateWorkflowHandler crea un nuevo workflow
func CreateWorkflowHandler(c *gin.Context) {
    var req transport.CreateWorkflowRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
        return
    }

    // Validar la estructura de la petición
    if err := validate.Struct(req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed: " + err.Error()})
        return
    }

    userID, _ := c.Get("userID") // Asumimos que el middleware AuthMiddleware ya lo puso
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

    // Verificar que el workflow existe y pertenece al usuario
    existingWorkflow, found := workflow.GetWorkflowByID(userID.(string), workflowID)
    if !found {
        c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found or access denied"})
        return
    }

    var req transport.CreateWorkflowRequest // Reutilizamos el DTO de creación para la actualización
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
        return
    }
    if err := validate.Struct(req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed: " + err.Error()})
        return
    }

    // Actualizar campos del workflow existente
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