# 🚀 Orquestador de Tareas (Task Orchestrator)


Un poderoso orquestador de tareas de código abierto, auto-alojado y construido desde cero con una arquitectura moderna de microservicios. Este proyecto te permite automatizar tareas digitales conectando "disparadores" (como una programación horaria) con "acciones" (como registrar un mensaje o llamar a una API externa).

---

## ✨ Características Principales

* **Arquitectura de Microservicios:** Backend desacoplado en servicios independientes (`auth-service`, `task-orchestrator-service`) para mayor escalabilidad y mantenibilidad.
* **Autenticación Segura:** Sistema completo de registro e inicio de sesión con tokens JWT, hashing de contraseñas (bcrypt) y rutas protegidas.
* **Persistencia de Datos:** Toda la información de usuarios y workflows se almacena de forma persistente en una base de datos PostgreSQL.
* **Motor de Tareas Dinámico:**
    * Crea workflows complejos desde una interfaz de usuario amigable.
    * **Disparador Programado (Scheduler):** Ejecuta tareas en base a expresiones cron.
    * **Motor de Acciones Extensible:** Actualmente soporta acciones de `log_message` y `http_endpoint` (con métodos GET/POST y envío de datos).
* **Interfaz Reactiva y Moderna:** Frontend construido con Vue 3 (Composition API), TypeScript y Pinia para una gestión de estado predecible.
* **Contenerización Completa:** Toda la aplicación (frontend, backends, base de datos) corre en contenedores Podman aislados y conectados en red.
* **Diseño Responsivo:** La interfaz se adapta para una correcta visualización en dispositivos móviles, tablets y escritorio.

---

## 📖 Descripción Detallada de la Aplicación

El "Orquestador de Tareas" es un sistema integral diseñado para la automatización de procesos digitales. Su arquitectura de microservicios se compone de los siguientes elementos principales:

*   **Frontend (Vue.js):** Es la interfaz de usuario web que permite a los usuarios registrarse, iniciar sesión y gestionar visualmente los flujos de trabajo. Desde aquí se configuran los disparadores (eventos que inician tareas, como una programación horaria) y las acciones (operaciones a ejecutar, como realizar una llamada HTTP o registrar un mensaje). Está construido con Vue 3, TypeScript y Pinia.

*   **Servicio de Autenticación (Go):** Este microservicio gestiona todo lo relativo a la identidad y seguridad de los usuarios. Se encarga del registro, inicio de sesión, generación y validación de tokens JWT. Utiliza bcrypt para el hashing seguro de contraseñas y PostgreSQL para persistir los datos de los usuarios.

*   **Servicio Orquestador de Tareas (Go):** Es el cerebro de la aplicación. Define, programa y ejecuta los flujos de trabajo automatizados. Incluye un motor de disparadores (actualmente con soporte para tareas programadas vía cron) y un motor de acciones extensible (actualmente con `log_message` y `http_endpoint`). También utiliza PostgreSQL para almacenar las definiciones de los workflows y el estado de las tareas.

*   **Base de Datos (PostgreSQL):** Actúa como el almacén de datos centralizado y persistente para toda la información crítica, incluyendo usuarios, workflows, y logs de tareas.

La comunicación entre el frontend y los servicios de backend se realiza a través de APIs REST. El `task-orchestrator-service` también se comunica con el `auth-service` para validar la autenticación de las solicitudes. Todo el sistema está diseñado para ser desplegado mediante contenedores (Podman/Docker), facilitando su instalación y gestión.

---

## 🛠️ Stack Tecnológico y Arquitectura

Este proyecto fue construido utilizando un stack tecnológico moderno y robusto, enfocado en el rendimiento y las buenas prácticas.

| Capa | Tecnología | Propósito |
| :--- | :--- | :--- |
| **Frontend** | Vue 3 (Composition API), TypeScript, Vite, Pinia, Vue Router | Interfaz de usuario reactiva, moderna y con tipado seguro. |
| **Backend** | Go (Golang) | Alto rendimiento, concurrencia y binarios eficientes. |
| **API Framework**| Gin Web Framework | Enrutamiento y manejo de peticiones HTTP rápido y ligero. |
| **Base de Datos**| PostgreSQL | Almacenamiento de datos relacional, robusto y persistente. |
| **Driver de BD**| pgx | Driver de Go moderno y de alto rendimiento para PostgreSQL. |
| **Contenerización**| Podman / Docker | Aislamiento, despliegue y gestión de los servicios. |
| **Seguridad** | JWT, bcrypt, CORS | Autenticación, hashing de contraseñas y políticas de origen cruzado. |

### Diagrama de Arquitectura Simplificado
```
[Navegador del Usuario] -> [Frontend VUE (Contenedor, puerto 3003)]
        |                                     |
        | (Llama a)                           | (Llama a)
        V                                     V
[Auth Service (Go, Contenedor, puerto 5000)]  [Task Orchestrator (Go, Contenedor, puerto 9091)]
        |                                     |
        | (Valida token con) <----------------+
        |
        +-------------------------------------> [PostgreSQL (Contenedor, puerto 5432)]
```

---

## 🏁 Instalación y Ejecución

Para levantar el proyecto completo en tu máquina local, sigue estos pasos:

### Prerrequisitos
* Git
* Podman (o Docker)

### Pasos

1.  **Clona el repositorio:**
    ```bash
    git clone [https://github.com/guildmember145/orquestador-plataforma.git](https://github.com/tu-usuario/orquestador-plataforma.git)
    cd orquestador-plataforma
    ```

2.  **Configura las variables de entorno:**
    * Habrá un archivo `.env.example` en cada servicio (`auth-service` y `task-orchestrator-service`). Cópialos a un archivo `.env` en sus respectivos directorios y ajústalos si es necesario (especialmente las claves secretas de JWT).

3.  **Crea la red y el volumen de Podman:**
    ```bash
    podman network create orquestador-net
    podman volume create postgres_data
    ```

4.  **Construye las imágenes de los servicios:**
    ```bash
    # Construir auth-service
    podman build -t auth-service:latest -f services/auth-service/Containerfile ./services/auth-service

    # Construir task-orchestrator-service
    podman build -t task-orchestrator-service:latest -f services/task-orchestrator-service/Containerfile ./services/task-orchestrator-service

    # Construir el frontend
    podman build -t vue-frontend-vite:latest -f frontend/vue-app/Containerfile ./frontend/vue-app
    ```

5.  **Levanta los contenedores:**
    Ejecuta los siguientes comandos en orden:
    ```bash
    # 1. Base de Datos
    podman run -d --name postgres_db --network orquestador-net -p 5432:5432 -v postgres_data:/var/lib/postgresql/data -e POSTGRES_USER=miusuario -e POSTGRES_PASSWORD=micontraseñasegura -e POSTGRES_DB=orquestador_db docker.io/postgres:16-alpine

    # 2. Auth Service
    podman run -d --name auth_service --network orquestador-net -p 5000:5000 -e PORT="5000" -e JWT_SECRET_KEY="TU_CLAVE_SECRETA_JWT" -e DATABASE_URL="postgres://miusuario:micontraseñasegura@postgres_db:5432/orquestador_db" auth-service:latest

    # 3. Task Orchestrator Service
    podman run -d --name task_orchestrator_service --network orquestador-net -p 9091:9090 -e PORT="9090" -e AUTH_SERVICE_BASE_URL="http://auth_service:5000/api/baas/v1/auth" -e DATABASE_URL="postgres://miusuario:micontraseñasegura@postgres_db:5432/orquestador_db" task-orchestrator-service:latest

    # 4. Frontend
    podman run -d --name vue_frontend_app_vite --network orquestador-net -p 3003:4173 vue-frontend-vite:latest
    ```

6.  **¡Listo!** Abre tu navegador y ve a `http://localhost:3003`.

---

## 🚀 Próximos Pasos (Roadmap)
Este proyecto tiene un gran potencial para seguir creciendo. Algunas ideas para el futuro incluyen:
- [ ] **Suite de Pruebas Completa:** Añadir pruebas unitarias y de integración con Vitest y Go Test.
- [ ] **Más Tipos de Acciones:** Implementar acciones como "Enviar Email" o "Publicar en Slack/Discord".
- [ ] **Disparadores Basados en Eventos (Webhooks):** Permitir que los workflows se inicien mediante llamadas HTTP externas.
- [ ] **Historial y Logs de Ejecución Detallados:** Mejorar la interfaz y el almacenamiento para un seguimiento exhaustivo de las ejecuciones de workflows.
- [ ] **Gestión de Errores y Reintentos en Acciones:** Configurar políticas de reintento y manejo de fallos para acciones individuales.
- [ ] **Versionado de Workflows:** Implementar un sistema para guardar y revertir a versiones anteriores de los workflows.
- [ ] **Gestión Segura de Secretos:** Incorporar un sistema para almacenar y utilizar de forma segura credenciales o claves API en las acciones.
- [ ] **Documentación de API para Desarrolladores:** Generar y mantener documentación Swagger/OpenAPI para los servicios de backend.
- [ ] **CI/CD:** Configurar un pipeline con GitHub Actions para construir y testear el proyecto automáticamente.
- [ ] **Autenticación Social (OAuth):** Permitir inicio de sesión con Google, GitHub, etc.

---

## 📄 Licencia
Este proyecto está bajo la Licencia MIT.

---

## 👤 Contacto
* **Hector Leonardo Achucarro**
* **LinkedIn:** https://www.linkedin.com/in/hla-/