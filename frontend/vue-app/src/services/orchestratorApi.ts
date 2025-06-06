import axios from 'axios';
import { useAuthStore } from '../stores/auth'; // Usaremos el authStore para obtener el token

const orchestratorApi = axios.create({
    // La URL base de tu task-orchestrator-service, que corre en el puerto 9091 del host
    baseURL: 'http://localhost:9091/api/tasks/v1',
    headers: {
        'Content-Type': 'application/json',
    },
});

// Interceptor para añadir el token de autorización a CADA petición
// que se haga con esta instancia de Axios.
orchestratorApi.interceptors.request.use(
    (config) => {
        const authStore = useAuthStore();
        if (authStore.accessToken && config.headers) {
            // Añadimos la cabecera 'Authorization' con el Bearer Token
            config.headers['Authorization'] = `Bearer ${authStore.accessToken}`;
        }
        return config;
    },
    (error) => {
        // Manejar errores de la configuración de la petición
        return Promise.reject(error);
    }
);

// Opcional: Interceptor de respuesta para manejar errores 401 globalmente
// (similar al que tenemos en api.ts para el auth-service).
orchestratorApi.interceptors.response.use(
    (response) => response,
    async (error) => {
        const authStore = useAuthStore();
        if (error.response?.status === 401) {
            // Si el token es inválido o expiró, deslogueamos al usuario.
            console.error("Authorization error with orchestrator service. Logging out.");
            authStore.logout();
            // Aquí podrías redirigir al login si usas vue-router
            // router.push('/login');
        }
        return Promise.reject(error);
    }
);

export default orchestratorApi;