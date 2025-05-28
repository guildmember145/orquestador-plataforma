import { defineStore } from 'pinia'
import api from '../services/api'
import type { LoginRequest, RegisterRequest } from '../types/auth';

interface AuthState {
    accessToken: string | null;
    user: { id: string; username: string; email: string } | null;
    isLoading: boolean;
    error: string | null;
}

export const useAuthStore = defineStore('auth', {
    state: (): AuthState => ({
        accessToken: localStorage.getItem('accessToken') || null,
        user: JSON.parse(localStorage.getItem('user') || 'null'),
        isLoading: false,
        error: null,
    }),
    getters: {
        isAuthenticated: (state) => !!state.accessToken,
        currentUser: (state) => state.user,
    },
    actions: {
        setTokens(accessToken: string) {
            this.accessToken = accessToken;
            localStorage.setItem('accessToken', accessToken);
        },
        setUser(userData: any) { // Deberías tener un tipo específico para el usuario
            this.user = userData;
            localStorage.setItem('user', JSON.stringify(userData));
        },
        clearAuthData() {
            this.accessToken = null;
            this.user = null;
            localStorage.removeItem('accessToken');
            localStorage.removeItem('user');
            // También deberías limpiar la instancia de axios si configuraste un interceptor
            if (api.defaults.headers.common['Authorization']) {
                delete api.defaults.headers.common['Authorization'];
            }
        },
        async register(payload: RegisterRequest) {
            this.isLoading = true;
            this.error = null;
            try {
                // La URL completa, incluyendo el puerto si es necesario,
                // ya que el auth-service corre en localhost:8080
                const response = await api.post('/auth/register', payload);
                // Podrías querer hacer login automático o solo mostrar mensaje de éxito
                console.log('Registration successful:', response.data);
                return response.data;
            } catch (err: any) {
                this.error = err.response?.data?.error || 'Registration failed';
                throw err;
            } finally {
                this.isLoading = false;
            }
        },
        async login(payload: LoginRequest) {
            this.isLoading = true;
            this.error = null;
            try {
                const response = await api.post('/auth/login', payload);
                const { access_token } = response.data;
                this.setTokens(access_token);
                api.defaults.headers.common['Authorization'] = `Bearer ${access_token}`;
                // Después del login, obtener datos del usuario
                await this.fetchUser();
                return response.data;
            } catch (err: any) {
                this.clearAuthData();
                this.error = err.response?.data?.error || 'Login failed';
                throw err;
            } finally {
                this.isLoading = false;
            }
        },
        async fetchUser() {
            if (!this.accessToken) return;
            this.isLoading = true;
            this.error = null;
            try {
                // Asegúrate que el token ya está en los headers de la instancia de api
                if (!api.defaults.headers.common['Authorization'] && this.accessToken) {
                     api.defaults.headers.common['Authorization'] = `Bearer ${this.accessToken}`;
                }
                const response = await api.get('/users/me');
                this.setUser(response.data);
            } catch (err: any) {
                this.error = err.response?.data?.error || 'Failed to fetch user';
                this.clearAuthData(); // Si falla obtener el usuario, es probable que el token sea inválido
                throw err;
            } finally {
                this.isLoading = false;
            }
        },
        logout() {
            // Podrías llamar a un endpoint /auth/logout en el backend si lo implementaste
            // para invalidar el token en el servidor (ej. usando una blacklist)
            this.clearAuthData();
            // Redirigir al login (si usas router)
            // router.push('/login');
        },
        // Acción para inicializar el estado al cargar la app
        async init() {
            if (this.accessToken && !this.user) {
               try {
                   await this.fetchUser();
               } catch (error) {
                   console.error("Failed to init auth state:", error);
                   // El token podría haber expirado o ser inválido
                   this.logout();
               }
            }
        }
    },
});