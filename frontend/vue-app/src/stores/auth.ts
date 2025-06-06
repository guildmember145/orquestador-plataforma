// frontend/vue-app/src/stores/auth.ts

import { defineStore } from 'pinia';
import api from '@/services/api';
import { useWorkflowStore } from '@/stores';
import type { LoginRequest, RegisterRequest, UserProfile } from '@/types/auth';

export interface AuthState {
  accessToken: string | null;
  user: UserProfile | null;
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
    isAuthenticated: (state): boolean => !!state.accessToken,
    currentUser: (state): UserProfile | null => state.user,
  },
  actions: {
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
      const workflowStore = useWorkflowStore();
      workflowStore.clearWorkflows();
    },

    // --- INICIO DE LA SECCIÓN A VERIFICAR ---
    async register(payload: RegisterRequest) { // <-- La firma ahora acepta 1 argumento 'payload'
      this.isLoading = true;
      this.error = null;
      try {
        const response = await api.post('/auth/register', payload);
        console.log('Registration successful:', response.data);
        return response.data;
      } catch (err: any) {
        this.error = err.response?.data?.error || 'Registration failed';
        throw new Error(this.error);
      } finally {
        this.isLoading = false;
      }
    },
    // --- FIN DE LA SECCIÓN A VERIFICAR ---

    async login(payload: LoginRequest) {
      this.isLoading = true;
      this.error = null;
      try {
        const response = await api.post('/auth/login', payload);
        const { access_token } = response.data;
        this.setTokens(access_token);
        api.defaults.headers.common['Authorization'] = `Bearer ${access_token}`;
        await this.fetchUser();
      } catch (err: any) {
        this.clearAuthData();
        this.error = err.response?.data?.error || 'Login failed';
        throw new Error(this.error || 'An unknown login error occurred');
      } finally {
        this.isLoading = false;
      }
    },
    async fetchUser() {
      if (!this.accessToken) return;
      this.isLoading = true;
      this.error = null;
      try {
        const response = await api.get('/users/me');
        this.setUser(response.data);
      } catch (err: any) {
        this.error = err.response?.data?.error || 'Failed to fetch user';
        this.clearAuthData();
        throw err;
      } finally {
        this.isLoading = false;
      }
    },
    logout() {
      this.clearAuthData();
    },
    init() {
      if (this.accessToken && !this.user) {
        this.fetchUser().catch(() => {
          this.logout();
        });
      }
    }
  },
});