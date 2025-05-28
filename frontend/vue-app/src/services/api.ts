import axios from 'axios';
import { useAuthStore } from '../stores/auth'; // Ajusta la ruta si tu archivo auth.ts está en src/stores

const api = axios.create({
    // La URL base de tu auth-service.
    // Cuando el frontend corre en el navegador, llama a localhost y el puerto mapeado.
    baseURL: 'http://localhost:5000/api/baas/v1', // Asegúrate que este puerto coincide con cómo expones el auth-service
    headers: {
        'Content-Type': 'application/json',
    },
});

// Interceptor para añadir el token a las peticiones
api.interceptors.request.use(
    (config) => {
        const authStore = useAuthStore(); // Obtener la instancia del store aquí
        if (authStore.accessToken && config.headers) {
            config.headers['Authorization'] = `Bearer ${authStore.accessToken}`;
        }
        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);

// Opcional: Interceptor para manejar errores 401 (Unauthorized) globalmente
api.interceptors.response.use(
    (response) => response,
    async (error) => {
        const originalRequest = error.config;
        const authStore = useAuthStore();

        // Si el error es 401 y no es un reintento, podríamos intentar refrescar el token
        // (Esto requiere implementar lógica de refresh token en backend y store)
        if (error.response?.status === 401 && !originalRequest._retry) {
            originalRequest._retry = true;
            console.log('Token might be expired or invalid. Logging out.');
            authStore.logout(); // Simplemente desloguear por ahora
            // Aquí podrías redirigir al login
            // router.push('/login');
        }
        return Promise.reject(error);
    }
);

export default api;