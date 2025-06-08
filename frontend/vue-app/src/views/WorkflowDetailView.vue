<template>
  <div class="workflow-details">
    <header class="details-header">
      <div>
        <router-link to="/dashboard/workflows" class="back-link">&larr; Volver a Workflows</router-link>
        <h1>Historial de Ejecuciones</h1>
        <p v-if="workflow" class="workflow-name">Workflow: {{ workflow.name }}</p>
      </div>
    </header>

    <div v-if="workflowStore.isWorkflowsLoading" class="loading">Cargando historial...</div>
    <div v-else-if="workflowStore.getWorkflowError" class="error-message">
        {{ workflowStore.getWorkflowError }}
    </div>
    <div v-else-if="executions.length === 0" class="no-workflows">
      Este workflow no tiene ejecuciones todavía.
    </div>

    <div v-else class="executions-list">
      </div>

    </div>
</template>

<script setup lang="ts">
import { onMounted, computed, ref } from 'vue';
import { useRoute } from 'vue-router';
// --- CORRECCIÓN DE IMPORTS ---
import { useWorkflowStore } from '@/stores';
import type { ExecutionLog } from '@/types'; // Importamos el tipo desde el archivo central

const route = useRoute();
const workflowStore = useWorkflowStore();

const workflowId = ref(route.params.id as string);

// Obtenemos los datos directamente del store
const executions = computed(() => workflowStore.getExecutions);

// --- CORRECCIÓN DE COMPARACIÓN ---
// Comparamos el `id` del workflow con `workflowId.value` (el valor del ref)
const workflow = computed(() => workflowStore.allWorkflows.find(wf => String(wf.id) === workflowId.value));

const selectedExecution = ref<ExecutionLog | null>(null);

onMounted(() => {
  // --- CORRECCIÓN DE LLAMADA A LA ACCIÓN ---
  // La acción se llama 'fetchExecutionsForWorkflow'
  workflowStore.fetchExecutionsForWorkflow(workflowId.value);
});

const showLogs = (exec: ExecutionLog) => {
    selectedExecution.value = exec;
};

const formatLogs = (logs: any) => {
    try {
        // El backend puede devolver null si no hay logs
        if (!logs) return "No hay logs para esta ejecución.";
        const parsedLogs = typeof logs === 'string' ? JSON.parse(logs) : logs;
        return parsedLogs.map((log: any) => `[<span class="math-inline">\{new Date\(log\.timestamp\)\.toLocaleTimeString\(\)\}\] \[</span>{log.status}] ${log.message}`).join('\n');
    } catch (e) {
        return "No se pudieron parsear los logs.";
    }
};
</script>

<style scoped>
  /* Tus estilos existentes aquí... */
</style>