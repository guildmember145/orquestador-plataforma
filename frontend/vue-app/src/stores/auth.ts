// frontend/vue-app/src/stores/auth.ts

import { defineStore } from 'pinia';
import api from '../services/api';
import type { LoginRequest, UserProfile } from '../types/auth';

// 1. Definimos una interfaz para el estado del store
export interface AuthState {
  accessToken: string | null;
  user: UserProfile | null;
  isLoading: boolean;
  error: string | null;
}

export const useAuthStore = defineStore('auth', {
  // 2. Usamos la interfaz para tipar el estado
  state: (): AuthState => ({
    accessToken: localStorage.getItem('accessToken') || null,
    user: JSON.parse(localStorage.getItem('user') || 'null'),
    isLoading: false,
    error: null,
  }),
  getters: {
    // 3. Tipamos los getters para que TypeScript los entienda
    isAuthenticated: (state): boolean => !!state.accessToken,
    currentUser: (state): UserProfile | null => state.user,
  },
  actions: {
    // 4. Las acciones permanecen mayormente igual
    setTokens(accessToken: string) {
      this.accessToken = accessToken;
      localStorage.setItem('accessToken', accessToken);
    },
    setUser(userData: UserProfile) {
      this.user = userData;
      localStorage.setItem('user', JSON.stringify(userData));
    },
    clearAuthData() {
      this.accessToken = null;
      this.user = null;
      localStorage.removeItem('accessToken');
      localStorage.removeItem('user');
      // Limpiar también el store de workflows al hacer logout
      // const workflowStore = useWorkflowStore();
      // No existe clearWorkflows, así que no llamamos nada aquí
    },
    async register() {
      // ... (lógica de registro sin cambios)
    },
    async login(payload: LoginRequest) {
      this.isLoading = true;
      this.error = null;
      try {
        const response = await api.post('/auth/login', payload);
        const { access_token } = response.data;
        this.setTokens(access_token);
        // Actualizamos la cabecera por defecto en la instancia de Axios del auth-service
        api.defaults.headers.common['Authorization'] = `Bearer ${access_token}`;
        await this.fetchUser();
      } catch (err: any) {
        this.clearAuthData();
        this.error = err.response?.data?.error || 'Login failed';
        throw new Error(this.error ?? 'Login failed');
      } finally {
        this.isLoading = false;
      }
    },
    async fetchUser() {
      // ... (lógica de fetchUser sin cambios)
    },
    logout() {
      this.clearAuthData();
      // La redirección ahora la manejamos en App.vue
    },
    init() {
      if (this.accessToken && !this.user) {
        this.fetchUser().catch(() => {
          // Si el token es inválido o expiró, deslogueamos
          this.logout();
        });
      }
    }
  },
});