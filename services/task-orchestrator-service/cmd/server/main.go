// services/task-orchestrator-service/cmd/server/main.go
package main

import (
    "fmt" 
    "log" 
    "github.com/gin-gonic/gin" 

    "github.com/guildmember145/task-orchestrator-service/internal/handlers"
    "github.com/guildmember145/task-orchestrator-service/internal/middleware"
    "github.com/guildmember145/task-orchestrator-service/pkg/config"
)

func main() {
    config.LoadConfig()
    router := gin.Default()

    // CORS Middleware (igual que antes)
    router.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    })

    taskApiRoutes := router.Group("/api/tasks/v1")
    taskApiRoutes.Use(middleware.AuthMiddleware()) // Middleware de autenticación para todas estas rutas
    {
        taskApiRoutes.POST("/workflows", handlers.CreateWorkflowHandler)
        taskApiRoutes.GET("/workflows", handlers.GetWorkflowsHandler)
        taskApiRoutes.GET("/workflows/:workflow_id", handlers.GetWorkflowByIDHandler)
        taskApiRoutes.PUT("/workflows/:workflow_id", handlers.UpdateWorkflowHandler)
        taskApiRoutes.DELETE("/workflows/:workflow_id", handlers.DeleteWorkflowHandler)
        // Podrías añadir más rutas para logs, ejecutar manualmente, etc.
    }

    addr := fmt.Sprintf(":%s", config.AppConfig.Port)
    log.Printf("Task Orchestrator service starting on %s", addr)
    if err := router.Run(addr); err != nil {
        log.Fatalf("Failed to run Task Orchestrator server: %v", err)
    }
}