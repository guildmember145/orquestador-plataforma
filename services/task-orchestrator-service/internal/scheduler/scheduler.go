// services/task-orchestrator-service/internal/scheduler/scheduler.go
package scheduler

import (
    "log"

    "github.com/robfig/cron/v3"
    // Asegúrate que la ruta al modelo de workflow y su store sea correcta
    "github.com/guildmember145/task-orchestrator-service/internal/workflow"
)

// WorkflowExecutor define una interfaz para la función que ejecutará un workflow.
// Esto nos permitirá inyectar la lógica de ejecución más adelante.
type WorkflowExecutor func(wf workflow.Workflow)

type Scheduler struct {
    cronRunner       *cron.Cron
    workflowStore    WorkflowStore // Interfaz para acceder a los workflows
    executeWorkflow  WorkflowExecutor   // Función para ejecutar el workflow
}

// WorkflowStore es una interfaz para desacoplar el scheduler del store concreto.
// Por ahora, nuestro store en memoria en model.go puede implementar esto.
type WorkflowStore interface {
    GetAllEnabledScheduledWorkflows() ([]*workflow.Workflow, error)
}

// New crea una nueva instancia del Scheduler.
func New(store WorkflowStore, executor WorkflowExecutor) *Scheduler {
    // Usamos cron.WithSeconds() si queremos precisión de segundos, opcional.
    // Por defecto, la precisión es a nivel de minuto.
    c := cron.New(cron.WithChain(
        cron.SkipIfStillRunning(cron.DefaultLogger), // Evita que un job se ejecute si la instancia anterior aún corre
        cron.Recover(cron.DefaultLogger),           // Recupera de panics en los jobs
    ))

    return &Scheduler{
        cronRunner:      c,
        workflowStore:   store,
        executeWorkflow: executor,
    }
}

// Start inicia el planificador de cron y carga los workflows.
func (s *Scheduler) Start() {
    log.Println("Scheduler starting...")
    s.loadAndScheduleWorkflows() // Carga inicial
    s.cronRunner.Start()
    log.Println("Scheduler started and cron jobs running.")
}

// Stop detiene el planificador de cron.
func (s *Scheduler) Stop() {
    log.Println("Scheduler stopping...")
    s.cronRunner.Stop()
    log.Println("Scheduler stopped.")
}

// loadAndScheduleWorkflows carga todos los workflows habilitados y de tipo schedule
// y los añade al planificador cron.
// NOTA: En un sistema más complejo, querrías poder actualizar esto dinámicamente
// sin reiniciar, o al menos recargar periódicamente.
func (s *Scheduler) loadAndScheduleWorkflows() {
    log.Println("Loading and scheduling workflows...")

    // Limpiar jobs existentes antes de recargar (si se llama múltiples veces)
    for _, entry := range s.cronRunner.Entries() {
        s.cronRunner.Remove(entry.ID)
    }

    workflowsToSchedule, err := s.workflowStore.GetAllEnabledScheduledWorkflows()
    if err != nil {
        log.Printf("Error loading workflows for scheduler: %v", err)
        return
    }

    scheduledCount := 0
    for _, wf := range workflowsToSchedule {
        // Copiamos wf a una variable local para evitar problemas de clausura en el goroutine del job
        currentWorkflow := *wf 

        cronSpec, ok := currentWorkflow.Trigger.Config["cron"].(string)
        if !ok || cronSpec == "" {
            log.Printf("Workflow ID %s (%s) has schedule trigger but missing or invalid cron spec in config.", currentWorkflow.ID, currentWorkflow.Name)
            continue
        }

        _, err := s.cronRunner.AddFunc(cronSpec, func() {
            log.Printf("SCHEDULER: Triggering workflow ID %s Name: %s", currentWorkflow.ID, currentWorkflow.Name)
            // Aquí llamamos a la función que realmente ejecuta las acciones del workflow
            s.executeWorkflow(currentWorkflow)
        })

        if err != nil {
            log.Printf("Error adding cron job for workflow ID %s (%s) with spec '%s': %v", currentWorkflow.ID, currentWorkflow.Name, cronSpec, err)
        } else {
            log.Printf("Scheduled workflow ID %s (%s) with spec: %s", currentWorkflow.ID, currentWorkflow.Name, cronSpec)
            scheduledCount++
        }
    }
    log.Printf("%d workflows scheduled.", scheduledCount)
}

// ReloadAndRescheduleWorkflows permite recargar y replanificar.
// Podrías llamar a esto después de que un workflow se crea/actualiza/elimina vía API.
func (s *Scheduler) ReloadAndRescheduleWorkflows() {
    log.Println("Reloading and rescheduling workflows...")
    s.loadAndScheduleWorkflows()
}