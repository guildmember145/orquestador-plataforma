// services/task-orchestrator-service/cmd/server/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/guildmember145/task-orchestrator-service/internal/engine"
	"github.com/guildmember145/task-orchestrator-service/internal/handlers"
	"github.com/guildmember145/task-orchestrator-service/internal/middleware"
	"github.com/guildmember145/task-orchestrator-service/internal/scheduler"
	"github.com/guildmember145/task-orchestrator-service/internal/workflow"
	"github.com/guildmember145/task-orchestrator-service/pkg/config"
)

func main() {
	config.LoadConfig()
	router := gin.Default()

	workflowStoreInstance := &workflow.InMemoryWorkflowStore{}
	appScheduler := scheduler.New(workflowStoreInstance, engine.ExecuteWorkflow)
	go appScheduler.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("Shutdown signal received, stopping scheduler...")
		appScheduler.Stop()
	}()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	taskApiRoutes := router.Group("/api/tasks/v1")
	taskApiRoutes.Use(middleware.AuthMiddleware())
	{
		// Pasar la instancia de appScheduler a los handlers que lo necesitan
		taskApiRoutes.POST("/workflows", func(c *gin.Context) {
			handlers.CreateWorkflowHandler(c, appScheduler)
		})
		taskApiRoutes.GET("/workflows", handlers.GetWorkflowsHandler) // No necesita el scheduler
		taskApiRoutes.GET("/workflows/:workflow_id", handlers.GetWorkflowByIDHandler) // No necesita el scheduler
		taskApiRoutes.PUT("/workflows/:workflow_id", func(c *gin.Context) {
			handlers.UpdateWorkflowHandler(c, appScheduler)
		})
		taskApiRoutes.DELETE("/workflows/:workflow_id", func(c *gin.Context) {
			handlers.DeleteWorkflowHandler(c, appScheduler)
		})
	}

	addr := fmt.Sprintf(":%s", config.AppConfig.Port)
	log.Printf("Task Orchestrator service starting HTTP server on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to run Task Orchestrator HTTP server: %v", err)
	}
}