// src/stores/workflowStore.ts
import { defineStore } from 'pinia';
import orchestratorApi from '@/services/orchestratorApi';
import type { Workflow, ExecutionLog } from '@/types';

interface WorkflowState {
    workflows: Workflow[];
    executions: ExecutionLog[]; // Para el historial del workflow seleccionado
    isLoading: boolean;
    error: string | null;
}

export const useWorkflowStore = defineStore('workflows', {
    state: (): WorkflowState => ({
        workflows: [],
        executions: [],
        isLoading: false,
        error: null,
    }),
    getters: {
        allWorkflows: (state) => state.workflows,
        getExecutions: (state) => state.executions,
        isWorkflowsLoading: (state) => state.isLoading,
        getWorkflowError: (state) => state.error,
    },
    actions: {
        async fetchWorkflows() {
            this.isLoading = true;
            this.error = null;
            try {
                const response = await orchestratorApi.get<Workflow[]>('/workflows');
                this.workflows = response.data || [];
            } catch (err: any) {
                this.error = err.response?.data?.error || 'No se pudieron cargar los workflows';
            } finally {
                this.isLoading = false;
            }
        },

        // Acción para obtener el historial de un workflow específico
        async fetchExecutionsForWorkflow(workflowId: string) {
            this.isLoading = true;
            this.error = null;
            this.executions = []; // Limpiamos el historial anterior antes de cargar el nuevo
            try {
                const response = await orchestratorApi.get<ExecutionLog[]>(`/workflows/${workflowId}/executions`);
                this.executions = response.data || [];
            } catch (err: any) {
                this.error = err.response?.data?.error || 'No se pudo cargar el historial de ejecuciones';
            } finally {
                this.isLoading = false;
            }
        },

        async createWorkflow(payload: any) {
            // ... tu código existente ...
            await orchestratorApi.post<Workflow>('/workflows', payload);
            await this.fetchWorkflows(); // Forzar recarga
        },

        async updateWorkflow(id: string, payload: any) {
            // ... tu código existente ...
            const response = await orchestratorApi.put<Workflow>(`/workflows/${id}`, payload);
            const index = this.workflows.findIndex(wf => wf.id === id);
            if (index !== -1) { this.workflows[index] = response.data; }
        },

        async deleteWorkflow(workflowId: string) {
            // ... tu código existente ...
            await orchestratorApi.delete(`/workflows/${workflowId}`);
            this.workflows = this.workflows.filter(wf => wf.id !== workflowId);
        },

        clearWorkflows() {
            this.workflows = [];
            this.executions = [];
            this.isLoading = false;
            this.error = null;
        },
    },
});