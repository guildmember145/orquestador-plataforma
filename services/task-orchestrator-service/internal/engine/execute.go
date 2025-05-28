// services/task-orchestrator-service/internal/engine/executor.go
package engine

import (
    "log"
    "github.com/guildmember145/task-orchestrator-service/internal/workflow" // Ajusta el path
)

// ExecuteWorkflow es la función que se llamará cuando un workflow se dispare.
// Por ahora, es un placeholder.
func ExecuteWorkflow(wf workflow.Workflow) {
    log.Printf("ENGINE: Executing workflow ID %s, Name: '%s'", wf.ID, wf.Name)
    log.Printf("ENGINE: Trigger type: %s, Config: %+v", wf.Trigger.Type, wf.Trigger.Config)
    log.Println("ENGINE: Actions to perform:")
    for i, action := range wf.Actions {
        log.Printf("ENGINE:   Action %d: Name: '%s', Type: '%s', Config: %+v, DependsOn: %v",
            i+1, action.Name, action.Type, action.Config, action.DependsOn)
    }
    // TODO: Aquí irá la lógica para interpretar y ejecutar cada ActionDefinition.
    // Por ejemplo, si action.Type == workflow.ActionTypeLogMessage, imprimir action.Config["message"].
    // Si action.Type == workflow.ActionTypeHTTPEndpoint, hacer la llamada HTTP.

    // Simular la ejecución de la primera acción si es log_message
    if len(wf.Actions) > 0 {
        firstAction := wf.Actions[0]
        if firstAction.Type == workflow.ActionTypeLogMessage {
            if msg, ok := firstAction.Config["message"].(string); ok {
                log.Printf("ENGINE: EXECUTED ACTION '%s': %s", firstAction.Name, msg)
            }
        }
    }
    log.Printf("ENGINE: Finished (simulated) execution for workflow ID %s", wf.ID)
}