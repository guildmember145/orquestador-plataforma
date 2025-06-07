# Servicio Orquestador de Tareas (Task Orchestrator Service)

Este servicio es el núcleo funcional de la aplicación "Orquestador de Tareas". Se encarga de la creación, gestión, programación y ejecución de flujos de trabajo automatizados (workflows).

## Funcionalidades Principales

*   **Gestión de Workflows:** Permite a los usuarios (a través del frontend) crear, leer, actualizar y eliminar workflows. Cada workflow consiste en un disparador y una o más acciones.
*   **Motor de Disparadores (Triggers):** Inicia la ejecución de los workflows.
    *   **Disparador Programado (Scheduler):** Actualmente implementado, ejecuta tareas basándose en expresiones cron.
*   **Motor de Acciones:** Ejecuta las operaciones definidas dentro de un workflow.
    *   **Acciones Soportadas Actualmente:**
        *   `log_message`: Registra un mensaje especificado.
        *   `http_endpoint`: Realiza una llamada HTTP (GET/POST) a un endpoint externo, con capacidad de enviar datos.
*   **Persistencia de Datos:** Almacena las definiciones de los workflows, el historial de ejecución de tareas y otros datos relevantes en la base de datos PostgreSQL.

## Integración

*   Es consumido por el **Frontend** para que los usuarios puedan diseñar y gestionar sus workflows.
*   Interactúa con el **Servicio de Autenticación (Auth Service)** para validar los tokens JWT de las solicitudes entrantes, asegurando que solo los usuarios autenticados puedan acceder y modificar sus workflows.
*   Se conecta directamente a la **Base de Datos PostgreSQL** compartida para persistir y recuperar la información de los workflows y tareas.

## Stack Tecnológico

*   **Lenguaje:** Go (Golang)
*   **Framework API:** Gin Web Framework
*   **Base de Datos:** PostgreSQL
