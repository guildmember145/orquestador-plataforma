// frontend/vue-app/src/stores/workflowStore.ts

import { defineStore } from 'pinia';
import orchestratorApi from '@/services/orchestratorApi'; // Nuestra instancia de Axios para el orquestador

// 1. Interfaz que coincide con los datos REALES que envía tu backend Go
export interface Workflow {
    id: string; // El UUID del backend es un string
    user_id: string;
    name: string;
    description: string;
    trigger: any; // 'any' por ahora para flexibilidad
    actions: any[]; // 'any' por ahora para flexibilidad
    is_enabled: boolean;
    created_at: string;
    updated_at: string;
}

// Interfaz para definir la forma de nuestro estado
interface WorkflowState {
    workflows: Workflow[];
    isLoading: boolean;
    error: string | null;
}

// 2. Usamos el 'Options Store' que ya está preparado para la API
export const useWorkflowStore = defineStore('workflows', {
    // STATE: El estado reactivo de nuestro store
    state: (): WorkflowState => ({
        workflows: [],
        isLoading: false,
        error: null,
    }),
    
    // GETTERS: Propiedades computadas para acceder al estado
    getters: {
        allWorkflows: (state) => state.workflows,
        isWorkflowsLoading: (state) => state.isLoading,
        getWorkflowError: (state) => state.error,
    },

    // ACTIONS: Métodos para modificar el estado (ej. llamando a la API)
    actions: {
        // Obtiene los workflows del backend
        async fetchWorkflows() {
            this.isLoading = true;
            this.error = null;
            try {
                const response = await orchestratorApi.get<Workflow[]>('/workflows');
                this.workflows = response.data;
                console.log('Workflows obtenidos de la API:', this.workflows);
            } catch (err: any) {
                this.error = err.response?.data?.error || 'No se pudieron cargar los workflows';
                console.error(this.error);
            } finally {
                this.isLoading = false;
            }
        },

        // Crea un nuevo workflow en el backend
        async createWorkflow(payload: any) {
            this.isLoading = true;
            this.error = null;
            try {
                const response = await orchestratorApi.post<Workflow>('/workflows', payload);
                this.workflows.push(response.data); // Añade el nuevo workflow a la lista local
                console.log('Workflow creado exitosamente:', response.data);
                return response.data;
            } catch (err: any) {
                this.error = err.response?.data?.error || 'No se pudo crear el workflow';
                console.error(this.error);
                throw new Error(this.error); // Lanza el error para que el componente lo pueda atrapar
            } finally {
                this.isLoading = false;
            }
        },

        // Limpia el store (ej. al hacer logout)
        clearWorkflows() {
            this.workflows = [];
            this.isLoading = false;
            this.error = null;
        },

        // TODO: Implementar la lógica para actualizar y eliminar
        async updateWorkflow(workflowId: string, payload: any) {
        this.isLoading = true;
        this.error = null;
        try {
            // 1. Llamamos al endpoint PUT de la API con los nuevos datos
            const response = await orchestratorApi.put<Workflow>(`/workflows/${workflowId}`, payload);

            // 2. Buscamos el índice del workflow actualizado en nuestro array local
            const index = this.workflows.findIndex(wf => wf.id === workflowId);
            if (index !== -1) {
                // 3. Si lo encontramos, lo reemplazamos con los datos actualizados de la respuesta
                this.workflows[index] = response.data;
            }

            console.log(`Workflow ${workflowId} actualizado exitosamente.`);
            return response.data;

        } catch (err: any) {
            this.error = err.response?.data?.error || `No se pudo actualizar el workflow ${workflowId}`;
            console.error(this.error);
            throw new Error(this.error);
        } finally {
            this.isLoading = false;
        }
    },

        async deleteWorkflow(workflowId: string) {
        this.isLoading = true;
        this.error = null;
        try {
            // 1. Llamamos al endpoint DELETE de nuestra API
            await orchestratorApi.delete(`/workflows/${workflowId}`);

            // 2. Si la llamada es exitosa, actualizamos el estado local
            //    eliminando el workflow del array. Esto hace que la UI
            //    se actualice instantáneamente sin necesidad de volver a
            //    pedir toda la lista al servidor.
            this.workflows = this.workflows.filter((wf: Workflow) => wf.id !== workflowId);

            console.log(`Workflow ${workflowId} eliminado exitosamente.`);

        } catch (err: any) {
            this.error = err.response?.data?.error || `No se pudo eliminar el workflow ${workflowId}`;
            console.error(this.error);
            throw new Error(this.error); // Lanza el error para que el componente lo sepa
        } finally {
            this.isLoading = false;
        }
    },
},
 })
