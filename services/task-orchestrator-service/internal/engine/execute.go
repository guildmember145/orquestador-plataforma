// services/task-orchestrator-service/internal/engine/executor.go
package engine

import (
	"bytes" // Para manejar el cuerpo de la petición HTTP
	"encoding/json" // Para manejar el cuerpo JSON
	"fmt"
	"io" // Para leer el cuerpo de la respuesta
	"log"
	"net/http" // El paquete estrella para esta acción
	"strings"  // Para trabajar con strings como el método HTTP
	"time"     // Para timeouts
	"errors"
	// Asegúrate de que la ruta a tu paquete workflow sea correcta
	"github.com/guildmember145/task-orchestrator-service/internal/workflow"
)

// HTTPClient con un timeout por defecto. Puedes hacerlo configurable más adelante.
var httpClient = &http.Client{
	Timeout: 30 * time.Second, // Timeout general para las peticiones HTTP
}

// ExecuteWorkflow es la función que se llamará cuando un workflow se dispare.
func ExecuteWorkflow(wf workflow.Workflow) {
	log.Printf("ENGINE: >>> Starting execution for Workflow ID %s, Name: '%s' <<<", wf.ID, wf.Name)

	if len(wf.Actions) == 0 {
		log.Printf("ENGINE: No actions to execute for Workflow ID %s, Name: '%s'", wf.ID, wf.Name)
		log.Printf("ENGINE: >>> Finished execution for Workflow ID %s, Name: '%s' <<<", wf.ID, wf.Name)
		return
	}

	log.Printf("ENGINE: Processing %d action(s)...", len(wf.Actions))

	for i, action := range wf.Actions {
		log.Printf("ENGINE: --- Action %d/%d: Name: '%s', Type: '%s' ---",
			i+1, len(wf.Actions), action.Name, action.Type)

		var actionErr error
		actionSuccessful := true // Asumimos éxito hasta que algo falle

		switch action.Type {
		case workflow.ActionTypeLogMessage:
			msg, ok := action.Config["message"].(string)
			if !ok {
				errMsg := fmt.Sprintf("Action '%s': 'message' not found in config or is not a string. Config: %+v", action.Name, action.Config)
				log.Printf("ENGINE: ERROR - %s", errMsg)
				actionErr = errors.New(errMsg)
				actionSuccessful = false
			} else {
				log.Printf("ACTION_OUTPUT [%s - %s]: %s", wf.Name, action.Name, msg)
			}

		case workflow.ActionTypeHTTPEndpoint:
			// 1. Extraer configuración de la acción HTTP
			url, urlOk := action.Config["url"].(string)
			if !urlOk || url == "" {
				errMsg := fmt.Sprintf("Action '%s': 'url' is required and must be a string for http_endpoint. Config: %+v", action.Name, action.Config)
				log.Printf("ENGINE: ERROR - %s", errMsg)
				actionErr = errors.New(errMsg)
				actionSuccessful = false
				break // Salir del switch para esta acción
			}

			method, methodOk := action.Config["method"].(string)
			if !methodOk || method == "" {
				method = http.MethodGet // Por defecto GET
			}
			method = strings.ToUpper(method)

			var reqBody io.Reader
			if bodyData, bodyExists := action.Config["body"]; bodyExists && bodyData != nil {
				// Si el body es un string, lo usamos directamente.
				// Si es un mapa o slice, lo convertimos a JSON.
				if bodyStr, isStr := bodyData.(string); isStr {
					reqBody = strings.NewReader(bodyStr)
				} else {
					// Asumimos que es algo que se puede serializar a JSON (map, struct)
					jsonBody, err := json.Marshal(bodyData)
					if err != nil {
						errMsg := fmt.Sprintf("Action '%s': failed to marshal 'body' to JSON. Config.body: %+v, Error: %v", action.Name, bodyData, err)
						log.Printf("ENGINE: ERROR - %s", errMsg)
						actionErr = errors.New(errMsg)
						actionSuccessful = false
						break
					}
					reqBody = bytes.NewBuffer(jsonBody)
					log.Printf("ENGINE: INFO - Action '%s': Marshalled JSON body: %s", action.Name, string(jsonBody))
				}
			}

			// Crear la petición HTTP
			req, err := http.NewRequest(method, url, reqBody)
			if err != nil {
				errMsg := fmt.Sprintf("Action '%s': failed to create HTTP request. URL: %s, Method: %s, Error: %v", action.Name, url, method, err)
				log.Printf("ENGINE: ERROR - %s", errMsg)
				actionErr = errors.New(errMsg)
				actionSuccessful = false
				break
			}

			// Añadir cabeceras
			if headersData, headersExist := action.Config["headers"].(map[string]interface{}); headersExist {
				for key, val := range headersData {
					if valStr, isStr := val.(string); isStr {
						req.Header.Set(key, valStr)
					} else {
						log.Printf("ENGINE: WARNING - Action '%s': Header '%s' value is not a string, skipping. Value: %+v", action.Name, key, val)
					}
				}
			}
			// Asegurar Content-Type si hay cuerpo y no se especificó una por el usuario
			if reqBody != nil && req.Header.Get("Content-Type") == "" {
				// Si el cuerpo fue serializado desde JSON, asumimos application/json
				if _, isMapOrSlice := action.Config["body"].(map[string]interface{}); isMapOrSlice {
					req.Header.Set("Content-Type", "application/json")
				} else if _, isStr := action.Config["body"].(string); isStr {
					// Si es un string, podría ser text/plain o el usuario debe especificarlo
					// Por ahora no ponemos un default si es string para no adivinar
				}
			}


			log.Printf("ENGINE: INFO - Action '%s': Sending %s request to %s", action.Name, method, url)
			if reqBody != nil {
				log.Printf("ENGINE: INFO - Action '%s': With Headers: %+v, With Body (type): %T", action.Name, req.Header, action.Config["body"])
			} else {
				log.Printf("ENGINE: INFO - Action '%s': With Headers: %+v", action.Name, req.Header)
			}
			

			// Enviar la petición
			resp, err := httpClient.Do(req)
			if err != nil {
				errMsg := fmt.Sprintf("Action '%s': failed to execute HTTP request to %s. Error: %v", action.Name, url, err)
				log.Printf("ENGINE: ERROR - %s", errMsg)
				actionErr = errors.New(errMsg)
				actionSuccessful = false
				break
			}
			defer resp.Body.Close()

			// Leer el cuerpo de la respuesta (opcionalmente, podríamos limitarlo)
			respBodyBytes, errRead := io.ReadAll(resp.Body)
			if errRead != nil {
				log.Printf("ENGINE: WARNING - Action '%s': failed to read response body from %s. Error: %v", action.Name, url, errRead)
				// No consideramos esto un fallo de la acción en sí, pero sí un problema de lectura.
			}

			log.Printf("ACTION_OUTPUT [%s - %s]: HTTP Status: %s, Response Body (first 500 bytes): %.500s",
				wf.Name, action.Name, resp.Status, string(respBodyBytes))

			// Considerar un éxito si el código de estado es 2xx
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				errMsg := fmt.Sprintf("Action '%s': HTTP request to %s returned non-2xx status: %s", action.Name, url, resp.Status)
				log.Printf("ENGINE: ERROR - %s", errMsg)
				actionErr = errors.New(errMsg) // Consideramos esto un error de la acción
				actionSuccessful = false
			}


		default:
			errMsg := fmt.Sprintf("Action '%s': Unknown action type '%s'", action.Name, action.Type)
			log.Printf("ENGINE: WARNING - %s", errMsg)
			actionErr = errors.New(errMsg)
			actionSuccessful = false
		}

		if !actionSuccessful {
			log.Printf("ENGINE: --- Error processing Action %d: Name: '%s'. Error: %v ---", i+1, action.Name, actionErr)
			// TODO: Decidir si la ejecución del workflow debe detenerse si una acción falla.
			// Por ahora, continuamos con la siguiente acción. O podríamos añadir un 'break' aquí.
		} else {
			log.Printf("ENGINE: --- Successfully processed Action %d: Name: '%s' ---", i+1, action.Name)
		}
	}

	log.Printf("ENGINE: >>> All actions processed for Workflow ID %s, Name: '%s' <<<", wf.ID, wf.Name)
	// TODO: Implementar el logging de ejecución general del workflow (éxito/fallo basado en las acciones).
}