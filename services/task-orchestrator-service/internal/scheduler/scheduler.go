// services/task-orchestrator-service/internal/scheduler/scheduler.go
package scheduler

import (
	"log"

	"github.com/robfig/cron/v3"
	// Asegúrate que la ruta al modelo de workflow y su store sea correcta
	"github.com/guildmember145/task-orchestrator-service/internal/workflow"
)

// WorkflowExecutor define una interfaz para la función que ejecutará un workflow.
// La firma ahora espera un workflow.Store, que es la interfaz central.
type WorkflowExecutor func(wf workflow.Workflow, store workflow.Store)

type Scheduler struct {
	cronRunner      *cron.Cron
	// Usa directamente la interfaz del paquete workflow.
	workflowStore   workflow.Store
	executeWorkflow WorkflowExecutor
}

// New crea una nueva instancia del Scheduler.
func New(store workflow.Store, executor WorkflowExecutor) *Scheduler {
	c := cron.New(cron.WithChain(
		cron.SkipIfStillRunning(cron.DefaultLogger),
		cron.Recover(cron.DefaultLogger),
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
	s.loadAndScheduleWorkflows()
	s.cronRunner.Start()
	log.Println("Scheduler started and cron jobs running.")
}

// Stop detiene el planificador de cron.
func (s *Scheduler) Stop() {
	log.Println("Scheduler stopping...")
	s.cronRunner.Stop()
	log.Println("Scheduler stopped.")
}

func (s *Scheduler) loadAndScheduleWorkflows() {
	log.Println("Loading and scheduling workflows...")

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
		// --- INICIO DE LA CORRECCIÓN DEL BUG ---
		// Creamos una variable local para esta iteración del bucle.
		// Esto es CRÍTICO para que la closure (la función del cron) capture
		// el workflow correcto y no el último de la lista.
		workflowToRun := *wf
		// --- FIN DE LA CORRECCIÓN DEL BUG ---

		cronSpec, ok := workflowToRun.Trigger.Config["cron"].(string)
		if !ok || cronSpec == "" {
			log.Printf("Workflow ID %s (%s) has schedule trigger but missing or invalid cron spec in config.", workflowToRun.ID, workflowToRun.Name)
			continue
		}

		_, err := s.cronRunner.AddFunc(cronSpec, func() {
			// ✅ La closure ahora usa `workflowToRun`, que es la variable correcta y segura.
			log.Printf("SCHEDULER: Triggering workflow ID %s Name: %s", workflowToRun.ID, workflowToRun.Name)
			s.executeWorkflow(workflowToRun, s.workflowStore)
		})

		if err != nil {
			log.Printf("Error adding cron job for workflow ID %s (%s) with spec '%s': %v", workflowToRun.ID, workflowToRun.Name, cronSpec, err)
		} else {
			log.Printf("Scheduled workflow ID %s (%s) with spec: %s", workflowToRun.ID, workflowToRun.Name, cronSpec)
			scheduledCount++
		}
	}
	log.Printf("%d workflows scheduled.", scheduledCount)
}

// ReloadAndRescheduleWorkflows permite recargar y replanificar.
func (s *Scheduler) ReloadAndRescheduleWorkflows() {
	log.Println("Reloading and rescheduling workflows...")
	s.loadAndScheduleWorkflows()
}