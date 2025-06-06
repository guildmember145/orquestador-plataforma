import axios from 'axios';
import { useAuthStore } from '@/stores/auth';

const orchestratorApi = axios.create({
    baseURL: 'http://localhost:9091/api/tasks/v1',
    headers: {
        'Content-Type': 'application/json',
    },
});

// Interceptor para añadir el token de autorización
orchestratorApi.interceptors.request.use(
    (config) => {
        const authStore = useAuthStore();

        // --- INICIO DE LAS LÍNEAS DE DEPURACIÓN ---
        console.log(
            'Interceptor de OrchestratorAPI: Revisando token. Valor actual:', 
            authStore.accessToken
        );
        // --- FIN DE LAS LÍNEAS DE DEPURACIÓN ---

        if (authStore.accessToken && config.headers) {
            console.log('Interceptor de OrchestratorAPI: Token encontrado. Adjuntando cabecera...');
            config.headers['Authorization'] = `Bearer ${authStore.accessToken}`;
        } else {
            console.warn('Interceptor de OrchestratorAPI: No se encontró accessToken en el store.');
        }
        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);

// Interceptor de respuesta (sin cambios)
orchestratorApi.interceptors.response.use(
    (response) => response,
    async (error) => {
        const authStore = useAuthStore();
        if (error.response?.status === 401) {
            console.error("Authorization error with orchestrator service. Logging out.");
            authStore.logout();
        }
        return Promise.reject(error);
    }
);

export default orchestratorApi;