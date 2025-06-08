// src/stores/workflowStore.ts
import { defineStore } from 'pinia';
import orchestratorApi from '@/services/orchestratorApi';
import type { Workflow, ExecutionLog } from '@/types'; // <-- Importamos desde nuestro nuevo archivo de tipos

interface WorkflowState {
    workflows: Workflow[];
    executions: ExecutionLog[];
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
        async fetchExecutionsForWorkflow(workflowId: string) {
            this.isLoading = true;
            this.error = null;
            this.executions = [];
            try {
                const response = await orchestratorApi.get<ExecutionLog[]>(`/workflows/${workflowId}/executions`);
                this.executions = response.data || [];
            } catch (err: any) {
                this.error = err.response?.data?.error || 'No se pudo cargar el historial';
            } finally {
                this.isLoading = false;
            }
        },
        async createWorkflow(payload: any) { /* ...código existente... */ },

        async updateWorkflow(workflowId: string, payload: any) {
        this.isLoading = true;
        this.error = null;
        try {
            // 1. Llamamos al endpoint PUT de nuestra API con los nuevos datos
            const response = await orchestratorApi.put<Workflow>(`/workflows/${workflowId}`, payload);

            // 2. Buscamos el índice del workflow en nuestro array local
            const index = this.workflows.findIndex(wf => wf.id === workflowId);
            if (index !== -1) {
                // 3. Si lo encontramos, lo reemplazamos con los datos actualizados
                //    para que la UI se refresque instantáneamente.
                this.workflows[index] = response.data;
            }
            return response.data;
        } catch (err: any) {
            this.error = err.response?.data?.error || `No se pudo actualizar el workflow`;
            throw err;
        } finally {
            this.isLoading = false;
        }
    },
    // --- FIN DE LA NUEVA LÓGICA ---

    async deleteWorkflow(workflowId: string) {
        this.isLoading = true;
        this.error = null;
        try {
            await orchestratorApi.delete(`/workflows/${workflowId}`);
            this.workflows = this.workflows.filter(wf => wf.id !== workflowId);
        } catch (err: any) {
            this.error = `No se pudo eliminar el workflow`;
            throw err;
        } finally {
            this.isLoading = false;
        }
    },



        clearWorkflows() {
            this.workflows = [];
            this.executions = [];
            this.isLoading = false;
            this.error = null;
        },
        // ...resto de acciones como updateWorkflow
    },
});