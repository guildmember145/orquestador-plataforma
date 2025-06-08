<template>
  <div class="workflow-details-view">
    <header class="details-header">
      <div>
        <router-link to="/dashboard/workflows" class="back-link">&larr; Volver a Mis Workflows</router-link>
        <h1>Historial de Ejecuciones</h1>
        <p v-if="workflow" class="workflow-name">Mostrando ejecuciones para: <strong>{{ workflow.name }}</strong></p>
      </div>
      <button @click="refreshExecutions" class="primary-action" :disabled="workflowStore.isWorkflowsLoading">
        {{ workflowStore.isWorkflowsLoading ? 'Cargando...' : 'Refrescar' }}
      </button>
    </header>

    <div v-if="workflowStore.isWorkflowsLoading && executions.length === 0" class="loading">Cargando historial...</div>
    <div v-else-if="workflowStore.getWorkflowError" class="error-message">{{ workflowStore.getWorkflowError }}</div>
    <div v-else-if="executions.length === 0" class="no-workflows">Este workflow no tiene ejecuciones registradas todavía.</div>

    <div v-else class="executions-list">
      <table>
        <thead>
          <tr>
            <th>Estado</th>
            <th>Iniciado</th>
            <th>Finalizado</th>
            <th>Duración</th>
            <th>Logs</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="exec in executions" :key="exec.id">
            <td data-label="Estado"><span :class="['status-badge', exec.status]">{{ exec.status }}</span></td>
            <td data-label="Iniciado">{{ formatarFecha(exec.triggered_at) }}</td>
            <td data-label="Finalizado">{{ exec.completed_at ? formatarFecha(exec.completed_at) : '-' }}</td>
            <td data-label="Duración">{{ calcularDuracion(exec.triggered_at, exec.completed_at) }}</td>
            <td data-label="Logs"><button class="action-btn" @click="showLogs(exec)">Ver Logs</button></td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="selectedExecution" class="modal-overlay" @click="selectedExecution = null">
        <div class="modal-content" @click.stop>
            <h3>Logs para Ejecución <span class="exec-id">{{ selectedExecution.id.substring(0,8) }}...</span></h3>
            <pre class="logs-box">{{ formatLogs(selectedExecution.logs) }}</pre>
            <button @click="selectedExecution = null" class="primary-action">Cerrar</button>
        </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { onMounted, computed, ref } from 'vue';
import { useRoute } from 'vue-router';
import { useWorkflowStore } from '@/stores';
import type { ExecutionLog } from '@/types';

const route = useRoute();
const workflowStore = useWorkflowStore();
const workflowId = route.params.id as string;

const executions = computed(() => workflowStore.getExecutions);
const workflow = computed(() => workflowStore.allWorkflows.find(wf => wf.id === workflowId));
const selectedExecution = ref<ExecutionLog | null>(null);

const refreshExecutions = () => {
  workflowStore.fetchExecutionsForWorkflow(workflowId);
};

onMounted(() => {
  if (workflowStore.allWorkflows.length === 0) {
    workflowStore.fetchWorkflows();
  }
  refreshExecutions();
});

const showLogs = (exec: ExecutionLog) => { selectedExecution.value = exec; };
const formatarFecha = (dateString: string) => new Date(dateString).toLocaleString();
const calcularDuracion = (start: string, end: string | null | undefined) => {
    if (!end) return 'En curso...';
    const duration = new Date(end).getTime() - new Date(start).getTime();
    return `${(duration / 1000).toFixed(2)}s`;
};
const formatLogs = (logs: any) => {
    try {
        if (!logs) return "No hay logs para esta ejecución.";
        const parsedLogs = typeof logs === 'string' ? JSON.parse(logs) : logs;
        if (!Array.isArray(parsedLogs) || parsedLogs.length === 0) return "No hay logs detallados.";
        return parsedLogs.map((log: any) => `[<span class="math-inline">\{new Date\(log\.timestamp\)\.toLocaleTimeString\(\)\}\] \[</span>{log.status}] ${log.message}`).join('\n');
    } catch (e) {
        return "No se pudieron interpretar los logs.";
    }
};
</script>

<style scoped>
  .details-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
  .back-link { color: var(--color-accent); text-decoration: none; margin-bottom: 0.5rem; display: inline-block; }
  .workflow-name { color: var(--color-text-secondary); margin: 0; }
  .status-badge { padding: 5px 10px; border-radius: 12px; color: white; font-size: 0.85em; font-weight: bold; }
  .status-badge.completed { background-color: var(--color-success); }
  .status-badge.running { background-color: var(--color-accent); }
  .status-badge.failed { background-color: var(--color-error); }
  .modal-overlay { position: fixed; inset: 0; z-index: 100; background-color: rgba(0,0,0,0.7); display: flex; justify-content: center; align-items: center; }
  .modal-content { background-color: var(--color-surface); padding: 25px; border-radius: 8px; width: 90%; max-width: 800px; max-height: 80vh; display: flex; flex-direction: column; }
  .modal-content h3 { margin-top: 0; }
  .modal-content .exec-id { font-family: monospace; color: var(--color-text-secondary); }
  .logs-box { flex-grow: 1; background-color: var(--color-background); padding: 15px; border-radius: 5px; white-space: pre-wrap; word-break: break-all; overflow-y: auto; font-family: 'Courier New', Courier, monospace; margin: 1rem 0; }
  /* Estilos para la tabla y los mensajes de estado (puedes copiarlos de WorkflowsDashboard.vue) */
  .loading, .error-message, .no-workflows { text-align: center; margin-top: 50px; padding: 20px; background-color: var(--color-surface); border-radius: 8px; }
  .executions-list table { width: 100%; border-collapse: collapse; }
  .executions-list th, .executions-list td { padding: 15px; text-align: left; border-bottom: 1px solid var(--color-border); }
</style>