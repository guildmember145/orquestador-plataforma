// services/task-orchestrator-service/internal/engine/executor.go
package engine

import (
	"fmt" // Lo usaremos para formatear algunos logs
	"log"

	// Asegúrate de que la ruta a tu paquete workflow sea correcta
	"github.com/guildmember145/task-orchestrator-service/internal/workflow"
)

// ExecuteWorkflow es la función que se llamará cuando un workflow se dispare.
func ExecuteWorkflow(wf workflow.Workflow) {
	log.Printf("ENGINE: >>> Starting execution for Workflow ID %s, Name: '%s' <<<", wf.ID, wf.Name)
	// log.Printf("ENGINE: Trigger type: %s, Config: %+v", wf.Trigger.Type, wf.Trigger.Config) // Log opcional del trigger

	if len(wf.Actions) == 0 {
		log.Printf("ENGINE: No actions to execute for Workflow ID %s, Name: '%s'", wf.ID, wf.Name)
		log.Printf("ENGINE: >>> Finished execution for Workflow ID %s, Name: '%s' <<<", wf.ID, wf.Name)
		return
	}

	log.Printf("ENGINE: Processing %d action(s)...", len(wf.Actions))

	// Por ahora, las acciones se ejecutan secuencialmente en el orden definido.
	// TODO: Considerar el campo 'DependsOn' para flujos de ejecución más complejos (paralelos/condicionales).
	for i, action := range wf.Actions {
		log.Printf("ENGINE: --- Action %d/%d: Name: '%s', Type: '%s' ---",
			i+1, len(wf.Actions), action.Name, action.Type)

		// Variable para registrar el resultado de la acción, si es necesario
		var actionErr error

		switch action.Type {
		case workflow.ActionTypeLogMessage:
			// Extraer el mensaje de la configuración de la acción
			msg, ok := action.Config["message"].(string)
			if !ok {
				errMsg := fmt.Sprintf("Action '%s': 'message' not found in config or is not a string. Config: %+v", action.Name, action.Config)
				log.Printf("ENGINE: ERROR - %s", errMsg)
				actionErr = fmt.Errorf(errMsg) // Guardar el error
			} else {
				// Imprimir el mensaje en el log del task-orchestrator-service
				// Usamos un prefijo distintivo para identificar estos logs de acción.
				log.Printf("ACTION_OUTPUT [%s - %s]: %s", wf.Name, action.Name, msg)
			}

		case workflow.ActionTypeHTTPEndpoint:
			log.Printf("ENGINE: INFO - Action '%s': HTTP Endpoint action - Not yet implemented. Config: %+v", action.Name, action.Config)
			// TODO: Implementar la lógica para hacer llamadas HTTP aquí.
			// actionErr = fmt.Errorf("HTTP Endpoint action not implemented yet")

		default:
			errMsg := fmt.Sprintf("Action '%s': Unknown action type '%s'", action.Name, action.Type)
			log.Printf("ENGINE: WARNING - %s", errMsg)
			actionErr = fmt.Errorf(errMsg) // Guardar el error
		}

		if actionErr != nil {
			log.Printf("ENGINE: --- Error executing Action %d: Name: '%s'. Error: %v ---", i+1, action.Name, actionErr)
			// TODO: Decidir si la ejecución del workflow debe detenerse si una acción falla.
			// Por ahora, continuamos con la siguiente acción.
		} else {
			log.Printf("ENGINE: --- Successfully processed Action %d: Name: '%s' ---", i+1, action.Name)
		}
	}

	log.Printf("ENGINE: >>> All actions processed for Workflow ID %s, Name: '%s' <<<", wf.ID, wf.Name)
	// TODO: Implementar el logging de ejecución general del workflow (éxito/fallo basado en las acciones).
}