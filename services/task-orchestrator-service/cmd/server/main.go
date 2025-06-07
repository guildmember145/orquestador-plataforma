// services/task-orchestrator-service/cmd/server/main.go
package main

import (
	"context" // <-- Añadido para los chequeos de la base de datos
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/guildmember145/task-orchestrator-service/internal/engine"
	"github.com/guildmember145/task-orchestrator-service/internal/handlers"
	"github.com/guildmember145/task-orchestrator-service/internal/middleware"
	"github.com/guildmember145/task-orchestrator-service/internal/scheduler"
	"github.com/guildmember145/task-orchestrator-service/internal/workflow"
	"github.com/guildmember145/task-orchestrator-service/pkg/config"
	"github.com/guildmember145/task-orchestrator-service/pkg/database"
)

func main() {
	config.LoadConfig()
	dbPool := database.ConnectDB()
	defer dbPool.Close()
	database.RunMigrations(dbPool)

	// --- INICIO: PUNTO DE CHEQUEO 1 ---
	log.Println("--- POOL CHECK 1 (después de migraciones) ---")
	var one int
	errCheck1 := dbPool.QueryRow(context.Background(), "SELECT 1").Scan(&one)
	if errCheck1 != nil {
		log.Fatalf("Pool check 1 FALLÓ: %v", errCheck1)
	}
	log.Println("Pool check 1 PASÓ con éxito.")
	// --- FIN: PUNTO DE CHEQUEO 1 ---

	workflowStore := workflow.NewPostgresWorkflowStore(dbPool)
	appScheduler := scheduler.New(workflowStore, engine.ExecuteWorkflow)
	workflowHandler := handlers.NewWorkflowHandler(workflowStore, appScheduler)

	// --- INICIO: PUNTO DE CHEQUEO 2 ---
	log.Println("--- POOL CHECK 2 (antes de iniciar scheduler) ---")
	errCheck2 := dbPool.QueryRow(context.Background(), "SELECT 1").Scan(&one)
	if errCheck2 != nil {
		log.Fatalf("Pool check 2 FALLÓ: %v", errCheck2)
	}
	log.Println("Pool check 2 PASÓ con éxito.")
    // --- FIN: PUNTO DE CHEQUEO 2 ---

	go appScheduler.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("Shutdown signal received, stopping scheduler...")
		appScheduler.Stop()
	}()

	router := gin.Default()
	
    corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3003"}
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(corsConfig))

	taskApiRoutes := router.Group("/api/tasks/v1")
	taskApiRoutes.Use(middleware.AuthMiddleware())
	{
		taskApiRoutes.POST("/workflows", workflowHandler.CreateWorkflowHandler)
		taskApiRoutes.GET("/workflows", workflowHandler.GetWorkflowsHandler)
		taskApiRoutes.GET("/workflows/:workflow_id", workflowHandler.GetWorkflowByIDHandler)
		taskApiRoutes.PUT("/workflows/:workflow_id", workflowHandler.UpdateWorkflowHandler)
		taskApiRoutes.DELETE("/workflows/:workflow_id", workflowHandler.DeleteWorkflowHandler)
	}

	addr := fmt.Sprintf(":%s", config.AppConfig.Port)
	log.Printf("Task Orchestrator service starting HTTP server on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to run Task Orchestrator HTTP server: %v", err)
	}
}