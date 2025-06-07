# Servicio de Autenticación (Auth Service)

Este servicio es responsable de gestionar la autenticación y autorización de usuarios dentro de la aplicación "Orquestador de Tareas".

## Funcionalidades Principales

*   **Registro de Usuarios:** Permite a los nuevos usuarios crear una cuenta.
*   **Inicio de Sesión:** Verifica las credenciales de los usuarios y emite tokens de acceso.
*   **Gestión de Tokens JWT:** Genera JSON Web Tokens (JWT) para la autenticación basada en tokens y valida los tokens de las solicitudes entrantes.
*   **Hashing de Contraseñas:** Utiliza bcrypt para almacenar de forma segura las contraseñas de los usuarios.
*   **Persistencia de Datos:** Almacena la información de los usuarios en la base de datos PostgreSQL.

## Integración

*   Es consumido por el **Frontend** para todas las operaciones relacionadas con el usuario (registro, inicio de sesión).
*   Es utilizado por el **Servicio Orquestador de Tareas** para validar los tokens JWT y asegurar que las solicitudes a sus propios endpoints estén debidamente autenticadas y autorizadas.
*   Se conecta directamente a la **Base de Datos PostgreSQL** compartida para gestionar los datos de los usuarios.

## Stack Tecnológico

*   **Lenguaje:** Go (Golang)
*   **Framework API:** Gin Web Framework
*   **Base de Datos:** PostgreSQL
*   **Autenticación:** JWT, bcrypt
