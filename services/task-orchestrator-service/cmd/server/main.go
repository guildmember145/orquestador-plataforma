// services/task-orchestrator-service/cmd/server/main.go
package main

import (
    "fmt"
    "log"
    "net/http" // Necesario para c.AbortWithStatus en CORS
    "os"       // Para señales de interrupción
    "os/signal"// Para señales de interrupción
    "syscall"  // Para señales de interrupción

    "github.com/gin-gonic/gin"
    "github.com/guildmember145/task-orchestrator-service/internal/engine"    // Ajusta tu path
    "github.com/guildmember145/task-orchestrator-service/internal/handlers"   // Ajusta tu path
    "github.com/guildmember145/task-orchestrator-service/internal/middleware" // Ajusta tu path
    "github.com/guildmember145/task-orchestrator-service/internal/scheduler"  // Ajusta tu path
    "github.com/guildmember145/task-orchestrator-service/internal/workflow"   // Ajusta tu path
    "github.com/guildmember145/task-orchestrator-service/pkg/config"          // Ajusta tu path
)

func main() {
    config.LoadConfig()
    router := gin.Default()

    // --- INICIO: Configuración del Scheduler ---
    // Crear una instancia del workflow store en memoria
    workflowStoreInstance := &workflow.InMemoryWorkflowStore{}

    // Crear el scheduler, inyectando el store y la función placeholder de ejecución
    appScheduler := scheduler.New(workflowStoreInstance, engine.ExecuteWorkflow)

    // Arrancar el scheduler en una goroutine para que no bloquee el servidor HTTP
    go appScheduler.Start()

    // Configurar un canal para escuchar señales de interrupción (Ctrl+C)
    // para detener el scheduler elegantemente.
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-quit
        log.Println("Shutdown signal received, stopping scheduler...")
        appScheduler.Stop()
        // Podrías añadir os.Exit(0) aquí si quieres que la app termine
    }()
    // --- FIN: Configuración del Scheduler ---


    // CORS Middleware (igual que antes)
    router.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent) // Usar http.StatusNoContent para OPTIONS
            return
        }
        c.Next()
    })

    taskApiRoutes := router.Group("/api/tasks/v1")
    // Inyectar el scheduler en los handlers que lo necesiten (ej. para recargar)
    // Por ahora, CreateWorkflowHandler no lo usa, pero podría para llamar a ReloadAndRescheduleWorkflows
    // Para simplificar, la recarga se puede manejar manualmente o en otro momento.
    taskApiRoutes.Use(middleware.AuthMiddleware()) 
    {
        // Pasamos una referencia al scheduler para que los handlers puedan interactuar con él si es necesario
        // Por ejemplo, para llamar a ReloadAndRescheduleWorkflows después de crear/actualizar un workflow.
        // Esto requiere modificar los handlers. Por ahora lo mantenemos simple.
        taskApiRoutes.POST("/workflows", handlers.CreateWorkflowHandler) // Modificar para que llame a appScheduler.ReloadAndRescheduleWorkflows()
        taskApiRoutes.GET("/workflows", handlers.GetWorkflowsHandler)
        taskApiRoutes.GET("/workflows/:workflow_id", handlers.GetWorkflowByIDHandler)
        taskApiRoutes.PUT("/workflows/:workflow_id", handlers.UpdateWorkflowHandler) // Modificar para que llame a appScheduler.ReloadAndRescheduleWorkflows()
        taskApiRoutes.DELETE("/workflows/:workflow_id", handlers.DeleteWorkflowHandler) // Modificar para que llame a appScheduler.ReloadAndRescheduleWorkflows()
    }

    addr := fmt.Sprintf(":%s", config.AppConfig.Port)
    log.Printf("Task Orchestrator service starting HTTP server on %s", addr)
    if err := router.Run(addr); err != nil {
        log.Fatalf("Failed to run Task Orchestrator HTTP server: %v", err)
    }
}